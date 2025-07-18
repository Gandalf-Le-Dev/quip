package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
)

type FileHandler struct {
	fileService *services.FileService
	log         *slog.Logger
}

// File upload handler
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	logger := h.log.With("remote_addr", r.RemoteAddr)
	logger.Debug("Attempting to upload a file")

	// Parse multipart form
	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB max
		logger.Warn("File too large", "error", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Warn("Missing file in form", "error", err)
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse TTL
	ttl := 24 * time.Hour // default
	if ttlStr := r.FormValue("ttl"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil {
			ttl = parsed
		} else {
			logger.Debug("Invalid TTL format, using default", "ttl_provided", ttlStr)
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
		logger.Error("Failed to upload file", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]any{
		"id":       uploadedFile.ID,
		"filename": uploadedFile.OriginalName,
		"size":     uploadedFile.Size,
		"download": fmt.Sprintf("/api/file/%s", uploadedFile.ID),
		"view":     fmt.Sprintf("/api/view/%s", uploadedFile.ID),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Error("Failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	logger.Info("File uploaded successfully", "file_id", uploadedFile.ID)
}

// File download handler
func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger := h.log.With("file_id", id, "remote_addr", r.RemoteAddr)
	logger.Debug("Attempting to download a file")

	reader, file, err := h.fileService.Download(r.Context(), id)
	if err != nil {
		logger.Warn("Failed to download file", "error", err)
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "File not found", http.StatusNotFound)
		case domain.ErrExpired:
			http.Error(w, "File has expired", http.StatusGone)
		default:
			logger.Error("Internal server error during file download", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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
		logger.Error("Failed to stream file to response", "error", err)
		http.Error(w, "Failed to stream file", http.StatusInternalServerError)
	}
	logger.Info("File downloaded successfully")
}

// Get file info handler
func (h *FileHandler) GetFileInfo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger := h.log.With("file_id", id, "remote_addr", r.RemoteAddr)
	logger.Debug("Fetching file info")

	file, err := h.fileService.GetInfo(r.Context(), id)
	if err != nil {
		logger.Warn("Failed to get file info", "error", err)
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "File not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
	logger.Debug("File info sent successfully")
}

// Delete file handler
func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger := h.log.With("file_id", id, "remote_addr", r.RemoteAddr)
	logger.Debug("Attempting to delete a file")

	err := h.fileService.Delete(r.Context(), id)
	if err != nil {
		logger.Error("Failed to delete file", "error", err)
		switch err {
		case domain.ErrNotFound:
			http.Error(w, "File not found", http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
	logger.Info("File deleted successfully")
}
