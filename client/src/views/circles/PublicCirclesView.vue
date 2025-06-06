<template>
  <div class="pa-4">
    <v-text-field
      v-model="search"
      label="Search Circles"
      prepend-inner-icon="mdi-magnify"
      @input="onSearch"
      clearable
    />
    <v-list v-if="publicCircles &&publicCircles.length > 0">
      <v-list-item
        v-for="circle in publicCircles"
        :key="circle.name"
        :title="circle.title"
        :subtitle="circle.name"
        :to="{ name: 'publicCircle', params: { circleId: circle.name } }"
      />
    </v-list>
    <div v-else-if="search">No circles found.</div>
  </div>
  <v-btn
    color="primary"
    icon="mdi-plus"
    style="position: fixed; bottom: 16px; right: 16px"
    :to="{ name: 'circleCreate' }"
  ></v-btn>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useCirclesStore } from '@/stores/circles'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { storeToRefs } from 'pinia'
const search = ref('')
const circlesStore = useCirclesStore()
const { publicCircles } = storeToRefs(circlesStore)
const breadcrumbStore = useBreadcrumbStore()

onMounted(() => {
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'publicCircles' } },
  ])
})

let debounceTimer: ReturnType<typeof setTimeout> | null = null

function onSearch() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    if (search.value) {
      circlesStore.loadPublicCircles(`title = '${search.value}'`)
    } else {
      publicCircles.value.length = 0 // Clear results if search is cleared
    }
  }, 1000)
}
</script>

<style scoped>
.pa-4 {
  max-width: 500px;
  margin: 0 auto;
}
</style>
