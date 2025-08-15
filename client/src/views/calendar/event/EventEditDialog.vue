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

const props = withDefaults(defineProps<{
  modelValue: boolean
  event: Event | null
  calendars: Calendar[]
}>(), {
  modelValue: false,
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
  if (newValue && props.event) {
    console.log('Populating form from modelValue change')
    populateForm(props.event)
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
  if (event && internalOpen.value) {
    console.log('Populating form from event change')
    populateForm(event)
  }
}, { immediate: true })

// Also watch for when dialog opens to populate form if event is available
watch(internalOpen, (isOpen) => {
  console.log('EventEditDialog internalOpen changed (form population):', isOpen)
  if (isOpen && props.event) {
    console.log('Populating form from dialog open')
    populateForm(props.event)
  }
})

// Cleanup on unmount
onBeforeUnmount(() => {
  internalOpen.value = false
  updating.value = false
})

function populateForm(event: Event) {
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
    start: toLocalInput(event.startTime as unknown as string),
    end: toLocalInput(event.endTime as unknown as string),
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

async function update() {
  if (!form.value.calendarName || !props.event?.name) return
  
  updating.value = true
  try {
    const startIso = new Date(form.value.start).toISOString()
    const endIso = new Date(form.value.end || form.value.start).toISOString()
    
    const updatedEvent = await eventService.UpdateEvent({
      event: {
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
      },
      updateMask: undefined
    })
    
    emit('updated', updatedEvent)
    close()
  } finally {
    updating.value = false
  }
}
</script>

<style scoped>
</style>
