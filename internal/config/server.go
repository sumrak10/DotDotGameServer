package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
)

var (
	instanceServerConfig *ServerConfig
	onceServerConfig     sync.Once
)

type ServerConfig struct {
	ServerHost string `env:"SERVER_HOST,required"`
	ServerPort uint   `env:"SERVER_PORT,required"`
	TPS        uint   `env:"SERVER_TPS,required"`
}

func Server() *ServerConfig {
	onceServerConfig.Do(func() {
		instanceServerConfig = &ServerConfig{}
		if err := env.Parse(instanceServerConfig); err != nil {
			log.Fatalf("Fatal: %v", err)
		}
	})
	return instanceServerConfig
}

func (s *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.ServerHost, s.ServerPort)
}
