<script setup lang="ts">
import { VNavigationDrawer } from 'vuetify/components'

import { ref } from 'vue'

const drawer = ref(false)
</script>

<template>
  <v-app-bar>
    <v-app-bar-title>
      <div>Daylear</div>
      <div class="text-caption">Plan It Make It Do It</div>
    </v-app-bar-title>
    <template v-slot:append>
      <v-btn @click.stop="drawer = !drawer" stacked prepend-icon="mdi-account-circle">
        Jace Fugate
      </v-btn>
    </template>
  </v-app-bar>
  <v-navigation-drawer location="right" v-model="drawer">
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
        <v-btn block> Logout </v-btn>
      </div>
    </template>
  </v-navigation-drawer>
</template>
