<template>
  <v-container v-if="circle">
    <v-container max-width="600" class="pa-1">
      <v-row>
        <v-col align-self="auto" cols="12" sm="8">
          <div class="image-container">
            <v-img 
              v-if="circle.imageUri" 
              class="mt-1" 
              style="background-color: lightgray" 
              :src="circle.imageUri" 
              cover
              height="300"
            ></v-img>
            <div 
              v-else 
              class="mt-1 d-flex align-center justify-center"
              style="background-color: lightgray; height: 300px; border-radius: 4px;"
            >
              <div class="text-center">
                <v-icon size="64" color="grey-darken-1">mdi-image-outline</v-icon>
                <div class="text-grey-darken-1 mt-2">No image available</div>
              </div>
            </div>
          </div>
        </v-col>
      </v-row>
      <v-row>
        <v-col class="pt-5">
          <div class="text-h4">
            {{ circle.title }}
          </div>
        </v-col>
      </v-row>

      <!-- Visibility Section -->
      <v-row>
        <v-col cols="12">
          <div class="mt-4">
            <div v-if="selectedVisibilityDescription" class="mt-2">
              <v-alert
                :icon="selectedVisibilityIcon"
                density="compact"
                variant="tonal"
                :color="selectedVisibilityColor"
              >
                <div class="text-body-2">
                  <strong>{{ selectedVisibilityLabel }}:</strong> {{ selectedVisibilityDescription }}
                </div>
              </v-alert>
            </div>
          </div>
        </v-col>
      </v-row>
    </v-container>

    <!-- Speed Dial -->
    <v-fab location="bottom right" app color="primary" icon @click="speedDialOpen = !speedDialOpen">
      <v-icon>mdi-dots-vertical</v-icon>
      <v-speed-dial location="top" v-model="speedDialOpen" transition="slide-y-reverse-transition" activator="parent">
        <v-btn key="edit" v-if="hasWritePermission(circle.permission)" icon="mdi-pencil"
        @click="router.push({ name: 'circle-edit', params: { circleId: circle.name } })" color="primary"></v-btn>

        <v-btn key="share" v-if="hasWritePermission(circle.permission)" icon="mdi-share-variant"
          @click="showShareDialog = true" color="primary"></v-btn>
  
        <v-btn key="remove-access" icon="mdi-link-variant-off" @click="showRemoveAccessDialog = true" color="warning"></v-btn>
  
        <v-btn key="delete" v-if="hasWritePermission(circle.permission)" icon="mdi-delete"
          @click="showDeleteDialog = true" color="error"></v-btn>
      </v-speed-dial>
    </v-fab>
  </v-container>

  <!-- Remove Access Dialog -->
  <v-dialog v-model="showRemoveAccessDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Remove Access
      </v-card-title>
      <v-card-text>
        Are you sure you want to remove your access to this circle? You will no longer be able to view it.
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

  <!-- Share Dialog -->
  <v-dialog v-model="showShareDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Share Circle
      </v-card-title>
      <v-card-text>
        <v-text-field v-model="usernameInput" label="Enter Username" :rules="[validateUsername]"
          :prepend-inner-icon="getUsernameIcon" :color="getUsernameColor" :loading="isLoadingUsername"
          @update:model-value="handleUsernameInput"></v-text-field>
        <v-select
          v-model="selectedPermission"
          :items="permissionOptions"
          label="Permission Level"
          class="mt-2"
        ></v-select>
        <v-btn block color="primary" @click="shareCircle" :loading="sharing" :disabled="!isValidUsername" class="mt-2">
          Share with User
        </v-btn>

        <v-divider class="my-4"></v-divider>

        <div v-if="currentShares.length > 0">
          <div class="text-subtitle-1 mb-2">Current Shares</div>
          <v-list>
            <v-list-item v-for="share in currentShares" :key="share.id" :title="share.name" :subtitle="`${share.type} ${share.state === 'ACCESS_STATE_PENDING' ? '(Pending)' : ''}`">
              <template #append>
                <div class="d-flex align-center gap-2">
                  <v-chip size="small" :color="hasWritePermission(share.permission) ? 'primary' : 'grey'">
                    {{ hasWritePermission(share.permission) ? 'Read & Write' : 'Read Only' }}
                  </v-chip>
                  <v-chip v-if="share.state === 'ACCESS_STATE_PENDING'" size="small" color="warning" variant="outlined">
                    Pending
                  </v-chip>
                  <v-btn icon="mdi-delete" variant="text" @click="removeShare(share.id)"></v-btn>
                </div>
              </template>
            </v-list-item>
          </v-list>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showShareDialog = false">
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Delete Dialog -->
  <v-dialog v-model="showDeleteDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Delete Circle
      </v-card-title>
      <v-card-text>
        Are you sure you want to delete this circle? This action will also delete the circle for any users that
        can view it.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showDeleteDialog = false">
          Cancel
        </v-btn>
        <v-btn color="error" @click="handleDelete" :loading="deleting">
          Delete
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import type { apitypes_VisibilityLevel } from '@/genapi/api/circles/circle/v1alpha1'
import type { Access, CreateAccessRequest, ListAccessesRequest, DeleteAccessRequest } from '@/genapi/api/circles/circle/v1alpha1'
import type { PermissionLevel, AccessState } from '@/genapi/api/types'
import { useCirclesStore } from '@/stores/circles'
import { storeToRefs } from 'pinia'
import { onMounted, onBeforeUnmount, watch, computed, ref } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRoute, useRouter } from 'vue-router'
import { circleService, userService, circleAccessService } from '@/api/api'
import type { User, ListUsersRequest } from '@/genapi/api/users/user/v1alpha1'
import { useAuthStore } from '@/stores/auth'
import { hasWritePermission } from '@/utils/permissions'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const circlesStore = useCirclesStore()
const { circle } = storeToRefs(circlesStore)
const breadcrumbStore = useBreadcrumbStore()

// Visibility options with descriptions and icons
const visibilityOptions = [
  {
    value: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
    label: 'Public',
    icon: 'mdi-earth',
    color: 'success',
    description: 'Everyone can see this circle'
  },
  {
    value: 'VISIBILITY_LEVEL_RESTRICTED' as apitypes_VisibilityLevel,
    label: 'Restricted',
    icon: 'mdi-account-group',
    color: 'warning',
    description: 'Shared users and their connections can see this'
  },
  {
    value: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
    label: 'Private',
    icon: 'mdi-lock',
    color: 'info',
    description: 'Only specifically shared users can see this'
  },
  {
    value: 'VISIBILITY_LEVEL_HIDDEN' as apitypes_VisibilityLevel,
    label: 'Hidden',
    icon: 'mdi-eye-off',
    color: 'secondary',
    description: 'Only you can see this circle'
  }
]

// Computed properties for the selected visibility
const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === circle.value?.visibility)
})

const selectedVisibilityDescription = computed(() => {
  return selectedVisibility.value?.description || ''
})

const selectedVisibilityLabel = computed(() => {
  return selectedVisibility.value?.label || ''
})

const selectedVisibilityIcon = computed(() => {
  return selectedVisibility.value?.icon || 'mdi-help-circle'
})

const selectedVisibilityColor = computed(() => {
  return selectedVisibility.value?.color || 'primary'
})

async function loadAndSetBreadcrumbs(circleId: string | string[]) {
  await circlesStore.loadCircle(circleId as string)
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'circles' } },
    { title: circle.value?.title || 'Circle', to: { name: 'circle', params: { circleId: circle.value?.name } } },
  ])
}

onMounted(() => {
  loadAndSetBreadcrumbs(route.params.circleId)
})

watch(
  () => route.params.circleId,
  (newId) => {
    loadAndSetBreadcrumbs(newId)
  }
)

// Share dialog state
const showShareDialog = ref(false)
const usernameInput = ref('')
const selectedUser = ref<User | null>(null)
const sharing = ref(false)

// Validation states
const isValidUsername = ref(false)
const isLoadingUsername = ref(false)

// Current shares state
const currentShares = ref<Array<{ 
  id: string; 
  name: string; 
  type: 'user';
  permission: PermissionLevel;
  state: AccessState;
}>>([])

// Debounce timer
let usernameDebounceTimer: number | null = null

// Computed properties for icons and colors
const getUsernameIcon = computed(() => {
  if (isLoadingUsername.value) return 'mdi-loading'
  if (!usernameInput.value) return undefined
  return isValidUsername.value ? 'mdi-check-circle' : 'mdi-close-circle'
})

const getUsernameColor = computed(() => {
  if (isLoadingUsername.value) return undefined
  if (!usernameInput.value) return undefined
  return isValidUsername.value ? 'success' : 'error'
})

// Debounced API call
async function checkUsername(username: string) {
  if (!username) {
    isValidUsername.value = false
    selectedUser.value = null
    return
  }

  isLoadingUsername.value = true
  try {
    const request: ListUsersRequest = {
      filter: `username = "${username}"`,
      pageSize: 1,
      pageToken: undefined
    }
    const response = await userService.ListUsers(request)

    if (response.users?.length === 1 && response.users[0].name !== authStore.activeAccount?.name) {
      selectedUser.value = response.users[0]
      isValidUsername.value = true
    } else {
      selectedUser.value = null
      isValidUsername.value = false
    }
  } catch (error) {
    console.error('Error checking username:', error)
    selectedUser.value = null
    isValidUsername.value = false
  } finally {
    isLoadingUsername.value = false
  }
}

function handleUsernameInput(value: string) {
  if (usernameDebounceTimer) {
    clearTimeout(usernameDebounceTimer)
  }
  usernameDebounceTimer = window.setTimeout(() => {
    checkUsername(value)
  }, 300)
}

// Clean up timer when component is unmounted
onBeforeUnmount(() => {
  if (usernameDebounceTimer) {
    clearTimeout(usernameDebounceTimer)
  }
})

function validateUsername(value: string): boolean | string {
  if (!value) return true
  return true // Validation is now handled by the API call
}

// Permission options
const selectedPermission = ref<PermissionLevel>('PERMISSION_LEVEL_READ')

const permissionOptions = [
  { title: 'Read Only', value: 'PERMISSION_LEVEL_READ' as PermissionLevel },
  { title: 'Read & Write', value: 'PERMISSION_LEVEL_WRITE' as PermissionLevel },
]

// Share circle function
async function shareCircle() {
  if (!selectedUser.value) return
  if (!circle.value?.name) return

  sharing.value = true
  try {
    const access: Access = {
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      recipient: selectedUser.value?.name,
      level: selectedPermission.value,
      state: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: circle.value.name,
      access
    }

    await circleAccessService.CreateAccess(request)
    
    // Refresh the recipients list
    await fetchCircleRecipients()
    
    // Reset selections and inputs after sharing
    selectedUser.value = null
    usernameInput.value = ''
    isValidUsername.value = false
    selectedPermission.value = 'PERMISSION_LEVEL_READ' // Reset to default
  } catch (error) {
    console.error('Error sharing circle:', error)
    // You might want to show an error notification here
  } finally {
    sharing.value = false
  }
}

async function removeShare(shareId: string) {
  try {
    const request: DeleteAccessRequest = {
      name: shareId // shareId is actually the access name in the format circles/{circle}/accesses/{access}
    }
    
    await circleAccessService.DeleteAccess(request)
    
    // Remove from local state after successful API call
    currentShares.value = currentShares.value.filter(share => share.id !== shareId)
  } catch (error) {
    console.error('Error removing share:', error)
    // You might want to show an error notification here
  }
}

// Function to fetch circle recipients
async function fetchCircleRecipients() {
  if (!circle.value?.name) return

  try {
    const request: ListAccessesRequest = {
      parent: circle.value.name,
      filter: undefined,
      pageSize: undefined,
      pageToken: undefined
    }

    const response = await circleAccessService.ListAccesses(request)

    if (response.accesses) {
      currentShares.value = response.accesses.filter(access => {
        // Filter out the current user's own access to avoid showing it in the shares list
        const isCurrentUser = access.recipient === authStore.activeAccount?.name
        return !isCurrentUser
      }).map(access => {
        const recipientName = access.recipient || ''
        
        return {
          id: access.name || '',
          name: recipientName,
          type: 'user' as const,
          permission: access.level || 'PERMISSION_LEVEL_READ',
          state: access.state || 'ACCESS_STATE_PENDING'
        }
      })
    }
  } catch (error) {
    console.error('Error fetching circle recipients:', error)
  }
}

// Fetch recipients when circle is loaded or when share dialog is opened
watch(circle, (newCircle) => {
  if (newCircle) {
    fetchCircleRecipients()
  }
}, { immediate: true })

// Fetch recipients when share dialog is opened
watch(showShareDialog, (isOpen) => {
  if (isOpen && circle.value) {
    fetchCircleRecipients()
  }
})

// Remove access dialog state
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

async function handleRemoveAccess() {
  if (!circle.value?.name || !authStore.activeAccount?.name) return

  removingAccess.value = true
  try {
    // First, we need to find the current user's access to delete it
    const listRequest: ListAccessesRequest = {
      parent: circle.value.name,
      filter: undefined,
      pageSize: undefined,
      pageToken: undefined
    }
    
    const response = await circleAccessService.ListAccesses(listRequest)
    
    // Find the current user's access
    const userAccess = response.accesses?.find(access => 
      access.recipient === authStore.activeAccount?.name
    )
    
    if (userAccess?.name) {
      const deleteRequest: DeleteAccessRequest = {
        name: userAccess.name
      }
      
      await circleAccessService.DeleteAccess(deleteRequest)
      router.push({ name: 'circles' })
    } else {
      console.error('Could not find user access to remove')
    }
  } catch (error) {
    console.error('Error removing access:', error)
    alert(error instanceof Error ? error.message : String(error))
  } finally {
    removingAccess.value = false
    showRemoveAccessDialog.value = false
  }
}

// Speed dial state
const speedDialOpen = ref(false)
const showDeleteDialog = ref(false)
const deleting = ref(false)

async function handleDelete() {
  if (!circle.value?.name) return

  deleting.value = true
  try {
    await circleService.DeleteCircle({
      name: circle.value.name
    })
    router.push({ name: 'circles' })
  } catch (error) {
    console.error('Error deleting circle:', error)
    alert(error instanceof Error ? error.message : String(error))
  } finally {
    deleting.value = false
    showDeleteDialog.value = false
  }
}
</script>

<style scoped>
.image-container {
  position: relative;
}
</style>
