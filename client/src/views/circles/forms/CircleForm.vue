<template>
  <v-container v-if="circle" max-width="600" class="pa-1">
    <v-row>
      <v-col class="pt-5">
        <div class="text-h4">
          {{ isEditing ? 'Edit Circle' : 'Create New Circle' }}
        </div>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card class="mx-auto" max-width="600">
          <v-card-text>
            <v-form @submit.prevent="handleSave">
              <v-text-field
                v-model="circle.title"
                label="Circle Title"
                required
              ></v-text-field>

              <v-select
                v-model="circle.visibility"
                :items="visibilityOptions"
                item-title="label"
                item-value="value"
                label="Visibility"
                required
              />

              <!-- Visibility Description -->
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
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

  <!-- Close FAB -->
  <v-btn
    color="primary"
    icon="mdi-close"
    style="position: fixed; bottom: 16px; left: 16px"
    @click="$emit('close')"
  ></v-btn>

  <!-- Save FAB -->
  <v-btn
    color="primary"
    icon="mdi-content-save"
    style="position: fixed; bottom: 16px; right: 16px"
    @click="handleSave"
  ></v-btn>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Circle, apitypes_VisibilityLevel } from '@/genapi/api/circles/circle/v1alpha1'

const props = defineProps<{
  modelValue: Circle
  isEditing?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Circle): void
  (e: 'save'): void
  (e: 'close'): void
}>()

const circle = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const visibilityOptions = [
  {
    label: 'Public',
    value: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
    icon: 'mdi-earth',
    color: 'success',
    description: 'Everyone can see this circle'
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
    description: 'Only you can see this circle'
  }
]

// Computed properties for the selected visibility
const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === circle.value?.visibility)
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