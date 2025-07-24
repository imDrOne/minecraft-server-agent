package config

import (
	_ "embed"
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
}

type HTTPServer struct {
	Port string `yaml:"port"`
}

func New(env string) *Config {
	configDir := os.Getenv("APP_CONFIG_DIR")
	if configDir == "" {
		slog.Error("Undefined config directory", "env", env)
		os.Exit(1)
	}

	slog.Info("Using config directory", "dir", configDir)
	configPath := path.Join(configDir, fmt.Sprintf("application.%s.yaml", env))

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("Failed to parse config file", "error", err, "env", env)
		os.Exit(1)
	}

	return &cfg
}
