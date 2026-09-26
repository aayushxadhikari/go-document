package model

import "time"

type Document struct{
	ID int64 `json:"id"`
	Filename string `json:"filename"`
	StoragePath string `json:"-"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}