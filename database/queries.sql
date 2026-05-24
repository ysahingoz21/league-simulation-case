SELECT id, name, strength, created_at
FROM teams
ORDER BY id;

SELECT
    m.id,
    m.week,
    ht.name AS home_team,
    at.name AS away_team,
    m.home_goals,
    m.away_goals,
    m.status,
    m.played_at
FROM matches AS m
JOIN teams AS ht ON ht.id = m.home_team_id
JOIN teams AS at ON at.id = m.away_team_id
ORDER BY m.week, m.id;

SELECT
    m.id,
    m.week,
    ht.name AS home_team,
    at.name AS away_team,
    m.home_goals,
    m.away_goals,
    m.status
FROM matches AS m
JOIN teams AS ht ON ht.id = m.home_team_id
JOIN teams AS at ON at.id = m.away_team_id
WHERE m.week = :week
ORDER BY m.id;

SELECT
    s.rank,
    t.name AS team_name,
    s.played,
    s.wins,
    s.draws,
    s.losses,
    s.goals_for,
    s.goals_against,
    s.goal_difference,
    s.points,
    s.updated_at
FROM standings AS s
JOIN teams AS t ON t.id = s.team_id
ORDER BY s.rank, t.name;

SELECT
    p.week,
    t.name AS team_name,
    p.championship_probability,
    p.expected_points,
    p.projected_rank,
    p.created_at
FROM predictions AS p
JOIN teams AS t ON t.id = p.team_id
WHERE p.week = (SELECT MAX(week) FROM predictions)
ORDER BY p.projected_rank, t.name;

SELECT
    m.id,
    m.week,
    ht.name AS home_team,
    at.name AS away_team,
    m.home_goals,
    m.away_goals,
    m.played_at
FROM matches AS m
JOIN teams AS ht ON ht.id = m.home_team_id
JOIN teams AS at ON at.id = m.away_team_id
WHERE m.status = 'played'
ORDER BY m.week, m.id;

SELECT
    m.id,
    m.week,
    ht.name AS home_team,
    at.name AS away_team,
    m.status
FROM matches AS m
JOIN teams AS ht ON ht.id = m.home_team_id
JOIN teams AS at ON at.id = m.away_team_id
WHERE m.status = 'scheduled'
ORDER BY m.week, m.id;

SELECT id, current_week, total_weeks, is_completed, created_at, updated_at
FROM league_state
WHERE id = 1;

UPDATE matches
SET
    home_goals = :home_goals,
    away_goals = :away_goals,
    status = 'played',
    played_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = :match_id;

DELETE FROM predictions;
DELETE FROM standings;
DELETE FROM matches;
DELETE FROM league_state;
DELETE FROM teams;

INSERT INTO league_state (id, current_week, total_weeks, is_completed)
VALUES (1, 0, 6, 0);
