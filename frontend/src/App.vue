<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="sidebar-brand">
        <img src="/team-logos/app-logo.png"
          alt="WCS"
          class="sidebar-brand-logo" />
        <div>
          <div class="sidebar-brand-name">World Cup Group Simulator</div>
          <div class="sidebar-brand-season">2026</div>
        </div>
      </div>

      <div class="sidebar-section-title">Sim Manager</div>

      <nav class="sidebar-nav">
        <a href="#overview"
          class="sidebar-nav-item sidebar-nav-item--active">
          <span class="sidebar-nav-dot"></span>
          Overview
        </a>
        <a href="#standings"
          class="sidebar-nav-item">League Table</a>
        <a href="#fixtures"
          class="sidebar-nav-item">Match Center</a>
        <a href="#predictions"
          class="sidebar-nav-item">Predictions</a>
      </nav>

      <div class="sidebar-footer">
        <a href="https://github.com/ysahingoz21"
          target="_blank"
          rel="noopener noreferrer"
          class="sidebar-credit">© Yusuf Enes Şahingöz - 2026</a>
      </div>
    </aside>

    <div class="main-wrapper">
      <header class="top-bar">
        <span class="top-bar-title">World Cup Group Simulator</span>
        <div class="top-bar-actions">
          <span class="badge badge-success">
            <span class="badge-dot"></span>
            Backend Connected
          </span>
          <button class="btn btn-secondary btn-sm"
            type="button"
            :disabled="loading"
            @click="runAction(leagueStore.resetLeague)">
            Reset League
          </button>
          <button class="btn btn-primary btn-primary--dark-text btn-sm"
            type="button"
            :disabled="loading || isCompleted"
            @click="runAction(leagueStore.playNextWeek)">
            Play Next Week
          </button>
        </div>
      </header>

      <main class="page-content">
        <div id="overview"
          class="dashboard-grid-top">
          <div class="dashboard-left-col">
            <LeagueControls />
            <LeagueStatusCard />
          </div>
          <div id="standings">
            <StandingsTable />
          </div>
        </div>

        <div class="dashboard-grid-bottom">
          <div id="fixtures">
            <FixturesPanel @edit-match="openMatchEditor" />
          </div>
          <div id="predictions">
            <PredictionsPanel />
          </div>
        </div>
      </main>
    </div>

    <EditMatchModal v-if="selectedMatch"
      :match="selectedMatch"
      @close="closeMatchEditor" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'

import EditMatchModal from './components/EditMatchModal.vue'
import FixturesPanel from './components/FixturesPanel.vue'
import LeagueControls from './components/LeagueControls.vue'
import LeagueStatusCard from './components/LeagueStatusCard.vue'
import PredictionsPanel from './components/PredictionsPanel.vue'
import StandingsTable from './components/StandingsTable.vue'
import { useLeagueStore } from './stores/leagueStore'
import type { Match } from './types/league'

const leagueStore = useLeagueStore()
const { loading, league } = storeToRefs(leagueStore)
const selectedMatch = ref<Match | null>(null)

const isCompleted = computed(() => league.value?.is_completed ?? false)

onMounted(() => {
  leagueStore.refreshAll().catch(() => { })
})

async function runAction<T>(action: () => Promise<T>) {
  try {
    await action()
  } catch {
  }
}

function openMatchEditor(match: Match) {
  selectedMatch.value = match
  leagueStore.clearError()
}

function closeMatchEditor() {
  selectedMatch.value = null
}
</script>
