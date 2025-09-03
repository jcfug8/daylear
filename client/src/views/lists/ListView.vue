<template>
  <div v-if="list" style="padding-bottom: 84px;">
    <v-container max-width="600" class="pa-1">
      <v-row>
        <v-col class="pt-5">
          <div class="text-h4">
            {{ list.title }}
          </div>
          <div class="text-body-1">
            {{ list.description }}
          </div>
        </v-col>
      </v-row>

      <!-- List Settings -->
      <v-row>
        <v-col cols="12">
          <div v-if="list.showCompleted !== undefined">
            <strong>Show Completed Items:</strong> {{ list.showCompleted ? 'Yes' : 'No' }}
          </div>
          <div v-if="list.sections && list.sections.length">
            <strong>Sections:</strong> {{ list.sections.length }} section(s)
          </div>
        </v-col>
      </v-row>

      <!-- Visibility Section -->
      <v-row>
        <v-col cols="12">
          <div>
            <div v-if="selectedVisibilityDescription">
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

      <!-- List Items and Sections -->
      <v-row>
        <v-col cols="12">
          <!-- Items Without Section -->
          <div v-if="itemsWithoutSection.length > 0 || (hasWritePermission(list.listAccess?.permissionLevel) && !isCreatingNewItem)">
            <h3 class="text-h6 mt-6 mb-4">Items</h3>
            
            <!-- Existing Items Without Section -->
            <div v-for="(item, itemIndex) in itemsWithoutSection" :key="itemIndex" class="mb-2">
              <div class="d-flex align-center">
                <v-checkbox class="mr-2" hide-details density="compact" />
                <span>{{ item.title }}</span>
                <span v-if="item.points && item.points > 0" class="text-caption text-grey ml-2">
                  - {{ item.points }} points
                </span>
              </div>
            </div>
            
            <!-- Add New Item (without section) -->
            <div v-if="hasWritePermission(list.listAccess?.permissionLevel)" class="mt-4">
              <div v-if="!isCreatingNewItem" class="d-flex align-center text-grey cursor-pointer" @click="startCreatingItem">
                <v-icon class="mr-2">mdi-plus</v-icon>
                <span>Add new item...</span>
              </div>
              
              <v-text-field
                v-else
                ref="newItemInputRef"
                v-model="newItemText"
                placeholder="Enter item title..."
                variant="outlined"
                density="compact"
                @keydown="handleItemInputKeydown"
                @blur="handleItemInputBlur"
                autofocus
              />
            </div>
          </div>

          <!-- New Section Button -->
          <div v-if="hasWritePermission(list.listAccess?.permissionLevel)" class="mb-4">
            <v-btn
              v-if="!isCreatingNewSection"
              color="primary"
              variant="outlined"
              prepend-icon="mdi-plus"
              @click="startCreatingSection"
            >
              New Section
            </v-btn>
            
            <!-- New Section Input -->
            <v-text-field
              v-else
              ref="newSectionInputRef"
              v-model="newSectionText"
              placeholder="Enter section title..."
              variant="outlined"
              density="compact"
              @blur="handleSectionInputBlur"
              autofocus
            />
          </div>

          <!-- Sections with their Items -->
          <div v-for="(section, sectionIndex) in list.sections" :key="sectionIndex" class="mb-6">
            <div class="text-h6 mb-3">{{ section.title || 'Untitled Section' }}</div>
            
            <!-- Items in this section -->
            <div v-if="itemsBySection[section.name || ''] && itemsBySection[section.name || ''].length > 0" class="ml-4">
              <div v-for="(item, itemIndex) in itemsBySection[section.name || '']" :key="itemIndex" class="mb-2">
                <div class="d-flex align-center">
                  <v-checkbox class="mr-2" hide-details density="compact" />
                  <span>{{ item.title }}</span>
                  <span v-if="item.points && item.points > 0" class="text-caption text-grey ml-2">
                    - {{ item.points }} points
                  </span>
                </div>
              </div>
            </div>
            
            <!-- Empty section message -->
            <div v-else class="ml-4 text-body-2 text-grey">
              No items in this section yet.
            </div>
          </div>
        </v-col>
      </v-row>
    </v-container>

    <v-icon 
      size="24" 
      :color="list.favorited ? 'red' : 'black'"
      class="favorite-heart"
      @click="toggleFavorite"
      style="cursor: pointer;"
    >
      {{ list.favorited ? 'mdi-heart' : 'mdi-heart-outline' }}
    </v-icon>

    <!-- Floating action buttons container -->
    <div class="fab-container">
      <v-btn v-if="hasWritePermission(list.listAccess?.permissionLevel)" color="primary" density="compact" :to="'/'+list.name+'/edit'">
        <v-icon>mdi-pencil</v-icon>
        Edit
      </v-btn>
      <v-btn density="compact" v-if="!hasReadPermission(list.listAccess?.permissionLevel) && !list.listAccess" color="primary" @click="handleRequestAccess" :loading="requestingAccess">
        <v-icon>mdi-account-plus</v-icon>Request Access
      </v-btn>
      <v-btn v-if="hasWritePermission(list.listAccess?.permissionLevel) && list.visibility !== 'VISIBILITY_LEVEL_HIDDEN'" color="primary" density="compact" @click="showShareDialog = true">
        <v-icon>mdi-share-variant</v-icon>
        Manage Access
      </v-btn>
      <v-btn density="compact" v-if="list.listAccess?.state === 'ACCESS_STATE_PENDING' && list.listAccess?.acceptTarget !== 'ACCEPT_TARGET_RECIPIENT'" color="warning" @click="showCancelRequestDialog = true">
        <v-icon>mdi-close</v-icon>Cancel Request
      </v-btn>
      <v-btn v-if="!hasAdminPermission(list.listAccess?.permissionLevel) && list.listAccess?.state === 'ACCESS_STATE_ACCEPTED'" color="warning" density="compact" @click="showRemoveAccessDialog = true">
        <v-icon>mdi-link-variant-off</v-icon>
        Remove Access
      </v-btn>
      <v-btn v-if="hasAdminPermission(list.listAccess?.permissionLevel)" color="error" density="compact" @click="showDeleteDialog = true">
        <v-icon>mdi-delete</v-icon>
        Delete
      </v-btn>
    </div>
    
    <template v-if="list.listAccess?.acceptTarget === 'ACCEPT_TARGET_RECIPIENT' && list.listAccess?.state === 'ACCESS_STATE_PENDING'">
      <v-btn
        color="success"
        class="mb-2"
        block
        :loading="acceptingList"
        @click="acceptList(list.listAccess?.name)"
      >
        Accept List
      </v-btn>
      <v-btn
        color="error"
        class="mb-4"
        block
        :loading="decliningList"
        @click="declineList"
      >
        Decline
      </v-btn>
    </template>

    <!-- Share Dialog -->
    <ShareDialog
      v-model="showShareDialog"
      title="Share List"
      :allowCircleShare="true"
      :currentAccesses="currentAccesses"
      :sharing="sharing"
      :sharePermissionLoading="updatingPermission"
      :userPermissionLevel="list.listAccess?.permissionLevel"
      :allowPermissionOptions="allowPermissionOptions"
      @share-user="shareWithUser"
      @share-circle="shareWithCircle"
      @remove-access="unshareList"
      @permission-change="updatePermission"
      @approve-access="acceptListFromShareDialog"
    />
  </div>
  
  <!-- Remove Access Dialog -->
  <v-dialog v-model="showRemoveAccessDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Remove Access
      </v-card-title>
      <v-card-text>
        Are you sure you want to remove your access to this list? You may no longer be able to view it.
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

  <!-- Delete Dialog -->
  <v-dialog v-model="showDeleteDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Delete List
      </v-card-title>
      <v-card-text>
        Are you sure you want to delete this list? This action will also delete the list for any users or circles
        that can view it.
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
        Are you sure you want to cancel your access request to this list?
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
</template>

<script setup lang="ts">
import type { apitypes_VisibilityLevel } from '@/genapi/api/lists/list/v1alpha1'
import type { Access, CreateAccessRequest, ListAccessesRequest, DeleteAccessRequest } from '@/genapi/api/lists/list/v1alpha1'
import type { ListItem, CreateListItemRequest, ListListItemsRequest, ListListItemsResponse } from '@/genapi/api/lists/list/v1alpha1'
import type { PermissionLevel } from '@/genapi/api/types'
import { useListStore } from '@/stores/list'
import { storeToRefs } from 'pinia'
import { onMounted, watch, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { listService, listAccessService, listItemService } from '@/api/api'
import { useAuthStore } from '@/stores/auth'
import { hasWritePermission, hasAdminPermission, hasReadPermission } from '@/utils/permissions'
import ShareDialog from '@/components/common/ShareDialog.vue'
import { useAlertStore } from '@/stores/alerts'

const authStore = useAuthStore()
const listStore = useListStore()
const route = useRoute()
const router = useRouter()
const alertsStore = useAlertStore()

const trimmedListName = computed(() => {
  return route.path.substring(route.path.indexOf('lists/'))
})

const { list } = storeToRefs(listStore)

// List items state
const listItems = ref<ListItem[]>([])
const loadingListItems = ref(false)

// Computed properties for grouping items
const itemsWithoutSection = computed(() => {
  return listItems.value.filter(item => !item.listSection)
})

const itemsBySection = computed(() => {
  const grouped: Record<string, ListItem[]> = {}
  
  listItems.value.forEach(item => {
    if (item.listSection) {
      if (!grouped[item.listSection]) {
        grouped[item.listSection] = []
      }
      grouped[item.listSection].push(item)
    }
  })
  
  return grouped
})

// Editing state
const isCreatingNewItem = ref(false)
const isCreatingNewSection = ref(false)
const newItemText = ref('')
const newSectionText = ref('')
const newItemInputRef = ref()
const newSectionInputRef = ref()

// *** Visibility ***

const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === list.value?.visibility)
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

// Visibility options with descriptions and icons
const visibilityOptions = [
  {
    value: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
    label: 'Public',
    icon: 'mdi-earth',
    color: 'success',
    description: 'Everyone can see this list'
  },
  {
    value: 'VISIBILITY_LEVEL_RESTRICTED' as apitypes_VisibilityLevel,
    label: 'Restricted',
    icon: 'mdi-account-group',
    color: 'warning',
    description: 'Shared users, circles and their connections can see this'
  },
  {
    value: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
    label: 'Private',
    icon: 'mdi-lock',
    color: 'info',
    description: 'Only specifically shared users and circles can see this'
  },
  {
    value: 'VISIBILITY_LEVEL_HIDDEN' as apitypes_VisibilityLevel,
    label: 'Hidden',
    icon: 'mdi-eye-off',
    color: 'secondary',
    description: 'Only you can see this list'
  }
]

// *** Remove Access ***
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

async function handleRemoveAccess() {
  if (!list.value?.listAccess?.name) return

  removingAccess.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: list.value.listAccess.name
    }
    
    await listAccessService.DeleteAccess(deleteRequest)
    router.push(route.params.circleId ? { name: 'circle', params: { circleId: route.params.circleId } } : { name: 'lists' })
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    removingAccess.value = false
    showRemoveAccessDialog.value = false
  }
}

// *** Cancel Request ***
const showCancelRequestDialog = ref(false)
const cancelingRequest = ref(false)

async function handleCancelRequest() {
  if (!list.value?.listAccess?.name) return

  cancelingRequest.value = true
  try {
    const deleteRequest: DeleteAccessRequest = {
      name: list.value.listAccess.name
    }
    
    await listAccessService.DeleteAccess(deleteRequest)
    await listStore.loadList(list.value.name!)
    alertsStore.addAlert('Access request cancelled.', 'info')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    cancelingRequest.value = false
    showCancelRequestDialog.value = false
  }
}

// *** Delete List ***
const showDeleteDialog = ref(false)
const deleting = ref(false)

async function handleDelete() {
  if (!list.value?.name) return

  deleting.value = true
  try {
    await listService.DeleteList({
      name: list.value.name
    })
    router.push({ name: 'lists' })
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    deleting.value = false
    showDeleteDialog.value = false
  }
}

// *** Accept/Decline List ***
const acceptingList = ref(false)
const decliningList = ref(false)

async function acceptListFromShareDialog(listAccessName: string | undefined) {
  if (!list.value?.listAccess?.name) return
  if (!listAccessName) return
  await acceptList(listAccessName)
  await fetchListRecipients()
}

async function acceptList(listAccessName: string | undefined) {
  if (!list.value?.listAccess?.name || !authStore.user?.name) return
  if (!listAccessName) return
  acceptingList.value = true
  try {
    await listStore.acceptList(listAccessName)
    listStore.loadList(list.value?.name ?? '')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    acceptingList.value = false
  }
}

async function declineList() {
  if (!list.value?.listAccess?.name) return
  decliningList.value = true
  try {
    await listStore.deleteListAccess(list.value.listAccess.name)
    router.push(route.params.circleId ? { name: 'circle', params: { circleId: route.params.circleId } } : { name: 'lists' })
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  } finally {
    decliningList.value = false
  }
}

// *** List Sharing ***
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
  if (isOpen && list.value && hasWritePermission(list.value.listAccess?.permissionLevel)) {
    fetchListRecipients()
  }
})

// Function to fetch list recipients
async function fetchListRecipients() {
  if (!list.value?.name) return

  try {
    const request: ListAccessesRequest = {
      parent: trimmedListName.value,
      filter: undefined,
      pageSize: undefined,
      pageToken: undefined
    }

    const response = await listAccessService.ListAccesses(request)

    if (response.accesses) {
      currentAccesses.value = response.accesses.filter(access => {
        // Filter out the current user's own access to avoid showing it in the shares list
        return access.recipient?.user?.name !== authStore.user.name
      }).map(access => ({
        name: access.name || '',
        level: access.level || 'PERMISSION_LEVEL_READ',
        state: access.state || 'ACCESS_STATE_PENDING',
        recipient: access.recipient,
        requester: access.requester || undefined,
        acceptTarget: access.acceptTarget || 'ACCEPT_TARGET_UNSPECIFIED',
      }))
    }
  } catch (error) {
    console.error('Error fetching list recipients:', error)
  }
}

async function updatePermission({ access, newLevel }: { access: Access, newLevel: PermissionLevel }) {
  if (access.level === newLevel) return
  if (!access.name) return
  updatingPermission.value[access.name] = true
  try {
    await listAccessService.UpdateAccess({
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
    console.error('Error updating permission:', error)
  } finally {
    updatingPermission.value[access.name] = false
  }
}

async function unshareList(accessName: string) {
  unsharing.value[accessName] = true
  try {
    const request: DeleteAccessRequest = {
      name: accessName
    }
    
    await listAccessService.DeleteAccess(request)
    await fetchListRecipients()
  } catch (error) {
    console.error('Error removing share:', error)
  } finally {
    unsharing.value[accessName] = false
  }
}

async function shareWithUser({ userName, permission }: { userName: string, permission: PermissionLevel }) {
  if (!userName) return
  if (!list.value?.name) return

  sharing.value = true
  try {
    const access: Access = {
      recipient: {
        user: {
          name: userName,
          username: undefined,
          givenName: undefined,
          familyName: undefined,
        }
      },
      level: permission,
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: trimmedListName.value,
      access
    }

    await listAccessService.CreateAccess(request)
    await fetchListRecipients()
  } catch (error) {
    console.error('Error sharing list:', error)
  } finally {
    sharing.value = false
  }
}

async function shareWithCircle({ circleName, permission }: { circleName: string, permission: PermissionLevel }) {
  if (!circleName) return
  if (!list.value?.name) return

  sharing.value = true
  try {
    const access: Access = {
      recipient: {
        circle: { 
          name: circleName, 
          title: undefined, 
          handle: undefined,
        }
      },
      level: permission,
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: list.value.name,
      access
    }

    await listAccessService.CreateAccess(request)
    await fetchListRecipients()
  } catch (error) {
    console.error('Error sharing list:', error)
  } finally {
    sharing.value = false
  }
}

const requestingAccess = ref(false)

async function handleRequestAccess() {
  if (!list.value?.name) return
  if (!authStore.user?.name) return
  
  requestingAccess.value = true
  try {
    const access: Access = {
       recipient: {
        user: {
          name: authStore.user?.name,
          username: undefined,
          givenName: undefined,
          familyName: undefined,
        }
      },
      level: 'PERMISSION_LEVEL_READ',
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: list.value.name,
      access
    }

    await listAccessService.CreateAccess(request)
    await listStore.loadList(list.value.name)
    alertsStore.addAlert('Access request sent.', 'info')
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    requestingAccess.value = false
  }
}

// *** Favorite List ***
async function toggleFavorite() {
  if (!list.value?.name) return
  try {
    if (list.value.favorited) {
      await listService.UnfavoriteList({ name: list.value.name })
    } else {
      await listService.FavoriteList({ name: list.value.name })
    }
    list.value.favorited = !list.value.favorited
  } catch (error) {
    console.error('Error toggling favorite:', error)
  }
}

// Load list items
async function loadListItems() {
  if (!list.value?.name) return
  
  loadingListItems.value = true
  try {
    const request: ListListItemsRequest = {
      parent: list.value.name,
      pageSize: undefined,
      pageToken: undefined,
      filter: undefined,
    }
    
    const response: ListListItemsResponse = await listItemService.ListListItems(request)
    listItems.value = response.listItems ?? []
  } catch (error) {
    console.error('Error loading list items:', error)
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    loadingListItems.value = false
  }
}

// Create new list item
async function createListItem() {
  if (!list.value?.name || !newItemText.value.trim()) return
  
  try {
    const request: CreateListItemRequest = {
      parent: list.value.name,
      listItem: {
        name: undefined, // Will be set by server
        title: newItemText.value.trim(),
        listSection: undefined,
        points: 0,
        recurrenceRule: undefined,
        createTime: undefined,
        updateTime: undefined,
      }
    }
    
    const newItem = await listItemService.CreateListItem(request)
    listItems.value.push(newItem)
    newItemText.value = ''
    isCreatingNewItem.value = false
  } catch (error) {
    console.error('Error creating list item:', error)
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  }
}

// Create new section
async function createSection() {
  if (!newSectionText.value.trim()) return
  
  try {
    // Add new section to the list
    const newSection = {
      name: undefined, // Will be set by server
      title: newSectionText.value.trim(),
    }
    
    // Update the list with the new section
    if (list.value) {
      list.value.sections = [...(list.value.sections || []), newSection]
    }
    
    await listStore.updateList()
    newSectionText.value = ''
    isCreatingNewSection.value = false
  } catch (error) {
    console.error('Error creating section:', error)
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  }
}

// Handle new item input
function startCreatingItem() {
  isCreatingNewItem.value = true
  newItemText.value = ''
  setTimeout(() => {
    newItemInputRef.value?.focus()
  }, 100)
}

function handleItemInputKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter') {
    event.preventDefault()
    createListItem()
  }
}

function handleItemInputBlur() {
  if (newItemText.value.trim()) {
    createListItem()
  } else {
    isCreatingNewItem.value = false
    newItemText.value = ''
  }
}

// Handle new section input
function startCreatingSection() {
  isCreatingNewSection.value = true
  newSectionText.value = ''
  setTimeout(() => {
    newSectionInputRef.value?.focus()
  }, 100)
}

function handleSectionInputBlur() {
  if (newSectionText.value.trim()) {
    createSection()
  } else {
    isCreatingNewSection.value = false
    newSectionText.value = ''
  }
}

// Load list on mount
onMounted(async () => {
  await listStore.loadList(trimmedListName.value)
  await loadListItems()
})
</script>

<style scoped>
.list-item {
  min-height: auto;
}

.cursor-pointer {
  cursor: pointer;
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

.favorite-heart {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10000000;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.6));
  border-radius: 50%;
  padding: 4px;
  transition: all 0.2s ease-in-out;
  color: red;
}
</style>