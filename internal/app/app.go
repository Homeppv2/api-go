package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Homeppv2/api-go/internal/broker"
	"github.com/Homeppv2/api-go/internal/database"
	"github.com/Homeppv2/api-go/internal/middleware"
	"github.com/Homeppv2/api-go/internal/service"
	"github.com/Homeppv2/api-go/pkg/hasher"
	"github.com/gofiber/fiber/v2/log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run() {

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	l.Info("success initializing logger")

	// db
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("POSTGRESQL_USER"), os.Getenv("POSTGRESQL_PASSWORD"), os.Getenv("POSTGRESQL_HOST"), os.Getenv("POSTGRESQL_PORT"), os.Getenv("POSTGRESQL_BASE")))
	if err != nil {
		l.Info("failed to connect to postgresql: %s", err.Error())
		return
	}
	defer func(db *pgxpool.Pool) {
		db.Close()
	}(db)

	err = db.Ping(context.Background())
	if err != nil {
		l.Info("failed to ping to postgresql: %s", err.Error())
		return
	}

	l.Info("success connecting to postgresql")
	uriBroker := fmt.Sprintf("%s://%s:%s@%s:%s",
		os.Getenv("BROKER_PROTOCOL"),
		os.Getenv("BROKER_USERNAME"),
		os.Getenv("BROKER_PASSWORD"),
		os.Getenv("BROKER_HOST"),
		os.Getenv("BROKER_PORT"),
	)
	eventSubsripter, err := broker.NewEventSubsripter(uriBroker)
	if err != nil {
		l.Info("failed to ping to rabbit: ", err.Error())
		return
	}
	// hasher
	h := hasher.NewHasher()
	base := database.NewDatabase(db)
	hashAdmin, err := h.HashPassword("admin")
	if err != nil {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	row, err := db.Query(context.Background(), "insert into users(username, email, hash_password) values ($1, $2, $3) RETURNING id;", "admin", "admin@admin.com", hashAdmin)
	var iduser int
	if err != nil {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	var idc10 int
	var idc11 int
	var idc12 int
	if row.Next() {
		row.Scan(&iduser)
	} else {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	row.Close()
	row, err = db.Query(context.Background(), "insert into controllers(type_controller, number_controller) values ($1, $2) RETURNING id_controller;", 10, 1)
	if err != nil {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	if row.Next() {
		row.Scan(&idc10)
	} else {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	row.Close()
	row, err = db.Query(context.Background(), "insert into controllers(type_controller, number_controller) values ($1, $2) RETURNING id_controller;", 11, 1)
	if err != nil {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	if row.Next() {
		row.Scan(&idc11)
	} else {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	row.Close()
	row, err = db.Query(context.Background(), "insert into controllers(type_controller, number_controller) values ($1, $2) RETURNING id_controller;", 12, 1)
	if err != nil {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	if row.Next() {
		row.Scan(&idc12)
	} else {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	row.Close()
	row, err = db.Query(context.Background(), "insert into user_controllers(id_user, id_controller) values ($1, $2);", iduser, idc10)
	if err != nil {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	row.Close()
	row, err = db.Query(context.Background(), "insert into user_controllers(id_user, id_controller) values ($1, $2);", iduser, idc11)
	if err != nil {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	row.Close()
	row, err = db.Query(context.Background(), "insert into user_controllers(id_user, id_controller) values ($1, $2);", iduser, idc12)
	if err != nil {
		l.Debug("faled data insert: %s", err.Error())
		return
	}
	row.Close()
	// services
	controllerService := service.NewControllerService(base)
	userService := service.NewUserService(base, *h)

	// controllers

	r := middleware.NewRouter(os.Getenv("API_HOST"), os.Getenv("API_PORT"), userService, controllerService, eventSubsripter, h)

	var ctrl []int
	rows, err := db.Query(context.Background(), "select id_controller from controllers;")
	if err != nil {
		l.Debug("faled data insert:", err.Error())
		return
	}
	for rows.Next() {
		var id_ctrl int
		rows.Scan(&id_ctrl)
		ctrl = append(ctrl, id_ctrl)
	}
	rows.Close()
	if err != nil {
		l.Debug("ошибка извлечения контролеров")
		l.Debug(err.Error())
		return
	}
	log.Info("количество контролеров", len(ctrl))
	/*
		var buffer chan []byte = make(chan []byte, 100)
		for i := 0; i < len(ctrl); i++ {
			go eventSubsripter.SubscribeMessange(context.Background(), strconv.Itoa(ctrl[i]), buffer)
		}
		go func() {
			for tmp := range buffer {
				l.Info("сообщение загружается в базу")
				var msg entitys.MessangeTypeZiroJson
				json.Unmarshal(tmp, &msg)
				l.Info("%s ", msg)
				if msg.One != nil {
					controllerService.CreateMessageTypeOne(context.Background(), *msg.One)
				}
				if msg.Two != nil {
					controllerService.CreateMessageTypeTwo(context.Background(), *msg.Two)
				}
				if msg.Three != nil {
					controllerService.CreateMessageTypeThree(context.Background(), *msg.Three)
				}
			}
		}()
	*/
	go func() {
		r.ListenAndServe()
	}()

	// groups

	l.Debug("success starting http server")

	l.Debug("success starting application")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	l.Info("application has been shut down")

}
