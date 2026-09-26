package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChunkRepository struct {
	pool *pgxpool.Pool
}

func NewChunkRepository(pool *pgxpool.Pool) *ChunkRepository {
	return &ChunkRepository{pool: pool}
}

func (r *ChunkRepository) CreateMany(
	ctx context.Context,
	documentID int64,
	chunks []string,
) error {
	if len(chunks) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin chunk transaction: %w", err)
	}

	defer tx.Rollback(context.Background())

	const query = `
		INSERT INTO chunks (document_id, chunk_index, content)
		VALUES ($1,$2,$3)
	`

	for index, content := range chunks {
		if _, err := tx.Exec(ctx, query, documentID, index, content); err != nil {
			return fmt.Errorf("insert chunk %d:%w", index, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit chunks: %w", err)
	}
	return nil
}
