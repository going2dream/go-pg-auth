package sql

import (
	"github.com/ZeroDayDrake/go-pg-auth/src/app/store"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	pool           *pgxpool.Pool
	userRepository *UserRepository
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{
		pool: pool,
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
