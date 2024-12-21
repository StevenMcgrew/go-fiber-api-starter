package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	data := struct{ ShowLogin bool }{ShowLogin: true}
	filenames := []string{"rootlayout", "header"}
	renderAndSendHTML(w, data, filenames)
}

func AboutPage(c *fiber.Ctx) error {
	// Get views directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Server error when getting working directory: %v", err)
	}
	viewsDir := wd + "/internal/views"
	// Paths
	paths := []string{
		filepath.Join(viewsDir, "rootlayout.html"),
		filepath.Join(viewsDir, "header.html"),
	}
	// Parse
	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		log.Fatalf("Server error when parsing files: %v", err)
	}
	if tmpl == nil {
		log.Fatalf("Server error, tmpl is nil: %v", tmpl)
	}

	// Data
	data := struct{ ShowLogin bool }{ShowLogin: true}
	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c.Response().BodyWriter(), data)
}

func renderAndSendHTML(w http.ResponseWriter, data any, filenames []string) {
	// Get views directory
	wd, err := os.Getwd()
	if err != nil {
		sendRenderError(w, err, "Server error when getting working directory")
		return
	}
	viewsDir := wd + "/internal/views"

	// Get template file paths
	paths := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		paths = append(paths, filepath.Join(viewsDir, filename+".html"))
	}

	// Render
	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		sendRenderError(w, err, "Server error while parsing HTML templates")
		return
	}
	if tmpl == nil {
		sendRenderError(w, errors.New(""), "Parsing html templates produced a nil template")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		sendRenderError(w, err, "Server error while rendering HTML templates")
		return
	}
}

func sendRenderError(w http.ResponseWriter, err error, message string) {
	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	w.WriteHeader(http.StatusInternalServerError)
	encodeErr := json.NewEncoder(w).Encode(Response{
		Status:  "error",
		Message: message,
		Data:    err.Error()})
	if encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}
}
