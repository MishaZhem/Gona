package app

import (
	"context"
	"errors"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Program struct {
	repo       repository.Repository
	jwtService *JWTService
}

type App interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}

type JWTService struct {
	secretKey string
	ttl       time.Duration
}

var ErrEmailTaken = errors.New("email is already taken")
var ErrInvalid = errors.New("invalid email or password")

func NewApp(authRepository repository.Repository, jwtService *JWTService) App {
	return &Program{
		repo:       authRepository,
		jwtService: jwtService,
	}
}

func NewJWTService(secretKey string, ttl time.Duration) *JWTService {
	return &JWTService{secretKey: secretKey, ttl: ttl}
}

func (r *Program) Register(ctx context.Context, username, email, password string) error {
	taken, err := r.repo.IsEmailTaken(ctx, email)
	if err != nil {
		return err
	}
	if taken {
		return ErrEmailTaken
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:        uuid.Must(uuid.NewRandom()),
		Email:     email,
		Username:  username,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
	}

	return r.repo.CreateUser(ctx, user)
}

func (s *Program) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalid
	}

	token, err := s.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *JWTService) GenerateToken(userID string, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(j.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string), nil
}
