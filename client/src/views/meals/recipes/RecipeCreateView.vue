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

const router = useRouter()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()

function navigateBack() {
  router.push({ name: 'recipes' })
}

function saveRecipe() {
  recipesStore
    .createRecipe()
    .then(() => navigateBack())
    .catch((err) => alert(err.message || err))
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
