<template>
  <v-container>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>
    
    <div v-if="!loading && recipes.length === 0" class="text-center py-8">
      <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-book-open-page-variant</v-icon>
      <h3 class="text-grey-lighten-1 mb-2">No recipes found</h3>
      <p class="text-grey-lighten-1">Try switching to a different tab or check back later.</p>
    </div>

    <v-row v-if="!loading && recipes.length > 0">
      <v-col class="pa-1" md="3" sm="4" cols="6" v-for="recipe in recipes" :key="recipe.name">
        <v-card
          :to="'/'+recipe.name"
          style="aspect-ratio: 8/6;border-color: lightgrey;border-width: 1.5px;border-style: solid;"
          hover
          class="recipe-card"
        >
          <v-card-title style="font-size: 1rem;">
            {{ recipe.title }}
            <v-icon 
                v-if="recipe.favorited"
                size="24" 
                class="favorite-heart"
              >
              mdi-heart
              </v-icon>
          </v-card-title>
          <v-card-subtitle style="font-size: 0.8rem;">
            {{ recipe.description }}
          </v-card-subtitle>
          <v-img
            class="mt-2"
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
            variant="elevated"
            :color="getPermissionColor(recipe.recipeAccess?.permissionLevel)"
            class="permission-chip"
          >
            {{ getPermissionText(recipe.recipeAccess?.permissionLevel) }}
          </v-chip>

        </v-card>
        <!-- Accept button for pending recipes -->
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
  </v-container>
</template>

<script setup lang="ts">
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'

interface Props {
  recipes: Recipe[]
  loading?: boolean
}

defineProps<Props>()

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
  bottom: 8px;
  right: 8px;
  z-index: 1;
}

.favorite-heart {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.6));
  border-radius: 50%;
  padding: 4px;
  transition: all 0.2s ease-in-out;
  color: red;
}

.accept-btn {
  margin-top: 12px;
}
</style> 