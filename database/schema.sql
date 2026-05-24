PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    strength INTEGER NOT NULL CHECK (strength BETWEEN 1 AND 100),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS league_state (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    current_week INTEGER NOT NULL DEFAULT 0,
    total_weeks INTEGER NOT NULL DEFAULT 6,
    is_completed INTEGER NOT NULL DEFAULT 0 CHECK (is_completed IN (0, 1)),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    week INTEGER NOT NULL CHECK (week BETWEEN 1 AND 6),
    home_team_id INTEGER NOT NULL,
    away_team_id INTEGER NOT NULL,
    home_goals INTEGER NULL CHECK (home_goals IS NULL OR home_goals >= 0),
    away_goals INTEGER NULL CHECK (away_goals IS NULL OR away_goals >= 0),
    status TEXT NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'played')),
    played_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (home_team_id) REFERENCES teams(id),
    FOREIGN KEY (away_team_id) REFERENCES teams(id),
    CHECK (home_team_id <> away_team_id),
    CHECK (
        (status = 'scheduled' AND home_goals IS NULL AND away_goals IS NULL)
        OR (status = 'played' AND home_goals IS NOT NULL AND away_goals IS NOT NULL)
    ),
    UNIQUE (week, home_team_id, away_team_id)
);

CREATE TABLE IF NOT EXISTS standings (
    team_id INTEGER PRIMARY KEY REFERENCES teams(id),
    played INTEGER NOT NULL DEFAULT 0 CHECK (played >= 0),
    wins INTEGER NOT NULL DEFAULT 0 CHECK (wins >= 0),
    draws INTEGER NOT NULL DEFAULT 0 CHECK (draws >= 0),
    losses INTEGER NOT NULL DEFAULT 0 CHECK (losses >= 0),
    goals_for INTEGER NOT NULL DEFAULT 0 CHECK (goals_for >= 0),
    goals_against INTEGER NOT NULL DEFAULT 0 CHECK (goals_against >= 0),
    goal_difference INTEGER NOT NULL DEFAULT 0,
    points INTEGER NOT NULL DEFAULT 0 CHECK (points >= 0),
    rank INTEGER NOT NULL DEFAULT 0 CHECK (rank >= 0),
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS predictions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    week INTEGER NOT NULL CHECK (week BETWEEN 4 AND 6),
    team_id INTEGER NOT NULL REFERENCES teams(id),
    championship_probability REAL NOT NULL CHECK (championship_probability BETWEEN 0 AND 100),
    expected_points REAL NOT NULL CHECK (expected_points >= 0),
    projected_rank REAL NOT NULL CHECK (projected_rank BETWEEN 1 AND 4),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_matches_week ON matches (week);
CREATE INDEX IF NOT EXISTS idx_matches_status ON matches (status);
CREATE INDEX IF NOT EXISTS idx_predictions_week ON predictions (week);
CREATE INDEX IF NOT EXISTS idx_standings_rank ON standings (rank);
