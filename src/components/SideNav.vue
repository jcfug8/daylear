<template>
  <v-navigation-drawer permanent v-if="isLoggedIn">
    <v-list nav>
      <v-list-group value="Meals">
        <template v-slot:activator="{ props }">
          <v-list-item v-bind="props" title="Meals" prepend-icon="mdi-food"></v-list-item>
        </template>

        <v-list-item
          prepend-icon="mdi-book-open-page-variant"
          title="Recipes"
          value="recipes"
          :to="{ name: 'recipes' }"
        ></v-list-item>
        <v-list-item
          prepend-icon="mdi-food-apple"
          title="Ingredients"
          value="ingredients"
          :to="{ name: 'ingredients' }"
        ></v-list-item>
      </v-list-group>
      <v-list-item
        prepend-icon="mdi-calendar"
        title="Calendar"
        value="calendar"
        :to="{ name: 'calendar' }"
      ></v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'

const authStore = useAuthStore()
const { isLoggedIn } = storeToRefs(authStore)
</script>
