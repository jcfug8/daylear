import { ref } from 'vue'
import { defineStore } from 'pinia'
import { userService, userSettingsService, userAccessService } from '@/api/api'
import type { User, UserSettings } from '@/genapi/api/users/user/v1alpha1'

export const useUsersStore = defineStore('users', () => {
  const users = ref<User[]>([])
  const friends = ref<User[]>([])
  const pendingFriends = ref<User[]>([])
  const publicUsers = ref<User[]>([])
  const loading = ref(false)
  const currentUser = ref<User | null>(null)
  const currentUserSettings = ref<UserSettings | null>(null)

  async function loadUsers(filter?: string) {
    loading.value = true
    try {
      const res = await userService.ListUsers({ pageSize: 100, pageToken: undefined, filter })
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

  async function loadFriends() {
    const result = await loadUsers('state = 200')
    friends.value = result
  }

  async function loadPendingFriends() {
    const result = await loadUsers('state = 100')
    pendingFriends.value = result
  }

  async function loadPublicUsers() {
    const result = await loadUsers('visibility = 1')
    publicUsers.value = result
  }

  async function loadUser(name: string) {
    loading.value = true
    try {
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
          visibility: editUser.visibility,
          imageUri: editUser.imageUri,
          access: undefined,
          bio: editUser.bio,
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
      // Optionally reload pending friends
      await loadPendingFriends()
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
      // Optionally reload pending friends
      await loadPendingFriends()
    } catch (error) {
      console.error('Failed to decline user access:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    users,
    friends,
    pendingFriends,
    publicUsers,
    loading,
    loadUsers,
    loadFriends,
    loadPendingFriends,
    loadPublicUsers,
    currentUser,
    currentUserSettings,
    loadUser,
    loadUserSettings,
    updateUser,
    updateUserSettings,
    acceptUserAccess,
    declineUserAccess,
  }
}) 