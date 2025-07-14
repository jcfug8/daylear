<template>
  <v-container v-if="recipe">
    <v-app-bar>
      <v-tabs style="width: 100%" v-model="tab" center-active show-arrows fixed-tabs>
        <v-tab value="general">
          <v-icon left>mdi-information-outline</v-icon>
        </v-tab>
        <v-tab value="ingredients">
          <v-icon left>mdi-food-apple-outline</v-icon>
        </v-tab>
        <v-tab value="directions">
          <v-icon left>mdi-format-list-numbered</v-icon>
        </v-tab>
      </v-tabs>
      <template #append>
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
              <div class="image-container">
                <v-img 
                  v-if="recipe.imageUri" 
                  class="mt-1" 
                  style="background-color: lightgray" 
                  :src="recipe.imageUri" 
                  cover
                  height="300"
                >
                  <template #placeholder>
                    <v-row class="fill-height ma-0" align="center" justify="center">
                      <v-progress-circular indeterminate color="grey lighten-5"></v-progress-circular>
                    </v-row>
                  </template>
                </v-img>
                <div 
                  v-else 
                  class="mt-1 d-flex align-center justify-center"
                  style="background-color: lightgray; height: 300px; border-radius: 4px;"
                >
                  <div class="text-center">
                    <v-icon size="64" color="grey-darken-1">mdi-image-outline</v-icon>
                    <div class="text-grey-darken-1 mt-2">No image available</div>
                  </div>
                </div>
                <v-btn
                  v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel) && !recipe.imageUri"
                  class="generate-image-btn"
                  color="warning"
                  @click="openGenerateImageModal"
                  style="position: absolute; bottom: 8px; right: 8px; z-index: 2;"
                  :loading="generatingImage"
                  :disabled="generatingImage"
                  title="Generate Image"
                >
                  <v-icon>mdi-image-auto-adjust</v-icon>
                  Generate Image
                </v-btn>
              </div>
            </v-col>
            <v-spacer></v-spacer>
          </v-row>

          <!-- New Recipe Fields -->
          <v-row>
            <v-col cols="12">
              <div v-if="recipe.citation">
                <span v-if="isUrl(recipe.citation)">
                  <a :href="recipe.citation" target="_blank" rel="noopener">{{ recipe.citation }}</a>
                </span>
                <span v-else>
                  {{ recipe.citation }}
                </span>
              </div>
              <div v-if="recipe.prepDuration && recipe.prepDuration !== '0s'">
                <strong>Prep Time:</strong> {{ parseDuration(recipe.prepDuration ? recipe.prepDuration : "0") }}
              </div>
              <div v-if="recipe.cookDuration && recipe.cookDuration !== '0s'">
                <strong>Cook Time:</strong> {{ parseDuration(recipe.cookDuration ? recipe.cookDuration : "0") }}
              </div>
              <div v-if="recipe.totalDuration && recipe.totalDuration !== '0s'">
                <strong>Total Time:</strong> {{ parseDuration(recipe.totalDuration ? recipe.totalDuration : "0") }}
              </div>
              <div v-if="recipe.cookingMethod">
                <strong>Cooking Method:</strong> {{ recipe.cookingMethod }}
              </div>
              <div v-if="recipe.categories && recipe.categories.length">
                <strong>Categories:</strong> {{ recipe.categories.join(', ') }}
              </div>
              <div v-if="recipe.yieldAmount">
                <strong>Yield:</strong> {{ recipe.yieldAmount }}
              </div>
              <div v-if="recipe.cuisines && recipe.cuisines.length">
                <strong>Cuisines:</strong> {{ recipe.cuisines.join(', ') }}
              </div>
              <!-- <div v-if="createTimeString">
                <strong>Created:</strong> {{ formatDate(createTimeString) }}
              </div>
              <div v-if="updateTimeString">
                <strong>Updated:</strong> {{ formatDate(updateTimeString) }}
              </div> -->
            </v-col>
          </v-row>

          <!-- Visibility Section -->
          <v-row>
            <v-col cols="12">
              <div class="mt-4">
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
        </v-container>
      </v-tabs-window-item>
      <v-tabs-window-item value="ingredients">
        <v-container max-width="600">
          <v-card class="my-1" v-for="(ingredientGroup, i) in recipe.ingredientGroups" :key="i">
            <v-card-title v-if="ingredientGroup.title">{{ ingredientGroup.title }}</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item class="ingredient-item" v-for="(ingredient, j) in ingredientGroup.ingredients" :key="j">
                  <div class="d-flex align-center" style="width: 100%;">
                    <v-checkbox v-model="checkedIngredients[i][j]" class="mr-2" hide-details density="compact" style="margin-bottom: 0;" />
                    <span :style="checkedIngredients[i][j] ? 'text-decoration: line-through; color: #888;' : ''">
                      <strong>
                        <span v-if="ingredient.measurementAmount">{{ isFranctional(ingredient.measurementType) ? toFraction(ingredient.measurementAmount ?? 0) : ingredient.measurementAmount }}</span>
                        {{ measurementTypeLabel(ingredient.measurementType, ingredient.measurementAmount) }}
                        <template v-if="ingredient.measurementConjunction && ingredient.secondMeasurementAmount && ingredient.secondMeasurementType">
                          {{ renderConjunction(ingredient.measurementConjunction) }}
                          <span v-if="ingredient.secondMeasurementAmount">{{ isFranctional(ingredient.secondMeasurementType) ? toFraction(ingredient.secondMeasurementAmount ?? 0) : ingredient.secondMeasurementAmount }}</span>
                          {{ measurementTypeLabel(ingredient.secondMeasurementType, ingredient.secondMeasurementAmount) }}
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
                <v-list-item  v-for="(step, n) in direction.steps" :key="n">
                  <div class="d-flex align-start" style="width: 100%;">
                    <v-checkbox
                      v-model="checkedDirections[i][n]"
                      class="mr-2"
                      hide-details
                      density="compact"
                      style="margin-bottom: 0; align-self: flex-start;"
                    />
                    <span
                      :style="checkedDirections[i][n] ? 'text-decoration: line-through; color: #888;' : ''"
                      style="margin-left: 4px; flex: 1 1 0; word-break: break-word;"
                    >
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


    <!-- Speed Dial -->
    <v-fab location="bottom right" app color="primary"  icon @click="speedDialOpen = !speedDialOpen">
      <v-icon>mdi-dots-vertical</v-icon>
      <v-speed-dial location="top" v-model="speedDialOpen" transition="slide-y-reverse-transition" activator="parent">
        <v-btn key="edit" v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel)"
        :to="getRecipeEditRoute()" color="primary"><v-icon>mdi-pencil</v-icon>Edit</v-btn>

        <v-btn key="share" v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel) && recipe.visibility !== 'VISIBILITY_LEVEL_HIDDEN'"
          @click="showShareDialog = true" color="primary"><v-icon>mdi-share-variant</v-icon>Share</v-btn>
  
        <v-btn key="remove-access" v-if="!hasAdminPermission(recipe.recipeAccess?.permissionLevel)" @click="showRemoveAccessDialog = true" color="warning"><v-icon>mdi-link-variant-off</v-icon>Remove</v-btn>
  
        <v-btn key="delete" v-if="hasAdminPermission(recipe.recipeAccess?.permissionLevel)"
          @click="showDeleteDialog = true" color="error"><v-icon>mdi-delete</v-icon>Delete</v-btn>
      </v-speed-dial>
    </v-fab>

    <v-btn
      v-if="recipe.recipeAccess?.state === 'ACCESS_STATE_PENDING'"
      color="success"
      class="mb-2"
      block
      :loading="acceptingRecipe"
      @click="acceptRecipe"
    >
      Accept Recipe
    </v-btn>
    <v-btn
      v-if="recipe.recipeAccess?.state === 'ACCESS_STATE_PENDING'"
      color="error"
      class="mb-4"
      block
      :loading="decliningRecipe"
      @click="declineRecipe"
    >
      Decline
    </v-btn>
  </v-container>
  <!-- Remove Access Dialog -->
  <v-dialog v-model="showRemoveAccessDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Remove Access
      </v-card-title>
      <v-card-text>
        Are you sure you want to remove your access to this recipe? You will no longer be able to view it.
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

  <!-- Share Dialog -->
  <ShareDialog
    v-model="showShareDialog"
    title="Share Recipe"
    :allowCircleShare="true"
    :currentShares="currentShares"
    :sharing="sharing"
    :sharePermissionLoading="updatingPermission"
    :hasWritePermission="hasWritePermission"
    @share-user="shareWithUser"
    @share-circle="shareWithCircle"
    @remove-share="unshareCircle"
    @permission-change="updatePermission"
  />

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
          <v-btn v-if="generatedImageUrl" color="success" :disabled="!generatedImageBlob || updatingImage || generatingImage" @click="updateRecipeImage">
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
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRecipesStore } from '@/stores/recipes'
import { useRecipeFormStore } from '@/stores/recipeForm'
import { useCirclesStore } from '@/stores/circles'
import { storeToRefs } from 'pinia'
import { onMounted, onBeforeUnmount, watch, computed, ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { recipeService, recipeAccessService, fileService } from '@/api/api'
import { useAuthStore } from '@/stores/auth'
import { hasWritePermission, hasAdminPermission } from '@/utils/permissions'
import ShareDialog from '@/components/common/ShareDialog.vue'

const authStore = useAuthStore()
const breadcrumbStore = useBreadcrumbStore()
const recipesStore = useRecipesStore()
const circlesStore = useCirclesStore()
const recipeFormStore = useRecipeFormStore()
const route = useRoute()
const router = useRouter()

const { recipe } = storeToRefs(recipesStore)
const { circle } = storeToRefs(circlesStore)
const speedDialOpen = ref(false)

// Add local state for checked ingredients and directions
const checkedIngredients = reactive<Array<Array<boolean>>>([])
const checkedDirections = reactive<Array<Array<boolean>>>([])

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

function getRecipeEditRoute() {
  if (route.params.circleId) {
    return { name: 'circleRecipeEdit', params: { circleId: route.params.circleId, recipeId: recipe.value?.name } }
  }
  return { name: 'recipeEdit', params: { recipeId: recipe.value?.name } }
}

// *** Breadcrumbs ***

onMounted(async () => {
  // First check URL hash
  const currentHash = route.hash
  if (currentHash in hashToTab) {
    tab.value = hashToTab[currentHash]
  }

  // Load the recipe based on the route parameter
  const recipeId = route.params.recipeId as string
  await recipesStore.loadRecipe(recipeId)

  let firstCrumbs

  if (route.params.circleId) {
    await circlesStore.loadCircle(route.params.circleId as string)
    firstCrumbs = [
      { title: 'Circles', to: { name: 'circles' } },
      { title: circle.value?.title || '', to: { name: 'circle', params: { circleId: route.params.circleId } } },
    ]
  } else {
    firstCrumbs = [{ title: 'Recipes', to: { name: 'recipes' } }]
  }

  breadcrumbStore.setBreadcrumbs([
    ...firstCrumbs,
    {
      title: recipe.value?.title || '',
      to: { name: 'recipe', params: { recipeId: recipe.value?.name } },
    },
  ])
})

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

// *** Remove Access ***
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

async function handleRemoveAccess() {
  if (!recipe.value?.recipeAccess?.name || !authStore.activeAccount?.name) return

  removingAccess.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: recipe.value.recipeAccess.name
    }
    
    await recipeAccessService.DeleteAccess(deleteRequest)
    router.push(route.params.circleId ? { name: 'circle', params: { circleId: route.params.circleId } } : { name: 'recipes' })
  } catch (error) {
    console.error('Error removing access:', error)
    alert(error instanceof Error ? error.message : String(error))
  } finally {
    removingAccess.value = false
    showRemoveAccessDialog.value = false
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
    console.error('Error deleting recipe:', error)
    alert(error instanceof Error ? error.message : String(error))
  } finally {
    deleting.value = false
    showDeleteDialog.value = false
  }
}

// *** Accept/Decline Recipe ***

const acceptingRecipe = ref(false)
const decliningRecipe = ref(false)

async function acceptRecipe() {
  if (!recipe.value?.recipeAccess?.name || !authStore.user?.name) return
  acceptingRecipe.value = true
  try {
    await recipesStore.acceptRecipe(recipe.value.recipeAccess.name)
    recipesStore.loadRecipe(recipe.value?.name ?? '')
  } catch (error) {
    // Optionally show a notification
  } finally {
    acceptingRecipe.value = false
  }
}

async function declineRecipe() {
  if (!recipe.value?.recipeAccess?.name) return
  decliningRecipe.value = true
  try {
    await recipesStore.deleteRecipeAccess(recipe.value.recipeAccess.name)
    router.push({ name: 'recipes' })
  } catch (error) {
    // Optionally show a notification
  } finally {
    decliningRecipe.value = false
  }
}

// *** Measurements and Ingredients ***

function renderConjunction(conjunction: Recipe_Ingredient_MeasurementConjunction | undefined): string {
  switch (conjunction) {
    case 'MEASUREMENT_CONJUNCTION_AND':
      return '+';
    case 'MEASUREMENT_CONJUNCTION_TO':
      return '-';
    case 'MEASUREMENT_CONJUNCTION_OR':
      return 'or';
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
  let frac = amount - whole
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
  if (result === '') result = '0'
  return result
}

// *** Generate Image ***

const showGeneratedImageModal = ref(false)
const generatingImage = ref(false)
const generatedImageBlob = ref<Blob|null>(null)
const generatedImageUrl = ref<string|null>(null)
const generateImageError = ref<string|null>(null)

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

const showShareDialog = ref(false)
const currentShares = ref<Access[]>([])
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
      parent: recipe.value.name,
      filter: undefined,
      pageSize: undefined,
      pageToken: undefined
    }

    const response = await recipeAccessService.ListAccesses(request)

    if (response.accesses) {
      currentShares.value = response.accesses.filter(access => {
        // Filter out the current user's own access to avoid showing it in the shares list
        const isCurrentUser = (access.recipient?.user && access.recipient.user.name === authStore.activeAccount?.name) ||
                             (access.recipient?.circle && access.recipient.circle.name === authStore.activeAccount?.name)
        return !isCurrentUser
      }).map(access => ({
        name: access.name || '',
        level: access.level || 'PERMISSION_LEVEL_READ',
        state: access.state || 'ACCESS_STATE_PENDING',
        recipient: access.recipient,
        requester: access.requester || undefined,
      }))
    }
  } catch (error) {
    console.error('Error fetching recipe recipients:', error)
  }
}

async function updatePermission({ share, newLevel }: { share: Access, newLevel: PermissionLevel }) {
  if (share.level === newLevel) return
  if (!share.name) return
  updatingPermission.value[share.name] = true
  try {
    await recipeAccessService.UpdateAccess({
      access: {
        name: share.name,
        level: newLevel,
        state: undefined,
        recipient: undefined,
        requester: undefined,
      },
      updateMask: 'level',
    })
    // Update local state
    share.level = newLevel
  } catch (error) {
    console.error('Error updating permission:', error)
  } finally {
    updatingPermission.value[share.name] = false
  }
}

async function unshareCircle(accessName: string) {
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

// *** Utility Functions ***

function isUrl(str: string): boolean {
  return /^https?:\/\//.test(str);
}

function parseDuration(duration: string): number {
  if (!duration) return 0;

  if (duration.endsWith('s')) {
    return parseInt(duration.slice(0, -1))/60;
  }
  return 0;
}

</script>

<style scoped>
.image-container {
  position: relative;
}
.generate-image-btn {
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}
.ingredient-item {
  min-height: auto;
}
</style>
