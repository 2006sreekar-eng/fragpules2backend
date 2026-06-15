CREATE TABLE IF NOT EXISTS scores (
    id VARCHAR(36) PRIMARY KEY,
    player_name VARCHAR(100) NOT NULL,
    score INT NOT NULL,
    accuracy NUMERIC(5,2) NOT NULL,
    avg_reaction_time INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_scores_score_desc_created_at_asc 
ON scores (score DESC, created_at ASC);