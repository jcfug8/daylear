<template>
  <div class="pa-4">
    <div v-if="loading">Loading...</div>
    <div v-else-if="!circle">Circle not found.</div>
    <div v-else>
      <h1>{{ circle.title }}</h1>
      <p><strong>ID:</strong> {{ circle.name }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useCirclesStore } from '@/stores/circles'
import { storeToRefs } from 'pinia'

const route = useRoute()
const breadcrumbStore = useBreadcrumbStore()
const circlesStore = useCirclesStore()
const { publicCircle } = storeToRefs(circlesStore)
const loading = ref(true)

async function fetchCircle() {
  loading.value = true
  let circleId = route.params.circleId
  if (Array.isArray(circleId)) circleId = circleId[0]
  await circlesStore.loadPublicCircle(circleId as string)
  setCrumbs()
  loading.value = false
}

function setCrumbs() {
  let circleId = route.params.circleId
  if (Array.isArray(circleId)) circleId = circleId[0]
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'publicCircles' } },
    { title: publicCircle.value?.title || circleId, to: { name: 'publicCircle', params: { circleId } } },
  ])
}

onMounted(fetchCircle)
watch(() => route.params.circleId, fetchCircle)

const circle = publicCircle
</script>

<style scoped>
.pa-4 {
  max-width: 500px;
  margin: 0 auto;
}
</style>
