package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/prologic/bitcask"
)

const repoURL = "https://git.tcp.direct/S4D/tcp-wiki.git"
const localPath = "../data"

var commentsDB *bitcask.Bitcask

func main() {
	err := cloneRepository(repoURL, localPath)
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		log.Fatalf("Failed to clone repository: %v", err)
	}

	commentsDB, err = bitcask.Open("../comments.db")
	if err != nil {
		log.Fatalf("Failed to open comments database: %v", err)
	}
	defer commentsDB.Close()

	http.HandleFunc("/", handler)
	http.HandleFunc("/submit_comment", submitCommentHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	//for debugging
	log.Printf("LOCAL PATH: %q", localPath)

	//...

	if r.URL.Path == "../assets/favicon.ico" {
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
	log.Printf("Rendering file %q from path %q", filePath, r.URL.Path)

	err = renderPage(w, r, localPath, filePath, commentsDB)
	if err != nil {
		log.Printf("Comment loading? %q", commentsDB.Path())

		http.Error(w, "File not found", http.StatusNotFound)
	}
}
