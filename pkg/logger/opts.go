package logger

import "log/slog"

type Option func(*slog.HandlerOptions)

func WithLevel(lvl string) Option {
	return func(opts *slog.HandlerOptions) {
		switch lvl {
		case "debug":
			opts.Level = slog.LevelDebug
		case "error":
			opts.Level = slog.LevelError
		}
	}
}

func config(opts ...Option) slog.HandlerOptions {
	cfg := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}
