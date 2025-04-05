<template>
  <v-container>
    <v-row>
      <v-col lg="3" md="4" sm="6" cols="12" v-for="recipe in recipes" :key="recipe.name">
        <v-card
          :to="{ name: 'recipe', params: { recipeId: recipe.name } }"
          :title="recipe.title"
          style="aspect-ratio: 8/6"
        >
          <v-card-subtitle>
            {{ recipe.description }}
          </v-card-subtitle>
          <v-img
            class="mt-4"
            style="background-color: lightgray"
            height="100%"
            :src="recipe.imageUri"
            cover
          ></v-img>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { useRecipesStore } from '@/stores/recipes'
import { storeToRefs } from 'pinia'
import { onMounted } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'

const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()
breadcrumbStore.setBreadcrumbs([{ title: 'Recipes', href: '/meals/recipes' }])

onMounted(() => {
  recipesStore.loadRecipes()
})

const { recipes } = storeToRefs(recipesStore)
</script>

<style></style>
