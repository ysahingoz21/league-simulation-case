<template>
  <section class="card controls-card">
    <div class="section-heading">
      <div>
        <p class="section-kicker">Controls</p>
        <h2>League Actions</h2>
      </div>
      <button
        v-if="error"
        class="button button-ghost"
        type="button"
        @click="leagueStore.clearError()"
      >
        Dismiss Error
      </button>
    </div>

    <div v-if="error" class="error-banner">
      {{ error }}
    </div>

    <div class="button-grid">
      <button
        class="button"
        type="button"
        :disabled="loading"
        @click="runAction(leagueStore.initializeLeague)"
      >
        Initialize League
      </button>
      <button
        class="button"
        type="button"
        :disabled="loading"
        @click="runAction(leagueStore.resetLeague)"
      >
        Reset League
      </button>
      <button
        class="button button-accent"
        type="button"
        :disabled="loading || isCompleted"
        @click="runAction(leagueStore.playNextWeek)"
      >
        Play Next Week
      </button>
      <button
        class="button button-accent"
        type="button"
        :disabled="loading || isCompleted"
        @click="runAction(leagueStore.playAll)"
      >
        Play All Remaining
      </button>
      <button
        class="button button-ghost"
        type="button"
        :disabled="loading"
        @click="runAction(leagueStore.refreshAll)"
      >
        Refresh
      </button>
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
    // Store owns error state.
  }
}
</script>
