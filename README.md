# World Cup Group Simulator

A full-stack Go + Vue application for simulating a 4-team football group, updating standings, predicting championship probabilities, and editing match results.

---

## Overview

The application simulates a 4-team football group stage over 6 weeks:

- The **backend** initializes teams and fixtures, simulates weekly matches, recalculates standings, generates championship predictions after week 4, and supports result editing via a REST API.
- The **frontend** provides a live dashboard to control and visualize the simulation: run weeks, browse the standings table, inspect match cards, view predictions, and edit any played result.
- Both surfaces are independently usable — the backend is fully operable via `curl` or Postman without the frontend.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend language | Go 1.25 |
| HTTP framework | Gin |
| Database | SQLite via `database/sql` |
| Architecture | Handler → Service → Repository → Domain |
| Frontend framework | Vue 3 + Vite + TypeScript |
| State management | Pinia |
| Styling | Plain CSS (no UI library) |
| Tooling | npm, curl / Postman, SQLite CLI (optional) |

---

## Why Vue 3 and Pinia?

- **Vue 3** is the current stable Vue option and integrates naturally with Vite for fast development. Its component model maps cleanly to the dashboard panels.
- **Pinia** is the recommended state management library for Vue 3, replacing Vuex. It is simpler, fully typed, and works without boilerplate.
- All league state (`league`, `teams`, `fixtures`, `standings`, `predictions`, `loading`, `error`) lives in a single Pinia store. Components call store actions — they never make direct API calls — which keeps data flow predictable and components thin.

---

## Features

**Backend**
- League initialization with 4 default teams and configurable strengths
- Double round-robin fixture generation (6 weeks, 2 matches per week)
- Weekly simulation using team strength-weighted random outcomes
- Play a specific week or all remaining weeks in one request
- Premier League-style standings: points → goal difference → goals for → name
- Championship predictions via Monte Carlo simulation (available after week 4)
- Deterministic final predictions when the league is completed
- Match result editing with automatic standings and prediction recalculation
- Full league reset
- SQLite schema with CHECK constraints and foreign keys

**Frontend**
- Live dashboard with sidebar navigation and sticky top bar
- League controls: initialize, reset, play next week, play all remaining
- Standings table with team flags, rank highlighting, and live badge
- Match Center with per-week tabs and football scoreboard cards
- Season Predictions panel with probability bars and expected points
- Edit Match Result modal with score inputs and validation
- Team logo support with initials fallback
- Smooth section scroll navigation

**Testing**
- Backend test suites across all packages (domain, repository, service, app integration)
- All routes covered: initialization, simulation, standings, predictions, match patch
- Frontend Vite build validation

---

## League Rules

| Rule | Value |
|---|---|
| Teams | 4 |
| Total weeks | 6 |
| Matches per week | 2 |
| Format | Double round-robin (each pair plays twice, home/away reversed) |
| Win | 3 points |
| Draw | 1 point |
| Loss | 0 points |

**Ranking tiebreaker order:**
1. Points (descending)
2. Goal difference (descending)
3. Goals for (descending)
4. Team name (ascending, alphabetical fallback)

---

## Prediction Logic

Predictions are generated after week 4 because that is when enough match data exists to make meaningful probability estimates for the remaining weeks.

**Algorithm:**
- Monte Carlo simulation over all unplayed remaining fixtures
- Each simulation run samples match outcomes weighted by team strength
- Aggregates across 10,000 runs to produce:
  - `championship_probability` — percentage of runs where the team finished rank 1
  - `expected_points` — average final points
  - `projected_rank` — average final rank

**Completed league:**
- When all 6 weeks are played, predictions become deterministic
- Rank 1 team receives 100% championship probability
- All other teams receive 0%

---

## Match Editing Logic

Any played match can have its result updated via `PATCH /api/v1/matches/:id` or through the frontend Edit Result modal.

After an edit:
1. The match score is updated in the database
2. Standings are recalculated from scratch using all currently played matches
3. If predictions exist (week ≥ 4), predictions are regenerated
4. If the league is completed, final deterministic predictions are regenerated

---

## Project Structure

```
.
├── cmd/api/                     # Backend entrypoint (main.go)
├── internal/
│   ├── app/                     # Gin router setup, integration tests
│   ├── config/                  # Environment variable loading
│   ├── database/                # SQLite connection, schema migration
│   ├── domain/                  # Core models, standings helpers, match logic
│   ├── handler/                 # HTTP handlers (thin, delegates to service)
│   ├── repository/              # Repository interfaces + SQLite implementations
│   └── service/                 # Business logic workflows
├── database/
│   ├── schema.sql               # Table definitions with constraints
│   ├── seed.sql                 # Default teams and league state
│   └── queries.sql              # Reference queries for manual inspection
├── frontend/
│   ├── public/                  # Static assets (app logo, team flag images)
│   └── src/
│       ├── api/                 # Typed HTTP client (wraps fetch)
│       ├── components/          # Vue dashboard components
│       ├── stores/              # Pinia store (all league state)
│       ├── styles/              # Global CSS
│       ├── types/               # TypeScript type definitions
│       └── utils/               # Team color hash, initials, logo mapping
└── README.md
```

---

## Backend Architecture

The backend follows a strict layered architecture:

| Layer | Responsibility |
|---|---|
| `handler` | Parse HTTP request, call service, write JSON response |
| `service` | Coordinate business logic across repositories |
| `repository` | Persist and query data; defined as interfaces |
| `domain` | Database-independent models and pure helper functions |

**Key decisions:**
- Repository interfaces allow services to be tested with mock implementations without touching SQLite.
- The domain layer has no dependencies on any infrastructure package — it is pure Go.
- SQLite was chosen for zero-setup local review. The schema is applied on startup via the `database` package.
- Standings are always recalculated from scratch (not incrementally) to ensure correctness after edits.

---

## Frontend Architecture

The frontend is a single-page dashboard with no routing:

- **Pinia store** is the single source of truth for all league state: `league`, `teams`, `fixtures`, `standings`, `predictions`, `loading`, and `error`.
- **Components** are stateless consumers — they read from the store and dispatch actions. No component makes a direct API call.
- **TeamAvatar** is a reusable component that renders a team's flag image if a mapping exists, or falls back to initials inside a colored circle.
- **Plain CSS** with CSS custom properties (variables) provides a consistent design system without adding a UI framework dependency.

---

## Requirements

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Node.js | 18+ |
| npm | 9+ |
| SQLite CLI | Optional (for manual DB inspection) |
| jq | Optional (for pretty curl output) |

---

## Environment Variables

**Backend**

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `APP_ENV` | `development` | Environment label |
| `DB_PATH` | `league.db` | SQLite file path |

**Frontend**

| Variable | Default | Description |
|---|---|---|
| `VITE_API_BASE_URL` | `http://localhost:8080` | Backend base URL |

**Examples:**

```bash
PORT=9090 DB_PATH=local.db go run cmd/api/main.go
```

```bash
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

---

## Setup and Run

**Backend**

```bash
go mod tidy
go run cmd/api/main.go
```

The server starts on `http://localhost:8080`. The SQLite schema is applied automatically on first run.

**Frontend**

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:5173` in a browser.

**Production build**

```bash
cd frontend
npm run build
```

---

## API Reference

Base URL: `http://localhost:8080`

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/league/init` | Initialize league (create teams + fixtures) |
| `POST` | `/api/v1/league/reset` | Reset league to initial state |
| `GET` | `/api/v1/league` | Get current league state |
| `GET` | `/api/v1/teams` | List all teams |
| `GET` | `/api/v1/standings` | Get current standings |
| `GET` | `/api/v1/fixtures` | List all fixtures |
| `GET` | `/api/v1/fixtures/:week` | List fixtures for a specific week |
| `GET` | `/api/v1/matches/:id` | Get a single match |
| `PATCH` | `/api/v1/matches/:id` | Edit a played match result |
| `GET` | `/api/v1/predictions` | Get latest championship predictions |
| `POST` | `/api/v1/simulation/week/next` | Simulate the next week |
| `POST` | `/api/v1/simulation/week/:week` | Simulate a specific week |
| `POST` | `/api/v1/simulation/play-all` | Simulate all remaining weeks |

---

## curl Demo Flow

Run these in order after starting the backend:

```bash
# 1. Confirm the backend is running
curl -s http://localhost:8080/health | jq

# 2. Reset to a clean state (safe to run on first start or re-run)
curl -s -X POST http://localhost:8080/api/v1/league/reset | jq

# 3. Initialize the league
curl -s -X POST http://localhost:8080/api/v1/league/init | jq

# 5. Check league state
curl -s http://localhost:8080/api/v1/league | jq

# 6. List teams
curl -s http://localhost:8080/api/v1/teams | jq

# 7. View all fixtures (note the match IDs returned here)
curl -s http://localhost:8080/api/v1/fixtures | jq

# 8. Play week 1
curl -s -X POST http://localhost:8080/api/v1/simulation/week/next | jq

# 9. Check standings
curl -s http://localhost:8080/api/v1/standings | jq

# 10. Play weeks 2, 3, 4
curl -s -X POST http://localhost:8080/api/v1/simulation/week/next | jq
curl -s -X POST http://localhost:8080/api/v1/simulation/week/next | jq
curl -s -X POST http://localhost:8080/api/v1/simulation/week/next | jq

# 11. Get predictions (available from week 4)
curl -s http://localhost:8080/api/v1/predictions | jq

# 12. Edit a played match result
# Match IDs depend on database state. Use GET /api/v1/fixtures to find a valid played match ID.
curl -s -X PATCH http://localhost:8080/api/v1/matches/MATCH_ID \
  -H "Content-Type: application/json" \
  -d '{"home_goals": 3, "away_goals": 0}' | jq

# 13. Confirm standings updated
curl -s http://localhost:8080/api/v1/standings | jq

# 14. Play all remaining weeks
curl -s -X POST http://localhost:8080/api/v1/simulation/play-all | jq

# 15. Final predictions
curl -s http://localhost:8080/api/v1/predictions | jq

# 16. Reset for another run
curl -s -X POST http://localhost:8080/api/v1/league/reset | jq
```

---

## Database Schema

Five tables — no ORM:

| Table | Purpose |
|---|---|
| `teams` | Team names and strength values |
| `league_state` | Single-row state: current week, total weeks, completed flag |
| `matches` | Fixtures with optional goals; enforces `scheduled`/`played` consistency via CHECK |
| `standings` | Denormalized standings row per team, recalculated on every play/edit |
| `predictions` | Per-week Monte Carlo results, replaced on every prediction refresh |

**Manual inspection:**

```bash
sqlite3 league.db

.tables
SELECT * FROM teams;
SELECT rank, name FROM standings JOIN teams ON teams.id = standings.team_id ORDER BY rank;
SELECT week, championship_probability FROM predictions JOIN teams ON teams.id = predictions.team_id ORDER BY projected_rank;
.quit
```

Reference queries are also available in `database/queries.sql`.

---

## Tests

```bash
go test ./...
```

All four packages pass:

| Package | Coverage |
|---|---|
| `internal/app` | Integration tests for all HTTP routes |
| `internal/domain` | Unit tests for standings sort, rank assignment, match logic |
| `internal/repository/sqlite` | Repository tests against an in-memory SQLite instance |
| `internal/service` | Service-level tests with mock repositories |

Frontend build validation:

```bash
cd frontend && npm run build
```
