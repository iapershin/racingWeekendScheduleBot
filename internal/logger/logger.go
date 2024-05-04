package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

var ErrInvalidFormat = errors.New("invalid format")

type Format string

const (
	JSON Format = "json"
	TEXT Format = "text"
)

type Level = slog.Level

type Logger struct {
	*slog.Logger
}

type HandlerOptions struct {
	Format    Format
	Level     string
	AddSource bool
}

func Init(opt HandlerOptions) error {
	var handler slog.Handler

	var lvl slog.Level
	err := lvl.UnmarshalText([]byte(opt.Level))
	if err != nil {
		return fmt.Errorf("failed to unmarshal log level: %w", err)
	}

	switch opt.Format {
	case JSON:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: opt.AddSource,
		})
	case TEXT:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: opt.AddSource,
		})
	default:
		return fmt.Errorf("%w: %s", ErrInvalidFormat, opt.Format)
	}

	slog.SetDefault(slog.New(handler))
	return nil
}
