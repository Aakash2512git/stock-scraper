package main

import (
	"context"
	"fmt"
	_ "main/docs"
	"main/internal/redis"
	"main/internal/routes"
	"main/internal/scraper"
	"main/pkg/config"
	logger "main/pkg/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/swagger"
)

func main() {
	log := logger.GetLogger()
	conf := config.GetConfig()

	// Create cancellable context for scraper
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Redis
	redis.InitRedis(ctx)

	// Start scraper in background
	go scraper.Run(ctx, "24h") // or your cron expression

	// Set up Fiber server
	server := routes.SetUpRoutes()
	port := ":" + conf.AppPort
	fmt.Println(port + "<-this is my port")

	server.Get("/swagger/*", swagger.HandlerDefault)
	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run server in goroutine
	go func() {
		if err := server.Listen(port); err != nil {
			log.Error("Error starting server: ", err)
		}
	}()

	// Wait for signal
	<-quit
	fmt.Println("\nShutting down gracefully...")

	// Stop scraper
	cancel()
	if err := server.Shutdown(); err != nil {
		log.Error("Fiber shutdown error: ", err)
	}

	fmt.Println("All services stopped")
}
