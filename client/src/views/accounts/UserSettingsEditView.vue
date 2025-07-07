<template>
  <v-container>
    <v-card class="mx-auto" max-width="600">
      <v-card-title>Edit User Settings</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="saveSettings">
          <v-text-field
            v-model="editedUser.givenName"
            label="Given Name"
          ></v-text-field>

          <v-text-field
            v-model="editedUser.familyName"
            label="Family Name"
          ></v-text-field>

          <v-text-field
            disabled
            v-model="editedUser.email"
            label="Email"
            type="email"
            required
          ></v-text-field>

          <v-text-field
            v-model="editedUser.username"
            label="Username"
            required
          ></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="secondary"
          @click="navigateBack"
        >
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          @click="saveSettings"
        >
          Save Changes
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { onMounted, ref } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRouter } from 'vue-router'
import type { User, apitypes_VisibilityLevel } from '@/genapi/api/users/user/v1alpha1'

const router = useRouter()
const authStore = useAuthStore()
const { user } = storeToRefs(authStore)
const breadcrumbStore = useBreadcrumbStore()

// Create a copy of the user object for editing
const editedUser = ref<User>({
  name: user.value.name,
  email: user.value.email,
  username: user.value.username,
  givenName: user.value.givenName,
  familyName: user.value.familyName,
  visibility: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
})

function navigateBack() {
  router.push({ name: 'user-settings', params: { userId: user.value.name } })
}

async function saveSettings() {
  try {
    await authStore.updateAuthUser(editedUser.value)
    navigateBack()
  } catch (error) {
    console.error('Error saving settings:', error)
    alert('Failed to save settings')
  }
}

onMounted(async () => {
  await authStore.loadAuthUser()

  breadcrumbStore.setBreadcrumbs([
    { title: 'User Settings', to: { name: 'user-settings', params: { userId: user.value.name } } },
    { title: 'Edit', to: { name: 'user-settings-edit', params: { userId: user.value.name } } },
  ])
})
</script>

<style></style>
