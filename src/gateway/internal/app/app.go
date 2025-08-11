package app

import (
	"context"
	"time"

	"github.com/MishaZhem/Gona/src/gateway/genproto/userpb"
)

type App struct {
	userClient userpb.UserServiceClient
	timeout    time.Duration
}

func NewApp(userClient userpb.UserServiceClient, timeout time.Duration) App {
	return App{
		userClient: userClient,
		timeout:    timeout,
	}
}

func (a *App) Register(username, email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	_, err := a.userClient.Register(ctx, &userpb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})

	if err != nil {
		return err
	}
	return nil
}

func (a *App) Login(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	resp, err := a.userClient.Login(ctx, &userpb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (a *App) ValidateToken(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	resp, err := a.userClient.ValidateToken(ctx, &userpb.TokenRequest{
		Token: token,
	})
	if err != nil {
		return "", err
	}
	return resp.UserId, nil
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Worker    bool   `json:"worker"`
}

func (a *App) UserProfile(userId string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	resp, err := a.userClient.Profile(ctx, &userpb.ProfileRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        resp.UserId,
		Email:     resp.Email,
		Username:  resp.Email,
		AvatarURL: resp.Avatar,
		Worker:    resp.Worker,
	}, nil
}

func (a *App) UpdateStatusWorker(userId string, worker bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	_, err := a.userClient.UpdateStatusWork(ctx, &userpb.UpdateStatusWorkRequest{
		UserId: userId,
		Worker: worker,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *App) ChangeAvatar(userId string, fileBytes []byte, fileSize int64, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	resp, err := a.userClient.UploadAvatar(ctx, &userpb.UploadAvatarRequest{
		UserId:      userId,
		AvatarData:  fileBytes,
		FileSize:    fileSize,
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	return resp.Url, nil
}
