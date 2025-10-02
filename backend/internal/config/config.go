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
	MinioConfig    `yaml:"minio" env-required:"true"`
}

type RestConfig struct {
	Port        int    `yaml:"port" env-required:"true"`
	AllowOrigin string `yaml:"allow_origin" env-required:"true"`
}

type PostgresConfig struct {
	User     string `yaml:"user" env-required:"true" env:"POSTGRES_USER"`
	Password string `yaml:"password" env-required:"true" env:"POSTGRES_PASSWORD"`
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Database string `yaml:"database" env-required:"true" env:"POSTGRES_DB"`
}

type LoggerConfig struct {
	Prefix string `yaml:"prefix" env-required:"false"`
	Flag   int    `yaml:"flag" env-required:"false"`
}

type MinioConfig struct {
	Endpoint       string `yaml:"endpoint" env-required:"true"`
	RootUser       string `yaml:"root-user" env-required:"true" env:"MINIO_ROOT_USER"`
	RootPassword   string `yaml:"root-password" env-required:"true" env:"MINIO_ROOT_PASSWORD"`
	BucketName     string `yaml:"bucket-name" env-required:"true"`
	PublicEndpoint string `yaml:"public-endpoint" env-required:"true"`
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

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	return &cfg
}
