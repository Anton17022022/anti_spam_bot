package storage

import (
	"fmt"
	"telegram-antispam-bot/internal/infra/config"
	models_adds "telegram-antispam-bot/internal/models/adds"
	model_storage_tables "telegram-antispam-bot/internal/models/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Storage ..
type Storage struct {
	s *gorm.DB
}

// NewStorage ..
func NewStorage(conf *config.Config) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(conf.Storage.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}

	s := &Storage{
		s: db,
	}

	err = s.autoMigrations()
	if err != nil {
		return nil, fmt.Errorf("s.autoMigrations: %w", err)
	}

	err = s.insertBaseList()
	if err != nil {
		return nil, fmt.Errorf("s.insertBaseList: %w", err)
	}

	return s, nil
}

// autoMigrations ..
func (s *Storage) autoMigrations() error {
	err := s.s.AutoMigrate(&model_storage_tables.Word{})
	if err != nil {
		return fmt.Errorf("s.AutoMigrate: %w", err)
	}

	return nil
}

func (s *Storage) insertBaseList() error {
	words, err := s.GetListBadWords()
	if err != nil {
		return fmt.Errorf("s.GetListBadWords: %w", err)
	}

	wdDB := make(map[string]struct{})
	for _, word := range words {
		wdDB[word] = struct{}{}
	}

	for _, word := range models_adds.AdKeywords {
		if _, exists := wdDB[word]; !exists {
			err = s.InsertWordToBadWords(word)
			if err != nil {
				return fmt.Errorf("s.InsertWordToBadWords: %w", err)
			}
		}
	}

	return nil
}
