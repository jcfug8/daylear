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
              <div class="image-container">
                <v-img 
                  v-if="recipe.imageUri" 
                  class="mt-1" 
                  style="background-color: lightgray" 
                  :src="recipe.imageUri" 
                  cover
                  height="300"
                ></v-img>
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
              </div>
            </v-col>
            <v-spacer></v-spacer>
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
          <v-card class="my-4 mx-1 pa-2" v-for="(ingredientGroup, i) in recipe.ingredientGroups" :key="i">
            <v-card-title v-if="ingredientGroup.title">{{ ingredientGroup.title }}</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item slim prepend-icon="mdi-circle-small" v-for="(ingredient, j) in ingredientGroup.ingredients"
                  :key="j">
                  <span v-if="ingredient.measurementAmount">{{ isFranctional(ingredient.measurementType) ? toFraction(ingredient.measurementAmount ?? 0) : ingredient.measurementAmount }}</span>
                  {{ measurementTypeLabel(ingredient.measurementType, ingredient.measurementAmount) }}
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
      <v-icon>mdi-dots-vertical</v-icon>
      <v-speed-dial location="top" v-model="speedDialOpen" transition="slide-y-reverse-transition" activator="parent">
        <v-btn key="edit" v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel)" icon="mdi-pencil"
        @click="router.push({ name: 'recipeEdit', params: { recipeId: recipe.name } })" color="primary"></v-btn>

        <v-btn key="share" v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel) && recipe.visibility !== 'VISIBILITY_LEVEL_HIDDEN'" icon="mdi-share-variant"
          @click="showShareDialog = true" color="primary"></v-btn>
  
        <v-btn key="remove-access" icon="mdi-link-variant-off" @click="showRemoveAccessDialog = true" color="warning"></v-btn>
  
        <v-btn key="delete" v-if="hasWritePermission(recipe.recipeAccess?.permissionLevel)" icon="mdi-delete"
          @click="showDeleteDialog = true" color="error"></v-btn>
      </v-speed-dial>
    </v-fab>

    <v-btn
      v-if="recipe.recipeAccess?.state === 'ACCESS_STATE_PENDING'"
      color="success"
      class="mb-4"
      @click="acceptRecipe"
      block
    >
      Accept Recipe
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
              :prepend-inner-icon="getUsernameIcon" :color="getUsernameColor" :loading="isLoadingUsername"
              @update:model-value="handleUsernameInput"></v-text-field>
            <v-select
              v-model="selectedPermission"
              :items="permissionOptions"
              label="Permission Level"
              class="mt-2"
            ></v-select>
            <v-btn block color="primary" @click="shareRecipe" :loading="sharing" :disabled="!isValidUsername" class="mt-2">
              Share with User
            </v-btn>
          </v-window-item>

          <v-window-item value="circles">
            <v-text-field v-model="circleInput" label="Enter Circle Name" :rules="[validateCircle]"
              :prepend-inner-icon="getCircleIcon" :color="getCircleColor" :loading="isLoadingCircle"
              @update:model-value="handleCircleInput"></v-text-field>
            <v-select
              v-model="selectedPermission"
              :items="permissionOptions"
              label="Permission Level"
              class="mt-2"
            ></v-select>
            <v-btn block color="primary" @click="shareRecipe" :loading="sharing" :disabled="!isValidCircle" class="mt-2">
              Share with Circle
            </v-btn>
          </v-window-item>
        </v-window>

        <v-divider class="my-4"></v-divider>

        <div v-if="currentShares.length > 0">
          <div class="text-subtitle-1 mb-2">Current Shares</div>
          <v-list>
            <v-list-item v-for="share in currentShares" :key="share.id" :title="share.name" :subtitle="`${share.type} ${share.state === 'ACCESS_STATE_PENDING' ? '(Pending)' : ''}`">
              <template #append>
                <div class="d-flex align-center gap-2">
                  <v-chip size="small" :color="hasWritePermission(share.permission) ? 'primary' : 'grey'">
                    {{ hasWritePermission(share.permission) ? 'Read & Write' : 'Read Only' }}
                  </v-chip>
                  <v-chip v-if="share.state === 'ACCESS_STATE_PENDING'" size="small" color="warning" variant="outlined">
                    Pending
                  </v-chip>
                  <v-btn icon="mdi-delete" variant="text" @click="removeShare(share.id)"></v-btn>
                </div>
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
import type { Recipe_MeasurementType, apitypes_VisibilityLevel } from '@/genapi/api/meals/recipe/v1alpha1'
import type { Access, CreateAccessRequest, ListAccessesRequest, DeleteAccessRequest, Access_RequesterOrRecipient } from '@/genapi/api/meals/recipe/v1alpha1'
import type { PermissionLevel, AccessState } from '@/genapi/api/types'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRecipesStore } from '@/stores/recipes'
import { useRecipeFormStore } from '@/stores/recipeForm'
import { storeToRefs } from 'pinia'
import { onMounted, onBeforeUnmount, watch, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {   recipeService, userService, circleService, recipeAccessService } from '@/api/api'
import type { User, ListUsersRequest } from '@/genapi/api/users/user/v1alpha1'
import type { Circle, ListCirclesRequest } from '@/genapi/api/circles/circle/v1alpha1'
import { useAuthStore } from '@/stores/auth'
import { hasWritePermission } from '@/utils/permissions'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

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
  MEASUREMENT_TYPE_CUP: 'cups',
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
const selectedUser = ref<User | null>(null)
const selectedCircle = ref<Circle | null>(null)
const sharing = ref(false)

// Validation states
const isValidUsername = ref(false)
const isValidCircle = ref(false)
const isLoadingUsername = ref(false)
const isLoadingCircle = ref(false)

// Current shares state
const currentShares = ref<Array<{ 
  id: string; 
  name: string; 
  type: 'user' | 'circle';
  permission: PermissionLevel;
  state: AccessState;
}>>([])

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
    const request: ListUsersRequest = {
      filter: `username = "${username}"`,
      pageSize: 1,
      pageToken: undefined
    }
    const response = await userService.ListUsers(request)

    if (response.users?.length === 1 && response.users[0].name !== authStore.activeAccount?.name) {
      selectedUser.value = response.users[0]
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
    const request: ListCirclesRequest = {
      filter: `title = "${circleName}"`,
      pageSize: 1,
      pageToken: undefined
    }
    const response = await circleService.ListCircles(request)

    if (response.circles?.length === 1 && response.circles[0].name !== authStore.activeAccount?.name) {
      selectedCircle.value = response.circles[0]
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

// Add these new refs and constants
const selectedPermission = ref<PermissionLevel>('PERMISSION_LEVEL_READ')

const permissionOptions = [
  { title: 'Read Only', value: 'PERMISSION_LEVEL_READ' as PermissionLevel },
  { title: 'Read & Write', value: 'PERMISSION_LEVEL_WRITE' as PermissionLevel },
]

// Update the shareRecipe function
async function shareRecipe() {
  if (!selectedUser.value && !selectedCircle.value) return
  if (!recipe.value?.name) return

  sharing.value = true
  try {
    // Create the recipient object based on whether we're sharing with user or circle
    const recipient: Access_RequesterOrRecipient = shareTab.value === 'users'
      ? { user: { name: selectedUser.value?.name || '', username: selectedUser.value?.username || '' } }
      : { circle: { name: selectedCircle.value?.name || '', title: selectedCircle.value?.title || '' } }

    const access: Access = {
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      recipient,
      level: selectedPermission.value,
      state: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: recipe.value.name,
      access
    }

    await recipeAccessService.CreateAccess(request)
    
    // Refresh the recipients list
    await fetchRecipeRecipients()
    
    // Reset selections and inputs after sharing
    selectedUser.value = null
    selectedCircle.value = null
    usernameInput.value = ''
    circleInput.value = ''
    isValidUsername.value = false
    isValidCircle.value = false
    selectedPermission.value = 'PERMISSION_LEVEL_READ' // Reset to default
  } catch (error) {
    console.error('Error sharing recipe:', error)
    // You might want to show an error notification here
  } finally {
    sharing.value = false
  }
}

async function removeShare(shareId: string) {
  try {
    const request: DeleteAccessRequest = {
      name: shareId // shareId is actually the access name in the format recipes/{recipe}/accesses/{access}
    }
    
    await recipeAccessService.DeleteAccess(request)
    
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
      }).map(access => {
        // Determine if this is a user or circle access
        const isCircle = !!access.recipient?.circle
        const recipientName = isCircle
          ? (access.recipient?.circle && access.recipient.circle.title ? access.recipient.circle.title : '')
          : (access.recipient?.user && access.recipient.user.username ? access.recipient.user.username : '')
        return {
          id: access.name || '',
          name: recipientName,
          type: isCircle ? 'circle' as const : 'user' as const,
          permission: access.level || 'PERMISSION_LEVEL_READ',
          state: access.state || 'ACCESS_STATE_PENDING'
        }
      })
    }
  } catch (error) {
    console.error('Error fetching recipe recipients:', error)
  }
}

// Fetch recipients when recipe is loaded or when share dialog is opened
watch(recipe, (newRecipe) => {
  if (newRecipe && hasWritePermission(newRecipe.recipeAccess?.permissionLevel)) {
    fetchRecipeRecipients()
  }
}, { immediate: true })

// Fetch recipients when share dialog is opened
watch(showShareDialog, (isOpen) => {
  if (isOpen && recipe.value && hasWritePermission(recipe.value.recipeAccess?.permissionLevel)) {
    fetchRecipeRecipients()
  }
})

// Remove access dialog state
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

async function handleRemoveAccess() {
  if (!recipe.value?.name || !authStore.activeAccount?.name) return

  removingAccess.value = true
  try {
    // First, we need to find the current user's access to delete it
    const listRequest: ListAccessesRequest = {
      parent: recipe.value.name,
      filter: undefined,
      pageSize: undefined,
      pageToken: undefined
    }
    
    const response = await recipeAccessService.ListAccesses(listRequest)
    
    // Find the current user's access
    const userAccess = response.accesses?.find(access => 
      access.recipient?.user === authStore.activeAccount?.name
    )
    
    if (userAccess?.name) {
      const deleteRequest: DeleteAccessRequest = {
        name: userAccess.name
      }
      
      await recipeAccessService.DeleteAccess(deleteRequest)
      router.push({ name: 'recipes' })
    } else {
      console.error('Could not find user access to remove')
    }
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

async function acceptRecipe() {
  if (!recipe.value?.recipeAccess?.name || !authStore.user?.name) return
  try {
    await recipesStore.acceptRecipe(recipe.value.recipeAccess.name)
    recipesStore.loadRecipe(recipe.value?.name ?? '')
  } catch (error) {
    // Optionally show a notification
  }
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
</script>

<style scoped>
.image-container {
  position: relative;
}
</style>
