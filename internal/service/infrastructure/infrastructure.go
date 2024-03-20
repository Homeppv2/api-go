package infrastructure

import (
	"context"

	"api-go/internal/entity"
)

type (
	UserGateway interface {
		Create(ctx context.Context, req entity.User) (res entity.User, err error)
		GetByID(ctx context.Context, id int) (res entity.User, err error)
		GetByEmail(ctx context.Context, email string) (res entity.User, err error)
	}

	// todo реализовать
	// функционал удаления смс
	ControllerGateway interface {
		GetIdController(ctx context.Context, type_ int, number int) (id int, err error)
		CreateMessangeControllerTypeOne(ctx context.Context, id int, main entity.MainMessangesData, add entity.ContollersLeackMessangesData) (err error)
		CreateMessangeControllerTypeTwo(ctx context.Context, id int, main entity.MainMessangesData, add entity.ContollersModuleMessangesData) (err error)
		CreateMessangeControllerTypeThree(ctx context.Context, id int, main entity.MainMessangesData, add entity.ControlerEnviromentDataMessange) (err error)
		GetCountMessangesFromIdForUserId(ctx context.Context, count int, from int, userID int) (msg []entity.MessangeTypeZiroJson, err error)
		GetControllersByUserId(ctx context.Context, user_id int) (controllers []entity.ControllersData, err error)
	}

	EntityManager interface {
		Controller() ControllerGateway
		User() UserGateway
	}

	Registry interface {
		EntityManager
		WithTx(context.Context, func(EntityManager) error) error
	}
)
