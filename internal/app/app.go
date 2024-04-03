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

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run() {

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	l.Info("success initializing logger")

	l.Info("success initializing fiber")

	// db
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("POSTGRESQL_USER"), os.Getenv("POSTGRESQL_PASSWORD"), os.Getenv("POSTGRESQL_HOST"), os.Getenv("POSTGRESQL_PORT"), os.Getenv("POSTGRESQL_BASE")))
	if err != nil {
		l.Error("failed to connect to postgresql: ", err.Error())
		return
	}
	defer func(db *pgxpool.Pool) {
		db.Close()
	}(db)

	err = db.Ping(context.Background())
	if err != nil {
		l.Error("failed to ping to postgresql: ", err.Error())
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
		l.Error("failed to ping to rabbit: ", err.Error())
		return
	}

	// hasher
	h := hasher.NewHasher()
	base := database.NewDatabase(db)
	hashAdmin, err := h.HashPassword("admin")
	tx, err := db.Begin(context.Background())
	row, err := tx.Query(context.Background(), "insert into users (username, email, hash_password) values ($1, $2, $3) RETURNING id;", "admin", "admin@admin.com", hashAdmin)
	var iduser int
	var idc10 int
	var idc11 int
	var idc12 int
	if row.Next() {
		row.Scan(&iduser)
	}
	row, err = tx.Query(context.Background(), "insert into contollers(type_controller, number_controller) values ($1, $2) RETURNING id_contorller;", 10, 1)
	if row.Next() {
		row.Scan(&idc10)
	}
	row, err = tx.Query(context.Background(), "insert into contollers(type_controller, number_controller) values ($1, $2) RETURNING id_contorller;", 11, 1)
	if row.Next() {
		row.Scan(&idc11)
	}
	row, err = tx.Query(context.Background(), "insert into contollers(type_controller, number_controller) values ($1, $2) RETURNING id_contorller;", 12, 1)
	if row.Next() {
		row.Scan(&idc12)
	}
	tx.Query(context.Background(), "insert into user_conttollers(id_user, id_contorller) values ($1, $2);", iduser, idc10)
	tx.Query(context.Background(), "insert into user_conttollers(id_user, id_contorller) values ($1, $2);", iduser, idc11)
	tx.Query(context.Background(), "insert into user_conttollers(id_user, id_contorller) values ($1, $2);", iduser, idc12)

	tx.Commit(context.Background())

	// services
	controllerService := service.NewControllerService(base)
	userService := service.NewUserService(base, *h)

	// controllers

	r := middleware.NewRouter(os.Getenv("API_HOST"), os.Getenv("API_PORT"), userService, controllerService, eventSubsripter, h)

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
