package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	logger       *log.Logger
	currentLevel = INFO
)

func Init(debug bool) {
	if debug {
		currentLevel = DEBUG
	}
	logger = log.New(os.Stdout, "", 0)
}

func Debug(msg string, keyValues ...interface{}) {
	if currentLevel <= DEBUG {
		logger.Print(formatMessage("DEBUG", msg, keyValues...))
	}
}

func Info(msg string, keyValues ...interface{}) {
	if currentLevel <= INFO {
		logger.Print(formatMessage("INFO", msg, keyValues...))
	}
}

func Warn(msg string, keyValues ...interface{}) {
	if currentLevel <= WARN {
		logger.Print(formatMessage("WARN", msg, keyValues...))
	}
}

func Error(msg string, keyValues ...interface{}) {
	if currentLevel <= ERROR {
		logger.Print(formatMessage("ERROR", msg, keyValues...))
	}
}

func formatMessage(level, msg string, keyValues ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	kvStr := ""
	for i := 0; i < len(keyValues); i += 2 {
		if i+1 < len(keyValues) {
			kvStr += fmt.Sprintf(" %s=%v", keyValues[i], keyValues[i+1])
		}
	}
	return fmt.Sprintf("[%s] %s: %s%s", timestamp, level, msg, kvStr)
}
