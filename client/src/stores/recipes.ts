import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'

export const useRecipesStore = defineStore('recipes', () => {
  const recipes = ref<Recipe[]>([])
  const recipe = ref<Recipe | undefined>()

  function loadRecipes() {
    recipes.value = [
      {
        name: 'users/1/recipes/1',
        title: 'Spaghetti Bolognese',
        description: 'A classic Italian pasta dish with a rich meat sauce.',
        directions: undefined,
        ingredientGroups: undefined,
        imagePath: undefined,
        imageUri: 'https://healthfulblondie.com/wp-content/uploads/2022/05/Chicken-Bolognese.jpg',
      },
      {
        name: 'users/1/recipes/2',
        title: 'Chicken Curry',
        description: 'A spicy and flavorful chicken curry.',
        directions: undefined,
        ingredientGroups: undefined,
        imagePath: undefined,
        imageUri: 'https://upload.wikimedia.org/wikipedia/commons/7/76/Creamy_Chicken_Curry.jpg',
      },
      {
        name: 'users/1/recipes/3',
        title: 'Vegetable Stir Fry',
        description: 'A quick and healthy vegetable stir fry.',
        directions: undefined,
        ingredientGroups: undefined,
        imagePath: undefined,
        imageUri: '',
      },
      {
        name: 'users/1/recipes/4',
        title: 'Beef Tacos',
        description: 'Delicious beef tacos with fresh toppings.',
        directions: undefined,
        ingredientGroups: undefined,
        imagePath: undefined,
        imageUri:
          'https://cdn12.picryl.com/photo/2016/12/31/taco-mexican-beef-food-drink-256d7b-1024.jpg',
      },
      {
        name: 'users/1/recipes/5',
        title: 'Caesar Salad',
        description: 'A classic Caesar salad with romaine lettuce and croutons.',
        directions: undefined,
        ingredientGroups: undefined,
        imagePath: undefined,
        imageUri: '',
      },
      {
        name: 'users/1/recipes/6',
        title: 'Chocolate Cake',
        description: 'A rich and moist chocolate cake.',
        directions: undefined,
        ingredientGroups: undefined,
        imagePath: undefined,
        imageUri:
          'https://bakesbybrownsugar.com/wp-content/uploads/2024/02/Chocolate-Sour-Cream-Pound-Cake-22-360x360.jpg',
      },
      {
        name: 'users/1/recipes/7',
        title: 'Grilled Salmon',
        description: 'Perfectly grilled salmon with lemon and herbs.',
        directions: undefined,
        ingredientGroups: undefined,
        imagePath: undefined,
        imageUri:
          'https://freerangestock.com/sample/166079/grilled-salmon-with-veggies-on-stone-plate.jpg',
      },
    ]
  }

  function loadRecipe(recipeId: string) {
    loadRecipes()
    const recipeFound = recipes.value.find((r) => r.name === recipeId)
    if (recipeFound) {
      recipe.value = recipeFound
    } else {
      recipe.value = undefined
    }
  }

  return {
    loadRecipes,
    loadRecipe,
    recipes,
    recipe,
  }
})
