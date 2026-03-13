package api

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/dmi3midd/notter/internal/config"
	"github.com/dmi3midd/notter/internal/handlers"
	"github.com/dmi3midd/notter/internal/repositories"
	"github.com/dmi3midd/notter/internal/routers"
	"github.com/dmi3midd/notter/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	cfg    *config.Config
	db     *sqlx.DB
	logger *slog.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, logger *slog.Logger) *Server {
	return &Server{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

func (s *Server) Start() error {
	mux := s.setupRoutes()
	srv := &http.Server{
		Addr:         s.cfg.HttpServer.Address,
		Handler:      mux,
		WriteTimeout: s.cfg.HttpServer.WriteTimeout,
		ReadTimeout:  s.cfg.HttpServer.ReadTimeout,
		IdleTimeout:  s.cfg.HttpServer.IdleTimeout,
	}
	log.Printf("server is running on %s", s.cfg.HttpServer.Address)
	return srv.ListenAndServe()
}

func (s *Server) setupRoutes() *chi.Mux {
	mainRouter := chi.NewRouter()

	tokenRepo := repositories.NewTokenRepo(s.db)
	tokenService := services.NewTokenService(s.cfg.JWT, tokenRepo)

	userRepo := repositories.NewUserRepo(s.db)
	userService := services.NewUserService(userRepo, tokenService)
	userHandler := handlers.NewUserHandler(userService)
	userRouter := routers.NewUserRouter(userHandler)
	mainRouter.Mount("/users", userRouter)
	return mainRouter
}
