<template>
  <v-form>
    <v-container v-if="recipe">
      <v-app-bar>
        <v-tabs style="width: 100%" v-model="tab" center-active show-arrows fixed-tabs>
          <v-tab
            v-for="tabItem in tabs"
            :key="tabItem.value"
            :value="tabItem.value"
            :text="tabItem.text"
          ></v-tab>
        </v-tabs>
        <template #prepend>
          <v-btn icon="mdi-close" @click="$emit('close')"></v-btn>
        </template>
        <template #append>
          <v-btn icon="mdi-content-save" @click="$emit('save')"></v-btn>
        </template>
      </v-app-bar>
      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="general">
          <recipe-general-form v-model="recipe" />
        </v-tabs-window-item>
        <v-tabs-window-item value="ingredients">
          <recipe-ingredients-form v-model="recipe" />
        </v-tabs-window-item>
        <v-tabs-window-item value="directions">
          <recipe-directions-form v-model="recipe" />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-container>
  </v-form>
</template>

<script setup lang="ts">
import { computed, onMounted, watch, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'
import RecipeGeneralForm from './RecipeGeneralForm.vue'
import RecipeIngredientsForm from './RecipeIngredientsForm.vue'
import RecipeDirectionsForm from './RecipeDirectionsForm.vue'
import { useRecipeFormStore } from '@/stores/recipeForm'

const props = defineProps<{
  modelValue: Recipe
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Recipe): void
  (e: 'save'): void
  (e: 'close'): void
}>()

const recipe = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const route = useRoute()
const router = useRouter()
const recipeFormStore = useRecipeFormStore()
const { activeTab } = storeToRefs(recipeFormStore)

// Use the store's activeTab as our local tab
const tab = computed({
  get: () => activeTab.value,
  set: (value) => recipeFormStore.setActiveTab(value),
})

const tabs = [
  { value: 'general', text: 'General' },
  { value: 'ingredients', text: 'Ingredients' },
  { value: 'directions', text: 'Directions' },
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
</script>
