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
          @validation-change="onValidationChange"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="close()">Cancel</v-btn>
        <v-btn 
          v-if="canDelete" 
          color="error" 
          variant="outlined" 
          :loading="deleting" 
          @click="deleteEvent"
        >
          Delete
        </v-btn>
        <v-btn 
          color="primary" 
          :loading="updating" 
          :disabled="isUpdateButtonDisabled"
          @click="update()"
        >
          Update
        </v-btn>
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

  <!-- Recurring Event Delete Choice Modal -->
  <v-dialog v-model="showDeleteRecurringChoiceModal" max-width="400" persistent>
    <v-card>
      <v-card-title>Delete Recurring Event</v-card-title>
      <v-card-text>
        <p class="mb-4">This is a recurring event. How would you like to delete it?</p>
        <v-radio-group v-model="selectedDeleteOption" class="mt-4">
          <v-radio
            v-if="pendingDeleteData?.availableOptions?.includes('this-event')"
            value="this-event"
            label="This event only"
            class="mb-2"
          />
          <v-radio
            v-if="pendingDeleteData?.availableOptions?.includes('all-future')"
            value="all-future"
            label="This and all future events"
            class="mb-2"
          />
          <v-radio
            v-if="pendingDeleteData?.availableOptions?.includes('all-events')"
            value="all-events"
            label="All events in the series"
            class="mb-2"
          />
        </v-radio-group>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="cancelRecurringDelete()">Cancel</v-btn>
        <v-btn color="error" @click="confirmRecurringDelete()">Delete</v-btn>
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
import { datetime, RRule } from 'rrule'

type RecurringUpdateOption = 'this-event' | 'all-future' | 'all-events'
type RecurringDeleteOption = 'this-event' | 'all-future' | 'all-events' | 'single-event'

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
  (e: 'deleted', event: Event): void
}>()

const internalOpen = ref<boolean>(props.modelValue)
const updating = ref(false)
const deleting = ref(false)
const hasValidationErrors = ref(false)
const showRecurringChoiceModal = ref(false)
const showDeleteRecurringChoiceModal = ref(false)
const selectedUpdateOption = ref<'this-event' | 'all-future' | 'all-events'>('this-event')
const selectedDeleteOption = ref<'this-event' | 'all-future' | 'all-events'>('this-event')
const pendingUpdateData = ref<{
  newEvent: Event
  oldEvent: Event
  displayEvent: CalendarEventExternal
  availableOptions?: RecurringUpdateOption[]
} | null>(null)
const pendingDeleteData = ref<{
  availableOptions?: RecurringDeleteOption[]
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

// Check if user can delete this event
const canDelete = computed(() => {
  if (!props.event?.name) return false
  
  // Extract calendar name from event name
  const parts = props.event.name.split('/')
  const calendarIndex = parts.findIndex(p => p === 'calendars')
  if (calendarIndex < 0 || parts.length <= calendarIndex + 1) return false
  
  const calendarName = `calendars/${parts[calendarIndex + 1]}`
  
  // Find the calendar and check permissions
  const calendar = props.calendars.find(cal => cal.name === calendarName)
  if (!calendar?.calendarAccess?.permissionLevel) return false
  
  return calendar.calendarAccess.permissionLevel === 'PERMISSION_LEVEL_WRITE' || 
         calendar.calendarAccess.permissionLevel === 'PERMISSION_LEVEL_ADMIN'
})

// Computed property to check if the update button should be disabled
const isUpdateButtonDisabled = computed(() => {
  return updating.value || 
         !hasValidationErrors.value || 
         !form.value.calendarName || 
         !form.value.title || 
         !form.value.start
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
      parentEvent: props.event.parentEvent,
      alarms: props.event.alarms,
      geo: props.event.geo,
      recurrenceEndTime: props.event.recurrenceEndTime,
    }
    
    // normalize event times
    oldEvent.startTime = new Date(oldEvent.startTime as string).toISOString()
    oldEvent.endTime = new Date(oldEvent.endTime as string).toISOString()
    const diff = Object.keys(newEvent).filter(key => newEvent[key as keyof Event] !== oldEvent[key as keyof Event])
    console.log('diff', diff)

    if ((diff.includes('recurrenceRule') || diff.includes('startTime')) && newEvent.recurrenceRule) { // check if we should ensure that the "unit" has the correct time
      // get the rrule
      const rrule = RRule.fromString(newEvent.recurrenceRule)
      if (rrule.origOptions.until) {
        // set the until to the new start time
        rrule.origOptions.until = datetime(rrule.origOptions.until.getUTCFullYear(), rrule.origOptions.until.getUTCMonth()+1, rrule.origOptions.until.getUTCDate(), displayStart.getUTCHours(), displayStart.getUTCMinutes(), displayStart.getUTCSeconds())
        newEvent.recurrenceRule = RRule.optionsToString(rrule.origOptions)
      }
    }

    if (diff.length === 0) { // no changes

    } else if (props.event.recurrenceRule !== undefined && props.event.recurrenceRule !== '') { // this event was recurring
      // log out the fields that are different between the new event and the old event
      if (!newEvent.recurrenceRule) {
        // if this isn't the first event, then actually take the old recurrence rule
        // and tack on the until with the new start time
        if (displayStart.toISOString() != oldEvent.startTime) {
          const rrule = RRule.fromString(props.event.recurrenceRule)
          rrule.origOptions.until = datetime(displayStart.getUTCFullYear(), displayStart.getUTCMonth()+1, displayStart.getUTCDate(), displayStart.getUTCHours(), displayStart.getUTCMinutes(), displayStart.getUTCSeconds())
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
        let availableOptions: RecurringUpdateOption[] = ['all-future']
        if (displayStart.toISOString() == oldEvent.startTime) {
          availableOptions = ['all-events']
        }
        // Changing future-only fields - only show "all future events" option
        pendingUpdateData.value = {
          newEvent,
          oldEvent,
          displayEvent: props.displayEvent,
          availableOptions: availableOptions
        }
        selectedUpdateOption.value = 'all-future' // Set default for this case
        showRecurringChoiceModal.value = true
        updating.value = false // Reset updating state since we're waiting for user choice
        return // Don't proceed with update yet, wait for user choice
      } else {
        let availableOptions: RecurringUpdateOption[] = ['this-event', 'all-future', 'all-events']
        if (displayStart.toISOString() == oldEvent.startTime) {
          availableOptions = ['all-events']
        }
       
        // Other changes - show all three options
        pendingUpdateData.value = {
          newEvent,
          oldEvent,
          displayEvent: props.displayEvent,
          availableOptions: availableOptions
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

  if (!oldEvent.name) {
    console.error('No name found for old event')
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
        if (newEvent.parentEvent) {
          updatedEvent = await updateEvent(newEvent)
        } else {
          newEvent.startTime = new Date(form.value.start).toISOString()
          newEvent.endTime = new Date(form.value.end).toISOString()
          // update the old event's excluded times the the start time of
          // the new event
          if (!oldEvent.excludedTimes) {
            oldEvent.excludedTimes = []
          }
          oldEvent.excludedTimes.push(new Date(displayEvent.start).toISOString())
          updatedEvent = await updateEvent(oldEvent)

          newEvent.parentEvent = oldEvent.name
          newEvent.recurrenceRule = ''
          newEvent.excludedTimes = []
          newEvent.overridenStartTime = new Date(displayEvent.start).toISOString()
          const calendarName = newEvent.name?.substring(0, newEvent.name.indexOf('/events/'))
          updatedEvent = await createEvent(newEvent, calendarName)
        }
        break
        
      case 'all-future':
        updateRecurrenceRulesForAllFutureEvents(newEvent, oldEvent, displayEvent)

        // update old event
        updatedEvent = await updateEvent(oldEvent)

        // create new event
        newEvent.startTime = new Date(form.value.start).toISOString()
        newEvent.endTime = new Date(form.value.end).toISOString()
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

function updateRecurrenceRulesForAllFutureEvents(newEvent: Event, oldEvent: Event, displayEvent: { start: string }) {
  // ** first lets set the until to the occurence before the new start time
  const displayStartTime = new Date(displayEvent.start)
  const oldStartTime = new Date(oldEvent.startTime as string)
  
  // Update this and all future events by modifying the old recurrence rule to end
  // before the new start time
  let oldRrule = RRule.fromString(oldEvent.recurrenceRule as string)
  
  oldRrule.origOptions.dtstart = datetime(oldStartTime.getUTCFullYear(), oldStartTime.getUTCMonth() + 1, oldStartTime.getUTCDate(), oldStartTime.getUTCHours(), oldStartTime.getUTCMinutes(), oldStartTime.getUTCSeconds())
  const before = datetime(displayStartTime.getUTCFullYear(), displayStartTime.getUTCMonth() + 1, displayStartTime.getUTCDate(), displayStartTime.getUTCHours(), displayStartTime.getUTCMinutes(), displayStartTime.getUTCSeconds())
  oldRrule = new RRule(oldRrule.origOptions) // this is needed to get the before method to work with the newly set dtstart
  const previousOccurence = oldRrule.before(before)

  if (previousOccurence) {
    oldRrule.origOptions.until = previousOccurence
  } else {
    console.log('no previous occurence found, setting until to before')
    oldRrule.origOptions.until = before
  }
  
  // if there was a count we need to know how many are let so we can set them onto the new event
  if (oldRrule.origOptions.count) {
     // Get all occurrences from start to target date
    let index = 0
    oldRrule.between(new Date(oldEvent.startTime as string), new Date(displayEvent.start), true, (date, i) => {
      index = i
      return true // continue until we reach target
    })
    const newRrule = RRule.fromString(newEvent.recurrenceRule as string)
    newRrule.origOptions.count = oldRrule.origOptions.count - index
    newEvent.recurrenceRule = RRule.optionsToString(newRrule.origOptions)
    oldRrule.origOptions.count = null
  }

  // now we need to set the new recurrence rule onto the old event
  oldEvent.recurrenceRule = RRule.optionsToString(oldRrule.origOptions).split('\n')[1]
}

// Delete event functionality
async function deleteEvent() {
  if (!props.event || !props.displayEvent) return

  
  // If this is a recurring event, show the choice modal
  if (props.event.recurrenceRule && props.event.recurrenceRule !== '') {
    let availableOptions: RecurringDeleteOption[] = ['this-event', 'all-future', 'all-events']
    const eventStart = new Date(props.event.startTime as string)
    const displayStart = new Date(props.displayEvent.start)
    if (displayStart.toISOString() == eventStart.toISOString()) {
      availableOptions = ['all-events']
    }
    pendingDeleteData.value = {
      availableOptions: availableOptions
    }
    showDeleteRecurringChoiceModal.value = true
  } else {
    // Non-recurring event, delete directly
    await performDelete('single-event')
  }
}

function cancelRecurringDelete() {
  showDeleteRecurringChoiceModal.value = false
  selectedDeleteOption.value = 'this-event'
  deleting.value = false
}

async function confirmRecurringDelete() {
  if (!props.event) return
  
  showDeleteRecurringChoiceModal.value = false
  await performDelete(selectedDeleteOption.value)
}

async function performDelete(option: RecurringDeleteOption) {
  if (!props.event || !props.displayEvent) return
  
  deleting.value = true
  try {
    const eventToDelete = props.event
    
    switch (option) {
      case 'this-event':
        if (eventToDelete.parentEvent) {
          // Just delete this event if it has a parent
          await eventService.DeleteEvent({ name: eventToDelete.name })
        } else {
          // Add exdate to the old event and delete this one
          if (!eventToDelete.excludedTimes) {
            eventToDelete.excludedTimes = []
          }
          eventToDelete.excludedTimes.push(eventToDelete.startTime as string)
          await eventService.UpdateEvent({ event: eventToDelete, updateMask: undefined })
        }
        break
        
      case 'all-future':
        if (eventToDelete.parentEvent) {
          // Find the parent event and update its recurrence rule
          const parentEvent = await findParentEvent(eventToDelete.parentEvent)
          if (parentEvent) {
            updateRecurrenceRulesForAllFutureEventsDelete(parentEvent, eventToDelete, props.displayEvent)
            await eventService.UpdateEvent({ event: parentEvent, updateMask: undefined })
          }
          
          // Delete all future events with this parent
          await deleteFutureEventsWithParent(eventToDelete.parentEvent)
        } else {
          // This is the parent event, update its recurrence rule
          updateRecurrenceRulesForAllFutureEventsDelete(eventToDelete, eventToDelete, props.displayEvent)
          await eventService.UpdateEvent({ event: eventToDelete, updateMask: undefined })
          
          // Delete all future events with this event as parent
          await deleteFutureEventsWithParent(eventToDelete.name || '')
        }
        break
        
      case 'all-events':
        if (eventToDelete.parentEvent) {
          // Delete the parent event and all its children
          const parentEvent = await findParentEvent(eventToDelete.parentEvent)
          if (parentEvent) {
            await eventService.DeleteEvent({ name: parentEvent.name })
          }
        } else {
          // Delete this event and all its children
          await eventService.DeleteEvent({ name: eventToDelete.name || '' })
        }
        break
      case 'single-event':
        await eventService.DeleteEvent({ name: eventToDelete.name || '' })
        break
    }
    
    emit('deleted', eventToDelete)
    close()
  } catch (error) {
    console.error('Error deleting event:', error)
  } finally {
    deleting.value = false
  }
}

// Helper function to find parent event
async function findParentEvent(parentName: string): Promise<Event | null> {
  try {
    // Try to get the parent event by name
    const parentEvent = await eventService.GetEvent({ name: parentName })
    return parentEvent
  } catch (error) {
    console.warn('Could not find parent event:', error)
    return null
  }
}

// Helper function to delete future events with a specific parent
async function deleteFutureEventsWithParent(parentName: string) {
  // This would require a ListEvents API call with filtering
  // For now, we'll rely on the parent component to refresh the events list
  console.log('Deleting future events with parent:', parentName)
}

// Helper function to update recurrence rules for delete operations
function updateRecurrenceRulesForAllFutureEventsDelete(oldEvent: Event, currentEvent: Event, displayEvent: { start: string }) {
  const displayStartTime = new Date(displayEvent.start)
  const oldStartTime = new Date(oldEvent.startTime as string)
  
  let oldRrule = RRule.fromString(oldEvent.recurrenceRule as string)
  
  oldRrule.origOptions.dtstart = datetime(oldStartTime.getUTCFullYear(), oldStartTime.getUTCMonth() + 1, oldStartTime.getUTCDate(), oldStartTime.getUTCHours(), oldStartTime.getUTCMinutes(), oldStartTime.getUTCSeconds())
  const before = datetime(displayStartTime.getUTCFullYear(), displayStartTime.getUTCMonth() + 1, displayStartTime.getUTCDate(), displayStartTime.getUTCHours(), displayStartTime.getUTCMinutes(), displayStartTime.getUTCSeconds())
  oldRrule = new RRule(oldRrule.origOptions)
  const previousOccurence = oldRrule.before(before)

  if (previousOccurence) {
    oldRrule.origOptions.until = previousOccurence
  } else {
    oldRrule.origOptions.until = before
  }
  
  oldEvent.recurrenceRule = RRule.optionsToString(oldRrule.origOptions).split('\n')[1]
}

function onValidationChange(isValid: boolean) {
  hasValidationErrors.value = !isValid
}
</script>

<style scoped>
</style>
