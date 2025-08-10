package postgres

import (
	"github.com/MishaZhem/Gona/src/userservice/internal/adapters/postgres/queries"
	"github.com/MishaZhem/Gona/src/userservice/internal/app"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type repo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger log.FieldLogger
}

func NewUserRepository(pgxPool *pgxpool.Pool, logger log.FieldLogger) app.UserRepository {
	return &repo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
		logger:  logger,
	}
}
