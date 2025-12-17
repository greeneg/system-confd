package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greeneg/system-confd/globals"
)

var logFormatString = "2006/01/02 15:04:05"

type Logger struct {
	Level  string   `ini:"level"`
	Target *os.File `ini:"target"`
}

func setGinLogMode(level string) {
	switch level {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "info":
		gin.SetMode(gin.ReleaseMode)
	case "warn":
		gin.SetMode(gin.ReleaseMode)
	case "error":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode) // default to release mode
	}
}

func createLogDir(logFile string) error {
	dir := filepath.Dir(logFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// There is nothing exposed by our logs that would be sensitive
		// so we can create the directory with 0755 permissions
		// #nosec G301
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func setupLogger(config globals.Config) (Logger, error) {
	var logger Logger

	// Set log level
	logger.Level = config.General.LogLevel
	if logger.Level == "" {
		logger.Level = "info" // default log level
		setGinLogMode(logger.Level)
	}

	// Set log target
	if config.General.LogFile != "" {
		err := createLogDir(config.General.LogFile)
		if err != nil {
			return logger, err
		}
		// Our logs do not contain sensitive information
		// so we can create the log file with 0640 permissions
		// #nosec G302
		fh, err := os.OpenFile(config.General.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		if err != nil {
			return logger, err
		}
		logger.Target = fh
	} else {
		logger.Target = os.Stdout // default to stdout if no log file is specified
	}

	return logger, nil
}

func (l *Logger) Info(message string) {
	currTime := time.Now().Format(logFormatString)
	if l.Target != nil {
		_, err := l.Target.WriteString(string(currTime) + " INFO: " + message + "\n")
		if err != nil {
			log.Println("Failed to write info message to log:", err)
		}
	}
}

func (l *Logger) Error(message string) {
	currTime := time.Now().Format(logFormatString)
	if l.Target != nil {
		_, err := l.Target.WriteString(string(currTime) + " ERROR: " + message + "\n")
		if err != nil {
			log.Println("Failed to write error message to log:", err)
		}
		log.Fatal("ERROR " + message) // Log fatal error and exit
	}
}

func (l *Logger) Debug(message string) {
	currTime := time.Now().Format(logFormatString)
	if l.Level == "debug" && l.Target != nil {
		_, err := l.Target.WriteString(string(currTime) + " DEBUG: " + message + "\n")
		if err != nil {
			log.Println("Failed to write debug message to log:", err)
		}
	}
}

func (l *Logger) Warn(message string) {
	currTime := time.Now().Format(logFormatString)
	if l.Target != nil {
		_, err := l.Target.WriteString(string(currTime) + " WARN: " + message + "\n")
		if err != nil {
			log.Println("Failed to write warn message to log:", err)
		}
	}
}
