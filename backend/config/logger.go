package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initJSONLogger(app *AppConfig) (*zap.Logger, error) {

	if err := os.MkdirAll("./logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	jsonFile := "logs.json"

	jsonPath := filepath.Join("logs", jsonFile)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	consoleEncoderConfig := encoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	fileWriter, err := createFileWriter(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file writer: %w", err)
	}

	consoleWriter := zapcore.AddSync(os.Stdout)

	logLevel := zapcore.InfoLevel
	if app.IsDebug() {
		logLevel = zapcore.DebugLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, fileWriter, logLevel),
		zapcore.NewCore(consoleEncoder, consoleWriter, logLevel),
	)

	// Buat logger dengan sampling
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("app", app.Name),
			zap.String("env", app.Env),
		),
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				500, // initial
				50,  // thereafter
			)
		}),
	)

	zap.RedirectStdLog(logger)

	return logger, nil
}
func createFileWriter(path string) (zapcore.WriteSyncer, error) {
	writer := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}
	// Test apakah file bisa ditulis
	if _, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return nil, err
	}

	return zapcore.AddSync(writer), nil
}
func GetLogger(app *AppConfig) *zap.Logger {
	if Logger == nil {
		loggerOnce.Do(func() {
			var err error
			Logger, err = initJSONLogger(app)
			if err != nil {
				// Gunakan fallback logger dan log error ke stderr
				fallbackLogger.Error("Failed to initialize custom logger",
					zap.Error(err),
					zap.String("fallback", "using example logger"),
				)
				Logger = fallbackLogger
			}
		})
	}
	return Logger
}
func CloseLogger() {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	if Logger != nil {
		// Sync() penting untuk file writer, tapi bisa error di Windows
		err := Logger.Sync()
		if err != nil && !isHarmlessSyncError(err) {
			// Log error ke stderr jika logger sudah tidak tersedia
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}
}
func isHarmlessSyncError(err error) bool {
	return runtime.GOOS == "windows" &&
		strings.Contains(err.Error(), "The handle is invalid")
}
