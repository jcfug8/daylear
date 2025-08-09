<template>
  <calendar-form
    v-if="calendarsStore.calendar"
    v-model="calendarsStore.calendar"
    :is-editing="false"
    @save="saveCalendar"
    @close="navigateBack"
  />
</template>

<script setup lang="ts">
import { useCalendarsStore } from '@/stores/calendar'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import CalendarForm from '@/views/calendar/forms/CalendarForm.vue'
import type { apitypes_VisibilityLevel } from '@/genapi/api/calendars/calendar/v1alpha1'
import { useAlertStore } from '@/stores/alerts'

const router = useRouter()
const calendarsStore = useCalendarsStore()
const authStore = useAuthStore()
const alertsStore = useAlertStore()

function navigateBack() {
  router.push({ name: 'calendars' })
}

async function saveCalendar() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  try {
    const calendar = await calendarsStore.createCalendar(authStore.user.name)
    router.push('/'+calendar.name)
  } catch (err) {
    alertsStore.addAlert(err instanceof Error ? err.message : String(err),'error')
  }
}

onMounted(() => {
  calendarsStore.calendar = {
    name: '',
    title: '',
    visibility: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
    description: '',
    calendarAccess: undefined,
  }
})
</script>

<style scoped>
</style>
