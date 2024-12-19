package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pan/backend/utils"
	"path"
	"path/filepath"
	"strconv"
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
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"size":     utils.FormatFileSize(file.Size),
		"type":     utils.GetFileType(file.Filename),
	})
}

// IsPreviewable checks if a file can be previewed based on its extension
func IsPreviewable(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" || ext == ".pdf" || ext == ".epub"
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

// IsEPUB checks if a file is an EPUB based on its extension
func IsEPUB(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".epub"
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

    // Set appropriate headers based on file type
    if IsPDF(filename) {
        c.Header("Content-Type", "application/pdf")
    } else if IsEPUB(filename) {
        c.Header("Content-Type", "application/epub+zip")
        c.Header("Accept-Ranges", "bytes")
        c.Header("Cache-Control", "no-cache")
        c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
    } else if IsImage(filename) {
		// 根据图片扩展名设置正确的 MIME 类型
		ext := strings.ToLower(path.Ext(filename))
		switch ext {
		case ".jpg", ".jpeg":
			c.Header("Content-Type", "image/jpeg")
		case ".png":
			c.Header("Content-Type", "image/png")
		case ".gif":
			c.Header("Content-Type", "image/gif")
		case ".webp":
			c.Header("Content-Type", "image/webp")
		}
	}

	// Add CORS headers
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, HEAD")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Range, Content-Disposition, Content-Description, X-Requested-With")
	c.Header("Access-Control-Expose-Headers", "Content-Range, Content-Length")

	// Handle OPTIONS request
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusOK)
		return
	}

	// Handle range requests
	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		fileInfo, err := os.Stat(filepath)
		if err != nil {
			log.Error("Failed to get file info: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
			return
		}

		fileSize := fileInfo.Size()
		c.Header("Content-Length", fmt.Sprintf("%d", fileSize))

		// Parse range header
		ranges, err := parseRange(rangeHeader, fileSize)
		if err != nil {
			log.Error("Failed to parse range header: %v", err)
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		if len(ranges) > 0 {
			start, end := ranges[0].start, ranges[0].end
			c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
			c.Status(http.StatusPartialContent)

			// Open the file
			file, err := os.Open(filepath)
			if err != nil {
				log.Error("Failed to open file: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
				return
			}
			defer file.Close()

			// Seek to the start position
			file.Seek(start, 0)

			// Create a limited reader for the range
			limitReader := io.LimitReader(file, end-start+1)
			c.DataFromReader(http.StatusPartialContent, end-start+1, c.ContentType(), limitReader, nil)
			return
		}
	}

	// If no range request or invalid range, serve the entire file
	c.File(filepath)
}

// Range represents a byte range
type Range struct {
	start, end int64
}

// parseRange parses a Range header string as per RFC 7233
func parseRange(rangeHeader string, size int64) ([]Range, error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return nil, fmt.Errorf("invalid range format")
	}

	ranges := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), ",")
	parsedRanges := make([]Range, 0, len(ranges))

	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			continue
		}

		var start, end int64
		var err error

		if parts[0] == "" {
			// suffix length
			end = size - 1
			start = size - 1
			if parts[1] != "" {
				if start, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
					continue
				}
				start = size - start
			}
		} else {
			if start, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
				continue
			}
			if parts[1] != "" {
				if end, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
					continue
				}
			} else {
				end = size - 1
			}
		}

		if start > end || start < 0 || end >= size {
			continue
		}

		parsedRanges = append(parsedRanges, Range{start: start, end: end})
	}

	if len(parsedRanges) == 0 {
		return nil, fmt.Errorf("no valid ranges")
	}

	return parsedRanges, nil
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
			"name":     file.Name(),
			"size":     utils.FormatFileSize(info.Size()),
			"type":     utils.GetFileType(file.Name()),
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
