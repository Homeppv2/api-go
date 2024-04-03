package service

import (
	"context"

	"github.com/Homeppv2/api-go/internal/database"
	"github.com/Homeppv2/entitys"

	"github.com/Homeppv2/api-go/pkg/errors"
	"github.com/Homeppv2/api-go/pkg/hasher"
)

type (
	UserService struct {
		base database.DatabaseInterface
		hasher hasher.Hasher
	}
)

func NewUserService(base database.DatabaseInterface, hasher hasher.Hasher) *UserService {
	return &UserService{
		base:  base,
		hasher: hasher,
	}
}

func (s *UserService) Register(ctx context.Context, req entitys.User) (res entitys.User, err error) {
	_, err = s.base.GetUserByEmail(ctx, req.Email)
	if err != errors.ErrUserNotFound {
		if err == nil {
			return entitys.User{}, errors.HandleServiceError(errors.ErrUserEmailAlreadyExist)
		}
		return entitys.User{}, errors.HandleServiceError(err)
	}

	res, err = s.base.CreateUser(ctx, req)
	if err != nil {
		return entitys.User{}, errors.HandleServiceError(err)
	}
	return 
}

func (s *UserService) GetByID(ctx context.Context, userID int) (res entitys.User, err error) {
	res, err = s.base.GetUserByID(ctx, userID)
	if err != nil {
		return entitys.User{}, errors.HandleServiceError(err)
	}
	return 
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (res entitys.User, err error) {
	res, err = s.base.GetUserByEmail(ctx, email)
	if err != nil {
		return entitys.User{}, errors.HandleServiceError(err)
	}
	return 
}

func (s *UserService) GetControllersByUserId(ctx context.Context, userID int) (res []entitys.ControllersData, err error) {
	res, err = s.base.GetListControllersByUserId(ctx, userID)
	return
}
