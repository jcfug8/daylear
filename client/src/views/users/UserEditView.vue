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

          <v-select
            v-model="editedUser.visibility"
            :items="visibilityOptions"
            item-title="label"
            item-value="value"
            label="Profile Visibility"
            density="compact"
            variant="outlined"
          >
            <template #selection="{ item }">
              <div class="d-flex align-center">
                <v-icon :icon="item.raw.icon" class="me-2" size="small"></v-icon>
                {{ item.raw.label }}
              </div>
            </template>
            <template #item="{ props, item }">
              <v-list-item v-bind="props">
                <template #prepend>
                  <v-icon :icon="item.raw.icon" size="small"></v-icon>
                </template>
                <v-list-item-subtitle class="text-wrap">
                  {{ item.raw.description }}
                </v-list-item-subtitle>
              </v-list-item>
            </template>
          </v-select>
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
import { onMounted, ref, computed } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRouter } from 'vue-router'
import type { User, UserSettings, apitypes_VisibilityLevel } from '@/genapi/api/users/user/v1alpha1'
import { useAlertStore } from '@/stores/alerts'

const router = useRouter()
const authStore = useAuthStore()
const { user, userSettings } = storeToRefs(authStore)
const breadcrumbStore = useBreadcrumbStore()
const alertStore = useAlertStore()

// Create a copy of the user object for editing
const editedUser = ref<User & UserSettings>({
  name: user.value.name,
  email: userSettings.value.email,
  username: user.value.username,
  givenName: user.value.givenName,
  familyName: user.value.familyName,
  visibility: (user.value.visibility || 'VISIBILITY_LEVEL_PUBLIC') as apitypes_VisibilityLevel,
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

function navigateBack() {
  router.push({ name: 'user', params: { userId: user.value.name } })
}

async function saveSettings() {
  try {
    await authStore.updateAuthUser(editedUser.value)
    navigateBack()
  } catch (err) {
    console.log('Error saving settings:', err)
    alertStore.addAlert(err instanceof Error ? "Unable to save settings\n" + err.message : String(err), 'error')
  }
}

onMounted(async () => {
  await authStore.loadAuthUser()

  breadcrumbStore.setBreadcrumbs([
    { title: 'User Settings', to: { name: 'user', params: { userId: user.value.name } } },
    { title: 'Edit', to: { name: 'user-edit', params: { userId: user.value.name } } },
  ])
})
</script>

<style></style>
