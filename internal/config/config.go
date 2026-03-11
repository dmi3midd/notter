package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB         DBConfig         `yaml:"dbConfig"`
	HttpServer HttpServerConfig `yaml:"httpServer"`
	JWT        JWTConfig        `yaml:"jwt"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s, error: %v", configPath, err)
	}

	return &cfg
}
