package storage

import (
	"fmt"
	"telegram-antispam-bot/internal/infra/config"
	models_adds "telegram-antispam-bot/internal/models/adds"
	model_storage_tables "telegram-antispam-bot/internal/models/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Storage — слой доступа к данным через GORM.
type Storage struct {
	s *gorm.DB
}

// NewStorage создаёт подключение к БД, выполняет миграции и заполняет базовые списки.
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

	err = s.insertBaseWhitelist(conf)
	if err != nil {
		return nil, fmt.Errorf("s.insertBaseWhitelist: %w", err)
	}

	return s, nil
}

// autoMigrations создаёт/обновляет схему таблиц.
func (s *Storage) autoMigrations() error {
	err := s.s.AutoMigrate(
		&model_storage_tables.Word{},
		&model_storage_tables.WhitelistAuthor{},
		&model_storage_tables.WhitelistTag{},
	)
	if err != nil {
		return fmt.Errorf("s.AutoMigrate: %w", err)
	}

	return nil
}

// insertBaseList заполняет список запрещённых слов значениями по умолчанию, которых ещё нет в БД.
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

// insertBaseWhitelist при первом старте заполняет whitelist авторов и тэгов
// значениями из конфига, которых ещё нет в БД. Так конфиг остаётся источником
// начальных значений whitelist.
func (s *Storage) insertBaseWhitelist(conf *config.Config) error {
	authors, err := s.GetWhitelistAuthors()
	if err != nil {
		return fmt.Errorf("s.GetWhitelistAuthors: %w", err)
	}

	authorsDB := make(map[int64]struct{}, len(authors))
	for _, id := range authors {
		authorsDB[id] = struct{}{}
	}

	for _, id := range conf.BotAntiSpam.WhiteListAuthor {
		if _, exists := authorsDB[id]; !exists {
			err = s.InsertWhitelistAuthor(id)
			if err != nil {
				return fmt.Errorf("s.InsertWhitelistAuthor: %w", err)
			}
		}
	}

	tags, err := s.GetWhitelistTags()
	if err != nil {
		return fmt.Errorf("s.GetWhitelistTags: %w", err)
	}

	tagsDB := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagsDB[tag] = struct{}{}
	}

	for tag := range conf.BotAntiSpam.WhiteListTags {
		if _, exists := tagsDB[tag]; !exists {
			err = s.InsertWhitelistTag(tag)
			if err != nil {
				return fmt.Errorf("s.InsertWhitelistTag: %w", err)
			}
		}
	}

	return nil
}
