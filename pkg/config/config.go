// Package config provides application configuration settings.
package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	// this will automatically load .env file:.
	_ "github.com/joho/godotenv/autoload"
)

// PullInterval defines pull intervals for different sports from external service.
type PullInterval struct {
	Baseball time.Duration `env:"PULL_INTERVAL_BASEBALL" env-default:"1s"`
	Soccer   time.Duration `env:"PULL_INTERVAL_SOCCER"   env-default:"2s"`
	Football time.Duration `env:"PULL_INTERVAL_FOOTBALL" env-default:"3s"`
}

// RedisConfig defines configuration for redis.
type RedisConfig struct {
	Host     string `env:"REDIS_HOST"     env-default:"localhost"`
	Port     int    `env:"REDIS_PORT"     env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB"       env-default:"0"`
}

// HTTPConfig defines http server configuration for line sync checking.
type HTTPConfig struct {
	Port int    `env:"HTTP_PORT" env-default:"8080"`
	Host string `env:"HTTP_HOST" env-default:"localhost"`
}

// GrpcConfig defines grpc server configuration of coefficients streaming.
type GrpcConfig struct {
	Port int    `env:"GRPC_PORT" env-default:"8081"`
	Host string `env:"GRPC_HOST" env-default:"localhost"`
}

// LoggerConfig defines http configuration for line sync checking.
type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" env-default:"debug"`
}

func getAddr(host string, port int) string {
	return fmt.Sprintf("%s:%v", host, port)
}

// Addr returns full address.
func (c *GrpcConfig) Addr() string {
	return getAddr(c.Host, c.Port)
}

// Addr returns full address.
func (c *HTTPConfig) Addr() string {
	return getAddr(c.Host, c.Port)
}

// Addr returns full address.
func (c *LinesProviderConfig) Addr() string {
	return getAddr(c.Host, c.Port)
}

// LinesProviderConfig defines request configuration for coefficients pulling from external service.
type LinesProviderConfig struct {
	Port int    `env:"LINES_PROVIDER_PORT" env-default:"8000"`
	Host string `env:"LINES_PROVIDER_HOST" env-default:"localhost"`
}

// Config defines all configurations.
type Config struct {
	HTTP          HTTPConfig
	Grpc          GrpcConfig
	PullInterval  PullInterval
	LinesProvider LinesProviderConfig
	Redis         RedisConfig
	Logger        LoggerConfig
}

// InitConfig initialize configuration.
func InitConfig() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic("Read .env file error")
	}

	return cfg
}
