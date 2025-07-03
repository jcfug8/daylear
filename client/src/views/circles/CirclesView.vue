<template>
  <v-container>
    <v-tabs v-model="activeTab" align-tabs="center" color="primary" grow>
      <v-tab value="my">My Circles</v-tab>
      <v-tab value="shared">Shared Circles</v-tab>
      <v-tab value="explore">Explore Circles</v-tab>
    </v-tabs>
    <v-card-text>
      <v-tabs-window v-model="activeTab">
        <!-- My Circles Tab -->
        <v-tabs-window-item value="my">
          <div class="d-flex justify-space-between align-center mb-4">
            <h2>My Circles</h2>
          </div>
          <v-list v-if="myCircles.length > 0">
            <v-list-item
              v-for="circle in myCircles"
              :key="circle.name"
              :title="circle.title"
              :subtitle="circle.name"
              :to="{ name: 'circle', params: { circleId: circle.name } }"
            />
          </v-list>
          <div v-else>No circles found.</div>
        </v-tabs-window-item>
        <!-- Shared Circles Tab -->
        <v-tabs-window-item value="shared">
          <div class="mb-4">
            <h2 class="mb-4">Shared Circles</h2>
            <v-tabs v-model="sharedTab" density="compact" color="secondary">
              <v-tab value="accepted">Accepted</v-tab>
              <v-tab value="pending">Pending</v-tab>
            </v-tabs>
          </div>
          <v-tabs-window v-model="sharedTab">
            <v-tabs-window-item value="accepted">
              <v-list v-if="sharedAcceptedCircles.length > 0">
                <v-list-item
                  v-for="circle in sharedAcceptedCircles"
                  :key="circle.name"
                  :title="circle.title"
                  :subtitle="circle.name"
                  :to="{ name: 'circle', params: { circleId: circle.name } }"
                />
              </v-list>
              <div v-else>No accepted shared circles found.</div>
            </v-tabs-window-item>
            <v-tabs-window-item value="pending">
              <v-list v-if="sharedPendingCircles.length > 0">
                <v-list-item
                  v-for="circle in sharedPendingCircles"
                  :key="circle.name"
                  :title="circle.title"
                  :subtitle="circle.name"
                  :to="{ name: 'circle', params: { circleId: circle.name } }"
                />
              </v-list>
              <div v-else>No pending shared circles found.</div>
            </v-tabs-window-item>
          </v-tabs-window>
        </v-tabs-window-item>
        <!-- Explore Circles Tab -->
        <v-tabs-window-item value="explore">
          <div class="d-flex justify-space-between align-center mb-4">
            <h2>Explore Public Circles</h2>
          </div>
          <v-list v-if="publicCircles.length > 0">
            <v-list-item
              v-for="circle in publicCircles"
              :key="circle.name"
              :title="circle.title"
              :subtitle="circle.name"
              :to="{ name: 'circle', params: { circleId: circle.name } }"
            />
          </v-list>
          <div v-else>No public circles found.</div>
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card-text>
    <!-- Floating Action Button - only visible on My Circles tab -->
    <v-btn
      v-if="activeTab === 'my'"
      color="primary"
      icon="mdi-plus"
      style="position: fixed; bottom: 16px; right: 16px"
      :to="{ name: 'circleCreate' }"
    ></v-btn>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useCirclesStore } from '@/stores/circles'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'
const activeTab = ref('my')
const sharedTab = ref('accepted')
const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()

// Separate circle arrays for each view
const myCircles = ref<Circle[]>([])
const sharedAcceptedCircles = ref<Circle[]>([])
const sharedPendingCircles = ref<Circle[]>([])
const publicCircles = ref<Circle[]>([])

async function loadTab(tab: string) {
  switch (tab) {
    case 'my':
      await circlesStore.loadMyCircles()
      myCircles.value = [...circlesStore.circles]
      break
    case 'shared':
      await circlesStore.loadSharedCircles(200)
      sharedAcceptedCircles.value = [...circlesStore.circles]
      await circlesStore.loadSharedCircles(100)
      sharedPendingCircles.value = [...circlesStore.circles]
      break
    case 'explore':
      await circlesStore.loadPublicCircles()
      publicCircles.value = [...circlesStore.circles]
      break
  }
}

onMounted(async () => {
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'circles' } },
  ])
  await loadTab(activeTab.value)
})

watch(activeTab, async (newTab) => {
  await loadTab(newTab)
})
</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}
</style>
