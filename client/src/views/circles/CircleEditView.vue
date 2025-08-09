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
import { useRouter, useRoute } from 'vue-router'
import type { Circle, apitypes_VisibilityLevel } from '@/genapi/api/circles/circle/v1alpha1'
import CircleForm from '@/views/circles/forms/CircleForm.vue'
import { fileService } from '@/api/api'
import { computed } from 'vue'
import { useAlertStore } from '@/stores/alerts'

const router = useRouter()
const route = useRoute()
const circlesStore = useCirclesStore()
const { circle } = storeToRefs(circlesStore)
const alertsStore = useAlertStore()

const editedCircle = ref<Circle>({
  name: '',
  title: '',
  handle: '',
  imageUri: '',
  visibility: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
  circleAccess: undefined,
  description: '',
})

const pendingImageFile = ref<File | null>(null)

function navigateBack() {
  router.push('/'+circleName.value)
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
    navigateBack()
  } catch (err) {
    alertsStore.addAlert(err instanceof Error ? err.message : String(err),'error')
  }
}

function handleImageSelected(file: File | null, url: string | null) {
  pendingImageFile.value = file
  if (url) {
    editedCircle.value.imageUri = url
  }
}

const circleName = computed(() => {
  return route.path.replace('/edit', '').substring(1)
})

async function loadCircle() {
  await circlesStore.loadCircle(circleName.value)
  if (circle.value) {
    editedCircle.value = { ...circle.value }
  }
}

onMounted(async () => {
  await loadCircle()
})

watch(circle, (newVal) => {
  if (newVal) {
    editedCircle.value = { ...newVal }
  }
})

watch(
  () => route.path,
  async (newCircleName) => {
    if (newCircleName) {
      await loadCircle()
    }
  }
)
</script>

<style></style>
