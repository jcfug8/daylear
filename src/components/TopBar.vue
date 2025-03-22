<template>
  <v-app-bar>
    <v-app-bar-title>
      <div>Daylear</div>
      <div class="text-caption">Plan It Make It Do It</div>
    </v-app-bar-title>
    <template v-slot:append>
      <v-btn
        @click.stop="drawer = !drawer"
        stacked
        prepend-icon="mdi-account-circle"
        v-if="isLoggedIn"
      >
        Jace Fugate
      </v-btn>
    </template>
  </v-app-bar>
  <v-navigation-drawer location="right" v-model="drawer" v-if="isLoggedIn">
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
        <v-btn block @click="logOut"> Logout </v-btn>
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

const drawer = ref(false)
const router = useRouter()

const authStore = useAuthStore()
const logOut = () => {
  drawer.value = false
  authStore.logOut()
  router.push({ name: 'login' })
}
const { isLoggedIn } = storeToRefs(authStore)
</script>
