package queries

import (
	"context"
	"errors"

	"github.com/MishaZhem/Gona/src/userservice/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

const createUserQuery = `
    INSERT INTO users (id, email, username, password_hash, created_at, avatar)
    VALUES ($1, $2, $3, $4, $5, $6)`

func (q *Queries) CreateUser(ctx context.Context, user *domain.User) error {

	if _, err := q.pool.Exec(ctx, createUserQuery, uuid.Must(uuid.NewRandom()), user.Email, user.Username, user.Password, user.CreatedAt, ""); err != nil {
		return err
	}
	return nil
}

const getUserByEmailQuery = `SELECT id, email, username, password_hash, created_at, avatar
	FROM users
	WHERE email = $1`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := q.pool.QueryRow(ctx, getUserByEmailQuery, email)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.AvatarURL)
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

const getUserByIdQuery = `SELECT id, email, username, password_hash, created_at, avatar
	FROM users
	WHERE id = $1`

func (q *Queries) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	row := q.pool.QueryRow(ctx, getUserByIdQuery, id)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.AvatarURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

const updateUserAvatarQuery = `
	UPDATE users
	SET avatar = $1
	WHERE id = $2
`

func (q *Queries) UpdateUserAvatar(ctx context.Context, userID, avatarURL string) error {
	_, err := q.pool.Exec(ctx, updateUserAvatarQuery, avatarURL, userID)
	return err
}

const updatePasswordQuery = `
	UPDATE users
	SET password_hash = $1
	WHERE id = $2
`

func (q *Queries) UpdatePassword(ctx context.Context, userID, newHashedPassword string) error {
	_, err := q.pool.Exec(ctx, updatePasswordQuery, newHashedPassword, userID)
	return err
}

const updateUsernameQuery = `
	UPDATE users
	SET password_hash = $1
	WHERE id = $2
`

func (q *Queries) UpdateUsername(ctx context.Context, userID, newUsername string) error {
	_, err := q.pool.Exec(ctx, updateUsernameQuery, newUsername, userID)
	return err
}

const updateEmailQuery = `
	UPDATE users
	SET password_hash = $1
	WHERE id = $2
`

func (q *Queries) UpdateEmail(ctx context.Context, userID, newEmail string) error {
	_, err := q.pool.Exec(ctx, updateEmailQuery, newEmail, userID)
	return err
}
