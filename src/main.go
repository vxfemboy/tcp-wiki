package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
  "path/filepath"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/prologic/bitcask"
)

const repoURL = "https://git.tcp.direct/S4D/tcp-wiki.git"
const repoBRANCH = "main"
const localPath = "data"

var commentsDB *bitcask.Bitcask

func main() {
	err := cloneRepository(repoURL, localPath)
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	commentsDB, err = bitcask.Open("comments.db")
	if err != nil {
		log.Fatalf("Failed to open comments database: %v", err)
	}
	defer commentsDB.Close()

  fs := http.FileServer(http.Dir("./assets"))
  http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", handler)
	http.HandleFunc("/submit_comment", submitCommentHandler)

	srv := &http.Server{Addr: ":8080"}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe() failed: %v", err)
		}
	}()

	fmt.Println("Server running at http://127.0.0.1:8080")
	fmt.Println("Press Ctrl-C to stop the server")

	// Wait for interrupt signal to stop the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Shutdown the server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown failed: %v", err)
	}
	fmt.Println("Server stopped")
}

func handler(w http.ResponseWriter, r *http.Request) {
	// For debugging
	log.Printf("LOCAL PATH: %q", localPath)

	if r.URL.Path == "./assets/favicon.ico" {
		return
	}

	err := pullRepository(localPath, repoBRANCH)
	if err != nil {
		log.Printf("Failed to pull repository: %v", err)
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/")
	if filePath == "" {
		filePath = "README.md"
	}
	log.Printf("Rendering file %q from path %q", filePath, r.URL.Path)

	// Set the Content Security Policy
	csp := "default-src 'self'; img-src 'self'; script-src 'self'; style-src 'self';"
	w.Header().Set("Content-Security-Policy", csp)

  markdownFiles, err := listMarkdownFiles(localPath)
  if err != nil {
    log.Printf("Error listing markdown files: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }

  err = renderPage(w, r, localPath, filePath, commentsDB, markdownFiles)
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
            // Ensure the path uses web-friendly slashes
            relPath = strings.Replace(relPath, string(os.PathSeparator), "/", -1)
            files = append(files, relPath)
        }
        return nil
    })
    return files, err
}
