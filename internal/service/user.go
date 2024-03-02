package service

import (
	"context"

	controllerpkg "api-go/internal/controller"
	"api-go/internal/entity"
	infrastructure2 "api-go/internal/service/infrastructure"
	"api-go/pkg/errors"
	"api-go/pkg/hasher"

	"gopkg.in/guregu/null.v4"
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

func (s *UserService) Register(ctx context.Context, req entity.RegisterUserRequest) (res entity.User, err error) {
	controller, err := s.repos.Controller().GetByHwKey(ctx, req.HwKey)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	if *controller.IsUsed {
		return entity.User{}, errors.HandleServiceError(errors.ErrControllerAlreadyIsUsed)
	}

	_, err = s.repos.User().GetByEmail(ctx, req.Email)
	if err != errors.ErrUserNotFound {
		if err == nil {
			return entity.User{}, errors.HandleServiceError(errors.ErrUserEmailAlreadyExist)
		}
		return entity.User{}, errors.HandleServiceError(err)
	}

	hashPassword, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return entity.User{}, errors.HandleServiceError(err)
	}

	err = s.repos.WithTx(ctx, func(m infrastructure2.EntityManager) (err error) {
		res, err = s.repos.User().Create(ctx, entity.CreateUserDTO{
			Username:     req.Username,
			Email:        req.Email,
			HashPassword: hashPassword,
		})
		if err != nil {
			return err
		}

		trueIsUsed := true
		_, err = s.repos.Controller().UpdateIsUsed(ctx, entity.UpdateControllerIsUsedByRequest{
			ID:       controller.ID,
			IsUsed:   &trueIsUsed,
			IsUsedBy: null.NewInt(res.ID, true),
		})
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

func (s *UserService) GetByID(ctx context.Context, userID int64) (res entity.User, err error) {
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
