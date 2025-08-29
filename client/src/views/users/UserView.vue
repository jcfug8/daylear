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
          <v-btn density="compact" v-if="hasWriteOnlyPermission(user.access?.permissionLevel) && user.access?.state === 'ACCESS_STATE_ACCEPTED'" color="warning" @click="showRemoveAccessDialog = true">
            <v-icon>mdi-account-remove</v-icon>Remove Access
          </v-btn>
        </div>
        <v-icon 
          size="24" 
          :color="user.favorited ? 'red' : 'black'"
          class="favorite-heart"
          @click="toggleFavorite"
          style="cursor: pointer;"
        >
        {{ user.favorited ? 'mdi-heart' : 'mdi-heart-outline' }}
        </v-icon>
      </template>
      <template #accessKeys="{ items, loading }">
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-card>
                <v-card-title class="d-flex align-center justify-space-between">
                  <span>Access Keys</span>
                  <v-btn
                    color="primary"
                    @click="showCreateAccessKeyDialog = true"
                    prepend-icon="mdi-plus"
                  >
                    Create Access Key
                  </v-btn>
                </v-card-title>
                <v-card-text>
                  <div v-if="loading" class="text-center pa-4">
                    <v-progress-circular indeterminate color="primary"></v-progress-circular>
                  </div>
                  <div v-else-if="!items || items.length === 0" class="text-center pa-4 text-grey">
                    No access keys found.
                  </div>
                  <v-list v-else>
                    <v-list-item
                      v-for="accessKey in (items as AccessKey[])"
                      :key="accessKey.name"
                      class="mb-2"
                    >
                      <template v-slot:prepend>
                        <v-icon icon="mdi-key"></v-icon>
                      </template>
                      <v-list-item-title>{{ accessKey.title }}</v-list-item-title>
                      <v-list-item-subtitle v-if="accessKey.description">
                        {{ accessKey.description }}
                      </v-list-item-subtitle>
                      <template v-slot:append>
                        <v-btn
                          icon="mdi-pencil"
                          variant="text"
                          size="small"
                          @click="editAccessKey(accessKey)"
                        ></v-btn>
                        <v-btn
                          icon="mdi-delete"
                          variant="text"
                          size="small"
                          color="error"
                          @click="confirmDeleteAccessKey(accessKey)"
                        ></v-btn>
                      </template>
                    </v-list-item>
                  </v-list>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </v-container>
      </template>
      <template #recipes="{ items, loading }">
        <RecipeGrid :recipes="(items as Recipe[])" :loading="(loading as boolean)" />
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
        <UserGrid :users="(items as User[])" :loading="(loading as boolean)" empty-text="No friends found." />
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
        <CircleGrid :circles="(items as Circle[])" :loading="(loading as boolean)" empty-text="No circles found." />
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
      <template #calendars="{ items, loading }">
        <CalendarGrid v-if="viewMode === 'grid'" :calendars="(items as Calendar[])" :loading="(loading as boolean)" />
        <template v-else>
          <ScheduleCal 
            v-if="!loading" 
            :events="events" 
            :calendars="(items as Calendar[])" 
            :show-create-button="hasAdminPermission(user.access?.permissionLevel)" 
            @created="onEventCreated" 
            @updated="onEventUpdated" 
            @deleted="onEventDeleted"
          />
        </template>
        <!-- View mode toggle FAB -->
        <v-btn
          color="primary"
          density="compact"
          class="text-none"
          style="position: fixed; bottom: 16px; left: 16px; z-index: 10;"
          @click="toggleViewMode"
        >
          <v-icon class="mr-1">{{ viewMode === 'grid' ? 'mdi-calendar-month' : 'mdi-view-grid' }}</v-icon>
          <span>{{ viewMode === 'grid' ? 'Schedule' : 'Grid' }}</span>
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

  <!-- Create Access Key Dialog -->
  <v-dialog v-model="showCreateAccessKeyDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Create Access Key
      </v-card-title>
      <v-card-text>
        <v-form @submit.prevent="createAccessKey">
          <v-text-field
            v-model="accessKeyTitle"
            label="Title"
            required
            :rules="[v => !!v || 'Title is required']"
          ></v-text-field>
          <v-textarea
            v-model="accessKeyDescription"
            label="Description (optional)"
            rows="3"
          ></v-textarea>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showCreateAccessKeyDialog = false">
          Cancel
        </v-btn>
        <v-btn 
          color="primary" 
          @click="createAccessKey" 
          :loading="creatingAccessKey"
          :disabled="!accessKeyTitle.trim()"
        >
          Create
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Edit Access Key Dialog -->
  <v-dialog v-model="showEditAccessKeyDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        Edit Access Key
      </v-card-title>
      <v-card-text>
        <v-form @submit.prevent="updateAccessKey">
          <v-text-field
            v-model="accessKeyTitle"
            label="Title"
            required
            :rules="[v => !!v || 'Title is required']"
          ></v-text-field>
          <v-textarea
            v-model="accessKeyDescription"
            label="Description (optional)"
            rows="3"
          ></v-textarea>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showEditAccessKeyDialog = false">
          Cancel
        </v-btn>
        <v-btn 
          color="primary" 
          @click="updateAccessKey"
          :disabled="!accessKeyTitle.trim()"
        >
          Update
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Delete Access Key Confirmation Dialog -->
  <v-dialog v-model="showDeleteAccessKeyDialog" max-width="500">
    <v-card>
      <v-card-title class="text-h5 d-flex align-center">
        <v-icon color="error" class="mr-2">mdi-alert-circle</v-icon>
        Delete Access Key
      </v-card-title>
      <v-card-text>
        <div class="mb-4">
          <p class="text-body-1 mb-2">
            Are you sure you want to delete the access key <strong>"{{ deletingAccessKey?.title }}"</strong>?
          </p>
          <div class="text-caption text-grey">
            <v-icon size="small" class="mr-1">mdi-information</v-icon>
            This action cannot be undone. Any systems using this access key will no longer be able to authenticate.
          </div>
        </div>
        
        <v-card variant="outlined" class="pa-3" color="error">
          <div class="text-subtitle-2 mb-2">Access Key Details:</div>
          <div class="mb-2">
            <strong>Title:</strong> {{ deletingAccessKey?.title }}
          </div>
          <div v-if="deletingAccessKey?.description" class="mb-2">
            <strong>Description:</strong> {{ deletingAccessKey?.description }}
          </div>
        </v-card>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="showDeleteAccessKeyDialog = false">
          Cancel
        </v-btn>
        <v-btn 
          color="error" 
          @click="deleteAccessKey"
          :loading="isDeletingAccessKey"
        >
          Delete Access Key
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Created Access Key Dialog -->
  <v-dialog v-model="showCreatedKeyDialog" max-width="600" persistent>
    <v-card>
      <v-card-title class="text-h5 d-flex align-center">
        <v-icon color="success" class="mr-2">mdi-check-circle</v-icon>
        Access Key Created Successfully!
      </v-card-title>
      <v-card-text>
        <div class="mb-4">
          <p class="text-body-1 mb-2">
            Your access key has been created. <strong>Copy this key now</strong> - you won't be able to see it again!
          </p>
          <div class="text-caption text-grey">
            <v-icon size="small" class="mr-1">mdi-information</v-icon>
            This is the only time you'll see the unencrypted access key. Store it securely.
          </div>
        </div>
        
        <v-card variant="outlined" class="pa-3 mb-3">
          <div class="text-subtitle-2 mb-2">Access Key Details:</div>
          <div class="mb-2">
            <strong>Title:</strong> {{ newlyCreatedAccessKey?.title }}
          </div>
          <div class="mb-3" v-if="newlyCreatedAccessKey?.description">
            <strong>Description:</strong> {{ newlyCreatedAccessKey?.description }}
          </div>
          
          <div class="mb-2">
            <strong>Access Key:</strong>
          </div>
          <v-text-field
            :model-value="newlyCreatedAccessKey?.unencryptedAccessKey || ''"
            readonly
            variant="outlined"
            density="compact"
            class="mb-2"
            append-inner-icon="mdi-content-copy"
            @click:append-inner="copyAccessKey"
          ></v-text-field>
          
          <div class="text-caption text-grey">
            <v-icon size="small" class="mr-1">mdi-alert</v-icon>
            Keep this key secure and don't share it with others.
          </div>
        </v-card>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn 
          color="primary" 
          @click="closeCreatedKeyDialog"
        >
          Got it!
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
import { onMounted, ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import { useAuthStore } from '@/stores/auth'
import { hasWritePermission, hasAdminPermission, hasWriteOnlyPermission } from '@/utils/permissions'
import { userAccessService } from '@/api/api'
import type { DeleteAccessRequest, ListAccessesRequest } from '@/genapi/api/meals/recipe/v1alpha1'
import type { CreateAccessRequest, Access, User, AccessKey } from '@/genapi/api/users/user/v1alpha1'
import { useAlertStore } from '@/stores/alerts'
import UserGrid from '@/components/UserGrid.vue'
import CircleGrid from '@/components/CircleGrid.vue'
import { useCirclesStore } from '@/stores/circles'
import ShareDialog from '@/components/common/ShareDialog.vue'
import type { PermissionLevel } from '@/genapi/api/types'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'
import { useCalendarsStore } from '@/stores/calendar'
import CalendarGrid from '@/components/CalendarGrid.vue'
import ScheduleCal from '@/views/calendar/event/ScheduleCal.vue'
import type { Calendar, Event } from '@/genapi/api/calendars/calendar/v1alpha1'
import { userService } from '@/api/api'

const usersStore = useUsersStore()
const alertsStore = useAlertStore()
const recipesStore = useRecipesStore()
const circlesStore = useCirclesStore()
const calendarsStore = useCalendarsStore()
const { currentUser: user } = storeToRefs(usersStore)
const route = useRoute()
const tabsPage = ref()
// const speedDialOpen = ref(false)

// Calendars tab state
const viewMode = ref<'grid' | 'schedule'>('grid')
const events = ref<Event[]>([])

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

async function onEventCreated() {
  if (viewMode.value === 'schedule' && tabsPage.value?.activeTab?.value === 'calendars') {
    await loadEventsForCalendars()
  }
}

async function onEventUpdated() {
  if (viewMode.value === 'schedule' && tabsPage.value?.activeTab?.value === 'calendars') {
    await loadEventsForCalendars()
  }
}

async function onEventDeleted() {
  if (viewMode.value === 'schedule' && tabsPage.value?.activeTab?.value === 'calendars') {
    await loadEventsForCalendars()
  }
}

const tabs = computed(() => {
  const baseTabs: any[] = [
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

  // Only show Access Keys tab if user has admin permissions
  if (hasAdminPermission(user.value?.access?.permissionLevel)) {
    baseTabs.splice(1, 0, {
      label: 'Access Keys',
      value: 'accessKeys',
      loader: async () => {
        if (!user.value?.name) return []
        await usersStore.loadAccessKeys(user.value.name)
        return [...usersStore.accessKeys]
      },
    })
  }

  // Add remaining tabs
  baseTabs.push(
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
  {
    label: 'Calendars',
    value: 'calendars',
    loader: async () => {
      if (!user.value?.name) return []
      await calendarsStore.loadMyCalendars(user.value.name)
      return [...calendarsStore.myCalendars]
    },
  }
  )

  return baseTabs
})

const authStore = useAuthStore()

// Computed property to show access request section
const showAccessRequest = computed(() => {
  return user.value?.access?.state === 'ACCESS_STATE_PENDING' && 
         user.value?.access?.requester !== authStore.user?.name
})

function toggleViewMode() {
  viewMode.value = viewMode.value === 'grid' ? 'schedule' : 'grid'
  if (viewMode.value === 'schedule' && tabsPage.value?.activeTab?.value === 'calendars') {
    void loadEventsForCalendars()
  }
}

async function loadEventsForCalendars() {
  events.value = []
  for (const calendar of calendarsStore.myCalendars) {
    if (!calendar?.name) continue
    const es = await calendarsStore.loadEvents(calendar.name)
    events.value.push(...es)
  }
}

// *** Remove Access ***
const showRemoveAccessDialog = ref(false)
const removingAccess = ref(false)

// *** Cancel Request ***
const showCancelRequestDialog = ref(false)
const cancelingRequest = ref(false)

// *** Accept/Decline Request ***
const acceptingRequest = ref(false)
const decliningRequest = ref(false)

// *** Access Keys ***
const showCreateAccessKeyDialog = ref(false)
const showEditAccessKeyDialog = ref(false)
const editingAccessKey = ref<AccessKey | null>(null)
const creatingAccessKey = ref(false)
const accessKeyTitle = ref('')
const accessKeyDescription = ref('')
const newlyCreatedAccessKey = ref<AccessKey | null>(null)
const showCreatedKeyDialog = ref(false)
const showDeleteAccessKeyDialog = ref(false)
const deletingAccessKey = ref<AccessKey | null>(null)
const isDeletingAccessKey = ref(false)

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

// *** Access Key Methods ***
async function createAccessKey() {
  if (!user.value?.name || !accessKeyTitle.value.trim()) return
  
  creatingAccessKey.value = true
  try {
    const newAccessKey = await usersStore.createAccessKey(user.value.name, accessKeyTitle.value.trim(), accessKeyDescription.value.trim() || undefined)
    
    // Store the newly created key to show the unencrypted value
    newlyCreatedAccessKey.value = newAccessKey
    
    // Clear form and close create dialog
    accessKeyTitle.value = ''
    accessKeyDescription.value = ''
    showCreateAccessKeyDialog.value = false
    
    // Show the dialog with the generated key
    showCreatedKeyDialog.value = true
    
    // Refresh the access keys tab
    await tabsPage.value?.reloadTab('accessKeys')
    
    alertsStore.addAlert('Access key created successfully!', 'success')
  } catch (error) {
    const msg = `Error creating access key: ${error instanceof Error ? error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  } finally {
    creatingAccessKey.value = false
  }
}

function editAccessKey(accessKey: AccessKey) {
  editingAccessKey.value = accessKey
  accessKeyTitle.value = accessKey.title || ''
  accessKeyDescription.value = accessKey.description || ''
  showEditAccessKeyDialog.value = true
}

async function updateAccessKey() {
  if (!editingAccessKey.value || !accessKeyTitle.value.trim()) return
  
  try {
    const updateMask: string[] = []
    if (editingAccessKey.value.title !== accessKeyTitle.value.trim()) {
      updateMask.push('title')
    }
    if (editingAccessKey.value.description !== accessKeyDescription.value.trim()) {
      updateMask.push('description')
    }
    
    if (updateMask.length > 0) {
      await usersStore.updateAccessKey({
        ...editingAccessKey.value,
        title: accessKeyTitle.value.trim(),
        description: accessKeyDescription.value.trim(),
      }, updateMask)
      
      alertsStore.addAlert('Access key updated successfully!', 'success')
      showEditAccessKeyDialog.value = false
      editingAccessKey.value = null
      accessKeyTitle.value = ''
      accessKeyDescription.value = ''
      // Refresh the access keys tab
      await tabsPage.value?.reloadTab('accessKeys')
    }
  } catch (error) {
    const msg = `Error updating access key: ${error instanceof Error ? error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  }
}

function confirmDeleteAccessKey(accessKey: AccessKey) {
  deletingAccessKey.value = accessKey
  showDeleteAccessKeyDialog.value = true
}

async function deleteAccessKey() {
  if (!deletingAccessKey.value?.name) return
  
  isDeletingAccessKey.value = true
  try {
    await usersStore.deleteAccessKey(deletingAccessKey.value.name)
    alertsStore.addAlert('Access key deleted successfully!', 'success')
    
    // Close dialog and clear state
    showDeleteAccessKeyDialog.value = false
    deletingAccessKey.value = null
    
    // Refresh the access keys tab
    await tabsPage.value?.reloadTab('accessKeys')
  } catch (error) {
    const msg = `Error deleting access key: ${error instanceof Error ? error.message : String(error)}`
    console.error(msg)
    alertsStore.addAlert(msg, 'error')
  } finally {
    isDeletingAccessKey.value = false
  }
}

function copyAccessKey() {
  if (newlyCreatedAccessKey.value?.unencryptedAccessKey) {
    navigator.clipboard.writeText(newlyCreatedAccessKey.value.unencryptedAccessKey)
    alertsStore.addAlert('Access key copied to clipboard!', 'success')
  }
}

function closeCreatedKeyDialog() {
  showCreatedKeyDialog.value = false
  newlyCreatedAccessKey.value = null
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
      acceptTarget: undefined,
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
        acceptTarget: access.acceptTarget || 'ACCEPT_TARGET_UNSPECIFIED',
      }))
    }
  } catch (error) {
    alertsStore.addAlert(error instanceof Error ? error.message : String(error),'error')
  }
}

async function updatePermission({ access, newLevel }: { access: Access, newLevel: PermissionLevel }) {
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
      acceptTarget: undefined, // Will be set by the server
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

// Calendars tab helpers
watch(
  () => calendarsStore.myCalendars,
  async () => {
    if (viewMode.value === 'schedule' && tabsPage.value?.activeTab?.value === 'calendars') {
      await loadEventsForCalendars()
    }
  },
  { deep: true }
)

watch(
  () => tabsPage.value?.activeTab?.value,
  async (tab) => {
    if (tab === 'calendars' && viewMode.value === 'schedule') {
      await loadEventsForCalendars()
    }
  }
)

// *** Favorite User ***

async function toggleFavorite() {
  if (!user.value?.name) return
  try {
    if (user.value.favorited) {
      await userService.UnfavoriteUser({ name: user.value.name })
    } else {
      await userService.FavoriteUser({ name: user.value.name })
    }
    // Update the local state
    user.value.favorited = !user.value.favorited
  } catch (error) {
    console.error('Error toggling favorite:', error)
    alertsStore.addAlert('Failed to update favorite status.', 'error')
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

.favorite-heart {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10000000;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.6));
  /* background-color: rgba(255, 255, 255, 0.9); */
  border-radius: 50%;
  padding: 4px;
  transition: all 0.2s ease-in-out;
  color: red;
}
</style>
