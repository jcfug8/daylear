<template>
  <ListTabsPage
    :tabs="tabs"
  >
    <template #my="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
        <h2>My Circles</h2>
      </div>
      <v-list v-if="items.length > 0">
        <v-list-item
          v-for="circle in items"
          :key="circle.name"
          :title="circle.title"
          :subtitle="circle.name"
          :to="{ name: 'circle', params: { circleId: circle.name } }"
        />
      </v-list>
      <div v-else>No circles found.</div>
    </template>
    <template #shared-accepted="{ items, loading }">
      <v-list v-if="items.length > 0">
        <v-list-item
          v-for="circle in items"
          :key="circle.name"
          :title="circle.title"
          :subtitle="circle.name"
          :to="{ name: 'circle', params: { circleId: circle.name } }"
        />
      </v-list>
      <div v-else>No accepted shared circles found.</div>
    </template>
    <template #shared-pending="{ items, loading }">
      <v-list v-if="items.length > 0">
        <v-list-item
          v-for="circle in items"
          :key="circle.name"
          :title="circle.title"
          :subtitle="circle.name"
          :to="{ name: 'circle', params: { circleId: circle.name } }"
        >
          <template #append>
            <v-btn color="success" @click.stop.prevent="acceptCircleAccess(circle)" :loading="acceptingCircleId === circle.name">
              Accept
            </v-btn>
          </template>
        </v-list-item>
      </v-list>
      <div v-else>No pending shared circles found.</div>
    </template>
    <template #explore="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
        <h2>Explore Public Circles</h2>
      </div>
      <v-list v-if="items.length > 0">
        <v-list-item
          v-for="circle in items"
          :key="circle.name"
          :title="circle.title"
          :subtitle="circle.name"
          :to="{ name: 'circle', params: { circleId: circle.name } }"
        />
      </v-list>
      <div v-else>No public circles found.</div>
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
import { circleAccessService } from '@/api/api'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'

const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()

const acceptingCircleId = ref<string | null>(null)

const tabs = [
  {
    label: 'My Circles',
    value: 'my',
    loader: async () => {
      await circlesStore.loadMyCircles()
      return [...circlesStore.circles]
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
          return [...circlesStore.circles]
        },
      },
      {
        label: 'Pending',
        value: 'pending',
        loader: async () => {
          await circlesStore.loadSharedCircles(100)
          return [...circlesStore.circles]
        },
      },
    ],
  },
  {
    label: 'Explore Circles',
    value: 'explore',
    loader: async () => {
      await circlesStore.loadPublicCircles()
      return [...circlesStore.circles]
    },
  },
]

async function acceptCircleAccess(circle: Circle) {
  if (!circle.name) return
  acceptingCircleId.value = circle.name
  try {
    await circleAccessService.AcceptAccess({ name: circle.name })
    // Reload pending circles after accepting
    const pendingTab = tabs.find(t => t.value === 'shared')?.subTabs?.find(s => s.value === 'pending')
    if (pendingTab && pendingTab.loader) await pendingTab.loader()
  } catch (error) {
    // Optionally show a notification
  } finally {
    acceptingCircleId.value = null
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
