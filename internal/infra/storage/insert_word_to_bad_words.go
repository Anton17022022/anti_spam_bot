package storage

import (
	"fmt"
	model_storage_tables "telegram-antispam-bot/internal/models/storage"
)

func (s *Storage) InsertWordToBadWords(word string) error {
	newWord := model_storage_tables.Word{Word: word}

	result := s.s.Create(&newWord)
	if result.Error != nil {
		return fmt.Errorf("s.Create: %w", result.Error)
	}

	return nil
}
