package logger

import (
	"backend/internal/repository"
	"context"
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// InitializeLogger initializes the global logger.
func InitializeLogger(logRepo *repository.LogRepository) {
	// Create a custom MongoDB core
	mongoCore := NewMongoDBCore(logRepo)

	// Create the logger with the MongoDB core
	globalLogger = zap.New(mongoCore)
	zap.ReplaceGlobals(globalLogger)
}

// NewMongoDBCore creates a custom zapcore.Core for MongoDB.
func NewMongoDBCore(logRepo *repository.LogRepository) zapcore.Core {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&mongoWriter{logRepo: logRepo}),
		zapcore.DebugLevel,
	)
}

type mongoWriter struct {
	logRepo *repository.LogRepository
}

func (m *mongoWriter) Write(p []byte) (n int, err error) {
	var logEntry map[string]interface{}
	err = json.Unmarshal(p, &logEntry)
	if err != nil {
		return 0, err
	}

	// Insert log into the appropriate collection
	if level, ok := logEntry["level"].(string); ok && level == "error" {
		err = m.logRepo.InsertErrorLog(context.Background(), logEntry)
	} else {
		err = m.logRepo.InsertLog(context.Background(), logEntry)
	}

	return len(p), err
}

func (m *mongoWriter) Sync() error {
	return nil
}

// L returns the global logger for access in other files.
func L() *zap.Logger {
	return globalLogger
}
