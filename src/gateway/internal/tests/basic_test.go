package tests

import (
	"context"
	"testing"
	"time"

	"github.com/MishaZhem/Gona/src/gateway/genproto/userpb"
	"github.com/MishaZhem/Gona/src/gateway/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// / MockUserServiceClient mocks user.UserServiceClient
type MockUserServiceClient struct {
	mock.Mock
	userpb.UserServiceClient
}

func (m *MockUserServiceClient) Register(ctx context.Context, in *userpb.RegisterRequest, opts ...grpc.CallOption) (*userpb.Empty, error) {
	args := m.Called(ctx, in)
	return &userpb.Empty{}, args.Error(1)
}

func (m *MockUserServiceClient) Login(ctx context.Context, in *userpb.LoginRequest, opts ...grpc.CallOption) (*userpb.LoginResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userpb.LoginResponse), args.Error(1)
}

func (m *MockUserServiceClient) ValidateToken(ctx context.Context, in *userpb.TokenRequest, opts ...grpc.CallOption) (*userpb.TokenResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*userpb.TokenResponse), args.Error(1)
}

func TestRegister(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	a := app.NewApp(mockClient, 2*time.Second)

	req := &userpb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "pass",
	}

	mockClient.On("Register", mock.Anything, req).Return(&userpb.Empty{}, nil)

	err := a.Register("testuser", "test@example.com", "pass")
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	a := app.NewApp(mockClient, 2*time.Second)

	req := &userpb.LoginRequest{
		Email:    "test@example.com",
		Password: "pass",
	}
	mockResp := &userpb.LoginResponse{Token: "token123"}

	mockClient.On("Login", mock.Anything, req).Return(mockResp, nil)

	token, err := a.Login("test@example.com", "pass")
	assert.NoError(t, err)
	assert.Equal(t, "token123", token)
	mockClient.AssertExpectations(t)
}

func TestValidateToken(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	a := app.NewApp(mockClient, 2*time.Second)

	req := &userpb.TokenRequest{
		Token: "token123",
	}
	mockResp := &userpb.TokenResponse{UserId: "user-id"}

	mockClient.On("ValidateToken", mock.Anything, req).Return(mockResp, nil)

	userID, err := a.ValidateToken("token123")
	assert.NoError(t, err)
	assert.Equal(t, "user-id", userID)
	mockClient.AssertExpectations(t)
}
