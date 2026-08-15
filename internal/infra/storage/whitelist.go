package storage

import (
	"fmt"

	model_storage_tables "telegram-antispam-bot/internal/models/storage"
)

// GetWhitelistAuthors возвращает telegram id авторов из whitelist.
func (s *Storage) GetWhitelistAuthors() ([]int64, error) {
	var authors = make([]model_storage_tables.WhitelistAuthor, 0)

	res := s.s.Find(&authors)
	if res.Error != nil {
		return nil, fmt.Errorf("s.Find: %w", res.Error)
	}

	result := make([]int64, 0, len(authors))
	for _, author := range authors {
		result = append(result, author.TelegramID)
	}

	return result, nil
}

// InsertWhitelistAuthor добавляет автора в whitelist.
func (s *Storage) InsertWhitelistAuthor(id int64) error {
	newAuthor := model_storage_tables.WhitelistAuthor{TelegramID: id}

	result := s.s.Create(&newAuthor)
	if result.Error != nil {
		return fmt.Errorf("s.Create: %w", result.Error)
	}

	return nil
}

// DelWhitelistAuthor удаляет автора из whitelist.
func (s *Storage) DelWhitelistAuthor(id int64) error {
	result := s.s.Where("telegram_id = ?", id).Delete(&model_storage_tables.WhitelistAuthor{})
	if result.Error != nil {
		return fmt.Errorf("s.Delete: %w", result.Error)
	}

	return nil
}

// IsWhitelistAuthor проверяет, находится ли автор в whitelist.
func (s *Storage) IsWhitelistAuthor(id int64) (bool, error) {
	var count int64

	result := s.s.Model(&model_storage_tables.WhitelistAuthor{}).
		Where("telegram_id = ?", id).
		Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("s.Count: %w", result.Error)
	}

	return count > 0, nil
}

// GetWhitelistTags возвращает список тэгов из whitelist.
func (s *Storage) GetWhitelistTags() ([]string, error) {
	var tags = make([]model_storage_tables.WhitelistTag, 0)

	res := s.s.Find(&tags)
	if res.Error != nil {
		return nil, fmt.Errorf("s.Find: %w", res.Error)
	}

	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		result = append(result, tag.Tag)
	}

	return result, nil
}

// InsertWhitelistTag добавляет тэг в whitelist.
func (s *Storage) InsertWhitelistTag(tag string) error {
	newTag := model_storage_tables.WhitelistTag{Tag: tag}

	result := s.s.Create(&newTag)
	if result.Error != nil {
		return fmt.Errorf("s.Create: %w", result.Error)
	}

	return nil
}

// DelWhitelistTag удаляет тэг из whitelist.
func (s *Storage) DelWhitelistTag(tag string) error {
	result := s.s.Where("tag = ?", tag).Delete(&model_storage_tables.WhitelistTag{})
	if result.Error != nil {
		return fmt.Errorf("s.Delete: %w", result.Error)
	}

	return nil
}

// IsWhitelistTag проверяет, находится ли тэг в whitelist.
func (s *Storage) IsWhitelistTag(tag string) (bool, error) {
	var count int64

	result := s.s.Model(&model_storage_tables.WhitelistTag{}).
		Where("tag = ?", tag).
		Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("s.Count: %w", result.Error)
	}

	return count > 0, nil
}
