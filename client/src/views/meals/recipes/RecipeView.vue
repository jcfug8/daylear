<template>
  <div v-if="recipe" class="position-relative" style="padding-bottom: 84px;">
    <v-container max-width="600" class="pa-1">
      <v-row no-gutters>
        <v-col cols="10">
          <div class="text-h5">
            {{ recipe.title }}
          </div>
        </v-col>
        <v-col cols="2" class="text-right">
          <v-btn class="text-h5" id="recipe-menu-btn" variant="text">
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
          <v-menu activator="#recipe-menu-btn">
            <v-list>
              <v-list-item v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel)" color="primary"
                density="compact" :to="'/' + recipe.name + '/edit'">
                <v-icon>mdi-pencil</v-icon>
                Edit
              </v-list-item>
              <v-list-item color="success" density="compact" @click="openCreateEventDialog">
                <v-icon>mdi-calendar-plus</v-icon>
                Schedule Event
              </v-list-item>
              <v-list-item density="compact"
                v-if="!hasReadPermission(recipe.recipeAccess?.permissionLevel) && !recipe.recipeAccess" color="primary"
                @click="handleRequestAccess" :loading="requestingAccess">
                <v-icon>mdi-account-plus</v-icon>Request Access
              </v-list-item>
              <v-list-item
                v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel) && recipe.visibility !== 'VISIBILITY_LEVEL_HIDDEN'"
                color="primary" density="compact" @click="showShareDialog = true">
                <v-icon>mdi-share-variant</v-icon>
                Manage Access
              </v-list-item>
              <v-list-item density="compact"
                v-if="recipe.recipeAccess?.state === 'ACCESS_STATE_PENDING' && recipe.recipeAccess?.acceptTarget !== 'ACCEPT_TARGET_RECIPIENT'"
                color="warning" @click="showCancelRequestDialog = true">
                <v-icon>mdi-close</v-icon>Cancel Request
              </v-list-item>
              <v-list-item
                v-if="!hasAdminPermission(recipe.recipeAccess?.permissionLevel) && recipe.recipeAccess?.state === 'ACCESS_STATE_ACCEPTED'"
                color="warning" density="compact" @click="showRemoveAccessDialog = true">
                <v-icon>mdi-link-variant-off</v-icon>
                Remove Access
              </v-list-item>
              <v-list-item v-if="hasAdminPermission(recipe.recipeAccess?.permissionLevel)" color="error"
                density="compact" @click="showDeleteDialog = true">
                <v-icon>mdi-delete</v-icon>
                Delete
              </v-list-item>
            </v-list>
          </v-menu>
        </v-col>
        <v-col cols="12">
          <div class="image-container">
            <v-icon size="24" :color="recipe.favorited ? 'red' : 'black'" class="favorite-heart" @click="toggleFavorite"
              style="cursor: pointer;">
              {{ recipe.favorited ? 'mdi-heart' : 'mdi-heart-outline' }}
            </v-icon>
            <v-img v-if="recipe.imageUri" style="background-color: lightgray" :src="recipe.imageUri" cover height="150">
              <template #placeholder>
                <v-row class="fill-height ma-0" align="center" justify="center">
                  <v-progress-circular indeterminate color="grey lighten-5"></v-progress-circular>
                </v-row>
              </template>
            </v-img>
            <div v-else class="d-flex align-center justify-center"
              style="background-color: lightgray; height: 300px; border-radius: 4px;">
              <div class="text-center">
                <v-icon size="64" color="grey-darken-1">mdi-image-outline</v-icon>
                <div class="text-grey-darken-1 mt-2">No image available</div>
              </div>
            </div>
            <v-btn v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel) && !recipe.imageUri"
              class="generate-image-btn" color="warning" @click="openGenerateImageModal"
              style="position: absolute; bottom: 8px; right: 8px; z-index: 2;" :loading="generatingImage"
              :disabled="generatingImage" title="Generate Image">
              <v-icon>mdi-image-auto-adjust</v-icon>
              Generate Image
            </v-btn>
          </div>
        </v-col>
      </v-row>
      <v-row no-gutters>
        <v-col class="text-caption"  v-if="recipe.prepDuration && recipe.prepDuration !== '0s'">
          <v-icon>mdi-knife</v-icon> {{ parseDuration(recipe.prepDuration ? recipe.prepDuration : "0") }}mins
        </v-col>
        <v-col class="text-caption"  v-if="recipe.cookDuration && recipe.cookDuration !== '0s'">
          <v-icon>mdi-stove</v-icon> {{ parseDuration(recipe.cookDuration ? recipe.cookDuration : "0") }}mins
        </v-col>
        <v-col class="text-caption"  v-if="recipe.totalDuration && recipe.totalDuration !== '0s'">
          <v-icon>mdi-clock</v-icon> {{ parseDuration(recipe.totalDuration ? recipe.totalDuration : "0") }}mins
        </v-col>
        <v-col class="text-caption"  v-if="recipe.yieldAmount">
          <v-icon>mdi-food-apple-outline</v-icon> {{ recipe.yieldAmount }}
        </v-col>
      </v-row>
      <v-row no-gutters>
        <v-col cols="12" class="text-center mb-1">
          <v-chip class="ma-1" v-if="recipe.cookingMethod" label color="green"><v-icon icon="mdi-toaster-oven"
              start></v-icon>{{
                recipe.cookingMethod }}</v-chip>
          <v-chip class="ma-1" v-for="category in recipe.categories" :key="category" label color="blue"><v-icon
              icon="mdi-label" start></v-icon>{{ category }}</v-chip>
          <v-chip class="ma-1" v-for="cuisine in recipe.cuisines" :key="cuisine" label color="orange"><v-icon
              icon="mdi-chef-hat" start></v-icon>{{ cuisine }}</v-chip>
        </v-col>
        <v-col cols="12" class="text-caption">
          <div v-if="recipe.citation">
            <span v-if="isUrl(recipe.citation)">
              <a :href="recipe.citation" target="_blank" rel="noopener">{{ recipe.citation }}</a>
            </span>
            <span v-else>
              {{ recipe.citation }}
            </span>
          </div>
        </v-col>
        <v-col cols="12">
          <div class="text-body-2">
            {{ recipe.description }}
          </div>
        </v-col>
        <!-- <v-col cols="12">
          <div class="mt-4">
            <div v-if="selectedVisibilityDescription" class="mt-2">
              <v-alert :icon="selectedVisibilityIcon" density="compact" variant="tonal"
                :color="selectedVisibilityColor">
                <div class="text-body-2">
                  <strong>{{ selectedVisibilityLabel }}:</strong> {{ selectedVisibilityDescription }}
                </div>
              </v-alert>
            </div>
          </div>
        </v-col> -->
      </v-row>
    </v-container>
    <v-divider></v-divider>
    <v-tabs class="position-sticky" style="width: 100%;z-index: 2;top: 48px;" bg-color="white" v-model="tab" center-active show-arrows fixed-tabs>
      <v-tab value="ingredients">
        <v-icon>mdi-food-apple-outline</v-icon>
        <span class="text-caption">Ingredients</span>
      </v-tab>
      <v-tab value="directions">
        <v-icon>mdi-format-list-numbered</v-icon>
        <span class="text-caption">Directions</span>
      </v-tab>
    </v-tabs>
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="ingredients">
        <v-container max-width="600">
          <!-- Multiplier Card -->
          <v-card class="mb-2" outlined>
            <v-card-text class="d-flex align-center">
              <span class="mr-4 font-weight-medium">Multiplier</span>
              <v-select v-model="ingredientMultiplier" :items="multiplierOptions" :item-title="x => `${x}x`"
                :item-value="x => x" :menu-props="{ maxHeight: '200px' }" hide-details density="compact"
                style="max-width: 120px;" />
            </v-card-text>
          </v-card>
          <!-- Ingredients List -->
          <v-card class="my-1" v-for="(ingredientGroup, i) in recipe.ingredientGroups" :key="i">
            <v-card-title v-if="ingredientGroup.title">{{ ingredientGroup.title }}</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item class="ingredient-item" v-for="(ingredient, j) in ingredientGroup.ingredients" :key="j">
                  <div class="d-flex align-center" style="width: 100%;">
                    <v-checkbox v-model="checkedIngredients[i][j]" class="mr-2" hide-details density="compact"
                      style="margin-bottom: 0;" />
                    <span :style="checkedIngredients[i][j] ? 'text-decoration: line-through; color: #888;' : ''">
                      <strong>
                        <span v-if="ingredient.measurementAmount && ingredient.measurementType">
                          <span>
                            {{ isFranctional(ingredient.measurementType)
                              ? toFraction(getDisplayAmount(ingredient.measurementAmount, i, j,
                                ingredient.measurementType))
                              : formatAmountDisplay(getDisplayAmount(ingredient.measurementAmount, i, j,
                                ingredient.measurementType))
                            }}
                          </span>
                          <span
                            :style="'cursor: pointer; text-decoration: underline dotted; -webkit-text-decoration: underline dotted; margin-left: 2px;'"
                            @click="handleUnitClick(i, j, ingredient.measurementType)" title="Click to change unit">
                            {{ measurementTypeLabel(getDisplayUnit(i, j, ingredient.measurementType),
                              getDisplayAmount(ingredient.measurementAmount, i, j, ingredient.measurementType)) }}
                          </span>
                        </span>
                        <template
                          v-if="ingredient.measurementConjunction && ingredient.secondMeasurementAmount && ingredient.secondMeasurementType">
                          {{ renderConjunction(ingredient.measurementConjunction) }}
                          <span v-if="ingredient.secondMeasurementAmount && ingredient.secondMeasurementType">
                            <span>
                              {{ isFranctional(ingredient.secondMeasurementType)
                                ? toFraction(getDisplayAmount(ingredient.secondMeasurementAmount, i, j,
                                  ingredient.secondMeasurementType, true))
                                : formatAmountDisplay(getDisplayAmount(ingredient.secondMeasurementAmount, i, j,
                                  ingredient.secondMeasurementType, true))
                              }}
                            </span>
                            <span :style="'cursor: pointer; text-decoration: underline dotted; margin-left: 2px;'"
                              @click="handleUnitClick(i, j, ingredient.secondMeasurementType, true)"
                              title="Click to change unit">
                              {{ measurementTypeLabel(getDisplayUnit(i, j, ingredient.secondMeasurementType, true),
                                getDisplayAmount(ingredient.secondMeasurementAmount, i, j,
                                  ingredient.secondMeasurementType, true)) }}
                            </span>
                          </span>
                        </template>
                      </strong>
                      {{ ingredient.title }} <span v-if="ingredient.optional">(optional)</span>
                    </span>
                  </div>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-container>
      </v-tabs-window-item>
      <v-tabs-window-item value="directions">
        <v-container max-width="600">
          <v-card v-for="(direction, i) in recipe.directions" :key="i">
            <v-card-title v-if="direction.title">{{ direction.title }}</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item v-for="(step, n) in direction.steps" :key="n">
                  <div class="d-flex align-start" style="width: 100%;">
                    <v-checkbox v-model="checkedDirections[i][n]" class="mr-2" hide-details density="compact"
                      style="margin-bottom: 0; align-self: flex-start;" />
                    <span :style="checkedDirections[i][n] ? 'text-decoration: line-through; color: #888;' : ''"
                      style="margin-left: 4px; flex: 1 1 0; word-break: break-word;">
                      <span class="font-weight-bold" style="display: inline;">Step {{ n + 1 }}</span>
                      {{ step }}
                    </span>
                  </div>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-container>
      </v-tabs-window-item>
    </v-tabs-window>

    <template
      v-if="recipe.recipeAccess?.acceptTarget === 'ACCEPT_TARGET_RECIPIENT' && recipe.recipeAccess?.state === 'ACCESS_STATE_PENDING'">
      <v-btn color="success" class="mb-2" block :loading="acceptingRecipe"
        @click="acceptRecipe(recipe.recipeAccess?.name)">
        Accept Recipe
      </v-btn>
      <v-btn color="error" class="mb-4" block :loading="decliningRecipe" @click="declineRecipe">
        Decline
      </v-btn>
    </template>

    <!-- Share Dialog -->
    <ShareDialog v-model="showShareDialog" title="Share Recipe" :allowCircleShare="true"
      :currentAccesses="currentAccesses" :sharing="sharing" :sharePermissionLoading="updatingPermission"
      :userPermissionLevel="recipe.recipeAccess?.permissionLevel" :allowPermissionOptions="allowPermissionOptions"
      @share-user="shareWithUser" @share-circle="shareWithCircle" @remove-access="unshareRecipe"
      @permission-change="updatePermission" @approve-access="acceptRecipeFromShareDialog" />

    <!-- Event Create Dialog -->
    <EventCreateDialog v-model="showCreateEventDialog" :calendars="writableCalendars" @created="handleEventCreated" />
  </div>
  <!-- Remove Access Dialog -->
  <v-dialog v-model="showRemoveAccessDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Remove Access
      </v-card-title>
      <v-card-text>
        Are you sure you want to remove your access to this recipe? You may no longer be able to view it.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showRemoveAccessDialog = false">
          Cancel
        </v-btn>
        <v-btn color="error" @click="handleRemoveAccess" :loading="removingAccess">
          Remove Access
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Delete Dialog -->
  <v-dialog v-model="showDeleteDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Delete Recipe
      </v-card-title>
      <v-card-text>
        Are you sure you want to delete this recipe? This action will also delete the recipe for any users or circles
        that
        can view it.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showDeleteDialog = false">
          Cancel
        </v-btn>
        <v-btn color="error" @click="handleDelete" :loading="deleting">
          Delete
        </v-btn>
      </v-card-actions>
    </v-card>

  </v-dialog>
  <!-- Cancel Request Dialog -->
  <v-dialog v-model="showCancelRequestDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Cancel Access Request
      </v-card-title>
      <v-card-text>
        Are you sure you want to cancel your access request to this recipe?
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showCancelRequestDialog = false">
          Cancel
        </v-btn>
        <v-btn color="error" @click="handleCancelRequest" :loading="cancelingRequest">
          Cancel Request
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Add modal for generated image -->
  <v-dialog v-model="showGeneratedImageModal" max-width="600">
    <v-card>
      <v-card-title class="text-h5">Generated Recipe Image</v-card-title>
      <v-card-text>
        <div v-if="generatingImage" class="d-flex justify-center my-4">
          <v-progress-circular indeterminate color="primary" size="48" />
        </div>
        <div v-else-if="generatedImageUrl">
          <img :src="generatedImageUrl" style="max-width: 100%; max-height: 400px; display: block; margin: 0 auto;" />
        </div>
        <v-alert v-if="generateImageError" type="error" class="mt-2">{{ generateImageError }}</v-alert>
        <div class="d-flex justify-center align-center mt-4">
          <v-btn color="primary" class="mr-2" :disabled="generatingImage" @click="startImageGeneration">
            <v-icon left>mdi-image-auto-adjust</v-icon>
            Generate Image
          </v-btn>
          <v-btn v-if="generatedImageUrl" color="success"
            :disabled="!generatedImageBlob || updatingImage || generatingImage" @click="updateRecipeImage">
            <v-icon left>mdi-check-circle</v-icon>
            Use Image
          </v-btn>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="closeGenerateImageModal">Cancel</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import type { Recipe_MeasurementType, apitypes_VisibilityLevel, Recipe_Ingredient_MeasurementConjunction } from '@/genapi/api/meals/recipe/v1alpha1'
import type { Access, CreateAccessRequest, ListAccessesRequest, DeleteAccessRequest } from '@/genapi/api/meals/recipe/v1alpha1'
import type { PermissionLevel } from '@/genapi/api/types'
import type { Calendar, Event } from '@/genapi/api/calendars/calendar/v1alpha1'
import { useRecipesStore } from '@/stores/recipes'
import { useRecipeFormStore } from '@/stores/recipeForm'
import { storeToRefs } from 'pinia'
import { onMounted, onBeforeUnmount, watch, computed, ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { recipeService, recipeAccessService, fileService, eventRecipeService } from '@/api/api'
import { useAuthStore } from '@/stores/auth'
import { hasWritePermission, hasAdminPermission, hasReadPermission } from '@/utils/permissions'
import ShareDialog from '@/components/common/ShareDialog.vue'
import { useAlertStore } from '@/stores/alerts'
import EventCreateDialog from '@/views/calendar/event/EventCreateDialog.vue'
import { useCalendarsStore } from '@/stores/calendar'

const authStore = useAuthStore()
const recipesStore = useRecipesStore()
const recipeFormStore = useRecipeFormStore()
const calendarsStore = useCalendarsStore()
const route = useRoute()
const router = useRouter()
const alertsStore = useAlertStore()
const trimmedRecipeName = computed(() => {
  return route.path.substring(route.path.indexOf('recipes/'))
})

const { recipe } = storeToRefs(recipesStore)

// Add local state for checked ingredients and directions
const checkedIngredients = reactive<Array<Array<boolean>>>([])
const checkedDirections = reactive<Array<Array<boolean>>>([])

// Multiplier for ingredients
const ingredientMultiplier = ref(1)
const multiplierOptions = [0.5, 1, 2, 3, 4, 5, 6]

// --- Refactored unit conversion logic ---
// Volume: all multipliers are in terms of teaspoons (smallest)
const VOLUME_UNIT_MULTIPLIERS: Partial<Record<Recipe_MeasurementType, number>> = {
  MEASUREMENT_TYPE_TEASPOON: 1,
  MEASUREMENT_TYPE_TABLESPOON: 3,
  MEASUREMENT_TYPE_CUP: 48,
  MEASUREMENT_TYPE_MILLILITER: 0.202884,
  MEASUREMENT_TYPE_LITER: 202.884,
}
const VOLUME_UNITS: Recipe_MeasurementType[] = [
  'MEASUREMENT_TYPE_TEASPOON',
  'MEASUREMENT_TYPE_TABLESPOON',
  'MEASUREMENT_TYPE_CUP',
  'MEASUREMENT_TYPE_MILLILITER',
  'MEASUREMENT_TYPE_LITER',
]
// Weight: all multipliers are in terms of grams (smallest)
const WEIGHT_UNIT_MULTIPLIERS: Partial<Record<Recipe_MeasurementType, number>> = {
  MEASUREMENT_TYPE_GRAM: 1,
  MEASUREMENT_TYPE_OUNCE: 28.3495,
  MEASUREMENT_TYPE_POUND: 453.592,
}
const WEIGHT_UNITS: Recipe_MeasurementType[] = [
  'MEASUREMENT_TYPE_GRAM',
  'MEASUREMENT_TYPE_OUNCE',
  'MEASUREMENT_TYPE_POUND',
]

function getUnitGroupAndList(unit: Recipe_MeasurementType | undefined): { group: 'volume' | 'weight' | null, units: Recipe_MeasurementType[] } {
  if (!unit || unit === 'MEASUREMENT_TYPE_UNSPECIFIED') return { group: null, units: [] }
  if (VOLUME_UNIT_MULTIPLIERS[unit] !== undefined) {
    return { group: 'volume', units: VOLUME_UNITS }
  }
  if (WEIGHT_UNIT_MULTIPLIERS[unit] !== undefined) {
    return { group: 'weight', units: WEIGHT_UNITS }
  }
  return { group: null, units: [] }
}

function getNextUnit(unit: Recipe_MeasurementType | undefined) {
  if (!unit || unit === 'MEASUREMENT_TYPE_UNSPECIFIED') return unit
  const { units } = getUnitGroupAndList(unit)
  if (!units.length) return unit
  const idx = units.indexOf(unit)
  return units[(idx + 1) % units.length]
}

function convertAmount(amount: number, from: Recipe_MeasurementType, to: Recipe_MeasurementType): number {
  if (from === 'MEASUREMENT_TYPE_UNSPECIFIED' || to === 'MEASUREMENT_TYPE_UNSPECIFIED') return amount
  if (VOLUME_UNIT_MULTIPLIERS[from] !== undefined && VOLUME_UNIT_MULTIPLIERS[to] !== undefined) {
    const base = amount * (VOLUME_UNIT_MULTIPLIERS[from] as number)
    return base / (VOLUME_UNIT_MULTIPLIERS[to] as number)
  }
  if (WEIGHT_UNIT_MULTIPLIERS[from] !== undefined && WEIGHT_UNIT_MULTIPLIERS[to] !== undefined) {
    const base = amount * (WEIGHT_UNIT_MULTIPLIERS[from] as number)
    return base / (WEIGHT_UNIT_MULTIPLIERS[to] as number)
  }
  return amount
}

// function getUnitGroup(unit: string) {
//   // For cycling logic compatibility
//   if (VOLUME_UNIT_MULTIPLIERS[unit as Recipe_MeasurementType] !== undefined) {
//     return { units: VOLUME_UNITS }
//   }
//   if (WEIGHT_UNIT_MULTIPLIERS[unit as Recipe_MeasurementType] !== undefined) {
//     return { units: WEIGHT_UNITS }
//   }
//   return undefined
// }

function getDisplayUnitKey(i: number, j: number, secondary = false) {
  return `${i}_${j}${secondary ? '_2' : ''}`
}

function handleUnitClick(i: number, j: number, unit: Recipe_MeasurementType | undefined, secondary = false) {
  if (!unit || unit === 'MEASUREMENT_TYPE_UNSPECIFIED') return
  const key = getDisplayUnitKey(i, j, secondary)
  const current = String(ingredientDisplayUnits[key] ?? unit)
  const next = getNextUnit(current as Recipe_MeasurementType)
  ingredientDisplayUnits[key] = String(next)
}

function getDisplayAmount(amount: number | undefined, i: number, j: number, unit: Recipe_MeasurementType | undefined, secondary = false): number {
  if (typeof amount !== 'number' || !unit) return 0
  const key = getDisplayUnitKey(i, j, secondary)
  const displayUnit = (ingredientDisplayUnits[key] as Recipe_MeasurementType) || unit
  const converted = convertAmount(amount * ingredientMultiplier.value, unit, displayUnit)
  return converted
}

function getDisplayUnit(i: number, j: number, unit: Recipe_MeasurementType | undefined, secondary = false): Recipe_MeasurementType {
  if (!unit) return 'MEASUREMENT_TYPE_UNSPECIFIED'
  const key = getDisplayUnitKey(i, j, secondary)
  return (ingredientDisplayUnits[key] as Recipe_MeasurementType) || unit
}

function initializeCheckedState() {
  checkedIngredients.length = 0
  recipe.value?.ingredientGroups?.forEach((group, i) => {
    checkedIngredients[i] = []
    group.ingredients?.forEach((_, j) => {
      checkedIngredients[i][j] = false
    })
  })
  checkedDirections.length = 0
  recipe.value?.directions?.forEach((dir, i) => {
    checkedDirections[i] = []
    dir.steps?.forEach((_, n) => {
      checkedDirections[i][n] = false
    })
  })
}

watch(recipe, initializeCheckedState, { immediate: true })



// *** Tabs ***

const { activeTab } = storeToRefs(recipeFormStore)
onBeforeUnmount(() => {
  if (router.currentRoute.value.name !== 'recipeEdit' && router.currentRoute.value.name !== 'circleRecipeEdit') {
    recipeFormStore.setActiveTab('general')
  }
})

const tab = computed({
  get: () => activeTab.value,
  set: (value) => recipeFormStore.setActiveTab(value),
})

watch(tab, (newTab) => {
  const newHash = tabToHash[newTab]
  if (newHash && route.hash !== newHash) {
    router.replace({ hash: newHash }) // Update the URL hash without reloading the page
  }
})

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


// *** Visibility ***

const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === recipe.value?.visibility)
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

// Track current display unit for each ingredient (by group and index)
const ingredientDisplayUnits = reactive<{ [key: string]: string }>({})

// *** Remove Access ***
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

async function handleRemoveAccess() {
  if (!recipe.value?.recipeAccess?.name) return

  removingAccess.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: recipe.value.recipeAccess.name
    }

    await recipeAccessService.DeleteAccess(deleteRequest)
    router.push(route.params.circleId ? { name: 'circle', params: { circleId: route.params.circleId } } : { name: 'recipes' })
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    removingAccess.value = false
    showRemoveAccessDialog.value = false
  }
}

// *** Cancel Request ***

const showCancelRequestDialog = ref(false)
const cancelingRequest = ref(false)

async function handleCancelRequest() {
  if (!recipe.value?.recipeAccess?.name) return

  cancelingRequest.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: recipe.value.recipeAccess.name
    }

    await recipeAccessService.DeleteAccess(deleteRequest)
    await recipesStore.loadRecipe(recipe.value.name!)
    alertsStore.addAlert('Access request cancelled.', 'info')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    cancelingRequest.value = false
    showCancelRequestDialog.value = false
  }
}

// *** Delete Recipe ***

const showDeleteDialog = ref(false)
const deleting = ref(false)

async function handleDelete() {
  if (!recipe.value?.name) return

  deleting.value = true
  try {
    await recipeService.DeleteRecipe({
      name: recipe.value.name
    })
    router.push({ name: 'recipes' })
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    deleting.value = false
    showDeleteDialog.value = false
  }
}

// *** Accept/Decline Recipe ***

const acceptingRecipe = ref(false)
const decliningRecipe = ref(false)

async function acceptRecipeFromShareDialog(recipeAccessName: string | undefined) {
  if (!recipe.value?.recipeAccess?.name) return
  if (!recipeAccessName) return
  await acceptRecipe(recipeAccessName)
  await fetchRecipeRecipients()
}

async function acceptRecipe(recipeAccessName: string | undefined) {
  if (!recipe.value?.recipeAccess?.name || !authStore.user?.name) return
  if (!recipeAccessName) return
  acceptingRecipe.value = true
  try {
    await recipesStore.acceptRecipe(recipeAccessName)
    recipesStore.loadRecipe(recipe.value?.name ?? '')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    acceptingRecipe.value = false
  }
}

async function declineRecipe() {
  if (!recipe.value?.recipeAccess?.name) return
  decliningRecipe.value = true
  try {
    await recipesStore.deleteRecipeAccess(recipe.value.recipeAccess.name)
    router.push(route.params.circleId ? { name: 'circle', params: { circleId: route.params.circleId } } : { name: 'recipes' })
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    decliningRecipe.value = false
  }
}

// *** Measurements and Ingredients ***

function renderConjunction(conjunction: Recipe_Ingredient_MeasurementConjunction | undefined): string {
  switch (conjunction) {
    case 'MEASUREMENT_CONJUNCTION_AND':
      return ' +';
    case 'MEASUREMENT_CONJUNCTION_TO':
      return ' -';
    case 'MEASUREMENT_CONJUNCTION_OR':
      return ' or';
    default:
      return '';
  }
}

function measurementTypeLabel(type: Recipe_MeasurementType | undefined, amount: number | undefined): string {
  const singular: Record<string, string> = {
    MEASUREMENT_TYPE_CUP: 'cup',
    MEASUREMENT_TYPE_TABLESPOON: 'tablespoon',
    MEASUREMENT_TYPE_TEASPOON: 'teaspoon',
    MEASUREMENT_TYPE_OUNCE: 'ounce',
    MEASUREMENT_TYPE_POUND: 'pound',
    MEASUREMENT_TYPE_GRAM: 'gram',
    MEASUREMENT_TYPE_MILLILITER: 'milliliter',
    MEASUREMENT_TYPE_LITER: 'liter',
  }
  const plural: Record<string, string> = {
    MEASUREMENT_TYPE_CUP: 'cups',
    MEASUREMENT_TYPE_TABLESPOON: 'tablespoons',
    MEASUREMENT_TYPE_TEASPOON: 'teaspoons',
    MEASUREMENT_TYPE_OUNCE: 'ounces',
    MEASUREMENT_TYPE_POUND: 'pounds',
    MEASUREMENT_TYPE_GRAM: 'grams',
    MEASUREMENT_TYPE_MILLILITER: 'milliliters',
    MEASUREMENT_TYPE_LITER: 'liters',
  }
  if (!type || type === 'MEASUREMENT_TYPE_UNSPECIFIED') return ''
  if ((amount ?? 0) <= 1) return singular[type] || ''
  return plural[type] || ''
}

function isFranctional(type: Recipe_MeasurementType | undefined): boolean {
  return [
    'MEASUREMENT_TYPE_CUP',
    'MEASUREMENT_TYPE_TABLESPOON',
    'MEASUREMENT_TYPE_TEASPOON',
    'MEASUREMENT_TYPE_POUND',
  ].includes(type as string)
}

function toFraction(amount: number): string {
  if (isNaN(amount)) return ''
  const whole = Math.floor(amount)
  const frac = amount - whole
  // Find the closest common fraction denominator
  const denominators = [2, 3, 4, 8, 16]
  let best = { num: 0, den: 1, diff: 1 }
  for (const den of denominators) {
    const num = Math.round(frac * den)
    const diff = Math.abs(frac - num / den)
    if (num > 0 && diff < best.diff) {
      best = { num, den, diff }
    }
  }
  let result = ''
  if (whole > 0) result += whole
  if (best.num > 0) {
    if (whole > 0) result += ' '
    result += `${best.num}/${best.den}`
  }
  if (result === '') result = amount.toLocaleString(undefined, { maximumFractionDigits: 3 })
  return result
}

// Helper to format small values for display
function formatAmountDisplay(val: number): string {
  if (val === 0) return '0'
  if (val > 0 && val < 0.01) return '< 0.01'
  if (val > 0 && val < 1) return val.toLocaleString(undefined, { maximumFractionDigits: 3 })
  return val.toLocaleString(undefined, { maximumFractionDigits: 2 })
}

// *** Generate Image ***

const showGeneratedImageModal = ref(false)
const generatingImage = ref(false)
const generatedImageBlob = ref<Blob | null>(null)
const generatedImageUrl = ref<string | null>(null)
const generateImageError = ref<string | null>(null)

function openGenerateImageModal() {
  generateImageError.value = null
  generatedImageBlob.value = null
  generatedImageUrl.value = null
  showGeneratedImageModal.value = true
}

async function startImageGeneration() {
  generatingImage.value = true
  generateRecipeImage()
}

async function generateRecipeImage() {
  if (!recipe.value?.name) return
  generateImageError.value = null
  try {
    const blob = await fileService.GenerateRecipeImage({ name: recipe.value.name })
    generatedImageBlob.value = blob
    if (generatedImageUrl.value) {
      URL.revokeObjectURL(generatedImageUrl.value)
    }
    generatedImageUrl.value = URL.createObjectURL(blob)
  } catch (error) {
    generateImageError.value = error instanceof Error ? error.message : String(error)
    generatedImageBlob.value = null
    generatedImageUrl.value = null
  } finally {
    generatingImage.value = false
  }
}

function closeGenerateImageModal() {
  showGeneratedImageModal.value = false
  if (generatedImageUrl.value) {
    URL.revokeObjectURL(generatedImageUrl.value)
    generatedImageUrl.value = null
  }
  generatedImageBlob.value = null
  generateImageError.value = null
}

// *** Update Image Modal ***

const updatingImage = ref(false)

async function updateRecipeImage() {
  if (!recipe.value?.name || !generatedImageBlob.value) return
  updatingImage.value = true
  try {
    await fileService.UploadRecipeImage({ name: recipe.value.name, file: new File([generatedImageBlob.value], 'generated-image.png', { type: generatedImageBlob.value.type || 'image/png' }) })
    // Reload the recipe to show the new image
    await recipesStore.loadRecipe(recipe.value.name)
    closeGenerateImageModal()
  } catch (error) {
    generateImageError.value = error instanceof Error ? error.message : String(error)
  } finally {
    updatingImage.value = false
  }
}

// *** Recipe Sharing ***

const allowPermissionOptions: PermissionLevel[] = [
  'PERMISSION_LEVEL_READ',
  'PERMISSION_LEVEL_WRITE',
  'PERMISSION_LEVEL_ADMIN',
]
const showShareDialog = ref(false)
const currentAccesses = ref<Access[]>([])
const sharing = ref(false)
const updatingPermission = ref<Record<string, boolean>>({})
const unsharing = ref<Record<string, boolean>>({})


// Fetch recipients when share dialog is opened
watch(showShareDialog, (isOpen) => {
  if (isOpen && recipe.value && hasWritePermission(recipe.value.recipeAccess?.permissionLevel)) {
    fetchRecipeRecipients()
  }
})

// Function to fetch recipe recipients
async function fetchRecipeRecipients() {
  if (!recipe.value?.name) return

  try {
    const request: ListAccessesRequest = {
      parent: trimmedRecipeName.value,
      filter: undefined,
      pageSize: undefined,
      pageToken: undefined
    }

    const response = await recipeAccessService.ListAccesses(request)

    if (response.accesses) {
      currentAccesses.value = response.accesses.filter(access => {
        // Filter out the current user's own access to avoid showing it in the shares list
        return access.recipient?.user?.name !== authStore.user.name
      }).map(access => ({
        name: access.name || '',
        level: access.level || 'PERMISSION_LEVEL_READ',
        state: access.state || 'ACCESS_STATE_PENDING',
        recipient: access.recipient,
        requester: access.requester || undefined,
        acceptTarget: access.acceptTarget || 'ACCEPT_TARGET_UNSPECIFIED',
      }))
    }
  } catch (error) {
    console.error('Error fetching recipe recipients:', error)
  }
}

async function updatePermission({ access, newLevel }: { access: Access, newLevel: PermissionLevel }) {
  if (access.level === newLevel) return
  if (!access.name) return
  updatingPermission.value[access.name] = true
  try {
    await recipeAccessService.UpdateAccess({
      access: {
        name: access.name,
        level: newLevel,
        state: undefined,
        recipient: undefined,
        requester: undefined,
        acceptTarget: undefined,
      },
      updateMask: 'level',
    })
    // Update local state
    access.level = newLevel
  } catch (error) {
    console.error('Error updating permission:', error)
  } finally {
    updatingPermission.value[access.name] = false
  }
}

async function unshareRecipe(accessName: string) {
  unsharing.value[accessName] = true
  try {
    const request: DeleteAccessRequest = {
      name: accessName
    }

    await recipeAccessService.DeleteAccess(request)
    await fetchRecipeRecipients()
  } catch (error) {
    console.error('Error removing share:', error)
  } finally {
    unsharing.value[accessName] = false
  }
}

async function shareWithUser({ userName, permission }: { userName: string, permission: PermissionLevel }) {
  if (!userName) return
  if (!recipe.value?.name) return

  sharing.value = true
  try {
    const access: Access = {
      recipient: {
        user: {
          name: userName,
          username: undefined,
          givenName: undefined,
          familyName: undefined,
        }
      },
      level: permission,
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: trimmedRecipeName.value,
      access
    }

    await recipeAccessService.CreateAccess(request)
    await fetchRecipeRecipients()
  } catch (error) {
    console.error('Error sharing recipe:', error)
    // You might want to show an error notification here
  } finally {
    sharing.value = false
  }
}

async function shareWithCircle({ circleName, permission }: { circleName: string, permission: PermissionLevel }) {
  if (!circleName) return
  if (!recipe.value?.name) return

  sharing.value = true
  try {
    const access: Access = {
      recipient: {
        circle: {
          name: circleName,
          title: undefined,
          handle: undefined,
        }
      },
      level: permission,
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: recipe.value.name,
      access
    }

    await recipeAccessService.CreateAccess(request)
    await fetchRecipeRecipients()
  } catch (error) {
    console.error('Error sharing recipe:', error)
    // You might want to show an error notification here
  } finally {
    sharing.value = false
  }
}

const requestingAccess = ref(false)

// Event creation state
const showCreateEventDialog = ref(false)
const writableCalendars = ref<Calendar[]>([])

// Load writable calendars when component mounts
onMounted(async () => {
  // First check URL hash
  const currentHash = route.hash
  if (currentHash in hashToTab) {
    tab.value = hashToTab[currentHash]
  }

  await recipesStore.loadRecipe(trimmedRecipeName.value)

  // Load writable calendars for event creation
  await loadWritableCalendars()
})

async function loadWritableCalendars() {
  if (!authStore.user?.name) return

  try {
    await calendarsStore.loadMyCalendars(authStore.user.name)
    // Get the calendars from the store after loading
    const { myCalendars } = storeToRefs(calendarsStore)
    writableCalendars.value = myCalendars.value.filter((calendar: Calendar) => {
      // Check if user has write permission to this specific calendar
      const calendarAccess = calendar.calendarAccess
      if (!calendarAccess?.permissionLevel) return false

      return calendarAccess.permissionLevel === 'PERMISSION_LEVEL_WRITE' ||
        calendarAccess.permissionLevel === 'PERMISSION_LEVEL_ADMIN'
    })
  } catch (error) {
    console.error('Error loading writable calendars:', error)
    writableCalendars.value = []
  }
}

function openCreateEventDialog() {
  showCreateEventDialog.value = true
}

async function handleEventCreated(event: Event) {
  if (!recipe.value?.name || !event?.name) return

  try {
    // Create event_recipe relationship
    await eventRecipeService.CreateEventRecipe({
      parent: event.name, // Use the full event name as parent
      eventRecipe: {
        name: undefined, // Will be set by server
        recipe: recipe.value.name, // Use the full recipe name
        createTime: undefined, // Will be set by server
      }
    })

    alertsStore.addAlert('Event created and linked to recipe successfully!', 'success')
  } catch (error) {
    console.error('Error creating event-recipe relationship:', error)
    alertsStore.addAlert(error instanceof Error ? error.message : 'Failed to link event to recipe', 'error')
  }
}

async function handleRequestAccess() {
  if (!recipe.value?.name) return
  if (!authStore.user?.name) return

  requestingAccess.value = true
  try {
    const access: Access = {
      recipient: {
        user: {
          name: authStore.user?.name,
          username: undefined,
          givenName: undefined,
          familyName: undefined,
        }
      },
      level: 'PERMISSION_LEVEL_READ',
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: recipe.value.name,
      access
    }

    await recipeAccessService.CreateAccess(request)
    await recipesStore.loadRecipe(recipe.value.name)
    alertsStore.addAlert('Access request sent.', 'info')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    requestingAccess.value = false
  }
}


// *** Favorite Recipe ***

async function toggleFavorite() {
  if (!recipe.value?.name) return
  try {
    if (recipe.value.favorited) {
      await recipeService.UnfavoriteRecipe({ name: recipe.value.name })
    } else {
      await recipeService.FavoriteRecipe({ name: recipe.value.name })
    }
    recipe.value.favorited = !recipe.value.favorited
  } catch (error) {
    console.error('Error toggling favorite:', error)
  }
}

// *** Utility Functions ***

function isUrl(str: string): boolean {
  return /^https?:\/\//.test(str);
}

function parseDuration(duration: string): number {
  if (!duration) return 0;

  if (duration.endsWith('s')) {
    return parseInt(duration.slice(0, -1)) / 60;
  }
  return 0;
}

</script>

<style scoped>
.image-container {
  position: relative;
}

.generate-image-btn {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.ingredient-item {
  min-height: auto;
}

.favorite-heart {
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 4px !important;
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.6));
  /* background-color: rgba(255, 255, 255, 0.9); */
  border-radius: 50%;
  padding: 4px;
  transition: all 0.2s ease-in-out;
  color: red;
}
</style>
