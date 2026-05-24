<template>
  <main class="page-shell">
    <section class="hero">
      <p class="eyebrow">Frontend Scaffold</p>
      <h1>League Simulation Dashboard</h1>
      <p class="subtitle">Vue 3 frontend for the Go league simulation API</p>
      <p class="api-note">API base: {{ apiBaseUrl }}</p>
    </section>
  </main>
</template>

<script setup lang="ts">
import { API_BASE_URL } from './api/client'

const apiBaseUrl = API_BASE_URL
</script>
