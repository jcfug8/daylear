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
      hide-details
      density="compact"
      class="mb-4"
      :disabled="disabled"
      @update:model-value="updateField('title', $event)"
    />
    <v-text-field
      :model-value="modelValue.start"
      label="Start"
      type="datetime-local"
      hide-details
      density="compact"
      class="mb-4"
      :disabled="disabled"
      @update:model-value="updateField('start', $event)"
    />
    <v-text-field
      :model-value="modelValue.end"
      label="End"
      type="datetime-local"
      hide-details
      density="compact"
      class="mb-4"
      :disabled="disabled"
      @update:model-value="updateField('end', $event)"
    />
    
    <!-- Recurrence Rule Section -->
    <div class="mb-4">
      <RecurrenceRuleForm
        :model-value="modelValue.recurrenceRule"
        :disabled="disabled"
        @update:model-value="updateField('recurrenceRule', $event)"
      />
    </div>
    
    <v-textarea
      :model-value="modelValue.description"
      label="Description"
      rows="3"
      hide-details
      density="compact"
      placeholder="Optional event description..."
      :disabled="disabled"
      @update:model-value="updateField('description', $event)"
    />
  </v-form>
</template>

<script setup lang="ts">
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
}>()

// Function to update individual fields without causing recursive updates
function updateField(field: keyof EventFormData, value: EventFormData[keyof EventFormData]) {
  const updatedData = { ...props.modelValue, [field]: value }
  emit('update:modelValue', updatedData)
}
</script>

<style scoped>
</style>
