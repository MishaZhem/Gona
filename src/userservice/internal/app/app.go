package app

import (
	"context"
	"errors"
	"time"
	"userservice/internal/domain"
	"userservice/internal/repository"

	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Program struct {
	repo       repository.Repository
	jwtService *JWTService
	logger     *log.Logger
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

func NewApp(authRepository repository.Repository, jwtService *JWTService, logger *log.Logger) App {
	return &Program{
		repo:       authRepository,
		jwtService: jwtService,
		logger:     logger,
	}
}

func NewJWTService(secretKey string, ttl time.Duration) *JWTService {
	return &JWTService{secretKey: secretKey, ttl: ttl}
}

func (r *Program) Register(ctx context.Context, username, email, password string) error {
	r.logger.Infof("Trying to register user: %s", email)
	taken, err := r.repo.IsEmailTaken(ctx, email)
	if err != nil {
		r.logger.Warnf("Problem occurs when checking email: %s", email)
		return err
	}
	if taken {
		r.logger.Warnf("Email is already taken: %s", email)
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

	r.logger.Infof("User successfully registered: %s", email)
	return r.repo.CreateUser(ctx, user)
}

func (r *Program) Login(ctx context.Context, email, password string) (string, error) {
	r.logger.Infof("Trying to login user: %s", email)
	user, err := r.repo.GetUserByEmail(ctx, email)
	if err != nil {
		r.logger.Warnf("User not found: %s", email)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		r.logger.Warnf("Problem with comparing passwords: %s", email)
		return "", ErrInvalid
	}

	token, err := r.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		r.logger.Warnf("Failed to generate token: %s", email)
		return "", err
	}

	r.logger.Infof("User successfully logged in: %s", email)
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
