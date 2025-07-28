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
import { useRouter, useRoute } from 'vue-router'
import RecipeForm from '@/views/meals/recipes/forms/RecipeForm.vue'
import { useAuthStore } from '@/stores/auth'
import { useAlertStore } from '@/stores/alerts'
import { computed } from 'vue'
const router = useRouter()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()
const alertStore = useAlertStore()
const route = useRoute()

function navigateBack() {
  if (route.params.circleId) {
    router.push({ name: 'circle', params: { circleId: route.params.circleId } })
  } else {
    router.push({ name: 'recipes' })
  }
}

const circleName = computed(() => {
  return route.path.replace('/recipes/create', '').slice(1)
})

async function saveRecipe() {
  if (!authStore.user.name && !route.params.circleId) {
    throw new Error('User not authenticated')
  }
  try {
    const recipe = await recipesStore.createRecipe(circleName.value ? circleName.value : authStore.user.name)
    router.push('/'+recipe.name)
  } catch (err) {
    alertStore.addAlert(err instanceof Error ? "Unable to create recipe\n" + err.message : String(err), 'error')
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
