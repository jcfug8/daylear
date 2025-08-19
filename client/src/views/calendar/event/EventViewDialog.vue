<template>
  <v-dialog 
    :key="`view-${event?.name || 'unknown'}`"
    v-model="internalOpen" 
    max-width="600" 
    @update:modelValue="onDialogModel"
  >
    <v-card v-if="event && calendar">
      <v-card-title class="d-flex align-center">
        <span>{{ event.title || 'Event Details' }}</span>
        <v-spacer />
        <v-btn icon="mdi-close" variant="text" @click="close()" />
      </v-card-title>
      
      <v-card-text>
        <!-- Calendar Info -->
        <div class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">Calendar</div>
          <div class="text-body-1">{{ calendarTitle }}</div>
        </div>

        <!-- Time Info -->
        <div class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">Time</div>
          <div class="text-body-1">
            <div>{{ formatDateTime(eventOccurrence?.start) }}</div>
            <div v-if="eventOccurrence?.end && eventOccurrence?.end !== eventOccurrence?.start">
              to {{ formatDateTime(eventOccurrence?.end) }}
            </div>
          </div>
        </div>

        <!-- Description -->
        <div v-if="event.description" class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">Description</div>
          <div class="text-body-1">{{ event.description }}</div>
        </div>

        <!-- geo -->
        <div v-if="event.geo" class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">geo</div>
          <div class="text-body-1">
            <div v-if="event.geo.latitude && event.geo.longitude">
              {{ event.geo.latitude.toFixed(6) }}, {{ event.geo.longitude.toFixed(6) }}
            </div>
          </div>
        </div>

        <!-- location -->
        <div v-if="event.location" class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">Location</div>
          <div class="text-body-1">{{ event.location }}</div>
        </div>

        <!-- URL -->
        <div v-if="event.uri" class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">URL</div>
          <div class="text-body-1">
            <a :href="event.uri" target="_blank" class="text-decoration-none">
              {{ event.uri }}
            </a>
          </div>
        </div>

        <!-- Recurrence -->
        <div v-if="event.recurrenceRule" class="mb-4">
          <RecurrenceRuleDisplay 
            :recurrence-rule="event.recurrenceRule"
            :start-date="event.startTime as unknown as string"
          />
        </div>

        <!-- Alarms -->
        <div v-if="event.alarms && event.alarms.length > 0" class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">Alarms</div>
          <div v-for="alarm in event.alarms" :key="alarm.alarmId" class="text-body-1 mb-1">
            <div v-if="alarm.trigger?.dateTime">
              {{ formatDateTime(alarm.trigger.dateTime) }}
            </div>
            <div v-else-if="alarm.trigger?.duration">
              {{ alarm.trigger.duration }} before event
            </div>
          </div>
        </div>

        <!-- Additional Info -->
        <div v-if="event.excludedTimes && event.excludedTimes.length > 0" class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">Excluded Dates</div>
          <div class="text-body-1">
            <div v-for="time in event.excludedTimes" :key="time" class="mb-1">
              {{ formatDateTime(time) }}
            </div>
          </div>
        </div>

        <div v-if="event.additionalTimes && event.additionalTimes.length > 0" class="mb-4">
          <div class="text-caption text-medium-emphasis mb-1">Additional Dates</div>
          <div class="text-body-1">
            <div v-for="time in event.additionalTimes" :key="time" class="mb-1">
              {{ formatDateTime(time) }}
            </div>
          </div>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="close()">Close</v-btn>
        <v-btn 
          v-if="canEdit" 
          color="primary" 
          @click="editEvent"
        >
          Edit
        </v-btn>
      </v-card-actions>
    </v-card>
    
    <!-- Loading state when event or calendar is not available -->
    <v-card v-else>
      <v-card-title>Loading...</v-card-title>
      <v-card-text class="text-center py-8">
        <v-progress-circular indeterminate />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="close()">Close</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  
  <!-- Event Edit Dialog - Only render when needed -->
  <EventEditDialog
    v-if="showEditDialog && event && calendar"
    :key="`edit-${event?.name || 'unknown'}`"
    v-model="showEditDialog"
    :display-event="eventOccurrence"
    :event="event"
    :calendars="[calendar]"
    @updated="handleEventUpdated"
    @deleted="handleEventDeleted"
  />
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import type { Calendar, Event } from '@/genapi/api/calendars/calendar/v1alpha1'
import EventEditDialog from './EventEditDialog.vue'
import { RecurrenceRuleDisplay } from '@/components/calendar'
import type { CalendarEventExternal } from '@schedule-x/calendar'

const props = withDefaults(defineProps<{
  modelValue: boolean
  event: Event | null
  eventOccurrence: CalendarEventExternal | null
  calendar: Calendar | null
}>(), {
  modelValue: false,
  event: null,
  calendar: null,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'edit', event: Event): void
  (e: 'updated', event: Event): void
  (e: 'deleted', event: Event): void
}>()

const internalOpen = ref<boolean>(props.modelValue)
const showEditDialog = ref<boolean>(false)

// Watch for prop changes and update internal state
watch(() => props.modelValue, (newValue) => {
  internalOpen.value = newValue
  // Reset edit dialog when main dialog closes
  if (!newValue) {
    showEditDialog.value = false
  }
})

// Watch internal state and emit changes
watch(internalOpen, (newValue) => {
  emit('update:modelValue', newValue)
  // Reset edit dialog when main dialog closes
  if (!newValue) {
    showEditDialog.value = false
  }
})

// Calendar title for display
const calendarTitle = computed(() => {
  if (!props.calendar) return 'Unknown Calendar'
  return props.calendar.title || 'Untitled Calendar'
})

// Check if user can edit this event
const canEdit = computed(() => {
  if (!props.calendar?.calendarAccess?.permissionLevel) return false
  
  return props.calendar.calendarAccess.permissionLevel === 'PERMISSION_LEVEL_WRITE' || 
         props.calendar.calendarAccess.permissionLevel === 'PERMISSION_LEVEL_ADMIN'
})

// Format timestamp for display
function formatDateTime(timestamp: string | undefined): string {
  if (!timestamp) return 'Not specified'
  
  try {
    const date = new Date(timestamp)
    if (isNaN(date.getTime())) return 'Invalid date'
    
    return date.toLocaleString(undefined, {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZoneName: 'short'
    })
  } catch {
    return 'Invalid date'
  }
}

function close() {
  internalOpen.value = false
}

function onDialogModel(value: boolean) {
  internalOpen.value = value
}

function editEvent() {
  if (props.event && props.calendar) {
    console.log('Opening edit dialog for event:', props.event)
    console.log('Calendar:', props.calendar)
    showEditDialog.value = true
  } else {
    console.warn('Cannot open edit dialog - missing event or calendar:', {
      event: props.event,
      calendar: props.calendar
    })
  }
}

function handleEventUpdated(event: Event) {
  emit('updated', event)
  showEditDialog.value = false
}

function handleEventDeleted(event: Event) {
  emit('deleted', event)
  showEditDialog.value = false
}
</script>

<style scoped>
.text-decoration-none {
  text-decoration: none;
}
</style>
