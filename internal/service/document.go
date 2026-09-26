package service

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aayushxadhikari/go-document-ai/internal/model"
	"github.com/aayushxadhikari/go-document-ai/internal/repository"
)

type DocumentService struct {
	repo      *repository.DocumentRepository
	chunkRepo *repository.ChunkRepository
	uploadDir string
}

func NewDocumentService(
	repo *repository.DocumentRepository,
	chunkRepo *repository.ChunkRepository,
	uploadDir string,
) *DocumentService {
	return &DocumentService{
		repo:      repo,
		chunkRepo: chunkRepo,
		uploadDir: uploadDir,
	}
}

func (s *DocumentService) Upload(
	ctx context.Context,
	filename string,
	content io.Reader,
) (*model.Document, error) {
	if err := os.MkdirAll(s.uploadDir, 0750); err != nil {
		return nil, fmt.Errorf("create uploads directory:%w", err)
	}
	file, err := os.CreateTemp(s.uploadDir, "document-*")

	if err != nil {
		return nil, fmt.Errorf("create file:%w", err)
	}

	saved := false
	defer func() {
		if !saved {
			file.Close()
			os.Remove(file.Name())
		}
	}()

	if _, err := io.Copy(file, content); err != nil {
		return nil, fmt.Errorf("write file: %w", err)

	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close file:%w", err)
	}

	document, err := s.repo.Create(ctx, filename, file.Name())
	if err != nil {
		return nil, err
	}

	saved = true
	return document, nil
}
