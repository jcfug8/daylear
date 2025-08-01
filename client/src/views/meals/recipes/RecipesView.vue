<template>
  <div>
    <ListTabsPage
      ref="tabsPage"
      :tabs="tabs"
    >
    <template #filter style="max-width: 1000px; margin: 0 auto;">
      <v-row class="align-center" style="max-width: 600px; margin: 0 auto;">
        <!-- Left: Search -->
        <v-col cols="3" class="pa-0 d-flex align-center justify-start">
          <template v-if="searchExpanded || searchQuery">
            <v-text-field
              v-model="searchQuery"
              label="Search recipes"
              prepend-inner-icon="mdi-magnify"
              clearable
              hide-details
              density="compact"
              class="mt-1 search-bar"
              :class="{ expanded: searchExpanded || searchQuery, collapsed: !searchExpanded && !searchQuery }"
              :style="searchBarStyle"
              @focus="onSearchFocus"
              @blur="onSearchBlur"
              @keydown.enter="onSearchEnter"
              ref="searchInput"
            />
          </template>
          <template v-else>
            <v-btn icon variant="text" class="search-icon-btn" @click="expandSearch">
              <v-icon>mdi-magnify</v-icon>
            </v-btn>
          </template>
        </v-col>
        <!-- Center: Active account title/name -->
        <v-col cols="6" class="pa-0 d-flex align-center justify-center">
          <div
            class="active-account-title clickable text-center"
            style="cursor: pointer; user-select: none; font-weight: 500; font-size: 1.1rem; display: flex; align-items: center; justify-content: center; gap: 4px;"
            @click="showFilterModal = true"
          >
            <v-icon size="18">{{ selectedAccount?.icon || 'mdi-account-circle' }}</v-icon>
            <span class="account-title-ellipsis">{{ selectedAccount?.label || 'My Recipes' }}</span>
          </div>
        </v-col>
        <!-- Right: Filter button -->
        <v-col cols="3" class="pa-0 d-flex align-center justify-end">
          <v-btn class="filter-button mr-2" :color="selectedCuisines.length === 0 && selectedCategories.length === 0 ? 'white' : 'grey'" variant="flat" @click="showFilterModal = true" title="Filter recipes">
            <v-icon>mdi-filter-variant</v-icon>
          </v-btn>
        </v-col>
      </v-row>
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
          v-if="selectedAccount?.value === authStore.user.name"
          color="primary"
          density="compact"
          style="position: fixed; bottom: 16px; right: 16px"
          :to="{ name: 'recipeCreate' }"
        >
          <v-icon>mdi-plus</v-icon>
          <span>Create Recipe</span>
        </v-btn>
      </template>
    </ListTabsPage>
  </div>
  <v-dialog v-model="showFilterModal" max-width="500">
    <v-card>
      <v-card-title>Filter Recipes</v-card-title>
      <v-card-text>
        <!-- User/Circle select at the top -->
        <div class="mb-4">
          <div class="font-weight-bold mb-2">Account</div>
          <v-autocomplete
            ref="accountSelectRef"
            v-model="selectedAccount"
            :items="accountOptions"
            item-title="label"
            item-value="value"
            return-object
            hide-details
            density="compact"
            class="mb-4"
            :prepend-inner-icon="selectedAccount?.icon"
            :menu-props="{ maxHeight: '300px' }"
          >
            <template #item="{ props, item }">
              <v-list-item v-bind="props">
                <template #prepend>
                  <v-icon :icon="item.raw.icon" size="small" class="mr-2"></v-icon>
                </template>
              </v-list-item>
            </template>
          </v-autocomplete>
        </div>
        <div class="mb-4">
          <div class="font-weight-bold mb-2">Cuisines</div>
          <v-autocomplete
            v-model="selectedCuisines"
            :items="allCuisines"
            label="Select cuisines"
            multiple
            chips
            clearable
            hide-details
            density="compact"
          />
        </div>
        <div>
          <div class="font-weight-bold mb-2">Categories</div>
          <v-autocomplete
            v-model="selectedCategories"
            :items="allCategories"
            label="Select categories"
            multiple
            chips
            clearable
            hide-details
            density="compact"
          />
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text @click="clearFilters">Clear</v-btn>
        <v-btn color="primary" @click="showFilterModal = false">Ok</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch, onMounted } from 'vue'
import { useRecipesStore } from '@/stores/recipes'
import { useAuthStore } from '@/stores/auth'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'
import { useAlertStore } from '@/stores/alerts'
import Fuse from 'fuse.js'
import { useCirclesStore } from '@/stores/circles'
import { useUsersStore } from '@/stores/users'

const circlesStore = useCirclesStore()
const authStore = useAuthStore()
const recipesStore = useRecipesStore()
const alertStore = useAlertStore()
const usersStore = useUsersStore()

const acceptingRecipeId = ref<string | null>(null)
const tabsPage = ref()

const searchQuery = ref('')
const searchExpanded = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)

// Dropdown for account/circle selection
const selectedAccount = ref<any>(null)

function onSearchEnter() {
  searchExpanded.value = false
}

onMounted(() => {
  circlesStore.loadMyCircles()
  usersStore.loadFriends()
})

const accountOptions = computed(() => {
  const options = []
  if (authStore.user && authStore.user.name) {
    options.push({
      label: 'My Recipes',
      value: authStore.user.name,
      account: "",
      icon: 'mdi-account-circle',
    })
    if (Array.isArray(circlesStore.myCircles)) {
      for (const circle of circlesStore.myCircles) {
        options.push({
          label: circle.title || 'Untitled Circle',
          value: circle.name,
          account: circle,
          icon: 'mdi-account-group',
        })
      }
    }
    if (Array.isArray(usersStore.friends)) {
      for (const user of usersStore.friends) {
        let label = ''
        if (user.givenName || user.familyName) { // user full name
          label = user.givenName + ' ' + user.familyName
          label = label.trim()
        } else if (user.username) { // user username
          label = user.username
        } 
        options.push({
          label: label,
          value: user.name,
          account: user,
          icon: 'mdi-account-circle',
        })
      }
    }
  }
  return options
})

// Set default selectedAccount to user on mount or when user changes
watch(
  () => authStore.user?.name,
  (newUserName) => {
    if (newUserName && (!selectedAccount.value || selectedAccount.value.value !== newUserName)) {
      selectedAccount.value = accountOptions.value[0]
    }
  },
  { immediate: true }
)

// When selectedAccount changes, reload the current tab
watch(
  selectedAccount,
  () => {
    // Only reload if tabsPage is ready
    nextTick(() => {
      tabsPage.value?.reloadActiveTab()
    })
  }
)

function expandSearch() {
  searchExpanded.value = true
  nextTick(() => {
    if (searchInput.value && searchInput.value.focus) {
      searchInput.value.focus()
    }
  })
}
function onSearchFocus() {
  searchExpanded.value = true
}
function onSearchBlur() {
  if (!searchQuery.value) {
    searchExpanded.value = false
  }
}

const searchBarStyle = computed(() => {
  return searchExpanded.value || searchQuery.value
    ? { maxWidth: '350px', width: '100%', transition: 'max-width 0.3s cubic-bezier(0.4,0,0.2,1)' }
    : { maxWidth: '44px', width: '44px', transition: 'max-width 0.3s cubic-bezier(0.4,0,0.2,1)' }
})

const showFilterModal = ref(false)
const selectedCuisines = ref<string[]>([])
const selectedCategories = ref<string[]>([])

function clearFilters() {
  selectedCuisines.value = []
  selectedCategories.value = []
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
      // Use selectedAccount for context
      const account = selectedAccount.value?.account
      await recipesStore.loadMyRecipes(account.name)
      return [...recipesStore.myRecipes]
    },
  },
  {
    label: 'Pending',
    value: 'pending',
    icon: 'mdi-email-arrow-left-outline',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      const account = selectedAccount.value?.account
      await recipesStore.loadPendingRecipes(account.name)
      return [...recipesStore.sharedPendingRecipes]
    },
  },
  {
    label: 'Explore',
    value: 'explore',
    icon: 'mdi-card-search-outline',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      const account = selectedAccount.value?.account
      await recipesStore.loadPublicRecipes(account.name)
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
    alertStore.addAlert(err instanceof Error ? "Unable to decline recipe\n" + err.message : String(err), 'error')
  }
}

</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}

.search-bar {
  transition: max-width 0.3s cubic-bezier(0.4,0,0.2,1), width 0.3s cubic-bezier(0.4,0,0.2,1);
  min-width: 44px;
}
.search-bar.collapsed {
  max-width: 44px !important;
  width: 44px !important;
  padding-left: 0 !important;
}
.search-bar.expanded {
  max-width: 350px !important;
  width: 100% !important;
}
.search-icon-btn {
  min-width: 44px;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.account-title-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
  display: inline-block;
  vertical-align: middle;
}
</style>
