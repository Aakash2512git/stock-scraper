package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"
)

type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
}

var (
	instance *Logger
	once     sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			logger: log.New(os.Stdout, "", 0),
		}
	})
	return instance
}

func (l *Logger) log(level string, msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	output := fmt.Sprintf("[%s] [%s] %s", timestamp, level, formatted)

	l.logger.Println(output)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(LevelInfo, msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.log(LevelWarning, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(LevelError, msg, args...)
}
