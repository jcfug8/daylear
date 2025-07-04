<template>
  <v-container>
    <v-tabs v-model="activeTab" align-tabs="center" color="primary" grow>
      <v-tab v-for="tab in tabs" :key="tab.value" :value="tab.value">
        {{ tab.label }}
      </v-tab>
    </v-tabs>
    <v-card-text>
      <v-tabs-window v-model="activeTab">
        <v-tabs-window-item v-for="tab in tabs" :key="tab.value" :value="tab.value">
          <template v-if="tab.subTabs">
            <v-tabs v-model="subTab[tab.value]" density="compact" color="secondary">
              <v-tab v-for="sub in tab.subTabs" :key="sub.value" :value="sub.value">{{ sub.label }}</v-tab>
            </v-tabs>
            <v-tabs-window v-model="subTab[tab.value]">
              <v-tabs-window-item v-for="sub in tab.subTabs" :key="sub.value" :value="sub.value">
                <slot :name="`${tab.value}-${sub.value}`" :items="items[tab.value]?.[sub.value] || []" :loading="loading[tab.value]?.[sub.value] || false" />
              </v-tabs-window-item>
            </v-tabs-window>
          </template>
          <template v-else>
            <slot :name="tab.value" :items="items[tab.value] || []" :loading="loading[tab.value] || false" />
          </template>
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card-text>
    <slot name="fab" :tab="activeTab" />
  </v-container>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'

interface TabDef {
  label: string
  value: string
  loader?: () => Promise<any>
  subTabs?: Array<{ label: string; value: string; loader?: () => Promise<any> }>
}

const props = defineProps<{
  tabs: TabDef[]
  initialTab?: string
}>()

const activeTab = ref(props.initialTab ?? props.tabs[0]?.value ?? '')
const subTab = ref<Record<string, string>>({})
const items = ref<Record<string, any>>({})
const loading = ref<Record<string, any>>({})

function loadTab(tabValue: string) {
  const tab = props.tabs.find(t => t.value === tabValue)
  if (!tab) return
  if (tab.subTabs) {
    if (!subTab.value[tabValue]) subTab.value[tabValue] = tab.subTabs[0].value
    for (const sub of tab.subTabs) {
      if (sub.loader) {
        loading.value[tabValue] = loading.value[tabValue] || {}
        loading.value[tabValue][sub.value] = true
        sub.loader().then(data => {
          items.value[tabValue] = items.value[tabValue] || {}
          items.value[tabValue][sub.value] = data
        }).finally(() => {
          loading.value[tabValue][sub.value] = false
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
  loadTab(activeTab.value)
})

watch(activeTab, (newTab) => {
  loadTab(newTab)
})
watch(subTab, (newSubTabs) => {
  const tab = props.tabs.find(t => t.value === activeTab.value)
  if (tab && tab.subTabs) {
    const subValue = newSubTabs[activeTab.value]
    const sub = tab.subTabs.find(s => s.value === subValue)
    if (sub && sub.loader) {
      loading.value[activeTab.value] = loading.value[activeTab.value] || {}
      loading.value[activeTab.value][subValue] = true
      sub.loader().then(data => {
        items.value[activeTab.value] = items.value[activeTab.value] || {}
        items.value[activeTab.value][subValue] = data
      }).finally(() => {
        loading.value[activeTab.value][subValue] = false
      })
    }
  }
})
</script> 