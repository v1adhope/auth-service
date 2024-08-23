package services

import (
	"context"
	"fmt"

	"github.com/v1adhope/auth-service/internal/models"
)

type Auth struct {
	Validator    Validater
	TokenManager TokenManager
	Hash         Hasher
	AuthRepo     AuthRepo
	Allert       Allerter
}

func (s *Auth) GenerateTokenPair(ctx context.Context, id string, ip string) (models.TokenPair, error) {
	// TODO: might be ErrNotValidGuid
	if err := s.Validator.ValidateGuid(id); err != nil {
		return models.TokenPair{}, err
	}

	tp, err := s.TokenManager.GeneratePair(id, ip)
	if err != nil {
		return models.TokenPair{}, err
	}

	storeToken, err := s.makeAuthStoreToken(tp)
	if err != nil {
		return models.TokenPair{}, err
	}

	if err := s.AuthRepo.Store(ctx, storeToken); err != nil {
		return models.TokenPair{}, err
	}

	return tp, nil
}

func (s *Auth) RefreshTokenPair(ctx context.Context, tp models.TokenPair, ip string) (models.TokenPair, error) {
	storeToken, err := s.makeAuthStoreToken(tp)
	if err != nil {
		return models.TokenPair{}, err
	}

	// TODO: might be internal or ErrNotValidTokens
	if err := s.AuthRepo.Check(ctx, storeToken); err != nil {
		return models.TokenPair{}, err
	}

	// TODO: might be ErrNotValidTokens
	newTp, isIpChanged, err := s.TokenManager.RefreshPair(tp, ip)
	if err != nil {
		return models.TokenPair{}, err
	}

	// INFO: should be req getEmailByUserId
	if isIpChanged {
		if err := s.Allert.Do("<SOME_EMAIL>", "<SOME_MSG>"); err != nil {
			return models.TokenPair{}, err
		}
	}

	if err := s.AuthRepo.Destroy(ctx, storeToken); err != nil {
		return models.TokenPair{}, err
	}

	return newTp, nil
}

func (s *Auth) makeAuthStoreToken(tp models.TokenPair) (string, error) {
	str := fmt.Sprintf("%s:%s", tp.Access, tp.Refresh)

	storeToken, err := s.Hash.Do(str)
	if err != nil {
		return "", err
	}

	return storeToken, nil
}
