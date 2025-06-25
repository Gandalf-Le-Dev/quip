package web

import (
	"net/http"
)

func NewRouter(handlers *Handlers) http.Handler {
	mux := http.NewServeMux()

	// File routes
	mux.HandleFunc("POST /upload", handlers.UploadFile)
	mux.HandleFunc("GET /files/{id}", handlers.DownloadFile)

	// Paste routes
	mux.HandleFunc("POST /paste", handlers.CreatePaste)
	mux.HandleFunc("GET /paste/{id}", handlers.GetPaste)
	mux.HandleFunc("GET /paste/{id}/raw", handlers.GetRawPaste)

	// Universal viewer
	mux.HandleFunc("GET /view/{id}", handlers.ViewContent)
	mux.HandleFunc("GET /{id}", handlers.ViewContent)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
	})

	// Apply CORS middleware
	return corsMiddleware(mux)
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
