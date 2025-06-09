package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelStrings = [...]string{"DEBUG", "INFO", "WARN", "ERROR"}

type Logger struct {
	level  LogLevel
	logger *log.Logger
	mu     sync.Mutex
}

func NewLogger(level LogLevel, out io.Writer) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(out, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) logf(lv LogLevel, format string, v ...interface{}) {
	if lv < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	prefix := fmt.Sprintf("[%s] ", levelStrings[lv])
	l.logger.Output(3, prefix+fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(format string, v ...interface{}) { l.logf(DEBUG, format, v...) }
func (l *Logger) Info(format string, v ...interface{})  { l.logf(INFO, format, v...) }
func (l *Logger) Warn(format string, v ...interface{})  { l.logf(WARN, format, v...) }
func (l *Logger) Error(format string, v ...interface{}) { l.logf(ERROR, format, v...) }

// Helper to create logger that writes to file
func NewFileLogger(level LogLevel, filePath string) (*Logger, error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return NewLogger(level, f), nil
}

// Helper to create logger that writes to both file and stdout
func NewMultiLogger(level LogLevel, filePath string) (*Logger, error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	mw := io.MultiWriter(os.Stdout, f)
	return NewLogger(level, mw), nil
}
