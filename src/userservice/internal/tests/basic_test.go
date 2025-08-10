package app_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/MishaZhem/Gona/src/userservice/internal/domain"

	"github.com/MishaZhem/Gona/src/userservice/internal/app"

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

func (m *MockRepo) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepo) UpdatePassword(ctx context.Context, userID, newHashedPassword string) error {
	args := m.Called(ctx, userID, newHashedPassword)
	return args.Error(0)
}

func (m *MockRepo) UpdateUsername(ctx context.Context, userID, newUsername string) error {
	args := m.Called(ctx, userID, newUsername)
	return args.Error(0)
}

func (m *MockRepo) UpdateEmail(ctx context.Context, userID, newEmail string) error {
	args := m.Called(ctx, userID, newEmail)
	return args.Error(0)
}

func (m *MockRepo) UpdateUserAvatar(ctx context.Context, userID, avatarURL string) error {
	args := m.Called(ctx, userID, avatarURL)
	return args.Error(0)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) UploadFile(ctx context.Context, bucket, objectName string, data io.Reader, size int64, contentType string) error {
	args := m.Called(ctx, bucket, objectName, data, size, contentType)
	return args.Error(0)
}

func (m *MockStorage) GetFileURL(ctx context.Context, bucket, objectId string) (string, error) {
	args := m.Called(ctx, bucket, objectId)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) DeleteFile(ctx context.Context, bucket, objectId string) error {
	args := m.Called(ctx, bucket, objectId)
	return args.Error(0)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func TestRegister(t *testing.T) {
	mockRepo := new(MockRepo)
	mockStorage := new(MockStorage)
	jwtService := app.NewJWTService("secret", time.Hour)
	logger := log.New()
	testApp := app.NewApp(mockRepo, jwtService, logger, mockStorage, "test-bucket")

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
	mockStorage := new(MockStorage)
	jwtService := app.NewJWTService("test-secret", time.Hour)
	logger := log.New()
	testApp := app.NewApp(mockRepo, jwtService, logger, mockStorage, "test-bucket")

	email := "test@example.com"
	password := "password123"
	hashed, _ := HashPassword(password)

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

func TestUploadAvatar(t *testing.T) {
	t.Setenv("MINIO_PUBLIC_ENDPOINT", "http://localhost")

	mockRepo := new(MockRepo)
	mockStorage := new(MockStorage)
	logger := log.New()
	jwtService := app.NewJWTService("secret", time.Hour)

	testApp := app.NewApp(mockRepo, jwtService, logger, mockStorage, "test-bucket")

	data := bytes.NewReader([]byte("image-bytes"))
	userID := "user123"
	fileSize := int64(data.Len())
	contentType := "image/jpeg"

	expectedObject := userID
	expectedURL := "http://localhost/test-bucket/" + userID

	mockStorage.
		On("UploadFile", mock.Anything, "test-bucket", expectedObject, mock.Anything, fileSize, contentType).
		Return(nil)

	mockRepo.
		On("UpdateUserAvatar", mock.Anything, userID, expectedURL).
		Return(nil)

	url, err := testApp.UploadAvatar(context.Background(), userID, data, fileSize, contentType)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL, url)

	mockStorage.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestDeleteAvatar(t *testing.T) {
	mockRepo := new(MockRepo)
	mockStorage := new(MockStorage)
	jwtService := app.NewJWTService("secret", time.Hour)
	logger := log.New()

	testApp := app.NewApp(mockRepo, jwtService, logger, mockStorage, "test-bucket")

	objectId := "avatars/user123"

	mockStorage.On("DeleteFile", mock.Anything, "test-bucket", objectId).Return(nil)

	err := testApp.RemoveAvatar(context.Background(), objectId)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
}
