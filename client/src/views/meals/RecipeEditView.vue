<template>
  <v-form>
    <v-container v-if="editRecipe">
      <v-app-bar>
        <v-tabs style="width: 100%" v-model="tab" center-active show-arrows fixed-tabs>
          <v-tab value="general" text="General"></v-tab>
          <v-tab value="ingredients" text="Ingredients"></v-tab>
          <v-tab value="directions" text="Directions"></v-tab>
        </v-tabs>
        <template #prepend>
          <v-btn
            icon="mdi-close"
            @click="router.push({ name: 'recipe', params: { recipeId: route.params.recipeId } })"
          ></v-btn>
        </template>
        <template #append>
          <v-btn
            icon="mdi-content-save"
            @click="router.push({ name: 'recipe', params: { recipeId: route.params.recipeId } })"
          ></v-btn>
        </template>
      </v-app-bar>
      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="general">
          <v-container max-width="600" class="pa-1">
            <v-row>
              <v-col class="pt-5">
                <div class="text-h4">
                  <v-text-field density="compact" v-model="editRecipe.title"></v-text-field>
                </div>
                <div class="text-body-1">
                  <v-textarea density="compact" v-model="editRecipe.description"></v-textarea>
                </div>
              </v-col>
            </v-row>
            <v-row>
              <v-spacer></v-spacer>
              <v-col align-self="auto" cols="12" sm="8">
                <v-img
                  class="mt-1"
                  style="background-color: lightgray"
                  :src="editRecipe.imageUri"
                  cover
                ></v-img>
              </v-col>
              <v-spacer></v-spacer>
            </v-row>
          </v-container>
        </v-tabs-window-item>
        <v-tabs-window-item value="ingredients">
          <v-container max-width="600">
            <div
              v-for="(ingredientGroup, i) in editRecipe.ingredientGroups"
              :key="i"
              class="d-flex gap-2 mb-4"
            >
              <div class="button-column d-flex flex-column gap-1 mt-4">
                <v-btn
                  icon="mdi-arrow-up"
                  size="small"
                  @click="moveGroupUp(i)"
                  variant="text"
                  v-show="i > 0"
                ></v-btn>
                <v-btn
                  icon="mdi-arrow-down"
                  size="small"
                  @click="moveGroupDown(i)"
                  variant="text"
                  v-show="i < (editRecipe.ingredientGroups?.length || 0) - 1"
                ></v-btn>
              </div>
              <v-card class="flex-grow-1">
                <v-card-title>
                  <v-text-field v-model="ingredientGroup.title"></v-text-field>
                </v-card-title>
                <v-card-text>
                  <div
                    v-for="(ingredient, j) in ingredientGroup.ingredients"
                    :key="j"
                    class="d-flex gap-2 mb-2"
                  >
                    <div class="button-column d-flex flex-column gap-1">
                      <v-btn
                        icon="mdi-arrow-up"
                        size="x-small"
                        @click="moveIngredientUp(i, j)"
                        variant="text"
                        v-show="j > 0"
                      ></v-btn>
                      <v-btn
                        icon="mdi-arrow-down"
                        size="x-small"
                        @click="moveIngredientDown(i, j)"
                        variant="text"
                        v-show="j < (ingredientGroup.ingredients?.length || 0) - 1"
                      ></v-btn>
                    </div>
                    <div class="flex-grow-1">
                      <v-row dense>
                        <v-col cols="4" sm="2">
                          <v-text-field
                            density="compact"
                            variant="outlined"
                            hide-details
                            v-model="ingredient.measurementAmount"
                            class="mt-0"
                          ></v-text-field>
                        </v-col>
                        <v-col cols="8" sm="4">
                          <v-select
                            density="compact"
                            variant="outlined"
                            hide-details
                            v-model="ingredient.measurementType"
                            :items="measurementSelect"
                            item-title="title"
                            item-value="value"
                            class="mt-0"
                          ></v-select>
                        </v-col>
                        <v-col cols="6" sm="6">
                          <v-text-field
                            density="compact"
                            variant="outlined"
                            hide-details
                            v-model="ingredient.title"
                            class="mt-0"
                          ></v-text-field>
                        </v-col>
                        <v-col cols="6" sm="12" class="d-flex justify-end">
                          <v-checkbox
                            density="compact"
                            hide-details
                            v-model="ingredient.optional"
                            label="Optional"
                            class="mt-0"
                          ></v-checkbox>
                        </v-col>
                      </v-row>
                    </div>
                  </div>
                </v-card-text>
              </v-card>
            </div>
          </v-container>
        </v-tabs-window-item>
        <v-tabs-window-item value="directions">
          <v-container max-width="600">
            <div v-for="(direction, i) in editRecipe.directions" :key="i" class="d-flex gap-2 mb-4">
              <div class="button-column d-flex flex-column gap-1 mt-4">
                <v-btn
                  icon="mdi-arrow-up"
                  size="small"
                  @click="moveDirectionUp(i)"
                  variant="text"
                  v-show="i > 0"
                ></v-btn>
                <v-btn
                  icon="mdi-arrow-down"
                  size="small"
                  @click="moveDirectionDown(i)"
                  variant="text"
                  v-show="i < (editRecipe.directions?.length || 0) - 1"
                ></v-btn>
              </div>
              <v-card class="flex-grow-1">
                <v-card-title><v-text-field v-model="direction.title"></v-text-field></v-card-title>
                <v-card-text>
                  <v-list>
                    <div
                      v-for="(step, n) in getSteps(direction)"
                      :key="n"
                      class="d-flex gap-2 mb-2"
                    >
                      <div class="button-column d-flex flex-column gap-1">
                        <v-btn
                          icon="mdi-arrow-up"
                          size="x-small"
                          @click="moveStepUp(i, n)"
                          variant="text"
                          v-show="n > 0"
                        ></v-btn>
                        <v-btn
                          icon="mdi-arrow-down"
                          size="x-small"
                          @click="moveStepDown(i, n)"
                          variant="text"
                          v-show="n < getSteps(direction).length - 1"
                        ></v-btn>
                      </div>
                      <v-list-item class="flex-grow-1">
                        <div class="font-weight-bold">Step {{ n + 1 }}</div>
                        <v-textarea
                          :model-value="getSteps(direction)[n]"
                          @update:model-value="updateStep(direction, n, $event)"
                        ></v-textarea>
                      </v-list-item>
                    </div>
                  </v-list>
                </v-card-text>
              </v-card>
            </div>
          </v-container>
        </v-tabs-window-item>
      </v-tabs-window>
    </v-container>
  </v-form>
</template>

<script setup lang="ts">
import type {
  Recipe,
  Recipe_Direction,
  Recipe_MeasurementType,
} from '@/genapi/api/meals/recipe/v1alpha1'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRecipesStore } from '@/stores/recipes'
import { storeToRefs } from 'pinia'
import { onMounted, watch, ref } from 'vue'
import type { Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const tab = ref('general') // Default tab value

const route = useRoute()
const router = useRouter()

const recipesStore = useRecipesStore()
const editRecipe: Ref<Recipe | undefined, Recipe | undefined> = ref(undefined)
const editDirections = ref<Recipe_Direction[]>([])

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

// Map string to MeasurmentType
const measurementSelect = [
  { title: '', value: 'MEASUREMENT_TYPE_UNSPECIFIED' },
  { title: 'tablespoons', value: 'MEASUREMENT_TYPE_TABLESPOON' },
  { title: 'teaspoons', value: 'MEASUREMENT_TYPE_TEASPOON' },
  { title: 'ounces', value: 'MEASUREMENT_TYPE_OUNCE' },
  { title: 'pounds', value: 'MEASUREMENT_TYPE_POUND' },
  { title: 'grams', value: 'MEASUREMENT_TYPE_GRAM' },
  { title: 'milliliters', value: 'MEASUREMENT_TYPE_MILLILITER' },
  { title: 'liters', value: 'MEASUREMENT_TYPE_LITER' },
]

function convertMeasurementTypeToString(type: Recipe_MeasurementType | undefined): string {
  return type ? measurementTypeToString[type] || '' : ''
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
  // make a deap copy of editRecipe.value = recipe.value
  editRecipe.value = JSON.parse(JSON.stringify(recipesStore.recipe))
  // Set the steps for the directions ignoring any that are undefined
  editDirections.value = editRecipe.value?.directions || []

  breadcrumbStore.setBreadcrumbs([
    { title: 'Recipes', to: { name: 'recipes' } },
    {
      title: recipesStore.recipe?.title ? recipesStore.recipe.title : '',
      to: { name: 'recipe', params: { recipeId: recipeId } },
    },
    { title: 'Edit', to: { name: 'recipeEdit', params: { recipeId: recipeId } } },
  ])
})

// Watch for tab changes and update the URL hash
watch(tab, (newTab) => {
  const newHash = tabToHash[newTab]
  if (newHash && route.hash !== newHash) {
    router.replace({ hash: newHash }) // Update the URL hash without reloading the page
  }
})

function moveGroupUp(index: number) {
  if (!editRecipe.value?.ingredientGroups) return
  if (index > 0) {
    const groups = editRecipe.value.ingredientGroups
    const temp = groups[index]
    groups[index] = groups[index - 1]
    groups[index - 1] = temp
  }
}

function moveGroupDown(index: number) {
  if (!editRecipe.value?.ingredientGroups) return
  const groups = editRecipe.value.ingredientGroups
  if (index < groups.length - 1) {
    const temp = groups[index]
    groups[index] = groups[index + 1]
    groups[index + 1] = temp
  }
}

function moveIngredientUp(groupIndex: number, ingredientIndex: number) {
  if (!editRecipe.value?.ingredientGroups?.[groupIndex]?.ingredients) return
  const ingredients = editRecipe.value.ingredientGroups[groupIndex].ingredients
  if (ingredientIndex > 0) {
    const temp = ingredients[ingredientIndex]
    ingredients[ingredientIndex] = ingredients[ingredientIndex - 1]
    ingredients[ingredientIndex - 1] = temp
  }
}

function moveIngredientDown(groupIndex: number, ingredientIndex: number) {
  if (!editRecipe.value?.ingredientGroups?.[groupIndex]?.ingredients) return
  const ingredients = editRecipe.value.ingredientGroups[groupIndex].ingredients
  if (ingredientIndex < ingredients.length - 1) {
    const temp = ingredients[ingredientIndex]
    ingredients[ingredientIndex] = ingredients[ingredientIndex + 1]
    ingredients[ingredientIndex + 1] = temp
  }
}

function moveDirectionUp(index: number) {
  if (!editRecipe.value?.directions) return
  if (index > 0) {
    const directions = editRecipe.value.directions
    const temp = directions[index]
    directions[index] = directions[index - 1]
    directions[index - 1] = temp
  }
}

function moveDirectionDown(index: number) {
  if (!editRecipe.value?.directions) return
  const directions = editRecipe.value.directions
  if (index < directions.length - 1) {
    const temp = directions[index]
    directions[index] = directions[index + 1]
    directions[index + 1] = temp
  }
}

function moveStepUp(directionIndex: number, stepIndex: number) {
  if (!editRecipe.value?.directions?.[directionIndex]?.steps) return
  const steps = editRecipe.value.directions[directionIndex].steps
  if (stepIndex > 0) {
    const temp = steps[stepIndex]
    steps[stepIndex] = steps[stepIndex - 1]
    steps[stepIndex - 1] = temp
  }
}

function moveStepDown(directionIndex: number, stepIndex: number) {
  if (!editRecipe.value?.directions?.[directionIndex]?.steps) return
  const steps = editRecipe.value.directions[directionIndex].steps
  if (stepIndex < steps.length - 1) {
    const temp = steps[stepIndex]
    steps[stepIndex] = steps[stepIndex + 1]
    steps[stepIndex + 1] = temp
  }
}

function getSteps(direction: any): string[] {
  return direction?.steps || []
}

function updateStep(direction: any, index: number, value: string) {
  if (!direction.steps) {
    direction.steps = []
  }
  direction.steps[index] = value
}
</script>

<style scoped>
.button-column {
  width: 32px;
  align-items: center;
}

.gap-1 {
  gap: 4px;
}

.gap-2 {
  gap: 8px;
}
</style>
