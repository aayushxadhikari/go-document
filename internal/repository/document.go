package repository

import (
	"context"
	"fmt"

	"github.com/aayushxadhikari/go-document-ai/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentRepository(pool *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{pool: pool}
}

func (r *DocumentRepository) Create(
	ctx context.Context,
	filename string,
	storagePath string,
) (*model.Document, error) {
	const query = `
	INSERT INTO documents (filename, storage_path)
	VALUES ($1,$2)
	RETURNING id, filename, storage_path, status, created_at
	`

	var document model.Document

	err := r.pool.QueryRow(ctx, query, filename, storagePath).Scan(
		&document.ID,
		&document.Filename,
		&document.StoragePath,
		&document.Status,
		&document.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}
	return &document, nil
}

func (r *DocumentRepository) UpdateStatus(
	ctx context.Context,
	documentID int64,
	status string,
) error {
	const query = `
		UPDATE documents 
		SET status = $1
		WHERE id = $2
	`

	result, err := r.pool.Exec(ctx, query, status, documentID)
	if err != nil {
		return fmt.Errorf("update document status:%w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("document %d not found", documentID)
	}
	return nil
}
