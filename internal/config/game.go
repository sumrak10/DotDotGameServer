package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
)

var (
	instanceGameConfig *GameConfig
	onceGameConfig     sync.Once
)

type GameConfig struct {
	ValuesScaleCoef uint `env:"GAME_VALUES_SCALE_COEF,required"`
	ArmySpeed       uint `env:"GAME_ARMY_SPEED,required"`
}

func Game() *GameConfig {
	onceGameConfig.Do(func() {
		instanceGameConfig = &GameConfig{}
		if err := env.Parse(instanceGameConfig); err != nil {
			log.Fatalf("Fatal: %v", err)
		}
	})
	return instanceGameConfig
}
