package models

import "errors"

var (
	ErrNotValidGuid   = errors.New("Not valid guid")
	ErrNotValidTokens = errors.New("Not valid token(s)")
)
