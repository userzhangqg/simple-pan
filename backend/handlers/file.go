package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"pan/backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

var log = utils.Log()

// FileHandler handles file-related operations
type FileHandler struct {
	uploadPath string
}

// NewFileHandler creates a new FileHandler instance
func NewFileHandler(uploadPath string) *FileHandler {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Fatal("Failed to create upload directory: %v", err)
	}
	log.Debug("Upload directory created: %s", uploadPath)
	return &FileHandler{
		uploadPath: uploadPath,
	}
}

// Upload handles file upload
func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file
	if err := utils.ValidateFile(file); err != nil {
		log.Warn("File validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Debug("File validation passed: %s (size: %d)", file.Filename, file.Size)

	filename := filepath.Join(h.uploadPath, file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		log.Error("Failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	log.Info("File uploaded successfully: %s", filename)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"filename": file.Filename,
		"size": utils.FormatFileSize(file.Size),
		"type": utils.GetFileType(file.Filename),
	})
}

// IsPreviewable checks if a file can be previewed based on its extension
func IsPreviewable(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" || ext == ".pdf"
}

// IsImage checks if a file is an image based on its extension
func IsImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
}

// IsPDF checks if a file is a PDF based on its extension
func IsPDF(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".pdf"
}

// Preview handles file preview
func (h *FileHandler) Preview(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join(h.uploadPath, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Warn("Preview failed: file not found: %s", filepath)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	if !IsPreviewable(filename) {
		log.Warn("Preview failed: file format not supported for preview: %s", filepath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File format not supported for preview"})
		return
	}

	// Set appropriate headers
	if IsPDF(filename) {
		c.Header("Content-Type", "application/pdf")
	}
	
	// Add CORS headers
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
	
	// Add cache control headers
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.File(filepath)
}

// Download handles file download
func (h *FileHandler) Download(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join(h.uploadPath, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Warn("Download failed: file not found: %s", filepath)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	log.Debug("Starting file download: %s", filepath)

	c.File(filepath)
}

// List returns a paginated and filtered list of uploaded files
func (h *FileHandler) List(c *gin.Context) {
	log.Debug("Starting to list files in directory: %s", h.uploadPath)
	
	// Get query parameters
	page := utils.StringToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StringToInt(c.DefaultQuery("pageSize", "10"))
	search := strings.ToLower(c.DefaultQuery("search", ""))
	
	// Read all files
	files, err := os.ReadDir(h.uploadPath)
	if err != nil {
		log.Error("Failed to read directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory"})
		return
	}

	var fileList []map[string]interface{}
	var filteredFiles []map[string]interface{}

	// Build file list with filtering
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Warn("Failed to get file info for %s: %v", file.Name(), err)
			continue
		}
		
		fileData := map[string]interface{}{
			"name": file.Name(),
			"size": utils.FormatFileSize(info.Size()),
			"type": utils.GetFileType(file.Name()),
			"modified": info.ModTime(),
		}
		
		fileList = append(fileList, fileData)
		
		// Apply search filter
		if search == "" || strings.Contains(strings.ToLower(file.Name()), search) {
			filteredFiles = append(filteredFiles, fileData)
		}
	}

	// Calculate pagination
	total := len(filteredFiles)
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	if start > total {
		start = total
	}

	// Get paginated results
	var paginatedFiles []map[string]interface{}
	if start < total {
		paginatedFiles = filteredFiles[start:end]
	}

	log.Info("Listed %d files (filtered from %d total files)", len(paginatedFiles), len(fileList))
	// 格式化文件列表以确保 JSON 格式正确
	formattedFiles := make([]map[string]interface{}, len(paginatedFiles))
	for i, file := range paginatedFiles {
		formattedFiles[i] = map[string]interface{}{
			"name":     file["name"],
			"size":     file["size"],
			"type":     file["type"],
			"modified": file["modified"],
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"files":    formattedFiles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// Delete handles file deletion
func (h *FileHandler) Delete(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join(h.uploadPath, filename)
	log.Debug("Attempting to delete file: %s", filepath)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Warn("Delete failed: file not found: %s", filepath)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	fileInfo, err := os.Stat(filepath)
	if err == nil {
		log.Debug("File found: %s (size: %s)", filename, utils.FormatFileSize(fileInfo.Size()))
	}

	if err := os.Remove(filepath); err != nil {
		log.Error("Failed to delete file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}
	log.Info("File deleted successfully: %s", filepath)

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
