<template>
  <v-container max-width="600" class="pa-1">
    <!-- Scrape Recipe by URL (only in create mode) -->
    <template v-if="isCreate">
      <v-row>
        <v-col cols="9">
          <v-text-field
            v-model="scrapeUrl"
            label="Import Recipe from URL"
            placeholder="Paste recipe URL (e.g. https://...)"
            density="compact"
            :disabled="scrapeLoading"
          />
        </v-col>
        <v-col cols="3" class="d-flex align-end">
          <v-btn color="primary" :loading="scrapeLoading" @click="scrapeRecipe" :disabled="!scrapeUrl || scrapeLoading">
            Scrape
          </v-btn>
        </v-col>
      </v-row>
      <v-row v-if="scrapeError">
        <v-col cols="12">
          <v-alert type="error" density="compact">{{ scrapeError }}</v-alert>
        </v-col>
      </v-row>
      <!-- OCR Recipe from Image -->
      <v-row class="mt-2">
        <v-col cols="9">
          <v-file-input
            multiple
            v-model="ocrImageFiles"
            label="Import Recipe from Image (OCR)"
            accept="image/jpeg, image/jpg, image/png, image/gif, image/tiff, image/tif, image/webp, image/bmp, image/svg, image/pdf"
            density="compact"
            :disabled="ocrLoading"
            prepend-icon="mdi-camera"
          />
        </v-col>
        <v-col cols="3" class="d-flex align-end">
          <v-btn color="primary" :loading="ocrLoading" @click="ocrRecipe" :disabled="!ocrImageFiles || ocrLoading">
            OCR
          </v-btn>
        </v-col>
      </v-row>
      <v-row v-if="ocrError">
        <v-col cols="12">
          <v-alert type="error" density="compact">{{ ocrError }}</v-alert>
        </v-col>
      </v-row>
    </template>
    <v-row>
      <v-col class="pt-5">
        <div class="text-h4">
          <v-text-field
            density="compact"
            v-model="recipe.title"
            placeholder="Recipe Title"
          ></v-text-field>
        </div>
        <div class="text-body-1">
          <v-textarea
            density="compact"
            v-model="recipe.description"
            placeholder="Recipe Description"
          ></v-textarea>
        </div>
      </v-col>
    </v-row>
    <v-row>
      <v-spacer></v-spacer>
      <v-col align-self="auto" cols="12" sm="8">
        <div class="image-container">
          <v-img
            class="mt-1"
            style="background-color: lightgray"
            :src="previewImage || recipe.imageUri"
            cover
            height="300"
          ></v-img>
          <v-btn
            icon="mdi-camera"
            color="primary"
            class="image-upload-btn"
            @click="showImageDialog = true"
          ></v-btn>
        </div>
      </v-col>
      <v-spacer></v-spacer>
    </v-row>

    <!-- Visibility Section -->
    <v-row>
      <v-col cols="12">
        <div class="mt-4">
          <v-select
            v-model="recipe.visibility"
            :items="visibilityOptions"
            item-title="label"
            item-value="value"
            label="Recipe Visibility"
            density="compact"
            variant="outlined"
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

    <!-- Image Upload Dialog -->
    <v-dialog v-model="showImageDialog" max-width="500">
      <v-card>
        <v-card-title>Add Recipe Image</v-card-title>
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
                accept="image/*"
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
  </v-container>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Recipe, apitypes_VisibilityLevel } from '@/genapi/api/meals/recipe/v1alpha1'
import { recipeService, fileService } from '@/api/api'

const props = defineProps<{
  modelValue: Recipe,
  isCreate?: boolean
}>()

const isCreate = computed(() => props.isCreate ?? false)

const emit = defineEmits<{
  (e: 'update:modelValue', value: Recipe): void
  (e: 'imageSelected', file: File | null, url: string | null): void
}>()

const recipe = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

// Scrape logic
const scrapeUrl = ref('')
const scrapeLoading = ref(false)
const scrapeError = ref('')

async function scrapeRecipe() {
  scrapeError.value = ''
  scrapeLoading.value = true
  try {
    const resp = await recipeService.ScrapeRecipe({ uri: scrapeUrl.value })
    if (resp && resp.recipe) {
      emit('update:modelValue', resp.recipe)
    } else {
      scrapeError.value = 'No recipe found at that URL.'
    }
  } catch (err: any) {
    scrapeError.value = err?.message || 'Failed to scrape recipe.'
  } finally {
    scrapeLoading.value = false
  }
}

// Visibility options with descriptions and icons
const visibilityOptions = [
  {
    value: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
    label: 'Public',
    icon: 'mdi-earth',
    color: 'success',
    description: 'Everyone can see this recipe'
  },
  {
    value: 'VISIBILITY_LEVEL_RESTRICTED' as apitypes_VisibilityLevel,
    label: 'Restricted',
    icon: 'mdi-account-group',
    color: 'warning',
    description: 'Shared users, circles and their connections can see this'
  },
  {
    value: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
    label: 'Private',
    icon: 'mdi-lock',
    color: 'info',
    description: 'Only specifically shared users and circles can see this'
  },
  {
    value: 'VISIBILITY_LEVEL_HIDDEN' as apitypes_VisibilityLevel,
    label: 'Hidden',
    icon: 'mdi-eye-off',
    color: 'secondary',
    description: 'Only you can see this recipe'
  }
]

// Computed properties for the selected visibility
const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === recipe.value.visibility)
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
      recipe.value.imageUri = imageUrl.value
      emit('imageSelected', null, imageUrl.value)
    }
  } else if (imageFile.value) {
    emit('imageSelected', imageFile.value, null)
  }
  
  showImageDialog.value = false
  imageUrl.value = ''
  imageFile.value = null
}

// OCR logic
const ocrImageFiles = ref<File[] | null>(null)
const ocrLoading = ref(false)
const ocrError = ref('')

async function ocrRecipe() {
  ocrError.value = ''
  ocrLoading.value = true
  try {
    if (!ocrImageFiles.value) throw new Error('No image selected')
    const resp = await fileService.OCRRecipe({ files: ocrImageFiles.value })
    if (resp && resp.recipe) {
      emit('update:modelValue', resp.recipe)
    } else {
      ocrError.value = 'No recipe found in image.'
    }
  } catch (err: any) {
    ocrError.value = err?.message || 'Failed to OCR recipe.'
  } finally {
    ocrLoading.value = false
  }
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
</style>
