package model_storage_tables

import "gorm.io/gorm"

// Word — запрещённое слово, по которому сообщение признаётся спамом.
type Word struct {
	gorm.Model
	Word string
}

// WhitelistAuthor — пользователь, чьи сообщения игнорируются проверками антиспама.
type WhitelistAuthor struct {
	gorm.Model
	TelegramID int64 `gorm:"uniqueIndex"`
}

// WhitelistTag — тэг, наличие которого заставляет игнорировать сообщение.
type WhitelistTag struct {
	gorm.Model
	Tag string `gorm:"uniqueIndex"`
}
