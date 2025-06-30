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

          <v-select
            v-model="editedCircle.visibility"
            :items="visibilityOptions"
            item-title="label"
            item-value="value"
            label="Visibility"
            required
          />
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
import type { Circle, apitypes_VisibilityLevel, apitypes_PermissionLevel } from '@/genapi/api/circles/circle/v1alpha1'
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
  visibility: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
  permission: 'PERMISSION_LEVEL_UNSPECIFIED' as apitypes_PermissionLevel,
})

const visibilityOptions = [
  { label: 'Public', value: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel },
  { label: 'Restricted', value: 'VISIBILITY_LEVEL_RESTRICTED' as apitypes_VisibilityLevel },
  { label: 'Private', value: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel },
  { label: 'Hidden', value: 'VISIBILITY_LEVEL_HIDDEN' as apitypes_VisibilityLevel },
]

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
