package config

import (
	// this will automatically load .env file:
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

// ADDR_HTTP=localhost:8080
// ADDR_GRPC=localhost:8081
// PULL_INTERVAL_BASEBALL=1s
// PULL_INTERVAL_SOCCER=2s
// PULL_INTERVAL_FOOTBALL=3s
// LOG_LEVEL=debug

type PullInterval struct {
	Baseball time.Duration `env:"PULL_INTERVAL_BASEBALL" env-default:"1s"`
	Soccer   time.Duration `env:"PULL_INTERVAL_SOCCER" env-default:"2s"`
	Football time.Duration `env:"PULL_INTERVAL_BASEBALL" env-default:"3s"`
}

// type Config struct {
// 	PullIntervals PullInterval
// }

type Config struct {
	HttpAddr     string `env:"ADDR_HTTP" env-default:"localhost:8080"`
	GrpcAddr     string `env:"ADDR_GRPC" env-default:"localhost:8081"`
	PullInterval PullInterval
	// PullIntervalBaseball time.Duration `env:"PULL_INTERVAL_BASEBALL" env-default:"1s"`
	// PullIntervalSoccer   time.Duration `env:"PULL_INTERVAL_SOCCER" env-default:"2s"`
	// PullIntervalFootball time.Duration `env:"PULL_INTERVAL_BASEBALL" env-default:"3s"`
}

func InitConfig() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic("Read .env file error")
	}

	return cfg
}
