package database

import (
	"context"

	"github.com/Homeppv2/entitys"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

func (r *Database) CreateUser(ctx context.Context, req entitys.User) (res entitys.User, err error) {
	q := `
		INSERT INTO users (username, email, hash_password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, hash_password;
		`
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return entitys.User{}, err
	}
	rows, err := tx.Query(context.Background(), q, req.Username, req.Email, req.HashPassword)
	if err != nil {
		tx.Rollback(context.Background())
		return entitys.User{}, err
	}
	defer rows.Close()
	err = tx.Commit(context.Background())
	if err != nil {
		return entitys.User{}, err
	}
	if rows.Next() {
		rows.Scan(&res.ID, &res.Username, &res.Email, &res.HashPassword)
	}
	return res, nil
}

func (r *Database) GetUserByID(ctx context.Context, id int) (res entitys.User, err error) {
	q := `
		SELECT id, username, email, hash_password
		FROM users
		WHERE id = $1`

	rows, err := r.pool.Query(context.Background(), q, id)
	if err != nil {
		return entitys.User{}, err
	}
	if rows.Next() {
		rows.Scan(&res.ID, &res.Username, &res.Email, &res.HashPassword)
	}
	defer rows.Close()
	return res, nil
}

func (r *Database) GetUserByEmail(ctx context.Context, email string) (res entitys.User, err error) {
	q := `
		SELECT id, username, email, hash_password
		FROM users
		WHERE email = $1;`

	rows, err := r.pool.Query(context.Background(), q, email)
	if err != nil {
		return entitys.User{}, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&res.ID, &res.Username, &res.Email, &res.HashPassword)
	}
	return res, nil
}

func (r *Database) GetIdUserControllers(ctx context.Context, type_ int, number int, id_user int) (id int, err error) {
	q := `
		SELECT id_controller from controllers where type_controller=$1 and number_controller=$2;
		`
	var id_ctr int
	rows, err := r.pool.Query(context.Background(), q, type_, number)
	if err != nil {
		id = -1
		return -1, err
	}
	if rows.Next() {
		rows.Scan(&id_ctr)
		return
	}
	rows.Close()
	q1 := `
	SELECT id from user_controllers where id_user = $1 and id_controller = $2;
	`
	rows, err = r.pool.Query(context.Background(), q1, id_user, id_ctr)
	if err != nil {
		id = -1
		return -1, err
	}
	if rows.Next() {
		rows.Scan(&id)
		return
	}
	rows.Close()
	return
}

func (r *Database) CreateMessangeControllerTypeOne(ctx context.Context, idController int, main entitys.MainMessangesData, add entitys.ContollersLeackMessangesData) (err error) {
	q := `
	INSERT INTO messanges (id_user_controllers, status_controller, charge_controller, temperature_MK_controller, t)
	VALUES ($1, $2, $3, $4, $5) RETURNING id_messange;

	`
	var id int
	rows, err := r.pool.Query(context.Background(), q, idController, main.Status_controller, main.Charge_controller, main.Temperature_MK_controller, main.Tm)
	if err != nil {
		return
	}
	if rows.Next() {
		rows.Scan(&id)
	}
	rows.Close()
	q = `
	INSERT INTO messanges_controllers_leack (id_messange, leack)
	VALUES ($1, $2);
	`
	rows, err = r.pool.Query(context.Background(), q, id, add.Leack)
	if err != nil {
		return
	}
	rows.Close()
	return
}

func (r *Database) CreateMessangeControllerTypeTwo(ctx context.Context, idController int, main entitys.MainMessangesData, add entitys.ContollersModuleMessangesData) (err error) {
	q := `
	INSERT INTO messanges (id_user_controllers, status_controller, charge_controller, temperature_MK_controller, t)
	VALUES ($1, $2, $3, $4, $5) RETURNING id_messange;

	`
	var id int
	rows, err := r.pool.Query(context.Background(), q, idController, main.Status_controller, main.Charge_controller, main.Temperature_MK_controller, main.Tm)
	if err != nil {
		return
	}
	if rows.Next() {
		rows.Scan(&id)
	}
	rows.Close()
	q = `
	INSERT INTO messanges_controllers_module (id_messange, temperature, humidity, pressure, gas)
	VALUES ($1, $2, $3, $4, $5);
	`
	rows, err = r.pool.Query(context.Background(), q, id, add.Temperature, add.Humidity, add.Pressure, add.Gas)
	if err != nil {
		return
	}
	rows.Close()
	return
}

func (r *Database) CreateMessangeControllerTypeThree(ctx context.Context, idController int, main entitys.MainMessangesData, add entitys.ControlerEnviromentDataMessange) (err error) {
	q := `
	INSERT INTO messanges (id_user_controllers, status_controller, charge_controller, temperature_MK_controller, t)
	VALUES ($1, $2, $3, $4, $5) RETURNING id_messange;

	`
	var id int
	rows, err := r.pool.Query(context.Background(), q, idController, main.Status_controller, main.Charge_controller, main.Temperature_MK_controller, main.Tm)
	if err != nil {
		return
	}
	if rows.Next() {
		rows.Scan(&id)
	}
	rows.Close()
	q = `
	INSERT INTO messanges_controllers_enviroment (id_messange, temperature, humidity, pressure, VOC, gas1, gas2, gas3, pm1, pm25, pm10, fire, smoke)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);
	`
	rows, err = r.pool.Query(context.Background(), q, id, add.Temperature, add.Humidity, add.Pressure, add.Voc, add.Gas1, add.Gas2, add.Gas3, add.Pm1, add.Pm25, add.Pm10, add.Fire, add.Smoke)
	rows.Close()
	return
}

func (r *Database) GetListMessangesFromIdForUserId(ctx context.Context, count int, from int, userID int) (msg []entitys.MessangeTypeZiroJson, err error) {
	q := `
	WITH ctrls as (
		SELECT id_controller from user_controllers where id_user = $1
	)
	SELECT id_messange, id_user_controllers, status_controller, charge_controller, temperature_MK_controller, t from messanges where id_messange >= $2 and id_contorller = ctrl.id_contorller ORDER BY id_messange limit $3;
	`
	conn, err := r.pool.Acquire(context.Background())
	rows, err := conn.Query(context.Background(), q, userID, from, count)
	var msgs []entitys.MainMessangesData
	for rows.Next() {
		var main entitys.MainMessangesData
		rows.Scan(&main.Id_messange, &main.Id_contorller, &main.Status_controller, &main.Charge_controller, &main.Temperature_MK_controller, &main.Tm)
		msgs = append(msgs, main)
	}
	rows.Close()
	for _, v := range msgs {
		q = `
		SELECT type_controller, number_controller from contollers where id_controller = $1;
		`
		rows, err = conn.Query(context.Background(), q, v.Id_contorller)
		var type_ctrl, nmbr_ctrl int
		if rows.Next() {
			rows.Scan(&type_ctrl, &nmbr_ctrl)
		}
		rows.Close()
		switch type_ctrl {
		case 10:
			q = `
				SELECT leack from messanges_contollers_leack where id_messange = $1;
				`
			var leak entitys.ControlerLeackDataMessange
			rows, err = conn.Query(context.Background(), q, v.Id_messange)
			if rows.Next() {
				rows.Scan(&leak.Leack)
			}
			rows.Close()
			var json entitys.MessangeTypeZiroJson
			json.One = new(entitys.MessageTypeOneJSON)
			json.One.MainMessageJSON.Type = type_ctrl
			json.One.MainMessageJSON.Number = nmbr_ctrl
			json.One.MainMessageJSON.Charge = v.Charge_controller
			json.One.MainMessageJSON.Status = v.Status_controller
			json.One.MainMessageJSON.Temperature_MK = v.Temperature_MK_controller
			json.One.MainMessageJSON.Data = v.Tm
			json.One.Controlerleack = &leak
			msg = append(msg, json)
			break
		case 11:
			q = `
				SELECT temperature, humidity, pressure, gas from messanges_controllers_module where id_messange = $1;
				`
			var module entitys.ControlerModuleDataMessange
			rows, err = conn.Query(context.Background(), q, v.Id_messange)
			if rows.Next() {
				rows.Scan(&module.Humidity, &module.Humidity, &module.Pressure, &module.Gas)
			}
			rows.Close()
			var json entitys.MessangeTypeZiroJson
			json.Two = new(entitys.MessageTypeTwoJSON)
			json.Two.MainMessageJSON.Type = type_ctrl
			json.Two.MainMessageJSON.Number = nmbr_ctrl
			json.Two.MainMessageJSON.Charge = v.Charge_controller
			json.Two.MainMessageJSON.Status = v.Status_controller
			json.Two.MainMessageJSON.Temperature_MK = v.Temperature_MK_controller
			json.Two.MainMessageJSON.Data = v.Tm
			json.Two.Controlermodule = &module
			msg = append(msg, json)
			break
		case 12:
			q = `
				SELECT temperature, humidity, pressure, VOC, gas1, gas2, gas3, pm1, pm25, pm10, fire, smoke from messanges_controllers_enviroment where id_messange = $1;
				`
			var env entitys.ControlerEnviromentDataMessange
			rows, err = conn.Query(context.Background(), q, v.Id_messange)
			if rows.Next() {
				rows.Scan(&env.Temperature, &env.Humidity, &env.Pressure, &env.Voc, &env.Gas1, &env.Gas1, &env.Gas2, &env.Gas3, &env.Pm1, &env.Pm25, &env.Pm10, &env.Fire, &env.Smoke)
			}
			rows.Close()
			var json entitys.MessangeTypeZiroJson
			json.Three = new(entitys.MessageTypeThreeJSON)
			json.Three.MainMessageJSON.Type = type_ctrl
			json.Three.MainMessageJSON.Number = nmbr_ctrl
			json.Three.MainMessageJSON.Charge = v.Charge_controller
			json.Three.MainMessageJSON.Status = v.Status_controller
			json.Three.MainMessageJSON.Temperature_MK = v.Temperature_MK_controller
			json.Three.MainMessageJSON.Data = v.Tm
			json.Three.Controlerenviroment = &env
			break
		default:
			continue
		}
	}
	return
}

func (r *Database) GetControllerById(ctx context.Context, controller_id int) (controller entitys.ControllersData, err error) {
	q := `
	SELECT id_controller, type_controller, number_controller from controllers where id_controller = $1;
	`
	rows, err := r.pool.Query(ctx, q, controller_id)
	if err != nil {
		return entitys.ControllersData{}, err
	}
	if rows.Next() {
		rows.Scan(&controller.Id_contorller, &controller.Number_controller, &controller.Type_controller)
	}
	rows.Close()
	return
}

func (r *Database) GetListControllersByUserId(ctx context.Context, user_id int) (controllers []entitys.ControllersData, err error) {
	q := `
	SELECT id_controller from user_controllers where id_user = $1;
	`
	var id_ctrls []int
	rows, err := r.pool.Query(ctx, q, user_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		id_ctrls = append(id_ctrls, id)
	}
	rows.Close()
	for i := 0; i < len(id_ctrls); i++ {
		tmp, err := r.GetControllerById(ctx, id_ctrls[i])
		if err != nil {
			return nil, err
		}
		controllers = append(controllers, tmp)
	}
	return
}

func (r *Database) GetListContorllers(ctx context.Context) (controllers []entitys.ControllersData, err error) {
	q := `
	SELECT id_controller from controllers;
	`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ctrl entitys.ControllersData
		rows.Scan(&ctrl.Id_contorller)
		controllers = append(controllers, ctrl)
	}
	rows.Close()
	return
}
func (r *Database) GetIdControllers(ctx context.Context, type_, number_ int) (id int, err error) {
	q := `
	SELECT id_controller from controllers where type_controller=$1 and number_controller=$2;
	`
	var id_ctr int
	rows, err := r.pool.Query(context.Background(), q, type_, number_)
	if err != nil {
		id = -1
		return -1, err
	}
	if rows.Next() {
		rows.Scan(&id_ctr)
		return
	}
	rows.Close()
	return
}
