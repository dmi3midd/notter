package config

import "time"

type JWTConfig struct {
	AccessSecret  string        `yaml:"accessSecret"`
	RefreshSecret string        `yaml:"refreshSecret"`
	AccessExpiry  time.Duration `yaml:"accessExpiry"`
	RefreshExpiry time.Duration `yaml:"refreshExpiry"`
}
