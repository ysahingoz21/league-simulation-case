<template>
  <section class="card">
    <div class="card-header">
      <div>
        <h2 class="card-title">Season Predictions</h2>
        <p class="card-subtitle">Post-Week 4 Predictive Model</p>
      </div>
    </div>

    <div v-if="predictions.length === 0"
      class="empty-state">
      Predictions are available after week 4.
    </div>

    <div v-else
      class="prediction-list">
      <p v-if="isCompleted"
        class="prediction-note">
        Final probabilities are based on the completed league table.
      </p>

      <article v-for="prediction in predictions"
        :key="prediction.team_id"
        class="prediction-item">
        <div class="prediction-item__top">
          <div class="prediction-team-display">
            <TeamAvatar :team-name="prediction.team_name" />
            <span class="prediction-team-name">{{ prediction.team_name }}</span>
          </div>
          <strong class="prediction-prob">{{ formatPercent(prediction.championship_probability) }}</strong>
        </div>

        <div class="prediction-bar-track"
          aria-hidden="true">
          <div class="prediction-bar-fill"
            :style="{ width: `${prediction.championship_probability}%` }" />
        </div>

        <div class="prediction-meta">
          <span>Expected Points: {{ formatNumber(prediction.expected_points) }}</span>
          <span>Estimated Rank: #{{ formatNumber(prediction.projected_rank) }}</span>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'

import TeamAvatar from './TeamAvatar.vue'
import { useLeagueStore } from '../stores/leagueStore'

const leagueStore = useLeagueStore()
const { league, predictions } = storeToRefs(leagueStore)

const isCompleted = computed(() => league.value?.is_completed ?? false)

function formatPercent(value: number) {
  return `${value.toFixed(1)}%`
}

function formatNumber(value: number) {
  return value.toFixed(1)
}
</script>
