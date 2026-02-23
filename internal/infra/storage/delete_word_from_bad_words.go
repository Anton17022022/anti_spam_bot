package storage

import (
	"fmt"
	model_storage_tables "telegram-antispam-bot/internal/models/storage"
)

func (s *Storage) DelWordFromBadWords(word string) error {
	result := s.s.Where("word = ?", word).Delete(&model_storage_tables.Word{})
	if result.Error != nil {
		return fmt.Errorf("s.Delete: %w", result.Error)
	}

	return nil
}
