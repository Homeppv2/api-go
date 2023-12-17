package controller

import (
	"context"

	"homepp/api-go/internal/entity"
)

type (
	AuthService interface {
		Login(ctx context.Context, req entity.LoginRequest) (res entity.LoginResponse, err error)
		CheckSession(ctx context.Context, token string) (res entity.Session, err error)
		Logout(ctx context.Context, token string) (err error)
	}

	UserService interface {
		Register(ctx context.Context, req entity.RegisterUserRequest) (res entity.User, err error)
		GetByID(ctx context.Context, userID int64) (res entity.User, err error)
		GetByEmail(ctx context.Context, email string) (res entity.User, err error)
	}

	ControllerService interface {
		Create(ctx context.Context, hwKey string) (res entity.Controller, err error)
		GetByID(ctx context.Context, id int64) (res entity.Controller, err error)
		GetByHwKey(ctx context.Context, hwKey string) (res entity.Controller, err error)
		GetByIsUsedBy(ctx context.Context, isUsedBy int64) (res []entity.Controller, err error)
		UpdateIsUsed(ctx context.Context, req entity.UpdateControllerIsUsedByRequest) (res entity.Controller, err error)
		Delete(ctx context.Context, id int64) (err error)
	}
)
