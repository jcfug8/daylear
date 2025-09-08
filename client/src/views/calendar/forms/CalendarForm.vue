<template>
  <v-container v-if="calendar" max-width="600" class="pa-1 pb-16">
    <v-row>
      <v-col class="pt-5">
        <div class="text-h4">
          {{ isEditing ? 'Edit Calendar' : 'Create New Calendar' }}
        </div>
      </v-col>
    </v-row>

    <v-row>
      <v-col class="pt-5">
        <div class="text-body-1">
          <v-text-field
            density="compact"
            v-model="calendar.title"
            placeholder="Calendar Title"
            required
          ></v-text-field>
        </div>
        <div class="text-body-1">
          <v-textarea
            density="compact"
            v-model="calendar.description"
            placeholder="Describe this calendar..."
            rows="3"
            auto-grow
          ></v-textarea>
        </div>
      </v-col>
    </v-row>

    <!-- Visibility Section -->
    <v-row>
      <v-col cols="12">
        <div class="mt-4">
          <v-select
            v-model="calendar.visibility"
            :items="visibilityOptions"
            item-title="label"
            item-value="value"
            label="Calendar Visibility"
            density="compact"
            variant="outlined"
            required
          >
            <template #selection="{ item }">
              <div class="d-flex align-center">
                <v-icon :icon="item.raw.icon" class="me-2" size="small"></v-icon>
                {{ item.raw.label }}
              </div>
            </template>
            <template #item="{ props, item }">
              <v-list-item v-bind="props">
                <template #prepend>
                  <v-icon :icon="item.raw.icon" size="small"></v-icon>
                </template>
                <v-list-item-subtitle class="text-wrap">
                  {{ item.raw.description }}
                </v-list-item-subtitle>
              </v-list-item>
            </template>
          </v-select>
          
          <!-- Current selection description -->
          <div v-if="selectedVisibilityDescription" class="mt-2">
            <v-alert
              :icon="selectedVisibilityIcon"
              density="compact"
              variant="tonal"
              :color="selectedVisibilityColor"
            >
              <div class="text-body-2">
                <strong>{{ selectedVisibilityLabel }}:</strong> {{ selectedVisibilityDescription }}
              </div>
            </v-alert>
          </div>
        </div>
      </v-col>
    </v-row>
  </v-container>

  <!-- Close FAB -->
  <v-btn
    color="error"
    density="compact"
    style="position: fixed; bottom: 56px; left: 16px"
    @click="$emit('close')"
  >
    <v-icon>mdi-close</v-icon>
    Cancel
  </v-btn>

  <!-- Save FAB -->
  <v-btn
    color="success"
    density="compact"
    style="position: fixed; bottom: 56px; right: 16px"
    @click="handleSave"
  >
    <v-icon>mdi-content-save</v-icon>
    Save
  </v-btn>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Calendar, apitypes_VisibilityLevel } from '@/genapi/api/calendars/calendar/v1alpha1'

const props = defineProps<{
  modelValue: Calendar
  isEditing?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Calendar): void
  (e: 'save'): void
  (e: 'close'): void
}>()

const calendar = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const visibilityOptions = [
  {
    label: 'Public',
    value: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
    icon: 'mdi-earth',
    color: 'success',
    description: 'Everyone can see this calendar'
  },
  {
    label: 'Restricted',
    value: 'VISIBILITY_LEVEL_RESTRICTED' as apitypes_VisibilityLevel,
    icon: 'mdi-account-group',
    color: 'warning',
    description: 'Shared users and their connections can see this'
  },
  {
    label: 'Private',
    value: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
    icon: 'mdi-lock',
    color: 'info',
    description: 'Only specifically shared users can see this'
  },
  {
    label: 'Hidden',
    value: 'VISIBILITY_LEVEL_HIDDEN' as apitypes_VisibilityLevel,
    icon: 'mdi-eye-off',
    color: 'secondary',
    description: 'Only you can see this calendar'
  }
]

// Computed properties for the selected visibility
const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === calendar.value?.visibility)
})

const selectedVisibilityDescription = computed(() => {
  return selectedVisibility.value?.description || ''
})

const selectedVisibilityLabel = computed(() => {
  return selectedVisibility.value?.label || ''
})

const selectedVisibilityIcon = computed(() => {
  return selectedVisibility.value?.icon || 'mdi-help-circle'
})

const selectedVisibilityColor = computed(() => {
  return selectedVisibility.value?.color || 'primary'
})

function handleSave() {
  emit('save')
}
</script>

<style scoped>
</style>
