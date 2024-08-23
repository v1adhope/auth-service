package httpserver

import "time"

type Option func(*Config)

type Config struct {
	Socket          string
	ShutdownTimeout time.Duration
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
}

func WithSocket(socket string) Option {
	return func(cfg *Config) {
		cfg.Socket = socket
	}
}

func WithShutdownTimeout(st time.Duration) Option {
	return func(cfg *Config) {
		cfg.ShutdownTimeout = st
	}
}

func WithWriteTimeout(wt time.Duration) Option {
	return func(cfg *Config) {
		cfg.WriteTimeout = wt
	}
}

func WithReadTimeout(rt time.Duration) Option {
	return func(cfg *Config) {
		cfg.ReadTimeout = rt
	}
}

func config(opts ...Option) Config {
	cfg := Config{
		Socket: ":8080",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}
