<template>
  <section class="card status-card">
    <div class="section-heading">
      <div>
        <p class="section-kicker">League Status</p>
        <h2>Current Snapshot</h2>
      </div>
      <span class="status-pill" :class="statusClass">{{ statusLabel }}</span>
    </div>

    <div class="stats-grid">
      <article class="stat-item">
        <span class="stat-label">Current Week</span>
        <strong class="stat-value">{{ league?.current_week ?? '-' }}</strong>
      </article>
      <article class="stat-item">
        <span class="stat-label">Total Weeks</span>
        <strong class="stat-value">{{ league?.total_weeks ?? 6 }}</strong>
      </article>
      <article class="stat-item">
        <span class="stat-label">Teams</span>
        <strong class="stat-value">{{ teams.length }}</strong>
      </article>
      <article class="stat-item">
        <span class="stat-label">Fixtures</span>
        <strong class="stat-value">{{ fixtures.length }}</strong>
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
  if (!league.value) {
    return 'Not initialized'
  }

  return league.value.is_completed ? 'Completed' : 'In progress'
})

const statusClass = computed(() => {
  if (!league.value) {
    return 'status-pill-neutral'
  }

  return league.value.is_completed ? 'status-pill-complete' : 'status-pill-live'
})
</script>
