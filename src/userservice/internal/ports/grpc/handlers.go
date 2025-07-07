package grpc

import (
	context "context"
	"errors"

	"github.com/MishaZhem/Gona/src/userservice/internal/app"

	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedUserServiceServer
	app app.App
}

func Register(gRPC *grpc.Server, app app.App) {
	RegisterUserServiceServer(gRPC, &Server{app: app})
}

func (s *Server) Register(ctx context.Context, req *RegisterRequest) (*Empty, error) {
	err := s.app.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(getStatusByError(err), err.Error())
	}
	return &Empty{}, nil
}

func (s *Server) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	token, err := s.app.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &LoginResponse{Token: ""}, status.Error(getStatusByError(err), err.Error())
	}
	return &LoginResponse{Token: token}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	userId, err := s.app.ValidateToken(req.Token)
	if err != nil {
		return &TokenResponse{UserId: ""}, status.Error(getStatusByError(err), err.Error())
	}
	return &TokenResponse{UserId: userId}, nil
}

func getStatusByError(err error) codes.Code {
	switch {
	case errors.Is(err, app.ErrInvalid):
		return codes.AlreadyExists
	case errors.Is(err, app.ErrInvalid):
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}
