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
import { storeToRefs } from 'pinia'
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import RecipeForm from '@/views/meals/recipes/forms/RecipeForm.vue'
import { fileService } from '@/api/api'
import { useCirclesStore } from '@/stores/circles'
import { useAlertStore } from '@/stores/alerts'

const route = useRoute()
const router = useRouter()
const recipesStore = useRecipesStore()
const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()
const alertStore = useAlertStore()

const { circle } = storeToRefs(circlesStore)

function navigateBack() {
  if (route.params.circleId) {
    router.push({ name: 'circleRecipe', params: { circleId: route.params.circleId, recipeId: route.params.recipeId } })
  } else {
    router.push({ name: 'recipe', params: { recipeId: route.params.recipeId } })
  }
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
    console.log("Error saving recipe:", err)
    alertStore.addAlert(err instanceof Error ? "Unable to save recipe\n" + err.message : String(err), 'error')
  }
}

onMounted(async () => {
  // Load the recipe based on the route parameter
  const recipeId = route.params.recipeId as string
  await recipesStore.loadRecipe(recipeId)

  let firstCrumbs

  if (route.params.circleId) {
    await circlesStore.loadCircle(route.params.circleId as string)
    firstCrumbs = [
      { title: 'Circles', to: { name: 'circles' } },
      { title: circle.value?.title || '', to: { name: 'circle', params: { circleId: route.params.circleId } } },
      { title: recipesStore.recipe?.title ? recipesStore.recipe.title : '', to: { name: 'circleRecipe', params: { circleId: route.params.circleId, recipeId: recipeId } } }
    ]
  } else {
    firstCrumbs = [
      { title: 'Recipes', to: { name: 'recipes' } },
      { title: recipesStore.recipe?.title ? recipesStore.recipe.title : '', to: { name: 'recipe', params: { recipeId: recipeId } } }
    ]
  }

  breadcrumbStore.setBreadcrumbs([
    ...firstCrumbs,
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
