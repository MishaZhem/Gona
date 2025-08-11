package app

import (
	"context"
	"errors"
	"io"
	"os"
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
	UploadAvatar(ctx context.Context, userID string, file io.Reader, fileSize int64, contentType string) (string, error)
	GetAvatarUrl(ctx context.Context, userID string) string
	RemoveAvatar(ctx context.Context, userID string) error
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	ChangeEmail(ctx context.Context, userID, newEmail string) error
	ChangeUsername(ctx context.Context, userID, newUsername string) error
	UpdateStatusWorker(ctx context.Context, userId string, worker bool) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	UpdatePassword(ctx context.Context, userID, newHashedPassword string) error
	UpdateUsername(ctx context.Context, userID, newUsername string) error
	UpdateEmail(ctx context.Context, userID, newEmail string) error
	UpdateUserAvatar(ctx context.Context, userID, avatarURL string) error
	UpdateStatusWork(ctx context.Context, userID string, newStatus bool) error
}

type StorageRepository interface {
	UploadFile(ctx context.Context, bucket string, objectName string, data io.Reader, size int64, contentType string) error
	DeleteFile(ctx context.Context, bucket string, objectName string) error
}

type JWTService struct {
	secretKey string
	ttl       time.Duration
}

var ErrEmailTaken = errors.New("email is already taken")
var ErrInvalid = errors.New("invalid email or password")

func NewApp(authRepository UserRepository, jwtService *JWTService, logger *log.Logger, minioClient StorageRepository, bucketAvatars string) App {
	return &Service{
		repo:          authRepository,
		jwtService:    jwtService,
		logger:        logger,
		minio:         minioClient,
		bucketAvatars: bucketAvatars,
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
		AvatarURL: "",
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

func (r *Service) UploadAvatar(ctx context.Context, userID string, file io.Reader, fileSize int64, contentType string) (string, error) {
	objectName := userID
	err := r.minio.UploadFile(ctx, r.bucketAvatars, objectName, file, fileSize, contentType)
	if err != nil {
		r.logger.Warnf("Failed to upload avatar: %s", err)
		return "", err
	}
	url := r.GetAvatarUrl(ctx, userID)
	r.logger.Infof("get avatar url: %v", url)
	err = r.repo.UpdateUserAvatar(ctx, userID, url)
	if err != nil {
		r.logger.Errorf("UpdateUserAvatar failed: %v", err)
		return "", nil
	}
	r.logger.Infof("upload user avatar url: %v", url)
	return url, nil
}

func (r *Service) GetAvatarUrl(ctx context.Context, userID string) string {
	host := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	objectName := userID
	return host + "/" + r.bucketAvatars + "/" + objectName
}

func (r *Service) RemoveAvatar(ctx context.Context, userID string) error {
	objectName := userID
	err := r.minio.DeleteFile(ctx, r.bucketAvatars, objectName)
	if err != nil {
		r.logger.Warnf("Failed to remove avatar: %s", err)
		return err
	}
	return nil
}

func (s *Service) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalid
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, userID, string(hashed))
}

func (s *Service) ChangeUsername(ctx context.Context, userID, newUsername string) error {
	return s.repo.UpdateUsername(ctx, userID, newUsername)
}

func (s *Service) ChangeEmail(ctx context.Context, userID, newEmail string) error {
	return s.repo.UpdateEmail(ctx, userID, newEmail)
}

func (s *Service) UpdateStatusWorker(ctx context.Context, userID string, worker bool) error {
	return s.repo.UpdateStatusWork(ctx, userID, worker)
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
