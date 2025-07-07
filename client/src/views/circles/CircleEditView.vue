<template>
  <circle-form
    v-if="circlesStore.circle"
    v-model="editedCircle"
    :is-editing="true"
    @save="saveSettings"
    @close="navigateBack"
    @imageSelected="handleImageSelected"
  />
</template>

<script setup lang="ts">
import { useCirclesStore } from '@/stores/circles'
import { storeToRefs } from 'pinia'
import { onMounted, ref, watch } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRouter, useRoute } from 'vue-router'
import type { Circle, apitypes_VisibilityLevel, apitypes_PermissionLevel, apitypes_AccessState } from '@/genapi/api/circles/circle/v1alpha1'
import { useAuthStore } from '@/stores/auth'
import CircleForm from '@/views/circles/forms/CircleForm.vue'
import { fileService } from '@/api/api'

const router = useRouter()
const route = useRoute()
const circlesStore = useCirclesStore()
const { circle } = storeToRefs(circlesStore)
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()

const editedCircle = ref<Circle>({
  name: '',
  title: '',
  imageUri: '',
  visibility: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
  circleAccess: undefined,
})

const pendingImageFile = ref<File | null>(null)

function navigateBack() {
  router.push({ name: 'circle', params: { circleId: editedCircle.value.name } })
}

async function saveSettings() {
  try {
    circlesStore.circle = editedCircle.value
    await circlesStore.updateCircle()
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
  } catch (error) {
    console.error('Error saving settings:', error)
    alert('Failed to save settings')
  }
}

function handleImageSelected(file: File | null, url: string | null) {
  pendingImageFile.value = file
  if (url) {
    editedCircle.value.imageUri = url
  }
}

onMounted(async () => {
  await circlesStore.loadCircle(route.params.circleId as string)
  if (circle.value) {
    editedCircle.value = { ...circle.value }
  }
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'circles' } },
    { title: circle.value?.title || 'Circle', to: { name: 'circle', params: { circleId: circle.value?.name } } },
    { title: 'Edit', to: { name: 'circle-edit', params: { circleId: circle.value?.name } } },
  ])
})

watch(circle, (newVal) => {
  if (newVal) {
    editedCircle.value = { ...newVal }
  }
})
</script>

<style></style>
