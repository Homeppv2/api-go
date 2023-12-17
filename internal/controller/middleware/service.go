package middleware

import (
	"context"

	"homepp/api-go/internal/entity"
)

type AuthService interface {
	CheckSession(ctx context.Context, token string) (res entity.Session, err error)
}
