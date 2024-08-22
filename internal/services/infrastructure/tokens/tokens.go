package tokens

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/v1adhope/auth-service/internal/models"
)

type Config struct {
	Access  Acccess
	Refresh Refresh
	Issuer  string
}

type Tokens struct {
	access  Acccess
	refresh Refresh
	issuer  string
}

type Acccess struct {
	Ttl int
	Key string
}

type Refresh struct {
	Key string
}

type accessClaims struct {
	Ip string `json:"ip"`
	jwt.RegisteredClaims
}

func New(cfg Config) *Tokens {
	return &Tokens{
		access: Acccess{
			Ttl: cfg.Access.Ttl,
			Key: cfg.Access.Key,
		},
		refresh: Refresh{
			Key: cfg.Refresh.Key,
		},
		issuer: cfg.Issuer,
	}
}

func (t *Tokens) GeneratePair(id string, ip string) (models.TokenPair, error) {
	accessT, err := t.generateAccess(id, ip)
	if err != nil {
		return models.TokenPair{}, err
	}

	refreshT, err := t.generateRefresh()
	if err != nil {
		return models.TokenPair{}, err
	}

	return models.TokenPair{accessT, refreshT}, nil
}

func (t *Tokens) generateAccess(id string, ip string) (string, error) {
	claims := accessClaims{
		Ip: ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.access.Ttl) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    t.issuer,
			Subject:   id,
		},
	}

	accessT, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(t.access.Key)
	if err != nil {
		return "", fmt.Errorf("tokens: tokens: generateAccess: SignedString: %w", err)
	}

	return accessT, nil
}

// INFO: NotBefore brute force protection
func (t *Tokens) generateRefresh() (string, error) {
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		// NotBefore: jwt.NewNumericDate(time.Now().Add(time.Duration(t.access.ttl) * time.Second)),
		Issuer: t.issuer,
	}

	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims).SignedString(t.refresh.Key)
	if err != nil {
		return "", fmt.Errorf("tokens: tokens: generateRefresh: SignedString: %w", err)
	}

	refreshT := base64.StdEncoding.EncodeToString([]byte(ss))

	return refreshT, nil
}

func (t *Tokens) RefreshPair(tp models.TokenPair, ip string) (newTp models.TokenPair, isIpChanged bool, err error) {
	claims, err := t.ParseAccess(tp.Access)
	if err != nil {
		return models.TokenPair{}, false, err
	}

	err = t.ParseRefresh(tp.Refresh)
	if err != nil {
		return models.TokenPair{}, false, err
	}

	id, tokenIp := t.extractUsefulClaims(claims)
	if tokenIp != ip {
		isIpChanged = true
	}

	newTp, err = t.GeneratePair(id, ip)
	if err != nil {
		return models.TokenPair{}, false, err
	}

	return newTp, false, nil
}

func (t *Tokens) ParseAccess(target string) (accessClaims, error) {
	accessT, err := jwt.Parse(target, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrNotValidTokens
		}

		return []byte(t.access.Key), nil
	})
	if err != nil {
		return accessClaims{}, fmt.Errorf("tokens: tokens: ParseAccess: Parse: %w", err)
	}

	claims, ok := accessT.Claims.(*accessClaims)
	if !ok {
		return accessClaims{}, models.ErrNotValidTokens
	}

	return *claims, nil
}

func (t *Tokens) extractUsefulClaims(claims accessClaims) (id, ip string) {
	return claims.Subject, claims.Ip
}

func (t *Tokens) ParseRefresh(target string) error {
	decodeTarget, err := base64.StdEncoding.DecodeString(target)
	if err != nil {
		return models.ErrNotValidTokens
	}

	_, err = jwt.Parse(string(decodeTarget), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrNotValidTokens
		}

		return []byte(t.refresh.Key), nil
	})
	if err != nil {
		return fmt.Errorf("tokens: tokens: ParseRefresh: Parse: %w", err)
	}

	return nil
}
