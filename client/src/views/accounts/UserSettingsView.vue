<template>
  <v-container>
    <v-card class="mx-auto" max-width="600">
      <v-card-title>User Settings</v-card-title>
      <v-card-text>
        <v-list>
          <v-list-item>
            <template v-slot:prepend>
              <v-icon icon="mdi-account"></v-icon>
            </template>
            <v-list-item-title>Given Name</v-list-item-title>
            <v-list-item-subtitle>{{ user.givenName || 'Not set' }}</v-list-item-subtitle>
          </v-list-item>

          <v-list-item>
            <template v-slot:prepend>
              <v-icon icon="mdi-account"></v-icon>
            </template>
            <v-list-item-title>Family Name</v-list-item-title>
            <v-list-item-subtitle>{{ user.familyName || 'Not set' }}</v-list-item-subtitle>
          </v-list-item>

          <v-list-item>
            <template v-slot:prepend>
              <v-icon icon="mdi-email"></v-icon>
            </template>
            <v-list-item-title>Email</v-list-item-title>
            <v-list-item-subtitle>{{ user.email }}</v-list-item-subtitle>
          </v-list-item>

          <v-list-item>
            <template v-slot:prepend>
              <v-icon icon="mdi-account-circle"></v-icon>
            </template>
            <v-list-item-title>Username</v-list-item-title>
            <v-list-item-subtitle>{{ user.username }}</v-list-item-subtitle>
          </v-list-item>
        </v-list>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          :to="{ name: 'user-settings-edit', params: { userId: $route.params.userId } }"
        >
          Edit Settings
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { onMounted } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'

const authStore = useAuthStore()
const { user } = storeToRefs(authStore)
const breadcrumbStore = useBreadcrumbStore()

onMounted(async () => {
  await authStore.loadAuthUser()

  breadcrumbStore.setBreadcrumbs([
    { title: 'User Settings', to: { name: 'user-settings', params: { userId: user.value.name } } },
  ])
})
</script>

<style></style>
