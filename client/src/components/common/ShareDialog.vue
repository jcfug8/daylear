<template>
  <v-dialog :model-value="modelValue" @update:model-value="emitClose" max-width="500">
    <v-card>
      <v-card-title class="text-h5">
        {{ title }}
      </v-card-title>
      <v-card-text>
        <v-tabs v-if="allowCircleShare" v-model="accessTab" class="mb-4">
          <v-tab value="users">Share with User</v-tab>
          <v-tab value="circles">Share with Circle</v-tab>
        </v-tabs>
        <v-window v-if="allowCircleShare" v-model="accessTab">
          <v-window-item value="users">
            <v-text-field v-model="usernameInput" label="Enter Username" :rules="[validateUsername]"
              clearable
              :input-attrs="{ autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' }"
              :prepend-inner-icon="getUsernameIcon" :color="getUsernameColor" :loading="isLoadingUsername"
              @update:model-value="handleUsernameInput"></v-text-field>
            <v-select
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
            <v-text-field v-model="circleInput" label="Enter Circle Handle" :rules="[validateCircle]"
              clearable
              :input-attrs="{ autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' }"
              :prepend-inner-icon="getCircleIcon" :color="getCircleColor" :loading="isLoadingCircle"
              @update:model-value="handleCircleInput"></v-text-field>
            <v-select
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
          <v-text-field v-model="usernameInput" label="Enter Username" :rules="[validateUsername]"
            :prepend-inner-icon="getUsernameIcon" :color="getUsernameColor" :loading="isLoadingUsername"
            @update:model-value="handleUsernameInput"></v-text-field>
          <v-select
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
        <div v-if="currentAccesses.length > 0">
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
                  >
                    <template #activator="{ props }">
                      <v-chip
                        v-bind="props"
                        size="small"
                        :color="hasWritePermission(access.level) ? 'primary' : 'grey'"
                        class="permission-chip d-flex align-center"
                        :disabled="accessPermissionLoading[access.name || '']"
                        style="cursor: pointer; min-width: 120px;"
                      >
                        <span>{{ hasWritePermission(access.level) ? 'Read & Write' : 'Read Only' }}</span>
                        <v-icon end size="18" class="ml-1">mdi-chevron-down</v-icon>
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
                  <v-chip v-if="access.state === 'ACCESS_STATE_PENDING'" size="small" color="warning" variant="outlined">
                    Pending
                  </v-chip>
                  <v-btn icon="mdi-delete" variant="text" @click="emitRemoveAccess(access.name || '')"></v-btn>
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
import { ref, computed, onBeforeUnmount } from 'vue'
import type { User } from '@/genapi/api/users/user/v1alpha1'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'
import type { ListUsersRequest } from '@/genapi/api/users/user/v1alpha1'
import type { ListCirclesRequest } from '@/genapi/api/circles/circle/v1alpha1'
import { userService, circleService } from '@/api/api'
import { useAuthStore } from '@/stores/auth'
import type { PermissionLevel } from '@/genapi/api/types'
import { useAlertStore } from '@/stores/alerts'

const authStore = useAuthStore()
const alertStore = useAlertStore()

const props = defineProps({
  modelValue: Boolean,
  title: { type: String, default: 'Share' },
  allowCircleShare: { type: Boolean, default: false },
  currentAccesses: { type: Array, required: true },
  sharing: Boolean,
  accessPermissionLoading: { type: Object, default: () => ({}) },
  hasWritePermission: { type: Function, required: true },
})

const emit = defineEmits([
  'update:modelValue',
  'share-user',
  'share-circle',
  'remove-access',
  'permission-change',
  'close',
])

const accessTab = ref('users')
const selectedPermission = ref<PermissionLevel>('PERMISSION_LEVEL_READ')
const accessMenuOpen = ref<Record<string, boolean>>({})
const permissionOptions = [
  { title: 'Read Only', value: 'PERMISSION_LEVEL_READ' as PermissionLevel },
  { title: 'Read & Write', value: 'PERMISSION_LEVEL_WRITE' as PermissionLevel },
]

function emitRemoveAccess(name: string) {
  emit('remove-access', name)
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

// *** Circle Sharing ***

const selectedCircle = ref<Circle | null>(null)
const circleInput = ref('')

var isValidCircle = ref(false)
var isLoadingCircle = ref(false)

let circleDebounceTimer: number | null = null

onBeforeUnmount(() => {
  if (circleDebounceTimer) {
    clearTimeout(circleDebounceTimer)
  }
})

const getCircleIcon = computed(() => {
  if (isLoadingCircle.value) return 'mdi-loading'
  if (!circleInput.value) return undefined
  return isValidCircle.value ? 'mdi-check-circle' : 'mdi-close-circle'
})

const getCircleColor = computed(() => {
  if (isLoadingCircle.value) return undefined
  if (!circleInput.value) return undefined
  return isValidCircle.value ? 'success' : 'error'
})

function validateCircle(value: string): boolean | string {
  return true
}

function emitShareCircle() {
    circleInput.value = circleInput.value
    selectedPermission.value = selectedPermission.value
    checkCircle(circleInput.value).then(() => {
        if (isValidCircle.value) {
            emit('share-circle', { circleName: selectedCircle.value?.name || '', permission: selectedPermission.value })
            selectedCircle.value = null
            circleInput.value = ''
            isValidCircle.value = false
            isLoadingCircle.value = false
        }
    })
}

function handleCircleInput(value: string) {
    if (circleDebounceTimer) {
    clearTimeout(circleDebounceTimer)
  }
  circleDebounceTimer = window.setTimeout(() => {
    checkCircle(value)
  }, 300)
}

async function checkCircle(circleHandle: string) {
  if (!circleHandle) {
    isValidCircle.value = false
    selectedCircle.value = null
    return
  }

  isLoadingCircle.value = true
  try {
    const request: ListCirclesRequest = {
      filter: `handle = "${circleHandle}"`,
      pageSize: 1,
      pageToken: undefined
    }
    const response = await circleService.ListCircles(request)

    if (response.circles?.length === 1) {
      selectedCircle.value = response.circles[0]
      isValidCircle.value = true
    } else {
      selectedCircle.value = null
      isValidCircle.value = false
    }
  } catch (error) {
    alertStore.addAlert(error instanceof Error ? "Unable to check circle\n" + error.message : String(error), 'error')
    selectedCircle.value = null
    isValidCircle.value = false
  } finally {
    isLoadingCircle.value = false
  }
}

// *** User Sharing ***

const selectedUser = ref<User | null>(null)
const usernameInput = ref('')

let isValidUsername = ref(false)
let isLoadingUsername = ref(false)

let usernameDebounceTimer: number | null = null

onBeforeUnmount(() => {
  if (usernameDebounceTimer) {
    clearTimeout(usernameDebounceTimer)
  }
})

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

function validateUsername(value: string): boolean | string {
  return true
}

function emitShareUser() {
    usernameInput.value = usernameInput.value
    selectedPermission.value = selectedPermission.value
    checkUsername(usernameInput.value).then(() => {
        if (isValidUsername.value) {
            emit('share-user', { userName: selectedUser.value?.name || '', permission: selectedPermission.value })
            selectedUser.value = null
            usernameInput.value = ''
            isValidUsername.value = false
            isLoadingUsername.value = false
        }
    })
}

function handleUsernameInput(value: string) {
  if (usernameDebounceTimer) {
    clearTimeout(usernameDebounceTimer)
  }
  usernameDebounceTimer = window.setTimeout(() => {
    checkUsername(value)
  }, 300)
}

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

    if (response.users?.length === 1 && response.users[0].name !== authStore.user.name) {
      selectedUser.value = response.users[0]
      isValidUsername.value = true
    } else {
      selectedUser.value = null
      isValidUsername.value = false
    }
  } catch (error) {
    alertStore.addAlert(error instanceof Error ? "Unable to check username\n" + error.message : String(error), 'error')
    selectedUser.value = null
    isValidUsername.value = false
  } finally {
    isLoadingUsername.value = false
  }
}

</script> 