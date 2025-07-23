<template>
  <v-container v-if="user">
    <ListTabsPage :tabs="tabs" ref="tabsPage">
      <template #general>
        <v-card class="mx-auto" max-width="600">
          <div class="image-container">
            <v-img
              v-if="user.imageUri"
              class="mt-1"
              style="background-color: lightgray"
              :src="user.imageUri"
              cover
              height="300"
            ></v-img>
            <div v-else class="mt-1 d-flex align-center justify-center" style="background-color: lightgray; height: 300px; border-radius: 4px;">
              <div class="text-center">
                <v-icon size="64" color="grey-darken-1">mdi-image-outline</v-icon>
                <div class="text-grey-darken-1 mt-2">No image available</div>
              </div>
            </div>
          </div>
          <div class="bio-section pa-4">
            <div class="text-subtitle-1 font-weight-bold mb-1">Bio</div>
            <div class="text-body-1" style="white-space: pre-line;">
              {{ user.bio || 'No bio set.' }}
            </div>
          </div>
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
          </v-card-actions>
        </v-card>
      </template>
      <template #recipes="{ items, loading }">
        <RecipeGrid :recipes="items" :loading="loading" />
      </template>
    </ListTabsPage>
    <!-- Speed Dial -->
    <v-fab location="bottom right" app color="primary"  icon @click="speedDialOpen = !speedDialOpen">
      <v-icon>mdi-dots-vertical</v-icon>
      <v-speed-dial location="top" v-model="speedDialOpen" transition="slide-y-reverse-transition" activator="parent">
        <v-btn key="edit" v-if="hasAdminPermission(user.access?.permissionLevel)"
        :to="{ name: 'user-edit', params: { userId: user.name } }" color="primary"><v-icon>mdi-pencil</v-icon>Edit</v-btn>

        <v-btn key="share" v-if="!hasWritePermission(user.access?.permissionLevel)"
          @click="handleConnect" :loading="connecting" color="primary"><v-icon>mdi-share-variant</v-icon>Connect</v-btn>
        
        <v-btn key="remove-access" v-if="hasWritePermission(user.access?.permissionLevel) && user.access?.state === 'ACCESS_STATE_PENDING'"
          @click="showRemoveAccessDialog = true" color="warning"><v-icon>mdi-link-variant-off</v-icon>Remove Connection</v-btn>
      </v-speed-dial>
    </v-fab>

    <!-- Remove Access Dialog -->
  <v-dialog v-model="showRemoveAccessDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Remove Access
      </v-card-title>
      <v-card-text>
        Are you sure you want to remove your access to this user? You may no longer be able to view them.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showRemoveAccessDialog = false">
          Cancel
        </v-btn>
        <v-btn color="error" @click="handleRemoveAccess" :loading="removingAccess">
          Remove Access
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { useUsersStore } from '@/stores/users'
import { useRecipesStore } from '@/stores/recipes'
import { storeToRefs } from 'pinia'
import { onMounted, ref, computed } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRoute } from 'vue-router'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { hasWritePermission, hasAdminPermission } from '@/utils/permissions'
import { userAccessService } from '@/api/api'
import type { DeleteAccessRequest } from '@/genapi/api/meals/recipe/v1alpha1'
import type { CreateAccessRequest, Access, Access_User } from '@/genapi/api/users/user/v1alpha1'
import type { PermissionLevel } from '@/genapi/api/types'
import { useAlertStore } from '@/stores/alerts'

const usersStore = useUsersStore()
const alertsStore = useAlertStore()
const recipesStore = useRecipesStore()
const { currentUser: user } = storeToRefs(usersStore)
const breadcrumbStore = useBreadcrumbStore()
const route = useRoute()
const tabsPage = ref()
const speedDialOpen = ref(false)

const userId = computed(() => String(route.params.userId || ''))

onMounted(async () => {
  await Promise.all([
    usersStore.loadUser(userId.value),
  ])
  breadcrumbStore.setBreadcrumbs([
    { title: 'User Settings', to: { name: 'user', params: { userId: userId.value } } },
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
  return visibilityOptions.find(option => option.value === user.value?.visibility)
})
const selectedVisibilityLabel = computed(() => selectedVisibility.value?.label || '')
const selectedVisibilityDescription = computed(() => selectedVisibility.value?.description || '')
const selectedVisibilityIcon = computed(() => selectedVisibility.value?.icon || 'mdi-help-circle')

const tabs = [
  {
    label: 'General',
    value: 'general',
  },
  {
    label: 'Recipes',
    value: 'recipes',
    loader: async () => {
      if (!user.value?.name) return []
      // The parent resource for user recipes is the user's resource name
      await recipesStore.loadMyRecipes(user.value.name)
      return [...recipesStore.myRecipes]
    },
  },
]

const authStore = useAuthStore()
const router = useRouter()

// *** Remove Access ***
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

async function handleRemoveAccess() {
  if (!user.value?.access?.name || !authStore.activeAccount?.name) return

  removingAccess.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: user.value.access.name
    }
    
    await userAccessService.DeleteAccess(deleteRequest)
    router.push({ name: 'users' })
  } catch (error) {
    const msg = `Error removing access: ${error instanceof Error ? 
    error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  } finally {
    removingAccess.value = false
    showRemoveAccessDialog.value = false
  }
}

// *** Connect ***
const connecting = ref(false)

async function handleConnect() {
  if (!user.value?.name || !authStore.activeAccount?.name) return
  connecting.value = true
  try {
    // Build the access object
    const access: Access = {
      name: undefined, // Will be set by backend
      requester: undefined, // Will be set by backend
      recipient: {
        name: user.value.name,
        username: undefined,
        givenName: undefined,
        familyName: undefined,
      },
      level: 'PERMISSION_LEVEL_WRITE',
      state: undefined,
    }
    const req: CreateAccessRequest = {
      parent: user.value.name,
      access,
    }
    await userAccessService.CreateAccess(req)
    // Optionally reload user or show a message
    await usersStore.loadUser(user.value.name)
    alertsStore.addAlert('Connection request sent.', 'info')
  } catch (error) {
    const msg = `Error connecting: ${error instanceof Error ? 
    error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  } finally {
    connecting.value = false
  }
}
</script>

<style></style>
