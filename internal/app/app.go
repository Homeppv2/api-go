package app

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-go/internal/config"
	"api-go/internal/controller/middleware"
	"api-go/internal/infrastructure"
	"api-go/internal/service"
	"api-go/pkg/hasher"

	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config", err.Error())
	}

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Logger.Level}))
	l.Info("success initializing logger")

	l.Info("success initializing fiber")

	// db
	db, err := sqlx.Connect("pgx", cfg.Postgres.ConnString)
	if err != nil {
		l.Error("failed to connect to postgresql: ", err.Error())
		return
	}

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			l.Error(err.Error())
			return
		}
	}(db)

	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime * time.Second)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Postgres.ConnMaxIdleTime * time.Second)

	err = db.Ping()
	if err != nil {
		l.Error("failed to ping to postgresql: ", err.Error())
		return
	}

	l.Info("success connecting to postgresql")

	eventSubsripter, err := infrastructure.NewEventSubsripter("uri")
	if err != nil {
		l.Error("failed to ping to rabbit: ", err.Error())
		return
	}

	// hasher
	h := hasher.NewHasher()

	// infrastructures
	registryRepo := infrastructure.NewPGRegistry(db, l)
	userRepo := infrastructure.NewUserRepo(db, l)
	controllerRepo := infrastructure.NewControllerRepo(db, l)

	// services
	controllerService := service.NewControllerService(controllerRepo, userRepo)
	userService := service.NewUserService(registryRepo, *h)

	// controllers

	serverlogin := &http.Server{
		Addr: "/login",
		Handler: middleware.ServerLogin{
			Logf:             log.Printf,
			EventSubsripter:  eventSubsripter,
			UserService:      userService,
			ControlerService: controllerService,
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	serverregister := &http.Server{
		Addr: "/register",
		Handler: middleware.ServerRegister{
			Logf:        log.Printf,
			UserService: userService,
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	go func() {
		serverregister.ListenAndServe()
	}()
	go func() {
		serverlogin.ListenAndServe()
	}()

	// groups

	l.Debug("success starting http server")

	l.Debug("success starting application")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	l.Info("application has been shut down")

}
