package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config
	Config struct {
		App   `yaml:"app"`
		DB    `yaml:"db"`
		Redis `yaml:"redis"`
		Auth  `yaml:"auth"`
	}

	// App
	App struct {
		Name    string `env-required:"true" yaml:"name"`
		Version string `env-required:"true" yaml:"version"`
		Env     string `env-required:"true" yaml:"env" env:"ENV"`
	}

	// DB
	DB struct {
		URI  string `env-required:"true" yaml:"uri"`
		Name string `env-required:"true" yaml:"name"`
	}

	// Redis
	Redis struct {
		Address  string `env-required:"true" yaml:"address"`
		DBNumber int    `yaml:"db_number"`
		Password string `yaml:"password"`
	}

	// Auth
	Auth struct {
		JwtSecret string `env-required:"true" yaml:"jwt_secret"`
		JwtTTL    int    `env-required:"true" yaml:"jwt_ttl"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configFileName(), cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func configFileName() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "development"
	}
	return fmt.Sprintf("./config/config.%s.yml", env)
}
