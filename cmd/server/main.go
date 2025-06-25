package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freekobie/hazel/handlers"
	"github.com/freekobie/hazel/mail"
	"github.com/freekobie/hazel/postgres"
	"github.com/freekobie/hazel/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := loadConfig()

	db, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	mailer := mail.NewMailer(cfg.MailConfig)
	userService := services.NewUserService(postgres.NewUserStore(db), mailer)
	workspaceService := services.NewWorkspaceService(postgres.NewWorkspaceStore(db))

	handler := handlers.NewHandler(userService, workspaceService)

	app := newApplication(handler, cfg.ServerAddress)

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Starting server")
		serverErr <- app.start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			panic(err)
		}
	case sig := <-stop:
		slog.Info("Shutting down server", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.shutdown(ctx); err != nil {
			slog.Error("Graceful shutdown failed", "error", err)
		}
	}
}
