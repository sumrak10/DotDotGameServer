package main

import (
	"OnlineGame/database"
	"OnlineGame/game/world"
	"OnlineGame/manager"
	"OnlineGame/server"
	"fmt"
	"log"
)

func main() {
	_ = database.GetDB()
	fmt.Println("DB initialized...")

	_ = world.GetPresetVault()
	fmt.Println("Preset Vault initialized...")

	_ = manager.GetManager()
	fmt.Println("Manager initialized...")

	s := server.NewServer()
	log.Fatal(s.Start(":8080"))
}
