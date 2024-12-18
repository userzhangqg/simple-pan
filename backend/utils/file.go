package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// MaxFileSize 最大文件大小 (100MB)
	MaxFileSize = 100 << 20
)

// 允许的文件类型
var AllowedFileTypes = map[string]bool{
	".txt":  true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".mp4":  true,
	".zip":  true,
	".rar":  true,
}

// ValidateFile checks if the file is valid (size and type)
func ValidateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > MaxFileSize {
		return ErrFileTooLarge
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedFileTypes[ext] {
		return ErrInvalidFileType
	}

	return nil
}

// GetFileType returns the general type of the file based on its extension
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return "image"
	case ".mp4", ".avi", ".mov":
		return "video"
	case ".doc", ".docx", ".pdf", ".txt":
		return "document"
	case ".zip", ".rar":
		return "archive"
	default:
		return "other"
	}
}

// FormatFileSize converts file size to human readable format
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// StringToInt converts a string to an integer with a default value of 1 if conversion fails
func StringToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 1
	}
	return val
}
