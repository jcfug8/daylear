<template>
  <v-container max-width="600" class="pa-1">
    <!-- Centered Import Button -->
    <v-row v-if="isCreate" justify="center">
      <v-col cols="auto">
        <v-btn
          color="primary"
          @click="showScrapeOcrDialog = true"
          prepend-icon="mdi-import"
        >
          Import Recipe
        </v-btn>
      </v-col>
    </v-row>
    <!-- Scrape/OCR Modal with Tabs -->
    <v-dialog v-model="showScrapeOcrDialog" max-width="600">
      <v-card>
        <v-card-title>Import Recipe</v-card-title>
        <v-card-subtitle>
          Importing uses AI and may take up to a minute to complete and often makes mistakes. After the import is finished, you may want to edit the recipe to make it more accurate.
        </v-card-subtitle>
        <v-card-text>
          <v-alert v-if="scrapeLoading || ocrLoading" type="info" density="compact" class="mb-2">
            Import in progress: {{ formatImportTimer() }} elapsed
          </v-alert>
          <v-tabs v-model="importTab">
            <v-tab value="scrape">URL</v-tab>
            <v-tab value="ocr">Image</v-tab>
          </v-tabs>
          <v-window v-model="importTab">
            <v-window-item value="scrape">
              <v-row class="mt-2">
                <v-col cols="12" sm="9">
                  <v-text-field
                    v-model="scrapeUrl"
                    clearable
                    label="Import Recipe from URL"
                    placeholder="Paste recipe URL (e.g. https://...)"
                    density="compact"
                    :disabled="scrapeLoading"
                  />
                </v-col>
                <v-col cols="12" sm="3" class="d-flex align-end">
                  <v-btn color="primary" block :loading="scrapeLoading" @click="scrapeRecipe" :disabled="!scrapeUrl || scrapeLoading">
                    Import
                  </v-btn>
                </v-col>
              </v-row>
            </v-window-item>
            <v-window-item value="ocr">
              <v-row class="mt-2">
                <v-col cols="12" sm="9">
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
                <v-col cols="12" sm="3" class="d-flex align-end">
                  <v-btn color="primary" block :loading="ocrLoading" @click="ocrRecipe" :disabled="!ocrImageFiles || ocrLoading">
                    Import
                  </v-btn>
                </v-col>
              </v-row>
              <v-row v-if="ocrError">
                <v-col cols="12">
                  <v-alert type="error" density="compact">{{ ocrError }}</v-alert>
                </v-col>
              </v-row>
            </v-window-item>
          </v-window>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" variant="text" @click="handleCloseScrapeOcrDialog">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="showPostImportDialog" max-width="400">
      <v-card>
        <v-card-title>Recipe Imported</v-card-title>
        <v-card-text>
          <div class="mb-2">Your recipe has been imported. Would you like to save it now, or review and edit it first?</div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="secondary" @click="showPostImportDialog = false">Review Import</v-btn>
          <v-btn color="success" @click="handleSaveAfterImport">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
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
      <v-col cols="12" sm="4">
        <v-text-field
          v-model.number="prepDurationMinutes"
          label="Prep Duration (minutes)"
          type="number"
          min="0"
          density="compact"
        />
      </v-col>
      <v-col cols="12" sm="4">
        <v-text-field
          v-model.number="cookDurationMinutes"
          label="Cook Duration (minutes)"
          type="number"
          min="0"
          density="compact"
        />
      </v-col>
      <v-col cols="12" sm="4">
        <v-text-field
          v-model.number="totalDurationMinutes"
          label="Total Duration (minutes)"
          type="number"
          min="0"
          density="compact"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="6">
        <v-text-field
          v-model="recipe.cookingMethod"
          label="Cooking Method"
          placeholder="e.g. Frying, Steaming"
          density="compact"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <v-text-field
          v-model="recipe.yieldAmount"
          label="Yield Amount"
          placeholder="e.g. 4 servings, 1 loaf"
          density="compact"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="6">
        <v-combobox
          v-model="categoriesInput"
          :items="[]"
          label="Categories"
          multiple
          chips
          clearable
          density="compact"
          placeholder="Add categories"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <v-combobox
          v-model="cuisinesInput"
          :items="[]"
          label="Cuisines"
          multiple
          chips
          clearable
          density="compact"
          placeholder="Add cuisines"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-spacer></v-spacer>
      <v-col align-self="auto" cols="12" sm="8">
        <div class="image-container">
          <v-img
            class="mt-1"
            style="background-color: lightgray"
            :src="imageCleared ? '' : (previewImage || recipe.imageUri)"
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
            v-if="!imageCleared && (previewImage || recipe.imageUri)"
            icon="mdi-close"
            color="warning"
            class="image-x-btn"
            @click="clearImage"
            title="Remove Image"
          ></v-btn>
          <v-btn
            v-if="imageCleared"
            icon="mdi-arrow-u-left-top"
            color="info"
            class="image-undo-btn"
            @click="undoClearImage"
            title="Undo Remove Image"
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
            v-model="visibilityValue"
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
                clearable
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
  </v-container>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue'
import type { Recipe, apitypes_VisibilityLevel } from '@/genapi/api/meals/recipe/v1alpha1'
import { recipeService, fileService } from '@/api/api'
import { useAlertStore } from '@/stores/alerts'

const props = defineProps<{
  modelValue: Recipe,
  isCreate?: boolean,
  showScrapeOcrDialog?: boolean
}>()

const isCreate = computed(() => props.isCreate ?? false)
const alertStore = useAlertStore()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Recipe): void
  (e: 'imageSelected', file: File | null, url: string | null): void
  (e: 'close-scrape-ocr-dialog'): void
  (e: 'save', pendingImageFile: File | null): void
}>()

const showScrapeOcrDialog = ref(false)
const importTab = ref('scrape')

// Scrape logic
const scrapeUrl = ref('')
const scrapeLoading = ref(false)

const importAbortController = ref<AbortController | null>(null)

async function scrapeRecipe() {
  scrapeLoading.value = true
  importAbortController.value = new AbortController()
  try {
    const resp = await recipeService.ScrapeRecipe({ uri: scrapeUrl.value }, /* importAbortController.value.signal */)
    if (resp && resp.recipe) {
      // Preserve current visibility if not set in scraped recipe
      if (!resp.recipe.visibility) {
        resp.recipe.visibility = recipe.value.visibility
      }
      emit('update:modelValue', resp.recipe)
      showScrapeOcrDialog.value = false
      showPostImportDialog.value = true // <-- show modal
    } else {
      alertStore.addAlert('No recipe found at that URL.', 'error')
    }
  } catch (err: any) {
    if (err?.name !== 'AbortError') {
      alertStore.addAlert(err instanceof Error ? "Unable to scrape recipe\n" + err.message : String(err), 'error')
    }
  } finally {
    scrapeLoading.value = false
    importAbortController.value = null
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

const visibilityValue = computed({
  get() {
    const valid = visibilityOptions.some(opt => opt.value === recipe.value.visibility)
    return valid ? recipe.value.visibility : 'VISIBILITY_LEVEL_PUBLIC'
  },
  set(val: apitypes_VisibilityLevel) {
    recipe.value.visibility = val
    emit('update:modelValue', recipe.value)
  }
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
      recipe.value.imageUri = imageUrl.value
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

// OCR logic
const ocrImageFiles = ref<File[] | null>(null)
const ocrLoading = ref(false)
const ocrError = ref('')

async function ocrRecipe() {
  ocrError.value = ''
  ocrLoading.value = true
  importAbortController.value = new AbortController()
  try {
    if (!ocrImageFiles.value) throw new Error('No image selected')
    const resp = await fileService.OCRRecipe({ files: ocrImageFiles.value }, importAbortController.value.signal)
    if (resp && resp.recipe) {
      // Preserve current visibility if not set in OCR'd recipe
      if (!resp.recipe.visibility) {
        resp.recipe.visibility = recipe.value.visibility
      }
      emit('update:modelValue', resp.recipe)
      showScrapeOcrDialog.value = false
      showPostImportDialog.value = true // <-- show modal
    } else {
      alertStore.addAlert('No recipe found in image.', 'error')
    }
  } catch (err: any) {
    if (err?.name !== 'AbortError') {
      alertStore.addAlert(err instanceof Error ? "Unable to OCR recipe\n" + err.message : String(err), 'error')
    }
  } finally {
    ocrLoading.value = false
    importAbortController.value = null
  }
}

function handleCloseScrapeOcrDialog() {
  showScrapeOcrDialog.value = false
  if (importAbortController.value) {
    importAbortController.value.abort()
    importAbortController.value = null
  }
  emit('close-scrape-ocr-dialog')
}

// Timer for import
const importTimer = ref(0)
let timerInterval: number | null = null

function startImportTimer() {
  if (timerInterval !== null) return
  importTimer.value = 0
  timerInterval = window.setInterval(() => {
    importTimer.value++
  }, 1000)
}
function stopImportTimer() {
  if (timerInterval !== null) {
    clearInterval(timerInterval)
    timerInterval = null
  }
}
function resetImportTimer() {
  importTimer.value = 0
}
function formatImportTimer() {
  const min = Math.floor(importTimer.value / 60)
  const sec = importTimer.value % 60
  return `${min}:${sec.toString().padStart(2, '0')}`
}

watch([scrapeLoading, ocrLoading], ([scrape, ocr]) => {
  if (scrape || ocr) {
    startImportTimer()
  } else {
    stopImportTimer()
    resetImportTimer()
  }
})

const recipe = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

// *** Clear Image ***


function clearImage() {
  if (!imageCleared.value) {
    originalImageUri.value = previewImage.value || recipe.value.imageUri || null
    imageCleared.value = true
    previewImage.value = null
    recipe.value.imageUri = ''
    emit('update:modelValue', recipe.value)
  }
}

function undoClearImage() {
  if (imageCleared.value && originalImageUri.value) {
    recipe.value.imageUri = originalImageUri.value
    previewImage.value = originalImageUri.value
    imageCleared.value = false
    emit('update:modelValue', recipe.value)
  }
}

// Reset clear state if a new image is selected
watch([previewImage, () => recipe.value.imageUri], (vals) => {
  const [newPreview, newUri] = vals
  if ((newPreview || newUri) && imageCleared.value) {
    imageCleared.value = false
  }
})


// Categories and cuisines as tag inputs
const categoriesInput = ref(recipe.value.categories ?? [])
const cuisinesInput = ref(recipe.value.cuisines ?? [])
watch(categoriesInput, (val) => {
  recipe.value.categories = val
})
watch(() => recipe.value.categories, (val) => {
  categoriesInput.value = val ?? []
})
watch(cuisinesInput, (val) => {
  recipe.value.cuisines = val
})
watch(() => recipe.value.cuisines, (val) => {
  cuisinesInput.value = val ?? []
})

// Cook duration in minutes (convert to/from nanoseconds)
const prepDurationMinutes = computed({
  get() {
    return recipe.value.prepDuration ? parseDuration(recipe.value.prepDuration) : 0
  },
  set(val: number) {
    recipe.value.prepDuration = formatDuration(val)
  }
})
const totalDurationMinutes = computed({
  get() {
    return recipe.value.totalDuration ? parseDuration(recipe.value.totalDuration) : 0
  },
  set(val: number) {
    recipe.value.totalDuration = formatDuration(val)
  }
})
const cookDurationMinutes = computed({
  get() {
    return recipe.value.cookDuration ? parseDuration(recipe.value.cookDuration) : 0
  },
  set(val: number) {
    recipe.value.cookDuration = formatDuration(val)
  }
})

onMounted(() => {
  if (!recipe.value.visibility) {
    recipe.value.visibility = 'VISIBILITY_LEVEL_PUBLIC'
    emit('update:modelValue', recipe.value)
  }
})

function formatDuration(duration: number): string {
  if (!duration) return '';
  return String(duration*60) + 's';
}
function parseDuration(duration: string): number {
  if (!duration) return 0;

  if (duration.endsWith('s')) {
    return parseInt(duration.slice(0, -1))/60;
  }
  return 0;
}

const showPostImportDialog = ref(false)

function handleSaveAfterImport() {
  emit('save', null)
  showPostImportDialog.value = false
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
  top: 16px;
  right: 16px;
}
</style>
