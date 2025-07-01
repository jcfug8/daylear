<template>
  <div class="pa-4">
    <v-list v-if="circles &&circles.length > 0">
      <v-list-item
        v-for="circle in circles"
        :key="circle.name"
        :title="circle.title"
        :subtitle="circle.name"
        :to="{ name: 'circle', params: { circleId: circle.name } }"
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
const { circles } = storeToRefs(circlesStore)
const breadcrumbStore = useBreadcrumbStore()

onMounted(() => {
  breadcrumbStore.setBreadcrumbs([
    { title: 'Circles', to: { name: 'circles' } },
  ])
  circlesStore.loadCircles()
})

</script>

<style scoped>
.pa-4 {
  max-width: 500px;
  margin: 0 auto;
}
</style>
