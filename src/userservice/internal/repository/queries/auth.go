package queries

import (
	"context"
	"errors"
	"fmt"
	"userservice/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

const createUserQuery = `
    INSERT INTO users (id, email, username, password_hash, created_at)
    VALUES ($1, $2, $3, $4, $5)`

func (q *Queries) CreateUser(ctx context.Context, user *domain.User) error {

	if _, err := q.pool.Exec(ctx, createUserQuery, uuid.Must(uuid.NewRandom()), user.Email, user.Username, user.Password, user.CreatedAt); err != nil {
		return err
	}
	return nil
}

const getUserByEmailQuery = `SELECT id, email, username, password_hash, created_at
	FROM users
	WHERE email = $1`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := q.pool.QueryRow(ctx, getUserByEmailQuery, email)
	fmt.Print(row)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

const isEmailTakenQuery = `
	SELECT EXISTS (
		SELECT 1 FROM users WHERE email = $1
	)
`

func (q *Queries) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := q.pool.QueryRow(ctx, isEmailTakenQuery, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
