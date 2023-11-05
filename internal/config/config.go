package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env"`
	LogsPath   string `yaml:"logs_path"`
	StatesPath string `yaml:"states_path"`
	HTTPServer `yaml:"http_server"`
	LifeConfig `yaml:"life_config"`
}

type HTTPServer struct {
	Addres      string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type LifeConfig struct {
	Height int `yaml:"height"`
	Width  int `yaml:"width"`
	Fill   int `yaml:"fill"`
}

func MustLoad() *Config {
	config_path := os.Getenv("CONFIG_PATH")
	if config_path == "" {
		log.Fatal("CONFIG_PATH is no set")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(config_path, &cfg); err != nil {
		log.Fatalf("Cannot read config: %v\n", err)
	}

	return &cfg
}
