package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prologic/bitcask"
)

type Comment struct {
	Author  string
	Content string
	Date    time.Time
}

func getComments(db *bitcask.Bitcask, key string) ([]Comment, error) {
	data, err := db.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	var comments []Comment
	err = json.Unmarshal(data, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func saveComment(db *bitcask.Bitcask, key string, comment Comment) error {
	comments, err := getComments(db, key)
	if err != nil && err != bitcask.ErrKeyNotFound {
		return err
	}

	comment.Date = time.Now()
	comments = append(comments, comment)

	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	return db.Put([]byte(key), data)
}

func submitCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	pagePath := r.FormValue("path")
	if pagePath == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	author := strings.TrimSpace(r.FormValue("author"))
	content := strings.TrimSpace(r.FormValue("content"))

	// Check if author and content fields are not empty
	if author == "" || content == "" {
		http.Error(w, "Author and content fields must not be empty", http.StatusBadRequest)
		return
	}

	comment := Comment{
		Author:  r.FormValue("author"),
		Content: r.FormValue("content"),
	}

	err = saveComment(commentsDB, pagePath, comment)
	if err != nil {
		http.Error(w, "Failed to save comment", http.StatusInternalServerError)
		return
	}

	log.Printf("Saved comment: %+v", comment)
	http.Redirect(w, r, pagePath, http.StatusSeeOther)
}
