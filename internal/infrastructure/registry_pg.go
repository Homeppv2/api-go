package infrastructure

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"homepp/api-go/internal/service/infrastructure"
)

type (
	PGRegistry struct {
		db *sqlx.DB
		m  *pgManager
		l  *slog.Logger
	}
)

func (r *PGRegistry) Controller() infrastructure.ControllerGateway {
	return r.m.controllerRepo
}

func (r *PGRegistry) User() infrastructure.UserGateway {
	return r.m.userRepo
}

func (r *PGRegistry) WithTx(ctx context.Context, f func(manager infrastructure.EntityManager) error) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = f(&pgManager{
		userRepo:       &UserRepo{tx, r.l},
		controllerRepo: &ControllerRepo{tx, r.l},
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

var _ infrastructure.Registry = &PGRegistry{}

func NewPGRegistry(db *sqlx.DB, l *slog.Logger) *PGRegistry {
	return &PGRegistry{
		db: db,
		m: &pgManager{
			userRepo:       &UserRepo{db, l},
			controllerRepo: &ControllerRepo{db, l},
		},
	}
}
