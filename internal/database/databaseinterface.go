package database

import (
	"context"
	"time"

	"github.com/Homeppv2/entitys"
)

type DatabaseInterface interface {
	CreateUser(ctx context.Context, req entitys.User) (res entitys.User, err error)
	GetUserByID(ctx context.Context, id int) (res entitys.User, err error)
	GetUserByEmail(ctx context.Context, email string) (res entitys.User, err error)
	GetIdUserControllers(ctx context.Context, type_ int, number int, id_user int) (id int, err error)
	CreateMessangeControllerTypeOne(ctx context.Context, id int, main entitys.MainMessangesData, add entitys.ContollersLeackMessangesData) (err error)
	CreateMessangeControllerTypeTwo(ctx context.Context, id int, main entitys.MainMessangesData, add entitys.ContollersModuleMessangesData) (err error)
	CreateMessangeControllerTypeThree(ctx context.Context, id int, main entitys.MainMessangesData, add entitys.ControlerEnviromentDataMessange) (err error)
	GetListMessangesFromIdForUserId(ctx context.Context, count int, from int, userID int) (msg []entitys.MessangeTypeZiroJson, err error)
	// GetListControllersByUserId(ctx context.Context, user_id int) (controllers []entitys.ControllersData, err error)
	GetListControllersByUserId(ctx context.Context, user_id int) (controllers []entitys.ControllersData, err error)
	GetListContorllers(ctx context.Context) (controllers []entitys.ControllersData, err error)
	GetIdControllers(ctx context.Context, type_, number_ int) (id int, err error)
	GetListMessangesFromDateForUserId(ctx context.Context, from, to time.Time, UserId int) (msg []entitys.MessangeTypeZiroJson, err error)
}
