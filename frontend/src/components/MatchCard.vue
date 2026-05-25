<template>
  <article class="match-card">
    <div class="match-card__top">
      <span class="match-week">Week {{ match.week }}</span>
      <span class="match-badge" :class="badgeClass">{{ badgeLabel }}</span>
    </div>

    <div class="match-card__teams">
      <div class="match-team-row">
        <span class="match-team-name">{{ match.home_team_name }}</span>
        <strong class="match-score">{{ scoreLabel.home }}</strong>
      </div>
      <div class="match-team-row">
        <span class="match-team-name">{{ match.away_team_name }}</span>
        <strong class="match-score">{{ scoreLabel.away }}</strong>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import type { Match } from '../types/league'

const props = defineProps<{
  match: Match
}>()

const badgeLabel = computed(() =>
  props.match.status === 'played' ? 'Played' : 'Scheduled',
)

const badgeClass = computed(() =>
  props.match.status === 'played' ? 'match-badge-played' : 'match-badge-scheduled',
)

const scoreLabel = computed(() => {
  if (props.match.status !== 'played') {
    return { home: '-', away: '-' }
  }

  return {
    home: String(props.match.home_goals ?? 0),
    away: String(props.match.away_goals ?? 0),
  }
})
</script>
