<template>
  <v-form ref="form" v-model="valid" @submit.prevent="saveList">
    <v-container v-if="listStore.list" max-width="600" class="pb-16">
      <v-row>
        <v-col cols="12">
          <v-text-field
            v-model="listStore.list.title"
            label="List Title"
            :rules="titleRules"
            required
            variant="outlined"
            prepend-inner-icon="mdi-format-title"
            counter="100"
            maxlength="100"
          />
        </v-col>
        
        <v-col cols="12">
          <v-textarea
            v-model="listStore.list.description"
            label="Description"
            variant="outlined"
            prepend-inner-icon="mdi-text"
            rows="3"
            counter="500"
            maxlength="500"
            hint="Optional description for your list"
            persistent-hint
          />
        </v-col>
        
        <v-col cols="12">
          <v-checkbox
            v-model="listStore.list.showCompleted"
            label="Show completed items"
            color="primary"
            hint="When enabled, completed items will be visible in the list"
            persistent-hint
          />
        </v-col>
        
        <v-col cols="12">
          <v-select
            v-model="listStore.list.visibility"
            :items="visibilityOptions"
            label="Visibility"
            variant="outlined"
            prepend-inner-icon="mdi-eye"
            hint="Control who can see this list"
            persistent-hint
          />
        </v-col>
      </v-row>
    </v-container>

    <!-- Close FAB -->
    <v-btn
      color="error"
      density="compact"
      style="position: fixed; bottom: 16px; left: 16px"
      @click="navigateBack"
    >
      <v-icon>mdi-close</v-icon>
      Cancel
    </v-btn>

    <!-- Save FAB -->
    <v-btn
      color="success"
      density="compact"
      style="position: fixed; bottom: 16px; right: 16px"
      :loading="saving"
      @click="saveList"
    >
      <v-icon>mdi-content-save</v-icon>
      Save Changes
    </v-btn>
  </v-form>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useListStore } from '@/stores/list'
import { useAuthStore } from '@/stores/auth'
import { useAlertStore } from '@/stores/alerts'

const router = useRouter()
const route = useRoute()
const listStore = useListStore()
const authStore = useAuthStore()
const alertStore = useAlertStore()

const form = ref()
const valid = ref(false)
const saving = ref(false)

// Form validation rules
const titleRules = [
  (v: string) => !!v || 'Title is required',
  (v: string) => (v && v.length >= 3) || 'Title must be at least 3 characters',
  (v: string) => (v && v.length <= 100) || 'Title must be less than 100 characters',
]

// Visibility options
const visibilityOptions = [
  { title: 'Public', value: 'VISIBILITY_LEVEL_PUBLIC', subtitle: 'Anyone can see this list' },
  { title: 'Restricted', value: 'VISIBILITY_LEVEL_RESTRICTED', subtitle: 'Only people with access can see this list' },
  { title: 'Private', value: 'VISIBILITY_LEVEL_PRIVATE', subtitle: 'Only you can see this list' },
  { title: 'Hidden', value: 'VISIBILITY_LEVEL_HIDDEN', subtitle: 'List is hidden from everyone' },
]

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
