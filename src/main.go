package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-git/go-git/v5"
	"github.com/prologic/bitcask"
)

type Config struct {
	Server struct {
		Port string
	}
	Git struct {
		UseGit    bool
		RepoURL   string
		Branch    string
		LocalPath string
	}
	Database struct {
		Path string
	}
}

var commentsDB *bitcask.Bitcask

func main() {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if config.Git.UseGit {
		err := cloneRepository(config.Git.RepoURL, config.Git.LocalPath)
		if err != nil && err != git.ErrRepositoryAlreadyExists {
			log.Fatalf("Failed to clone repository: %v", err)
		}
	} else {
		if _, err := os.Stat(config.Git.LocalPath); os.IsNotExist(err) {
			os.MkdirAll(config.Git.LocalPath, os.ModePerm)
		}
	}

	var err error

	commentsDB, err = bitcask.Open(config.Database.Path)
	if err != nil {
		log.Fatalf("Failed to open comments database: %v", err)
	}
	defer commentsDB.Close()

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(&config, w, r)
	})
	http.HandleFunc("/submit_comment", func(w http.ResponseWriter, r *http.Request) {
		submitCommentHandler(w, r)
	})

	srv := &http.Server{Addr: config.Server.Port}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe() failed: %v", err)
		}
	}()

	fmt.Println("Server running at http://127.0.0.1:8080")
	fmt.Println("Press Ctrl-C to stop the server")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown failed: %v", err)
	}
	fmt.Println("Server stopped")
}

func handler(config *Config, w http.ResponseWriter, r *http.Request) {
	// For debugging
	log.Printf("Local Path: %q", config.Git.LocalPath)

	if r.URL.Path == "./assets/favicon.ico" {
		return
	}

	if config.Git.UseGit {
		err := pullRepository(config.Git.LocalPath, config.Git.Branch)
		if err != nil {
			log.Printf("Failed to pull repository: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/")
	if filePath == "" {
		filePath = "README.md"
	}
	log.Printf("Rendering file %q from path %q", filePath, r.URL.Path)

	// Set the Content Security Policy
	csp := "default-src 'self'; font-src 'self' data:; frame-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self';"
	w.Header().Set("Content-Security-Policy", csp)

	markdownFiles, err := listMarkdownFiles(config.Git.LocalPath)
	if err != nil {
		log.Printf("Error listing markdown files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = renderPage(w, r, config, filePath, commentsDB, markdownFiles)
	if err != nil {
		log.Printf("Failed to render page: %v", err)
		http.Error(w, "File not found", http.StatusNotFound)
	}
}

func listMarkdownFiles(localPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			relPath, err := filepath.Rel(localPath, path)
			if err != nil {
				return err
			}

			relPath = strings.Replace(relPath, string(os.PathSeparator), "/", -1)
			files = append(files, relPath)
		}
		return nil
	})
	return files, err
}
