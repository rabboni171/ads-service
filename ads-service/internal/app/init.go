package app

import (
	"ads-service/internal/config"
	"ads-service/internal/controllers"
	"ads-service/internal/repository"
	"ads-service/internal/service"
	"ads-service/logger"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	_ "github.com/lib/pq"
)

func Initialization() {
	// Initialize config
	config.MustLoad()

	// Initialize logger
	logger.MustLoad()

	// New context
	ctx := context.Background()

	// Initialize repository
	repository := repository.NewRepository(ctx)

	// Initialize controllers
	router := controllers.New(ctx, services.NewService(repository))

	// Create the HTTP server
	serverAddr := config.Cfg.HTTPServerAddress
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		slog.Info("starting server at port: " + serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("unable to start server", "err", err.Error())
		}
	}()

	// Set up channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	slog.Info("Shutting down server...")

	// Create a context with a timeout for the shutdown
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctxShutDown); err != nil {
		slog.Error("Server forced to shutdown", "err", err.Error())
	} else {
		slog.Info("Server exited gracefully")
	}

}
