package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
)

func renderPage(w http.ResponseWriter, localPath, filePath string) error {
	content, err := readFileFromRepo(localPath, filePath)
	if err != nil {
		return err
	}

	ext := filepath.Ext(filePath)
	switch ext {
	case ".md":
		renderMarkdown(w, content)
	case ".html", ".css":
		renderStatic(w, content, ext)
	default:
		return fmt.Errorf("unsupported file format")
	}
	return nil
}

func renderMarkdown(w http.ResponseWriter, content []byte) {
	md := blackfriday.Run(content)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(md)
}

func renderStatic(w http.ResponseWriter, content []byte, ext string) {
	contentType := ""
	switch ext {
	case ".html":
		contentType = "text/html; charset=utf-8"
	case ".css":
		contentType = "text/css; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}
