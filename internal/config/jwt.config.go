package config

import "time"

type JWTConfig struct {
	JWT_ACCESS_SECRET  string        `yaml:"JWT_ACCESS_SECRET"`
	JWT_REFRESH_SECRET string        `yaml:"JWT_REFRESH_SECRET"`
	AccessExpiry       time.Duration `yaml:"accessExpiry"`
	RefreshExpiry      time.Duration `yaml:"refreshExpiry"`
}
