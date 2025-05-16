<template>
  <v-card v-if="circle" class="mx-auto my-12" max-width="400">
    <v-card-title>Create New Circle</v-card-title>
    <v-card-text>
      <v-text-field
        v-model="circle.title"
        label="Circle Title"
        required
      />
      <v-checkbox v-model="circle.isPublic" label="Public" />
    </v-card-text>
    <v-card-actions>
      <v-btn color="primary" @click="saveCircle">Save</v-btn>
      <v-btn text @click="navigateBack">Cancel</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { useCirclesStore } from '@/stores/circles'
import { useRouter } from 'vue-router'
import { onMounted, computed } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const circlesStore = useCirclesStore()
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()

const { circle } = storeToRefs(circlesStore)

function navigateBack() {
  router.push({ name: 'publicCircles' })
}

function saveCircle() {
  if (!authStore.user || !authStore.user.name) {
    throw new Error('User not authenticated')
  }
  circlesStore
    .createCircle(authStore.user.name)
    .then(() => {
      authStore.loadAuthCircles()
      navigateBack()
    })
    .catch((err) => alert(err.message || err))
}

onMounted(() => {
  circlesStore.initEmptyCircle()
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'publicCircles' } },
    { title: 'Create New Circle', to: { name: 'circleCreate' } },
  ])
})
</script>

<style scoped>
.mx-auto {
  margin-left: auto;
  margin-right: auto;
}
.my-12 {
  margin-top: 48px;
  margin-bottom: 48px;
}
</style>
