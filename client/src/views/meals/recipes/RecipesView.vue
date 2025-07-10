<template>
  <ListTabsPage
    ref="tabsPage"
    :tabs="tabs"
  >
    <template #my="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
        <h2>My Recipes</h2>
      </div>
      <RecipeGrid :recipes="items" :loading="loading" />
    </template>
    <template #shared-accepted="{ items, loading }">
      <RecipeGrid :recipes="items" :loading="loading" />
    </template>
    <template #shared-pending="{ items, loading }">
      <RecipeGrid :recipes="items" :loading="loading" @accept="onAcceptRecipe" @decline="onDeclineRecipe" />
      <div v-if="!loading && items.length === 0">No pending shared recipes found.</div>
    </template>
    <template #explore="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
        <h2>Explore Public Recipes</h2>
      </div>
      <RecipeGrid :recipes="items" :loading="loading" />
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
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRecipesStore } from '@/stores/recipes'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useAuthStore } from '@/stores/auth'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'

const authStore = useAuthStore()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()

const acceptingRecipeId = ref<string | null>(null)
const tabsPage = ref()

const tabs = [
  {
    label: 'My Recipes',
    value: 'my',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      await recipesStore.loadMyRecipes(authStore.activeAccountName)
      return [...recipesStore.myRecipes]
    },
  },
  {
    label: 'Shared Recipes',
    value: 'shared',
    subTabs: [
      {
        label: 'Accepted',
        value: 'accepted',
        loader: async () => {
          if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
          await recipesStore.loadSharedRecipes(authStore.activeAccountName, 200)
          return [...recipesStore.sharedAcceptedRecipes]
        },
      },
      {
        label: 'Pending',
        value: 'pending',
        loader: async () => {
          if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
          await recipesStore.loadSharedRecipes(authStore.activeAccountName, 100)
          return [...recipesStore.sharedPendingRecipes]
        },
      },
    ],
  },
  {
    label: 'Explore Recipes',
    value: 'explore',
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
    // Reload both accepted and pending subtabs
    tabsPage.value?.reloadTab('shared', 'accepted')
    tabsPage.value?.reloadTab('shared', 'pending')
  } catch (error) {
    // Optionally show a notification
  } finally {
    acceptingRecipeId.value = null
  }
}

async function onDeclineRecipe(recipe: Recipe) {
  if (!recipe.recipeAccess?.name) return
  try {
    await recipesStore.deleteRecipeAccess(recipe.recipeAccess.name)
    // Reload only the pending subtab
    tabsPage.value?.reloadTab('shared', 'pending')
  } catch (error) {
    // Optionally show a notification
  }
}

breadcrumbStore.setBreadcrumbs([
  { title: 'Recipes', to: { name: 'recipes' } },
])
</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style>
