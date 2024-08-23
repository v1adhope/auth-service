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

func (t *Tokens) GeneratePair(userId string, ip string) (models.TokenPair, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("tokens: tokens: GeneratePair: NewV6: %w", err)
	}

	accessT, err := t.generateAccess(userId, ip, id.String())
	if err != nil {
		return models.TokenPair{}, err
	}

	refreshT, err := t.generateRefresh(id.String())
	if err != nil {
		return models.TokenPair{}, err
	}

	return models.TokenPair{id.String(), accessT, refreshT}, nil
}

func (t *Tokens) generateAccess(userId string, ip string, id string) (string, error) {
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
		return "", models.ErrNotValidTokens
	}

	return string(text), nil
}

func (t *Tokens) ExtractAccessPayload(token string) (userId, id, ip string, err error) {
	claims, err := t.parseAccess(token)
	if err != nil {
		return "", "", "", err
	}

	userId, id, ip = t.extractUsefulClaims(claims)

	return userId, id, ip, nil
}

//	func (t *Tokens) RefreshPair(tp models.TokenPair, ip string) (newTp models.TokenPair, isIpChanged bool, err error) {
//		claims, err := t.parseAccess(tp.Access)
//		if err != nil {
//			return models.TokenPair{}, false, err
//		}
//
//		refreshT, err := t.parseRefresh(tp.Refresh)
//		if err != nil {
//			return models.TokenPair{}, false, err
//		}
//
//		id, tokenIp, tpId := t.extractUsefulClaims(claims)
//		if refreshT != tpId {
//			return models.TokenPair{}, false, models.ErrNotValidTokens
//		}
//
//		if tokenIp != ip {
//			isIpChanged = true
//		}
//
//		newTp, err = t.GeneratePair(id, ip)
//		if err != nil {
//			return models.TokenPair{}, false, err
//		}
//
//		return newTp, false, nil
//	}
func (t *Tokens) parseAccess(target string) (accessClaims, error) {
	accessT, err := jwt.ParseWithClaims(target, &accessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrNotValidTokens
		}

		return []byte(t.access.key), nil
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

func (t *Tokens) extractUsefulClaims(claims accessClaims) (userId, id, ip string) {
	return claims.Subject, claims.ID, claims.Ip
}

//
// func (t *Tokens) parseRefresh(target string) (string, error) {
// 	text, err := serialization.DecryptByGcm([]byte(target), t.refresh.key)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return string(text), nil
// }
