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

func (a *App) userProfile(userId string) (string, error) {
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
