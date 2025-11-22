package models_err_app

import "errors"

var (
	ErrInitApp     = errors.New("failed init App")
	ErrInitService = errors.New("failed init service")
)
