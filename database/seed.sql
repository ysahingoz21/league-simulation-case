PRAGMA foreign_keys = ON;

INSERT INTO teams (name, strength) VALUES
    ('Turkey', 95),
    ('USA', 85),
    ('Australia', 80),
    ('Paraguay', 75);

INSERT INTO league_state (id, current_week, total_weeks, is_completed)
VALUES (1, 0, 6, 0);
