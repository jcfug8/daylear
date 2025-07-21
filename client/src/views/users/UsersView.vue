<template>
  <ListTabsPage :tabs="tabs" ref="tabsPage">
    <template #friends="{ items, loading }">
      <UserGrid :users="items" :loading="usersStore.loading" empty-text="No friends found." />
    </template>
    <template #pending="{ items, loading }">
      <UserGrid :users="items" :loading="usersStore.loading" empty-text="No pending requests found." />
    </template>
    <template #explore="{ items, loading }">
      <UserGrid :users="items" :loading="usersStore.loading" empty-text="No users found." />
    </template>
  </ListTabsPage>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import { useUsersStore } from '@/stores/users'
import UserGrid from '@/components/UserGrid.vue'

const usersStore = useUsersStore()
const tabsPage = ref()

const tabs = [
  {
    label: 'Friends',
    value: 'friends',
    icon: 'mdi-account-multiple',
    loader: async () => {
      await usersStore.loadFriends()
      return [...usersStore.friends]
    },
  },
  {
    label: 'Pending',
    value: 'pending',
    icon: 'mdi-account-clock',
    loader: async () => {
      await usersStore.loadPendingFriends()
      return [...usersStore.pendingFriends]
    },
  },
  {
    label: 'Explore',
    value: 'explore',
    icon: 'mdi-compass-outline',
    loader: async () => {
      await usersStore.loadPublicUsers()
      return [...usersStore.publicUsers]
    },
  },
]
</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style> 