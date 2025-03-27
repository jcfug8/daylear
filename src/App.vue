<script setup lang="ts">
import { RouterView } from 'vue-router'
import TopBar from './components/MainNav.vue'
</script>

<template>
  <v-app>
    <TopBar />
    <v-main>
      <RouterView />
    </v-main>
  </v-app>
</template>

<style scoped></style>
