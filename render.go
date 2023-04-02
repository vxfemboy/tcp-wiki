package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

type Page struct {
	Content template.HTML
}

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

	page := &Page{Content: template.HTML(mdBuf.String())}
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
