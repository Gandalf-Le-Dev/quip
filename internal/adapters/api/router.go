package api

import (
	"log/slog"
	"net/http"
	"time"
)

func NewRouter(handlers *Handlers) http.Handler {
	mux := http.NewServeMux()

	// File routes
	fileHandler := handlers.fileHandler
	mux.HandleFunc("POST /api/file", fileHandler.UploadFile)
	mux.HandleFunc("GET /api/file/{id}", fileHandler.DownloadFile)
	mux.HandleFunc("GET /api/file/{id}/info", fileHandler.GetFileInfo)
	mux.HandleFunc("DELETE /api/file/{id}", fileHandler.DeleteFile)

	// Paste routes
	pasteHandler := handlers.pasteHandler
	mux.HandleFunc("POST /api/paste", pasteHandler.CreatePaste)
	mux.HandleFunc("GET /api/paste/{id}", pasteHandler.GetPaste)
	mux.HandleFunc("GET /api/paste/{id}/raw", pasteHandler.GetRawPaste)
	mux.HandleFunc("DELETE /api/paste/{id}", pasteHandler.DeletePaste)

	// Universal viewer
	viewerHandler := handlers.viewHandler
	mux.HandleFunc("GET /api/{id}", viewerHandler.GetContent)
	mux.HandleFunc("GET /api/view/{id}", viewerHandler.ViewContent)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK\n"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Apply middlewares
	loggedMux := requestLogger(handlers.log, mux)
	return corsMiddleware(loggedMux)
}

// Manual CORS implementation
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Request logging middleware
func requestLogger(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
			"remote_addr", r.RemoteAddr,
		)
	})
}
