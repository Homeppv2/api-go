package service

import (
	"context"

	controllerpkg "api-go/internal/controller"
	"api-go/internal/entity"
	infrastructure2 "api-go/internal/service/infrastructure"
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

func (c *ControllerService) GetCountMessangesFromIdForUserId(ctx context.Context, count int, from int, userID int) (msg []entity.MessangeTypeZiroJson, err error) {
	// достать все смс начная с определенного для такогоо юзер айди
	q2 := "select id_contorller from user_conttollers where id_user=$1;"
	q1 := `select (id_messange, id_contorller, status_controller, charge_controller, temperature_MK_controller) from messanges where id_contorller=$1 and id_messange > $2 limit $3;`
	q3 := "select (id_messange, leack) from messanges_contollers_leack where id_messange=$1;"
	q4 := "select (id_messange, temperature, humidity, pressure, gas) from messanges_contollers_module where id_messange=$1;"
	q5 := "select (id_messange, temperature, humidity, humidity, VOC, gas1, gas2, gas3, pm1, pm25, pm10, fire, smoke) where id_messange=$1;"

	return nil, nil
}

func (c *ControllerService) GetControllersByUserId(ctx context.Context, user_id int) (controllers []entity.ControllersData, err error) {
	// из таблицы достать контролеры по юзер айди
	q := `with tmp as (
		select id_contorller from user_conttollers where id_user=$1
	) select (id_contorller, type_controller, number_controller) from contollers where id_contorller=tmp.id;
	`

	return nil, nil
}

var _ controllerpkg.ControllerService = (*ControllerService)(nil)
