package config

import (
	"log/slog"
	"time"
)

type (
	Config struct {
		Logger   Logger
		HTTP     HTTP
		Postgres Postgres
		Redis    Redis
		Auth     Auth
	}

	Logger struct {
		Level slog.Level
	}

	HTTP struct {
		Host string
		Port string
	}

	Postgres struct {
		ConnString      string
		MaxOpenConns    int
		ConnMaxLifetime time.Duration
		MaxIdleConns    int
		ConnMaxIdleTime time.Duration
		AutoMigrate     bool
		MigrationsPath  string
	}

	Redis struct {
		Host     string `koanf:"host" validate:"required"`
		Port     string `koanf:"port" validate:"required"`
		Password string `koanf:"password"`
		DB       int    `koanf:"db"`
	}

	Auth struct {
		Token Token
	}

	Token struct {
		ExpiresIn time.Duration
	}
)

func LoadConfig() (*Config, error) {

	defaultLogLevel := slog.LevelInfo

	cfg := &Config{
		HTTP: HTTP{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Logger: Logger{
			Level: defaultLogLevel,
		},
		Postgres: Postgres{
			ConnString:      "postgresql://root:pass@localhost:5432/api?sslmode=disable",
			MaxOpenConns:    10,
			ConnMaxLifetime: 20,
			MaxIdleConns:    15,
			ConnMaxIdleTime: 30,
			AutoMigrate:     true,
			MigrationsPath:  "db/migration",
		},
		Redis: Redis{
			Host:     "127.0.0.1",
			Port:     "6380",
			Password: "",
			DB:       0,
		},
		Auth: Auth{
			Token: Token{
				ExpiresIn: 24 * time.Hour,
			},
		},
	}

	return cfg, nil
}
