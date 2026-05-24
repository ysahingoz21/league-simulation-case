package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"league-simulation-case/internal/config"
	dbpkg "league-simulation-case/internal/database"
)

func TestLeagueRoutesInitializeAndListFixturesByWeek(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app-test.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	defer db.Close()

	if err := dbpkg.ApplySchema(db, filepath.Join("..", "..", "database", "schema.sql")); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	router := NewRouter(config.Config{AppEnv: "test", Port: "8080", DBPath: dbPath}, db)

	initRecorder := httptest.NewRecorder()
	initRequest := httptest.NewRequest(http.MethodPost, "/api/v1/league/init", nil)
	router.ServeHTTP(initRecorder, initRequest)

	if initRecorder.Code != http.StatusOK {
		t.Fatalf("expected init status 200, got %d body=%s", initRecorder.Code, initRecorder.Body.String())
	}

	var initResponse struct {
		Teams    []json.RawMessage `json:"teams"`
		Fixtures []json.RawMessage `json:"fixtures"`
	}
	if err := json.Unmarshal(initRecorder.Body.Bytes(), &initResponse); err != nil {
		t.Fatalf("decode init response: %v", err)
	}

	if len(initResponse.Teams) != 4 {
		t.Fatalf("expected 4 teams from init, got %d", len(initResponse.Teams))
	}

	if len(initResponse.Fixtures) != 12 {
		t.Fatalf("expected 12 fixtures from init, got %d", len(initResponse.Fixtures))
	}

	fixtureRecorder := httptest.NewRecorder()
	fixtureRequest := httptest.NewRequest(http.MethodGet, "/api/v1/fixtures/1", nil)
	router.ServeHTTP(fixtureRecorder, fixtureRequest)

	if fixtureRecorder.Code != http.StatusOK {
		t.Fatalf("expected fixture status 200, got %d body=%s", fixtureRecorder.Code, fixtureRecorder.Body.String())
	}

	var fixtureResponse struct {
		Week     int               `json:"week"`
		Fixtures []json.RawMessage `json:"fixtures"`
	}
	if err := json.Unmarshal(fixtureRecorder.Body.Bytes(), &fixtureResponse); err != nil {
		t.Fatalf("decode fixture response: %v", err)
	}

	if fixtureResponse.Week != 1 {
		t.Fatalf("expected fixture week 1, got %d", fixtureResponse.Week)
	}

	if len(fixtureResponse.Fixtures) != 2 {
		t.Fatalf("expected 2 fixtures for week 1, got %d", len(fixtureResponse.Fixtures))
	}
}

func TestLeagueRoutesRejectInvalidWeek(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app-invalid-week.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	defer db.Close()

	if err := dbpkg.ApplySchema(db, filepath.Join("..", "..", "database", "schema.sql")); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	router := NewRouter(config.Config{AppEnv: "test", Port: "8080", DBPath: dbPath}, db)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/fixtures/7", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestSimulationRoutesPlayNextWeek(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app-simulation.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	defer db.Close()

	if err := dbpkg.ApplySchema(db, filepath.Join("..", "..", "database", "schema.sql")); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	router := NewRouter(config.Config{AppEnv: "test", Port: "8080", DBPath: dbPath}, db)

	initRecorder := httptest.NewRecorder()
	initRequest := httptest.NewRequest(http.MethodPost, "/api/v1/league/init", nil)
	router.ServeHTTP(initRecorder, initRequest)
	if initRecorder.Code != http.StatusOK {
		t.Fatalf("expected init status 200, got %d body=%s", initRecorder.Code, initRecorder.Body.String())
	}

	playRecorder := httptest.NewRecorder()
	playRequest := httptest.NewRequest(http.MethodPost, "/api/v1/simulation/week/next", nil)
	router.ServeHTTP(playRecorder, playRequest)

	if playRecorder.Code != http.StatusOK {
		t.Fatalf("expected simulation status 200, got %d body=%s", playRecorder.Code, playRecorder.Body.String())
	}

	var playResponse struct {
		Week        int               `json:"week"`
		Matches     []json.RawMessage `json:"matches"`
		Standings   []json.RawMessage `json:"standings"`
		Predictions []json.RawMessage `json:"predictions"`
	}
	if err := json.Unmarshal(playRecorder.Body.Bytes(), &playResponse); err != nil {
		t.Fatalf("decode simulation response: %v", err)
	}

	if playResponse.Week != 1 {
		t.Fatalf("expected simulated week 1, got %d", playResponse.Week)
	}

	if len(playResponse.Matches) != 2 {
		t.Fatalf("expected 2 simulated matches, got %d", len(playResponse.Matches))
	}

	if len(playResponse.Standings) != 4 {
		t.Fatalf("expected 4 standings in simulation response, got %d", len(playResponse.Standings))
	}

	if len(playResponse.Predictions) != 0 {
		t.Fatalf("expected no predictions before week 4, got %d", len(playResponse.Predictions))
	}
}

func TestStandingsRouteReflectsPlayedWeek(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app-standings.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	defer db.Close()

	if err := dbpkg.ApplySchema(db, filepath.Join("..", "..", "database", "schema.sql")); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	router := NewRouter(config.Config{AppEnv: "test", Port: "8080", DBPath: dbPath}, db)

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/league/init", nil))
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/simulation/week/next", nil))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/standings", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected standings status 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Standings []struct {
			Played int `json:"played"`
		} `json:"standings"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode standings response: %v", err)
	}

	if len(response.Standings) != 4 {
		t.Fatalf("expected 4 standings, got %d", len(response.Standings))
	}

	for _, standing := range response.Standings {
		if standing.Played != 1 {
			t.Fatalf("expected each team to have played 1 match, got %d", standing.Played)
		}
	}
}

func TestSimulationRoutePlayAll(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app-play-all.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	defer db.Close()

	if err := dbpkg.ApplySchema(db, filepath.Join("..", "..", "database", "schema.sql")); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	router := NewRouter(config.Config{AppEnv: "test", Port: "8080", DBPath: dbPath}, db)

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/league/init", nil))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/simulation/play-all", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected play-all status 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		League struct {
			CurrentWeek int  `json:"current_week"`
			IsCompleted bool `json:"is_completed"`
		} `json:"league"`
		Weeks     []json.RawMessage `json:"weeks"`
		Standings []struct {
			Played int `json:"played"`
		} `json:"standings"`
		Predictions []json.RawMessage `json:"predictions"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode play-all response: %v", err)
	}

	if response.League.CurrentWeek != 6 || !response.League.IsCompleted {
		t.Fatalf("unexpected league state: %+v", response.League)
	}

	if len(response.Weeks) != 6 {
		t.Fatalf("expected 6 weeks in play-all response, got %d", len(response.Weeks))
	}

	if len(response.Standings) != 4 {
		t.Fatalf("expected 4 standings, got %d", len(response.Standings))
	}

	if len(response.Predictions) != 4 {
		t.Fatalf("expected 4 predictions after play-all, got %d", len(response.Predictions))
	}

	for _, standing := range response.Standings {
		if standing.Played != 6 {
			t.Fatalf("expected each team to have played 6 matches, got %d", standing.Played)
		}
	}
}

func TestPredictionRouteAtWeekFour(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "app-predictions.db")
	db, err := dbpkg.OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	defer db.Close()

	if err := dbpkg.ApplySchema(db, filepath.Join("..", "..", "database", "schema.sql")); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	router := NewRouter(config.Config{AppEnv: "test", Port: "8080", DBPath: dbPath}, db)
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/league/init", nil))
	for i := 0; i < 4; i++ {
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/simulation/week/next", nil))
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/predictions", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected predictions status 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Week        int               `json:"week"`
		Predictions []json.RawMessage `json:"predictions"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode predictions response: %v", err)
	}

	if response.Week != 4 {
		t.Fatalf("expected prediction week 4, got %d", response.Week)
	}

	if len(response.Predictions) != 4 {
		t.Fatalf("expected 4 predictions, got %d", len(response.Predictions))
	}
}
