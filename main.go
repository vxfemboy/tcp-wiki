package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-git/go-git/v5"
)

const repoURL = "https://git.tcp.direct/S4D/tcp-wiki.git"
const localPath = "./data"

func main() {
	err := cloneRepository(repoURL, localPath)
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

	err := pullRepository(localPath)
	if err != nil {
		log.Printf("Failed to pull repository: %v", err)
	}

	filePath := strings.TrimPrefix(r.URL.Path, "/")
	if filePath == "" {
		filePath = "README.md"
	}

	err = renderPage(w, localPath, filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
	}
}
