package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/scheduler"
)

func waitForShutdown(server *http.Server, manager *scheduler.Manager, logger *logrus.Logger, db *gorm.DB) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if logger != nil {
		logger.Info("shutdown signal received, stopping services...")
	}
	if manager != nil {
		manager.Stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if server != nil {
		if err := server.Shutdown(ctx); err != nil && logger != nil {
			logger.WithError(err).Error("server shutdown failed")
		}
	}

	if db != nil {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	if logger != nil {
		logger.Info("shutdown completed")
	}
}
