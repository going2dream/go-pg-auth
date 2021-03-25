package pgsql

import (
	"context"
	"github.com/going2dream/go-pg-auth/src/app/models"
	"go.uber.org/zap"
)

type RefreshTokenRepository struct {
	store *Store
}

func (r *RefreshTokenRepository) Find(id string) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{}

	connection, err := r.store.pool.Acquire(context.Background())
	if err != nil {
		log.Info("Unable to acquire a database connection", zap.String("details", err.Error()))
		return nil, err
	}
	defer connection.Release()

	if err := connection.QueryRow(
		context.Background(),
		"SELECT id, user_id, ua, fingerprint, ip, expires_in, created_at, updated_at FROM refresh_tokens WHERE id = $1",
		id,
	).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.UA,
		&refreshToken.Fingerprint,
		&refreshToken.IP,
		&refreshToken.ExpiresIn,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) Create(rt *models.RefreshToken) error {
	connection, err := r.store.pool.Acquire(context.Background())
	if err != nil {
		log.Info("Unable to acquire a database connection", zap.String("details", err.Error()))
		return err
	}
	defer connection.Release()

	if err := connection.QueryRow(
		context.Background(),
		"INSERT INTO refresh_tokens (user_id, ua, fingerprint, ip, expires_in, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		rt.UserID,
		rt.UA,
		rt.Fingerprint,
		rt.IP,
		rt.ExpiresIn,
		rt.CreatedAt,
		rt.UpdatedAt,
	).Scan(&rt.ID); err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) Delete(id string) error {
	connection, err := r.store.pool.Acquire(context.Background())
	if err != nil {
		log.Info("Unable to acquire a database connection", zap.String("details", err.Error()))
		return err
	}
	defer connection.Release()

	if _, err := connection.Query(
		context.Background(),
		"DELETE FROM refresh_tokens WHERE id = $1",
		id,
	); err != nil {
		return err
	}

	return nil
}
