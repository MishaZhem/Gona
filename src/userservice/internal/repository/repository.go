package repository

import (
	"context"

	"github.com/MishaZhem/Gona/src/userservice/internal/domain"
	"github.com/MishaZhem/Gona/src/userservice/internal/repository/queries"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type repo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger log.FieldLogger
}

func NewRepository(pgxPool *pgxpool.Pool, logger log.FieldLogger) Repository {
	return &repo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
		logger:  logger,
	}
}

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}
