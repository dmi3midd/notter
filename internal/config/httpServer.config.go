package config

import "time"

type HttpServerConfig struct {
	Address      string        `yaml:"address"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
}
