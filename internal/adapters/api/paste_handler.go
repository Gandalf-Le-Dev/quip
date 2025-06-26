package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
)

type PasteHandler struct {
	pasteService *services.PasteService
	log          *slog.Logger
}

// Create paste handler
func (h *PasteHandler) CreatePaste(w http.ResponseWriter, r *http.Request) {
	logger := h.log.With("remote_addr", r.RemoteAddr)
	logger.Debug("Attempting to create a new paste")

	var req struct {
		Content  string `json:"content"`
		Language string `json:"language"`
		Title    string `json:"title"`
		TTL      string `json:"ttl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body", "error", err)
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
		logger.Error("Failed to create paste", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"id":       paste.ID,
		"language": paste.Language,
		"raw":      fmt.Sprintf("/api/paste/%s/raw", paste.ID),
		"view":     fmt.Sprintf("/api/view/%s", paste.ID),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Error("Failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logger.Info("Paste created successfully", "paste_id", paste.ID)
}

// Get paste handler
func (h *PasteHandler) GetPaste(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger := h.log.With("paste_id", id, "remote_addr", r.RemoteAddr)
	logger.Debug("Attempting to get paste")

	paste, err := h.pasteService.Get(r.Context(), id)
	if err != nil {
		logger.Warn("Failed to get paste", "error", err)
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "Paste not found", http.StatusNotFound)
		case domain.ErrExpired:
			http.Error(w, "Paste has expired", http.StatusGone)
		default:
			logger.Error("Internal server error while getting paste", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paste)
	if err != nil {
		logger.Error("Failed to encode paste response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	logger.Info("Successfully retrieved paste")
}

// Get raw paste handler
func (h *PasteHandler) GetRawPaste(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger := h.log.With("paste_id", id, "remote_addr", r.RemoteAddr)
	logger.Debug("Attempting to retrieve raw paste")

	content, err := h.pasteService.GetRaw(r.Context(), id)
	if err != nil {
		logger.Warn("Failed to retrieve raw paste", "error", err)
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "Paste not found", http.StatusNotFound)
		case domain.ErrExpired:
			http.Error(w, "Paste has expired", http.StatusGone)
		default:
			logger.Error("Internal server error while retrieving raw paste", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	logger.Info("Successfully retrieved raw paste")
}

// Delete paste handler
func (h *PasteHandler) DeletePaste(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger := h.log.With("paste_id", id, "remote_addr", r.RemoteAddr)
	logger.Debug("Attempting to delete paste")

	err := h.pasteService.Delete(r.Context(), id)
	if err != nil {
		logger.Warn("Failed to delete paste", "error", err)
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "Paste not found", http.StatusNotFound)
		case domain.ErrExpired:
			http.Error(w, "Paste has expired", http.StatusGone)
		default:
			logger.Error("Internal server error while deleting paste", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
	logger.Info("Paste deleted successfully")
}
