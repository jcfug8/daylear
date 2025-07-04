<template>
  <v-container>
    <v-tabs v-model="internalTab" align-tabs="center" color="primary" grow>
      <v-tab v-for="tab in tabs" :key="tab.value" :value="tab.value">{{ tab.label }}</v-tab>
    </v-tabs>
    <v-card-text>
      <v-tabs-window v-model="internalTab">
        <v-tabs-window-item v-for="tab in tabs" :key="tab.value" :value="tab.value">
          <slot
            :name="tab.value"
            :items="itemsMap[tab.value] || []"
            :loading="loadingMap[tab.value] || false"
            :tab="tab.value"
          />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card-text>
    <slot name="fab" :tab="internalTab" />
  </v-container>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'

interface TabDef {
  label: string
  value: string
}

const props = defineProps<{
  tabs: TabDef[]
  modelValue?: string
  itemsMap: Record<string, any[]>
  loadingMap: Record<string, boolean>
}>()

const emit = defineEmits(['update:modelValue'])

const internalTab = ref(props.modelValue ?? props.tabs[0]?.value ?? '')

watch(() => props.modelValue, (val) => {
  if (val !== undefined) internalTab.value = val
})

watch(internalTab, (val) => {
  emit('update:modelValue', val)
})
</script> 