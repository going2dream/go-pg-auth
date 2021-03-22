package sql

import (
	"context"
	"github.com/ZeroDayDrake/go-pg-auth/src/http/models"
	"github.com/ZeroDayDrake/go-pg-auth/src/logger"
	"go.uber.org/zap"
)

var log = logger.New()

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Find(id string) (*models.User, error) {
	u := &models.User{}

	connection, err := r.store.pool.Acquire(context.Background())
	if err != nil {
		log.Info("Unable to acquire a database connection", zap.String("details", err.Error()))
		return nil, err
	}
	defer connection.Release()

	if err := connection.QueryRow(
		context.Background(),
		"SELECT id, login, password FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.Login, &u.Password); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByLogin(login string) (*models.User, error) {
	u := &models.User{}

	connection, err := r.store.pool.Acquire(context.Background())
	if err != nil {
		log.Info("Unable to acquire a database connection", zap.String("details", err.Error()))
		return nil, err
	}
	defer connection.Release()

	if err := connection.QueryRow(
		context.Background(),
		"SELECT id, login, password FROM users WHERE login = $1",
		login,
	).Scan(&u.ID, &u.Login, &u.Password); err != nil {
		return nil, err
	}

	return u, nil
}
