package handlers

import (
	"encoding/json"
	"log" // Added for error tracking
	"net/http"
	"strconv"
	"fragpulse/internal/models"
	"fragpulse/internal/services"
)

type LeaderboardHandler struct {
	service *services.LeaderboardService
}

func NewLeaderboardHandler(service *services.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

func (h *LeaderboardHandler) CreateScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid route action", http.StatusMethodNotAllowed)
		return
	}
	var req models.CreateScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	score, err := h.service.AddScore(r.Context(), req)
	if err != nil {
		log.Printf("ERROR: AddScore failed: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(score)
}

func (h *LeaderboardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	scores, err := h.service.GetAllScores(r.Context())
	if err != nil {
		log.Printf("ERROR: GetLeaderboard failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(scores)
}

func (h *LeaderboardHandler) GetTopLeaderboard(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	scores, err := h.service.GetTopScores(r.Context(), limit)
	if err != nil {
		log.Printf("ERROR: GetTopLeaderboard failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(scores)
}