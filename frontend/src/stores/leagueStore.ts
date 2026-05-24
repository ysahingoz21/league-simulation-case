import { defineStore } from 'pinia'

import * as leagueApi from '../api/leagueApi'
import type {
  LeagueState,
  Match,
  Prediction,
  Standing,
  Team,
} from '../types/league'

type LeagueStoreState = {
  league: LeagueState | null
  teams: Team[]
  fixtures: Match[]
  standings: Standing[]
  predictions: Prediction[]
  selectedWeek: number | null
  selectedMatch: Match | null
  loading: boolean
  error: string | null
}

function mergeMatchIntoFixtures(fixtures: Match[], updatedMatch: Match) {
  const existingIndex = fixtures.findIndex((fixture) => fixture.id === updatedMatch.id)
  if (existingIndex === -1) {
    return [...fixtures, updatedMatch].sort((left, right) => {
      if (left.week !== right.week) {
        return left.week - right.week
      }
      return left.id - right.id
    })
  }

  const nextFixtures = [...fixtures]
  nextFixtures[existingIndex] = updatedMatch
  return nextFixtures
}

export const useLeagueStore = defineStore('league', {
  state: (): LeagueStoreState => ({
    league: null,
    teams: [],
    fixtures: [],
    standings: [],
    predictions: [],
    selectedWeek: null,
    selectedMatch: null,
    loading: false,
    error: null,
  }),
  actions: {
    clearError() {
      this.error = null
    },

    setSelectedWeek(week: number | null) {
      this.selectedWeek = week
    },

    async runAction<T>(action: () => Promise<T>) {
      this.loading = true
      this.error = null

      try {
        return await action()
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Unexpected error'
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchLeague() {
      return this.runAction(async () => {
        const response = await leagueApi.getLeague()
        this.league = response.league
        return response.league
      })
    },

    async fetchTeams() {
      return this.runAction(async () => {
        const response = await leagueApi.getTeams()
        this.teams = response.teams
        return response.teams
      })
    },

    async fetchFixtures() {
      return this.runAction(async () => {
        const response = await leagueApi.getFixtures()
        this.fixtures = response.fixtures
        return response.fixtures
      })
    },

    async fetchFixturesByWeek(week: number) {
      return this.runAction(async () => {
        const response = await leagueApi.getFixturesByWeek(week)
        this.selectedWeek = week
        return response.fixtures
      })
    },

    async fetchStandings() {
      return this.runAction(async () => {
        const response = await leagueApi.getStandings()
        this.standings = response.standings
        return response.standings
      })
    },

    async fetchPredictions() {
      return this.runAction(async () => {
        const response = await leagueApi.getPredictions()
        this.predictions = response.predictions
        return response
      })
    },

    async fetchMatch(id: number) {
      return this.runAction(async () => {
        const response = await leagueApi.getMatch(id)
        this.selectedMatch = response.match
        this.fixtures = mergeMatchIntoFixtures(this.fixtures, response.match)
        return response.match
      })
    },

    async initializeLeague() {
      return this.runAction(async () => {
        const response = await leagueApi.initLeague()
        this.league = response.league
        this.teams = response.teams
        this.fixtures = response.fixtures
        this.standings = response.standings
        this.predictions = []
        this.selectedMatch = null
        this.selectedWeek = null
        return response
      })
    },

    async resetLeague() {
      return this.runAction(async () => {
        const response = await leagueApi.resetLeague()
        this.league = response.league
        this.teams = response.teams
        this.fixtures = response.fixtures
        this.standings = response.standings
        this.predictions = []
        this.selectedMatch = null
        this.selectedWeek = null
        return response
      })
    },

    async playNextWeek() {
      return this.runAction(async () => {
        const response = await leagueApi.playNextWeek()
        this.league = response.league
        this.standings = response.standings
        this.predictions = response.predictions
        for (const match of response.matches) {
          this.fixtures = mergeMatchIntoFixtures(this.fixtures, match)
        }
        return response
      })
    },

    async playWeek(week: number) {
      return this.runAction(async () => {
        const response = await leagueApi.playWeek(week)
        this.league = response.league
        this.standings = response.standings
        this.predictions = response.predictions
        this.selectedWeek = week
        for (const match of response.matches) {
          this.fixtures = mergeMatchIntoFixtures(this.fixtures, match)
        }
        return response
      })
    },

    async playAll() {
      return this.runAction(async () => {
        const response = await leagueApi.playAll()
        this.league = response.league
        this.standings = response.standings
        this.predictions = response.predictions
        await this.fetchFixtures()
        return response
      })
    },

    async updateMatchResult(id: number, homeGoals: number, awayGoals: number) {
      return this.runAction(async () => {
        const response = await leagueApi.updateMatchResult(id, homeGoals, awayGoals)
        this.selectedMatch = response.match
        this.league = response.league
        this.standings = response.standings
        this.predictions = response.predictions
        this.fixtures = mergeMatchIntoFixtures(this.fixtures, response.match)
        return response
      })
    },

    async refreshAll() {
      return this.runAction(async () => {
        const [leagueResponse, teamsResponse, fixturesResponse, standingsResponse, predictionsResponse] =
          await Promise.all([
            leagueApi.getLeague(),
            leagueApi.getTeams(),
            leagueApi.getFixtures(),
            leagueApi.getStandings(),
            leagueApi.getPredictions(),
          ])

        this.league = leagueResponse.league
        this.teams = teamsResponse.teams
        this.fixtures = fixturesResponse.fixtures
        this.standings = standingsResponse.standings
        this.predictions = predictionsResponse.predictions

        return {
          league: leagueResponse.league,
          teams: teamsResponse.teams,
          fixtures: fixturesResponse.fixtures,
          standings: standingsResponse.standings,
          predictions: predictionsResponse.predictions,
        }
      })
    },
  },
})
