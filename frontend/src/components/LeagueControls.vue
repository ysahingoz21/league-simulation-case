<template>
  <section class="card">
    <div class="card-header">
      <h2 class="card-title">League Controls</h2>
      <button
        v-if="error"
        class="btn btn-ghost btn-sm"
        type="button"
        @click="leagueStore.clearError()"
      >
        Dismiss
      </button>
    </div>

    <div v-if="error" class="error-banner">{{ error }}</div>

    <div class="controls-actions">
      <button
        class="btn btn-primary btn-primary--dark-text btn-block controls-btn-lg"
        type="button"
        :disabled="loading || isCompleted"
        @click="runAction(leagueStore.playNextWeek)"
      >
        ▶ Play Next Week
      </button>

      <button
        class="btn btn-secondary btn-block controls-btn-lg"
        type="button"
        :disabled="loading"
        @click="runAction(leagueStore.initializeLeague)"
      >
        Initialize League
      </button>

      <button
        class="btn btn-secondary btn-block controls-btn-lg"
        type="button"
        :disabled="loading || isCompleted"
        @click="runAction(leagueStore.playAll)"
      >
        Play All Remaining
      </button>

      <div class="controls-actions-row">
        <button
          class="btn btn-ghost btn-sm"
          type="button"
          :disabled="loading"
          @click="runAction(leagueStore.resetLeague)"
        >
          Reset League
        </button>
        <button
          class="btn btn-ghost btn-sm"
          type="button"
          :disabled="loading"
          @click="runAction(leagueStore.refreshAll)"
        >
          Refresh
        </button>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'

import { useLeagueStore } from '../stores/leagueStore'

const leagueStore = useLeagueStore()
const { league, loading, error } = storeToRefs(leagueStore)

const isCompleted = computed(() => league.value?.is_completed ?? false)

async function runAction<T>(action: () => Promise<T>) {
  try {
    await action()
  } catch {
  }
}
</script>
