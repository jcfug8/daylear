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
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import RecipeForm from '@/views/meals/recipes/forms/RecipeForm.vue'
import { fileService } from '@/api/api'
import { useAlertStore } from '@/stores/alerts'
import { computed } from 'vue'

const route = useRoute()
const router = useRouter()
const recipesStore = useRecipesStore()
const alertStore = useAlertStore()

const recipeName = computed(() => {
  return route.path.replace('/edit', '')
})

const trimmedRecipeName = computed(() => {
  return recipeName.value.substring(recipeName.value.indexOf('recipes/'))
})

  
function navigateBack() {
  router.push(recipeName.value)
}

async function saveRecipe(pendingImageFile: File | null) {
  try {
    await recipesStore.updateRecipe()
    
    // Upload image if there's a pending file
    if (pendingImageFile && recipesStore.recipe?.name) {
      const response = await fileService.UploadRecipeImage({
        name: trimmedRecipeName.value,
        file: pendingImageFile,
      })
      
      recipesStore.recipe.imageUri = response.imageUri
    }
    
    navigateBack()
  } catch (err) {
    alertStore.addAlert(err instanceof Error ? "Unable to save recipe\n" + err.message : String(err), 'error')
  }
}

onMounted(async () => {
  // Load the recipe based on the route parameter
  await recipesStore.loadRecipe(recipeName.value)
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
