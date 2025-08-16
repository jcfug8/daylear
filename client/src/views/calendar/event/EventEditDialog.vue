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

  <!-- Recurring Event Update Choice Modal -->
  <v-dialog v-model="showRecurringChoiceModal" max-width="400" persistent>
    <v-card>
      <v-card-title>Update Recurring Event</v-card-title>
      <v-card-text>
        <p class="mb-4">This is a recurring event. How would you like to apply the changes?</p>
        <v-radio-group v-model="selectedUpdateOption" class="mt-4">
          <v-radio
            v-if="pendingUpdateData?.availableOptions?.includes('this-event')"
            value="this-event"
            label="This event only"
            class="mb-2"
          />
          <v-radio
            v-if="pendingUpdateData?.availableOptions?.includes('all-future')"
            value="all-future"
            label="This and all future events"
            class="mb-2"
          />
          <v-radio
            v-if="pendingUpdateData?.availableOptions?.includes('all-events')"
            value="all-events"
            label="All events in the series"
            class="mb-2"
          />
        </v-radio-group>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="cancelRecurringUpdate()">Cancel</v-btn>
        <v-btn color="primary" @click="confirmRecurringUpdate()">Confirm</v-btn>
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
const showRecurringChoiceModal = ref(false)
const selectedUpdateOption = ref<'this-event' | 'all-future' | 'all-events'>('this-event')
const pendingUpdateData = ref<{
  newEvent: Event
  oldEvent: Event
  displayEvent: CalendarEventExternal
  availableOptions?: ('this-event' | 'all-future' | 'all-events')[]
} | null>(null)

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

    } else if (props.event.recurrenceRule !== undefined && props.event.recurrenceRule !== '') { // this event was recurring
      // log out the fields that are different between the new event and the old event
      if (!newEvent.recurrenceRule) {
        // if this isn't the first event, then actually take the old recurrence rule
        // and tack on the until with the new start time
        if (displayStart.toISOString() != oldEvent.startTime) {
          const rrule = RRule.fromString(props.event.recurrenceRule)
          rrule.origOptions.until = new Date(displayStart.getUTCFullYear(), displayStart.getUTCMonth(), displayStart.getUTCDate(), displayStart.getUTCHours(), displayStart.getUTCMinutes(), displayStart.getUTCSeconds())
          rrule.origOptions.count = null
          newEvent.recurrenceRule = RRule.optionsToString(rrule.origOptions)
        }
        
        // Removing recurrence rule - only show "all events" option
        pendingUpdateData.value = {
          newEvent,
          oldEvent,
          displayEvent: props.displayEvent,
          availableOptions: ['all-events']
        }
        selectedUpdateOption.value = 'all-events' // Set default for this case
        showRecurringChoiceModal.value = true
        updating.value = false // Reset updating state since we're waiting for user choice
        return // Don't proceed with update yet, wait for user choice
      } else if (futureOnlyRecurringEditFields.some(field => diff.includes(field))) {
        // Changing future-only fields - only show "all future events" option
        pendingUpdateData.value = {
          newEvent,
          oldEvent,
          displayEvent: props.displayEvent,
          availableOptions: ['all-future']
        }
        selectedUpdateOption.value = 'all-future' // Set default for this case
        showRecurringChoiceModal.value = true
        updating.value = false // Reset updating state since we're waiting for user choice
        return // Don't proceed with update yet, wait for user choice
      } else {
        // Other changes - show all three options
        pendingUpdateData.value = {
          newEvent,
          oldEvent,
          displayEvent: props.displayEvent,
          availableOptions: ['this-event', 'all-future', 'all-events']
        }
        selectedUpdateOption.value = 'this-event' // Set default for this case
        showRecurringChoiceModal.value = true
        updating.value = false // Reset updating state since we're waiting for user choice
        return // Don't proceed with update yet, wait for user choice
      }
    } else {
      // Non-recurring event or no recurrence rule, proceed with normal update
      updatedEvent = await updateEvent(newEvent)
    }

    
    if (updatedEvent) {
      emit('updated', updatedEvent)
    }
    close()
  } catch (error) {
    console.error('Error updating event:', error)
  } finally {
    updating.value = false
  }
}


async function updateEvent(event: Event): Promise<Event> {
  try {
    const updatedEvent = await eventService.UpdateEvent({
        event: event,
        updateMask: undefined
      })
      return updatedEvent
  } catch (error) {
    throw error
  }
}

async function createEvent(event: Event, parent: string): Promise<Event> {
  try {
    const createdEvent = await eventService.CreateEvent({
      event: event,
      parent: parent,
    })
    return createdEvent
  } catch (error) {
    throw error
  }
}

function cancelRecurringUpdate() {
  showRecurringChoiceModal.value = false
  pendingUpdateData.value = null
  selectedUpdateOption.value = 'this-event'
  updating.value = false // Reset updating state
}

async function confirmRecurringUpdate() {
  if (!pendingUpdateData.value) return
  
  updating.value = true // Set updating state for the actual update
  
  const { newEvent, oldEvent, displayEvent } = pendingUpdateData.value
  let updatedEvent: Event | null = null

  if (!oldEvent.recurrenceRule) {
    console.error('No recurrence rule found for event')
    return
  }

  if (!newEvent.name) {
    console.error('No name found for event')
    return
  }

  if (!newEvent.startTime || !newEvent.endTime) {
    console.log('no start or end time found for event')
    return
  }
  
  try {
    switch (selectedUpdateOption.value) {
      case 'this-event':
        // if the parent event id is set then just update the event
        if (newEvent.parentEventId) {
          updatedEvent = await updateEvent(newEvent)
        } else {
          newEvent.startTime = new Date(displayEvent.start).toISOString()
          newEvent.endTime = new Date(displayEvent.end).toISOString()
          // update the old event's excluded times the the start time of
          // the new event
          oldEvent.excludedTimes?.push(newEvent.startTime)
          updatedEvent = await updateEvent(oldEvent)

          // if the parent event id is not set then create a new event
          const calendarName = newEvent.name?.substring(0, newEvent.name.indexOf('/events/'))
          updatedEvent = await createEvent(newEvent, calendarName)
        }
        break
        
      case 'all-future':
        // TODO: this is close. Need to fix
        // - be sure end of old recurrence rule is before the new start time so the event being editted doesn't "show" twice
        // - what to do if the one being editted is the first occurence?
        const displayStartTime = new Date(displayEvent.start)
        // Update this and all future events by modifying the old recurrence rule to end
        // before the new start time
        const rrule = RRule.fromString(oldEvent.recurrenceRule)
        const options = rrule.origOptions
        options.dtstart = new Date(oldEvent.startTime as string)
        const before = new Date(displayStartTime.getUTCFullYear(), displayStartTime.getUTCMonth(), displayStartTime.getUTCDate(), displayStartTime.getUTCHours(), displayStartTime.getUTCMinutes(), displayStartTime.getUTCSeconds())
        const previousOccurence = rrule.before(before)
        if (previousOccurence) {
          options.until = previousOccurence
        } else {
          console.log('no previous occurence found, setting until to before')
          options.until = before
        }
        options.count = null
        oldEvent.recurrenceRule = RRule.optionsToString(options)
        return
        updatedEvent = await updateEvent(oldEvent)

        newEvent.startTime = new Date(displayEvent.start).toISOString()
        newEvent.endTime = new Date(displayEvent.end).toISOString()

        // create new event
        const calendarName = newEvent.name?.substring(0, newEvent.name.indexOf('/events/'))
        await createEvent(newEvent, calendarName)
        break
        
      case 'all-events':
        // Update all events in the series
        updatedEvent = await updateEvent(newEvent)
        break
    }
    
    if (updatedEvent) {
      emit('updated', updatedEvent)
    }
    close()
  } catch (error) {
    console.error('Error updating recurring event:', error)
  } finally {
    updating.value = false
    showRecurringChoiceModal.value = false
    pendingUpdateData.value = null
    selectedUpdateOption.value = 'this-event'
  }
}
</script>

<style scoped>
</style>
