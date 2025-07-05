package config

import (
	// this will automatically load .env file:
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type PullInterval struct {
	Baseball time.Duration `env:"PULL_INTERVAL_BASEBALL" env-default:"1s"`
	Soccer   time.Duration `env:"PULL_INTERVAL_SOCCER" env-default:"2s"`
	Football time.Duration `env:"PULL_INTERVAL_FOOTBALL" env-default:"3s"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     int    `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

type HttpConfig struct {
	Port int    `env:"HTTP_PORT" env-default:"8080"`
	Host string `env:"HTTP_HOST" env-default:"localhost"`
}

type GrpcConfig struct {
	Port int    `env:"GRPC_PORT" env-default:"8081"`
	Host string `env:"GRPC_HOST" env-default:"localhost"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" env-default:"debug"`
}

func getAddr(host string, port int) string {
	return fmt.Sprintf("%s:%v", host, port)
}

func (c *GrpcConfig) Addr() string {
	return getAddr(c.Host, c.Port)
}

func (c *HttpConfig) Addr() string {
	return getAddr(c.Host, c.Port)
}

func (c *LinesProviderConfig) Addr() string {
	return getAddr(c.Host, c.Port)
}

type LinesProviderConfig struct {
	Port int    `env:"LINES_PROVIDER_PORT" env-default:"8000"`
	Host string `env:"LINES_PROVIDER_HOST" env-default:"localhost"`
}

type Config struct {
	Http          HttpConfig
	Grpc          GrpcConfig
	PullInterval  PullInterval
	LinesProvider LinesProviderConfig
	Redis         RedisConfig
	Logger        LoggerConfig
}

func InitConfig() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic("Read .env file error")
	}

	return cfg
}
