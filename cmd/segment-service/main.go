package main

import (
	"avito-internship/internal/config"
	"avito-internship/internal/lib/logger/handlers/slogpretty"
	"avito-internship/internal/lib/logger/slogger"
	"net/http"

	"avito-internship/internal/storage/postgres"
	"os"

	"avito-internship/internal/http-server/handlers/segments/del"
	"avito-internship/internal/http-server/handlers/segments/save"
	delsegments "avito-internship/internal/http-server/handlers/users/del_segments"
	getactiveseg "avito-internship/internal/http-server/handlers/users/get-active-seg"
	"avito-internship/internal/http-server/handlers/users/save/saveuser"
	save_seg_user "avito-internship/internal/http-server/handlers/users/save_seg_user"
	mwLogger "avito-internship/internal/http-server/middleware/logger"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustConfigLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting segment service", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	storage, err := postgres.New(cfg.PostgresPath)
	if err != nil {
		log.Error("failed to init storage", slogger.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)

	// Create and delete segments
	router.Post("/segment", save.New(log, storage))
	router.Delete("/segment/{id}", del.DelSeg(log, storage))

	router.Post("/users", saveuser.New(log, storage))

	// Get active users segments, save segments to user, delete segments from user
	router.Get("/users/{id}/segments", getactiveseg.GetActiveSegmentsForUser(log, storage))
	router.Post("/users/{id}/segments", save_seg_user.AddUserToSegments(log, storage))
	router.Delete("/segment/{segmentName}", delsegments.DelSeg(log, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {

	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
