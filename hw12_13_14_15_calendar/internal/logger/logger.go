package logger

import (
	"log/slog"
	"os"
	"strings"
)

const (
	lInfo  = "info"
	lDebug = "debug"
	lWarn  = "warn"
	lError = "error"
)

var levelMap = map[string]slog.Level{
	lDebug: slog.LevelDebug,
	lInfo:  slog.LevelInfo,
	lWarn:  slog.LevelWarn,
	lError: slog.LevelError,
}

type Logger struct {
	level string
}

func New(level string) *Logger {
	var programLevel = new(slog.LevelVar)
	programLevel.Set(levelMap[level])
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	return &Logger{
		level: strings.ToLower(level),
	}
}

func (l Logger) Info(msg string, args... any) {
		slog.Info(msg, args...)
}

func (l Logger) Error(msg string, args... any) {
	slog.Error(msg, args...)
}

func (l Logger) Warn(msg string, args... any) {
	slog.Warn(msg, args...)
}

func (l Logger) Debug(msg string, args... any) {
	slog.Debug(msg, args...)
}