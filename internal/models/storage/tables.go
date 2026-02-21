package model_storage_tables

import "gorm.io/gorm"

type Word struct {
	gorm.Model
	Word string
}
