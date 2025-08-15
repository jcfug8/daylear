<template>
  <div id="schedule-cal-wrapper" ref="wrapperRef" :style="wrapperStyle">
    <div class="d-flex align-center mb-2">
      <v-spacer />
      <v-btn 
        v-if="canCreateEvents"
        color="primary" 
        prepend-icon="mdi-plus"
        @click="showCreateDialog = true"
      >
        Create Event
      </v-btn>
    </div>
    <ScheduleXCalendar :calendar-app="calendarApp" />
  </div>
  
  <!-- Event View Dialog -->
  <EventViewDialog
    v-model="showEventDialog"
    :event="selectedEvent"
    :calendar="selectedCalendar"
    :event-occurrence="selectedEventOccurrence"
    @updated="handleEventUpdated"
  />
  
  <!-- Event Create Dialog -->
  <EventCreateDialog
    v-if="showCreateDialog"
    v-model="showCreateDialog"
    :calendars="writableCalendars"
    @created="handleEventCreated"
  />
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
import { createEventRecurrencePlugin, createEventsServicePlugin } from "@schedule-x/event-recurrence";
import EventViewDialog from './EventViewDialog.vue'
import EventCreateDialog from './EventCreateDialog.vue'
// Removed RRule and Luxon imports - using Schedule-X built-in recurrence

const eventsServicePlugin = createEventsServicePlugin();
const eventRecurrencePlugin = createEventRecurrencePlugin();

const props = withDefaults(defineProps<{
  events?: Event[]
  loading?: boolean
  calendars?: Calendar[]
}>(), {
  events: () => [],
  calendars: () => [],
  loading: false,
})

// Event dialog state
const showEventDialog = ref(false)
const selectedEvent = ref<Event | null>(null)
const selectedEventOccurrence = ref<CalendarEventExternal | null>(null)
const selectedCalendar = ref<Calendar | null>(null)

// Event create dialog state
const showCreateDialog = ref(false)

const emit = defineEmits<{
  (e: 'updated', event: Event): void
  (e: 'created', event: Event): void
}>()

// Handle event updated from dialog
function handleEventUpdated(event: Event) {
  console.log('Event updated:', event)
  // Close the dialog after update
  showEventDialog.value = false
  selectedEvent.value = null
  selectedCalendar.value = null

  emit('updated', event)
}

// Handle event created from dialog
function handleEventCreated(event: Event) {
  console.log('Event created:', event)
  // Close the dialog after creation
  showCreateDialog.value = false
  emit('created', event)
}

// Clean up dialog state when dialog closes
watch(showEventDialog, (isOpen) => {
  if (!isOpen) {
    // Reset state when dialog closes
    selectedEvent.value = null
    selectedCalendar.value = null
  }
})

function toScheduleXDateTime(value?: string | null): string | undefined {
  if (!value || typeof value !== 'string') return undefined
  
  try {
    const d = new Date(value)
    if (isNaN(d.getTime())) return undefined
    
    const pad = (n: number) => String(n).padStart(2, '0')
    const Y = d.getFullYear()
    const M = pad(d.getMonth() + 1)
    const D = pad(d.getDate())
    const h = pad(d.getHours())
    const m = pad(d.getMinutes())
    return `${Y}-${M}-${D} ${h}:${m}`
  } catch (error) {
    console.warn('Error parsing date:', value, error)
    return undefined
  }
}

function parseName(name?: string | null): { id?: string; calendarId?: string } {
  if (!name || typeof name !== 'string') return {}
  const parts = name.split('/').filter(Boolean)
  if (parts.length < 2) return {}
  
  const id = parts[parts.length - 1]
  const calIdx = parts.findIndex(p => p === 'calendars')
  const calendarId = calIdx >= 0 && parts.length > calIdx + 1 ? parts[calIdx + 1] : undefined
  
  return { id, calendarId }
}

// Deterministic color themes per calendar id
type ColorDef = { main: string; container: string; onContainer: string }

function fnv1aHash(str: string): number {
  if (!str || typeof str !== 'string') return 0
  
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
  if (!id || typeof id !== 'string') {
    // Fallback to a default theme
    return {
      lightColors: {
        main: 'hsl(0 68% 50%)',
        container: 'hsl(0 38% 96%)',
        onContainer: '#1a1a1a',
      },
      darkColors: {
        main: 'hsl(0 72% 62%)',
        container: 'hsl(0 26% 22%)',
        onContainer: '#f0f0f0',
      }
    }
  }
  
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
    .filter((event): event is Event => Boolean(event))
    .map((event) => {
      const { id, calendarId } = parseName(event.name)
      const start = toScheduleXDateTime(event.startTime as string)
      const end = toScheduleXDateTime(event.endTime as string) || ''

      if (!id || !calendarId || !start) return null

      const rrule = event.recurrenceRule ? event.recurrenceRule.replace('RRULE:', '') : null
      
      const baseEvent: CalendarEventExternal = {
        name: event.name,
        id,
        start,
        end,
        title: event.title ?? '',
        description: event.description ?? '',
        calendarId,
        rrule,
      }
      
      return baseEvent
    })
    .filter((event): event is CalendarEventExternal => event !== null)
})

const scheduleCalendars = computed<Record<string, CalendarType>>(() => {
  const source = props.calendars ?? []
  const result: Record<string, CalendarType> = {}
  source
    .filter((calendar): calendar is Calendar => Boolean(calendar) && Boolean(calendar.name))
    .forEach(calendar => {
      const parts = calendar.name!.split('/')
      const id = parts[parts.length - 1]
      if (!id) return
      const theme = colorThemeForId(id)
      result[id] = {
        colorName: id,
        label: calendar.title || 'Untitled Calendar',
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
  firstDayOfWeek: 0,
  events: scheduleEvents.value,
  plugins: [
    eventsServicePlugin,
    eventRecurrencePlugin,
  ],
  calendars: scheduleCalendars.value,
  callbacks: {
    onEventClick: (event) => {
      console.log('event clicked', event)
      
      try {
        // Find the corresponding Event and Calendar from our data
        const eventName = event.name
        const calendarId = event.calendarId
        
        if (!eventName || !calendarId) {
          console.warn('Event click: Missing event name or calendar ID', event)
          return
        }
        
        // Find the event in our events array
        const foundEvent = props.events.find(e => {
          return e && e.name === eventName
        })
        
        // Find the calendar in our calendars array
        const foundCalendar = props.calendars.find(c => {
          if (!c || !c.name) return false
          const parts = c.name.split('/')
          const calId = parts[parts.length - 1]
          return calId === calendarId
        })
        
        if (foundEvent && foundCalendar) {
          console.log('foundEvent', foundEvent)
          console.log('foundCalendar', foundCalendar)
          selectedEvent.value = foundEvent
          selectedEventOccurrence.value = event
          selectedCalendar.value = foundCalendar
          showEventDialog.value = true
        } else {
          console.warn('Event click: Could not find event or calendar', { 
            eventName, 
            calendarId, 
            foundEvent: !!foundEvent, 
            foundCalendar: !!foundCalendar 
          })
        }
      } catch (error) {
        console.error('Error in event click handler:', error)
      }
    }
  }
})

watch(scheduleEvents, (newEvents) => {
  try {
    if (newEvents && Array.isArray(newEvents)) {
      eventsServicePlugin.set(newEvents)
    }
  } catch (error) {
    console.error('Error updating events service plugin:', error)
  }
})

// Dynamically size wrapper to fill available viewport space
const wrapperRef = ref<HTMLElement | null>(null)
const wrapperHeight = ref<string>('auto')

function updateWrapperHeight() {
  try {
    if (!wrapperRef.value) return
    const rect = wrapperRef.value.getBoundingClientRect()
    if (!rect) return
    
    const available = Math.max(window.innerHeight - rect.top, 0)
    wrapperHeight.value = `${available}px`
  } catch (error) {
    console.warn('Error updating wrapper height:', error)
    wrapperHeight.value = 'auto'
  }
}

onMounted(() => {
  try {
    updateWrapperHeight()
    window.addEventListener('resize', updateWrapperHeight)
    window.addEventListener('orientationchange', updateWrapperHeight)
    window.addEventListener('scroll', updateWrapperHeight)
  } catch (error) {
    console.error('Error in onMounted:', error)
  }
})

onBeforeUnmount(() => {
  try {
    // Clean up event listeners
    window.removeEventListener('resize', updateWrapperHeight)
    window.removeEventListener('orientationchange', updateWrapperHeight)
    window.removeEventListener('scroll', updateWrapperHeight)
    
    // Clean up dialog state
    showEventDialog.value = false
    selectedEvent.value = null
    selectedCalendar.value = null
    showCreateDialog.value = false
    
    // Clean up calendar app if needed
    if (calendarApp && typeof calendarApp.destroy === 'function') {
      calendarApp.destroy()
    }
  } catch (error) {
    console.error('Error in onBeforeUnmount:', error)
  }
})

const wrapperStyle = computed<CSSProperties>(() => {
  try {
    return {
      width: '100%',
      height: wrapperHeight.value || 'auto',
      overflow: 'auto',
      overflowX: 'hidden',
      display: 'flex',
      flexDirection: 'column',
      minHeight: 0,
    }
  } catch (error) {
    console.warn('Error computing wrapper style:', error)
    return {
      width: '100%',
      height: 'auto',
      overflow: 'auto',
      overflowX: 'hidden',
      display: 'flex',
      flexDirection: 'column',
      minHeight: 0,
    }
  }
})

const canCreateEvents = computed(() => {
  return writableCalendars.value.length > 0
})

const writableCalendars = computed(() => {
  return props.calendars?.filter(calendar => {
    if (!calendar) return false
    
    // Check if user has write permission to this specific calendar
    const calendarAccess = calendar.calendarAccess
    if (!calendarAccess?.permissionLevel) return false
    
    return calendarAccess.permissionLevel === 'PERMISSION_LEVEL_WRITE' || 
           calendarAccess.permissionLevel === 'PERMISSION_LEVEL_ADMIN'
  }) ?? []
})
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
