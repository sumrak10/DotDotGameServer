package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
)

var (
	instanceJWTConfig *JWTConfig
	onceJWTConfig     sync.Once
)

type JWTConfig struct {
	Secret string `env:"JWT_SECRET,required"`
}

func JWT() *JWTConfig {
	onceJWTConfig.Do(func() {
		instanceJWTConfig = &JWTConfig{}
		if err := env.Parse(instanceJWTConfig); err != nil {
			log.Fatalf("Fatal: %v", err)
		}
	})
	return instanceJWTConfig
}
