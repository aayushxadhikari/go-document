package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aayushxadhikari/go-document-ai/internal/service"
)

type DocumentHandler struct {
	service *service.DocumentService
}

func NewDocumentHandler(s *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: s}
}

func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	err := r.ParseMultipartForm(2 << 20)
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Upload exceeds 10MB", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Invalid multipart upload", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "A file field named 'file' is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	document, err := h.service.Upload(r.Context(), header.Filename, file)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		http.Error(w, "Could not save document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(document); err != nil {
		log.Printf("Write upload response: %v", err)
	}
}
