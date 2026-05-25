<template>
  <section class="card">
    <div class="card-header">
      <h2 class="card-title">Match Center</h2>
      <div class="week-tabs">
        <button
          class="week-tab"
          :class="{ 'week-tab--active': selectedWeek === null }"
          type="button"
          @click="leagueStore.setSelectedWeek(null)"
        >
          All
        </button>
        <button
          v-for="week in weeks"
          :key="week"
          class="week-tab"
          :class="{ 'week-tab--active': selectedWeek === week }"
          type="button"
          @click="leagueStore.setSelectedWeek(week)"
        >
          Week {{ week }}
        </button>
      </div>
    </div>

    <div v-if="groupedFixtures.length === 0" class="empty-state">
      Initialize the league to view fixtures.
    </div>

    <div v-else class="fixture-groups">
      <section
        v-for="group in groupedFixtures"
        :key="group.week"
        class="fixture-group"
      >
        <p class="fixture-week-label">Week {{ group.week }}</p>
        <div class="fixture-list">
          <MatchCard
            v-for="match in group.matches"
            :key="match.id"
            :match="match"
            @edit-match="$emit('edit-match', $event)"
          />
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'

import MatchCard from './MatchCard.vue'
import { useLeagueStore } from '../stores/leagueStore'

defineEmits<{
  (event: 'edit-match', match: import('../types/league').Match): void
}>()

const leagueStore = useLeagueStore()
const { fixtures, selectedWeek } = storeToRefs(leagueStore)

const weeks = computed(() => {
  const weekSet = new Set(fixtures.value.map(f => f.week))
  return Array.from(weekSet).sort((a, b) => a - b)
})

const filteredFixtures = computed(() => {
  if (selectedWeek.value === null) return fixtures.value
  return fixtures.value.filter(f => f.week === selectedWeek.value)
})

const groupedFixtures = computed(() => {
  const groups = new Map<number, typeof fixtures.value>()

  for (const fixture of filteredFixtures.value) {
    const current = groups.get(fixture.week) ?? []
    groups.set(fixture.week, [...current, fixture])
  }

  return Array.from(groups.entries())
    .sort(([a], [b]) => a - b)
    .map(([week, matches]) => ({
      week,
      matches: matches.sort((a, b) => a.id - b.id),
    }))
})
</script>
