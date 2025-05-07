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
import { useRoute, useRouter } from 'vue-router'
import RecipeForm from '@/views/meals/recipes/forms/RecipeForm.vue'

const route = useRoute()
const router = useRouter()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()

function navigateBack() {
  router.push({ name: 'recipe', params: { recipeId: route.params.recipeId } })
}

function saveRecipe() {
  recipesStore
    .updateRecipe()
    .then(() => navigateBack())
    .catch((err) => alert(err.message || err))
}

onMounted(() => {
  // Load the recipe based on the route parameter
  const recipeId = route.params.recipeId as string
  recipesStore.loadRecipe(recipeId)

  breadcrumbStore.setBreadcrumbs([
    { title: 'Recipes', to: { name: 'recipes' } },
    {
      title: recipesStore.recipe?.title ? recipesStore.recipe.title : '',
      to: { name: 'recipe', params: { recipeId: recipeId } },
    },
    { title: 'Edit', to: { name: 'recipeEdit', params: { recipeId: recipeId } } },
  ])
})
</script>

<style scoped>
.gap-1 {
  gap: 4px;
}

.gap-2 {
  gap: 8px;
}
</style>
