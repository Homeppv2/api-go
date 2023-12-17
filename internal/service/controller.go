package service

import (
	"context"

	controllerpkg "homepp/api-go/internal/controller"
	"homepp/api-go/internal/entity"
	infrastructure2 "homepp/api-go/internal/service/infrastructure"
	"homepp/api-go/pkg/errors"
)

type (
	ControllerService struct {
		controllerRepo infrastructure2.ControllerGateway
		userRepo       infrastructure2.UserGateway
	}
)

func NewControllerService(controllerRepo infrastructure2.ControllerGateway, userRepo infrastructure2.UserGateway) *ControllerService {
	return &ControllerService{
		controllerRepo: controllerRepo,
		userRepo:       userRepo,
	}
}

func (s *ControllerService) Create(ctx context.Context, hwKey string) (res entity.Controller, err error) {
	_, err = s.controllerRepo.GetByHwKey(ctx, hwKey)
	if err != errors.ErrControllerNotFound {
		if err == nil {
			return entity.Controller{}, errors.HandleServiceError(errors.ErrControllerAlreadyExist)
		}
		return entity.Controller{}, errors.HandleServiceError(err)
	}

	res, err = s.controllerRepo.Create(ctx, entity.CreateControllerDTO{HwKey: hwKey, IsUsed: false})
	if err != nil {
		return entity.Controller{}, errors.HandleServiceError(err)
	}

	return res, nil

}

func (s *ControllerService) GetByID(ctx context.Context, id int64) (res entity.Controller, err error) {
	res, err = s.controllerRepo.GetByID(ctx, id)
	if err != nil {
		return entity.Controller{}, errors.HandleServiceError(err)
	}

	return res, nil
}

func (s *ControllerService) GetByHwKey(ctx context.Context, hwKey string) (res entity.Controller, err error) {
	res, err = s.controllerRepo.GetByHwKey(ctx, hwKey)
	if err != nil {
		return entity.Controller{}, errors.HandleServiceError(err)
	}

	return res, nil
}

func (s *ControllerService) GetByIsUsedBy(ctx context.Context, isUsedBy int64) (res []entity.Controller, err error) {
	_, err = s.userRepo.GetByID(ctx, isUsedBy)
	if err != nil {
		return []entity.Controller{}, errors.HandleServiceError(err)
	}

	res, err = s.controllerRepo.GetByIsUsedBy(ctx, isUsedBy)
	if err != nil {
		return []entity.Controller{}, errors.HandleServiceError(err)
	}

	return res, nil
}

func (s *ControllerService) UpdateIsUsed(ctx context.Context, req entity.UpdateControllerIsUsedByRequest) (res entity.Controller, err error) {
	if req.IsUsedBy.Valid {
		_, err = s.userRepo.GetByID(ctx, req.IsUsedBy.ValueOrZero())
		if err != nil {
			return entity.Controller{}, errors.HandleServiceError(err)
		}
	}

	res, err = s.controllerRepo.GetByID(ctx, req.ID)
	if err != nil {
		return entity.Controller{}, errors.HandleServiceError(err)
	}

	res, err = s.controllerRepo.UpdateIsUsed(ctx, req)
	if err != nil {
		return entity.Controller{}, errors.HandleServiceError(err)
	}

	return res, nil
}

func (s *ControllerService) Delete(ctx context.Context, id int64) (err error) {
	err = s.controllerRepo.Delete(ctx, id)
	if err != nil {
		return errors.HandleServiceError(err)
	}

	return nil
}

var _ controllerpkg.ControllerService = (*ControllerService)(nil)
