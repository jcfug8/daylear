<template>
  <v-form>
    <v-container v-if="recipe" class="pb-16">
      <v-app-bar>
        <v-tabs style="width: 100%" v-model="tab" center-active show-arrows fixed-tabs>
          <v-tab
            v-for="tabItem in tabs"
            :key="tabItem.value"
            :value="tabItem.value"
            :prepend-icon="tabItem.icon"
          ></v-tab>
        </v-tabs>
      </v-app-bar>
      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="general">
          <recipe-general-form 
            v-model="recipe" 
            @image-selected="handleImageSelected"
            :is-create="isCreate"
            :show-scrape-ocr-dialog="props.showImportDialog"
            :import-tab="props.importTab"
            @close-scrape-ocr-dialog="handleCloseImportDialog"
            @save="handleSave"
          />
        </v-tabs-window-item>
        <v-tabs-window-item value="ingredients">
          <recipe-ingredients-form v-model="recipe" />
        </v-tabs-window-item>
        <v-tabs-window-item value="directions">
          <recipe-directions-form v-model="recipe" />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-container>

    <!-- Close FAB -->
    <v-btn
      color="error"
      density="compact"
      style="position: fixed; bottom: 16px; left: 16px"
      @click="$emit('close')"
    >
      <v-icon>mdi-close</v-icon>
      Cancel
    </v-btn>

    <!-- Save FAB -->
    <v-btn
      color="success"
      density="compact"
      style="position: fixed; bottom: 16px; right: 16px"
      @click="handleSave"
    >
      <v-icon>mdi-content-save</v-icon>
      Save
    </v-btn>
  </v-form>
</template>

<script setup lang="ts">
import { computed, onMounted, watch, onBeforeUnmount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'
import RecipeGeneralForm from './RecipeGeneralForm.vue'
import RecipeIngredientsForm from './RecipeIngredientsForm.vue'
import RecipeDirectionsForm from './RecipeDirectionsForm.vue'
import { useRecipeFormStore } from '@/stores/recipeForm'
import { fileService } from '@/api/api'

const props = defineProps<{
  modelValue: Recipe
  showImportDialog?: boolean
  importTab?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Recipe): void
  (e: 'save', pendingImageFile: File | null): void
  (e: 'close'): void
  (e: 'close-import-dialog'): void
}>()

const recipe = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const route = useRoute()
const router = useRouter()
const recipeFormStore = useRecipeFormStore()
const { activeTab } = storeToRefs(recipeFormStore)

// Store the selected image file for later upload
const pendingImageFile = ref<File | null>(null)

// Use the store's activeTab as our local tab
const tab = computed({
  get: () => activeTab.value,
  set: (value) => recipeFormStore.setActiveTab(value),
})

const tabs = [
  { value: 'general', text: 'General', icon: 'mdi-information-outline' },
  { value: 'ingredients', text: 'Ingredients', icon: 'mdi-food-apple-outline' },
  { value: 'directions', text: 'Directions', icon: 'mdi-format-list-numbered' },
]

// Map hash values to tab values
const hashToTab: Record<string, string> = {
  '#general': 'general',
  '#ingredients': 'ingredients',
  '#directions': 'directions',
}

// Map tab values to hash values
const tabToHash: Record<string, string> = {
  general: '#general',
  ingredients: '#ingredients',
  directions: '#directions',
}

function handleImageSelected(file: File | null, url: string | null) {
  pendingImageFile.value = file
}

async function handleSave() {
  // Emit save with the pending image file
  emit('save', pendingImageFile.value)
  pendingImageFile.value = null
}

function handleCloseImportDialog() {
  emit('close-import-dialog')
}

onMounted(() => {
  // First check URL hash
  const currentHash = route.hash
  if (currentHash in hashToTab) {
    tab.value = hashToTab[currentHash]
  }
})

// Reset tab state when leaving the form
onBeforeUnmount(() => {
  // Only reset if we're not going to the recipe view
  if (router.currentRoute.value.name !== 'recipe') {
    recipeFormStore.setActiveTab('general')
  }
})

// Watch for tab changes and update the URL hash
watch(tab, (newTab) => {
  const newHash = tabToHash[newTab]
  if (newHash && route.hash !== newHash) {
    router.replace({ hash: newHash })
  }
})

const isCreate = computed(() => !recipe.value?.name)
</script>
