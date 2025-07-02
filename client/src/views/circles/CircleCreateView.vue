<template>
  <circle-form
    v-if="circlesStore.circle"
    v-model="circlesStore.circle"
    :is-editing="false"
    @save="saveCircle"
    @close="navigateBack"
  />
</template>

<script setup lang="ts">
import { useCirclesStore } from '@/stores/circles'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useAuthStore } from '@/stores/auth'
import CircleForm from '@/views/circles/forms/CircleForm.vue'

const router = useRouter()
const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()

function navigateBack() {
  router.push({ name: 'circles' })
}

async function saveCircle() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  try {
    await circlesStore.createCircle(authStore.user.name)
    authStore.loadAuthCircles()
    navigateBack()
  } catch (err) {
    alert(err instanceof Error ? err.message : String(err))
  }
}

onMounted(() => {
  circlesStore.initEmptyCircle()
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'circles' } },
    { title: 'Create New Circle', to: { name: 'circleCreate' } },
  ])
})
</script>

<style scoped>
</style>
