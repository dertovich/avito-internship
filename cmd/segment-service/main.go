package main

import (
	"avito-internship/internal/config"
	"avito-internship/internal/lib/logger/handlers/slogpretty"
	"avito-internship/internal/lib/logger/slogger"
	"net/http"

	//"avito-internship/internal/storage"
	"avito-internship/internal/storage/postgres"
	"os"

	"avito-internship/internal/http-server/handlers/segments/save"
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

	router.Post("/segment", save.New(log, storage))

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
