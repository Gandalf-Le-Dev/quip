package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
)

type Handlers struct {
	fileService  *services.FileService
	pasteService *services.PasteService
}

func NewHandlers(fileService *services.FileService, pasteService *services.PasteService) *Handlers {
	return &Handlers{
		fileService:  fileService,
		pasteService: pasteService,
	}
}

// File upload handler
func (h *Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB max
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse TTL
	ttl := 24 * time.Hour // default
	if ttlStr := r.FormValue("ttl"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil {
			ttl = parsed
		}
	}

	// Upload file
	uploadedFile, err := h.fileService.Upload(
		r.Context(),
		file,
		header.Filename,
		header.Size,
		header.Header.Get("Content-Type"),
		ttl,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"id":       uploadedFile.ID,
		"filename": uploadedFile.OriginalName,
		"size":     uploadedFile.Size,
		"download": fmt.Sprintf("/files/%s", uploadedFile.ID),
		"view":     fmt.Sprintf("/view/%s", uploadedFile.ID),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// File download handler
func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	reader, file, err := h.fileService.Download(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "File not found", http.StatusNotFound)
		case domain.ErrExpired:
			http.Error(w, "File has expired", http.StatusGone)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer reader.Close()

	// Set headers
	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.OriginalName))

	// Stream file to response
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Failed to stream file", http.StatusInternalServerError)
	}
}

// Create paste handler
func (h *Handlers) CreatePaste(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content  string `json:"content"`
		Language string `json:"language"`
		Title    string `json:"title"`
		TTL      string `json:"ttl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Parse TTL
	ttl := 24 * time.Hour // default
	if req.TTL != "" {
		if parsed, err := time.ParseDuration(req.TTL); err == nil {
			ttl = parsed
		}
	}

	// Create paste
	paste, err := h.pasteService.Create(
		r.Context(),
		req.Content,
		req.Language,
		req.Title,
		ttl,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"id":       paste.ID,
		"language": paste.Language,
		"raw":      fmt.Sprintf("/paste/%s/raw", paste.ID),
		"view":     fmt.Sprintf("/view/%s", paste.ID),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get paste handler
func (h *Handlers) GetPaste(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	paste, err := h.pasteService.Get(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "Paste not found", http.StatusNotFound)
		case domain.ErrExpired:
			http.Error(w, "Paste has expired", http.StatusGone)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paste)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get raw paste handler
func (h *Handlers) GetRawPaste(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	content, err := h.pasteService.GetRaw(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "Paste not found", http.StatusNotFound)
		case domain.ErrExpired:
			http.Error(w, "Paste has expired", http.StatusGone)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Universal content viewer handler
func (h *Handlers) ViewContent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Try as paste first
	paste, err := h.pasteService.Get(r.Context(), id)
	if err == nil {
		// Render paste view
		// In real implementation, this would serve the React app
		// with the paste data
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><pre>%s</pre></body></html>", paste.Content)
		return
	}

	// Try as file
	file, err := h.fileService.GetInfo(r.Context(), id)
	if err == nil {
		// Render file view
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><h1>%s</h1><p>Size: %d bytes</p><a href='/files/%s'>Download</a></body></html>",
			file.OriginalName, file.Size, file.ID)
		return
	}

	http.Error(w, "Content not found", http.StatusNotFound)
}
