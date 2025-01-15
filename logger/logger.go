package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// InitLogger initializes the global logger
func InitLogger(debug bool) {
	var cfg zap.Config
	if debug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.LevelKey = "level"
		cfg.EncoderConfig.NameKey = "logger"
		cfg.EncoderConfig.MessageKey = "message"
		cfg.EncoderConfig.StacktraceKey = "stacktrace"
	}

	// Build the logger
	var err error
	log, err = cfg.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

// Info logs an informational message
func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zapcore.Field) {
	log.Debug(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() {
	_ = log.Sync() // Flushes any buffered logs; ignore the error for simplicity
}
