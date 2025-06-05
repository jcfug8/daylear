<template>
  <v-container max-width="600" class="pa-1">
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
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'
import { recipeService } from '@/api/api'

const props = defineProps<{
  modelValue: Recipe
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Recipe): void
  (e: 'imageSelected', file: File | null, url: string | null): void
}>()

const recipe = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
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
