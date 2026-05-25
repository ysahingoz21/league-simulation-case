<template>
  <section class="card">
    <div class="card-header">
      <h2 class="card-title">Current Season Meta</h2>
      <span class="badge" :class="statusBadgeClass">{{ statusLabel }}</span>
    </div>

    <div class="stats-grid">
      <article class="stat-tile">
        <span class="stat-label">Week</span>
        <strong class="stat-value">
          {{ currentWeek }}<span class="stat-fraction">/{{ totalWeeks }}</span>
        </strong>
      </article>

      <article class="stat-tile">
        <span class="stat-label">Teams</span>
        <strong class="stat-value">{{ teams.length }}</strong>
      </article>

      <article class="stat-tile">
        <span class="stat-label">Played</span>
        <strong class="stat-value">
          {{ playedCount }}<span class="stat-fraction">/{{ fixtures.length }}</span>
        </strong>
      </article>

      <article class="stat-tile">
        <span class="stat-label">Status</span>
        <strong class="stat-value-sm">{{ statusLabel }}</strong>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'

import { useLeagueStore } from '../stores/leagueStore'

const leagueStore = useLeagueStore()
const { league, teams, fixtures } = storeToRefs(leagueStore)

const statusLabel = computed(() => {
  if (!league.value) return 'Not Initialized'
  return league.value.is_completed ? 'Completed' : 'In Progress'
})

const statusBadgeClass = computed(() => {
  if (!league.value) return 'badge-neutral'
  return league.value.is_completed ? 'badge-complete' : 'badge-live'
})

const currentWeek = computed(() => league.value?.current_week ?? '-')
const totalWeeks = computed(() => league.value?.total_weeks ?? 6)

const playedCount = computed(() =>
  fixtures.value.filter(f => f.status === 'played').length,
)
</script>
