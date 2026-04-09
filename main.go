package main

import (
	"OnlineGame/config"
	"OnlineGame/database"
	"OnlineGame/game/world"
	"OnlineGame/manager"
	"OnlineGame/server"
	"fmt"
	"log"
)

func main() {
	config.InitAll()
	fmt.Println("Config initialized...")

	_ = database.GetDB()
	fmt.Println("DB initialized...")

	_ = world.GetPresetVault()
	fmt.Println("Preset Vault initialized...")

	_ = manager.GetManager()
	fmt.Println("Manager initialized...")

	s := server.NewServer()
	log.Fatal(s.Start(":8080"))
}
