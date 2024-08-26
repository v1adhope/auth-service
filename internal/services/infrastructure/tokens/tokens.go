package tokens

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/v1adhope/auth-service/internal/models"
	"github.com/v1adhope/auth-service/pkg/serialization"
)

type Tokens struct {
	access  Acccess
	refresh Refresh
	issuer  string
}

type Acccess struct {
	ttl time.Duration
	key []byte
}

type Refresh struct {
	key []byte
}

type accessClaims struct {
	Ip string `json:"ip"`
	jwt.RegisteredClaims
}

// INFO: panic if access or refresh keys not defined
func New(opts ...Option) *Tokens {
	t := &Tokens{
		access: Acccess{
			ttl: time.Duration(20 * time.Second),
		},
		refresh: Refresh{},
		issuer:  "auth-service",
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.access.key == nil {
		panic("tokens: define access key")
	}

	if t.refresh.key == nil {
		panic("tokens: define refresh key")
	}

	return t
}

// INFO: not invariant values might be used as deps for testing
func (t *Tokens) GeneratePair(ip string, userId string) (models.TokenPair, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("tokens: tokens: GeneratePair: NewV6: %w", err)
	}

	accessT, err := t.generateAccess(id.String(), ip, userId)
	if err != nil {
		return models.TokenPair{}, err
	}

	refreshT, err := t.generateRefresh(id.String())
	if err != nil {
		return models.TokenPair{}, err
	}

	return models.TokenPair{id.String(), accessT, refreshT}, nil
}

func (t *Tokens) generateAccess(id string, ip string, userId string) (string, error) {
	claims := accessClaims{
		Ip: ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.access.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    t.issuer,
			Subject:   userId,
			ID:        id,
		},
	}

	accessT, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(t.access.key)
	if err != nil {
		return "", fmt.Errorf("tokens: tokens: generateAccess: SignedString: %w", err)
	}

	return accessT, nil
}

func (t *Tokens) generateRefresh(id string) (string, error) {
	ciphertext, err := serialization.EncryptByGcm([]byte(id), t.refresh.key)
	if err != nil {
		return "", err
	}

	return string(ciphertext), err
}

func (t *Tokens) ExtractRefreshPayload(token string) (string, error) {
	text, err := serialization.DecryptByGcm([]byte(token), t.refresh.key)
	if err != nil {
		return "", fmt.Errorf("tokens: tokens: ExtractRefreshPayload: DecryptByGcm: %w", models.ErrNotValidTokens)
	}

	return string(text), nil
}

func (t *Tokens) ExtractAccessPayload(token string) (userId, id, ip string, err error) {
	claims, err := t.parseAccess(token)
	if err != nil {
		return "", "", "", err
	}

	id, ip, userId = t.extractUsefulClaims(claims)

	return id, ip, userId, nil
}

func (t *Tokens) parseAccess(target string) (accessClaims, error) {
	accessT, err := jwt.ParseWithClaims(target, &accessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("tokens: tokens: parseAccess: ParseWithClaims: %w", models.ErrNotValidTokens)
		}

		return []byte(t.access.key), nil
	})
	if err != nil {
		return accessClaims{}, fmt.Errorf("tokens: tokens: parseAccess: Parse: %w", err)
	}

	claims, ok := accessT.Claims.(*accessClaims)
	if !ok {
		return accessClaims{}, fmt.Errorf("tokens: tokens: parseAccess: Claims: %w", models.ErrNotValidTokens)
	}

	return *claims, nil
}

func (t *Tokens) extractUsefulClaims(claims accessClaims) (id, ip, userId string) {
	return claims.ID, claims.Ip, claims.Subject
}
