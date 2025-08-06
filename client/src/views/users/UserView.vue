<template>
  <v-container v-if="user" class="pb-16">
    <ListTabsPage :tabs="tabs" ref="tabsPage">
      <template #general>
         <!-- Access Request Section -->
         <v-card v-if="showAccessRequest" class="mx-auto mt-4 mb-2" max-width="600">
          <v-card-title class="text-h6">
            <v-icon left>mdi-account-clock</v-icon>
            Pending Friend Request
          </v-card-title>
          <v-card-text>
            <p class="text-body-1 mb-4">
              You have a pending friend request from this user. Would you like to accept or decline?
            </p>
            
              <v-btn
                color="success"
                @click="handleAcceptRequest"
                :loading="acceptingRequest"
                block
              >
                <v-icon left>mdi-check</v-icon>
                Accept
              </v-btn>
              <v-btn
                color="error"
                @click="handleDeclineRequest"
                :loading="decliningRequest"
                block
              >
                <v-icon left>mdi-close</v-icon>
                Decline
            </v-btn>
          </v-card-text>
        </v-card>
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
            </v-list>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
          </v-card-actions>
        </v-card>
        <!-- Floating action buttons container -->
        <div class="fab-container">
          <v-btn density="compact" v-if="hasAdminPermission(user.access?.permissionLevel)" color="primary" :to="'/'+user.name+'/edit'">
            <v-icon>mdi-pencil</v-icon>Edit
          </v-btn>
          <v-btn density="compact" v-if="!hasWritePermission(user.access?.permissionLevel)" color="primary"  @click="handleConnect" :loading="connecting">
            <v-icon>mdi-account-plus</v-icon>Friend Request
          </v-btn>
          <v-btn density="compact" v-if="user.access?.state === 'ACCESS_STATE_PENDING'" color="warning" @click="showCancelRequestDialog = true">
            <v-icon>mdi-close</v-icon>Cancel Request
          </v-btn>
          <v-btn density="compact" v-if="hasWritePermission(user.access?.permissionLevel) && user.access?.state === 'ACCESS_STATE_ACCEPTED'" color="warning" @click="showRemoveAccessDialog = true">
            <v-icon>mdi-account-remove</v-icon>Remove Access
          </v-btn>
        </div>
      </template>
      <template #recipes="{ items, loading }">
        <RecipeGrid :recipes="items" :loading="loading" />
        <!-- Recipe Tab FAB -->
        <v-btn
          v-if="hasAdminPermission(user.access?.permissionLevel)"
          color="primary"
          density="compact"
          style="position: fixed; bottom: 16px; right: 16px"
          :to="{ name: 'recipeCreate' }"
        >
          <v-icon>mdi-plus</v-icon>
          Create Recipe
        </v-btn>
      </template>
      <template #friends="{ items, loading }">
        <UserGrid :users="items" :loading="loading" empty-text="No friends found." />
        <!-- Friends Tab FAB -->
        <v-btn
          v-if="hasAdminPermission(user.access?.permissionLevel)"
          color="primary"
          density="compact"
          style="position: fixed; bottom: 16px; right: 16px"
          @click="showShareDialog = true"
        >
          <v-icon>mdi-share-variant</v-icon>
          Manage Friends
        </v-btn>
      </template>
      <template #circles="{ items, loading }">
        <CircleGrid :circles="items" :loading="loading" empty-text="No circles found." />
        <!-- Circle Tab FAB -->
        <v-btn
          v-if="hasAdminPermission(user.access?.permissionLevel)"
          color="primary"
          density="compact"
          style="position: fixed; bottom: 16px; right: 16px"
          :to="{ name: 'circleCreate' }"
        >
          <v-icon>mdi-plus</v-icon>
          Create Circle
        </v-btn>
      </template>
    </ListTabsPage>

    <!-- Remove Access Dialog -->
  <v-dialog v-model="showRemoveAccessDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Disconnect Friend
      </v-card-title>
      <v-card-text>
        Are you sure you want to disconnect from this friend? You may no longer be able to view them.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showRemoveAccessDialog = false">
          Cancel
        </v-btn>
        <v-btn color="error" @click="handleRemoveAccess" :loading="removingAccess">
          Disconnect
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Cancel Request Dialog -->
  <v-dialog v-model="showCancelRequestDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Cancel Friend Request
      </v-card-title>
      <v-card-text>
        Are you sure you want to cancel your friend request to this user?
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showCancelRequestDialog = false">
          Cancel
        </v-btn>
        <v-btn color="error" @click="handleCancelRequest" :loading="cancelingRequest">
          Cancel Request
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Share Dialog for Friends -->
  <ShareDialog
    v-model="showShareDialog"
    title="Friends"
    :allowCircleShare="false"
    :currentAccesses="currentAccesses"
    :sharing="sharing"
    :sharePermissionLoading="updatingPermission"
    :userPermissionLevel="user.access?.permissionLevel"
    :allowPermissionOptions="allowPermissionOptions"
    :disableCreateShare="true"
    @share-user="shareWithUser"
    @remove-access="unshareUser"
    @permission-change="updatePermission"
  />
  </v-container>
</template>

<script setup lang="ts">
import { useUsersStore } from '@/stores/users'
import { useRecipesStore } from '@/stores/recipes'
import { storeToRefs } from 'pinia'
import { onMounted, ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { hasWritePermission, hasAdminPermission } from '@/utils/permissions'
import { userAccessService } from '@/api/api'
import type { DeleteAccessRequest, ListAccessesRequest } from '@/genapi/api/meals/recipe/v1alpha1'
import type { CreateAccessRequest, Access, Access_User } from '@/genapi/api/users/user/v1alpha1'
import { useAlertStore } from '@/stores/alerts'
import UserGrid from '@/components/UserGrid.vue'
import CircleGrid from '@/components/CircleGrid.vue'
import { useCirclesStore } from '@/stores/circles'
import ShareDialog from '@/components/common/ShareDialog.vue'
import type { PermissionLevel } from '@/genapi/api/types'

const usersStore = useUsersStore()
const alertsStore = useAlertStore()
const recipesStore = useRecipesStore()
const circlesStore = useCirclesStore()
const { currentUser: user } = storeToRefs(usersStore)
const route = useRoute()
const tabsPage = ref()
const speedDialOpen = ref(false)

const trimmedUserName = computed(() => {
  return route.path.substring(route.path.indexOf('users/'))
})

async function loadUser() {
  await Promise.all([
    usersStore.loadUser(trimmedUserName.value),
  ])
  tabsPage.value.activeTab = 'general'
}

onMounted(async () => {
  await loadUser()
})

watch(
  () => route.path,
  async (newUserName) => {
    if (newUserName) {
      await loadUser()
    }
  }
)

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
  {
    label: 'Friends',
    value: 'friends',
    loader: async () => {
      if (!user.value?.name) return []
      await usersStore.loadFriends(user.value.name)
      return [...usersStore.friends]
    },
  },
  {
    label: 'Circles',
    value: 'circles',
    loader: async () => {
      if (!user.value?.name) return []
      await circlesStore.loadMyCircles(user.value.name)
      return [...circlesStore.myCircles]
    },
  },
]

const authStore = useAuthStore()
const router = useRouter()

// Computed property to show access request section
const showAccessRequest = computed(() => {
  return user.value?.access?.state === 'ACCESS_STATE_PENDING' && 
         user.value?.access?.requester !== authStore.user?.name
})

// *** Remove Access ***
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

// *** Cancel Request ***
const showCancelRequestDialog = ref(false)
const cancelingRequest = ref(false)

// *** Accept/Decline Request ***
const acceptingRequest = ref(false)
const decliningRequest = ref(false)

async function handleRemoveAccess() {
  if (!user.value?.access?.name) return

  removingAccess.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: user.value.access.name
    }
    
    await userAccessService.DeleteAccess(deleteRequest)
    await usersStore.loadUser(user.value.name!)
    alertsStore.addAlert('Friend removed.', 'info')
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

async function handleCancelRequest() {
  if (!user.value?.access?.name) return

  cancelingRequest.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: user.value.access.name
    }
    
    await userAccessService.DeleteAccess(deleteRequest)
    // Reload the user to update the access state
    await usersStore.loadUser(user.value.name!)
    alertsStore.addAlert('Friend request cancelled.', 'info')
  } catch (error) {
    const msg = `Error cancelling friend request: ${error instanceof Error ? 
    error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  } finally {
    cancelingRequest.value = false
    showCancelRequestDialog.value = false
  }
}

async function handleAcceptRequest() {
  if (!user.value?.access?.name) return

  acceptingRequest.value = true
  try {
    await userAccessService.AcceptAccess({ name: user.value.access.name })
    // Reload the user to update the access state
    await usersStore.loadUser(user.value.name!)
    alertsStore.addAlert('Friend request accepted.', 'success')
  } catch (error) {
    const msg = `Error accepting friend request: ${error instanceof Error ? 
    error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  } finally {
    acceptingRequest.value = false
  }
}

async function handleDeclineRequest() {
  if (!user.value?.access?.name) return

  decliningRequest.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: user.value.access.name
    }
    
    await userAccessService.DeleteAccess(deleteRequest)
    // Reload the user to update the access state
    await usersStore.loadUser(user.value.name!)
    alertsStore.addAlert('Friend request declined.', 'info')
  } catch (error) {
    const msg = `Error declining friend request: ${error instanceof Error ? 
    error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  } finally {
    decliningRequest.value = false
  }
}

// *** Connect ***
const connecting = ref(false)

async function handleConnect() {
  if (!user.value?.name) return
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
    alertsStore.addAlert('Friend request sent.', 'info')
  } catch (error) {
    const msg = `Error connecting: ${error instanceof Error ? 
    error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  } finally {
    connecting.value = false
  }
}

// *** User Sharing ***
const allowPermissionOptions: PermissionLevel[] = [
  'PERMISSION_LEVEL_WRITE',
]

const showShareDialog = ref(false)
const currentAccesses = ref<Access[]>([]) 
const sharing = ref(false)
const updatingPermission = ref<Record<string, boolean>>({})
const unsharing = ref<Record<string, boolean>>({})

// Fetch recipients when share dialog is opened
watch(showShareDialog, (isOpen) => {
  if (isOpen && user.value && hasAdminPermission(user.value.access?.permissionLevel)) {
    fetchUserRecipients()
  }
})

// Fetch user recipients
async function fetchUserRecipients() {
  if (!user.value?.name) return

  try {
    const request: ListAccessesRequest = {
      parent: user.value.name,
      filter: undefined,
      pageSize: undefined,
      pageToken: undefined
    }

    const response = await userAccessService.ListAccesses(request)

    if (response.accesses) {
      currentAccesses.value = response.accesses.filter(access => {
        // Filter out the current user's own access to avoid showing it in the shares list
        const isCurrentUser = access.recipient?.name === authStore.user?.name
        return !isCurrentUser
      }).map(access => ({
        name: access.name || '',
        recipient: access.recipient,
        level: access.level || 'PERMISSION_LEVEL_READ',
        state: access.state || 'ACCESS_STATE_PENDING',
        requester: access.requester || undefined,
      }))
    }
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  }
}

async function updatePermission({ access, newLevel }: { access: any, newLevel: PermissionLevel }) {
  if (access.level === newLevel) return
  if (!access.name) return
  updatingPermission.value[access.name] = true
  try {
    // await userAccessService.UpdateAccess({
    //   access: {
    //     name: access.name,
    //     level: newLevel,
    //     state: undefined,
    //     recipient: undefined,
    //     requester: undefined,
    //   },
    //   updateMask: 'level',
    // })
    // // Update local state
    // access.level = newLevel
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    updatingPermission.value[access.name] = false
  }
}

async function unshareUser(accessName: string) {
  if (!accessName) return

  unsharing.value[accessName] = true
  try {
    const request: DeleteAccessRequest = {
      name: accessName 
    }
    await userAccessService.DeleteAccess(request)
    await fetchUserRecipients()
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    unsharing.value[accessName] = false
  }
}

async function shareWithUser({ userName, permission }: { userName: string, permission: PermissionLevel }) {
  if (!userName) return
  if (!user.value?.name) return

  sharing.value = true
  try {
    const access: Access = {
      recipient: {
        name: userName,
        username: undefined,
        givenName: undefined,
        familyName: undefined,
      },
      level: permission,
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: user.value.name,
      access
    }

    await userAccessService.CreateAccess(request)
    await fetchUserRecipients()
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    sharing.value = false
  }
}
</script>

<style scoped>
.fab-container {
  position: fixed;
  bottom: 16px;
  right: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 1000;
}
</style>
