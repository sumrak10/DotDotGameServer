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
	ServerSSL  bool   `env:"SERVER_SSL,required"`
	TPS        uint   `env:"SERVER_TPS,required"`

	BindHost string `env:"BIND_HOST" envDefault:"0.0.0.0"`
	BindPort uint   `env:"BIND_PORT" envDefault:"8080"`
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

func (s *ServerConfig) GetBindAddress() string {
	return fmt.Sprintf("%s:%d", s.BindHost, s.BindPort)
}

func (s *ServerConfig) GetBaseURL() string {
	var result string
	if s.ServerSSL {
		result = fmt.Sprintf("https://%s:%d", s.ServerHost, s.ServerPort)
	} else {
		result = fmt.Sprintf("http://%s:%d", s.ServerHost, s.ServerPort)
	}
	return result
}

func (s *ServerConfig) GetWSBaseURL() string {
	var result string
	if s.ServerSSL {
		result = fmt.Sprintf("wss://%s:%d", s.ServerHost, s.ServerPort)
	} else {
		result = fmt.Sprintf("ws://%s:%d", s.ServerHost, s.ServerPort)
	}
	return result
}
