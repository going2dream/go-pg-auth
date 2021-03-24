package pgsql

import (
	"github.com/going2dream/go-pg-auth/src/app/logger"
	"github.com/going2dream/go-pg-auth/src/app/store"
	"github.com/jackc/pgx/v4/pgxpool"
)

var log = logger.New()

type Store struct {
	pool           *pgxpool.Pool
	userRepository *UserRepository
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
