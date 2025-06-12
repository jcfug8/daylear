<template>
  <v-container v-if="recipe">
    <v-app-bar>
      <v-tabs style="width: 100%" v-model="tab" center-active show-arrows fixed-tabs>
        <v-tab value="general" text="General"></v-tab>
        <v-tab value="ingredients" text="Ingredients"></v-tab>
        <v-tab value="directions" text="Directions"></v-tab>
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
              <v-img class="mt-1" style="background-color: lightgray" :src="recipe.imageUri" cover></v-img>
            </v-col>
            <v-spacer></v-spacer>
          </v-row>
        </v-container>
      </v-tabs-window-item>
      <v-tabs-window-item value="ingredients">
        <v-container max-width="600">
          <v-card class="my-4 mx-1 pa-2" v-for="(ingredientGroup, i) in recipe.ingredientGroups" :key="i">
            <v-card-title v-if="ingredientGroup.title">{{ ingredientGroup.title }}</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item slim prepend-icon="mdi-circle-small" v-for="(ingredient, j) in ingredientGroup.ingredients"
                  :key="j">
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
                <v-list-item slim prepend-icon="mdi-circle-small" v-for="(step, n) in direction.steps">
                  <div class="font-weight-bold">Step {{ n + 1 }}</div>
                  {{ step }}
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-container>
      </v-tabs-window-item>
    </v-tabs-window>


    <!-- Speed Dial -->
    <v-fab location="bottom right" app color="primary"  icon @click="speedDialOpen = !speedDialOpen">
      <v-icon>mdi-dots-horizontal</v-icon>
      <v-speed-dial location="top" v-model="speedDialOpen" transition="slide-y-reverse-transition" activator="parent">
        <v-btn key="edit" v-if="permissionLevel === 'RESOURCE_PERMISSION_WRITE' && recipe && recipe.name" icon="mdi-pencil"
        @click="router.push({ name: 'recipeEdit', params: { recipeId: recipe.name } })" color="primary"></v-btn>

        <v-btn key="share" v-if="permissionLevel === 'RESOURCE_PERMISSION_WRITE'" icon="mdi-share-variant"
          @click="showShareDialog = true" color="primary"></v-btn>
  
        <v-btn key="remove-access" icon="mdi-link-variant-off" @click="showRemoveAccessDialog = true" color="warning"></v-btn>
  
        <v-btn key="delete" v-if="permissionLevel === 'RESOURCE_PERMISSION_WRITE' && recipe && recipe.name" icon="mdi-delete"
          @click="showDeleteDialog = true" color="error"></v-btn>
      </v-speed-dial>
    </v-fab>
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
  <v-dialog v-model="showShareDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Share Recipe
      </v-card-title>
      <v-card-text>
        <v-tabs v-model="shareTab" class="mb-4">
          <v-tab value="users">Share with User</v-tab>
          <v-tab value="circles">Share with Circle</v-tab>
        </v-tabs>

        <v-window v-model="shareTab">
          <v-window-item value="users">
            <v-text-field v-model="usernameInput" label="Enter Username" :rules="[validateUsername]"
              :append-inner-icon="getUsernameIcon" :color="getUsernameColor" :loading="isLoadingUsername"
              @update:model-value="handleUsernameInput"></v-text-field>
          </v-window-item>

          <v-window-item value="circles">
            <v-text-field v-model="circleInput" label="Enter Circle Name" :rules="[validateCircle]"
              :append-inner-icon="getCircleIcon" :color="getCircleColor" :loading="isLoadingCircle"
              @update:model-value="handleCircleInput"></v-text-field>
          </v-window-item>
        </v-window>

        <v-divider class="my-4"></v-divider>

        <div v-if="currentShares.length > 0">
          <div class="text-subtitle-1 mb-2">Current Shares</div>
          <v-list>
            <v-list-item v-for="share in currentShares" :key="share.id" :title="share.name" :subtitle="share.type">
              <template #append>
                <v-btn icon="mdi-delete" variant="text" @click="removeShare(share.id)"></v-btn>
              </template>
            </v-list-item>
          </v-list>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showShareDialog = false">
          Close
        </v-btn>
        <v-btn color="primary" @click="shareRecipe" :loading="sharing" :disabled="!isValidUsername && !isValidCircle">
          Share
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
</template>

<script setup lang="ts">
import type { Recipe_MeasurementType } from '@/genapi/api/meals/recipe/v1alpha1'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRecipesStore } from '@/stores/recipes'
import { useRecipeFormStore } from '@/stores/recipeForm'
import { storeToRefs } from 'pinia'
import { onMounted, onBeforeUnmount, watch, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { publicUserService, publicCircleService, recipeService, recipeRecipientsService } from '@/api/api'
import type { PublicUser, ListPublicUsersRequest } from '@/genapi/api/users/user/v1alpha1'
import type { PublicCircle, ListPublicCirclesRequest } from '@/genapi/api/circles/circle/v1alpha1'
import { useAuthStore } from '@/stores/auth'
import type { PermissionLevel } from '@/genapi/api/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const recipesStore = useRecipesStore()
const recipeFormStore = useRecipeFormStore()
const { recipe } = storeToRefs(recipesStore)
const { activeTab } = storeToRefs(recipeFormStore)
var permissionLevel = ref<PermissionLevel | undefined>(undefined)

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

// Share dialog state
const showShareDialog = ref(false)
const shareTab = ref('users')
const usernameInput = ref('')
const circleInput = ref('')
const selectedUser = ref<PublicUser | null>(null)
const selectedCircle = ref<PublicCircle | null>(null)
const sharing = ref(false)

// Validation states
const isValidUsername = ref(false)
const isValidCircle = ref(false)
const isLoadingUsername = ref(false)
const isLoadingCircle = ref(false)

// Current shares state
const currentShares = ref<Array<{ id: string; name: string; type: 'user' | 'circle' }>>([])

// Debounce timers
let usernameDebounceTimer: number | null = null
let circleDebounceTimer: number | null = null

// Computed properties for icons and colors
const getUsernameIcon = computed(() => {
  if (isLoadingUsername.value) return 'mdi-loading'
  if (!usernameInput.value) return undefined
  return isValidUsername.value ? 'mdi-check-circle' : 'mdi-close-circle'
})

const getCircleIcon = computed(() => {
  if (isLoadingCircle.value) return 'mdi-loading'
  if (!circleInput.value) return undefined
  return isValidCircle.value ? 'mdi-check-circle' : 'mdi-close-circle'
})

const getUsernameColor = computed(() => {
  if (isLoadingUsername.value) return undefined
  if (!usernameInput.value) return undefined
  return isValidUsername.value ? 'success' : 'error'
})

const getCircleColor = computed(() => {
  if (isLoadingCircle.value) return undefined
  if (!circleInput.value) return undefined
  return isValidCircle.value ? 'success' : 'error'
})

// Debounced API calls
async function checkUsername(username: string) {
  if (!username) {
    isValidUsername.value = false
    selectedUser.value = null
    return
  }

  isLoadingUsername.value = true
  try {
    const request: ListPublicUsersRequest = {
      filter: `username = "${username}"`,
      pageSize: 1,
      pageToken: undefined
    }
    const response = await publicUserService.ListPublicUsers(request)

    if (response.publicUsers?.length === 1 && response.publicUsers[0].name !== authStore.activeAccount?.publicName) {
      selectedUser.value = response.publicUsers[0]
      isValidUsername.value = true
    } else {
      selectedUser.value = null
      isValidUsername.value = false
    }
  } catch (error) {
    console.error('Error checking username:', error)
    selectedUser.value = null
    isValidUsername.value = false
  } finally {
    isLoadingUsername.value = false
  }
}

async function checkCircle(circleName: string) {
  if (!circleName) {
    isValidCircle.value = false
    selectedCircle.value = null
    return
  }

  isLoadingCircle.value = true
  try {
    const request: ListPublicCirclesRequest = {
      filter: `title = "${circleName}"`,
      pageSize: 1,
      pageToken: undefined
    }
    const response = await publicCircleService.ListPublicCircles(request)

    if (response.publicCircles?.length === 1 && response.publicCircles[0].name !== authStore.activeAccount?.publicName) {
      selectedCircle.value = response.publicCircles[0]
      isValidCircle.value = true
    } else {
      selectedCircle.value = null
      isValidCircle.value = false
    }
  } catch (error) {
    console.error('Error checking circle:', error)
    selectedCircle.value = null
    isValidCircle.value = false
  } finally {
    isLoadingCircle.value = false
  }
}

function handleUsernameInput(value: string) {
  if (usernameDebounceTimer) {
    clearTimeout(usernameDebounceTimer)
  }
  usernameDebounceTimer = window.setTimeout(() => {
    checkUsername(value)
  }, 300)
}

function handleCircleInput(value: string) {
  if (circleDebounceTimer) {
    clearTimeout(circleDebounceTimer)
  }
  circleDebounceTimer = window.setTimeout(() => {
    checkCircle(value)
  }, 300)
}

// Clean up timers when component is unmounted
onBeforeUnmount(() => {
  if (usernameDebounceTimer) {
    clearTimeout(usernameDebounceTimer)
  }
  if (circleDebounceTimer) {
    clearTimeout(circleDebounceTimer)
  }
})

function validateUsername(value: string): boolean | string {
  if (!value) return true
  return true // Validation is now handled by the API call
}

function validateCircle(value: string): boolean | string {
  if (!value) return true
  return true // Validation is now handled by the API call
}

async function shareRecipe() {
  if (!selectedUser.value && !selectedCircle.value) return

  sharing.value = true
  try {
    const request = {
      name: recipe.value?.name || '',
      recipients: [selectedUser.value?.name || selectedCircle.value?.name || ''],
      permission: 'RESOURCE_PERMISSION_READ' as const
    }
    await recipeService.ShareRecipe(request)
    // Refresh the recipients list
    await fetchRecipeRecipients()
    // Reset selections and inputs after sharing
    selectedUser.value = null
    selectedCircle.value = null
    usernameInput.value = ''
    circleInput.value = ''
    isValidUsername.value = false
    isValidCircle.value = false
    showShareDialog.value = false
  } catch (error) {
    console.error('Error sharing recipe:', error)
    // You might want to show an error notification here
  } finally {
    sharing.value = false
  }
}

async function removeShare(shareId: string) {
  try {
    const request = {
      name: recipe.value?.name || '',
      recipients: [shareId]
    }
    await recipeService.UnshareRecipe(request)
    // Remove from local state after successful API call
    currentShares.value = currentShares.value.filter(share => share.id !== shareId)
  } catch (error) {
    console.error('Error removing share:', error)
    // You might want to show an error notification here
  }
}

// Function to fetch recipe recipients
async function fetchRecipeRecipients() {
  if (!recipe.value?.name) return

  try {
    const response = await recipeRecipientsService.GetRecipeRecipients({
      name: recipe.value.name
    })

    if (response.recipients) {
      currentShares.value = response.recipients.filter(recipient => {
        if (recipient.name === authStore.activeAccount?.name) {
          permissionLevel.value = recipient.permission
          return false
        }
        return true
      }).map(recipient => ({
        id: recipient.name || '',
        name: recipient.title || '',
        type: recipient.name?.includes('circles') ? 'circle' : 'user'
      }))
    }
  } catch (error) {
    console.error('Error fetching recipe recipients:', error)
  }
}

// Fetch recipients when recipe is loaded
watch(recipe, (newRecipe) => {
  if (newRecipe) {
    fetchRecipeRecipients()
  }
}, { immediate: true })

// Remove access dialog state
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

async function handleRemoveAccess() {
  if (!recipe.value?.name) return

  removingAccess.value = true
  try {
    const request = {
      name: recipe.value.name,
      recipients: [authStore.activeAccount?.name || '']
    }
    await recipeService.UnshareRecipe(request)
    router.push({ name: 'recipes' })
  } catch (error) {
    console.error('Error removing access:', error)
    alert(error instanceof Error ? error.message : String(error))
  } finally {
    removingAccess.value = false
    showRemoveAccessDialog.value = false
  }
}

// Speed dial state
const speedDialOpen = ref(false)
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
</script>
