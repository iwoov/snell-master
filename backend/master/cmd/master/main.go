package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/iwoov/snell-master/pkg/database"
	"github.com/iwoov/snell-master/backend/master/internal/api"
	"github.com/iwoov/snell-master/backend/master/internal/api/handler"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
	"github.com/iwoov/snell-master/backend/master/internal/scheduler"
	"github.com/iwoov/snell-master/backend/master/internal/service"
	"github.com/iwoov/snell-master/pkg/config"
	"github.com/iwoov/snell-master/pkg/logger"
)

func main() {
	configPath := flag.String("config", "backend/master/configs/master.example.yaml", "path to the master config file")
	migrationsDir := flag.String("migrations", "backend/master/migrations", "directory with migration files")
	flag.Parse()

	startTime := time.Now()
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logInstance, err := logger.Init(cfg.Log)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}

	db, err := database.InitDB(cfg.Database)
	if err != nil {
		logInstance.Fatalf("init database: %v", err)
	}

	if err := database.RunMigrations(cfg.Database.Path, *migrationsDir); err != nil {
		logInstance.Fatalf("run migrations: %v", err)
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.ServiceDeps{
		Repositories: repos,
		Logger:       logInstance,
		Config:       cfg,
		DB:           db,
	})
	handlers := handler.NewHandlers(services, db, startTime)

	manager := scheduler.NewManager()
	manager.Add(scheduler.ScheduleDailyReset(repos.User, logInstance))
	manager.Add(scheduler.ScheduleMonthlyReset(repos.User, logInstance))
	manager.Add(scheduler.ScheduleHealthCheck(db, logInstance, 30*time.Second))

	engine := api.SetupRouter(cfg, handlers, services.Log, db)

	server := &http.Server{
		Addr:    cfg.Server.Address(),
		Handler: engine,
	}

	logInstance.Infof("master service listening on %s", cfg.Server.Address())
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logInstance.Fatalf("start server: %v", err)
		}
	}()

	waitForShutdown(server, manager, logInstance, db)
}
