<template>
  <div>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>
    
    <div v-if="!loading && recipes.length === 0" class="text-center py-8">
      <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-book-open-page-variant</v-icon>
      <h3 class="text-grey-lighten-1 mb-2">No recipes found</h3>
      <p class="text-grey-lighten-1">Try switching to a different tab or check back later.</p>
    </div>

    <v-row v-if="!loading && recipes.length > 0">
      <v-col lg="3" md="4" sm="6" cols="12" v-for="recipe in recipes" :key="recipe.name">
        <v-card
          :to="{ name: 'recipe', params: { recipeId: recipe.name } }"
          :title="recipe.title"
          style="aspect-ratio: 8/6"
          hover
          class="recipe-card"
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
          >
            <template v-slot:placeholder>
              <v-row class="fill-height ma-0" align="center" justify="center">
                <v-icon size="48" color="grey-lighten-1">mdi-food</v-icon>
              </v-row>
            </template>
          </v-img>
          
          <!-- Permission level indicator -->
          <v-chip
            v-if="recipe.recipeAccess?.permissionLevel"
            size="small"
            :color="getPermissionColor(recipe.recipeAccess?.permissionLevel)"
            class="permission-chip"
          >
            {{ getPermissionText(recipe.recipeAccess?.permissionLevel) }}
          </v-chip>

          <!-- Accept button for pending recipes -->
        </v-card>
        <v-btn
          v-if="recipe.recipeAccess?.state === 'ACCESS_STATE_PENDING'"
          color="success"
          class="accept-btn"
          @click.stop.prevent="$emit('accept', recipe)"
          block
        >
          Accept
        </v-btn>
        <v-btn
          v-if="recipe.recipeAccess?.state === 'ACCESS_STATE_PENDING'"
          color="error"
          class="decline-btn"
          @click.stop.prevent="$emit('decline', recipe)"
          block
        >
          Decline
        </v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'
import { defineEmits } from 'vue'

interface Props {
  recipes: Recipe[]
  loading?: boolean
}

defineProps<Props>()
const emit = defineEmits(['accept', 'decline'])

function getPermissionColor(permission: string) {
  switch (permission) {
    case 'PERMISSION_LEVEL_ADMIN':
      return 'primary'
    case 'PERMISSION_LEVEL_WRITE':
      return 'orange'
    case 'PERMISSION_LEVEL_READ':
      return 'blue-grey'
    case 'PERMISSION_LEVEL_PUBLIC':
      return 'green'
    default:
      return 'grey'
  }
}

function getPermissionText(permission: string) {
  switch (permission) {
    case 'PERMISSION_LEVEL_ADMIN':
      return 'Admin'
    case 'PERMISSION_LEVEL_WRITE':
      return 'Write'
    case 'PERMISSION_LEVEL_READ':
      return 'Read'
    case 'PERMISSION_LEVEL_PUBLIC':
      return 'Public'
    default:
      return permission.replace('PERMISSION_LEVEL_', '').toLowerCase()
  }
}
</script>

<style scoped>
.recipe-card {
  transition: transform 0.2s ease-in-out;
  position: relative;
}

.recipe-card:hover {
  transform: translateY(-4px);
}

.permission-chip {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 1;
}

.accept-btn {
  margin-top: 12px;
}
</style> 