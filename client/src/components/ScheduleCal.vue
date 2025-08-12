<template>
  <div id="schedule-cal-wrapper" ref="wrapperRef" :style="wrapperStyle">
    <ScheduleXCalendar :calendar-app="calendarApp" />
  </div>
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
import { computed, watch, ref, onMounted, onBeforeUnmount } from 'vue'
import type { CSSProperties } from 'vue'
import type { Event, Calendar } from '@/genapi/api/calendars/calendar/v1alpha1/index'
import type { CalendarEventExternal } from '@schedule-x/calendar'
import { createEventsServicePlugin } from '@schedule-x/events-service'

const eventsServicePlugin = createEventsServicePlugin();

const props = withDefaults(defineProps<{
  events?: Event[]
  loading?: boolean
  calendars?: Calendar[]
}>(), {
  events: () => [],
  calendars: () => [],
  loading: false,
})

function toScheduleXDateTime(value?: string | null): string | undefined {
  if (!value) return undefined
  const d = new Date(value)
  if (isNaN(d.getTime())) return undefined
  const pad = (n: number) => String(n).padStart(2, '0')
  const Y = d.getFullYear()
  const M = pad(d.getMonth() + 1)
  const D = pad(d.getDate())
  const h = pad(d.getHours())
  const m = pad(d.getMinutes())
  return `${Y}-${M}-${D} ${h}:${m}`
}

function parseName(name?: string | null): { id?: string; calendarId?: string } {
  if (!name) return {}
  const parts = name.split('/').filter(Boolean)
  if (parts.length === 0) return {}
  const id = parts[parts.length - 1]
  const calIdx = parts.findIndex(p => p === 'calendars')
  const calendarId = calIdx >= 0 && parts.length > calIdx + 1 ? parts[calIdx + 1] : undefined
  return { id, calendarId }
}

const scheduleEvents = computed<CalendarEventExternal[]>(() => {
  const source = props.events ?? []
  return source
    .map((event) => {
      if (!event) return undefined
      const { id, calendarId } = parseName(event.name)
      const start = toScheduleXDateTime(event.startTime as unknown as string)
      const end = toScheduleXDateTime(event.endTime as unknown as string)
      if (!id || !calendarId || !start) return undefined
      return {
        id,
        start,
        end,
        title: event.title ?? '',
        description: event.description ?? '',
        calendarId,
      } as CalendarEventExternal
    })
    .filter((e): e is CalendarEventExternal => Boolean(e))
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
  events: scheduleEvents.value,
  plugins: [eventsServicePlugin],
  calendars: {
    "4": {
      colorName: "4",
      lightColors: {
        main: '#f9d71c',
        container: '#fff5aa',
        onContainer: '#594800',
      },
      darkColors: {
        main: '#fff5c0',
        onContainer: '#fff5de',
        container: '#a29742',
      },
    },
    "5": {
      colorName: "5",
      lightColors: {
        main: '#f91c45',
        container: '#ffd2dc',
        onContainer: '#59000d',
      },
      darkColors: {
        main: '#ffc0cc',
        onContainer: '#ffdee6',
        container: '#a24258',
      },
    },
  },
})

watch(scheduleEvents, (newEvents) => {
  eventsServicePlugin.set(newEvents)
})

// Dynamically size wrapper to fill available viewport space
const wrapperRef = ref<HTMLElement | null>(null)
const wrapperHeight = ref<string>('auto')

function updateWrapperHeight() {
  if (!wrapperRef.value) return
  const rect = wrapperRef.value.getBoundingClientRect()
  const available = Math.max(window.innerHeight - rect.top, 0)
  wrapperHeight.value = `${available}px`
}

onMounted(() => {
  updateWrapperHeight()
  window.addEventListener('resize', updateWrapperHeight)
  window.addEventListener('orientationchange', updateWrapperHeight)
  window.addEventListener('scroll', updateWrapperHeight)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateWrapperHeight)
  window.removeEventListener('orientationchange', updateWrapperHeight)
  window.removeEventListener('scroll', updateWrapperHeight)
})

const wrapperStyle = computed<CSSProperties>(() => ({
  width: '100%',
  height: wrapperHeight.value,
  overflow: 'auto',
  overflowX: 'hidden',
  display: 'flex',
  flexDirection: 'column',
  minHeight: 0,
}))
</script>

<style>
#schedule-cal-wrapper {
  width: 100%;
  height: auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
  max-height: calc(100dvh - var(--v-layout-top, 0px));
  overflow: auto;
  overflow-x: hidden;
}

.sx-vue-calendar-wrapper {
  width: 100% !important;
  height: 100% !important;
  max-width: 100%;
  max-height: 100%;
  min-height: 0;
}
</style>
