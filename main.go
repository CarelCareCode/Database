package main

import (
	"emergency-response-backend/internal/config"
	"emergency-response-backend/internal/database"
	"emergency-response-backend/internal/redis"
	"emergency-response-backend/internal/server"
	"emergency-response-backend/internal/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize server
	srv := server.New(cfg, db, redisClient, wsHub)

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
}
