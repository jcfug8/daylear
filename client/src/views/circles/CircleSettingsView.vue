<template>
  <v-container>
    <v-card class="mx-auto" max-width="600">
      <v-card-title>Circle Settings</v-card-title>
      <v-card-text>
        <v-list>
          <v-list-item>
            <template v-slot:prepend>
              <v-icon icon="mdi-circle"></v-icon>
            </template>
            <v-list-item-title>Title</v-list-item-title>
            <v-list-item-subtitle>{{ circle?.title || 'Not set' }}</v-list-item-subtitle>
          </v-list-item>

          <v-list-item>
            <template v-slot:prepend>
              <v-icon icon="mdi-eye"></v-icon>
            </template>
            <v-list-item-title>Public</v-list-item-title>
            <v-list-item-subtitle>{{ circle?.isPublic ? 'Yes' : 'No' }}</v-list-item-subtitle>
          </v-list-item>
        </v-list>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          :to="{ name: 'circle-settings-edit', params: { circleId: $route.params.circleId } }"
        >
          Edit Settings
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { useCirclesStore } from '@/stores/circles'
import { storeToRefs } from 'pinia'
import { onMounted, watch } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRoute } from 'vue-router'

const route = useRoute()
const circlesStore = useCirclesStore()
const { circle } = storeToRefs(circlesStore)
const breadcrumbStore = useBreadcrumbStore()

async function loadAndSetBreadcrumbs(circleId: string | string[]) {
  await circlesStore.loadCircle(circleId as string)
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circle Settings', to: { name: 'circle-settings', params: { circleId: circle.value?.name } } },
  ])
}

onMounted(() => {
  loadAndSetBreadcrumbs(route.params.circleId)
})

watch(
  () => route.params.circleId,
  (newId) => {
    loadAndSetBreadcrumbs(newId)
  }
)
</script>

<style></style>
