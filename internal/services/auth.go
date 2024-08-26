package services

import (
	"context"
	"fmt"
	"time"

	"github.com/v1adhope/auth-service/internal/models"
)

func (s *Services) GenerateTokenPair(ctx context.Context, userId string, ip string) (models.TokenPair, error) {
	if err := s.Validator.ValidateGuid(userId); err != nil {
		return models.TokenPair{}, err
	}

	tp, err := s.TokenManager.GeneratePair(ip, userId)
	if err != nil {
		return models.TokenPair{}, err
	}

	storeT, err := s.Hash.Do(tp.Refresh)
	if err != nil {
		return models.TokenPair{}, err
	}

	if err := s.AuthRepo.StoreToken(ctx, tp.Id, storeT, time.Now()); err != nil {
		return models.TokenPair{}, err
	}

	tp.Refresh = EncodeBase64(tp.Refresh)

	return tp, nil
}

func (s *Services) RefreshTokenPair(ctx context.Context, tp models.TokenPair, ip string) (models.TokenPair, error) {
	var err error

	tp.Refresh, err = DecodeBase64(tp.Refresh)
	if err != nil {
		return models.TokenPair{}, err
	}

	tp.Id, err = s.TokenManager.ExtractRefreshPayload(tp.Refresh)
	if err != nil {
		return models.TokenPair{}, err
	}

	storeT, err := s.AuthRepo.GetToken(ctx, tp.Id)
	if err != nil {
		return models.TokenPair{}, err
	}

	if err := s.Hash.Check(storeT, tp.Refresh); err != nil {
		return models.TokenPair{}, err
	}

	idAccessT, ipAccessT, userId, err := s.TokenManager.ExtractAccessPayload(tp.Access)
	if err != nil {
		return models.TokenPair{}, err
	}

	if idAccessT != tp.Id {
		return models.TokenPair{}, fmt.Errorf("services: auth: RefreshTokenPair: not equal ids: %w", models.ErrNotValidTokens)
	}

	if ip != ipAccessT {
		if err := s.Alert.Do("<SOME_EMAIL>", "<SOME_MSG>"); err != nil {
			return models.TokenPair{}, err
		}
	}

	if err := s.AuthRepo.DestroyToken(ctx, tp.Id); err != nil {
		return models.TokenPair{}, err
	}

	return s.GenerateTokenPair(ctx, userId, ip)
}
