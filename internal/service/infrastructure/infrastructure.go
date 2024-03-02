package infrastructure

import (
	"context"

	"api-go/internal/entity"
)

type (
	UserGateway interface {
		Create(ctx context.Context, req entity.CreateUserDTO) (res entity.User, err error)
		GetByID(ctx context.Context, id int64) (res entity.User, err error)
		GetByEmail(ctx context.Context, email string) (res entity.User, err error)
	}

	SessionGateway interface {
		Create(ctx context.Context, req entity.Session) (err error)
		GetByToken(ctx context.Context, token string) (res entity.Session, err error)
		Delete(ctx context.Context, token string) (err error)
	}

	ControllerGateway interface {
		Create(ctx context.Context, req entity.CreateControllerDTO) (res entity.Controller, err error)
		GetByID(ctx context.Context, id int64) (res entity.Controller, err error)
		GetByHwKey(ctx context.Context, hwKey string) (res entity.Controller, err error)
		GetByIsUsedBy(ctx context.Context, isUsedBy int64) (res []entity.Controller, err error)
		UpdateIsUsed(ctx context.Context, req entity.UpdateControllerIsUsedByRequest) (res entity.Controller, err error)
		Delete(ctx context.Context, id int64) (err error)
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
