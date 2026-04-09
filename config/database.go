package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
)

var (
	instanceDatabaseConfig *DatabaseConfig
	onceDatabaseConfig     sync.Once
)

type DatabaseConfig struct {
	Path string `env:"DB_PATH" envDefault:":memory:"`
}

func Database() *DatabaseConfig {
	onceDatabaseConfig.Do(func() {
		instanceDatabaseConfig = &DatabaseConfig{}
		if err := env.Parse(instanceDatabaseConfig); err != nil {
			log.Fatalf("Fatal: %v", err)
		}
	})
	return instanceDatabaseConfig
}
