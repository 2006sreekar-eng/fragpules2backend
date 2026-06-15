package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fragpulse/internal/config"
	"fragpulse/internal/database"
	"fragpulse/internal/handlers"
	"fragpulse/internal/repositories"
	"fragpulse/internal/routes"
	"fragpulse/internal/services"
	"fragpulse/internal/websocket"
)

func main() {
	log.Println("Starting FragPulse Core Architecture Engine...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config load failure: %v", err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Postgres link crash: %v", err)
	}
	defer db.Close()

	wsManager := websocket.NewManager()
	go wsManager.Run()

	lbRepo := repositories.NewPostgresLeaderboardRepository(db)
	lbService := services.NewLeaderboardService(lbRepo, wsManager)
	lbHandler := handlers.NewLeaderboardHandler(lbService)

	appRouter := routes.RegisterRoutes(cfg, lbHandler, wsManager)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      appRouter,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server core loop failure: %v", err)
		}
	}()

	log.Printf("FragPulse Microservice running locally on port %s", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully draining connections...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}