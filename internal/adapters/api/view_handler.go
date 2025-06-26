package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
)

type ViewHandler struct {
	pasteService *services.PasteService
	fileService  *services.FileService
	log          *slog.Logger
}

// Get content
func (h *ViewHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	// [TODO] Implement
	h.log.Warn("GetContent handler is not implemented")
	fmt.Fprintf(w, "Not implemented yet") // Placeholder for now
}

// Universal content viewer handler
func (h *ViewHandler) ViewContent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger := h.log.With("content_id", id, "remote_addr", r.RemoteAddr)
	logger.Info("Serving content view")

	// Try as paste first
	paste, err := h.pasteService.Get(r.Context(), id)
	if err == nil {
		// Render paste view
		// In real implementation, this would serve the React app
		// with the paste data
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><pre>%s</pre></body></html>", paste.Content)
		logger.Debug("Serving as paste")
		return
	}

	// Try as file
	file, err := h.fileService.GetInfo(r.Context(), id)
	if err == nil {
		// Render file view
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><h1>%s</h1><p>Size: %d bytes</p><a href='/files/%s'>Download</a></body></html>",
			file.OriginalName, file.Size, file.ID)
		logger.Debug("Serving as file")
		return
	}

	logger.Warn("Content not found")
	http.Error(w, "Content not found", http.StatusNotFound)
}
