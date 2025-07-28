<template>
  <circle-form
    v-if="circlesStore.circle"
    v-model="circlesStore.circle"
    :is-editing="false"
    @save="saveCircle"
    @close="navigateBack"
    @imageSelected="handleImageSelected"
  />
</template>

<script setup lang="ts">
import { useCirclesStore } from '@/stores/circles'
import { useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useAuthStore } from '@/stores/auth'
import CircleForm from '@/views/circles/forms/CircleForm.vue'
import { fileService } from '@/api/api'
import type { Circle, apitypes_VisibilityLevel, apitypes_PermissionLevel, apitypes_AccessState } from '@/genapi/api/circles/circle/v1alpha1'
import { useAlertStore } from '@/stores/alerts'

const router = useRouter()
const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()
const alertsStore = useAlertStore()

const pendingImageFile = ref<File | null>(null)

function navigateBack() {
  router.push({ name: 'circles' })
}

async function saveCircle() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  try {
    const circle = await circlesStore.createCircle(authStore.user.name)
    // Upload image if there's a pending file
    if (pendingImageFile.value && circlesStore.circle?.name) {
      const response = await fileService.UploadCircleImage({
        name: circlesStore.circle.name,
        file: pendingImageFile.value,
      })
      circlesStore.circle.imageUri = response.imageUri
    }
    router.push('/'+circle.name)
  } catch (err) {
    alertsStore.addAlert({
      message: err instanceof Error ? err.message : String(err),
      type: 'error',
    })
  }
}

function handleImageSelected(file: File | null, url: string | null) {
  pendingImageFile.value = file
  if (url) {
    if (circlesStore.circle) circlesStore.circle.imageUri = url
  }
}

onMounted(() => {
  circlesStore.circle = {
    name: '',
    title: '',
    handle: '',
    imageUri: '',
    visibility: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
    circleAccess: undefined,
    description: '',
  }
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'circles' } },
    { title: 'Create New Circle', to: { name: 'circleCreate' } },
  ])
})
</script>

<style scoped>
</style>
