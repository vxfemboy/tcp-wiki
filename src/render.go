package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/prologic/bitcask"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

type Page struct {
	Content          template.HTML
	Comments         []Comment
	Path             string
	Author           string
	AuthoredDate     *time.Time
	LastModifier     string
	LastModifiedDate *time.Time
  Pages            []string
  UseGit           bool
}

func renderPage(w http.ResponseWriter, r *http.Request, config *Config, filePath string, commentsDB *bitcask.Bitcask, pages []string) error {
    var content []byte
    var err error

    if config.Git.UseGit {
      content, err = readFileFromRepo(config.Git.LocalPath, filePath)
      if err != nil {
        return err
      }
    } else {
      fullPath := filepath.Join(config.Git.LocalPath, filePath)
      content, err = ioutil.ReadFile(fullPath)
      if err != nil {
        return err
      }
    }

    ext := filepath.Ext(filePath)
    switch ext {
    case ".md":
        renderMarkdown(w, r, content, commentsDB, filePath, pages, config)
    case ".html", ".css":
        renderStatic(w, content, ext)
    default:
        return fmt.Errorf("unsupported file format")
    }
    return nil
}

func renderMarkdown(w http.ResponseWriter, r *http.Request, content []byte, commentsDB *bitcask.Bitcask, filePath string, pages []string, config *Config) {

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
      highlighting.NewHighlighting(
			  highlighting.WithStyle("monokai"),
			),
		),
	)

  var author, lastModifier string 
  var authoredDate, lastModifiedDate *time.Time
  var err error

  if config.Git.UseGit {
    var ad, lmd time.Time
    author, ad, lastModifier, lmd, err = getAuthorAndLastModification(config.Git.LocalPath, filePath)
	  if err != nil {
		  http.Error(w, "Error fetching author and last modification date", http.StatusInternalServerError)
		  return
	  }
    authoredDate = &ad 
    lastModifiedDate = &lmd
  }

	var mdBuf bytes.Buffer
	err = md.Convert(content, &mdBuf)
	if err != nil {
		http.Error(w, "Error converting Markdown", http.StatusInternalServerError)
		return
	}


  layout, err := ioutil.ReadFile("assets/_layout.html")
  if err != nil {
    log.Printf("Error reading _layout.html: %v", err)
    http.Error(w, "Layout not found", http.StatusInternalServerError)
    return
  }

	comments, err := getComments(commentsDB, r.URL.Path)
	if err != nil && err != bitcask.ErrKeyNotFound {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	page := &Page{
		Content:          template.HTML(mdBuf.String()),
		Comments:         comments,
		Path:             r.URL.Path,
		Author:           author,
		AuthoredDate:     authoredDate,
		LastModifier:     lastModifier,
		LastModifiedDate: lastModifiedDate,
    Pages:            pages,
    UseGit:           config.Git.UseGit,
	}

	t, err := template.New("layout").Parse(string(layout))
	if err != nil {
		http.Error(w, "Error parsing layout", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, page)
	if err != nil {
		log.Printf("Error executing template: %v", err) 
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
