<script setup lang="ts">
import { RouterView } from 'vue-router'
import TopBar from './components/TopBar.vue'
import SideNav from './components/SideNav.vue'
</script>

<template>
  <v-app>
    <TopBar />
    <SideNav />
    <v-main>
      <RouterView />
    </v-main>
  </v-app>
</template>

<style scoped></style>
