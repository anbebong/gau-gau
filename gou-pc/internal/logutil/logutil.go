package logutil

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	ERROR
)

var (
	logger   *log.Logger
	logLevel Level = INFO
)

func Init(logfile string, level Level) error {
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	logger = log.New(io.MultiWriter(f, os.Stdout), "", log.LstdFlags|log.Lshortfile)
	logLevel = level
	return nil
}

func Debug(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		logger.Output(2, "[DEBUG] "+formatMsg(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	if logLevel <= INFO {
		logger.Output(2, "[INFO] "+formatMsg(format, v...))
	}
}

func Error(format string, v ...interface{}) {
	if logLevel <= ERROR {
		logger.Output(2, "[ERROR] "+formatMsg(format, v...))
	}
}

func formatMsg(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}
