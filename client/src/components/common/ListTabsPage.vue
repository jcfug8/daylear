<template>

  <v-system-bar elevation="0" height="35" color="white" v-if="slots.filter">
    <div style="width: 100%;">
      <slot name="filter" />
    </div>
  </v-system-bar>
  <v-system-bar elevation="0" color="white" height="40">
    <v-tabs center-active v-model="internalActiveTab" align-tabs="center" color="primary" grow>
      <v-tab min-width="50" density="compact" v-for="tab in tabs" :key="tab.value" :value="tab.value">
        <template v-if="!tab.disabled">
          <v-icon v-if="tab.icon" left>{{ tab.icon }}</v-icon>
          <span class="text-caption">{{ tab.label }}</span>
        </template>
      </v-tab>
    </v-tabs>
  </v-system-bar>
  <v-tabs-window v-model="internalActiveTab">
    <v-tabs-window-item v-for="tab in tabs" :key="tab.value" :value="tab.value">
      <template v-if="tab.subTabs">
        <v-tabs v-model="subTab[tab.value]" density="compact" color="secondary">
          <v-tab v-for="sub in tab.subTabs" :key="sub.value" :value="sub.value">{{ sub.label }}</v-tab>
        </v-tabs>
        <v-tabs-window v-model="subTab[tab.value]">
          <v-tabs-window-item v-for="sub in tab.subTabs" :key="sub.value" :value="sub.value">
            <slot :name="`${tab.value}-${sub.value}`" :items="getItems(tab.value, sub.value)"
              :loading="getLoading(tab.value, sub.value)" />
          </v-tabs-window-item>
        </v-tabs-window>
      </template>
      <template v-else>
        <slot :name="tab.value" :items="getItems(tab.value)" :loading="getLoading(tab.value)" />
      </template>
    </v-tabs-window-item>
  </v-tabs-window>
  <slot name="fab" :tab="internalActiveTab" />

</template>

<script setup lang="ts">
import { ref, watch, onMounted, defineExpose, useSlots } from 'vue'

interface TabDef {
  label: string | undefined
  value: string
  loader?: () => Promise<unknown>
  subTabs?: Array<{ label: string; value: string; loader?: () => Promise<unknown> }>
  icon?: string
  disabled?: boolean
}

const props = defineProps<{
  tabs: TabDef[]
  initialTab?: string
  activeTab?: string
}>()
const emit = defineEmits<{
  (e: 'update:activeTab', value: string): void
}>()

const internalActiveTab = ref(props.activeTab ?? props.initialTab ?? props.tabs[0]?.value ?? '')
const subTab = ref<Record<string, string>>({})
type ItemsMap = Record<string, unknown | Record<string, unknown>>
const items = ref<ItemsMap>({})
type LoadingMap = Record<string, boolean | Record<string, boolean>>
const loading = ref<LoadingMap>({})
const slots = useSlots()

// Helper functions to safely access nested properties with proper typing
function getItems(tabValue: string, subValue?: string): unknown[] {
  if (subValue) {
    const tabItems = items.value[tabValue]
    if (typeof tabItems === 'object' && tabItems !== null) {
      const subItems = (tabItems as Record<string, unknown>)[subValue]
      return Array.isArray(subItems) ? subItems : []
    }
    return []
  }
  const tabItems = items.value[tabValue]
  return Array.isArray(tabItems) ? tabItems : []
}

function getLoading(tabValue: string, subValue?: string): boolean {
  if (subValue) {
    const tabLoading = loading.value[tabValue]
    if (typeof tabLoading === 'object' && tabLoading !== null) {
      return (tabLoading as Record<string, boolean>)[subValue] || false
    }
    return false
  }
  const tabLoading = loading.value[tabValue]
  return typeof tabLoading === 'boolean' ? tabLoading : false
}

function loadTab(tabValue: string) {
  const tab = props.tabs.find(t => t.value === tabValue)
  if (!tab) return
  if (tab.subTabs) {
    if (!subTab.value[tabValue]) subTab.value[tabValue] = tab.subTabs[0].value
    for (const sub of tab.subTabs) {
      if (sub.loader) {
        if (typeof loading.value[tabValue] !== 'object') {
          loading.value[tabValue] = {}
        }
        ; (loading.value[tabValue] as Record<string, boolean>)[sub.value] = true
        sub.loader().then(data => {
          if (typeof items.value[tabValue] !== 'object') {
            items.value[tabValue] = {}
          }
          ; (items.value[tabValue] as Record<string, unknown>)[sub.value] = data
        }).finally(() => {
          ; (loading.value[tabValue] as Record<string, boolean>)[sub.value] = false
        })
      }
    }
  } else if (tab.loader) {
    loading.value[tabValue] = true
    tab.loader().then(data => {
      items.value[tabValue] = data
    }).finally(() => {
      loading.value[tabValue] = false
    })
  }
}

onMounted(() => {
  loadTab(internalActiveTab.value)
})

watch(internalActiveTab, (newTab) => {
  emit('update:activeTab', newTab)
  loadTab(newTab)
})
watch(() => props.activeTab, (newVal) => {
  if (typeof newVal === 'string' && newVal !== internalActiveTab.value) {
    internalActiveTab.value = newVal
  }
})
watch(subTab, (newSubTabs) => {
  const tab = props.tabs.find(t => t.value === internalActiveTab.value)
  if (tab && tab.subTabs) {
    const subValue = newSubTabs[internalActiveTab.value]
    const sub = tab.subTabs.find(s => s.value === subValue)
    if (sub && sub.loader) {
      if (typeof loading.value[internalActiveTab.value] !== 'object') {
        loading.value[internalActiveTab.value] = {}
      }
      ; (loading.value[internalActiveTab.value] as Record<string, boolean>)[subValue] = true
      sub.loader().then(data => {
        if (typeof items.value[internalActiveTab.value] !== 'object') {
          items.value[internalActiveTab.value] = {}
        }
        ; (items.value[internalActiveTab.value] as Record<string, unknown>)[subValue] = data
      }).finally(() => {
        ; (loading.value[internalActiveTab.value] as Record<string, boolean>)[subValue] = false
      })
    }
  }
})

// Expose a method to reload the current active tab and subtab
function reloadActiveTab() {
  const tab = props.tabs.find(t => t.value === internalActiveTab.value)
  if (!tab) return
  if (tab.subTabs) {
    const subValue = subTab.value[internalActiveTab.value] || tab.subTabs[0].value
    const sub = tab.subTabs.find(s => s.value === subValue)
    if (sub && sub.loader) {
      if (typeof loading.value[internalActiveTab.value] !== 'object') {
        loading.value[internalActiveTab.value] = {}
      }
      ; (loading.value[internalActiveTab.value] as Record<string, boolean>)[subValue] = true
      sub.loader().then(data => {
        if (typeof items.value[internalActiveTab.value] !== 'object') {
          items.value[internalActiveTab.value] = {}
        }
        ; (items.value[internalActiveTab.value] as Record<string, unknown>)[subValue] = data
      }).finally(() => {
        ; (loading.value[internalActiveTab.value] as Record<string, boolean>)[subValue] = false
      })
    }
  } else {
    loadTab(internalActiveTab.value)
  }
}

function reloadTab(tabValue: string, subTabValue?: string) {
  const tab = props.tabs.find(t => t.value === tabValue)
  if (!tab) return
  if (tab.subTabs && subTabValue) {
    const sub = tab.subTabs.find(s => s.value === subTabValue)
    if (sub && sub.loader) {
      if (typeof loading.value[tabValue] !== 'object') {
        loading.value[tabValue] = {}
      }
      ; (loading.value[tabValue] as Record<string, boolean>)[subTabValue] = true
      sub.loader().then(data => {
        if (typeof items.value[tabValue] !== 'object') {
          items.value[tabValue] = {}
        }
        ; (items.value[tabValue] as Record<string, unknown>)[subTabValue] = data
      }).finally(() => {
        ; (loading.value[tabValue] as Record<string, boolean>)[subTabValue] = false
      })
    }
  } else {
    loadTab(tabValue)
  }
}
defineExpose({ reloadActiveTab, reloadTab, activeTab: internalActiveTab })
</script>