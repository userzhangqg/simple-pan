package utils

import (
	"fmt"
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	logger    *log.Logger
	logLevel  LogLevel
	levelTags map[LogLevel]string
}

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(DEBUG) // 默认使用 DEBUG 级别
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{
		logger:   log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		logLevel: level,
		levelTags: map[LogLevel]string{
			DEBUG: "DEBUG",
			INFO:  "INFO ",
			WARN:  "WARN ",
			ERROR: "ERROR",
			FATAL: "FATAL",
		},
	}
}

// Log 获取默认日志实例
func Log() *Logger {
	return defaultLogger
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.logLevel = level
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() LogLevel {
	return l.logLevel
}

func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	if level >= l.logLevel {
		l.logger.Printf("[%s] %s", l.levelTags[level], fmt.Sprintf(format, v...))
	}
}

// Debug 打印调试日志
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

// Info 打印信息日志
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

// Warn 打印警告日志
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

// Error 打印错误日志
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

// Fatal 打印致命错误日志并退出程序
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FATAL, format, v...)
	os.Exit(1)
}
