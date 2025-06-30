package logutil

import (
	"fmt"
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	ERROR
)

type Logger struct {
	logger   *log.Logger
	logLevel Level
}

var (
	coreLogger *Logger
	apiLogger  *Logger
)

func InitCoreLogger(logfile string, level Level) error {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	coreLogger = &Logger{
		logger:   log.New(f, "", log.LstdFlags|log.Lshortfile), // chỉ ghi ra file core
		logLevel: level,
	}
	return nil
}

func InitAPILogger(logfile string, level Level) error {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	apiLogger = &Logger{
		logger:   log.New(f, "", log.LstdFlags|log.Lshortfile), // chỉ ghi ra file api
		logLevel: level,
	}
	return nil
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l != nil && l.logLevel <= DEBUG {
		l.logger.Output(3, "[DEBUG] "+formatMsg(format, v...))
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l != nil && l.logLevel <= INFO {
		l.logger.Output(3, "[INFO] "+formatMsg(format, v...))
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l != nil && l.logLevel <= ERROR {
		l.logger.Output(3, "[ERROR] "+formatMsg(format, v...))
	}
}

// Hàm tiện cho code cũ
func CoreDebug(format string, v ...interface{}) { coreLogger.Debug(format, v...) }
func CoreInfo(format string, v ...interface{})  { coreLogger.Info(format, v...) }
func CoreError(format string, v ...interface{}) { coreLogger.Error(format, v...) }
func APIDebug(format string, v ...interface{})  { apiLogger.Debug(format, v...) }
func APIInfo(format string, v ...interface{})   { apiLogger.Info(format, v...) }
func APIError(format string, v ...interface{})  { apiLogger.Error(format, v...) }

func formatMsg(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}
