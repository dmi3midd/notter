package api

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/dmi3midd/notter/internal/config"
	"github.com/dmi3midd/notter/internal/handlers"
	"github.com/dmi3midd/notter/internal/middlewares"
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

	noteRepo := repositories.NewNoteRepo(s.db)
	noteService := services.NewNoteService(noteRepo)
	noteHandler := handlers.NewNoteHandler(noteService)
	noteRouter := routers.NewNoteRouter(noteHandler)

	mainRouter.Mount("/api/users", userRouter)

	mainRouter.Group(func(r chi.Router) {
		r.Use(middlewares.Authorization(tokenService, userRepo))

		r.Mount("/api/notes", noteRouter)
	})

	// mainRouter.Mount("/api/notes", noteRouter)
	return mainRouter
}
