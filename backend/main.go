package main

import (
	"fmt"
	"os"
	"pan/backend/handlers"
	"pan/backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

var log = utils.Log()

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 记录请求开始
		log.Debug("Request started | Method: %s | Path: %s | Client: %s",
			c.Request.Method, c.Request.URL.Path, c.ClientIP())

		c.Next()

		// 计算请求处理时间
		latencyTime := time.Since(startTime)
		statusCode := c.Writer.Status()

		// 根据状态码使用不同的日志级别
		if statusCode >= 500 {
			log.Error("Request completed | Status: %d | Latency: %v | Method: %s | Path: %s | Client: %s | Errors: %v",
				statusCode, latencyTime, c.Request.Method, c.Request.URL.Path, c.ClientIP(), c.Errors.String())
		} else if statusCode >= 400 {
			log.Warn("Request completed | Status: %d | Latency: %v | Method: %s | Path: %s | Client: %s",
				statusCode, latencyTime, c.Request.Method, c.Request.URL.Path, c.ClientIP())
		} else {
			log.Info("Request completed | Status: %d | Latency: %v | Method: %s | Path: %s | Client: %s",
				statusCode, latencyTime, c.Request.Method, c.Request.URL.Path, c.ClientIP())
		}
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // 默认使用8081端口
	}
	return fmt.Sprintf(":%s", port)
}

func main() {
	// 设置为调试模式
	gin.SetMode(gin.DebugMode)
	log.SetLevel(utils.DEBUG)

	r := gin.New()
	r.Use(RequestLogger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Set maximum file upload size (100MB)
	r.MaxMultipartMemory = 100 << 20

	// 创建上传目录
	uploadPath := "./uploads"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Fatal("Failed to create upload directory: %v", err)
	}
	log.Debug("Upload directory created/verified: %s", uploadPath)

	// 初始化文件处理器
	fileHandler := handlers.NewFileHandler(uploadPath)
	log.Debug("File handler initialized with upload path: %s", uploadPath)

	// File routes
	fileGroup := r.Group("/api/files")
	{
		fileGroup.POST("/upload", fileHandler.Upload)
		fileGroup.GET("/download/:filename", fileHandler.Download)
		fileGroup.GET("/preview/:filename", fileHandler.Preview)
		fileGroup.GET("/list", fileHandler.List)
		fileGroup.DELETE("/delete/:filename", fileHandler.Delete)
	}

	// 启动服务器
	port := getPort()
	log.Info("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
