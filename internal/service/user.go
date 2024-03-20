package service

import (
	"context"

	controllerpkg "api-go/internal/controller"
	"api-go/internal/entity"
	infrastructure2 "api-go/internal/service/infrastructure"
	"api-go/pkg/errors"
	"api-go/pkg/hasher"
)

type (
	UserService struct {
		repos  infrastructure2.Registry
		hasher hasher.Hasher
	}
)

func NewUserService(repos infrastructure2.Registry, hasher hasher.Hasher) *UserService {
	return &UserService{
		repos:  repos,
		hasher: hasher,
	}
}

func (s *UserService) Register(ctx context.Context, req entity.User) (res entity.User, err error) {
	_, err = s.repos.User().GetByEmail(ctx, req.Email)
	if err != errors.ErrUserNotFound {
		if err == nil {
			return entity.User{}, errors.HandleServiceError(errors.ErrUserEmailAlreadyExist)
		}
		return entity.User{}, errors.HandleServiceError(err)
	}

	err = s.repos.WithTx(ctx, func(m infrastructure2.EntityManager) (err error) {
		res, err = s.repos.User().Create(ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}
	return res, nil
}

func (s *UserService) GetByID(ctx context.Context, userID int) (res entity.User, err error) {
	res, err = s.repos.User().GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	return res, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (res entity.User, err error) {
	res, err = s.repos.User().GetByEmail(ctx, email)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	return res, nil
}

var _ controllerpkg.UserService = (*UserService)(nil)
