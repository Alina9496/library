// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/Alina9496/library/config"
	"github.com/Alina9496/library/internal/api"
	"github.com/Alina9496/library/internal/repo"
	"github.com/Alina9496/library/internal/service"
	"github.com/Alina9496/tool/pkg/httpserver"
	"github.com/Alina9496/tool/pkg/logger"
	"github.com/Alina9496/tool/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	runDatabaseMigration(cfg.PG.URL)

	// Repository
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(err)
	}
	defer pg.Close()

	// Use case
	service := service.New(
		repo.New(pg, l),
		l,
	)

	// HTTP Server
	handler := gin.New()
	api.NewServer(handler, l, service)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
