package utils

import "errors"

var (
	ErrFileTooLarge    = errors.New("file size exceeds maximum limit")
	ErrInvalidFileType = errors.New("invalid file type")
)
