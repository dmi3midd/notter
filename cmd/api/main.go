package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/dmi3midd/notter/internal/api"
	"github.com/dmi3midd/notter/internal/config"
	"github.com/dmi3midd/notter/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	cfg := config.LoadConfig()

	logger := setupLoger("txt")
	logger.Info("logger is active")

	db, err := db.NewDB(
		cfg.GetDsn(),
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("database connection established")

	server := api.NewServer(cfg, db, logger)
	log.Fatal(server.Start())

}

func setupLoger(lf string) *slog.Logger {
	var logger *slog.Logger
	switch lf {
	case "json":
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{},
			),
		)
	case "txt":
		logger = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{},
			),
		)
	}
	return logger
}
