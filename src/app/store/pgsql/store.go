package pgsql

import (
	"github.com/going2dream/go-pg-auth/src/app/logger"
	"github.com/going2dream/go-pg-auth/src/app/store"
	"github.com/jackc/pgx/v4/pgxpool"
)

var log = logger.New()

type Store struct {
	pool                   *pgxpool.Pool
	userRepository         *UserRepository
	refreshTokenRepository *RefreshTokenRepository
}

func NewStore() *Store {
	return &Store{
		pool: NewPoolInstance(),
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) RefreshToken() store.RefreshTokenRepository {
	if s.refreshTokenRepository != nil {
		return s.refreshTokenRepository
	}

	s.refreshTokenRepository = &RefreshTokenRepository{
		store: s,
	}

	return s.refreshTokenRepository
}
