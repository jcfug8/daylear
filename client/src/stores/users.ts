import { ref } from 'vue'
import { defineStore } from 'pinia'
import { userService } from '@/api/api'
import type { User } from '@/genapi/api/users/user/v1alpha1'

export const useUsersStore = defineStore('users', () => {
  const users = ref<User[]>([])
  const friends = ref<User[]>([])
  const pendingFriends = ref<User[]>([])
  const publicUsers = ref<User[]>([])
  const loading = ref(false)

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
  }
}) 