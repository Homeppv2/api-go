package infrastructure

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"homepp/api-go/internal/entity"
	"homepp/api-go/internal/service/infrastructure"
	"homepp/api-go/pkg/errors"
)

type (
	AuthRepo struct {
		rds       *redis.Client
		l         *slog.Logger
		expiresIn time.Duration
	}
)

func NewSessionRepo(rds *redis.Client, l *slog.Logger, expiresIn time.Duration) *AuthRepo {
	return &AuthRepo{
		rds:       rds,
		l:         l,
		expiresIn: expiresIn,
	}
}

func (r *AuthRepo) Create(ctx context.Context, req entity.Session) (err error) {
	res := r.rds.Set(ctx, req.Token, req.UserID, r.expiresIn)
	if res.Err() != nil {
		r.l.Error("failed in AuthRepo.Create", err.Error())
		return errors.ErrUnknown
	}
	return nil
}

func (r *AuthRepo) GetByToken(ctx context.Context, token string) (res entity.Session, err error) {
	userID, err := r.rds.Get(ctx, token).Result()
	if err != nil {
		return entity.Session{}, errors.ErrSessionNotFound
	}

	resUserID, err := strconv.Atoi(userID)
	if err != nil {
		r.l.Error("failed in AuthRepo.GetByToken", err.Error())
		return entity.Session{}, errors.ErrUnknown
	}

	return entity.Session{Token: token, UserID: int64(resUserID)}, nil
}

func (r *AuthRepo) Delete(ctx context.Context, token string) (err error) {
	err = r.rds.Del(ctx, token).Err()
	if err != nil {
		r.l.Error("failed in AuthRepo.Delete", err.Error())
		return err
	}

	return nil
}

var _ infrastructure.SessionGateway = (*AuthRepo)(nil)
