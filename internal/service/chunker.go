package service

import (
	"fmt"
	"strings"
)

func ChunkText(text string, chunkSize, overlap int) ([]string, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunk size must be positive")
	}
	if overlap < 0 || overlap >= chunkSize {
		return nil, fmt.Errorf("overlap must be between 0 and chunk size - 1")
	}

	words := strings.Fields(text)
	var chunks []string

	for start := 0; start < len(words); {
		end := start + chunkSize
		if end > len(words) {
			end = len(words)
		}

		chunks = append(chunks, strings.Join(words[start:end], " "))

		if end == len(words) {
			break
		}
		start = end - overlap
	}
	return chunks, nil
}
