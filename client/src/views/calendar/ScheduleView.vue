<template>
  <ScheduleXCalendar :calendar-app="calendarApp" />
  <v-app-bar density="compact">
    <v-app-bar-title></v-app-bar-title>
    <v-menu>
      <template v-slot:activator="{ props }">
        <v-btn density="compact" icon="mdi-cog" v-bind="props" />
      </template>
      <v-list>
        <v-list-item
          v-for="calendar in myCalendars" :key="calendar.name"
          prepend-icon="mdi-calendar"
          :title="calendar.title"
        />
        <v-list-item
          prepend-icon="mdi-calendar-plus"
          title="Create Calendar"
        />
      </v-list>
    </v-menu>
  </v-app-bar>
</template>

<script setup lang="ts">
import { ScheduleXCalendar } from '@schedule-x/vue'
import {
  createCalendar,
  createViewDay,
  createViewMonthAgenda,
  createViewMonthGrid,
  createViewWeek,
  createViewList,
} from '@schedule-x/calendar'
import '@schedule-x/theme-default/dist/index.css'
import { useCalendarsStore } from '@/stores/calendar'
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'

const calendarsStore = useCalendarsStore()

const { myCalendars } = storeToRefs(calendarsStore)

onMounted(async () => {
  await calendarsStore.loadMyCalendars('')
})

// Do not use a ref here, as the calendar instance is not reactive, and doing so might cause issues
// For updating events, use the events service plugin
const calendarApp = createCalendar({
  views: [
    createViewDay(),
    createViewWeek(),
    createViewMonthGrid(),
    createViewMonthAgenda(),
    createViewList(),
  ],
})
</script>

<style>

</style>
