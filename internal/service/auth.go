package service

import (
	"context"

	"github.com/google/uuid"
	controllerpkg "homepp/api-go/internal/controller"
	"homepp/api-go/internal/entity"
	"homepp/api-go/internal/service/infrastructure"
	"homepp/api-go/pkg/errors"
	"homepp/api-go/pkg/hasher"
)

type (
	AuthService struct {
		sessionRepo infrastructure.SessionGateway
		userRepo    infrastructure.UserGateway
		hasher      hasher.Hasher
	}
)

func NewAuthService(sessionRepo infrastructure.SessionGateway, userRepo infrastructure.UserGateway, hasher hasher.Hasher,
) *AuthService {
	return &AuthService{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		hasher:      hasher,
	}
}

func (s *AuthService) Login(ctx context.Context, req entity.LoginRequest) (res entity.LoginResponse, err error) {
	res.User, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return entity.LoginResponse{}, errors.HandleServiceError(err)
	}

	if !s.hasher.CompareAndHash(res.User.HashPassword, req.Password) {
		return entity.LoginResponse{}, errors.HandleServiceError(errors.ErrInvalidLogin)
	}

	token := uuid.New().String()
	err = s.sessionRepo.Create(ctx, entity.Session{UserID: res.User.ID, Token: token})
	if err != nil {
		return entity.LoginResponse{}, errors.HandleServiceError(err)
	}

	res.Session = entity.Session{
		UserID: res.User.ID,
		Token:  token,
	}

	return res, nil
}

func (s *AuthService) CheckSession(ctx context.Context, token string) (res entity.Session, err error) {
	res, err = s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return entity.Session{}, errors.HandleServiceError(err)
	}

	return res, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) (err error) {
	err = s.sessionRepo.Delete(ctx, token)
	if err != nil {
		return errors.HandleServiceError(err)
	}

	return nil
}

var _ controllerpkg.AuthService = (*AuthService)(nil)
