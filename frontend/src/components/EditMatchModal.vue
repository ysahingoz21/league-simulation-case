<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <section class="modal-card" role="dialog" aria-modal="true">
      <div class="modal-header">
        <h2 class="modal-title">Edit Match Result</h2>
        <button class="modal-close-btn" type="button" @click="$emit('close')">
          ×
        </button>
      </div>

      <div class="modal-body">
        <div class="modal-teams-row">
          <div class="modal-team">
            <TeamAvatar :team-name="match.home_team_name" class="team-avatar-lg" />
            <span class="modal-team-name">{{ match.home_team_name }}</span>
          </div>

          <span class="modal-vs-label">VS</span>

          <div class="modal-team">
            <TeamAvatar :team-name="match.away_team_name" class="team-avatar-lg" />
            <span class="modal-team-name">{{ match.away_team_name }}</span>
          </div>
        </div>

        <div class="modal-score-inputs">
          <div class="field">
            <span class="field-label">{{ match.home_team_name }}</span>
            <input
              v-model.number="homeGoals"
              class="field-input"
              type="number"
              min="0"
              max="20"
              inputmode="numeric"
            />
          </div>

          <div class="field">
            <span class="field-label">{{ match.away_team_name }}</span>
            <input
              v-model.number="awayGoals"
              class="field-input"
              type="number"
              min="0"
              max="20"
              inputmode="numeric"
            />
          </div>
        </div>

        <div class="info-callout">
          <svg class="info-callout__icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          <span>Updating the score will automatically recalculate standings and championship predictions.</span>
        </div>

        <p v-if="validationError" class="error-inline">{{ validationError }}</p>
        <p v-else-if="storeError" class="error-inline">{{ storeError }}</p>
      </div>

      <div class="modal-footer">
        <button
          class="btn btn-ghost"
          type="button"
          :disabled="loading"
          @click="$emit('close')"
        >
          Cancel
        </button>
        <button
          class="btn btn-primary"
          type="button"
          :disabled="loading"
          @click="save"
        >
          Save Result
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'

import type { Match } from '../types/league'
import TeamAvatar from './TeamAvatar.vue'
import { useLeagueStore } from '../stores/leagueStore'

const props = defineProps<{
  match: Match
}>()

const emit = defineEmits<{
  (event: 'close'): void
}>()

const leagueStore = useLeagueStore()
const { loading, error: storeError } = storeToRefs(leagueStore)

const homeGoals = ref(props.match.home_goals ?? 0)
const awayGoals = ref(props.match.away_goals ?? 0)
const validationError = ref<string | null>(null)

watch(
  () => props.match,
  (match) => {
    homeGoals.value = match.home_goals ?? 0
    awayGoals.value = match.away_goals ?? 0
    validationError.value = null
  },
  { deep: true },
)

const isValid = computed(() =>
  Number.isInteger(homeGoals.value) &&
  Number.isInteger(awayGoals.value) &&
  homeGoals.value >= 0 &&
  awayGoals.value >= 0 &&
  homeGoals.value <= 20 &&
  awayGoals.value <= 20,
)

async function save() {
  validationError.value = null

  if (!isValid.value) {
    validationError.value = 'Goals must be integers between 0 and 20.'
    return
  }

  try {
    await leagueStore.updateMatchResult(props.match.id, homeGoals.value, awayGoals.value)
    emit('close')
  } catch {
  }
}
</script>
