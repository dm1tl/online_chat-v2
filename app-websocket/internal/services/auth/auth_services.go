package auth

import (
	"app-websocket/internal/domain/auth"
	"context"
	"fmt"
)

type SSOProvider interface {
	Register(ctx context.Context, req auth.CreateUserReq) error
	Login(ctx context.Context, req auth.LoginReq) (*auth.LoginResp, error)
	Validate(ctx context.Context, req auth.ValidateTokenReq) (*auth.ValidateTokenResp, error)
	Delete(ctx context.Context, id int64) error
}

type AuthManager struct {
	ssoClient SSOProvider
}

func NewAuthManager(ssoCl SSOProvider) *AuthManager {
	return &AuthManager{
		ssoClient: ssoCl,
	}
}

func (u *AuthManager) Create(ctx context.Context, req auth.CreateUserReq) error {
	const op = "internal.services.Create()"
	err := u.ssoClient.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *AuthManager) Login(ctx context.Context, req auth.LoginReq) (auth.LoginResp, error) {
	var resp auth.LoginResp
	const op = "internal.services.Login()"
	ssoResp, err := u.ssoClient.Login(ctx, req)
	if err != nil {
		return resp, fmt.Errorf("%s: %w", op, err)
	}
	return auth.LoginResp{
		Token: ssoResp.Token,
	}, nil
}

func (u *AuthManager) Validate(ctx context.Context, req auth.ValidateTokenReq) (auth.ValidateTokenResp, error) {
	var resp auth.ValidateTokenResp
	const op = "internal.service.Validate()"
	ssoResp, err := u.ssoClient.Validate(ctx, req)
	if err != nil {
		return resp, fmt.Errorf("%s: %w", op, err)
	}
	return auth.ValidateTokenResp{
		ID: ssoResp.ID,
	}, nil
}
