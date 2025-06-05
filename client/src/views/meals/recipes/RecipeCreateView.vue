<template>
  <recipe-form
    v-if="recipesStore.recipe"
    v-model="recipesStore.recipe"
    @save="saveRecipe"
    @close="navigateBack"
  />
</template>

<script setup lang="ts">
import { useRecipesStore } from '@/stores/recipes'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import RecipeForm from '@/views/meals/recipes/forms/RecipeForm.vue'
import { useAuthStore } from '@/stores/auth'
const router = useRouter()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()

function navigateBack() {
  router.push({ name: 'recipes' })
}

async function saveRecipe() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  try {
    await recipesStore.createRecipe(authStore.activeAccountName)
    navigateBack()
  } catch (err) {
    alert(err instanceof Error ? err.message : String(err))
  }
}

onMounted(() => {
  // Initialize an empty recipe
  recipesStore.initEmptyRecipe()

  breadcrumbStore.setBreadcrumbs([
    { title: 'Recipes', to: { name: 'recipes' } },
    { title: 'Create New Recipe', to: { name: 'recipeCreate' } },
  ])
})
</script>
