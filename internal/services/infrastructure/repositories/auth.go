package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/v1adhope/auth-service/internal/models"
)

func (r *Repos) StoreToken(ctx context.Context, id, token string, now time.Time) error {
	sql, args, err := r.Builder.Insert("auth_whitelist").
		SetMap(squirrel.Eq{
			"id":         id,
			"created_at": now,
			"token":      token,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("repositories: auth: Store: ToSql: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("repositories: auth: Store: Exec: %w", err)
	}

	return nil
}

func (r *Repos) GetToken(ctx context.Context, id string) (string, error) {
	sql, args, err := r.Builder.Select("token").
		From("auth_whitelist").
		Where(squirrel.Eq{
			"id": id,
		}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("repositories: auth: Get: ToSql: %w", err)
	}

	token := ""

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&token); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrNotValidTokens
		}

		return "", fmt.Errorf("repositories: auth: Get: ToSql: %w", err)
	}

	return token, nil
}

func (r *Repos) DestroyToken(ctx context.Context, id string) error {
	sql, args, err := r.Builder.Delete("auth_whitelist").
		Where(squirrel.Eq{
			"id": id,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("repositories: auth: Destroy: ToSql: %w", err)
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("repositories: auth: Destroy: Exec: %w", err)
	}

	return nil
}
