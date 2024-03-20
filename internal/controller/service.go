package controller

import (
	"context"

	"api-go/internal/entity"
)

type (
	UserService interface {
		Register(ctx context.Context, req entity.User) (res entity.User, err error)
		GetByID(ctx context.Context, userID int) (res entity.User, err error)
		GetByEmail(ctx context.Context, email string) (res entity.User, err error)
	}

	ControllerService interface {
		GetCountMessangesFromIdForUserId(ctx context.Context, count int, from int, userID int) (msg []entity.MessangeTypeZiroJson, err error)
		GetControllersByUserId(ctx context.Context, user_id int) (controllers []entity.ControllersData, err error)
	}
)
