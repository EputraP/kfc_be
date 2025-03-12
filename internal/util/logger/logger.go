package logger

import (
	"log/slog"
	"os"
	"time"
)

var logger *slog.Logger

func Init(logFilePath string) error {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	cleanupOldLogs()

	return nil
}

func cleanupOldLogs() {
	// Define the maximum age for logs (e.g., delete logs older than 7 days)
	maxAge := 30 * 24 * time.Hour

	// Check if the log file exists
	logFile := "app.log"
	fileInfo, err := os.Stat(logFile)
	if err != nil {
		logger.Error("logger", "Error reading log file:", map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Check if the file is older than the max age
	if time.Since(fileInfo.ModTime()) > maxAge {
		err := os.Remove(logFile)
		if err != nil {
			logger.Error("logger", "Error deleting old log file:", map[string]string{
				"error": err.Error(),
			})
			return
		} else {
			logger.Info("logger", "Old log file deleted.", nil)
		}
	}
}

func Info(msg string, process string, details interface{}) {
	logger.Info(msg,
		slog.String("process", process),
		slog.Any("details", details))
}

func Debug(msg string, process string, details interface{}) {
	logger.Debug(msg,
		slog.String("process", process),
		slog.Any("details", details))
}

func Warn(msg string, process string, details interface{}) {
	logger.Warn(msg,
		slog.String("process", process),
		slog.Any("details", details))
}

func Error(msg string, process string, details interface{}) {
	logger.Error(msg,
		slog.String("process", process),
		slog.Any("details", details))
}
