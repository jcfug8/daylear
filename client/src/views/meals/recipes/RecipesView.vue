<template>
  <div>
    <ListTabsPage
      ref="tabsPage"
      :tabs="tabs"
    >
    <template #filter style="max-width: 1000px; margin: 0 auto;">
      <div class="d-flex align-center ml-2 mt-1" style="gap: 8px;">
        <v-text-field
          v-model="searchQuery"
          label="Search recipes"
          prepend-inner-icon="mdi-magnify"
          clearable
          hide-details
          density="compact"
          class="mt-1"
          style="max-width: 350px;"
        />
        <div class="d-flex align-center flex-wrap" style="gap: 4px;">
          <template v-for="cuisine in selectedCuisines" :key="'cuisine-'+cuisine">
            <v-chip size="small" closable @click:close="removeCuisine(cuisine)">
              {{ cuisine }}
            </v-chip>
          </template>
          <template v-for="category in selectedCategories" :key="'category-'+category">
            <v-chip size="small" closable @click:close="removeCategory(category)">
              {{ category }}
            </v-chip>
          </template>
        </div>
        <div style="flex: 1 1 auto;"></div>
        <v-btn variant="flat" class="mr-2" @click="showFilterModal = true" title="Filter recipes">
          <v-icon>mdi-filter-variant</v-icon>
        </v-btn>
      </div>
    </template>
      <template #my="{ items, loading }">
        <RecipeGrid :recipes="getFilteredRecipes(items)" :loading="loading" />
      </template>
      <template #pending="{ items, loading }">
        <RecipeGrid :recipes="getFilteredRecipes(items)" :loading="loading" @accept="onAcceptRecipe" @decline="onDeclineRecipe" />
        <div v-if="!loading && getFilteredRecipes(items).length === 0">No pending shared recipes found.</div>
      </template>
      <template #explore="{ items, loading }">
        <div class="d-flex justify-space-between align-center mb-4">
        </div>
        <RecipeGrid :recipes="getFilteredRecipes(items)" :loading="loading" />
      </template>
      <template #fab>
        <v-btn
          color="primary"
          icon="mdi-plus"
          style="position: fixed; bottom: 16px; right: 16px"
          :to="{ name: 'recipeCreate' }"
        ></v-btn>
      </template>
    </ListTabsPage>
  </div>
  <v-dialog v-model="showFilterModal" max-width="500">
      <v-card>
        <v-card-title>Filter Recipes</v-card-title>
        <v-card-text>
          <div class="mb-4">
            <div class="font-weight-bold mb-2">Cuisines</div>
            <v-chip-group v-model="selectedCuisines" multiple column>
              <v-chip
                v-for="cuisine in allCuisines"
                :key="cuisine"
                :value="cuisine"
                filter
                class="ma-1"
                @click="toggleCuisine(cuisine)"
                :color="selectedCuisines.includes(cuisine) ? 'primary' : ''"
              >
                {{ cuisine }}
              </v-chip>
            </v-chip-group>
          </div>
          <div>
            <div class="font-weight-bold mb-2">Categories</div>
            <v-chip-group v-model="selectedCategories" multiple column>
              <v-chip
                v-for="category in allCategories"
                :key="category"
                :value="category"
                filter
                class="ma-1"
                @click="toggleCategory(category)"
                :color="selectedCategories.includes(category) ? 'primary' : ''"
              >
                {{ category }}
              </v-chip>
            </v-chip-group>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="clearFilters">Clear</v-btn>
          <v-btn color="primary" @click="showFilterModal = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRecipesStore } from '@/stores/recipes'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useAuthStore } from '@/stores/auth'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'
import { useAlertStore } from '@/stores/alerts'
import Fuse from 'fuse.js'

const authStore = useAuthStore()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()
const alertStore = useAlertStore()

const acceptingRecipeId = ref<string | null>(null)
const tabsPage = ref()

const searchQuery = ref('')
const showFilterModal = ref(false)
const selectedCuisines = ref<string[]>([])
const selectedCategories = ref<string[]>([])

function clearFilters() {
  selectedCuisines.value = []
  selectedCategories.value = []
}

function removeCuisine(cuisine: string) {
  selectedCuisines.value = selectedCuisines.value.filter(c => c !== cuisine)
}
function removeCategory(category: string) {
  selectedCategories.value = selectedCategories.value.filter(c => c !== category)
}
function toggleCuisine(cuisine: string) {
  if (selectedCuisines.value.includes(cuisine)) {
    removeCuisine(cuisine)
  } else {
    selectedCuisines.value = [...selectedCuisines.value, cuisine]
  }
}
function toggleCategory(category: string) {
  if (selectedCategories.value.includes(category)) {
    removeCategory(category)
  } else {
    selectedCategories.value = [...selectedCategories.value, category]
  }
}

const allCuisines = computed(() => {
  // Gather all unique cuisines from all loaded recipes in all tabs
  const recipes = [
    ...recipesStore.myRecipes,
    ...recipesStore.sharedAcceptedRecipes,
    ...recipesStore.sharedPendingRecipes,
    ...recipesStore.publicRecipes,
  ]
  const set = new Set<string>()
  for (const recipe of recipes) {
    if (Array.isArray(recipe.cuisines)) {
      for (const c of recipe.cuisines) set.add(c)
    }
  }
  return Array.from(set).sort()
})

const allCategories = computed(() => {
  const recipes = [
    ...recipesStore.myRecipes,
    ...recipesStore.sharedAcceptedRecipes,
    ...recipesStore.sharedPendingRecipes,
    ...recipesStore.publicRecipes,
  ]
  const set = new Set<string>()
  for (const recipe of recipes) {
    if (Array.isArray(recipe.categories)) {
      for (const c of recipe.categories) set.add(c)
    }
  }
  return Array.from(set).sort()
})

function getFilteredRecipes(items: Recipe[]) {
  let filtered = items
  // Fuzzy search
  if (searchQuery.value) {
    const fuse = new Fuse(filtered, { keys: ['title'], threshold: 0.4 })
    filtered = fuse.search(searchQuery.value).map(result => result.item)
  }
  // Cuisine filter
  if (selectedCuisines.value.length > 0) {
    filtered = filtered.filter(recipe =>
      Array.isArray(recipe.cuisines) && recipe.cuisines.some(c => selectedCuisines.value.includes(c))
    )
  }
  // Category filter
  if (selectedCategories.value.length > 0) {
    filtered = filtered.filter(recipe =>
      Array.isArray(recipe.categories) && recipe.categories.some(c => selectedCategories.value.includes(c))
    )
  }
  return filtered
}

const tabs = [
  {
    label: 'My Recipes',
    value: 'my',
    icon: 'mdi-book-open-variant',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      await recipesStore.loadMyRecipes(authStore.activeAccountName)
      return [...recipesStore.myRecipes]
    },
  },
  {
    label: 'Pending',
    value: 'pending',
    icon: 'mdi-email-arrow-left-outline',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      await recipesStore.loadSharedRecipes(authStore.activeAccountName, 100)
      return [...recipesStore.sharedPendingRecipes]
    },
  },
  {
    label: 'Explore Recipes',
    value: 'explore',
    icon: 'mdi-card-search-outline',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      await recipesStore.loadPublicRecipes(authStore.activeAccountName)
      return [...recipesStore.publicRecipes]
    },
  },
]

async function onAcceptRecipe(recipe: Recipe) {
  if (!recipe.recipeAccess?.name) return
  acceptingRecipeId.value = recipe.recipeAccess.name
  try {
    await recipesStore.acceptRecipe(recipe.recipeAccess.name)
    tabsPage.value?.reloadTab('pending')
  } catch (err) {
    console.log("Error accepting recipe:", err)
    alertStore.addAlert(err instanceof Error ? "Unable to accept recipe\n" + err.message : String(err), 'error')
  } finally {
    acceptingRecipeId.value = null
  }
}

async function onDeclineRecipe(recipe: Recipe) {
  if (!recipe.recipeAccess?.name) return
  try {
    await recipesStore.deleteRecipeAccess(recipe.recipeAccess.name)
    // Reload only the pending subtab
    tabsPage.value?.reloadTab('pending')
  } catch (err) {
    console.log("Error declining recipe:", err)
    alertStore.addAlert(err instanceof Error ? "Unable to decline recipe\n" + err.message : String(err), 'error')
  }
}

</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style>
