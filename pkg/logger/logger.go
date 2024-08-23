package logger

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	LvlDebug = int(slog.LevelDebug)
	LvlInfo  = int(slog.LevelInfo)
	LvlError = int(slog.LevelError)
)

type Log struct {
	*slog.Logger
}

type Config struct {
	Level string
}

func New(cfg Config) *Log {
	opts := slog.HandlerOptions{}

	switch cfg.Level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &opts))

	slog.SetDefault(log)

	return &Log{log}
}

func (l *Log) Info(format string, msg ...any) {
	l.Logger.Info(fmt.Sprintf(format, msg...))
}

func (l *Log) Debug(err error, format string, msg ...any) {
	l.Logger.Debug(
		fmt.Sprintf(format, msg...),
		handleErr(err),
	)
}

func (l *Log) Error(err error, format string, msg ...any) {
	l.Logger.Error(
		fmt.Sprintf(format, msg...),
		handleErr(err),
	)
}

func handleErr(err error) slog.Attr {
	if err != nil {
		return slog.String("err", err.Error())
	}

	return slog.Attr{}
}
