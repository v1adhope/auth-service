package tokens

import "time"

type Option func(*Tokens)

func WithAccessKey(k string) Option {
	return func(t *Tokens) {
		t.access.key = []byte(k)
	}
}

func WithAccessTtl(ttl time.Duration) Option {
	return func(t *Tokens) {
		t.access.ttl = ttl
	}
}

func WithRefreshKey(k string) Option {
	return func(t *Tokens) {
		t.refresh.key = []byte(k)
	}
}

func WithIssuer(i string) Option {
	return func(t *Tokens) {
		t.issuer = i
	}
}
