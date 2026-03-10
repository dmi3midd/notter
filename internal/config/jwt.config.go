package config

import "time"

type JWTConfig struct {
	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
	AccessExpiry       time.Duration
	RefreshExpiry      time.Duration
}
