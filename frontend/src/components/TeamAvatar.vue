<template>
  <div
    class="team-avatar"
    :style="showInitials ? { background: teamColor(teamName) } : { background: '#e2e8f0' }"
  >
    <img
      v-if="logoSrc && !imgFailed"
      :src="logoSrc"
      :alt="teamName"
      @error="imgFailed = true"
    />
    <template v-else>{{ getInitials(teamName) }}</template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { getInitials, teamColor } from '../utils/team'
import { getTeamLogo } from '../utils/teamLogos'

const props = defineProps<{
  teamName: string
}>()

const imgFailed = ref(false)

watch(() => props.teamName, () => { imgFailed.value = false })

const logoSrc = computed(() => getTeamLogo(props.teamName))
const showInitials = computed(() => !logoSrc.value || imgFailed.value)
</script>
