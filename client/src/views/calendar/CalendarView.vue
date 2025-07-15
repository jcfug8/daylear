<template>
  <v-container>
    <v-calendar ref="calendar" :v-model="[new Date()]" view-mode="week"></v-calendar>
  </v-container>
</template>

<script setup lang="ts"></script>

<style>
.v-calendar__container {
  overflow: auto;
}
</style>
