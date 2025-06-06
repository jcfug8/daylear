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
import { fileService } from '@/api/api'

const route = useRoute()
const router = useRouter()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()

function navigateBack() {
  router.push({ name: 'recipe', params: { recipeId: route.params.recipeId } })
}

async function saveRecipe(pendingImageFile: File | null) {
  try {
    await recipesStore.updateRecipe()
    
    // Upload image if there's a pending file
    if (pendingImageFile && recipesStore.recipe?.name) {
      const response = await fileService.UploadRecipeImage({
        name: recipesStore.recipe.name,
        file: pendingImageFile,
      })
      
      recipesStore.recipe.imageUri = response.imageUri
    }
    
    navigateBack()
  } catch (err) {
    alert(err instanceof Error ? err.message : String(err))
  }
}

onMounted(async () => {
  // Load the recipe based on the route parameter
  const recipeId = route.params.recipeId as string
  await recipesStore.loadRecipe(recipeId)

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
