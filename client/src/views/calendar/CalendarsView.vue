<template>
  <ListTabsPage
    ref="tabsPage"
    :tabs="tabs"
  >
    <template #my="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
      </div>
      <CalendarGrid :calendars="items" :loading="loading" />
    </template>
    <template #pending="{ items, loading }">
      <CalendarGrid :calendars="items" @accept="acceptCalendarAccess" @decline="onDeclinetCalendar" :acceptingCalendarId="acceptingCalendarId" :loading="loading" />
      <div v-if="!loading && items.length === 0">No pending shared calendars found.</div>
    </template>
    <template #explore="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
      </div>
      <CalendarGrid :calendars="items" :loading="loading" />
    </template>
    <template #fab>
      <v-btn
        color="primary"
        density="compact"
        style="position: fixed; bottom: 16px; right: 16px"
        :to="{ name: 'calendarCreate' }"
      >
        <v-icon>mdi-plus</v-icon>
        <span>Create Calendar</span>
      </v-btn>
    </template>
  </ListTabsPage>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCalendarsStore } from '@/stores/calendar'
import { useAuthStore } from '@/stores/auth'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import CalendarGrid from '@/components/CalendarGrid.vue'
import { useAlertStore } from '@/stores/alerts'
import type { Calendar } from '@/genapi/api/calendars/calendar/v1alpha1'

const calendarsStore = useCalendarsStore()
const authStore = useAuthStore()

const acceptingCalendarId = ref<string | null>(null)
const tabsPage = ref()
const alertsStore = useAlertStore()

const tabs = [
  { 
    label: 'My Calendars', 
    value: 'my', 
    icon: 'mdi-calendar-account',
    loader: async () => {
      await calendarsStore.loadMyCalendars(authStore.user?.name || '')
      return [...calendarsStore.myCalendars]
    }
  },
  { 
    label: 'Pending',  
    value: 'pending', 
    icon: 'mdi-calendar-clock',
    loader: async () => {
      await calendarsStore.loadPendingCalendars(authStore.user?.name || '')
      return [...calendarsStore.sharedPendingCalendars]
    }
  },
  { 
    label: 'Explore',  
    value: 'explore', 
    icon: 'mdi-compass',
    loader: async () => {
      await calendarsStore.loadPublicCalendars(authStore.user?.name || '')
      return [...calendarsStore.publicCalendars]
    }
  }
]

async function acceptCalendarAccess(calendar: Calendar) {
  if (!calendar.name) return
  
  acceptingCalendarId.value = calendar.name
  try {
    await calendarsStore.acceptCalendar(calendar.name)
    tabsPage.value?.reloadTab('pending')
  } catch (error) {
    alertsStore.addAlert(`Failed to accept calendar access: ${error}`)
  } finally {
    acceptingCalendarId.value = null
  }
}

async function onDeclinetCalendar(calendar: Calendar) {
  if (!calendar.name) return
  
  try {
    await calendarsStore.deleteCalendarAccess(calendar.name)
    // Refresh the calendars list
    tabsPage.value?.reloadTab('pending')
  } catch (error) {
    alertsStore.addAlert(`Failed to decline calendar access: ${error}`)
  }
}

</script>

<style scoped>
</style>