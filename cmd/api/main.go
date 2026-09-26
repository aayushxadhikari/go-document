package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aayushxadhikari/go-document-ai/internal/database"
	"github.com/aayushxadhikari/go-document-ai/internal/handler"
	"github.com/aayushxadhikari/go-document-ai/internal/repository"
	"github.com/aayushxadhikari/go-document-ai/internal/service"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL")

	mux := http.NewServeMux()

	documentRepo := repository.NewDocumentRepository(pool)
	chunkRepo := repository.NewChunkRepository(pool)
	documentService := service.NewDocumentService(documentRepo, chunkRepo, "uploads")
	documentHandler := handler.NewDocumentHandler(documentService)

	mux.HandleFunc("POST /documents", documentHandler.Upload)

	mux.HandleFunc("GET /health", healthHandler)

	log.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
