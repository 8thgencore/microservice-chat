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

// Env represents the environment in which the application is running.
type Env string

const (
	// Local environment.
	Local Env = "local"
	// Dev environment.
	Dev Env = "dev"
	// Prod environment.
	Prod Env = "prod"
)

// Config represents the configuration for the application.
type Config struct {
	Env        Env `env:"ENV" env-default:"local"`
	GRPC       GRPC
	TLS        TLSConfig
	Database   DatabaseConfig
	AuthClient AuthClient
}

// GRPC represents the configuration for the GRPC server.
type GRPC struct {
	Host      string        `env:"GRPC_SERVER_HOST" env-default:"localhost"`
	Port      int           `env:"GRPC_SERVER_PORT" env-default:"50051"`
	Transport string        `env:"GRPC_SERVER_TRANSPORT" env-default:"tcp"`
	Timeout   time.Duration `env:"GRPC_SERVER_TIMEOUT"`
}

// Address returns the address of the GRPC server in the format "host:port".
func (c *GRPC) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// DatabaseConfig represents the configuration for the Postgres database.
type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST"     env-required:"true"`
	Port     string `env:"POSTGRES_PORT"     env-required:"true"`
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Name     string `env:"POSTGRES_DB"       env-required:"true"`
}

// DSN returns the data source name (DSN) for the database
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Name, c.User, c.Password)
}

// AuthClient represents a client for authenticating users.
type AuthClient struct {
	Host     string `env:"AUTH_CLIENT_HOST" env-default:"auth"`
	Port     int    `env:"AUTH_CLIENT_PORT" env-default:"50052"`
	CertPath string `env:"AUTH_CERT_PATH"`
}

// Address returns the address of the authentication server in the format "host:port".
func (c *AuthClient) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// TLSConfig represents the configuration for the TLSConfig.
type TLSConfig struct {
	CertPath string `env:"TLS_CERT_PATH"`
	KeyPath  string `env:"TLS_KEY_PATH"`
}

// NewConfig creates a new instance of Config
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
