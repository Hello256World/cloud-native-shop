package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kkyr/fig"
)

type Config struct {
	Environment    string         `fig:"environment" validate:"required"`
	PostgresConfig PostgresConfig `fig:"postgres" validate:"required"`
	RedisConfig    RedisConfig    `fig:"redis" validate:"required"`
	HTTPConfig     ServerConfig   `fig:"http" validate:"required"`
}

func (c Config) IsProduction() bool {
	return c.Environment == "production"
}

type PostgresConfig struct {
	Host        string `fig:"host" validate:"required"`
	Port        string `fig:"port" validate:"required"`
	Username    string `fig:"username"`
	Password    string `fig:"password"`
	Params      string `fig:"params"`
	Database    string `fig:"database" validate:"required"`
	AppName     string `fig:"appName" validate:"required"`
	MinPoolSize int    `fig:"minPoolSize" validate:"required"`
	MaxPoolSize int    `fig:"maxPoolSize" validate:"required"`
}

type RedisConfig struct {
	Host     string `fig:"host" validate:"required"`
	Port     int    `fig:"port" validate:"required"`
	Password string `fig:"password"`
	Database int    `fig:"database"`
}

func (c PostgresConfig) GetConnectionURI() string {
	userInfo := ""
	if c.Username != "" && c.Password != "" {
		userInfo = fmt.Sprintf(
			"%s:%s@",
			url.QueryEscape(c.Username),
			url.QueryEscape(c.Password),
		)
	}

	uri := fmt.Sprintf(
		"postgres://%s%s:%s/%s",
		userInfo,
		c.Host,
		c.Port,
		c.Database,
	)

	if c.Params != "" {
		uri += "?" + c.Params
	}

	log.Println(uri)

	return uri
}

func (r RedisConfig) GetConnectionURI() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type ServerConfig struct {
	Port int `fig:"port" validate:"required"`
}

func GetConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var config Config

	opts := []fig.Option{}

	if path := os.Getenv("CONFIG_PATH"); path != "" {
		opts = append(opts, fig.File(filepath.Base(path)), fig.Dirs(filepath.Dir(path)))
	}

	if err := fig.Load(&config, opts...); err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	if config.PostgresConfig.Username == "" {
		config.PostgresConfig.Username = os.Getenv("POSTGRES_USER")
	}

	if config.PostgresConfig.Password == "" {
		config.PostgresConfig.Password = os.Getenv("POSTGRES_PASSWORD")
	}

	return &config, nil
}
