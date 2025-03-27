<template>
  <v-app-bar>
    <template #prepend>
      <v-app-bar-nav-icon
        @click.stop="navDrawer = !navDrawer"
        v-if="isLoggedIn"
      ></v-app-bar-nav-icon>
    </template>
    <v-app-bar-title class="d-flex justify-center align-center">
      <div class="d-flex justify-center align-center">Daylear</div>
      <div class="text-caption">Plan It Make It Do It</div>
    </v-app-bar-title>
    <template #append>
      <v-btn icon @click.stop="accountDrawer = !accountDrawer" v-if="isLoggedIn">
        <v-icon>mdi-account-circle</v-icon>
      </v-btn>
    </template>
  </v-app-bar>
  <v-navigation-drawer temporary v-if="isLoggedIn" v-model="navDrawer" density="compact">
    <v-list nav>
      <v-list-item
        prepend-icon="mdi-calendar"
        title="Calendar"
        value="calendar"
        :to="{ name: 'calendar' }"
      ></v-list-item>
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
    </v-list>
  </v-navigation-drawer>
  <v-navigation-drawer temporary location="right" v-model="accountDrawer" v-if="isLoggedIn">
    <v-list>
      <v-list-item
        prepend-icon="mdi-account-circle"
        title="Jace Fugate"
        :to="{ name: 'account', params: { accountId: 'jace-fugate' } }"
      >
        <template #append>
          <v-btn
            icon="mdi-cog"
            density="comfortable"
            :to="{ name: 'account-settings', params: { accountId: 'jace-fugate' } }"
          ></v-btn>
        </template>
      </v-list-item>
      <v-divider></v-divider>
      <v-list-item
        prepend-icon="mdi-account-group"
        title="Fugate Family"
        :to="{ name: 'circle', params: { circleId: 'fugate-family' } }"
      >
        <template #append>
          <v-btn
            icon="mdi-cog"
            density="comfortable"
            :to="{ name: 'circle-settings', params: { circleId: 'fugate-family' } }"
          ></v-btn>
        </template>
      </v-list-item>
      <v-list-item
        prepend-icon="mdi-account-group"
        title="Big Fugate Family"
        :to="{ name: 'circle', params: { circleId: 'big-fugate-family' } }"
      >
        <template #append>
          <v-btn
            icon="mdi-cog"
            density="comfortable"
            :to="{ name: 'circle-settings', params: { circleId: 'big-fugate-family' } }"
          ></v-btn>
        </template>
      </v-list-item>
    </v-list>
    <template v-slot:append>
      <v-divider></v-divider>
      <div class="pa-2">
        <v-btn prepend-icon="mdi-logout" block @click="logOut"> Logout </v-btn>
      </div>
    </template>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { VNavigationDrawer } from 'vuetify/components'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

import { ref } from 'vue'

const accountDrawer = ref(false)
const navDrawer = ref(false)
const router = useRouter()

const authStore = useAuthStore()
const logOut = () => {
  accountDrawer.value = false
  authStore.logOut()
  router.push({ name: 'login' })
}
const { isLoggedIn } = storeToRefs(authStore)
</script>
