package app_test

import (
	"context"
	"testing"
	"time"
	"userservice/internal/app"
	"userservice/internal/domain"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func TestRegister(t *testing.T) {
	mockRepo := new(MockRepo)
	jwtService := app.NewJWTService("secret", time.Hour)
	logger := log.New()
	testApp := app.NewApp(mockRepo, jwtService, logger)

	email := "test@example.com"
	username := "testuser"
	password := "password123"

	mockRepo.On("IsEmailTaken", mock.Anything, email).Return(false, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	err := testApp.Register(context.Background(), username, email, password)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockRepo)
	jwtService := app.NewJWTService("test-secret", time.Hour)
	logger := log.New()
	testApp := app.NewApp(mockRepo, jwtService, logger)

	email := "test@example.com"
	password := "password123"
	hashed, _ := HashPassword(password) // helper

	user := &domain.User{
		ID:        uuid.New(),
		Email:     email,
		Username:  "testuser",
		Password:  hashed,
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)

	token, err := testApp.Login(context.Background(), email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestValidateToken(t *testing.T) {
	jwtService := app.NewJWTService("secret", time.Minute*10)
	token, err := jwtService.GenerateToken("123", "test@example.com")
	assert.NoError(t, err)

	userID, err := jwtService.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "123", userID)
}
