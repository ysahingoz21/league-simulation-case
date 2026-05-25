<template>
  <main class="page-shell">
    <AppHeader />

    <section class="dashboard-meta">
      <p class="api-note">API base: {{ apiBaseUrl }}</p>
    </section>

    <div class="dashboard-grid">
      <div class="dashboard-main">
        <LeagueControls />
        <LeagueStatusCard />
        <StandingsTable />
      </div>

      <div class="dashboard-side">
        <FixturesPanel />
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

import { API_BASE_URL } from './api/client'
import AppHeader from './components/AppHeader.vue'
import FixturesPanel from './components/FixturesPanel.vue'
import LeagueControls from './components/LeagueControls.vue'
import LeagueStatusCard from './components/LeagueStatusCard.vue'
import StandingsTable from './components/StandingsTable.vue'
import { useLeagueStore } from './stores/leagueStore'

const apiBaseUrl = API_BASE_URL
const leagueStore = useLeagueStore()

onMounted(() => {
  leagueStore.refreshAll().catch(() => {
    // The store exposes the backend error so the user can initialize the league.
  })
})
</script>
