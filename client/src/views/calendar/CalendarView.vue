<template>
  <v-container v-if="calendar" class="pb-16">
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
      </template>
      
      <template #events>
        <div class="text-center py-8">
          <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-calendar</v-icon>
          <h3 class="text-grey-lighten-1 mb-2">Events</h3>
          <p class="text-grey-lighten-1">Calendar events will be displayed here.</p>
        </div>
      </template>
      
      <template #members>
        <div class="text-center py-8">
          <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-account-group</v-icon>
          <h3 class="text-grey-lighten-1 mb-2">Members</h3>
          <p class="text-grey-lighten-1">Calendar members will be displayed here.</p>
        </div>
        <v-btn
          v-if="hasWritePermission(calendar.calendarAccess?.permissionLevel)"
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
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCalendarsStore } from '@/stores/calendar'
import { storeToRefs } from 'pinia'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import type { apitypes_VisibilityLevel } from '@/genapi/api/calendars/calendar/v1alpha1'

const route = useRoute()
const router = useRouter()
const calendarsStore = useCalendarsStore()
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

const tabs = [
  {
    label: 'General',
    value: 'general',
    icon: 'mdi-information'
  },
  {
    label: 'Events',
    value: 'events',
    icon: 'mdi-calendar'
  },
  {
    label: 'Members',
    value: 'members',
    icon: 'mdi-account-group'
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
  // TODO: Implement request access functionality
  console.log('Request access not implemented yet')
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
    // TODO: Implement delete calendar functionality
    console.log('Delete calendar not implemented yet')
    router.push('/calendars')
  } catch (error) {
    console.error('Failed to delete calendar:', error)
  } finally {
    deleting.value = false
    showDeleteDialog.value = false
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
</style>
