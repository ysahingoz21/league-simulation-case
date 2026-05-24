PRAGMA foreign_keys = ON;

INSERT INTO teams (name, strength) VALUES
    ('Türkiye', 95),
    ('ABD', 85),
    ('Avustralya', 80),
    ('paraguay', 75);

INSERT INTO league_state (id, current_week, total_weeks, is_completed)
VALUES (1, 0, 6, 0);
