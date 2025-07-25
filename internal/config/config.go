package config

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"path"
)

type Config struct {
	App        `yaml:"app"`
	HTTPServer `yaml:"http-server"`
}

type App struct {
	Name    string `yaml:"name"`
	Profile string
	Logging
}

type Logging struct {
	Level  string `yaml:"level" env-default:"debug"`
	Format string `yaml:"format" env-default:"text"`
}

type HTTPServer struct {
	Port        string `yaml:"port" env-default:"8080"`
	PrintRoutes bool   `yaml:"print-routes" env-default:"true"`
}

func Load(env string) (*Config, error) {
	configDir := os.Getenv("APP_CONFIG_DIR")
	if configDir == "" {
		return nil, errors.New("undefined config directory")
	}

	slog.Info("Using config directory", "dir", configDir)
	configPath := path.Join(configDir, fmt.Sprintf("application.%s.yaml", env))

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func MustLoad(env string) *Config {
	value, err := Load(env)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return value
}
