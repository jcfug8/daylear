<template>
  <v-container v-if="recipe">
    <v-app-bar>
      <v-tabs style="width: 100%" v-model="tab" center-active show-arrows fixed-tabs>
        <v-tab value="general" text="General"></v-tab>
        <v-tab value="ingredients" text="Ingredients"></v-tab>
        <v-tab value="directions" text="Directions"></v-tab>
      </v-tabs>
      <template #append>
        <v-btn
          icon="mdi-pencil"
          @click="router.push({ name: 'recipeEdit', params: { recipeId: recipe.name } })"
        ></v-btn>
      </template>
    </v-app-bar>
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="general">
        <v-container max-width="600" class="pa-1">
          <v-row>
            <v-col class="pt-5">
              <div class="text-h4">
                {{ recipe.title }}
              </div>
              <div class="text-body-1">
                {{ recipe.description }}
              </div>
            </v-col>
          </v-row>
          <v-row>
            <v-spacer></v-spacer>
            <v-col align-self="auto" cols="12" sm="8">
              <v-img
                class="mt-1"
                style="background-color: lightgray"
                :src="recipe.imageUri"
                cover
              ></v-img>
            </v-col>
            <v-spacer></v-spacer>
          </v-row>
        </v-container>
      </v-tabs-window-item>
      <v-tabs-window-item value="ingredients">
        <v-container max-width="600">
          <v-card
            class="my-4 mx-1 pa-2"
            v-for="(ingredientGroup, i) in recipe.ingredientGroups"
            :key="i"
          >
            <v-card-title v-if="ingredientGroup.title">{{ ingredientGroup.title }}</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item
                  slim
                  prepend-icon="mdi-circle-small"
                  v-for="(ingredient, j) in ingredientGroup.ingredients"
                  :key="j"
                >
                  {{ ingredient.measurementAmount }}
                  {{ convertMeasurementTypeToString(ingredient.measurementType) }}
                  {{ ingredient.title }} <span v-if="ingredient.optional">(optional)</span>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-container>
      </v-tabs-window-item>
      <v-tabs-window-item value="directions">
        <v-container max-width="600">
          <v-card class="my-4 mx-1 pa-2" v-for="(direction, i) in recipe.directions" :key="i">
            <v-card-title v-if="direction.title">{{ direction.title }}</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item
                  slim
                  prepend-icon="mdi-circle-small"
                  v-for="(step, n) in direction.steps"
                >
                  <div class="font-weight-bold">Step {{ n + 1 }}</div>
                  {{ step }}
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-container>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-container>
</template>

<script setup lang="ts">
import type { Recipe_MeasurementType } from '@/genapi/api/meals/recipe/v1alpha1'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRecipesStore } from '@/stores/recipes'
import { useRecipeFormStore } from '@/stores/recipeForm'
import { storeToRefs } from 'pinia'
import { onMounted, onBeforeUnmount, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const recipesStore = useRecipesStore()
const recipeFormStore = useRecipeFormStore()
const { recipe } = storeToRefs(recipesStore)
const { activeTab } = storeToRefs(recipeFormStore)

// Use the store's activeTab as our local tab
const tab = computed({
  get: () => activeTab.value,
  set: (value) => recipeFormStore.setActiveTab(value),
})

const breadcrumbStore = useBreadcrumbStore()

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

// Map MeasurmentType to string
const measurementTypeToString: Record<Recipe_MeasurementType, string> = {
  MEASUREMENT_TYPE_UNSPECIFIED: '',
  MEASUREMENT_TYPE_TABLESPOON: 'tablespoons',
  MEASUREMENT_TYPE_TEASPOON: 'teaspoons',
  MEASUREMENT_TYPE_OUNCE: 'ounces',
  MEASUREMENT_TYPE_POUND: 'pounds',
  MEASUREMENT_TYPE_GRAM: 'grams',
  MEASUREMENT_TYPE_MILLILITER: 'milliliters',
  MEASUREMENT_TYPE_LITER: 'liters',
}

function convertMeasurementTypeToString(type: Recipe_MeasurementType | undefined): string {
  return type ? measurementTypeToString[type] || '' : ''
}

onMounted(async () => {
  // First check URL hash
  const currentHash = route.hash
  if (currentHash in hashToTab) {
    tab.value = hashToTab[currentHash]
  }

  // Load the recipe based on the route parameter
  const recipeId = route.params.recipeId as string
  await recipesStore.loadRecipe(recipeId)

  breadcrumbStore.setBreadcrumbs([
    { title: 'Recipes', to: { name: 'recipes' } },
    {
      title: recipe.value?.title || '',
      to: { name: 'recipe', params: { recipeId: recipe.value?.name } },
    },
  ])
})

// Reset tab state when leaving the view
onBeforeUnmount(() => {
  // Only reset if we're not going to the edit view
  if (router.currentRoute.value.name !== 'recipeEdit') {
    recipeFormStore.setActiveTab('general')
  }
})

// Watch for tab changes and update the URL hash
watch(tab, (newTab) => {
  const newHash = tabToHash[newTab]
  if (newHash && route.hash !== newHash) {
    router.replace({ hash: newHash }) // Update the URL hash without reloading the page
  }
})
</script>
