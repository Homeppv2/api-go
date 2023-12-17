package infrastructure

import (
	"homepp/api-go/internal/service/infrastructure"
)

type (
	pgManager struct {
		controllerRepo *ControllerRepo
		userRepo       *UserRepo
	}
)

func (m *pgManager) Controller() infrastructure.ControllerGateway {
	return m.controllerRepo
}

func (m *pgManager) User() infrastructure.UserGateway {
	return m.userRepo
}

var _ infrastructure.EntityManager = &pgManager{}
