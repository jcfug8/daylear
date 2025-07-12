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

              <v-text-field
                v-model="circle.handle"
                label="Handle (unique, for sharing)"
                hint="Optional. If left blank, one will be generated. Must be unique."
                persistent-hint
                :rules="[handleRule]"
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

    <v-row>
      <v-spacer></v-spacer>
      <v-col align-self="auto" cols="12" sm="8">
        <div class="image-container">
          <v-img
            class="mt-1"
            style="background-color: lightgray"
            :src="imageCleared ? '' : (previewImage || circle.imageUri)"
            cover
            height="300"
          ></v-img>
          <v-btn
            icon="mdi-camera"
            color="primary"
            class="image-upload-btn"
            @click="showImageDialog = true"
          ></v-btn>
          <v-btn
            v-if="!imageCleared && (previewImage || circle.imageUri)"
            icon="mdi-close"
            color="error"
            class="image-x-btn"
            style="position: absolute; top: 8px; right: 56px; z-index: 2;"
            @click="clearImage"
            title="Remove Image"
          ></v-btn>
          <v-btn
            v-if="imageCleared"
            icon="mdi-arrow-u-left-top"
            color="info"
            class="image-undo-btn"
            style="position: absolute; top: 8px; right: 56px; z-index: 2;"
            @click="undoClearImage"
            title="Undo Remove Image"
          ></v-btn>
        </div>
      </v-col>
      <v-spacer></v-spacer>
    </v-row>
  </v-container>

  <!-- Close FAB -->
  <v-btn
    color="error"
    icon="mdi-close"
    style="position: fixed; bottom: 16px; left: 16px"
    @click="$emit('close')"
  ></v-btn>

  <!-- Save FAB -->
  <v-btn
    color="success"
    icon="mdi-content-save"
    style="position: fixed; bottom: 16px; right: 16px"
    @click="handleSave"
  ></v-btn>

  <v-dialog v-model="showImageDialog" max-width="500">
    <v-card>
      <v-card-title>Add Circle Image</v-card-title>
      <v-card-text>
        <v-tabs v-model="imageTab">
          <v-tab value="url">Image URL</v-tab>
          <v-tab value="upload">Upload Image</v-tab>
        </v-tabs>
        <v-window v-model="imageTab">
          <v-window-item value="url">
            <v-text-field
              v-model="imageUrl"
              label="Image URL"
              placeholder="Enter image URL"
              class="mt-4"
              @update:model-value="updatePreview"
            ></v-text-field>
            <v-img
              v-if="previewImage"
              :src="previewImage"
              max-height="200"
              contain
              class="mt-4"
            ></v-img>
          </v-window-item>
          <v-window-item value="upload">
            <v-file-input
              v-model="imageFile"
              label="Choose Image"
              accept="image/jpeg, image/jpg, image/png, image/gif, image/webp, image/bmp, image/svg, image/heic, image/heif"
              class="mt-4"
              prepend-icon="mdi-camera"
              @update:model-value="handleFileSelect"
            ></v-file-input>
            <v-img
              v-if="previewImage"
              :src="previewImage"
              max-height="200"
              contain
              class="mt-4"
            ></v-img>
          </v-window-item>
        </v-window>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="error" @click="cancelImageDialog">Cancel</v-btn>
        <v-btn color="primary" @click="handleImageSubmit">OK</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { Circle, apitypes_VisibilityLevel } from '@/genapi/api/circles/circle/v1alpha1'

const props = defineProps<{
  modelValue: Circle
  isEditing?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Circle): void
  (e: 'save'): void
  (e: 'close'): void
  (e: 'imageSelected', file: File | null, url: string | null): void
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

const showImageDialog = ref(false)
const imageTab = ref('url')
const imageUrl = ref('')
const imageFile = ref<File | null>(null)
const previewImage = ref<string | null>(null)
const imageCleared = ref(false)
const originalImageUri = ref<string | null>(null)

function updatePreview() {
  if (imageUrl.value) {
    previewImage.value = imageUrl.value
  } else {
    previewImage.value = null
  }
}

function handleFileSelect(files: File | File[] | null) {
  if (files instanceof File) {
    const reader = new FileReader()
    reader.onload = (e) => {
      previewImage.value = e.target?.result as string
    }
    reader.readAsDataURL(files)
  } else {
    previewImage.value = null
  }
}

function cancelImageDialog() {
  showImageDialog.value = false
  imageUrl.value = ''
  imageFile.value = null
  previewImage.value = null
}

function handleImageSubmit() {
  if (imageTab.value === 'url') {
    if (imageUrl.value) {
      circle.value.imageUri = imageUrl.value
      emit('imageSelected', null, imageUrl.value)
      imageCleared.value = false
    }
  } else if (imageFile.value) {
    emit('imageSelected', imageFile.value, null)
    imageCleared.value = false
  }
  showImageDialog.value = false
  imageUrl.value = ''
  imageFile.value = null
}

function clearImage() {
  if (!imageCleared.value) {
    originalImageUri.value = previewImage.value || circle.value.imageUri || null
    imageCleared.value = true
    previewImage.value = null
    circle.value.imageUri = ''
    emit('update:modelValue', circle.value)
  }
}

function undoClearImage() {
  if (imageCleared.value && originalImageUri.value) {
    circle.value.imageUri = originalImageUri.value
    previewImage.value = originalImageUri.value
    imageCleared.value = false
    emit('update:modelValue', circle.value)
  }
}

// Reset clear state if a new image is selected
watch([previewImage, () => circle.value.imageUri], (vals) => {
  const [newPreview, newUri] = vals
  if ((newPreview || newUri) && imageCleared.value) {
    imageCleared.value = false
  }
})

function handleSave() {
  emit('save')
}

// Simple rule: allow only letters, numbers, dashes, underscores
const handleRule = (v: string) => {
  if (!v) return true
  if (!/^[a-zA-Z0-9_-]+$/.test(v)) return 'Handle can only contain letters, numbers, dashes, and underscores.'
  return true
}
</script>

<style scoped>
.image-container {
  position: relative;
}

.image-upload-btn {
  position: absolute;
  bottom: 16px;
  right: 16px;
}

.image-x-btn, .image-undo-btn {
  position: absolute;
  top: 8px;
  right: 56px;
  z-index: 2;
}
</style> 