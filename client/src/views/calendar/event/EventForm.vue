<template>
  <v-form>
    <v-select
      :model-value="modelValue.calendarName"
      :items="calendars"
      item-title="title"
      item-value="name"
      label="Calendar"
      density="compact"
      hide-details
      class="mb-4"
      :disabled="disabled"
      @update:model-value="updateField('calendarName', $event)"
    />
    <v-text-field
      :model-value="modelValue.title"
      label="Title"
      density="compact"
      class="mb-4"
      :disabled="disabled"
      @update:model-value="updateField('title', $event)"
      />
      <v-text-field
      :model-value="modelValue.start"
      label="Start"
      type="datetime-local"
      :max="modelValue.end"
      :error-messages="startTimeValidationError"
      :error="startTimeValidationError !== ''"
      persistent-hint
      density="compact"
      class="mb-4"
      :disabled="disabled"
      @update:model-value="updateField('start', $event)"
    />
    <v-text-field
      :model-value="modelValue.end"
      label="End"
      type="datetime-local"
      :min="modelValue.start"
      :error-messages="endTimeValidationError"
      :error="endTimeValidationError !== ''"
      persistent-hint
      density="compact"
      class="mb-4"
      :disabled="disabled"
      @update:model-value="updateField('end', $event)"
    />
    
    <v-textarea
      :model-value="modelValue.description"
      label="Description"
      rows="3"
      class="mb-4"
      hide-details
      density="compact"
      placeholder="Optional event description..."
      :disabled="disabled"
      @update:model-value="updateField('description', $event)"
    />
    
    <!-- Recurrence Rule Section -->
    <div >
      <RecurrenceRuleForm
      :model-value="modelValue.recurrenceRule"
      :disabled="disabled"
      @update:model-value="updateField('recurrenceRule', $event)"
      />
    </div>
  </v-form>
</template>

<script setup lang="ts">
import { computed, watch, ref } from 'vue'
import type { Calendar } from '@/genapi/api/calendars/calendar/v1alpha1'
import { RecurrenceRuleForm } from '@/components/calendar'

export interface EventFormData {
  calendarName: string | null
  title: string
  start: string
  end: string
  description: string
  recurrenceRule: string
}

const props = withDefaults(defineProps<{
  modelValue: EventFormData
  calendars: Calendar[]
  disabled?: boolean
}>(), {
  calendars: () => [],
  disabled: false,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: EventFormData): void
  (e: 'validationChange', hasErrors: boolean): void
}>()

// Track if end time was automatically adjusted
const endTimeAutoAdjusted = ref(false)

// Computed property to validate start and end times
const startTimeValidationError = computed(() => {
  // If start time is empty, no error yet
  if (!props.modelValue.start) {
    return ''
  }
  
  // If end time is empty but start time exists, no error yet
  if (!props.modelValue.end) {
    return ''
  }
  
  const startTime = new Date(props.modelValue.start)
  const endTime = new Date(props.modelValue.end)
  
  // Check if dates are valid
  if (isNaN(startTime.getTime()) || isNaN(endTime.getTime())) {
    return ''
  }
  
  if (startTime >= endTime) {
    return 'Start time must be before end time'
  }
  
  return ''
})

const endTimeValidationError = computed(() => {
  // If end time is empty, no error yet
  if (!props.modelValue.end) {
    return ''
  }
  
  // If start time is empty but end time exists, no error yet
  if (!props.modelValue.start) {
    return ''
  }
  
  const startTime = new Date(props.modelValue.start)
  const endTime = new Date(props.modelValue.end)
  
  // Check if dates are valid
  if (isNaN(startTime.getTime()) || isNaN(endTime.getTime())) {
    return ''
  }
  
  if (startTime >= endTime) {
    return 'Start time must be before end time'
  }
  
  return ''
})

// Computed property to check if form has validation errors
const hasValidationErrors = computed(() => {
  const hasErrors = startTimeValidationError.value !== '' || endTimeValidationError.value !== ''
  return hasErrors
})

// Watch for validation changes and emit to parent
watch(hasValidationErrors, (hasErrors) => {
  emit('validationChange', hasErrors)
}, { immediate: true })

// Watch for form value changes to debug validation
watch(() => [props.modelValue.start, props.modelValue.end], () => {
  // Form values changed - validation will automatically update
}, { immediate: true })

// Function to update individual fields without causing recursive updates
function updateField(field: keyof EventFormData, value: EventFormData[keyof EventFormData]) {
  const updatedData = { ...props.modelValue, [field]: value }
  
  // If start time is being updated and we have both start and end times,
  // automatically adjust the end time to maintain the same duration
  if (field === 'start' && props.modelValue.start && props.modelValue.end) {
    const oldStartTime = new Date(props.modelValue.start)
    const oldEndTime = new Date(props.modelValue.end)
    
    // Check if both times are valid
    if (!isNaN(oldStartTime.getTime()) && !isNaN(oldEndTime.getTime())) {
      const newStartTime = new Date(value as string)
      
      // Check if new start time is valid
      if (!isNaN(newStartTime.getTime())) {
        // Calculate the duration between old start and end
        const duration = oldEndTime.getTime() - oldStartTime.getTime()
        
        // Set new end time to maintain the same duration
        const newEndTime = new Date(newStartTime.getTime() + duration)
        
        // Update the end time in the form data
        updatedData.end = toLocalInput(newEndTime.toISOString())
        
        // Mark that end time was automatically adjusted
        endTimeAutoAdjusted.value = true
      }
    }
  } else if (field === 'end') {
    // Reset the auto-adjusted flag when end time is manually changed
    endTimeAutoAdjusted.value = false
  }
  
  emit('update:modelValue', updatedData)
}

// Helper function to convert timestamp to local datetime-local format
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
</script>

<style scoped>
</style>
