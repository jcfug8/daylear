<template>
  <ListTabsPage
    ref="tabsPage"
    :tabs="tabs"
  >
    <template #my="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
      </div>
      <CircleGrid :circles="items" :loading="loading" />
    </template>
    <template #pending="{ items, loading }">
      <CircleGrid :circles="items" :loading="loading" @accept="acceptCircleAccess" :acceptingCircleId="acceptingCircleId" @decline="onDeclineCircle" />
      <div v-if="!loading && items.length === 0">No pending shared circles found.</div>
    </template>
    <template #explore="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
      </div>
      <CircleGrid :circles="items" :loading="loading" />
    </template>
    <template #fab>
      <v-btn
        color="primary"
        icon="mdi-plus"
        style="position: fixed; bottom: 16px; right: 16px"
        :to="{ name: 'circleCreate' }"
      ></v-btn>
    </template>
  </ListTabsPage>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCirclesStore } from '@/stores/circles'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import CircleGrid from '@/components/CircleGrid.vue'
import { circleAccessService } from '@/api/api'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'

const circlesStore = useCirclesStore()

const acceptingCircleId = ref<string | null>(null)
const tabsPage = ref()

const tabs = [
  {
    label: 'My Circles',
    value: 'my',
    icon: 'mdi-account-circle',
    loader: async () => {
      await circlesStore.loadMyCircles()
      return [...circlesStore.myCircles]
    },
  },
  {
    label: 'Pending',
    value: 'pending',
    icon: 'mdi-clock-outline',
    loader: async () => {
      await circlesStore.loadPendingCircles()
      return [...circlesStore.sharedPendingCircles]
    },
  },
  {
    label: 'Explore',
    value: 'explore',
    icon: 'mdi-compass-outline',
    loader: async () => {
      await circlesStore.loadPublicCircles()
      return [...circlesStore.publicCircles]
    },
  },
]

async function acceptCircleAccess(circle: Circle) {
  if (!circle.name) return
  acceptingCircleId.value = circle.name
  try {
    await circleAccessService.AcceptAccess({ name: circle.circleAccess?.name })
    tabsPage.value?.reloadTab('pending')
  } catch (error) {
    // Optionally show a notification
  } finally {
    acceptingCircleId.value = null
  }
}

async function onDeclineCircle(circle: Circle) {
  if (!circle.circleAccess?.name) return
  try {
    await circleAccessService.DeleteAccess({ name: circle.circleAccess.name })
    // Reload only the pending subtab
    tabsPage.value?.reloadTab('pending')
  } catch (error) {
    // Optionally show a notification
  }
}
</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style>
