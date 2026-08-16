package main

import (
	router "core/api/router"
	config "core/internal/config"
	database "core/internal/database"
	logger "core/pkg/log"
	openai "core/pkg/openai"

	"net/http"
	"time"

	"github.com/ory/graceful"
)

// @title Armin Lens API
// @version 1.0
// @BasePath /api/v1
func main() {
	logger.Start()

	logger.Info("🧩 Configuring environment")
	cfg, err := config.Load()
	logger.Fatal(err, "Failed to load env")

	logger.Info("📡 Bootstrapping database")
	gormDB := database.Init(&cfg.Database)
	sqlDB, err := gormDB.DB()
	logger.Fatal(err, "Failed to create connection with database")
	defer sqlDB.Close()

	server := graceful.WithDefaults(
		&http.Server{
			Addr:         cfg.Server.Port,
			Handler:      router.Config(gormDB),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	)

	err = database.Migrate(&cfg.Database)
	logger.Fatal(err, "Failed to create migrations")

	logger.Info("🤖 OpenAI setup")
	openai.Configure(&cfg.OpenAI)

	logger.Info("🚀 Starting server on port %s", cfg.Server.Port)
	err = graceful.Graceful(
		server.ListenAndServe,
		server.Shutdown,
	)

	logger.Fatal(err, "Failed to gracefully shutdown")
	logger.Info("🛑 Server exited")
}
