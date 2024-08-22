package services

import "github.com/v1adhope/auth-service/internal/models"

type Allerter interface {
	Do(email, msg string) error
}

type AuthRepo interface {
	Store(token string) error
	Check(token string) error
	Destroy(token string) error
}

type Hasher interface {
	Do(target string) (string, error)
}

type TokenManager interface {
	GeneratePair(id string, ip string) (models.TokenPair, error)
	RefreshPair(tp models.TokenPair) (newTp models.TokenPair, isIpChanged bool, err error)
}

type Validater interface {
	ValideGuid(id string) error
}
