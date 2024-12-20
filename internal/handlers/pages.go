package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	viewsDir := wd + "/internal/views"
	paths := []string{
		filepath.Join(viewsDir, "/rootlayout.html"),
		filepath.Join(viewsDir, "/email-verification-success.html"),
	}
	tmpl := template.Must(template.ParseFiles(paths...))
	tmpl.Execute(w, nil)
}
