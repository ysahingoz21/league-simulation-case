import { defineStore } from 'pinia'

type LeagueState = {
  loading: boolean
  error: string | null
}

export const useLeagueStore = defineStore('league', {
  state: (): LeagueState => ({
    loading: false,
    error: null,
  }),
  actions: {
    clearError() {
      this.error = null
    },
  },
})
