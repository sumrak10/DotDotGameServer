package main

import (
	"OnlineGame/internal/config"
	"OnlineGame/internal/database"
	"OnlineGame/internal/game/world"
	"OnlineGame/internal/manager"
	"OnlineGame/internal/server"
	"fmt"
	"log"
	//"net/http"
	//_ "net/http/pprof"
)

func main() {
	//go func() {
	//	http.ListenAndServe(":6060", nil)
	//}()

	config.InitAll()
	fmt.Println("Config initialized...")

	_ = database.GetDB()
	fmt.Println("DB initialized...")

	_ = world.GetPresetVault()
	fmt.Println("Preset Vault initialized...")

	_ = manager.GetManager()
	fmt.Println("Manager initialized...")

	s := server.NewServer()
	log.Fatal(s.Start(fmt.Sprintf("0.0.0.0:%d", config.Server().ServerPort)))
}
