package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func NewLogger(level, format string) (*slog.Logger, error) {
	logLevel, err := parseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:       logLevel,
		AddSource:   logLevel == slog.LevelDebug,
		ReplaceAttr: replaceAttr,
	}

	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text", "pretty":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		return nil, fmt.Errorf("unknown log format: %s", format)
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger, nil
}

func MustCreateLogger(level, format string) *slog.Logger {
	value, err := NewLogger(level, format)
	if err != nil {
		panic(fmt.Errorf("failed to load logger: %w", err))
	}
	return value
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level: %s", level)
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = shortenPath(source.File)
		source.Function = shortenPath(source.Function)
	}

	return a
}

func shortenPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return path
}
