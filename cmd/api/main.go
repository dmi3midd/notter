package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/dmi3midd/notter/internal/config"
	"github.com/dmi3midd/notter/internal/db"

	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	cfg := config.LoadConfig()

	// logs := setupLoger("txt")
	// logs.Info("logger is active")

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

	mux := http.NewServeMux()
	log.Fatal(start(mux, cfg))

}

func start(mux *http.ServeMux, cfg *config.Config) error {
	srv := http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      mux,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	log.Printf("server is running on %s", os.Getenv("ADDR"))
	return srv.ListenAndServe()
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
