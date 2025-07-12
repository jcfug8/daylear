<template>
  <v-container v-if="circle">
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
            </v-col>
          </v-row>
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
          <v-row v-if="circle.circleAccess?.state === 'ACCESS_STATE_PENDING'">
            <v-col cols="12">
              <v-btn
                color="success"
                class="mb-2"
                block
                :loading="acceptingCircle"
                @click="acceptCircle"
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
      </template>
      <template #recipes-circleRecipes="{ items, loading }">
        <RecipeGrid :recipes="items" :loading="loading" />
      </template>
      <template #recipes-sharedRecipes="{ items, loading }">
        <RecipeGrid :recipes="items" :loading="loading" />
      </template>
    </ListTabsPage>

    <!-- Speed Dial and Dialogs remain outside the tabbed area -->
    <v-fab location="bottom right" app color="primary" icon @click="speedDialOpen = !speedDialOpen">
      <v-icon>mdi-dots-vertical</v-icon>
      <v-speed-dial location="top" v-model="speedDialOpen" transition="slide-y-reverse-transition" activator="parent">
        <v-btn key="edit" v-if="hasWritePermission(circle.circleAccess?.permissionLevel)"
        @click="router.push({ name: 'circle-edit', params: { circleId: circle.name } })" color="primary"><v-icon>mdi-pencil</v-icon>Edit</v-btn>
        <v-btn key="share" v-if="hasWritePermission(circle.circleAccess?.permissionLevel) && circle.visibility !== 'VISIBILITY_LEVEL_HIDDEN'"
          @click="showShareDialog = true" color="primary"><v-icon>mdi-share-variant</v-icon>Share</v-btn>
        <v-btn key="remove-access" v-if="!hasAdminPermission(circle.circleAccess?.permissionLevel)" @click="showRemoveAccessDialog = true" color="warning"><v-icon>mdi-link-variant-off</v-icon>Remove Access</v-btn>
        <v-btn key="delete" v-if="hasWritePermission(circle.circleAccess?.permissionLevel)"
          @click="showDeleteDialog = true" color="error"><v-icon>mdi-delete</v-icon>Delete</v-btn>
      </v-speed-dial>
    </v-fab>
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
    <ShareDialog
      v-model="showShareDialog"
      title="Share Circle"
      :allowCircleShare="false"
      :currentShares="currentShares"
      :sharing="sharing"
      :sharePermissionLoading="updatingPermission"
      :hasWritePermission="hasWritePermission"
      @share-user="shareCircle"
      @remove-share="unshareCircle"
      @permission-change="updatePermission"
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
  </v-container>
</template>

<script setup lang="ts">
import type { apitypes_VisibilityLevel, Access, CreateAccessRequest, ListAccessesRequest, DeleteAccessRequest } from '@/genapi/api/circles/circle/v1alpha1'
import type { PermissionLevel } from '@/genapi/api/types'
import { useCirclesStore } from '@/stores/circles'
import { storeToRefs } from 'pinia'
import { onMounted,  watch, computed, ref } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRoute, useRouter } from 'vue-router'
import { circleService, circleAccessService } from '@/api/api'
import { useAuthStore } from '@/stores/auth'
import { hasAdminPermission, hasWritePermission } from '@/utils/permissions'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import { useRecipesStore } from '@/stores/recipes'
import ShareDialog from '@/components/common/ShareDialog.vue'

const route = useRoute()
const router = useRouter()
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()
const recipesStore = useRecipesStore()
const circlesStore = useCirclesStore()

const { circle } = storeToRefs(circlesStore)
const speedDialOpen = ref(false)

// *** Breadcrumbs ***

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

// *** Tabs ***

const tabsPage = ref()
const tabs = [
  {
    label: 'General',
    value: 'general',
  },
  {
    label: 'Recipes',
    value: 'recipes',
    subTabs: [
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
        label: 'Shared Recipes',
        value: 'sharedRecipes',
        loader: async () => {
          if (!circle.value?.name) return []
          // Non-admin, non-pending recipes for this circle
          await recipesStore.loadSharedRecipes(circle.value.name, 200)
          return [...recipesStore.sharedAcceptedRecipes]
        },
      },
    ],
  },
]

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
      access.recipient?.name === authStore.activeAccount?.name
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
    console.error('Error deleting circle:', error)
    alert(error instanceof Error ? error.message : String(error))
  } finally {
    deleting.value = false
    showDeleteDialog.value = false
  }
}

// *** Circle Accept/Decline ***

const acceptingCircle = ref(false)
const decliningCircle = ref(false)

async function acceptCircle() {
  if (!circle.value?.circleAccess?.name) return
  acceptingCircle.value = true
  try {
    await circleAccessService.AcceptAccess({ name: circle.value.circleAccess.name })
    await circlesStore.loadCircle(circle.value.name!)
  } catch (error) {
    // Optionally show a notification
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

const showShareDialog = ref(false)
const currentShares = ref<Access[]>([]) 
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
      currentShares.value = response.accesses.filter(access => {
        // Filter out the current user's own access to avoid showing it in the shares list
        const isCurrentUser = access.recipient?.name === authStore.activeAccount?.name
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
    console.error('Error fetching circle recipients:', error)
  }
}

async function updatePermission({ access, newLevel }: { access: Access, newLevel: PermissionLevel }) {
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
      },
      updateMask: 'level',
    })
    // Update local state
    access.level = newLevel
  } catch (error) {
    console.error('Error updating permission:', error)
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
    console.error('Error removing share:', error)
  } finally {
    unsharing.value[accessName] = false
  }
}

async function shareCircle({ userName, permission }: { userName: string, permission: PermissionLevel }) {
  if (!userName) return
  if (!circle.value?.name) return

  sharing.value = true
  try {
    const access: Access = {
      recipient: {
        name: userName,
        username: undefined
      },
      level: permission,
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: circle.value.name,
      access
    }

    await circleAccessService.CreateAccess(request)
    await fetchCircleRecipients()
  } catch (error) {
    console.error('Error sharing circle:', error)
  } finally {
    sharing.value = false
  }
}


</script>

<style scoped>
.image-container {
  position: relative;
}
</style>
