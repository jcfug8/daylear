<template>
  <v-container>
    <v-card class="mx-auto" max-width="600">
      <v-card-title>Edit Circle Settings</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="saveSettings">
          <v-text-field
            v-model="editedCircle.title"
            label="Title"
            required
          ></v-text-field>

          <v-checkbox
            v-model="editedCircle.isPublic"
            label="Public"
          ></v-checkbox>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="secondary"
          @click="navigateBack"
        >
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          @click="saveSettings"
        >
          Save Changes
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { useCirclesStore } from '@/stores/circles'
import { storeToRefs } from 'pinia'
import { onMounted, ref, watch } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRouter, useRoute } from 'vue-router'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const circlesStore = useCirclesStore()
const { circle } = storeToRefs(circlesStore)
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()

const editedCircle = ref<Circle>({
  name: '',
  title: '',
  isPublic: false,
})

function navigateBack() {
  router.push({ name: 'circle-settings', params: { circleId: editedCircle.value.name } })
}

async function saveSettings() {
  try {
    circlesStore.circle = editedCircle.value
    await circlesStore.updateCircle()
    authStore.loadAuthCircles()
    navigateBack()
  } catch (error) {
    console.error('Error saving settings:', error)
    alert('Failed to save settings')
  }
}

onMounted(async () => {
  await circlesStore.loadCircle(route.params.circleId as string)
  if (circle.value) {
    editedCircle.value = { ...circle.value }
  }
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circle Settings', to: { name: 'circle-settings', params: { circleId: circle.value?.name } } },
    { title: 'Edit', to: { name: 'circle-settings-edit', params: { circleId: circle.value?.name } } },
  ])
})

watch(circle, (newVal) => {
  if (newVal) {
    editedCircle.value = { ...newVal }
  }
})
</script>

<style></style>
