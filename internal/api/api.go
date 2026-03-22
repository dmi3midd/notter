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
	// repositories
	tokenRepo := repositories.NewTokenRepo(s.db)
	userRepo := repositories.NewUserRepo(s.db)
	noteRepo := repositories.NewNoteRepo(s.db)
	boardRepo := repositories.NewBoardRepo(s.db)

	// services
	tokenService := services.NewTokenService(s.cfg.JWT, tokenRepo)
	userService := services.NewUserService(userRepo, tokenService)
	noteService := services.NewNoteService(noteRepo, userRepo, boardRepo)
	boardService := services.NewBoardRepo(boardRepo, userRepo)

	// handlers
	userHandler := handlers.NewUserHandler(userService)
	noteHandler := handlers.NewNoteHadnler(noteService)
	boardHandler := handlers.NewBoardHandler(boardService)

	// routers
	mainRouter := chi.NewRouter()
	userRouter := routers.NewUserRouter(userHandler)
<<<<<<< HEAD
	noteRouter := routers.NewNoteRouter(noteHandler)
	boardRouter := routers.NewBoardRouter(boardHandler)

	// routers setup
	mainRouter.Mount("/api/users", userRouter)
	mainRouter.Group(func(router chi.Router) {
		router.Use(middlewares.Authorization(tokenService, userRepo))

		router.Mount("/api/notes", noteRouter)
		router.Mount("/api/boards", boardRouter)
	})

=======
	mainRouter.Mount("/api/users", userRouter)
>>>>>>> master
	return mainRouter
}
