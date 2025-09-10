<template>
  <v-app-bar height="48" id="top-nav" color="primary" elevation="0" density="compact">
    <v-app-bar-title >
      Daylear
    </v-app-bar-title>
    <template #append>
      <v-app-bar-nav-icon @click.stop="toggleAccountDrawer" v-if="isLoggedIn">
      </v-app-bar-nav-icon>
    </template>
  </v-app-bar>
  <v-bottom-navigation height="48" bg-color="primary" id="bottom-nav" v-if="isLoggedIn" density="compact">
      <v-btn
        density="compact"
        :to="{ name: 'recipes' }"
      >
        <v-icon>mdi-food-apple-outline</v-icon>
      </v-btn>
      <v-btn
        density="compact"
        :to="{ name: 'calendars' }"
      >
        <v-icon>mdi-calendar</v-icon>
      </v-btn>
      <v-btn
        density="compact"
        :to="{ name: 'users' }"
      >
        <v-icon>mdi-account</v-icon>
      </v-btn>
      <v-btn
        density="compact"
        :to="{ name: 'circles' }"
      >
        <v-icon>mdi-account-group</v-icon>
      </v-btn>
      <v-btn
        density="compact"
        :to="{ name: 'lists' }"
      >
        <v-icon>mdi-format-list-bulleted</v-icon>
      </v-btn>
  </v-bottom-navigation>
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

const toggleAccountDrawer = () => {
  accountDrawer.value = !accountDrawer.value
  navDrawer.value = false
}

const { isLoggedIn, user } = storeToRefs(authStore)
</script>

<style>

</style>
