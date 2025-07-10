<template>
  <ListTabsPage
    ref="tabsPage"
    :tabs="tabs"
  >
    <template #my="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
        <h2>My Circles</h2>
      </div>
      <CircleGrid :circles="items" :loading="loading" />
    </template>
    <template #shared-accepted="{ items, loading }">
      <CircleGrid :circles="items" :loading="loading" />
    </template>
    <template #shared-pending="{ items, loading }">
      <CircleGrid :circles="items" :loading="loading" @accept="acceptCircleAccess" :acceptingCircleId="acceptingCircleId" @decline="onDeclineCircle" />
      <div v-if="!loading && items.length === 0">No pending shared circles found.</div>
    </template>
    <template #explore="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
        <h2>Explore Public Circles</h2>
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
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import CircleGrid from '@/components/CircleGrid.vue'
import { circleAccessService } from '@/api/api'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'

const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()

const acceptingCircleId = ref<string | null>(null)
const tabsPage = ref()

const tabs = [
  {
    label: 'My Circles',
    value: 'my',
    loader: async () => {
      await circlesStore.loadMyCircles()
      return [...circlesStore.myCircles]
    },
  },
  {
    label: 'Shared Circles',
    value: 'shared',
    subTabs: [
      {
        label: 'Accepted',
        value: 'accepted',
        loader: async () => {
          await circlesStore.loadSharedCircles(200)
          return [...circlesStore.sharedAcceptedCircles]
        },
      },
      {
        label: 'Pending',
        value: 'pending',
        loader: async () => {
          await circlesStore.loadSharedCircles(100)
          return [...circlesStore.sharedPendingCircles]
        },
      },
    ],
  },
  {
    label: 'Explore Circles',
    value: 'explore',
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
    // Reload both accepted and pending subtabs
    tabsPage.value?.reloadTab('shared', 'accepted')
    tabsPage.value?.reloadTab('shared', 'pending')
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
    tabsPage.value?.reloadTab('shared', 'pending')
  } catch (error) {
    // Optionally show a notification
  }
}

breadcrumbStore.setBreadcrumbs([
  { title: 'Circles', to: { name: 'circles' } },
])
</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style>
