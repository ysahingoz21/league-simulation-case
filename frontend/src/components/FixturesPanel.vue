<template>
  <section class="card fixtures-card">
    <div class="section-heading">
      <div>
        <p class="section-kicker">Fixtures</p>
        <h2>Match Schedule</h2>
      </div>

      <label class="week-filter">
        <span>Filter</span>
        <select :value="selectedWeekModel" @change="onWeekChange">
          <option value="all">All weeks</option>
          <option v-for="week in weeks" :key="week" :value="week">
            Week {{ week }}
          </option>
        </select>
      </label>
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
        <header class="fixture-group__header">
          <h3>Week {{ group.week }}</h3>
        </header>
        <div class="fixture-list">
          <MatchCard
            v-for="match in group.matches"
            :key="match.id"
            :match="match"
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

const weeks = [1, 2, 3, 4, 5, 6]

const leagueStore = useLeagueStore()
const { fixtures, selectedWeek } = storeToRefs(leagueStore)

const selectedWeekModel = computed(() =>
  selectedWeek.value === null ? 'all' : String(selectedWeek.value),
)

const filteredFixtures = computed(() => {
  if (selectedWeek.value === null) {
    return fixtures.value
  }

  return fixtures.value.filter((fixture) => fixture.week === selectedWeek.value)
})

const groupedFixtures = computed(() => {
  const groups = new Map<number, typeof fixtures.value>()

  for (const fixture of filteredFixtures.value) {
    const current = groups.get(fixture.week) ?? []
    groups.set(fixture.week, [...current, fixture])
  }

  return Array.from(groups.entries())
    .sort(([left], [right]) => left - right)
    .map(([week, matches]) => ({
      week,
      matches: matches.sort((left, right) => left.id - right.id),
    }))
})

function onWeekChange(event: Event) {
  const value = (event.target as HTMLSelectElement).value
  leagueStore.setSelectedWeek(value === 'all' ? null : Number(value))
}
</script>
