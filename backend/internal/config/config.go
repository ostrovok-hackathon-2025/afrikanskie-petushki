package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LoggerConfig   `yaml:"logger" env-required:"true"`
	RestConfig     `yaml:"rest" env-required:"true"`
	PostgresConfig `yaml:"postgres" env-required:"true"`
}

type RestConfig struct {
	Port        int    `yaml:"port" env-required:"true"`
	AllowOrigin string `yaml:"allow_origin" env-required:"true"`
}

type PostgresConfig struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Database string `yaml:"database" env-required:"true"`
}

type LoggerConfig struct {
	Prefix string `yaml:"prefix" env-required:"false"`
	Flag   int    `yaml:"flag" env-required:"false"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	log.Printf(configPath)

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	return &cfg
}
