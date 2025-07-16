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
            <v-list-item-subtitle>{{ userSettings.email }}</v-list-item-subtitle>
          </v-list-item>

          <v-list-item>
            <template v-slot:prepend>
              <v-icon icon="mdi-account-circle"></v-icon>
            </template>
            <v-list-item-title>Username</v-list-item-title>
            <v-list-item-subtitle>{{ user.username }}</v-list-item-subtitle>
          </v-list-item>

          <v-list-item>
            <template v-slot:prepend>
              <v-icon :icon="selectedVisibilityIcon"></v-icon>
            </template>
            <v-list-item-title>Visibility</v-list-item-title>
            <v-list-item-subtitle>
              <strong>{{ selectedVisibilityLabel }}</strong>: {{ selectedVisibilityDescription }}
            </v-list-item-subtitle>
          </v-list-item>
        </v-list>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          :to="{ name: 'user-edit', params: { userId: $route.params.userId } }"
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
import { computed } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'

const authStore = useAuthStore()
const { user, userSettings } = storeToRefs(authStore)
const breadcrumbStore = useBreadcrumbStore()

onMounted(async () => {
  await authStore.loadAuthUser()

  breadcrumbStore.setBreadcrumbs([
    { title: 'User Settings', to: { name: 'user', params: { userId: user.value.name } } },
  ])
})

const visibilityOptions = [
  {
    value: 'VISIBILITY_LEVEL_PUBLIC',
    label: 'Public',
    icon: 'mdi-earth',
    color: 'success',
    description: 'Everyone can see your profile.'
  },
  {
    value: 'VISIBILITY_LEVEL_RESTRICTED',
    label: 'Restricted',
    icon: 'mdi-account-group',
    color: 'warning',
    description: 'Only shared users and their connections can see your profile.'
  },
  {
    value: 'VISIBILITY_LEVEL_PRIVATE',
    label: 'Private',
    icon: 'mdi-lock',
    color: 'info',
    description: 'Only specifically shared users can see your profile.'
  },
  {
    value: 'VISIBILITY_LEVEL_HIDDEN',
    label: 'Hidden',
    icon: 'mdi-eye-off',
    color: 'secondary',
    description: 'Only you can see your profile.'
  }
]

const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === user.value.visibility)
})
const selectedVisibilityLabel = computed(() => selectedVisibility.value?.label || '')
const selectedVisibilityDescription = computed(() => selectedVisibility.value?.description || '')
const selectedVisibilityIcon = computed(() => selectedVisibility.value?.icon || 'mdi-help-circle')
</script>

<style></style>
