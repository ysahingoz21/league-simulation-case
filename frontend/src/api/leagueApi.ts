import { request } from './client'
import type {
  FixturesResponse,
  HealthResponse,
  LeagueInitResponse,
  LeagueResponse,
  MatchResponse,
  MatchUpdateResponse,
  PlayAllResponse,
  PredictionsResponse,
  SimulationWeekResponse,
  StandingsResponse,
  TeamsResponse,
} from '../types/league'

export function getHealth() {
  return request<HealthResponse>('/health')
}

export function initLeague() {
  return request<LeagueInitResponse>('/api/v1/league/init', {
    method: 'POST',
  })
}

export function resetLeague() {
  return request<LeagueInitResponse>('/api/v1/league/reset', {
    method: 'POST',
  })
}

export function getLeague() {
  return request<LeagueResponse>('/api/v1/league')
}

export function getTeams() {
  return request<TeamsResponse>('/api/v1/teams')
}

export function getFixtures() {
  return request<FixturesResponse>('/api/v1/fixtures')
}

export function getFixturesByWeek(week: number) {
  return request<FixturesResponse>(`/api/v1/fixtures/${week}`)
}

export function getMatch(id: number) {
  return request<MatchResponse>(`/api/v1/matches/${id}`)
}

export function updateMatchResult(id: number, homeGoals: number, awayGoals: number) {
  return request<MatchUpdateResponse>(`/api/v1/matches/${id}`, {
    method: 'PATCH',
    body: JSON.stringify({
      home_goals: homeGoals,
      away_goals: awayGoals,
    }),
  })
}

export function playNextWeek() {
  return request<SimulationWeekResponse>('/api/v1/simulation/week/next', {
    method: 'POST',
  })
}

export function playWeek(week: number) {
  return request<SimulationWeekResponse>(`/api/v1/simulation/week/${week}`, {
    method: 'POST',
  })
}

export function playAll() {
  return request<PlayAllResponse>('/api/v1/simulation/play-all', {
    method: 'POST',
  })
}

export function getStandings() {
  return request<StandingsResponse>('/api/v1/standings')
}

export function getPredictions() {
  return request<PredictionsResponse>('/api/v1/predictions')
}
