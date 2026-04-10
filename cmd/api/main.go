package main

import (
	"log"

	"github.com/dmi3midd/notter/internal/api"
	"github.com/dmi3midd/notter/internal/config"
	"github.com/dmi3midd/notter/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	cfg := config.LoadConfig()

	dbService, err := database.New(&cfg.DB)
	if err != nil {
		log.Panic(err)
	}
	defer dbService.Close()
	log.Println("database connection established")

	server := api.NewServer(cfg, dbService.GetDB())
	log.Fatal(server.Start())

}
