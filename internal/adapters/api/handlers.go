package api

import (
	"github.com/Gandalf-Le-Dev/quip/internal/core/services"
)

type Handlers struct {
	fileHandler *FileHandler
	pasteHandler *PasteHandler
	viewHandler *ViewHandler
}

func NewHandlers(fileService *services.FileService, pasteService *services.PasteService) *Handlers {
	return &Handlers{
		fileHandler: &FileHandler{fileService: fileService},
		pasteHandler: &PasteHandler{pasteService: pasteService},
		viewHandler: &ViewHandler{pasteService: pasteService, fileService: fileService},
	}
}
