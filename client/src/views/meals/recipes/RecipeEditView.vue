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
import { computed } from 'vue'

const route = useRoute()
const router = useRouter()
const recipesStore = useRecipesStore()
const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()
const alertStore = useAlertStore()

const { circle } = storeToRefs(circlesStore)

const recipeName = computed(() => {
  return route.path.replace('/edit', '')
})

const trimmedRecipeName = computed(() => {
  return recipeName.value.substring(recipeName.value.indexOf('recipes/'))
})

const circleName = computed(() => {
  return route.path.indexOf('/recipes/') !== -1 
    ? route.path.substring(0, route.path.indexOf('/recipes/'))
    : null
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

  let firstCrumbs

  if (route.params.circleId) {
    await circlesStore.loadCircle(circleName.value)
    firstCrumbs = [
      { title: 'Circles', to: { name: 'circles' } },
      { title: circle.value?.title || '', to: circleName.value },
      { title: recipesStore.recipe?.title ? recipesStore.recipe.title : '', to: recipeName.value }
    ]
  } else {
    firstCrumbs = [
      { title: 'Recipes', to: { name: 'recipes' } },
      { title: recipesStore.recipe?.title ? recipesStore.recipe.title : '', to: recipeName.value }
    ]
  }

  breadcrumbStore.setBreadcrumbs([
    ...firstCrumbs,
    { title: 'Edit'},
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
