<template>
  <v-container>
    {{ console.log(recipe) }}
  </v-container>
</template>

<script setup lang="ts">
import { useRecipesStore } from '@/stores/recipes'
import { storeToRefs } from 'pinia'
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const recipesStore = useRecipesStore()

const { recipe } = storeToRefs(recipesStore)

onMounted(() => {
  const recipeId = route.params.recipeId as string
  recipesStore.loadRecipe(recipeId)
})
</script>
