package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

// LogLevel represents the available log levels
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Config holds the logger configuration
type Config struct {
	Level     LogLevel `json:"level" env:"LOG_LEVEL" default:"info"`
	Format    string   `json:"format" env:"LOG_FORMAT" default:"text"`   // text, json
	Output    string   `json:"output" env:"LOG_OUTPUT" default:"stdout"` // stdout, stderr, file
	File      string   `json:"file" env:"LOG_FILE" default:"app.log"`
	Timestamp bool     `json:"timestamp" env:"LOG_TIMESTAMP" default:"true"`
	Caller    bool     `json:"caller" env:"LOG_CALLER" default:"false"`
	Colors    bool     `json:"colors" env:"LOG_COLORS" default:"true"`
}

// NewLogger creates a new structured logger with charmbracelet/log as handler
func NewLogger(config Config) *slog.Logger {
	// Create charmbracelet logger
	charmLogger := log.New(getOutput(config.Output, config.File))

	// Configure the charmbracelet logger
	charmLogger.SetLevel(mapLogLevel(config.Level))
	charmLogger.SetTimeFormat("2006-01-02 15:04:05")
		charmLogger.SetReportTimestamp(config.Timestamp)
	// Manually handle caller reporting via slog.Record.PC
	charmLogger.SetColorProfile(getColorProfile(config.Colors))

	// Set output format
	if config.Format == "json" {
		charmLogger.SetFormatter(log.JSONFormatter)
	} else {
		charmLogger.SetFormatter(log.TextFormatter)
	}

	// Create slog handler from charmbracelet logger
	handler := NewCharmSlogHandler(charmLogger, config.Caller)

	return slog.New(handler)
}

// NewDefaultLogger creates a logger with sensible defaults
func NewDefaultLogger() *slog.Logger {
	return NewLogger(Config{
		Level:     LevelInfo,
		Format:    "text",
		Output:    "stdout",
		Timestamp: true,
		Caller:    false,
		Colors:    true,
	})
}

// NewDevelopmentLogger creates a logger optimized for development
func NewDevelopmentLogger() *slog.Logger {
	return NewLogger(Config{
		Level:     LevelDebug,
		Format:    "text",
		Output:    "stdout",
		Timestamp: true,
		Caller:    true, // Caller is true for development
		Colors:    true,
	})
}

// NewProductionLogger creates a logger optimized for production
func NewProductionLogger() *slog.Logger {
	return NewLogger(Config{
		Level:     LevelInfo,
		Format:    "json",
		Output:    "stdout",
		Timestamp: true,
		Caller:    false,
		Colors:    false,
	})
}

func getOutput(output, file string) io.Writer {
	switch output {
	case "stderr":
		return os.Stderr
	case "file":
		f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Error("Failed to open log file, falling back to stdout", "error", err)
			return os.Stdout
		}
		return f
	default:
		return os.Stdout
	}
}

func mapLogLevel(level LogLevel) log.Level {
	switch level {
	case LevelDebug:
		return log.DebugLevel
	case LevelInfo:
		return log.InfoLevel
	case LevelWarn:
		return log.WarnLevel
	case LevelError:
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}

func getColorProfile(colors bool) termenv.Profile {
	if !colors {
		return termenv.Ascii // No color profile
	}
	return termenv.TrueColor
}

// CharmSlogHandler wraps charmbracelet/log to implement slog.Handler
type CharmSlogHandler struct {
	logger       *log.Logger
	reportCaller bool
}

func NewCharmSlogHandler(logger *log.Logger, reportCaller bool) *CharmSlogHandler {
	return &CharmSlogHandler{logger: logger, reportCaller: reportCaller}
}

func (h *CharmSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.mapCharmLevelToSlog(h.logger.GetLevel())
}

func (h *CharmSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Convert slog attributes to key-value pairs
	attrs := make([]interface{}, 0, record.NumAttrs()*2)
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a.Key, a.Value.Any())
		return true
	})

	// Prepend source information to the message if enabled and available
	msg := record.Message
	if h.reportCaller && record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()
		// Use filepath.Base to get just the file name, not the full path
		caller := fmt.Sprintf("\x1b[90m<%s:%d>\x1b[0m", filepath.Base(f.File), f.Line)
		msg = fmt.Sprintf("%s %s", caller, record.Message)
	}

	// Log with appropriate level
	switch record.Level {
	case slog.LevelDebug:
		h.logger.Debug(msg, attrs...)
	case slog.LevelInfo:
		h.logger.Info(msg, attrs...)
	case slog.LevelWarn:
		h.logger.Warn(msg, attrs...)
	case slog.LevelError:
		h.logger.Error(msg, attrs...)
	default:
		h.logger.Info(msg, attrs...)
	}

	return nil
}

func (h *CharmSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new logger with prefix attributes
	newLogger := h.logger.With()
	for _, attr := range attrs {
		newLogger = newLogger.With(attr.Key, attr.Value.Any())
	}
	return NewCharmSlogHandler(newLogger, h.reportCaller)
}

func (h *CharmSlogHandler) WithGroup(name string) slog.Handler {
	// Charmbracelet/log doesn't have native group support,
	// so we'll prefix keys with the group name
	return &CharmSlogHandler{logger: h.logger.WithPrefix(name + "."), reportCaller: h.reportCaller}
}

func (h *CharmSlogHandler) slogLevelToCharmLevel(level slog.Level) slog.Level {
	switch {
	case level < slog.LevelInfo:
		return slog.LevelDebug
	case level < slog.LevelWarn:
		return slog.LevelInfo
	case level < slog.LevelError:
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}

func (h *CharmSlogHandler) mapCharmLevelToSlog(level log.Level) slog.Level {
	switch level {
	case log.DebugLevel:
		return slog.LevelDebug
	case log.InfoLevel:
		return slog.LevelInfo
	case log.WarnLevel:
		return slog.LevelWarn
	case log.ErrorLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}