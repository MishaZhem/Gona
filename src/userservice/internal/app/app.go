package app

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/MishaZhem/Gona/src/userservice/internal/domain"

	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo          UserRepository
	jwtService    *JWTService
	logger        *log.Logger
	minio         StorageRepository
	bucketAvatars string
}

type App interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(token string) (string, error)
	Profile(ctx context.Context, userId string) (*domain.User, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}

type StorageRepository interface {
	UploadFile(ctx context.Context, bucket string, objectName string, data io.Reader, size int64) error
	GetFileURL(ctx context.Context, bucket string, objectName string) (string, error)
	DeleteFile(ctx context.Context, bucket string, objectName string) error
}

type JWTService struct {
	secretKey string
	ttl       time.Duration
}

var ErrEmailTaken = errors.New("email is already taken")
var ErrInvalid = errors.New("invalid email or password")

func NewApp(authRepository UserRepository, jwtService *JWTService, logger *log.Logger, minioClient StorageRepository) App {
	return &Service{
		repo:       authRepository,
		jwtService: jwtService,
		logger:     logger,
		minio:      minioClient,
	}
}

func NewJWTService(secretKey string, ttl time.Duration) *JWTService {
	return &JWTService{secretKey: secretKey, ttl: ttl}
}

func (r *Service) Register(ctx context.Context, username, email, password string) error {
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

func (r *Service) Login(ctx context.Context, email, password string) (string, error) {
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

func (r *Service) Profile(ctx context.Context, userId string) (*domain.User, error) {
	r.logger.Infof("Trying to take profile of user: %s", userId)
	user, err := r.repo.GetUserById(ctx, userId)
	if err != nil {
		r.logger.Warnf("User not found: %s", userId)
		return nil, err
	}
	r.logger.Infof("User successfully got profile: %s", userId)
	return user, nil
}

func (r *Service) ValidateToken(token string) (string, error) {
	return r.jwtService.ValidateToken(token)
}

func (r *Service) UploadAvatar(ctx context.Context, userID string, file io.Reader, fileSize int64) (string, error) {
	objectName := "avatar_" + userID
	err := r.minio.UploadFile(ctx, r.bucketAvatars, objectName, file, fileSize)
	if err != nil {
		r.logger.Warnf("Failed to upload avatar: %s", err)
		return "", err
	}
	url, err := r.GetAvatarPresignedUrl(ctx, userID)
	if err != nil {
		r.logger.Warnf("Failed to get avatar: %s", err)
		return "", err
	}
	return url, nil
}

func (r *Service) GetAvatarPresignedUrl(ctx context.Context, userID string) (string, error) {
	objectName := "avatar_" + userID
	url, err := r.minio.GetFileURL(ctx, r.bucketAvatars, objectName)
	if err != nil {
		r.logger.Warnf("Failed to get avatar: %s", err)
		return "", err
	}
	return url, nil
}

func (r *Service) RemoveAvatar(ctx context.Context, userID string) error {
	objectName := "avatar_" + userID
	err := r.minio.DeleteFile(ctx, r.bucketAvatars, objectName)
	if err != nil {
		r.logger.Warnf("Failed to remove avatar: %s", err)
		return err
	}
	return nil
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
