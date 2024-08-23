package postgresql

type Option func(*Config)

type Config struct {
	ConnStr string
}

func WithConnStr(connStr string) Option {
	return func(cfg *Config) {
		cfg.ConnStr = connStr
	}
}

// INFO: panic if connStr not defined
func config(opts ...Option) Config {
	cfg := Config{}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.ConnStr == "" {
		panic("postgresql: define connection string")
	}

	return cfg
}
