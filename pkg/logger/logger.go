package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// Init creates new Logger object.
func Init(env string) {
	atomic := zap.NewAtomicLevel()

	var enconfig zapcore.EncoderConfig
	var enc zapcore.Encoder

	switch env {
	case "prod":
		atomic.SetLevel(zapcore.InfoLevel)
		enconfig = zap.NewProductionEncoderConfig()
		enc = zapcore.NewJSONEncoder(enconfig)

	default:
		atomic.SetLevel(zapcore.DebugLevel)
		enconfig = zap.NewDevelopmentEncoderConfig()
		enc = zapcore.NewConsoleEncoder(enconfig)
	}

	globalLogger = zap.New(zapcore.NewCore(
		enc,
		zapcore.Lock(os.Stdout),
		atomic,
	))
}

// Debug prints debug-level message.
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info prints info-level message.
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn prints warning-level message.
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error prints error-level message.
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Fatal prints fatal-level message.
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// WithOptions clones the current Logger, applies the supplied Options, and returns the resulting Logger.
func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
