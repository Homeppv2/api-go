package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-go/internal/config"
	"api-go/internal/controller/http"
	"api-go/internal/controller/middleware"
	"api-go/internal/infrastructure"
	"api-go/internal/service"
	"api-go/pkg/hasher"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	f := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	f.Use(cors.New(cors.Config{
		AllowHeaders: "*",
	}))

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

	if cfg.Postgres.AutoMigrate {

		migrationDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			l.Error("failed to migrate to postgresql: ", err.Error())
			return
		}

		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", cfg.Postgres.MigrationsPath),
			"user",
			migrationDriver,
		)
		if err != nil {
			l.Error("failed to migrate to postgresql: ", err.Error())
			return
		}

		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			l.Error("failed to migrate to postgresql: ", err.Error())
			return
		}
	}

	l.Info("success connecting to postgresql")

	// redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil || pong != "PONG" {
		l.Error("failed to ping to redis: ", err.Error())
		return
	}

	l.Info("success connecting to redis")

	// hasher
	h := hasher.NewHasher()

	// infrastructures
	registryRepo := infrastructure.NewPGRegistry(db, l)
	sessionRepo := infrastructure.NewSessionRepo(redisClient, l, cfg.Auth.Token.ExpiresIn)
	userRepo := infrastructure.NewUserRepo(db, l)
	controllerRepo := infrastructure.NewControllerRepo(db, l)

	// services
	authService := service.NewAuthService(sessionRepo, userRepo, *h)
	controllerService := service.NewControllerService(controllerRepo, userRepo)
	userService := service.NewUserService(registryRepo, *h)

	// controllers
	middlewareManager := middleware.NewMiddlewareManager(authService)
	authHandler := http.NewAuthHandler(authService, *middlewareManager)
	userHandler := http.NewUserHandler(userService, *middlewareManager)
	controllerHandler := http.NewControllerHandler(controllerService, *middlewareManager)

	// groups
	apiGroup := f.Group("api")
	authGroup := apiGroup.Group("auth")
	usersGroup := apiGroup.Group("users")
	controllersGroup := apiGroup.Group("controllers")

	authHandler.Register(authGroup)
	userHandler.Register(usersGroup)
	controllerHandler.Register(controllersGroup)

	go func() {
		err = f.Listen(net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port))
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	l.Debug("success starting http server")

	l.Debug("success starting application")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	l.Info("application has been shut down")

}
