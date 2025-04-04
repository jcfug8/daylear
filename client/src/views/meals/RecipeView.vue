<template>
  <v-container v-if="recipe">
    <v-row>
      <v-col>
        <v-tabs v-model="tab" fixed-tabs>
          <v-tab value="general" text="General"></v-tab>
          <v-tab value="ingredients" text="Ingredients"></v-tab>
          <v-tab value="directions" text="Directions"></v-tab>
        </v-tabs>
        <v-tabs-window v-model="tab">
          <v-tabs-window-item value="general">
            <v-row>
              <v-col>
                <v-img
                  class="mt-4"
                  style="background-color: lightgray"
                  :src="recipe.imageUri"
                  cover
                ></v-img>
              </v-col>
              <v-col>
                <div>
                  {{ recipe.title }}
                </div>
                <div>
                  {{ recipe.description }}
                </div>
              </v-col>
            </v-row>
          </v-tabs-window-item>
          <v-tabs-window-item value="ingredients">
            <v-row>
              <v-col>
                <v-card v-for="(ingredientGroup, i) in recipe.ingredientGroups" :key="i">
                  <v-card-title v-if="ingredientGroup.title">{{
                    ingredientGroup.title
                  }}</v-card-title>
                  <v-card-text>
                    <div v-for="ingredient in ingredientGroup.ingredients">
                      {{ ingredient.title }}{{ ingredient.measurementAmount
                      }}{{ ingredient.measurementType }}{{ ingredient.optional }}
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-tabs-window-item>
          <v-tabs-window-item value="directions">
            <v-row>
              <v-col>
                <v-card v-for="(direction, i) in recipe.directions" :key="i">
                  <v-card-title v-if="direction.title">{{ direction.title }}</v-card-title>
                  <v-card-text>
                    <div v-for="step in direction.steps">
                      {{ step }}
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-tabs-window-item>
        </v-tabs-window>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { useRecipesStore } from '@/stores/recipes'
import { storeToRefs } from 'pinia'
import { onMounted, watch, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const tab = ref('general') // Default tab value

const route = useRoute()
const router = useRouter()

const recipesStore = useRecipesStore()
const { recipe } = storeToRefs(recipesStore)

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
  // Set the tab based on the current hash in the URL
  const currentHash = route.hash
  if (currentHash in hashToTab) {
    tab.value = hashToTab[currentHash]
  }

  // Load the recipe based on the route parameter
  const recipeId = route.params.recipeId as string
  recipesStore.loadRecipe(recipeId)
})

// Watch for tab changes and update the URL hash
watch(tab, (newTab) => {
  const newHash = tabToHash[newTab]
  if (newHash && route.hash !== newHash) {
    router.replace({ hash: newHash }) // Update the URL hash without reloading the page
  }
})
</script>
