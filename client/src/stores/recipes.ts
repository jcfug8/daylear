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
import { recipeAccessService } from '@/api/api'

export const useRecipesStore = defineStore('recipes', () => {
  // const recipes = ref<Recipe[]>([])
  const recipe = ref<Recipe | undefined>()
  const myRecipes = ref<Recipe[]>([])
  const sharedAcceptedRecipes = ref<Recipe[]>([])
  const sharedPendingRecipes = ref<Recipe[]>([])
  const publicRecipes = ref<Recipe[]>([])

  async function loadRecipes(parent: string, filter?: string): Promise<Recipe[]> {
    try {
      const request = {
        parent,
        pageSize: undefined,
        pageToken: undefined,
        filter,
      }
      const response = (await recipeService.ListRecipes(
        request as ListRecipesRequest,
      )) as ListRecipesResponse
      return response.recipes ?? []
    } catch (error) {
      console.error('Failed to load recipes:', error)
      return []
    }
  }

  // Load my recipes (recipes where I have admin permission)
  async function loadMyRecipes(parent: string) {
    const recipes = await loadRecipes(parent, 'permission = 300')
    myRecipes.value = recipes
  }

  // Load shared recipes (recipes shared with me - read or write permission)
  async function loadSharedRecipes(parent: string, state?: number) {
    let filter = 'permission = 100 OR permission = 200'
    if (state) {
      filter += ` AND state = ${state}`
    }
    const recipes = await loadRecipes(parent, filter)
    if (state === 200) {
      sharedAcceptedRecipes.value = recipes
    } else if (state === 100) {
      sharedPendingRecipes.value = recipes
    }
  }

  // Load public recipes (recipes with public visibility)
  async function loadPublicRecipes(parent: string) {
    const recipes = await loadRecipes(parent, 'visibility = 1')
    publicRecipes.value = recipes
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
      description: undefined,
      directions: undefined,
      ingredientGroups: undefined,
      imageUri: undefined,
      visibility: 'VISIBILITY_LEVEL_HIDDEN' as apitypes_VisibilityLevel,
      recipeAccess: undefined,
      citation: undefined,
      cookDuration: undefined,
      cookingMethod: undefined,
      categories: undefined,
      yieldAmount: undefined,
      cuisines: undefined,
      createTime: undefined,
      updateTime: undefined,
      prepDuration: undefined,
      totalDuration: undefined,
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

  async function acceptRecipe(accessName: string) {
    try {
      await recipeAccessService.AcceptRecipeAccess({ name: accessName })
    } catch (error) {
      console.error('Failed to accept recipe access:', error)
      throw error
    }
  }

  async function deleteRecipeAccess(accessName: string) {
    try {
      await recipeAccessService.DeleteAccess({ name: accessName })
    } catch (error) {
      console.error('Failed to decline recipe access:', error)
      throw error
    }
  }

  return {
    loadMyRecipes,
    loadSharedRecipes,
    loadPublicRecipes,
    loadRecipe,
    initEmptyRecipe,
    createRecipe,
    updateRecipe,
    acceptRecipe,
    deleteRecipeAccess,
    recipe,
    myRecipes,
    sharedAcceptedRecipes,
    sharedPendingRecipes,
    publicRecipes,
  }
})
