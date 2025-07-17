package routes

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateData struct {
	Title     string
	FileName  string
	FileURL   string
	MimeType  string
	MediaType string
}

func WatchHandler(w http.ResponseWriter, r *http.Request) {
	// Simulación: en producción, obtendrás esta info desde Telegram o BD
	fileName := "example.mp4"
	fileURL := "https://example.com/files/example.mp4"
	mimeType := "video/mp4"

	tmpl, err := template.ParseFiles(filepath.Join("templates", "player.html"))
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mediaType := "other"
	if mimeType[:5] == "video" {
		mediaType = "video"
	} else if mimeType[:5] == "audio" {
		mediaType = "audio"
	}

	data := TemplateData{
		Title:     "Stream " + fileName,
		FileName:  fileName,
		FileURL:   fileURL,
		MimeType:  mimeType,
		MediaType: mediaType,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
	}
}
