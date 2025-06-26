package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
)

type PasteHandler struct {
	pasteService *services.PasteService
}

// Create paste handler
func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
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
func (h *PasteHandler) GetPaste(w http.ResponseWriter, r *http.Request) {
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
func (h *PasteHandler) GetRawPaste(w http.ResponseWriter, r *http.Request) {
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

// Delete paste handler
func (h *PasteHandler) DeletePaste(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.pasteService.Delete(r.Context(), id)
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
	w.WriteHeader(http.StatusNoContent)
}