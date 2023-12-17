package infrastructure

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"homepp/api-go/internal/entity"
	"homepp/api-go/internal/service/infrastructure"
	"homepp/api-go/pkg/errors"
)

type (
	ControllerRepo struct {
		db queryRunner
		l  *slog.Logger
	}
)

func NewControllerRepo(db *sqlx.DB, l *slog.Logger) *ControllerRepo {
	return &ControllerRepo{
		db: db,
		l:  l,
	}
}

func (r *ControllerRepo) Create(ctx context.Context, req entity.CreateControllerDTO) (res entity.Controller, err error) {
	q := `
		INSERT INTO controllers (hw_key, is_used, is_used_by)
		VALUES ($1, $2, $3)
		RETURNING id, hw_key, is_used, is_used_by;
		`

	err = r.db.GetContext(ctx, &res, q, req.HwKey, &req.IsUsed, req.IsUsedBy)
	if err != nil {
		r.l.Error("failed in ControllerRepo.Create: ", err.Error())
		return entity.Controller{}, err
	}

	return res, nil
}

func (r *ControllerRepo) GetByID(ctx context.Context, id int64) (res entity.Controller, err error) {
	q := `
		SELECT id, hw_key, is_used, is_used_by
		FROM controllers
		WHERE id = $1`

	err = r.db.GetContext(ctx, &res, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Controller{}, errors.ErrControllerNotFound
		}
		r.l.Error("failed in ControllerRepo.GetByID: ", err.Error())
		return entity.Controller{}, err
	}

	return res, nil
}

func (r *ControllerRepo) GetByHwKey(ctx context.Context, hwKey string) (res entity.Controller, err error) {
	q := `
		SELECT id, hw_key, is_used, is_used_by
		FROM controllers
		WHERE hw_key = $1`

	err = r.db.GetContext(ctx, &res, q, hwKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Controller{}, errors.ErrControllerNotFound
		}
		r.l.Error("failed in ControllerRepo.GetByHwKey: ", err.Error())
		return entity.Controller{}, err
	}

	return res, nil
}

func (r *ControllerRepo) GetByIsUsedBy(ctx context.Context, isUsedBy int64) (res []entity.Controller, err error) {
	q := `
		SELECT id, hw_key, is_used, is_used_by
		FROM controllers
		WHERE is_used_by = $1`

	err = r.db.SelectContext(ctx, &res, q, isUsedBy)
	if err != nil {
		r.l.Error("failed in ControllerRepo.GetByIsUsedBy: ", err.Error())
		return []entity.Controller{}, err
	}

	return res, nil
}

func (r *ControllerRepo) UpdateIsUsed(ctx context.Context, req entity.UpdateControllerIsUsedByRequest) (res entity.Controller, err error) {
	q := `
		UPDATE controllers
		SET
		    is_used = $2,
		    is_used_by = $3
		WHERE id=$1
		RETURNING id, hw_key, is_used, is_used_by;
	`

	err = r.db.GetContext(ctx, &res, q, req.ID, &req.IsUsed, req.IsUsedBy)
	if err != nil {
		r.l.Error("failed in ControllerRepo.UpdateIsUsed: ", err.Error())
		return entity.Controller{}, err
	}

	return res, nil
}

func (r *ControllerRepo) Delete(ctx context.Context, id int64) (err error) {
	q := `
		DELETE
		FROM controllers
		WHERE id=$1
		RETURNING 1;
		`
	var res int64
	err = r.db.GetContext(ctx, &res, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrControllerNotFound
		}
		return err
	}
	return nil
}

var _ infrastructure.ControllerGateway = (*ControllerRepo)(nil)
