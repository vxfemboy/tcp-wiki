package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/prologic/bitcask"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

type Page struct {
	Content  template.HTML
	Comments []Comment
}

func renderPage(w http.ResponseWriter, r *http.Request, localPath, filePath string, commentsDB *bitcask.Bitcask) error {
	content, err := readFileFromRepo(localPath, filePath)
	if err != nil {
		return err
	}

	log.Printf("Read file content: %s", content)

	ext := filepath.Ext(filePath)
	switch ext {
	case ".md":
		renderMarkdown(w, r, content, commentsDB)
	case ".html", ".css":
		renderStatic(w, content, ext)
	default:
		return fmt.Errorf("unsupported file format")
	}
	return nil
}

func renderMarkdown(w http.ResponseWriter, r *http.Request, content []byte, commentsDB *bitcask.Bitcask) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavored Markdown
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
		),
	)

	var mdBuf bytes.Buffer
	err := md.Convert(content, &mdBuf)
	if err != nil {
		http.Error(w, "Error converting Markdown", http.StatusInternalServerError)
		return
	}

	layout, err := ioutil.ReadFile(filepath.Join(localPath, "assets/_layout.html"))
	if err != nil {
		http.Error(w, "Layout not found", http.StatusInternalServerError)
		return
	}

	comments, err := getComments(commentsDB, r.URL.Path)
	if err != nil && err != bitcask.ErrKeyNotFound {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	page := &Page{Content: template.HTML(mdBuf.String()), Comments: comments}
	t, err := template.New("layout").Parse(string(layout))
	if err != nil {
		http.Error(w, "Error parsing layout", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, page)
	if err != nil {
		http.Error(w, "Error rendering layout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
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
