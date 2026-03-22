package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dmi3midd/notter/internal/config"
	"github.com/dmi3midd/notter/internal/db"
	"github.com/joho/godotenv"
)

type testUser struct {
	Username string
	Email    string
	Password string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
	}

	cfg := config.LoadConfig()
	database, err := db.NewDB(
		cfg.GetDsn(),
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	users := []testUser{
		{Username: "test_user1", Email: "test1@example.com", Password: "password1234"},
		{Username: "test_user2", Email: "test2@example.com", Password: "password4321"},
	}

	PlantUsers(database, users)

	fmt.Println("Seed planted!")
	os.Exit(0)
}
