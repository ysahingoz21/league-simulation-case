<template>
  <section class="card">
    <div class="card-header">
      <h2 class="card-title">League Standings</h2>
      <span class="badge badge-live">
        <span class="badge-dot"></span>
        Live Updates
      </span>
    </div>

    <div v-if="standings.length === 0" class="empty-state">
      Initialize the league to view standings.
    </div>

    <div v-else class="table-wrap">
      <table class="standings-table">
        <thead>
          <tr>
            <th>RANK</th>
            <th class="team-col">Team</th>
            <th>P</th>
            <th>W</th>
            <th>D</th>
            <th>L</th>
            <th>GF</th>
            <th>GA</th>
            <th>GD</th>
            <th>PTS</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="standing in standings" :key="standing.team_id">
            <td :class="['rank-cell', { 'rank-cell--first': standing.rank === 1 }]">{{ standing.rank }}</td>
            <td class="team-cell">
              <div class="team-display">
                <TeamAvatar :team-name="standing.team_name" />
                <span class="team-name">{{ standing.team_name }}</span>
              </div>
            </td>
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

import TeamAvatar from './TeamAvatar.vue'
import { useLeagueStore } from '../stores/leagueStore'

const leagueStore = useLeagueStore()
const { standings } = storeToRefs(leagueStore)
</script>
