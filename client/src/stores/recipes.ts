import { ref } from 'vue'
import { defineStore } from 'pinia'
import { recipeService } from '@/api/api'
import type {
  Recipe,
  ListRecipesRequest,
  ListRecipesResponse,
  apitypes_VisibilityLevel,
  apitypes_PermissionLevel,
} from '@/genapi/api/meals/recipe/v1alpha1'

export const useRecipesStore = defineStore('recipes', () => {
  const recipes = ref<Recipe[]>([])
  const recipe = ref<Recipe | undefined>()

  async function loadRecipes(parent: string) {
    try {
      const request = {
        parent,
        pageSize: undefined,
        pageToken: undefined,
        filter: undefined,
      }
      const response = (await recipeService.ListRecipes(
        request as ListRecipesRequest,
      )) as ListRecipesResponse
      recipes.value = response.recipes ?? []
    } catch (error) {
      console.error('Failed to load recipes:', error)
      recipes.value = []
    }
  }

  async function loadRecipe(recipeName: string) {
    try {
      const result = await recipeService.GetRecipe({ name: recipeName })
      recipe.value = result
    } catch (error) {
      console.error('Failed to load recipe:', error)
      recipe.value = undefined
    }
  }

  function initEmptyRecipe() {
    recipe.value = {
      name: undefined,
      title: undefined,
      description: '',
      directions: [],
      ingredientGroups: [],
      imageUri: undefined,
      visibility: 'VISIBILITY_LEVEL_UNSPECIFIED' as apitypes_VisibilityLevel,
      permission: 'PERMISSION_LEVEL_UNSPECIFIED' as apitypes_PermissionLevel,
    }
  }

  async function createRecipe(parent: string) {
    if (!recipe.value) {
      throw new Error('No recipe to create')
    }
    if (recipe.value.name) {
      throw new Error('Recipe already has a name and cannot be created')
    }
    console.log('Creating recipe with data:', recipe.value)
    console.log('Parent path:', parent)
    try {
      const created = await recipeService.CreateRecipe({
        recipe: recipe.value,
        recipeId: crypto.randomUUID(),
      })
      recipe.value = created
      return created
    } catch (error) {
      console.error('Failed to create recipe:', error)
      throw error
    }
  }

  async function updateRecipe() {
    if (!recipe.value) {
      throw new Error('No recipe to update')
    }
    if (!recipe.value.name) {
      throw new Error('Recipe must have a name to be updated')
    }
    try {
      const updated = await recipeService.UpdateRecipe({
        recipe: recipe.value,
        updateMask: undefined, // Optionally specify fields to update
      })
      recipe.value = updated
      return updated
    } catch (error) {
      console.error('Failed to update recipe:', error)
      throw error
    }
  }

  return {
    loadRecipes,
    loadRecipe,
    initEmptyRecipe,
    createRecipe,
    updateRecipe,
    recipes,
    recipe,
  }
})
