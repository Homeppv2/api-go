package service

import (
	"context"
	"github.com/Homeppv2/api-go/internal/database"

	"github.com/Homeppv2/entitys"
)

type (
	ControllerService struct {
		base database.DatabaseInterface
	}
)

func NewControllerService(base database.DatabaseInterface) *ControllerService {
	return &ControllerService{
		base : base,
	}
}

func (c *ControllerService) GetListMessangesFromIdForUserId(ctx context.Context, count int, from int, userID int) (msg []entitys.MessangeTypeZiroJson, err error) {
	msg, err = c.base.GetListMessangesFromIdForUserId(ctx, count, from, userID)
	return
}

func (c *ControllerService) GetIdControllerByTypeAndNumber(ctx context.Context, type_, number_ int) (id int, err error) {
	id, err = c.base.GetIdController(ctx, type_, number_)
	return
}
func (c *ControllerService)	CreateMessageTypeOne(ctx context.Context, msg entitys.MessageTypeOneJSON) (err error) {
	id, err := c.GetIdControllerByTypeAndNumber(ctx, msg.Type, msg.Number)
	err = c.base.CreateMessangeControllerTypeOne(ctx, id, entitys.MainMessangesData{Id_contorller:id, Status_controller: msg.Status, Charge_controller: msg.Charge, Temperature_MK_controller: msg.Temperature_MK}, entitys.ContollersLeackMessangesData{Leack: msg.Controlerleack.Leack})
	return
}


func (c *ControllerService)	CreateMessageTypeTwo(ctx context.Context, msg entitys.MessageTypeTwoJSON) (err error) {
	id, err := c.GetIdControllerByTypeAndNumber(ctx, msg.Type, msg.Number)
	err = c.base.CreateMessangeControllerTypeTwo(ctx, id, entitys.MainMessangesData{Id_contorller:id, Status_controller: msg.Status, Charge_controller: msg.Charge, Temperature_MK_controller: msg.Temperature_MK}, 
		entitys.ContollersModuleMessangesData{Humidity: msg.Controlermodule.Humidity, Temperature: msg.Controlermodule.Temperature, Pressure: msg.Controlermodule.Pressure, Gas: msg.Controlermodule.Gas})
	return
}


func (c *ControllerService) CreateMessageTypeThree(ctx context.Context, msg entitys.MessageTypeThreeJSON) (err error) {
	id, err := c.GetIdControllerByTypeAndNumber(ctx, msg.Type, msg.Number)
	err = c.base.CreateMessangeControllerTypeThree(ctx, id, entitys.MainMessangesData{Id_contorller:id, Status_controller: msg.Status, Charge_controller: msg.Charge, Temperature_MK_controller: msg.Temperature_MK}, 
		entitys.ControlerEnviromentDataMessange{Humidity: msg.Controlerenviroment.Humidity,
			 Temperature: msg.Controlerenviroment.Temperature,
			 Pressure: msg.Controlerenviroment.Pressure,
			 Gas1: msg.Controlerenviroment.Gas1, 
			 Gas2: msg.Controlerenviroment.Gas2,
			 Gas3: msg.Controlerenviroment.Gas3,
			 Pm1: msg.Controlerenviroment.Pm1,
			 Pm10: msg.Controlerenviroment.Pm10,
			 Pm25: msg.Controlerenviroment.Pm25,
			 Fire: msg.Controlerenviroment.Fire,
			 Smoke: msg.Controlerenviroment.Smoke,
			})
	return
}


