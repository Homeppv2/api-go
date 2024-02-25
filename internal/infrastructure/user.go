package infrastructure

import (
	"context"
	"database/sql"
	"log/slog"

	"api-go/internal/entity"
	"api-go/internal/service/infrastructure"
	"api-go/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type (
	UserRepo struct {
		db queryRunner
		l  *slog.Logger
	}
)

func NewUserRepo(db *sqlx.DB, l *slog.Logger) *UserRepo {
	return &UserRepo{
		db: db,
		l:  l,
	}
}

func (r *UserRepo) Create(ctx context.Context, req entity.CreateUserDTO) (res entity.User, err error) {
	q := `
		INSERT INTO users (username, email, hash_password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, hash_password;
		`

	err = r.db.GetContext(ctx, &res, q, req.Username, req.Email, req.HashPassword)
	if err != nil {
		r.l.Error("failed in UserRepo.Create: ", err.Error())
		return entity.User{}, err
	}

	return res, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (res entity.User, err error) {
	q := `
		SELECT id, username, email, hash_password
		FROM users
		WHERE id = $1`

	err = r.db.GetContext(ctx, &res, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, errors.ErrUserNotFound
		}
		r.l.Error("failed in UserRepo.GetByID: ", err.Error())
		return entity.User{}, err
	}

	return res, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (res entity.User, err error) {
	q := `
		SELECT id, username, email, hash_password
		FROM users
		WHERE email = $1`

	err = r.db.GetContext(ctx, &res, q, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, errors.ErrUserNotFound
		}
		r.l.Error("failed in UserRepo.getByEmail: ", err.Error())
		return entity.User{}, err
	}

	return res, nil
}

var _ infrastructure.UserGateway = (*UserRepo)(nil)
