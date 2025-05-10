<template>Circle {{ $route.params.circleId }}</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useCirclesStore } from '@/stores/circles'

const route = useRoute()
const breadcrumbStore = useBreadcrumbStore()
const circlesStore = useCirclesStore()

function setCrumbs() {
  let circleId = route.params.circleId
  if (Array.isArray(circleId)) circleId = circleId[0]
  const circle = circlesStore.circles.find(c => c.name === circleId)
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'circles' } },
    { title: circle?.title || circleId, to: { name: 'circle', params: { circleId } } },
  ])
}

onMounted(() => {
  setCrumbs()
})

watch(() => route.params.circleId, setCrumbs)
</script>

<style></style>
