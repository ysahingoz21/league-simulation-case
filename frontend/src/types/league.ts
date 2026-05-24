export type NullableNumber = number | null
export type NullableString = string | null

export type LeagueState = {
  current_week: number
  total_weeks: number
  is_completed: boolean
  created_at?: string
  updated_at?: string
}

export type Team = {
  id: number
  name: string
  strength: number
  created_at: string
}

export type Match = {
  id: number
  week: number
  home_team_id: number
  home_team_name: string
  away_team_id: number
  away_team_name: string
  home_goals: NullableNumber
  away_goals: NullableNumber
  status: 'scheduled' | 'played'
  played_at: NullableString
  created_at: string
  updated_at: string
}

export type Standing = {
  team_id: number
  team_name: string
  played: number
  wins: number
  draws: number
  losses: number
  goals_for: number
  goals_against: number
  goal_difference: number
  points: number
  rank: number
}

export type Prediction = {
  id?: number
  week: number
  team_id: number
  team_name: string
  championship_probability: number
  expected_points: number
  projected_rank: number
  created_at?: string
}

export type WeekResult = {
  week: number
  matches: Match[]
}

export type HealthResponse = {
  status: string
  message: string
  environment: string
}

export type LeagueInitResponse = {
  league: LeagueState
  teams: Team[]
  fixtures: Match[]
  standings: Standing[]
}

export type LeagueResponse = {
  league: LeagueState
}

export type TeamsResponse = {
  teams: Team[]
}

export type FixturesResponse = {
  week?: number
  fixtures: Match[]
}

export type MatchResponse = {
  match: Match
}

export type MatchUpdateResponse = {
  match: Match
  league: LeagueState
  standings: Standing[]
  predictions: Prediction[]
}

export type SimulationWeekResponse = {
  week: number
  league: LeagueState
  matches: Match[]
  standings: Standing[]
  predictions: Prediction[]
}

export type PlayAllResponse = {
  league: LeagueState
  weeks: WeekResult[]
  standings: Standing[]
  predictions: Prediction[]
}

export type StandingsResponse = {
  standings: Standing[]
}

export type PredictionsResponse = {
  week: number
  message?: string
  predictions: Prediction[]
}

export type ApiErrorResponse = {
  error: string
}
