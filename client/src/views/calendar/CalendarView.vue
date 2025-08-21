<template>
  <template v-if="calendar">
    <ListTabsPage
      :tabs="tabs"
      ref="tabsPage"
    >
      <template #general>
        <v-container max-width="600" class="pa-1">
          <v-row>
            <v-col class="pt-5">
              <div class="text-h4">
                {{ calendar.title }}
              </div>
              <div class="bio-section pa-2">
                <div class="text-subtitle-1 font-weight-bold mb-1">Description</div>
                <div class="text-body-1" style="white-space: pre-line;">
                  {{ calendar.description || 'No description set.' }}
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
          <v-row v-if="calendar.calendarAccess?.state === 'ACCESS_STATE_PENDING'">
            <v-col cols="12">
              <v-btn
                color="success"
                class="mb-2"
                block
                :loading="acceptingCalendar"
                @click="acceptCalendar(calendar.calendarAccess?.name)"
              >
                Accept Calendar
              </v-btn>
              <v-btn
                color="error"
                class="mb-4"
                block
                :loading="decliningCalendar"
                @click="declineCalendar"
              >
                Decline
              </v-btn>
            </v-col>
          </v-row>
        </v-container>
        
        <!-- Floating action buttons container -->
        <div class="fab-container">
          <v-btn density="compact" v-if="hasWritePermission(calendar.calendarAccess?.permissionLevel)" color="primary" :to="'/'+calendar.name+'/edit'">
            <v-icon>mdi-pencil</v-icon>Edit
          </v-btn>
          <v-btn density="compact" v-if="!hasReadPermission(calendar.calendarAccess?.permissionLevel) && !calendar.calendarAccess" color="primary" @click="handleRequestAccess" :loading="requestingAccess">
            <v-icon>mdi-account-plus</v-icon>Request Access
          </v-btn>
          <v-btn v-if="hasWritePermission(calendar.calendarAccess?.permissionLevel) && calendar.visibility !== 'VISIBILITY_LEVEL_HIDDEN'" color="primary" density="compact" @click="showShareDialog = true">
            <v-icon>mdi-share-variant</v-icon>
            Manage Access
          </v-btn>
          <v-btn density="compact" v-if="calendar.calendarAccess?.state === 'ACCESS_STATE_PENDING'" color="warning" @click="showCancelRequestDialog = true">
            <v-icon>mdi-close</v-icon>Cancel Request
          </v-btn>
          <v-btn density="compact" v-if="!hasAdminPermission(calendar.calendarAccess?.permissionLevel) && calendar.calendarAccess?.state === 'ACCESS_STATE_ACCEPTED'" color="warning" @click="showRemoveAccessDialog = true">
            <v-icon>mdi-link-variant-off</v-icon>Remove Access
          </v-btn>
          <v-btn density="compact" v-if="hasAdminPermission(calendar.calendarAccess?.permissionLevel)" color="error" @click="showDeleteDialog = true">
            <v-icon>mdi-delete</v-icon>Delete
          </v-btn>
        </div>
        <v-icon 
          size="24" 
          :color="calendar.favorited ? 'red' : 'black'"
          class="favorite-heart"
          @click="toggleFavorite"
          style="cursor: pointer;"
        >
        {{ calendar.favorited ? 'mdi-heart' : 'mdi-heart-outline' }}
        </v-icon>
      </template>
      
      <template #events="{ items, loading }">
        <ScheduleCal 
          :events="(items as Event[])" 
          :calendars="[calendar]" 
          :loading="(loading as boolean)" 
          :show-create-button="true"
          @created="onEventCreated" 
          @updated="onEventUpdated" 
          @deleted="onEventDeleted"
        />
      </template>
      

    </ListTabsPage>

    <!-- Share Dialog -->
    <ShareDialog
      v-model="showShareDialog"
      title="Share Calendar"
      :allowCircleShare="false"
      :currentAccesses="currentAccesses"
      :sharing="sharing"
      :accessPermissionLoading="updatingPermission"
      :userPermissionLevel="calendar.calendarAccess?.permissionLevel"
      :allowPermissionOptions="allowPermissionOptions"
      @share-user="shareWithUser"
      @share-circle="shareWithCircle"
      @remove-access="unshareCalendar"
      @permission-change="updatePermission"
      @approve-access="acceptCalendarFromShareDialog"
    />

    <!-- Remove Access Dialog -->
    <v-dialog v-model="showRemoveAccessDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">
          Remove Access
        </v-card-title>
        <v-card-text>
          Are you sure you want to remove your access to this calendar? You may no longer be able to view it.
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

    <!-- Cancel Request Dialog -->
    <v-dialog v-model="showCancelRequestDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">
          Cancel Access Request
        </v-card-title>
        <v-card-text>
          Are you sure you want to cancel your access request to this calendar?
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" variant="text" @click="showCancelRequestDialog = false">
            Cancel
          </v-btn>
          <v-btn color="warning" @click="handleCancelRequest" :loading="requestingAccess">
            Cancel Request
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">
          Delete Calendar
        </v-card-title>
        <v-card-text>
          Are you sure you want to delete this calendar? This action will also delete the calendar for any users that
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
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCalendarsStore } from '@/stores/calendar'
import { storeToRefs } from 'pinia'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import ShareDialog from '@/components/common/ShareDialog.vue'
import type { apitypes_VisibilityLevel, Access, apitypes_PermissionLevel, CreateAccessRequest, ListAccessesRequest, DeleteAccessRequest } from '@/genapi/api/calendars/calendar/v1alpha1'
import { calendarAccessService } from '@/api/api'
import { useAlertStore } from '@/stores/alerts'
import { useAuthStore } from '@/stores/auth'
import ScheduleCal from '@/views/calendar/event/ScheduleCal.vue'
import type { Event } from '@/genapi/api/calendars/calendar/v1alpha1'
import { calendarService } from '@/api/api'

const route = useRoute()
const router = useRouter()
const calendarsStore = useCalendarsStore()

const alertsStore = useAlertStore()
const authStore = useAuthStore()
const { calendar } = storeToRefs(calendarsStore)

const acceptingCalendar = ref(false)
const decliningCalendar = ref(false)
const requestingAccess = ref(false)
const removingAccess = ref(false)
const deleting = ref(false)
const showRemoveAccessDialog = ref(false)
const showDeleteDialog = ref(false)
const showCancelRequestDialog = ref(false)
const showShareDialog = ref(false)

// *** Calendar Sharing ***
const allowPermissionOptions: apitypes_PermissionLevel[] = [
  'PERMISSION_LEVEL_READ',
  'PERMISSION_LEVEL_WRITE',
  'PERMISSION_LEVEL_ADMIN'
]
const currentAccesses = ref<Access[]>([])
const sharing = ref(false)
const updatingPermission = ref<Record<string, boolean>>({})

const tabs = [
  {
    label: 'General',
    value: 'general',
    icon: 'mdi-information'
  },
  {
    label: 'Events',
    value: 'events',
    icon: 'mdi-calendar',
    loader: async () => {
      return await calendarsStore.loadEvents(calendar.value?.name ?? '')
    }
  }
]

const visibilityOptions = [
  {
    label: 'Public',
    value: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
    icon: 'mdi-earth',
    color: 'success',
    description: 'Everyone can see this calendar'
  },
  {
    label: 'Restricted',
    value: 'VISIBILITY_LEVEL_RESTRICTED' as apitypes_VisibilityLevel,
    icon: 'mdi-account-group',
    color: 'warning',
    description: 'Shared users and their connections can see this'
  },
  {
    label: 'Private',
    value: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
    icon: 'mdi-lock',
    color: 'info',
    description: 'Only specifically shared users can see this'
  },
  {
    label: 'Hidden',
    value: 'VISIBILITY_LEVEL_HIDDEN' as apitypes_VisibilityLevel,
    icon: 'mdi-eye-off',
    color: 'secondary',
    description: 'Only you can see this calendar'
  }
]

// Computed properties for the selected visibility
const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === calendar.value?.visibility)
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

function hasAdminPermission(permissionLevel?: string): boolean {
  return permissionLevel === 'PERMISSION_LEVEL_ADMIN'
}

function hasWritePermission(permissionLevel?: string): boolean {
  return permissionLevel === 'PERMISSION_LEVEL_ADMIN' || permissionLevel === 'PERMISSION_LEVEL_WRITE'
}

function hasReadPermission(permissionLevel?: string): boolean {
  return permissionLevel === 'PERMISSION_LEVEL_ADMIN' || permissionLevel === 'PERMISSION_LEVEL_WRITE' || permissionLevel === 'PERMISSION_LEVEL_READ'
}

async function acceptCalendar(accessName?: string) {
  if (!accessName) return
  
  acceptingCalendar.value = true
  try {
    await calendarsStore.acceptCalendar(accessName)
    // Refresh the calendar
    await loadCalendar()
  } catch (error) {
    console.error('Failed to accept calendar:', error)
  } finally {
    acceptingCalendar.value = false
  }
}

async function declineCalendar() {
  if (!calendar.value?.name) return
  
  decliningCalendar.value = true
  try {
    await calendarsStore.deleteCalendarAccess(calendar.value.name)
    router.push('/calendars')
  } catch (error) {
    console.error('Failed to decline calendar:', error)
  } finally {
    decliningCalendar.value = false
  }
}

async function handleRequestAccess() {
  if (!calendar.value?.name || !authStore.user?.name) return
  
  requestingAccess.value = true
  try {
    const access: Access = {
      recipient: {
        user: {
          name: authStore.user.name,
          username: undefined,
          givenName: undefined,
          familyName: undefined,
        },
      },
      level: 'PERMISSION_LEVEL_READ',
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: calendar.value.name,
      access
    }

    await calendarAccessService.CreateAccess(request)
    alertsStore.addAlert('Access request sent successfully', 'success')
    // Refresh the calendar to show the new access state
    await loadCalendar()
  } catch (error) {
    console.error('Failed to request access:', error)
    alertsStore.addAlert('Failed to request access', 'error')
  } finally {
    requestingAccess.value = false
  }
}

async function handleRemoveAccess() {
  if (!calendar.value?.name) return
  
  removingAccess.value = true
  try {
    await calendarsStore.deleteCalendarAccess(calendar.value.name)
    router.push('/calendars')
  } catch (error) {
    console.error('Failed to remove access:', error)
  } finally {
    removingAccess.value = false
    showRemoveAccessDialog.value = false
  }
}

async function handleDelete() {
  if (!calendar.value?.name) return
  
  deleting.value = true
  try {
    await calendarsStore.deleteCalendar(calendar.value.name)
    router.push('/calendars')
  } catch (error) {
    console.error('Failed to delete calendar:', error)
  } finally {
    deleting.value = false
    showDeleteDialog.value = false
  }
}

async function handleCancelRequest() {
  if (!calendar.value?.calendarAccess?.name) return
  
  requestingAccess.value = true
  try {
    await calendarsStore.deleteCalendarAccess(calendar.value.calendarAccess.name)
    alertsStore.addAlert('Access request cancelled.', 'info')
    router.push('/calendars')
  } catch (error) {
    console.error('Failed to cancel access request:', error)
    alertsStore.addAlert('Failed to cancel access request', 'error')
  } finally {
    requestingAccess.value = false
    showCancelRequestDialog.value = false
  }
}

// *** Calendar Sharing Functions ***

// Function to fetch calendar recipients
async function fetchCalendarRecipients() {
  if (!calendar.value?.name) return

  try {
    const request: ListAccessesRequest = {
      parent: calendar.value.name,
      filter: undefined,
      pageSize: undefined,
      pageToken: undefined
    }

    const response = await calendarAccessService.ListAccesses(request)

    if (response.accesses) {
      currentAccesses.value = response.accesses.filter(access => {
        // Filter out the current user's own access to avoid showing it in the shares list
        return access.recipient?.user?.name !== authStore.user?.name
      })
    }
  } catch (error) {
    console.error('Error fetching calendar recipients:', error)
  }
}

async function shareWithUser({ userName, permission }: { userName: string, permission: apitypes_PermissionLevel }) {
  if (!userName) return
  if (!calendar.value?.name) return
  
  sharing.value = true
  try {
    const access: Access = {
      recipient: {
        user: {
          name: userName,
          username: undefined,
          givenName: undefined,
          familyName: undefined,
        },
      },
      level: permission,
      name: undefined, // Will be set by the server
      requester: undefined, // Will be set by the server
      state: undefined, // Will be set by the server
      acceptTarget: undefined, // Will be set by the server
    }

    const request: CreateAccessRequest = {
      parent: calendar.value.name,
      access
    }

    await calendarAccessService.CreateAccess(request)
    await fetchCalendarRecipients()
  } catch (error) {
    console.error('Error sharing calendar:', error)
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    sharing.value = false
  }
}

async function shareWithCircle() {
  // Calendar access API doesn't support circles, only users
  alertsStore.addAlert('Sharing with circles is not supported for calendars', 'warning')
}

async function unshareCalendar(accessName: string) {
  try {
    const request: DeleteAccessRequest = {
      name: accessName
    }
    
    await calendarAccessService.DeleteAccess(request)
    await fetchCalendarRecipients()
  } catch (error) {
    console.error('Error removing share:', error)
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  }
}

async function updatePermission({ access, newLevel }: { access: Access, newLevel: apitypes_PermissionLevel }) {
  if (access.level === newLevel) return
  if (!access.name) return
  
  updatingPermission.value[access.name] = true
  try {
    await calendarAccessService.UpdateAccess({
      access: {
        name: access.name,
        level: newLevel,
        state: undefined,
        recipient: undefined,
        requester: undefined,
        acceptTarget: undefined,
      },
      updateMask: 'permissionLevel',
    })
    // Update local state
    access.level = newLevel
  } catch (error) {
    console.error('Error updating permission:', error)
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  } finally {
    updatingPermission.value[access.name] = false
  }
}

async function acceptCalendarFromShareDialog(accessName: string) {
  if (!accessName) return
  
  try {
    await calendarsStore.acceptCalendar(accessName)
    await fetchCalendarRecipients()
  } catch (error) {
    console.error('Error accepting calendar access:', error)
    alertsStore.addAlert(error instanceof Error ? error.message : String(error), 'error')
  }
}

async function loadCalendar() {
  const calendarName = route.path.substring(1)
  await calendarsStore.loadCalendar(calendarName)
}

onMounted(async () => {
  await loadCalendar()
})

watch(
  () => route.path,
  async (newPath) => {
    if (newPath) {
      await loadCalendar()
    }
  }
)

// Watch for share dialog opening to load current accesses
watch(showShareDialog, (isOpen) => {
  if (isOpen && calendar.value && hasWritePermission(calendar.value.calendarAccess?.permissionLevel)) {
    fetchCalendarRecipients()
  }
})

async function onEventCreated() {
  if (!calendar.value?.name) return
  await calendarsStore.loadEvents(calendar.value.name)
}

async function onEventUpdated() {
  if (!calendar.value?.name) return
  await calendarsStore.loadEvents(calendar.value.name)
}

async function onEventDeleted() {
  if (!calendar.value?.name) return
  await calendarsStore.loadEvents(calendar.value.name)
}

// *** Favorite User ***

async function toggleFavorite() {
  if (!calendar.value?.name) return
  try {
    if (calendar.value.favorited) {
      await calendarService.UnfavoriteCalendar({ name: calendar.value.name })
    } else {
      await calendarService.FavoriteCalendar({ name: calendar.value.name })
    }
    // Update the local state
    calendar.value.favorited = !calendar.value.favorited
  } catch (error) {
    console.error('Error toggling favorite:', error)
    alertsStore.addAlert('Failed to update favorite status.', 'error')
  }
}

</script>

<style scoped>
.bio-section {
  background-color: rgb(var(--v-theme-surface));
  border-radius: 4px;
  margin-top: 16px;
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
  /* background-color: rgba(255, 255, 255, 0.9); */
  border-radius: 50%;
  padding: 4px;
  transition: all 0.2s ease-in-out;
  color: red;
}
</style>
