package routes

import (
	"net/http"
	"fragpulse/internal/config"
	"fragpulse/internal/handlers"
	"fragpulse/internal/middleware"
	"fragpulse/internal/websocket"
)

func RegisterRoutes(cfg *config.Config, lbHandler *handlers.LeaderboardHandler, wsManager *websocket.Manager) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/api/scores", lbHandler.CreateScore)
	mux.HandleFunc("/api/leaderboard", lbHandler.GetLeaderboard)
	mux.HandleFunc("/api/leaderboard/top", lbHandler.GetTopLeaderboard)
	mux.HandleFunc("/ws", wsManager.ServeWS)

	return middleware.CORS(cfg)(mux)
}