<template>
  <v-container>
    <v-tabs v-model="activeTab" align-tabs="center" color="primary" grow>
      <v-tab value="my">My Recipes</v-tab>
      <v-tab value="shared">Shared Recipes</v-tab>
      <v-tab value="explore">Explore Recipes</v-tab>
    </v-tabs>

    <v-card-text>
      <v-tabs-window v-model="activeTab">
        <!-- My Recipes Tab -->
        <v-tabs-window-item value="my">
          <div class="d-flex justify-space-between align-center mb-4">
            <h2>My Recipes</h2>
          </div>
          <RecipeGrid :recipes="myRecipes" :loading="loadingMyRecipes" />
        </v-tabs-window-item>

        <!-- Shared Recipes Tab -->
        <v-tabs-window-item value="shared">
          <div class="mb-4">
            <h2 class="mb-4">Shared Recipes</h2>
            <v-tabs v-model="sharedTab" density="compact" color="secondary">
              <v-tab value="accepted">Accepted</v-tab>
              <v-tab value="pending">Pending</v-tab>
            </v-tabs>
          </div>
          
          <v-tabs-window v-model="sharedTab">
            <v-tabs-window-item value="accepted">
              <RecipeGrid :recipes="sharedAcceptedRecipes" :loading="loadingSharedRecipes" />
            </v-tabs-window-item>
            <v-tabs-window-item value="pending">
              <RecipeGrid :recipes="sharedPendingRecipes" :loading="loadingSharedRecipes" />
            </v-tabs-window-item>
          </v-tabs-window>
        </v-tabs-window-item>

        <!-- Explore Recipes Tab -->
        <v-tabs-window-item value="explore">
          <div class="d-flex justify-space-between align-center mb-4">
            <h2>Explore Public Recipes</h2>
          </div>
          <RecipeGrid :recipes="publicRecipes" :loading="loadingPublicRecipes" />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card-text>

    <!-- Floating Action Button - only visible on My Recipes tab -->
    <v-btn
      v-if="activeTab === 'my' && hasWritePermission(authStore.activeAccountPermissionLevel)"
      color="primary"
      icon="mdi-plus"
      style="position: fixed; bottom: 16px; right: 16px"
      :to="{ name: 'recipeCreate' }"
    ></v-btn>
  </v-container>
</template>

<script setup lang="ts">
import { useRecipesStore } from '@/stores/recipes'
import { ref, onMounted, watch } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useAuthStore } from '@/stores/auth'
import { hasWritePermission } from '@/utils/permissions'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'
import RecipeGrid from '@/components/RecipeGrid.vue'

const authStore = useAuthStore()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()
breadcrumbStore.setBreadcrumbs([{ title: 'Recipes', to: { name: 'recipes' } }])

const activeTab = ref('my')
const sharedTab = ref('accepted')

// Separate recipe arrays for each view
const myRecipes = ref<Recipe[]>([])
const sharedAcceptedRecipes = ref<Recipe[]>([])
const sharedPendingRecipes = ref<Recipe[]>([])
const publicRecipes = ref<Recipe[]>([])

// Loading states
const loadingMyRecipes = ref(false)
const loadingSharedRecipes = ref(false)
const loadingPublicRecipes = ref(false)

async function loadMyRecipes() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  loadingMyRecipes.value = true
  try {
    await recipesStore.loadMyRecipes(authStore.activeAccountName)
    myRecipes.value = [...recipesStore.recipes]
  } catch (error) {
    console.error('Failed to load my recipes:', error)
    myRecipes.value = []
  } finally {
    loadingMyRecipes.value = false
  }
}

async function loadSharedRecipes() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  loadingSharedRecipes.value = true
  try {
    // Load accepted shared recipes
    await recipesStore.loadSharedRecipes(authStore.activeAccountName, 200)
    sharedAcceptedRecipes.value = [...recipesStore.recipes]
    // Load pending shared recipes
    await recipesStore.loadSharedRecipes(authStore.activeAccountName, 100)
    sharedPendingRecipes.value = [...recipesStore.recipes]
  } catch (error) {
    console.error('Failed to load shared recipes:', error)
    sharedAcceptedRecipes.value = []
    sharedPendingRecipes.value = []
  } finally {
    loadingSharedRecipes.value = false
  }
}

async function loadPublicRecipes() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  loadingPublicRecipes.value = true
  try {
    await recipesStore.loadPublicRecipes(authStore.activeAccountName)
    publicRecipes.value = [...recipesStore.recipes]
  } catch (error) {
    console.error('Failed to load public recipes:', error)
    publicRecipes.value = []
  } finally {
    loadingPublicRecipes.value = false
  }
}

// Load data when component mounts
onMounted(async () => {
  await loadMyRecipes()
})

// Watch for tab changes and load data accordingly
watch(activeTab, async (newTab) => {
  switch (newTab) {
    case 'my':
      if (myRecipes.value.length === 0) {
        await loadMyRecipes()
      }
      break
    case 'shared':
      if (sharedAcceptedRecipes.value.length === 0) {
        await loadSharedRecipes()
      }
      break
    case 'explore':
      if (publicRecipes.value.length === 0) {
        await loadPublicRecipes()
      }
      break
  }
})
</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style>
