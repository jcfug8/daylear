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
import type { CalendarEventExternal, CalendarType } from '@schedule-x/calendar'
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

// Deterministic color themes per calendar id
type ColorDef = { main: string; container: string; onContainer: string }

function fnv1aHash(str: string): number {
  // 32-bit FNV-1a
  let hash = 0x811c9dc5
  for (let i = 0; i < str.length; i += 1) {
    hash ^= str.charCodeAt(i)
    hash = (hash >>> 0) * 0x01000193
  }
  return hash >>> 0
}

const DISTINCT_HUES = [
  0, 18, 35, 50, 70, 95, 120, 145, 165, 190, 210, 230, 255, 275, 300, 320, 340
]

function colorThemeForId(id: string): { lightColors: ColorDef; darkColors: ColorDef } {
  const h32 = fnv1aHash(id)
  const hue = DISTINCT_HUES[h32 % DISTINCT_HUES.length]
  const satMain = [68, 76, 84][(h32 >>> 8) % 3]
  const lightMain = [42, 50, 58][(h32 >>> 16) % 3]

  const lightColors: ColorDef = {
    main: `hsl(${hue} ${satMain}% ${lightMain}%)`,
    container: `hsl(${hue} 38% 96%)`,
    onContainer: '#1a1a1a',
  }
  const darkColors: ColorDef = {
    main: `hsl(${hue} ${Math.min(satMain + 4, 90)}% ${Math.min(lightMain + 12, 70)}%)`,
    container: `hsl(${hue} 26% 22%)`,
    onContainer: '#f0f0f0',
  }
  return { lightColors, darkColors }
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

const scheduleCalendars = computed<Record<string, CalendarType>>(() => {
  const source = props.calendars ?? []
  const result: Record<string, CalendarType> = {}
  source.forEach(calendar => {
    if (!calendar.name) return
    const parts = calendar.name.split('/')
    const id = parts[parts.length - 1]
    const theme = colorThemeForId(id)
    result[id] = {
      colorName: id,
      label: calendar.title,
      lightColors: theme.lightColors,
      darkColors: theme.darkColors,
    }
  })
  return result
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
  calendars: scheduleCalendars.value,
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
