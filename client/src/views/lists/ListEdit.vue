<template>
  <ListForm
    :is-edit-mode="true"
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

const trimmedListName = computed(() => {
  return route.path.substring(route.path.indexOf('lists/')).replace('/edit', '')
})

function navigateBack() {
  if (route.params.circleId) {
    router.push({ name: 'circle', params: { circleId: route.params.circleId } })
  } else {
    router.push('/' + listStore.list?.name || '/lists')
  }
}

async function saveList() {
  if (!authStore.user?.name && !route.params.circleId) {
    throw new Error('User not authenticated')
  }
  
  if (!listStore.list?.name) {
    throw new Error('List not found')
  }
  
  saving.value = true
  try {
    const list = await listStore.updateList()
    router.push('/' + list.name!)
  } catch (err) {
    alertStore.addAlert(err instanceof Error ? "Unable to update list\n" + err.message : String(err), 'error')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  // Load the existing list
  await listStore.loadList(trimmedListName.value)
})
</script>
