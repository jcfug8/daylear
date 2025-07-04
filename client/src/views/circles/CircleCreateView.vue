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

const router = useRouter()
const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()

const pendingImageFile = ref<File | null>(null)

function navigateBack() {
  router.push({ name: 'circles' })
}

async function saveCircle() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  try {
    await circlesStore.createCircle(authStore.user.name)
    // Upload image if there's a pending file
    if (pendingImageFile.value && circlesStore.circle?.name) {
      const response = await fileService.UploadCircleImage({
        name: circlesStore.circle.name,
        file: pendingImageFile.value,
      })
      circlesStore.circle.imageUri = response.imageUri
    }
    authStore.loadAuthCircles()
    navigateBack()
  } catch (err) {
    alert(err instanceof Error ? err.message : String(err))
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
    imageUri: '',
    visibility: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
    permission: 'PERMISSION_LEVEL_UNSPECIFIED' as apitypes_PermissionLevel,
    state: 'ACCESS_STATE_UNSPECIFIED' as apitypes_AccessState,
  }
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'circles' } },
    { title: 'Create New Circle', to: { name: 'circleCreate' } },
  ])
})
</script>

<style scoped>
</style>
