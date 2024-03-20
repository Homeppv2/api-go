package infrastructure

import (
	"context"
	"log/slog"

	"api-go/internal/entity"
	"api-go/internal/service/infrastructure"

	"github.com/jmoiron/sqlx"
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

func (r *ControllerRepo) GetIdController(ctx context.Context, type_ int, number int) (id int, err error) {
	q := `
		SELECT id_contorller from contollers where type_controller=$1 and number_controller=$2;
		`
	err = r.db.GetContext(ctx, &id, q, type_, number)
	if err == nil {
		return id, err
	}
	q = `
	INSERT INTO controllers (type_controller, number_controller)
	VALUES ($1, $2)
	RETURNING id;
	`
	err = r.db.GetContext(ctx, &id, q, type_, number)
	return id, err
}

func (r *ControllerRepo) CreateMessangeControllerTypeOne(ctx context.Context, id int, main entity.MainMessangesData, add entity.ContollersLeackMessangesData) (err error) {
	q := `
	INSERT INTO messanges (id_contorller, status_controller, charge_controller, temperature_MK_controller)
	VALUES ($1, $2, $3, $4) RETURNING id_messange;

	`
	id_mess_sql, err := r.db.ExecContext(ctx, q, id, main.Status_controller, main.Charge_controller, main.Temperature_MK_controller)
	if err != nil {
		return err
	}
	id_mess, _ := id_mess_sql.LastInsertId()
	q = `
	INSERT INTO messanges (id_messange, leack)
	VALUES ($1, $2);
	`
	_, err = r.db.ExecContext(ctx, q, id_mess, add.Leack)
	return err
}

func (r *ControllerRepo) CreateMessangeControllerTypeTwo(ctx context.Context, id int, main entity.MainMessangesData, add entity.ContollersModuleMessangesData) (err error) {
	q := `
	INSERT INTO messanges (id_contorller, status_controller, charge_controller, temperature_MK_controller)
	VALUES ($1, $2, $3, $4) RETURNING id_messange;

	`
	id_mess_sql, err := r.db.ExecContext(ctx, q, id, main.Status_controller, main.Charge_controller, main.Temperature_MK_controller)
	if err != nil {
		return err
	}
	id_mess, _ := id_mess_sql.LastInsertId()
	q = `
	INSERT INTO messanges_contollers_module (id_messange, temperature, humidity, pressure, gas)
	VALUES ($1, $2, $3, $4, $5);
	`
	_, err = r.db.ExecContext(ctx, q, id_mess, add.Temperature, add.Humidity, add.Pressure, add.Gas)
	return err
}

func (r *ControllerRepo) CreateMessangeControllerTypeThree(ctx context.Context, id int, main entity.MainMessangesData, add entity.ControlerEnviromentDataMessange) (err error) {
	q := `
	INSERT INTO messanges (id_contorller, status_controller, charge_controller, temperature_MK_controller)
	VALUES ($1, $2, $3, $4) RETURNING id_messange;

	`
	id_mess_sql, err := r.db.ExecContext(ctx, q, id, main.Status_controller, main.Charge_controller, main.Temperature_MK_controller)
	if err != nil {
		return err
	}
	id_mess, _ := id_mess_sql.LastInsertId()
	q = `
	INSERT INTO messanges (id_messange, temperature, humidity, pressure, VOC, gas1, gas2, gas3, pm1, pm25, pm10, fire, smoke)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`
	_, err = r.db.ExecContext(ctx, q, id_mess, add.Temperature, add.Humidity, add.Pressure, add.Voc, add.Gas1, add.Gas2, add.Gas3, add.Pm1, add.Pm25, add.Pm10, add.Fire, add.Smoke)
	return err

}

func (r *ControllerRepo) GetCountMessangesFromIdForUserId(ctx context.Context, count int, from int, userID int) (msg []entity.MessangeTypeZiroJson, err error) {
	q := `
	SELECT * from messanges where id_messange >= $1 and id_contorller = $2 ORDER BY id_messange limit $3;
	`
	var msgs []entity.MainMessangesData
	err = r.db.SelectContext(ctx, &msgs, q, from, userID, count)
	if err != nil {
		return nil, err
	}
	for _, v := range msgs {
		q = `
		SELECT * from contollers where id_contorller = $1;
		`
		var ctr entity.ControllersData
		err = r.db.GetContext(ctx, &ctr, q, v.Id_contorller)
		if err != nil {
			r.l.Error("ошибка получения данных", err)
			continue
		}
		switch ctr.Type_controller {
		case 10:
			q = `
				SELECT * from messanges_contollers_leack where id_messange = $1;
				`
			var leak entity.ControlerLeackDataMessange
			err = r.db.GetContext(ctx, &leak, q, v.Id_messange)
			if err != nil {
				r.l.Error("ошибка получения данных", err)
				continue
			}
			var json entity.MessangeTypeZiroJson
			json.One.MainMessageJSON.Type = ctr.Type_controller
			json.One.MainMessageJSON.Number = ctr.Number_controller
			json.One.MainMessageJSON.Charge = v.Charge_controller
			json.One.MainMessageJSON.Status = v.Status_controller
			json.One.MainMessageJSON.Temperature_MK = v.Temperature_MK_controller
			json.One.Controlerleack.Leack = leak.Leack
			msg = append(msg, json)
			break
		case 11:
			q = `
				SELECT * from messanges_contollers_module where id_messange = $1;
				`
			var module entity.ContollersModuleMessangesData
			err = r.db.GetContext(ctx, &module, q, v.Id_messange)
			if err != nil {
				r.l.Error("ошибка получения данных", err)
				continue
			}
			var json entity.MessangeTypeZiroJson
			json.Two.MainMessageJSON.Type = ctr.Type_controller
			json.Two.MainMessageJSON.Number = ctr.Number_controller
			json.Two.MainMessageJSON.Charge = v.Charge_controller
			json.Two.MainMessageJSON.Status = v.Status_controller
			json.Two.MainMessageJSON.Temperature_MK = v.Temperature_MK_controller
			json.Two.Controlermodule.Gas = module.Gas
			json.Two.Controlermodule.Humidity = module.Humidity
			json.Two.Controlermodule.Pressure = module.Pressure
			json.Two.Controlermodule.Temperature = module.Temperature
			msg = append(msg, json)
			break
		case 12:
			q = `
				SELECT * from messanges_contollers_enviroment where id_messange = $1;
				`
			var env entity.ControlerEnviromentDataMessange
			err = r.db.GetContext(ctx, &env, q, v.Id_messange)
			if err != nil {
				r.l.Error("ошибка получения данных", err)
				continue
			}
			var json entity.MessangeTypeZiroJson
			json.Three.MainMessageJSON.Type = ctr.Type_controller
			json.Three.MainMessageJSON.Number = ctr.Number_controller
			json.Three.MainMessageJSON.Charge = v.Charge_controller
			json.Three.MainMessageJSON.Status = v.Status_controller
			json.Three.MainMessageJSON.Temperature_MK = v.Temperature_MK_controller
			json.Three.Controlerenviroment.Fire = env.Fire
			json.Three.Controlerenviroment.Gas1 = env.Gas1
			json.Three.Controlerenviroment.Gas2 = env.Gas2
			json.Three.Controlerenviroment.Gas3 = env.Gas3
			json.Three.Controlerenviroment.Humidity = env.Humidity
			json.Three.Controlerenviroment.Pm1 = env.Pm1
			json.Three.Controlerenviroment.Pm10 = env.Pm10
			json.Three.Controlerenviroment.Pm25 = env.Pm25
			json.Three.Controlerenviroment.Pressure = env.Pressure
			json.Three.Controlerenviroment.Smoke = env.Pressure
			json.Three.Controlerenviroment.Temperature = env.Temperature
			json.Three.Controlerenviroment.Voc = env.Voc
			break
		default:
			r.l.Error("ошибка получения данных - такого типа контролера нет", err)
			continue
		}
	}
	return msg, err
}

func (r *ControllerRepo) GetControllerById(ctx context.Context, controller_id int) (controller entity.ControllersData, err error) {
	q := `
	SELECT * from contollers where id_contorller = $1;
	`
	var ctrl entity.ControllersData
	err = r.db.GetContext(ctx, &ctrl, q, controller_id)
	if err != nil {
		return entity.ControllersData{}, err
	}
	return ctrl, nil
}

func (r *ControllerRepo) GetControllersByUserId(ctx context.Context, user_id int) (controllers []entity.ControllersData, err error) {
	q := `
	SELECT id_contorller from user_conttollers where id = $1;
	`
	var id_ctrls []int
	err = r.db.SelectContext(ctx, id_ctrls, q, user_id)
	if err != nil {
		return nil, err
	}
	var ans []entity.ControllersData
	for i := 0; i < len(id_ctrls); i++ {
		tmp, err := r.GetControllerById(ctx, id_ctrls[i])
		if err != nil {
			return nil, err
		}
		ans = append(ans, tmp)
	}
	return ans, nil
}

var _ infrastructure.ControllerGateway = (*ControllerRepo)(nil)
