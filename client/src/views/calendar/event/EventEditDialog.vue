<template>
  <v-dialog 
    :key="`edit-dialog-${event?.name || 'unknown'}`"
    v-model="internalOpen" 
    max-width="520" 
    @update:modelValue="onDialogModel"
  >
    <v-card>
      <v-card-title>Edit Event</v-card-title>
      <v-card-text>
        <EventForm
          :key="`form-${event?.name || 'unknown'}`"
          v-model="form"
          :calendars="writableCalendars"
          :disabled="updating"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="close()">Cancel</v-btn>
        <v-btn color="primary" :loading="updating" @click="update()">Update</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed, onBeforeUnmount } from 'vue'
import type { Calendar, Event } from '@/genapi/api/calendars/calendar/v1alpha1'
import { eventService } from '@/api/api'
import EventForm, { type EventFormData } from './EventForm.vue'
import type { CalendarEventExternal } from '@schedule-x/calendar'
import { RRule } from 'rrule'

const props = withDefaults(defineProps<{
  modelValue: boolean
  displayEvent: CalendarEventExternal | null
  event: Event | null
  calendars: Calendar[]
}>(), {
  modelValue: false,
  displayEvent: null,
  event: null,
  calendars: () => [],
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'updated', event: Event): void
}>()

const internalOpen = ref<boolean>(props.modelValue)
const updating = ref(false)
const form = ref<EventFormData>({
  calendarName: null,
  title: '',
  start: '',
  end: '',
  description: '',
  recurrenceRule: '',
})

// Filter calendars to only show those where user has write access
const writableCalendars = computed(() => {
  return props.calendars.filter(calendar => {
    if (!calendar) return false
    
    // Check if user has write permission to this specific calendar
    const calendarAccess = calendar.calendarAccess
    if (!calendarAccess?.permissionLevel) return false
    
    return calendarAccess.permissionLevel === 'PERMISSION_LEVEL_WRITE' || 
           calendarAccess.permissionLevel === 'PERMISSION_LEVEL_ADMIN'
  })
})

// Watch for prop changes and update internal state
watch(() => props.modelValue, (newValue) => {
  console.log('EventEditDialog modelValue changed:', newValue)
  internalOpen.value = newValue
  if (props.event && props.displayEvent && newValue) {
    console.log('Populating form from modelValue change')
    populateForm(props.event, props.displayEvent)
  }
})

// Watch internal state and emit changes
watch(internalOpen, (newValue) => {
  console.log('EventEditDialog internalOpen changed:', newValue)
  emit('update:modelValue', newValue)
})

// Watch for event changes and populate form
watch(() => props.event, (event) => {
  console.log('EventEditDialog event changed:', event)
  if (event && props.displayEvent && internalOpen.value) {
    console.log('Populating form from event change')
    populateForm(event, props.displayEvent)
  }
}, { immediate: true })

// Also watch for when dialog opens to populate form if event is available
watch(internalOpen, (isOpen) => {
  console.log('EventEditDialog internalOpen changed (form population):', isOpen)
  if (props.event && props.displayEvent && isOpen) {
    console.log('Populating form from dialog open')
    populateForm(props.event, props.displayEvent)
  }
})

// Cleanup on unmount
onBeforeUnmount(() => {
  internalOpen.value = false
  updating.value = false
})

function populateForm(event: Event, displayEvent: CalendarEventExternal) {
  console.log('Populating form with event:', event)
  
  // Extract calendar name from event name (assuming format: calendars/{calendarId}/events/{eventId})
  let calendarName = null
  if (event.name) {
    const parts = event.name.split('/')
    const calendarIndex = parts.findIndex(p => p === 'calendars')
    if (calendarIndex >= 0 && parts.length > calendarIndex + 1) {
      calendarName = `calendars/${parts[calendarIndex + 1]}`
    }
  }
  
  // If we couldn't extract from event name, try to find it from the calendars prop
  if (!calendarName && props.calendars.length > 0) {
    // Find the calendar that contains this event
    const eventCalendar = props.calendars.find(cal => {
      if (!cal.name) return false
      return event.name?.startsWith(cal.name)
    })
    if (eventCalendar) {
      calendarName = eventCalendar.name
    }
  }

  // Convert timestamps to local datetime-local format
  const toLocalInput = (timestamp: string | undefined) => {
    if (!timestamp) return ''
    try {
      const date = new Date(timestamp)
      if (isNaN(date.getTime())) return ''
      
      const pad = (n: number) => String(n).padStart(2, '0')
      const yyyy = date.getFullYear()
      const mm = pad(date.getMonth() + 1)
      const dd = pad(date.getDate())
      const hh = pad(date.getHours())
      const mi = pad(date.getMinutes())
      return `${yyyy}-${mm}-${dd}T${hh}:${mi}`
    } catch (error) {
      console.warn('Error converting timestamp to local input:', timestamp, error)
      return ''
    }
  }

  const formData = {
    calendarName: calendarName || null,
    title: event.title || '',
    start: toLocalInput(displayEvent.start),
    end: toLocalInput(displayEvent.end),
    description: event.description || '',
    recurrenceRule: event.recurrenceRule || '',
  }
  
  console.log('Setting form data:', formData)
  form.value = formData
}

function close() {
  internalOpen.value = false
}

function onDialogModel(value: boolean) {
  internalOpen.value = value
}

const futureOnlyRecurringEditFields: string[] = [
  'recurrenceRule'
]

async function update() {
  if (!form.value.calendarName || !props.event ||!props.event?.name || !props.displayEvent) return
  const oldEvent = props.event
  updating.value = true
  try {
    // determine the time difference old start and end with the new start and end
    const displayStart = new Date(props.displayEvent.start)
    const displayEnd = new Date(props.displayEvent.end)
    const newStart = new Date(form.value.start)
    const newEnd = new Date(form.value.end)
    const endTimeDiff = newEnd.getTime() - displayEnd.getTime()
    const startTimeDiff = newStart.getTime() - displayStart.getTime()

    // create new event start and event end times by adding the time difference to the props.event start and end
    const calculatedStart = new Date(props.event.startTime as string)
    const calculatedEnd = new Date(props.event.endTime as string)
    calculatedStart.setTime(calculatedStart.getTime() + startTimeDiff)
    calculatedEnd.setTime(calculatedEnd.getTime() + endTimeDiff)

    const startIso = new Date(calculatedStart).toISOString()
    const endIso = new Date(calculatedEnd).toISOString()
    
    let updatedEvent: Event | null = null

    const newEvent = {
      name: props.event.name,
      title: form.value.title || 'Updated Event',
      startTime: startIso as unknown as string,
      endTime: endIso as unknown as string,
      description: form.value.description || undefined,
      // Update recurrence rule if specified
      recurrenceRule: form.value.recurrenceRule || undefined,
      // Preserve other fields from original event
      location: props.event.location,
      uri: props.event.uri,
      overridenStartTime: props.event.overridenStartTime,
      excludedTimes: props.event.excludedTimes,
      additionalTimes: props.event.additionalTimes,
      parentEventId: props.event.parentEventId,
      alarms: props.event.alarms,
      geo: props.event.geo,
      recurrenceEndTime: props.event.recurrenceEndTime,
    }
    
    // normalize event times
    oldEvent.startTime = new Date(oldEvent.startTime as string).toISOString()
    oldEvent.endTime = new Date(oldEvent.endTime as string).toISOString()
    const diff = Object.keys(newEvent).filter(key => newEvent[key as keyof Event] !== oldEvent[key as keyof Event])
    console.log('diff', diff)
    if (diff.length === 0) { // no changes

    } else if (props.event.recurrenceRule !== undefined) { // this event was recurring
      // log out the fields that are different between the new event and the old event
      if (!newEvent.recurrenceRule) {
        alert('only allow all events')
        // if this isn't the first event, then actually take the old recurrence rule
        // and tack on the until with the new start time
        if (displayStart.toISOString() != oldEvent.startTime) {
          const rrule = RRule.fromString(props.event.recurrenceRule)
          rrule.options.until = displayStart
          newEvent.recurrenceRule = rrule.toString().replace('RRULE:', '')
        }
      } else if (futureOnlyRecurringEditFields.some(field => diff.includes(field))) {
        alert('only allow all future events')
      } else {
        alert('allow this event or all future events')
      }
    }

    updatedEvent = await eventService.UpdateEvent({
      event: newEvent,
      updateMask: undefined
    })
    
    if (updatedEvent) {
      emit('updated', updatedEvent)
    }
    close()
  } finally {
    updating.value = false
  }
}
</script>

<style scoped>
</style>
