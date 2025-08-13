<template>
  <v-app-bar density="compact">
    <template #prepend>
      <v-app-bar-nav-icon @click.stop="toggleNavDrawer" v-if="isLoggedIn"></v-app-bar-nav-icon>
    </template>
    <v-app-bar-title class="d-flex justify-center align-center">
      <div class="d-flex justify-center align-center">Daylear</div>
      <div class="text-caption">Plan It Make It Do It</div>
    </v-app-bar-title>
    <template #append>
      <v-btn rounded="0" width="100" icon @click.stop="toggleAccountDrawer" v-if="isLoggedIn" stacked>
        <v-icon>mdi-account-circle</v-icon>
        <span class="text-caption">{{ user?.username }}</span>
      </v-btn>
    </template>
  </v-app-bar>
  <v-navigation-drawer temporary v-if="isLoggedIn" v-model="navDrawer" density="compact">
    <v-list nav>
      <v-list-item
        prepend-icon="mdi-book-open-page-variant"
        title="Recipes"
        value="recipes"
        :to="{ name: 'recipes' }"
      ></v-list-item>
      <v-list-item
        prepend-icon="mdi-calendar"
        title="Calendars"
        value="calendars"
        :to="{ name: 'calendars' }"
      ></v-list-item>
      <v-list-item
        prepend-icon="mdi-account"
        title="Users"
        value="users"
        :to="{ name: 'users' }"
      ></v-list-item>
      <v-list-item
        prepend-icon="mdi-account-group"
        title="Circles"
        value="circles"
        :to="{ name: 'circles' }"
      ></v-list-item>
    </v-list>
  </v-navigation-drawer>
  <v-navigation-drawer temporary location="right" v-model="accountDrawer" v-if="isLoggedIn">
    <v-list>
      <v-list-item
        prepend-icon="mdi-account-circle"
        title="Profile"
        :to="'/'+user?.name"
      >
      </v-list-item>
    </v-list>
    <template v-slot:append>
      <v-divider></v-divider>
      <div class="pa-2">
        <v-btn prepend-icon="mdi-logout" block @click="logOut"> Logout </v-btn>
      </div>
    </template>
  </v-navigation-drawer>
  <AlertStack />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { VNavigationDrawer } from 'vuetify/components'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

import { ref } from 'vue'
import AlertStack from './common/AlertStack.vue'

const accountDrawer = ref(false)
const navDrawer = ref(false)
const router = useRouter()

const authStore = useAuthStore()

const logOut = () => {
  accountDrawer.value = false
  authStore.logOut()
  router.push({ name: 'login' })
}

const toggleNavDrawer = () => {
  navDrawer.value = !navDrawer.value
  accountDrawer.value = false
}

const toggleAccountDrawer = () => {
  accountDrawer.value = !accountDrawer.value
  navDrawer.value = false
}

const { isLoggedIn, user } = storeToRefs(authStore)
</script>

<style>

</style>
