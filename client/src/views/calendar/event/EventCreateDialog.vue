<template>
  <v-dialog 
    v-model="internalOpen" 
    max-width="520" 
    @update:modelValue="onDialogModel"
  >
    <v-card>
      <v-card-title>Create New Event</v-card-title>
      <v-card-text>
        <EventForm
          v-model="form"
          :calendars="writableCalendars"
          :disabled="creating"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="close()">Cancel</v-btn>
        <v-btn color="primary" :loading="creating" @click="create()">Create</v-btn>
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
  defaultStartTime?: Date | null
  defaultEndTime?: Date | null
  calendars: Calendar[]
}>(), {
  modelValue: false,
  defaultStartTime: null,
  defaultEndTime: null,
  calendars: () => [],
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'created', event: Event): void
}>()

const internalOpen = ref<boolean>(props.modelValue)
const creating = ref(false)
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
  internalOpen.value = newValue
  if (newValue) {
    populateForm()
  }
})

// Watch internal state and emit changes
watch(internalOpen, (newValue) => {
  emit('update:modelValue', newValue)
})

// Cleanup on unmount
onBeforeUnmount(() => {
  internalOpen.value = false
  creating.value = false
})

function populateForm() {
  console.log('Populating form for new event')
  
  // Set default times
  let startTime = ''
  let endTime = ''
  
  if (props.defaultStartTime) {
    startTime = toLocalInput(props.defaultStartTime.toISOString())
  } else {
    // Default to current time rounded to next 15 minutes
    const now = new Date()
    now.setMinutes(Math.ceil(now.getMinutes() / 15) * 15, 0, 0)
    startTime = toLocalInput(now.toISOString())
  }
  
  if (props.defaultEndTime) {
    endTime = toLocalInput(props.defaultEndTime.toISOString())
  } else if (startTime) {
    // Default to 1 hour after start time
    const startDate = new Date(startTime)
    const endDate = new Date(startDate.getTime() + 60 * 60 * 1000)
    endTime = toLocalInput(endDate.toISOString())
  }

  const formData = {
    calendarName: null,
    title: '',
    start: startTime,
    end: endTime,
    description: '',
    recurrenceRule: '',
  }
  
  console.log('Setting form data:', formData)
  form.value = formData
}

// Convert timestamp to local datetime-local format
function toLocalInput(timestamp: string): string {
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
  } catch {
    return ''
  }
}

function close() {
  internalOpen.value = false
}

function onDialogModel(value: boolean) {
  internalOpen.value = value
}

async function create() {
  if (!form.value.calendarName || !form.value.title || !form.value.start) {
    // TODO: Show validation error
    return
  }
  
  creating.value = true
  try {
    const startIso = new Date(form.value.start).toISOString()
    const endIso = new Date(form.value.end || form.value.start).toISOString()
    
    // Extract calendar ID from calendar name
    const calendarId = form.value.calendarName.split('/').pop()
    if (!calendarId) {
      throw new Error('Invalid calendar ID')
    }
    
    const newEvent = await eventService.CreateEvent({
      event: {
        title: form.value.title,
        startTime: startIso as unknown as string,
        endTime: endIso as unknown as string,
        description: form.value.description || undefined,
        // Add recurrence rule if specified
        recurrenceRule: form.value.recurrenceRule || undefined,
        // Add required fields with default values
        name: undefined,
        location: undefined,
        uri: undefined,
        overridenStartTime: undefined,
        excludedTimes: undefined,
        additionalTimes: undefined,
        parentEvent: undefined,
        alarms: undefined,
        geo: undefined,
        recurrenceEndTime: undefined,
      },
      parent: form.value.calendarName
    })
    
    emit('created', newEvent)
    close()
  } catch (error) {
    console.error('Error creating event:', error)
    // TODO: Show error message to user
  } finally {
    creating.value = false
  }
}
</script>

<style scoped>
</style>


