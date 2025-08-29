import { ref } from 'vue'
import { defineStore } from 'pinia'
import { userService, userSettingsService, userAccessService, accessKeyService } from '@/api/api'
import type { User, UserSettings, AccessKey } from '@/genapi/api/users/user/v1alpha1'

export const useUsersStore = defineStore('users', () => {
  const users = ref<User[]>([])
  const friends = ref<User[]>([])
  const pendingFriends = ref<User[]>([])
  const publicUsers = ref<User[]>([])
  const loading = ref(false)
  const currentUser = ref<User | null>(null)
  const currentUserSettings = ref<UserSettings | null>(null)
  const accessKeys = ref<AccessKey[]>([])
  const accessKeysLoading = ref(false)

  async function loadUsers(parent: string, filter?: string) {
    loading.value = true
    try {
      users.value = []
      const res = await userService.ListUsers({ parent, pageSize: 100, pageToken: undefined, filter })
      users.value = res.users || []
      return users.value
    } catch (error) {
      users.value = []
      console.error('Failed to load users:', error)
      return []
    } finally {
      loading.value = false
    }
  }

  async function loadFriends(parent: string) {
    friends.value = []
    const result = await loadUsers(parent, 'state = 200 OR favorited = true')
    friends.value = result.reduce((acc, user) => {
      if (user.name !== parent) {
        acc.push(user)
      }
      return acc
    }, [] as User[])
  }

  async function loadPendingFriends(parent: string) {
    pendingFriends.value = []
    const result = await loadUsers(parent, 'state = 100')
    pendingFriends.value = result
  }

  async function loadPublicUsers(parent: string) {
    publicUsers.value = []
    const result = await loadUsers(parent, '')
    publicUsers.value = result
  }

  async function loadUser(name: string) {
    loading.value = true
    try {
      currentUser.value = null
      currentUser.value = await userService.GetUser({ name })
    } catch (error) {
      currentUser.value = null
      console.error('Failed to load user:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function loadUserSettings(name: string) {
    loading.value = true
    try {
      currentUserSettings.value = await userSettingsService.GetUserSettings({ name: `${name}/settings` })
    } catch (error) {
      currentUserSettings.value = null
      console.error('Failed to load user settings:', error)
    } finally {
      loading.value = false
    }
  }

  async function updateUser(editUser: User) {
    loading.value = true
    try {
      currentUser.value = await userService.UpdateUser({
        user: {
          name: editUser.name,
          username: editUser.username,
          givenName: editUser.givenName,
          familyName: editUser.familyName,
          imageUri: editUser.imageUri,
          access: undefined,
          bio: editUser.bio,
          favorited: editUser.favorited,
        },
        updateMask: undefined,
      })
      return currentUser.value
    } catch (error) {
      console.error('Failed to update user:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function updateUserSettings(editUserSettings: UserSettings) {
    loading.value = true
    try {
      currentUserSettings.value = await userSettingsService.UpdateUserSettings({
        userSettings: {
          name: editUserSettings.name,
          email: editUserSettings.email,
        },
        updateMask: undefined,
      })
      return currentUserSettings.value
    } catch (error) {
      console.error('Failed to update user settings:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function acceptUserAccess(accessName: string) {
    loading.value = true
    try {
      await userAccessService.AcceptAccess({ name: accessName })
    } catch (error) {
      console.error('Failed to accept user access:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function declineUserAccess(accessName: string) {
    loading.value = true
    try {
      await userAccessService.DeleteAccess({ name: accessName })
    } catch (error) {
      console.error('Failed to decline user access:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function loadAccessKeys(parent: string) {
    accessKeysLoading.value = true
    try {
      accessKeys.value = []
      const res = await accessKeyService.ListAccessKeys({ parent, pageSize: 100, pageToken: undefined, filter: undefined })
      accessKeys.value = res.accessKeys || []
      return accessKeys.value
    } catch (error) {
      accessKeys.value = []
      console.error('Failed to load access keys:', error)
      return []
    } finally {
      accessKeysLoading.value = false
    }
  }

  async function createAccessKey(parent: string, title: string, description?: string) {
    accessKeysLoading.value = true
    try {
      const newAccessKey = await accessKeyService.CreateAccessKey({
        parent,
        accessKey: {
          name: undefined,
          title,
          description,
          unencryptedAccessKey: undefined,
        }
      })
      // Refresh the access keys list
      await loadAccessKeys(parent)
      return newAccessKey
    } catch (error) {
      console.error('Failed to create access key:', error)
      throw error
    } finally {
      accessKeysLoading.value = false
    }
  }

  async function deleteAccessKey(name: string) {
    accessKeysLoading.value = true
    try {
      await accessKeyService.DeleteAccessKey({ name })
      // Remove from local state
      accessKeys.value = accessKeys.value.filter(key => key.name !== name)
    } catch (error) {
      console.error('Failed to delete access key:', error)
      throw error
    } finally {
      accessKeysLoading.value = false
    }
  }

  async function updateAccessKey(accessKey: AccessKey, updateMask?: string[]) {
    accessKeysLoading.value = true
    try {
      const updatedAccessKey = await accessKeyService.UpdateAccessKey({
        accessKey,
        updateMask: updateMask ? updateMask.join(',') : undefined,
      })
      // Update in local state
      const index = accessKeys.value.findIndex(key => key.name === accessKey.name)
      if (index !== -1) {
        accessKeys.value[index] = updatedAccessKey
      }
      return updatedAccessKey
    } catch (error) {
      console.error('Failed to update access key:', error)
      throw error
    } finally {
      accessKeysLoading.value = false
    }
  }

  return {
    users,
    friends,
    pendingFriends,
    publicUsers,
    loading,
    currentUser,
    currentUserSettings,
    accessKeys,
    accessKeysLoading,
    loadUsers,
    loadFriends,
    loadPendingFriends,
    loadPublicUsers,
    loadUser,
    loadUserSettings,
    updateUser,
    updateUserSettings,
    acceptUserAccess,
    declineUserAccess,
    loadAccessKeys,
    createAccessKey,
    deleteAccessKey,
    updateAccessKey,
  }
}) 