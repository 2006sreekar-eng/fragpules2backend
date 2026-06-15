package repositories

import (
	"context"
	"database/sql"
	"fragpulse/internal/models"
)

type LeaderboardRepository interface {
	Create(ctx context.Context, score *models.Score) error
	GetAll(ctx context.Context) ([]models.Score, error)
	GetTop(ctx context.Context, limit int) ([]models.Score, error)
}

type PostgresLeaderboardRepository struct {
	db *sql.DB
}

func NewPostgresLeaderboardRepository(db *sql.DB) *PostgresLeaderboardRepository {
	return &PostgresLeaderboardRepository{db: db}
}

func (r *PostgresLeaderboardRepository) Create(ctx context.Context, s *models.Score) error {
	query := `INSERT INTO scores (id, player_name, score, accuracy, avg_reaction_time, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.PlayerName, s.Score, s.Accuracy, s.AvgReactionTime, s.CreatedAt)
	return err
}

func (r *PostgresLeaderboardRepository) GetAll(ctx context.Context) ([]models.Score, error) {
	query := `SELECT id, player_name, score, accuracy, avg_reaction_time, created_at 
			  FROM scores ORDER BY score DESC, created_at ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.Score
	for rows.Next() {
		var s models.Score
		if err := rows.Scan(&s.ID, &s.PlayerName, &s.Score, &s.Accuracy, &s.AvgReactionTime, &s.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, nil
}

func (r *PostgresLeaderboardRepository) GetTop(ctx context.Context, limit int) ([]models.Score, error) {
	query := `SELECT id, player_name, score, accuracy, avg_reaction_time, created_at 
			  FROM scores ORDER BY score DESC, created_at ASC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.Score
	for rows.Next() {
		var s models.Score
		if err := rows.Scan(&s.ID, &s.PlayerName, &s.Score, &s.Accuracy, &s.AvgReactionTime, &s.CreatedAt); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, nil
}