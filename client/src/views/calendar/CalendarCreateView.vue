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
import { useRouter, useRoute } from 'vue-router'
import { onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import CalendarForm from '@/views/calendar/forms/CalendarForm.vue'
import type { apitypes_VisibilityLevel } from '@/genapi/api/calendars/calendar/v1alpha1'
import { useAlertStore } from '@/stores/alerts'

const router = useRouter()
const route = useRoute()
const calendarsStore = useCalendarsStore()
const authStore = useAuthStore()
const alertsStore = useAlertStore()

function navigateBack() {
  if (route.params.circleId) {
    router.push({ name: 'circle', params: { circleId: route.params.circleId } })
  } else {
    router.push({ name: 'calendars' })
  }
}

async function saveCalendar() {
  if (!authStore.user?.name && !route.params.circleId) {
    throw new Error('User not authenticated')
  }
  try {
    const parent = circleName.value ? circleName.value : authStore.user!.name!
    const calendar = await calendarsStore.createCalendar(parent)
    router.push('/'+calendar.name)
  } catch (err) {
    alertsStore.addAlert(err instanceof Error ? err.message : String(err),'error')
  }
}

const circleName = computed(() => {
  return route.path.replace('/calendars/create', '').slice(1)
})

onMounted(() => {
  calendarsStore.calendar = {
    name: '',
    title: '',
    visibility: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
    description: '',
    calendarAccess: undefined,
    favorited: false,
  }
})
</script>

<style scoped>
</style>
