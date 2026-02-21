package storage

import (
	"fmt"
	model_storage_tables "telegram-antispam-bot/internal/models/storage"
)

// GetListBadWords ..
func (s *Storage) GetListBadWords() ([]string, error) {
	var words = make([]model_storage_tables.Word, 0)

	res := s.s.Find(&words)
	if res.Error != nil {
		return nil, fmt.Errorf("s.Find: %w", res.Error)
	}

	result := make([]string, 0, len(words))

	for _, word := range words {
		result = append(result, word.Word)
	}

	return result, nil
}
