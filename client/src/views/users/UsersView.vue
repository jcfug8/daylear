<template>
  <ListTabsPage :tabs="tabs" ref="tabsPage">
    <template #friends="{ items, loading }">
      <UserGrid :users="items" :loading="loading" empty-text="No friends found." />
    </template>
    <template #pending="{ items, loading }">
      <UserGrid :users="items" :loading="loading" empty-text="No pending requests found." @accept="acceptAccess" @decline="declineAccess" :show-actions="true" />
    </template>
    <template #explore="{ items, loading }">
      <UserGrid :users="items" :loading="loading" empty-text="No users found." />
    </template>
  </ListTabsPage>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import { userService, userAccessService } from '@/api/api'
import type { User, Access } from '@/genapi/api/users/user/v1alpha1'
import UserGrid from '@/components/UserGrid.vue'

const tabsPage = ref()
const acceptingUserId = ref<string | null>(null)

const tabs = [
  {
    label: 'Friends',
    value: 'friends',
    icon: 'mdi-account-multiple',
    loader: async () => {
      // List accesses with state ACCEPTED
      const res = await userAccessService.ListAccesses({ parent: 'users/-', filter: 'state=ACCESS_STATE_ACCEPTED', pageSize: undefined, pageToken: undefined })
      // Extract recipient user info
      return (res.accesses || []).map(a => a.recipient)
    },
  },
  {
    label: 'Pending',
    value: 'pending',
    icon: 'mdi-account-clock',
    loader: async () => {
      const res = await userAccessService.ListAccesses({ parent: 'users/-', filter: 'state=ACCESS_STATE_PENDING', pageSize: undefined, pageToken: undefined })
      // Attach access name for accept/decline
      return (res.accesses || []).map(a => ({ ...a.recipient, accessName: a.name }))
    },
  },
  {
    label: 'Explore',
    value: 'explore',
    icon: 'mdi-compass-outline',
    loader: async () => {
      const res = await userService.ListUsers({ pageSize: 100, pageToken: undefined, filter: undefined })
      return res.users || []
    },
  },
]

async function acceptAccess(user: any) {
  if (!user.accessName) return
  acceptingUserId.value = user.name
  try {
    await userAccessService.AcceptAccess({ name: user.accessName })
    tabsPage.value?.reloadTab('pending')
    tabsPage.value?.reloadTab('friends')
  } finally {
    acceptingUserId.value = null
  }
}

async function declineAccess(user: any) {
  if (!user.accessName) return
  try {
    await userAccessService.DeleteAccess({ name: user.accessName })
    tabsPage.value?.reloadTab('pending')
  } catch (e) {}
}
</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style> 