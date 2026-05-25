<template>
  <article class="match-card">
    <div v-if="match.status !== 'played'" class="match-upcoming-badge">Upcoming</div>

    <div class="match-card__scoreboard">
      <div class="match-team">
        <TeamAvatar :team-name="match.home_team_name" />
        <span class="match-team-name">{{ match.home_team_name }}</span>
        <span class="match-team-role">Home</span>
      </div>

      <div class="match-score-center">
        <div v-if="match.status === 'played'" class="match-score-row">
          <span class="match-score-num">{{ match.home_goals ?? 0 }}</span>
          <span class="match-ft-label">FT</span>
          <span class="match-score-num">{{ match.away_goals ?? 0 }}</span>
        </div>
        <div v-else class="match-scheduled-state">
          <svg
            class="match-scheduled-icon"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.8"
            stroke-linecap="round"
            stroke-linejoin="round"
            aria-hidden="true"
          >
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
            <line x1="16" y1="2" x2="16" y2="6" />
            <line x1="8" y1="2" x2="8" y2="6" />
            <line x1="3" y1="10" x2="21" y2="10" />
          </svg>
          <span class="match-scheduled-text">Scheduled</span>
        </div>
      </div>

      <div class="match-team">
        <TeamAvatar :team-name="match.away_team_name" />
        <span class="match-team-name">{{ match.away_team_name }}</span>
        <span class="match-team-role">Away</span>
      </div>
    </div>

    <div v-if="match.status === 'played'" class="match-card__footer">
      <span class="match-metadata">Week {{ match.week }} · Played</span>
      <button
        class="btn btn-ghost btn-sm"
        type="button"
        @click="$emit('edit-match', match)"
      >
        Edit Result
      </button>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { Match } from '../types/league'
import TeamAvatar from './TeamAvatar.vue'

defineProps<{
  match: Match
}>()

defineEmits<{
  (event: 'edit-match', match: Match): void
}>()
</script>
