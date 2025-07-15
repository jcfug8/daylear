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
        <v-icon v-if="authStore.activeAccountType === AccountType.USER">mdi-account-circle</v-icon>
        <v-icon v-if="authStore.activeAccountType === AccountType.CIRCLE">mdi-account-group</v-icon>
        <span class="text-caption">{{ authStore.activeAccountTitle }}</span>
      </v-btn>
    </template>
  </v-app-bar>
  <v-navigation-drawer temporary v-if="isLoggedIn" v-model="navDrawer" density="compact">
    <v-list nav>
      <!-- <v-list-item
        prepend-icon="mdi-calendar"
        title="Calendar"
        value="calendar"
        :to="{ name: 'calendar' }"
      ></v-list-item> -->
      <v-list-item
        v-if="authStore.activeAccountType === AccountType.USER"
        prepend-icon="mdi-account-group"
        title="Circles"
        value="circles"
        :to="{ name: 'circles' }"
      ></v-list-item>
      <!-- <v-list-group value="Meals">
        <template v-slot:activator="{ props }">
          <v-list-item v-bind="props" title="Meals" prepend-icon="mdi-food"></v-list-item>
        </template> -->

        <v-list-item
          prepend-icon="mdi-book-open-page-variant"
          title="Recipes"
          value="recipes"
          :to="{ name: 'recipes' }"
        ></v-list-item>
        <!-- <v-list-item
          prepend-icon="mdi-food-apple"
          title="Ingredients"
          value="ingredients"
          :to="{ name: 'ingredients' }"
        ></v-list-item> -->
      <!-- </v-list-group> -->
    </v-list>
  </v-navigation-drawer>
  <v-app-bar v-if="$route.meta.breadcrumbs && breadcrumbs.length > 0" flat density="compact">
    <v-breadcrumbs density="compact"
      :items="breadcrumbs.map((crumb) => ({ title: crumb.title, to: crumb.to }))"
    ></v-breadcrumbs>
  </v-app-bar>
  <v-navigation-drawer temporary location="right" v-model="accountDrawer" v-if="isLoggedIn">
    <v-list>
      <v-list-item
        :active="user?.name === activeAccount?.name"
        prepend-icon="mdi-account-circle"
        :title="user?.username"
        @click="activateAccount(user)"
      >
        <template #append>
          <v-btn
            icon="mdi-cog"
            density="comfortable"
            @click.stop
            :to="{ name: 'user-settings', params: { userId: user?.name } }"
          ></v-btn>
        </template>
      </v-list-item>
      <v-divider></v-divider>
      <v-list-item
        v-for="circle in circles"
        :key="circle.name"
        :active="circle.name === activeAccount?.name"
        prepend-icon="mdi-account-group"
        :title="circle.title"
        @click="activateAccount(circle)"
      >
        <template #append>
          <v-btn
            icon="mdi-cog"
            density="comfortable"
            @click.stop
            :to="{ name: 'circle', params: { circleId: circle.name } }"
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
  <AlertStack />
</template>

<script setup lang="ts">
import { useAuthStore, AccountType } from '@/stores/auth'
import { VNavigationDrawer } from 'vuetify/components'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import type { User } from '@/genapi/api/users/user/v1alpha1'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'

import { ref } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import AlertStack from './common/AlertStack.vue'

const accountDrawer = ref(false)
const navDrawer = ref(false)
const router = useRouter()

const authStore = useAuthStore()
const breadcrumbStore = useBreadcrumbStore()

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

const activateAccount = (account: User | Circle) => {
  authStore.setActiveAccount(account)
  toggleAccountDrawer()
  router.push({ name: 'calendar' })
}

const { isLoggedIn, activeAccount, user, circles } = storeToRefs(authStore)
const { breadcrumbs } = storeToRefs(breadcrumbStore)
console.log('breadcrumbs', breadcrumbs.value)
</script>

<style>
.v-breadcrumbs-item--link {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: inline-block;
  max-width: 33vw;
  vertical-align: bottom;
}
</style>
