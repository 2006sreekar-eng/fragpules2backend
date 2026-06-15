package services

import (
	"context"
	"errors"
	"strings"
	"time"
	"fragpulse/internal/models"
	"fragpulse/internal/repositories"
	"fragpulse/internal/websocket"
	"github.com/google/uuid"
)

type LeaderboardService struct {
	repo      repositories.LeaderboardRepository
	wsManager *websocket.Manager
}

func NewLeaderboardService(repo repositories.LeaderboardRepository, wsManager *websocket.Manager) *LeaderboardService {
	return &LeaderboardService{repo: repo, wsManager: wsManager}
}

func (s *LeaderboardService) AddScore(ctx context.Context, req models.CreateScoreRequest) (*models.Score, error) {
	if strings.TrimSpace(req.PlayerName) == "" {
		return nil, errors.New("invalid player name")
	}

	score := &models.Score{
		ID:              uuid.New().String(),
		PlayerName:      strings.TrimSpace(req.PlayerName),
		Score:           req.Score,
		Accuracy:        req.Accuracy,
		AvgReactionTime: req.AvgReactionTime,
		CreatedAt:       time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, score); err != nil {
		return nil, err
	}

	topScores, err := s.repo.GetTop(ctx, 10)
	if err == nil {
		s.wsManager.BroadcastEvent("leaderboard_update", topScores)
	}
	return score, nil
}

func (s *LeaderboardService) GetAllScores(ctx context.Context) ([]models.Score, error) {
	return s.repo.GetAll(ctx)
}

func (s *LeaderboardService) GetTopScores(ctx context.Context, limit int) ([]models.Score, error) {
	return s.repo.GetTop(ctx, limit)
}