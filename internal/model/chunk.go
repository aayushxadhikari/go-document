package model

import "time"

type Chunk struct{
	ID int64 `json:"id`
	DocumentID int64 `json:"document_id"`
	ChunkIndex int `json:"chunk_index"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}