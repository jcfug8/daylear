<template>
  <calendar-form
    v-if="calendarsStore.calendar"
    v-model="editedCalendar"
    :is-editing="true"
    @save="saveSettings"
    @close="navigateBack"
  />
</template>

<script setup lang="ts">
import { useCalendarsStore } from '@/stores/calendar'
import { storeToRefs } from 'pinia'
import { onMounted, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { Calendar, apitypes_VisibilityLevel } from '@/genapi/api/calendars/calendar/v1alpha1'
import CalendarForm from '@/views/calendar/forms/CalendarForm.vue'
import { computed } from 'vue'
import { useAlertStore } from '@/stores/alerts'

const router = useRouter()
const route = useRoute()
const calendarsStore = useCalendarsStore()
const { calendar } = storeToRefs(calendarsStore)
const alertsStore = useAlertStore()

const editedCalendar = ref<Calendar>({
  name: '',
  title: '',
  visibility: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
  description: '',
  calendarAccess: undefined,
})

function navigateBack() {
  router.push('/'+calendarName.value)
}

async function saveSettings() {
  try {
    calendarsStore.calendar = editedCalendar.value
    await calendarsStore.updateCalendar()
    navigateBack()
  } catch (err) {
    alertsStore.addAlert(err instanceof Error ? err.message : String(err),'error')
  }
}

const calendarName = computed(() => {
  return route.path.replace('/edit', '').substring(1)
})

async function loadCalendar() {
  await calendarsStore.loadCalendar(calendarName.value)
  if (calendar.value) {
    editedCalendar.value = { ...calendar.value }
  }
}

onMounted(async () => {
  await loadCalendar()
})

watch(calendar, (newVal) => {
  if (newVal) {
    editedCalendar.value = { ...newVal }
  }
})

watch(
  () => route.path,
  async (newCalendarName) => {
    if (newCalendarName) {
      await loadCalendar()
    }
  }
)
</script>

<style></style>
