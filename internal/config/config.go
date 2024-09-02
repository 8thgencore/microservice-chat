package config

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Env string

const (
	Local Env = "local"
	Dev   Env = "dev"
	Prod  Env = "prod"
)

type Config struct {
	Env        Env `env:"ENV" env-default:"local"`
	GRPC       GRPC
	Database   DatabaseConfig
	AuthClient AuthClient
}

type GRPC struct {
	Host      string        `env:"GRPC_SERVER_HOST" env-default:"localhost"`
	Port      int           `env:"GRPC_SERVER_PORT" env-default:"50051"`
	Transport string        `env:"GRPC_SERVER_TRANSPORT" env-default:"tcp"`
	Timeout   time.Duration `env:"GRPC_SERVER_TIMEOUT"`
}

func (c *GRPC) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST"     env-required:"true"`
	Port     string `env:"POSTGRES_PORT"     env-required:"true"`
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Name     string `env:"POSTGRES_DB"       env-required:"true"`
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Name, c.User, c.Password)
}

type AuthClient struct {
	Host string `env:"AUTH_CLIENT_HOST" env-default:"auth"`
	Port int    `env:"AUTH_CLIENT_PORT" env-default:"50052"`
}

func (c *AuthClient) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func NewConfig() (*Config, error) {
	configPath := fetchConfigPath()

	cfg := &Config{}
	var err error

	if configPath != "" {
		err = godotenv.Load(configPath)
	} else {
		err = godotenv.Load()
	}
	if err != nil {
		log.Printf("No loading .env file: %v", err)
	}

	if err = cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("error reading env: %w", err)
	}

	return cfg, nil
}

func fetchConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", ".env", "Path to config file")

	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
