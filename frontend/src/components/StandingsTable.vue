<template>
  <section class="card table-card">
    <div class="section-heading">
      <div>
        <p class="section-kicker">Standings</p>
        <h2>League Table</h2>
      </div>
    </div>

    <div v-if="standings.length === 0" class="empty-state">
      Initialize the league to view standings.
    </div>

    <div v-else class="table-wrap">
      <table class="standings-table">
        <thead>
          <tr>
            <th>Rank</th>
            <th>Team</th>
            <th>P</th>
            <th>W</th>
            <th>D</th>
            <th>L</th>
            <th>GF</th>
            <th>GA</th>
            <th>GD</th>
            <th>Pts</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="standing in standings" :key="standing.team_id">
            <td>{{ standing.rank }}</td>
            <td class="team-cell">{{ standing.team_name }}</td>
            <td>{{ standing.played }}</td>
            <td>{{ standing.wins }}</td>
            <td>{{ standing.draws }}</td>
            <td>{{ standing.losses }}</td>
            <td>{{ standing.goals_for }}</td>
            <td>{{ standing.goals_against }}</td>
            <td>{{ standing.goal_difference }}</td>
            <td class="points-cell">{{ standing.points }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'

import { useLeagueStore } from '../stores/leagueStore'

const leagueStore = useLeagueStore()
const { standings } = storeToRefs(leagueStore)
</script>
