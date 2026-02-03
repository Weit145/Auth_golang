package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `yaml:"env" env-default:"local"`
	GRPC     Grpc   `yaml:"grpc"`
	JWT      JWT
	TokenTTL TokenTTL `yaml:"token_ttl"`
}

type Grpc struct {
	Address string `yaml::"address" env-default:"auth-service:50051"`
}

type JWT struct {
	Secret    string `env:"SECRET_JWT" env-required:"true"`
	Algorithm string `env:"ALGORITHM_JWT" env-required:"true"`
}

type TokenTTL struct {
	Access  time.Duration `yaml:"access" env-default:"1h"`
	Refresh time.Duration `yaml:"refresh" env-default:"72h"`
}

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("cannot read .env file: %s", err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// log.Fatal("CONFIG_PATH is not set")
		configPath = "config/local.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read env: %s", err)
	}

	return &cfg
}
