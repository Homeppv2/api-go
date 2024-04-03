package service

import (
	"context"

	"github.com/Homeppv2/entitys"
)

type (
	UserServiceInterface interface {
		Register(ctx context.Context, req entitys.User) (res entitys.User, err error)
		GetByID(ctx context.Context, userID int) (res entitys.User, err error)
		GetByEmail(ctx context.Context, email string) (res entitys.User, err error)
		GetControllersByUserId(ctx context.Context, userID int) (res []entitys.ControllersData, err error)
	}


	ControllerServiceInterface interface {
		GetListMessangesFromIdForUserId(ctx context.Context, count int, from int, userID int) (msg []entitys.MessangeTypeZiroJson, err error)
		// GetControllersByUserId(ctx context.Context, user_id int) (controllers []entitys.ControllersData, err error)
		GetIdControllerByTypeAndNumber(ctx context.Context, type_, number_ int) (id int, err error)
		CreateMessageTypeOne(ctx context.Context, msg entitys.MessageTypeOneJSON) (err error)
		CreateMessageTypeTwo(ctx context.Context, msg entitys.MessageTypeTwoJSON) (err error)
		CreateMessageTypeThree(ctx context.Context, msg entitys.MessageTypeThreeJSON) (err error)
	}
)
