package grpc

import (
	context "context"
	"errors"

	"github.com/MishaZhem/Gona/src/userservice/genproto/userpb"
	"github.com/MishaZhem/Gona/src/userservice/internal/app"

	"google.golang.org/grpc"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
	app app.App
}

func Register(gRPC *grpc.Server, app app.App) {
	userpb.RegisterUserServiceServer(gRPC, &Server{app: app})
}

func (s *Server) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.Empty, error) {
	err := s.app.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(getStatusByError(err), err.Error())
	}
	return &userpb.Empty{}, nil
}

func (s *Server) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	token, err := s.app.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &userpb.LoginResponse{Token: ""}, status.Error(getStatusByError(err), err.Error())
	}
	return &userpb.LoginResponse{Token: token}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *userpb.TokenRequest) (*userpb.TokenResponse, error) {
	userId, err := s.app.ValidateToken(req.Token)
	if err != nil {
		return &userpb.TokenResponse{UserId: ""}, status.Error(getStatusByError(err), err.Error())
	}
	return &userpb.TokenResponse{UserId: userId}, nil
}

func (s *Server) Profile(ctx context.Context, req *userpb.ProfileRequest) (*userpb.ProfileResponse, error) {
	user, err := s.app.Profile(ctx, req.UserId)
	if err != nil {
		return &userpb.ProfileResponse{UserId: ""}, status.Error(getStatusByError(err), err.Error())
	}
	return &userpb.ProfileResponse{UserId: user.ID.String(), Username: user.Username, Email: user.Email}, nil
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
