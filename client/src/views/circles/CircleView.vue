<template>
  <v-container v-if="circle" class="pb-16">
    <ListTabsPage
      :tabs="tabs"
      ref="tabsPage"
    >
      <template #general>
        <v-container max-width="600" class="pa-1">
          <v-row>
            <v-col class="pt-5">
              <div class="text-h4">
                {{ circle.title }}
              </div>
              <div class="bio-section pa-2">
                <div class="text-subtitle-1 font-weight-bold mb-1">Description</div>
                <div class="text-body-1" style="white-space: pre-line;">
                  {{ circle.description || 'No description set.' }}
                </div>
              </div>
              <div v-if="circle.handle" class="text-body-2 mt-1">
                <strong>Handle:</strong> <code>{{ circle.handle }}</code>
              </div>
            </v-col>
          </v-row>
          <v-row>
            <v-spacer></v-spacer>
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
            <v-spacer></v-spacer>
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
          <!-- Accept/Decline Buttons for Pending Access -->
          <v-row v-if="circle.circleAccess?.acceptTarget === 'ACCEPT_TARGET_RECIPIENT' && circle.circleAccess?.state === 'ACCESS_STATE_PENDING'">
            <v-col cols="12">
              <v-btn
                color="success"
                class="mb-2"
                block
                :loading="acceptingCircle"
                @click="acceptCircle(circle.circleAccess?.name)"
              >
                Accept Circle
              </v-btn>
              <v-btn
                color="error"
                class="mb-4"
                block
                :loading="decliningCircle"
                @click="declineCircle"
              >
                Decline
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
        <!-- Floating action buttons container -->
        <div class="fab-container">
          <v-btn density="compact" v-if="hasWritePermission(circle.circleAccess?.permissionLevel)" color="primary" :to="'/'+circle.name+'/edit'">
            <v-icon>mdi-pencil</v-icon>Edit
          </v-btn>
          <v-btn density="compact" v-if="!hasReadPermission(circle.circleAccess?.permissionLevel) && !circle.circleAccess" color="primary" @click="handleRequestAccess" :loading="requestingAccess">
            <v-icon>mdi-account-plus</v-icon>Request Access
          </v-btn>
          <v-btn density="compact" v-if="circle.circleAccess?.state === 'ACCESS_STATE_PENDING'" color="warning" @click="showCancelRequestDialog = true">
            <v-icon>mdi-close</v-icon>Cancel Request
          </v-btn>
          <v-btn density="compact" v-if="!hasAdminPermission(circle.circleAccess?.permissionLevel) && circle.circleAccess?.state === 'ACCESS_STATE_ACCEPTED'" color="warning" @click="showRemoveAccessDialog = true">
            <v-icon>mdi-link-variant-off</v-icon>Remove Access
          </v-btn>
          <v-btn density="compact" v-if="hasAdminPermission(circle.circleAccess?.permissionLevel)" color="error" @click="showDeleteDialog = true">
            <v-icon>mdi-delete</v-icon>Delete
          </v-btn>
        </div>
      </template>
      <template #recipes-circleRecipes="{ items, loading }">
        <RecipeGrid :recipes="items" :loading="loading" />
        <v-btn
          v-if="hasWritePermission(circle.circleAccess?.permissionLevel)"
          color="primary"
          density="compact"
          style="position: fixed; bottom: 16px; right: 16px"
          :to="'/'+circle.name+'/recipes/create'"
        >
          <v-icon>mdi-plus</v-icon>
          Create Recipe
        </v-btn>
      </template>
      <template #recipes-pending="{ items, loading }">
        <RecipeGrid :recipes="items" :loading="loading" />
      </template>
      <template #recipes="{ items, loading }">
        <RecipeGrid :recipes="items" :loading="loading" />
      </template>
      <template #members="{ items, loading }">
        <UserGrid :users="items" :loading="loading" empty-text="No members found." />
        <v-btn
          v-if="hasWritePermission(circle.circleAccess?.permissionLevel)"
          color="primary"
          density="compact"
          style="position: fixed; bottom: 16px; right: 16px"
          @click="showShareDialog = true"
        >
          <v-icon>mdi-share-variant</v-icon>
          Manage Members
        </v-btn>
      </template>
    </ListTabsPage>

    <!-- Remove Access Dialog -->
    <v-dialog v-model="showRemoveAccessDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">
          Remove Access
        </v-card-title>
        <v-card-text>
          Are you sure you want to remove your access to this circle? You may no longer be able to view it.
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
    <ShareDialog
      v-model="showShareDialog"
      title="Share Circle"
      :allowCircleShare="false"
      :currentAccesses="currentAccesses"
      :sharing="sharing"
      :sharePermissionLoading="updatingPermission"
      :userPermissionLevel="circle.circleAccess?.permissionLevel"
      :allowPermissionOptions="allowPermissionOptions"
      @share-user="shareWithUser"
      @remove-access="unshareCircle"
      @permission-change="updatePermission"
      @approve-access="acceptCircleFromShareDialog"
    />
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

    <!-- Cancel Request Dialog -->
    <v-dialog v-model="showCancelRequestDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">
          Cancel Access Request
        </v-card-title>
        <v-card-text>
          Are you sure you want to cancel your access request to this circle?
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
  </v-container>
</template>

<script setup lang="ts">
import type { apitypes_VisibilityLevel, Access, CreateAccessRequest, ListAccessesRequest, DeleteAccessRequest } from '@/genapi/api/circles/circle/v1alpha1'
import type { PermissionLevel } from '@/genapi/api/types'
import { useCirclesStore } from '@/stores/circles'
import { useUsersStore } from '@/stores/users'
import { storeToRefs } from 'pinia'
import { onMounted,  watch, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { circleService, circleAccessService } from '@/api/api'
import { useAuthStore } from '@/stores/auth'
import { hasAdminPermission, hasWritePermission, hasReadPermission } from '@/utils/permissions'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import UserGrid from '@/components/UserGrid.vue'
import { useRecipesStore } from '@/stores/recipes'
import ShareDialog from '@/components/common/ShareDialog.vue'
import { useAlertStore } from '@/stores/alerts'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const recipesStore = useRecipesStore()
const circlesStore = useCirclesStore()
const usersStore = useUsersStore()
const alertsStore = useAlertStore()

const { circle } = storeToRefs(circlesStore)
const { user } = storeToRefs(authStore)
const speedDialOpen = ref(false)


const trimmedCircleName = computed(() => {
  return route.path.substring(route.path.indexOf('circles/'))
})

onMounted(() => {
  circlesStore.loadCircle(trimmedCircleName.value)
})

watch(
  () => route.path,
  () => {
    circlesStore.loadCircle(trimmedCircleName.value)
  }
)

// *** Tabs ***

const tabsPage = ref()
const tabs = computed(() => [
  {
    label: 'General',
    value: 'general',
  },
  {
    label: 'Recipes',
    value: 'recipes',
    loader: async () => {
      if (!circle.value?.name) return []
      // Admin recipes for this circle
      await recipesStore.loadMyRecipes(circle.value.name)
      return [...recipesStore.myRecipes]
    },
    subTabs: hasWritePermission(circle.value?.circleAccess?.permissionLevel) ? [
      {
        label: 'Circle Recipes',
        value: 'circleRecipes',
        loader: async () => {
          if (!circle.value?.name) return []
          // Admin recipes for this circle
          await recipesStore.loadMyRecipes(circle.value.name)
          return [...recipesStore.myRecipes]
        },
      },
      {
        label: 'Pending Recipes',
        value: 'pending',
        loader: async () => {
          if (!circle.value?.name) return []
          // Non-admin, non-pending recipes for this circle
          await recipesStore.loadPendingRecipes(circle.value.name)
          return [...recipesStore.sharedAcceptedRecipes]
        },
      },
    ] : undefined,
  },
  {
    label: 'Members',
    value: 'members',
    loader: async () => {
      if (!circle.value?.name) return []
      await usersStore.loadFriends(circle.value.name)
      return [...usersStore.friends]
    },
  },
])

// *** Visibility ***

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

// *** Remove Access ***
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

async function handleRemoveAccess() {
  if (!circle.value?.circleAccess?.name) return

  removingAccess.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: circle.value.circleAccess.name
    }
    
    await circleAccessService.DeleteAccess(deleteRequest)
    router.push({ name: 'circles' })
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    removingAccess.value = false
    showRemoveAccessDialog.value = false
  }
}

// *** Delete Circle ***

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
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    deleting.value = false
    showDeleteDialog.value = false
  }
}

// *** Circle Accept/Decline ***

const acceptingCircle = ref(false)
const decliningCircle = ref(false)

// *** Access Request ***
const requestingAccess = ref(false)
const showCancelRequestDialog = ref(false)
const cancelingRequest = ref(false)

async function handleRequestAccess() {
  if (!circle.value?.name) return
  if (!authStore.user?.name) return
  
  requestingAccess.value = true
  try {
    const access: Access = {
      recipient: {
        name: authStore.user?.name,
        username: undefined,
        givenName: undefined,
        familyName: undefined,
      },
      level: 'PERMISSION_LEVEL_READ',
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: circle.value.name,
      access
    }

    await circleAccessService.CreateAccess(request)
    await circlesStore.loadCircle(circle.value.name)
    alertsStore.addAlert('Access request sent.', 'info')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    requestingAccess.value = false
  }
}

async function handleCancelRequest() {
  if (!circle.value?.circleAccess?.name) return

  cancelingRequest.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: circle.value.circleAccess.name
    }
    
    await circleAccessService.DeleteAccess(deleteRequest)
    await circlesStore.loadCircle(circle.value.name!)
    alertsStore.addAlert('Access request cancelled.', 'info')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    cancelingRequest.value = false
    showCancelRequestDialog.value = false
  }
}

async function acceptCircleFromShareDialog(circleAccessName: string | undefined) {
  if (!circle.value?.circleAccess?.name) return
  if (!circleAccessName) return
  await acceptCircle(circleAccessName)
  await fetchCircleRecipients()
}

async function acceptCircle(circleAccessName: string | undefined) {
  if (!circle.value?.circleAccess?.name) return
  if (!circleAccessName) return
  
  acceptingCircle.value = true
  try {
    await circleAccessService.AcceptAccess({ name: circleAccessName })
    await circlesStore.loadCircle(circle.value.name!)
    alertsStore.addAlert('Access accepted.', 'info')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    acceptingCircle.value = false
  }
}

async function declineCircle() {
  if (!circle.value?.circleAccess?.name) return
  decliningCircle.value = true
  try {
    await circleAccessService.DeleteAccess({ name: circle.value.circleAccess.name })
    router.push({ name: 'circles' })
  } catch (error) {
    // Optionally show a notification
  } finally {
    decliningCircle.value = false
  }
}

// *** Circle Sharing ***

const allowPermissionOptions: PermissionLevel[] = [
  'PERMISSION_LEVEL_READ',
  'PERMISSION_LEVEL_WRITE',
  'PERMISSION_LEVEL_ADMIN',
]
const showShareDialog = ref(false)
const currentAccesses = ref<Access[]>([]) 
const sharing = ref(false)
const updatingPermission = ref<Record<string, boolean>>({})
const unsharing = ref<Record<string, boolean>>({})

// Fetch recipients when share dialog is opened
watch(showShareDialog, (isOpen) => {
  if (isOpen && circle.value && hasWritePermission(circle.value.circleAccess?.permissionLevel)) {
    fetchCircleRecipients()
  }
})

// Fetch circle recipients
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
      currentAccesses.value = response.accesses.filter(access => {
        // Filter out the current user's own access to avoid showing it in the shares list
        return access.recipient?.name !== user.value.name
      }).map(access => ({
        name: access.name || '',
        recipient: access.recipient,
        level: access.level || 'PERMISSION_LEVEL_READ',
        state: access.state || 'ACCESS_STATE_PENDING',
        requester: access.requester || undefined,
        acceptTarget: access.acceptTarget || 'ACCEPT_TARGET_UNSPECIFIED',
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
    await circleAccessService.UpdateAccess({
      access: {
        name: access.name,
        level: newLevel,
        state: undefined,
        recipient: undefined,
        requester: undefined,
        acceptTarget: undefined,
      },
      updateMask: 'level',
    })
    // Update local state
    access.level = newLevel
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    updatingPermission.value[access.name] = false
  }
}

async function unshareCircle(accessName: string) {
  if (!accessName) return

  unsharing.value[accessName] = true
  try {
    const request: DeleteAccessRequest = {
      name: accessName 
    }
    await circleAccessService.DeleteAccess(request)
    await fetchCircleRecipients()
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    unsharing.value[accessName] = false
  }
}

async function shareWithUser({ userName, permission }: { userName: string, permission: PermissionLevel }) {
  if (!userName) return
  if (!circle.value?.name) return

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
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: circle.value.name,
      access
    }

    await circleAccessService.CreateAccess(request)
    await fetchCircleRecipients()
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    sharing.value = false
  }
}


</script>

<style scoped>
.image-container {
  position: relative;
}

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
