<template>
  <v-form>
    <v-select
      v-model="formData.calendarName"
      :items="calendars"
      item-title="title"
      item-value="name"
      label="Calendar"
      density="compact"
      hide-details
      class="mb-4"
      :disabled="disabled"
    />
    <v-text-field
      v-model="formData.title"
      label="Title"
      hide-details
      density="compact"
      class="mb-4"
      :disabled="disabled"
    />
    <v-text-field
      v-model="formData.start"
      label="Start"
      type="datetime-local"
      hide-details
      density="compact"
      class="mb-4"
      :disabled="disabled"
    />
    <v-text-field
      v-model="formData.end"
      label="End"
      type="datetime-local"
      hide-details
      density="compact"
      class="mb-4"
      :disabled="disabled"
    />
    <v-textarea
      v-model="formData.description"
      label="Description"
      rows="3"
      hide-details
      density="compact"
      placeholder="Optional event description..."
      :disabled="disabled"
    />
  </v-form>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Calendar } from '@/genapi/api/calendars/calendar/v1alpha1'

export interface EventFormData {
  calendarName: string | null
  title: string
  start: string
  end: string
  description: string
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

// Simple computed property for two-way binding
const formData = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})
</script>

<style scoped>
</style>
