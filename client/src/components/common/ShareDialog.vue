<template>
  <v-dialog :model-value="modelValue" @update:model-value="emitClose" max-width="500">
    <v-card style="max-height: 80vh; display: flex; flex-direction: column;">
      <v-card-title class="text-h5">
        {{ title }}
      </v-card-title>
      <v-card-text style="flex: 1; overflow: hidden; display: flex; flex-direction: column;">
        <!-- Fixed input section -->
        <div v-if="!disableCreateShare" style="flex-shrink: 0;">
          <v-tabs v-if="allowCircleShare" v-model="accessTab" class="mb-4">
            <v-tab value="users">Share with User</v-tab>
            <v-tab value="circles">Share with Circle</v-tab>
          </v-tabs>
          <v-window v-if="allowCircleShare" v-model="accessTab">
            <v-window-item value="users">
              <v-autocomplete
                density="compact"
                v-model="selectedUser"
                :items="userSuggestions"
                :loading="isLoadingUsers"
                :search-input.sync="usernameInput"
                label="Search Users"
                item-title="displayName"
                item-value="name"
                :rules="[validateUsername]"
                clearable
                :input-attrs="{ autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' }"
                :prepend-inner-icon="getUsernameIcon"
                :color="getUsernameColor"
                @update:search="handleUserSearch"
                @update:model-value="handleUserSelection"
                no-data-text="No friends found"
              >
                <template #item="{ props, item }">
                  <v-list-item v-bind="props">
                    <v-list-item-subtitle>
                      @{{ item.raw.username }}
                    </v-list-item-subtitle>
                  </v-list-item>
                </template>
              </v-autocomplete>
              <v-select
                density="compact"
                v-model="selectedPermission"
                :items="permissionOptions"
                label="Permission Level"
                class="mt-2"
              ></v-select>
              <v-btn block color="primary" @click="emitShareUser" :loading="sharing" :disabled="!isValidUsername" class="mt-2">
                Share with User
              </v-btn>
            </v-window-item>
            <v-window-item value="circles">
              <v-autocomplete
                density="compact"
                v-model="selectedCircle"
                :items="circleSuggestions"
                :loading="isLoadingCircles"
                :search-input.sync="circleInput"
                label="Search Circles"
                item-title="displayName"
                item-value="handle"
                :rules="[validateCircle]"
                clearable
                :input-attrs="{ autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' }"
                :prepend-inner-icon="getCircleIcon"
                :color="getCircleColor"
                @update:search="handleCircleSearch"
                @update:model-value="handleCircleSelection"
                no-data-text="No circles found"
              >
                <template #item="{ props, item }">
                  <v-list-item v-bind="props">
                    <v-list-item-subtitle>
                      @{{ item.raw.handle }}
                    </v-list-item-subtitle>
                  </v-list-item>
                </template>
              </v-autocomplete>
              <v-select
                density="compact"
                v-model="selectedPermission"
                :items="permissionOptions"
                label="Permission Level"
                class="mt-2"
              ></v-select>
              <v-btn block color="primary" @click="emitShareCircle" :loading="sharing" :disabled="!isValidCircle" class="mt-2">
                Share with Circle
              </v-btn>
            </v-window-item>
          </v-window>
          <template v-else>
            <v-autocomplete
              density="compact"
              v-model="selectedUser"
              :items="userSuggestions"
              :loading="isLoadingUsers"
              :search-input.sync="usernameInput"
              label="Search Users"
              item-title="displayName"
              item-value="name"
              :rules="[validateUsername]"
              clearable
              :input-attrs="{ autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' }"
              :prepend-inner-icon="getUsernameIcon"
              :color="getUsernameColor"
              @update:search="handleUserSearch"
              @update:model-value="handleUserSelection"
              no-data-text="No friends found"
            >
              <template #item="{ props, item }">
                <v-list-item v-bind="props">
                  <v-list-item-subtitle>
                    @{{ item.raw.username }}
                  </v-list-item-subtitle>
                </v-list-item>
              </template>
            </v-autocomplete>
            <v-select
              density="compact"
              v-model="selectedPermission"
              :items="permissionOptions"
              label="Permission Level"
              class="mt-2"
            ></v-select>
            <v-btn block color="primary" @click="emitShareUser" :loading="sharing" :disabled="!isValidUsername" class="mt-2">
              Share with User
            </v-btn>
          </template>
          <v-divider class="my-4"></v-divider>
        </div>
        <!-- Scrollable current shares section -->
        <div v-if="currentAccesses.length > 0" style="flex: 1; overflow-y: auto; min-height: 0;">
          <div class="text-subtitle-1 mb-2">Current Shares</div>
          <v-list>
            <v-list-item v-for="access in currentAccesses as any[]" :key="access.name || ''" :title="accessTitle(access)" :prependIcon="getAccessIcon(access)" :subtitle="accessSubtitle(access)">
              <template #append>
                <div class="d-flex align-center gap-2">
                  <v-menu
                    v-model="accessMenuOpen[access.name || '']"
                    :close-on-content-click="false"
                    location="bottom"
                    offset-y
                    :disabled="!isAccessEditable(access)"
                  >
                    <template #activator="{ props }">
                      <v-chip
                        v-bind="props"
                        size="small"
                        :color="getPermissionDetails(access.level)?.color"
                        class="permission-chip d-flex align-center"
                        :disabled="accessPermissionLoading[access.name || ''] || !isAccessEditable(access)"
                        style="cursor: pointer; min-width: 120px;"
                      >
                        <span>{{ getPermissionDetails(access.level)?.title }}</span>
                        <v-icon end size="18" class="ml-1" v-if="isAccessEditable(access)">mdi-chevron-down</v-icon>
                        <v-progress-circular
                          v-if="accessPermissionLoading[access.name || '']"
                          indeterminate
                          size="16"
                          color="primary"
                          class="ml-1"
                        />
                      </v-chip>
                    </template>
                    <v-list>
                      <v-list-item
                        v-for="option in permissionOptions as any[]"
                        :key="option.value"
                        :value="option.value"
                        @click="emitPermissionChange(access, option.value); accessMenuOpen[access.name || ''] = false"
                        :disabled="access.level === option.value"
                      >
                        <v-list-item-title>{{ option.title }}</v-list-item-title>
                      </v-list-item>
                    </v-list>
                  </v-menu>
                  <v-btn 
                    v-if="access.state === 'ACCESS_STATE_PENDING' && access.acceptTarget === 'ACCEPT_TARGET_RESOURCE'" 
                    size="small" 
                    color="success" 
                    variant="outlined"
                    @click="emitApproveAccess(access.name || '')"
                    :disabled="!isAccessEditable(access)"
                  >
                    Approve
                  </v-btn>
                  <v-chip v-else-if="access.state === 'ACCESS_STATE_PENDING'" size="small" color="warning" variant="outlined">
                    Pending
                  </v-chip>
                  <v-btn 
                    icon="mdi-delete" 
                    variant="text" 
                    @click="emitRemoveAccess(access.name || '')"
                    :disabled="!isAccessEditable(access)"
                  ></v-btn>
                </div>
              </template>
            </v-list-item>
          </v-list>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="emitClose">
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import type { User } from '@/genapi/api/users/user/v1alpha1'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'
import type { ListUsersRequest } from '@/genapi/api/users/user/v1alpha1'
import type { ListCirclesRequest } from '@/genapi/api/circles/circle/v1alpha1'
import { userService, circleService } from '@/api/api'
import { useAuthStore } from '@/stores/auth'
import type { PermissionLevel } from '@/genapi/api/types'
import { useAlertStore } from '@/stores/alerts'
import { getPermissionValue } from '@/utils/permissions'
import Fuse from 'fuse.js'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const props = defineProps({
  modelValue: Boolean,
  title: { type: String, default: 'Share' },
  allowCircleShare: { type: Boolean, default: false },
  currentAccesses: { type: Array, required: true },
  userPermissionLevel: { type: String },
  sharing: Boolean,
  accessPermissionLoading: { type: Object, default: () => ({}) },
  allowPermissionOptions: { type: Array, required: true },
  disableCreateShare: { type: Boolean, default: false },
})

const emit = defineEmits([
  'update:modelValue',
  'share-user',
  'share-circle',
  'remove-access',
  'permission-change',
  'approve-access',
  'close',
])

const accessTab = ref('users')
const selectedPermission = ref<PermissionLevel>('PERMISSION_LEVEL_READ')
const accessMenuOpen = ref<Record<string, boolean>>({})
const permissions = [
  { title: 'Read Only', value: 'PERMISSION_LEVEL_READ' as PermissionLevel, color: 'grey' },
  { title: 'Read & Write', value: 'PERMISSION_LEVEL_WRITE' as PermissionLevel, color: 'primary' },
  { title: 'Admin', value: 'PERMISSION_LEVEL_ADMIN' as PermissionLevel, color: 'warning' },
]

const permissionOptions = computed(() => {
  // First filter by allowed permissions, then by user's permission level
  return props.allowPermissionOptions
    .filter(permission => getPermissionValue(permission as PermissionLevel) <= getPermissionValue(props.userPermissionLevel as PermissionLevel))
    .map(permission => {
      for (const p of permissions) {
        if (p.value === permission) {
          return p
        }
      }
      return null
    })
    .filter(Boolean)
})

function getPermissionDetails(permission: PermissionLevel) {
  for (const p of permissions) {
    if (p.value === permission) {
      return p
    }
  }
}

// Check if an access entry is editable based on permission levels
function isAccessEditable(access: any): boolean {
  const accessLevel = access.level
  const userLevel = props.userPermissionLevel
  return getPermissionValue(accessLevel as PermissionLevel) <= getPermissionValue(userLevel as PermissionLevel)
}

// Autocomplete state
const allUsers = ref<User[]>([])
const allCircles = ref<Circle[]>([])
const userSuggestions = ref<User[]>([])
const circleSuggestions = ref<Circle[]>([])
const isLoadingUsers = ref(false)
const isLoadingCircles = ref(false)
const userFuse = ref<Fuse<User> | null>(null)
const circleFuse = ref<Fuse<Circle> | null>(null)

// Selection state
const selectedUser = ref<string | null>(null)
const selectedCircle = ref<string | null>(null)
const usernameInput = ref('')
const circleInput = ref('')

let isValidUsername = ref(false)
let isValidCircle = ref(false)

function emitRemoveAccess(name: string) {
  emit('remove-access', name)
}

function emitApproveAccess(name: string) {
  emit('approve-access', name)
}
function emitPermissionChange(access: any, newLevel: string) {
  emit('permission-change', { access, newLevel })
}
function emitClose() {
  emit('update:modelValue', false)
  emit('close')
}
function getAccessIcon(access: any) {
  if (access.recipient?.user || access.recipient?.username) { // user
    return 'mdi-account'
  } else if (access.recipient?.circle) { // circle
    return 'mdi-account-group'
  }
  return ''
}
function accessTitle(access: any) {
    let title = ''

    if (access.recipient?.user) { // user
      if (access.recipient.user.givenName || access.recipient.user.familyName) { // user full name
        title = access.recipient.user.givenName + ' ' + access.recipient.user.familyName
        title = title.trim()
      } else if (access.recipient.user.username) { // user username
        title = access.recipient.user.username
      } 
    } else if (access.recipient?.circle) { // circle
        title = access.recipient.circle.title || access.recipient.circle.handle || ''
    } else if (access.recipient) {
      if (access.recipient.givenName || access.recipient.familyName) { // user full name
        title = access.recipient.givenName + ' ' + access.recipient.familyName
        title = title.trim()
      } else if (access.recipient.username) { // user username
        title = access.recipient.username
      } else if (access.recipient.title) { // circle title or handle
        title = access.recipient.title || access.recipient.handle || ''
      }
    }
  return title
}
function accessSubtitle(access: any) {
  // if it is user but there is a given/family name use the username else empty string
  // if it is a circle use the handle
  let subtitle = ''
  if (access.recipient?.user) { // user
    if (access.recipient.user.givenName || access.recipient.user.familyName) { // user username
      subtitle = access.recipient.user.username || ''
    }
  } else if (access.recipient?.circle) { // circle
      subtitle = access.recipient.circle.handle || ''
  } else if (access.recipient) {
    if (access.recipient.givenName || access.recipient.familyName) { // user full name
      subtitle = access.recipient.username || ''
    } else if (access.recipient.handle) { // circle title or handle
      subtitle = access.recipient.handle || ''
    }
  }
  
  return subtitle
}

// Format display names for autocomplete
function formatUserDisplay(user: User): User & { displayName: string } {
  const fullName = [user.givenName, user.familyName].filter(Boolean).join(' ')
  const displayName = fullName || user.username || 'Unknown User'
  return { ...user, displayName }
}

function formatCircleDisplay(circle: Circle): Circle & { displayName: string } {
  const displayName = circle.title || circle.handle || 'Unknown Circle'
  return { ...circle, displayName }
}

// Initialize Fuse instances
function initializeFuseInstances() {
  // User search configuration
  const userSearchKeys = [
    { name: 'username', weight: 0.4 },
    { name: 'givenName', weight: 0.3 },
    { name: 'familyName', weight: 0.3 }
  ]

  // Circle search configuration
  const circleSearchKeys = [
    { name: 'handle', weight: 0.5 },
    { name: 'title', weight: 0.5 }
  ]

  const fuseOptions = {
    threshold: 0.3,
    includeScore: true,
    minMatchCharLength: 2,
    shouldSort: true
  }

  userFuse.value = new Fuse(allUsers.value, {
    ...fuseOptions,
    keys: userSearchKeys
  })

  circleFuse.value = new Fuse(allCircles.value, {
    ...fuseOptions,
    keys: circleSearchKeys
  })
}

// Fetch all users and circles
async function fetchUsers() {
  isLoadingUsers.value = true
  try {
    const request: ListUsersRequest = {
      pageSize: 1000,
      pageToken: undefined,
      filter: 'state = 200',
      parent: undefined
    }
    const response = await userService.ListUsers(request)
    
    if (response.users) {
      // Filter out current user and format for display
      allUsers.value = response.users
        .filter(user => user.name !== authStore.user.name)
        .map(formatUserDisplay)
    }
  } catch (error) {
    alertStore.addAlert(error instanceof Error ? "Unable to fetch users\n" + error.message : String(error), 'error')
  } finally {
    isLoadingUsers.value = false
  }
}

async function fetchCircles() {
  isLoadingCircles.value = true
  try {
    const request: ListCirclesRequest = {
      pageSize: 1000,
      pageToken: undefined,
      filter: 'state = 200',
      parent: undefined
    }
    const response = await circleService.ListCircles(request)
    
    if (response.circles) {
      allCircles.value = response.circles.map(formatCircleDisplay)
    }
  } catch (error) {
    alertStore.addAlert(error instanceof Error ? "Unable to fetch circles\n" + error.message : String(error), 'error')
  } finally {
    isLoadingCircles.value = false
  }
}

// Search functions
function searchUsers(query: string) {
  if (!userFuse.value || !query?.trim()) {
    userSuggestions.value = allUsers.value.slice(0, 15)
    return
  }

  const results = userFuse.value.search(query.trim())
  userSuggestions.value = results
    .slice(0, 15)
    .map(result => result.item)
}

function searchCircles(query: string) {
  if (!circleFuse.value || !query.trim()) {
    circleSuggestions.value = allCircles.value.slice(0, 15)
    return
  }

  const results = circleFuse.value.search(query.trim())
  circleSuggestions.value = results
    .slice(0, 15)
    .map(result => result.item)
}

// Input handlers
function handleUserSearch(query: string) {
  searchUsers(query)
}

function handleCircleSearch(query: string) {
  searchCircles(query)
}

function handleUserSelection(user: string | null) {
  selectedUser.value = user
  isValidUsername.value = !!user
}

function handleCircleSelection(circle: string | null) {
  selectedCircle.value = circle
  isValidCircle.value = !!circle
}

// Validation functions
function validateUsername(value: string): boolean | string {
  return true
}

function validateCircle(value: string): boolean | string {
  return true
}

// Computed properties for icons and colors
const getUsernameIcon = computed(() => {
  if (isLoadingUsers.value) return 'mdi-loading'
  if (!selectedUser.value) return undefined
  return isValidUsername.value ? 'mdi-check-circle' : 'mdi-close-circle'
})

const getUsernameColor = computed(() => {
  if (isLoadingUsers.value) return undefined
  if (!selectedUser.value) return undefined
  return isValidUsername.value ? 'success' : 'error'
})

const getCircleIcon = computed(() => {
  if (isLoadingCircles.value) return 'mdi-loading'
  if (!selectedCircle.value) return undefined
  return isValidCircle.value ? 'mdi-check-circle' : 'mdi-close-circle'
})

const getCircleColor = computed(() => {
  if (isLoadingCircles.value) return undefined
  if (!selectedCircle.value) return undefined
  return isValidCircle.value ? 'success' : 'error'
})

// Share functions
function emitShareUser() {
  if (selectedUser.value) {
    emit('share-user', { userName: selectedUser.value || '', permission: selectedPermission.value })
    selectedUser.value = null
    usernameInput.value = ''
    isValidUsername.value = false
  }
}

function emitShareCircle() {
  if (selectedCircle.value) {
    emit('share-circle', { circleName: selectedCircle.value || '', permission: selectedPermission.value })
    selectedCircle.value = null
    circleInput.value = ''
    isValidCircle.value = false
  }
}

// Initialize data on mount
onMounted(async () => {
  await Promise.all([
    fetchUsers(),
    props.allowCircleShare ? fetchCircles() : Promise.resolve()
  ])
  initializeFuseInstances()
  
  // Set initial suggestions
  userSuggestions.value = allUsers.value.slice(0, 15)
  if (props.allowCircleShare) {
    circleSuggestions.value = allCircles.value.slice(0, 15)
  }
})

</script> 