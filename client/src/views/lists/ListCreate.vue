<template>
  <ListForm
    :is-edit-mode="false"
    :saving="saving"
    @save="saveList"
    @cancel="navigateBack"
  />
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useListStore } from '@/stores/list'
import { useAuthStore } from '@/stores/auth'
import { useAlertStore } from '@/stores/alerts'
import ListForm from '@/views/lists/ListForm.vue'

const router = useRouter()
const route = useRoute()
const listStore = useListStore()
const authStore = useAuthStore()
const alertStore = useAlertStore()

const saving = ref(false)

const circleName = computed(() => {
  return route.path.replace('/lists/create', '').slice(1)
})

function navigateBack() {
  if (route.params.circleId) {
    router.push({ name: 'circle', params: { circleId: route.params.circleId } })
  } else {
    router.push({ name: 'lists' })
  }
}

async function saveList() {
  if (!authStore.user?.name && !route.params.circleId) {
    throw new Error('User not authenticated')
  }
  
  saving.value = true
  try {
    const list = await listStore.createList(circleName.value ? circleName.value : authStore.user!.name!)
    router.push('/' + list.name!)
  } catch (err) {
    alertStore.addAlert(err instanceof Error ? "Unable to create list\n" + err.message : String(err), 'error')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  // Initialize an empty list
  listStore.initEmptyList()
})
</script>

