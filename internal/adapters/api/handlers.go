package api

import (
	"log/slog"

	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
)

type Handlers struct {
	fileHandler  *FileHandler
	pasteHandler *PasteHandler
	viewHandler  *ViewHandler
	log          *slog.Logger
}

func NewHandlers(fileService *services.FileService, pasteService *services.PasteService, log *slog.Logger) *Handlers {
	return &Handlers{
		fileHandler:  &FileHandler{fileService: fileService, log: log.With("handler", "file")},
		pasteHandler: &PasteHandler{pasteService: pasteService, log: log.With("handler", "paste")},
		viewHandler:  &ViewHandler{pasteService: pasteService, fileService: fileService, log: log.With("handler", "view")},
		log:          log,
	}
}
