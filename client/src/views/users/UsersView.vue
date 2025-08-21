<template>
  <ListTabsPage :tabs="tabs" ref="tabsPage">
    <template #friends="{ items }">
      <UserGrid :users="items as User[]" :loading="usersStore.loading" empty-text="No friends found." />
    </template>
    <template #pending="{ items }">
      <UserGrid
        :users="items as User[]"
        :loading="usersStore.loading"
        empty-text="No pending requests found."
        show-actions
        :acceptingUserId="acceptingUserId || undefined"
        @accept="onAcceptUser"
        @decline="onDeclineUser"
      />
    </template>
    <template #explore="{ items }">
      <UserGrid :users="items as User[]" :loading="usersStore.loading" empty-text="No users found." />
    </template>
  </ListTabsPage>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import { useUsersStore } from '@/stores/users'
import UserGrid from '@/components/UserGrid.vue'
import { useAuthStore } from '@/stores/auth'
import type { User } from '@/genapi/api/users/user/v1alpha1'

const usersStore = useUsersStore()
const tabsPage = ref()
const acceptingUserId = ref<string | null>(null)
const authStore = useAuthStore()

const tabs = [
  {
    label: 'Friends',
    value: 'friends',
    icon: 'mdi-account-multiple',
    loader: async () => {
      await usersStore.loadFriends(authStore.user?.name || '')
      return [...usersStore.friends]
    },
  },
  {
    label: 'Pending',
    value: 'pending',
    icon: 'mdi-account-clock',
    loader: async () => {
      await usersStore.loadPendingFriends(authStore.user?.name || '')
      return [...usersStore.pendingFriends]
    },
  },
  {
    label: 'Explore',
    value: 'explore', 
    icon: 'mdi-compass-outline',
    loader: async () => {
      await usersStore.loadPublicUsers(authStore.user?.name || '')
      return [...usersStore.publicUsers]
    },
  },
]

async function onAcceptUser(user: User) {
  if (!user.access?.name || !user.name) return
  acceptingUserId.value = user.name
  try {
    await usersStore.acceptUserAccess(user.access.name)
    tabsPage.value?.reloadTab('pending')
  } catch (err) {
    // Optionally show a notification
  } finally {
    acceptingUserId.value = null
  }
}

async function onDeclineUser(user: User) {
  if (!user.access?.name) return
  try {
    await usersStore.declineUserAccess(user.access.name)
    tabsPage.value?.reloadTab('pending')
  } catch (err) {
    // Optionally show a notification
  }
}
</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style> 