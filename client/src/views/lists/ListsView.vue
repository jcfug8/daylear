<template>
  <div>
    <ListTabsPage
      ref="tabsPage"
      :tabs="tabs"
    >
    <template #filter>
      <v-row class="align-center" style="max-width: 600px; margin: 0 auto;">
        <!-- Left: Search -->
        <v-col cols="3" class="pa-0 d-flex align-center justify-start">
          <template v-if="searchExpanded || searchQuery">
            <v-text-field
              v-model="searchQuery"
              label="Search lists"
              prepend-inner-icon="mdi-magnify"
              clearable
              hide-details
              density="compact"
              class="mt-1 search-bar"
              :class="{ expanded: searchExpanded || searchQuery, collapsed: !searchExpanded && !searchQuery }"
              :style="searchBarStyle"
              @focus="onSearchFocus"
              @blur="onSearchBlur"
              @keydown.enter="onSearchEnter"
              ref="searchInput"
            />
          </template>
          <template v-else>
            <v-btn density="compact" icon variant="text" class="search-icon-btn" @click="expandSearch">
              <v-icon>mdi-magnify</v-icon>
            </v-btn>
          </template>
        </v-col>
        <!-- Center: Active account title/name -->
        <v-col cols="6" class="pa-0 d-flex align-center justify-center">
          <div
            class="active-account-title clickable text-center"
            style="cursor: pointer; user-select: none; font-weight: 500; font-size: 1.1rem; display: flex; align-items: center; justify-content: center; gap: 4px;"
            @click="showFilterModal = true"
          >
            <v-icon size="18">{{ selectedAccount?.icon || 'mdi-account-circle' }}</v-icon>
            <span class="account-title-ellipsis">{{ selectedAccount?.label || 'My Lists' }}</span>
          </div>
        </v-col>
        <!-- Right: Filter button -->
        <v-col cols="3" class="pa-0 d-flex align-center justify-end">
          <v-btn class="filter-button mr-2" :color="selectedVisibility.length === 0 ? 'white' : 'grey'" variant="flat" @click="showFilterModal = true" title="Filter lists">
            <v-icon>mdi-filter-variant</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </template>
      <template #my="{ items, loading }">
        <ListGrid :lists="getFilteredLists(items as List[])" :loading="loading as boolean" />
      </template>
      <template #pending="{ items, loading }">
        <ListGrid :lists="getFilteredLists(items as List[])" :loading="loading as boolean" @accept="onAcceptList" @decline="onDeclineList" />
      </template>
      <template #explore="{ items, loading }">
        <div class="d-flex justify-space-between align-center mb-4">
        </div>
        <ListGrid :lists="getFilteredLists(items as List[])" :loading="loading as boolean" />
      </template>
      <template #fab>
        <v-btn
          v-if="selectedAccount?.value === authStore.user.name"
          color="primary"
          density="compact"
          style="position: fixed; bottom: 46px; right: 16px"
          :to="{ name: 'listCreate' }"
        >
          <v-icon>mdi-plus</v-icon>
          <span>Create List</span>
        </v-btn>
      </template>
    </ListTabsPage>
  </div>
  <v-dialog v-model="showFilterModal" max-width="500">
    <v-card>
      <v-card-title>Filter Lists</v-card-title>
      <v-card-text>
        <!-- User/Circle select at the top -->
        <div class="mb-4">
          <div class="font-weight-bold mb-2">Account</div>
          <v-autocomplete
            ref="accountSelectRef"
            v-model="selectedAccount"
            :items="accountOptions"
            item-title="label"
            item-value="value"
            return-object
            hide-details
            density="compact"
            class="mb-4"
            :prepend-inner-icon="selectedAccount?.icon"
            :menu-props="{ maxHeight: '300px' }"
          >
            <template #item="{ props, item }">
              <v-list-item v-bind="props">
                <template #prepend>
                  <v-icon :icon="item.raw.icon" size="small" class="mr-2"></v-icon>
                </template>
              </v-list-item>
            </template>
          </v-autocomplete>
        </div>
        <div>
          <div class="font-weight-bold mb-2">Visibility</div>
          <v-autocomplete
            v-model="selectedVisibility"
            :items="allVisibility"
            label="Select visibility"
            multiple
            chips
            clearable
            hide-details
            density="compact"
          />
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text @click="clearFilters">Clear</v-btn>
        <v-btn color="primary" @click="showFilterModal = false">Ok</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch, onMounted } from 'vue'
import { useListStore } from '@/stores/list'
import { useAuthStore } from '@/stores/auth'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import ListGrid from '@/components/ListGrid.vue'
import type { List } from '@/genapi/api/lists/list/v1alpha1'
import { useAlertStore } from '@/stores/alerts'
import Fuse from 'fuse.js'
import { useCirclesStore } from '@/stores/circles'
import { useUsersStore } from '@/stores/users'

const circlesStore = useCirclesStore()
const authStore = useAuthStore()
const listStore = useListStore()
const alertStore = useAlertStore()
const usersStore = useUsersStore()

const acceptingListId = ref<string | null>(null)
const tabsPage = ref()

const searchQuery = ref('')
const searchExpanded = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)

// Dropdown for account/circle selection
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const selectedAccount = ref<any>(null)

function onSearchEnter() {
  searchExpanded.value = false
}

onMounted(() => {
  if (authStore.user && authStore.user.name) {
    circlesStore.loadMyCircles(authStore.user.name)
    usersStore.loadFriends(authStore.user.name)
  }
})

const accountOptions = computed(() => {
  const options = []
  if (authStore.user && authStore.user.name) {
    options.push({
      label: 'My Lists',
      value: authStore.user.name,
      account: "",
      icon: 'mdi-account-circle',
    })
    if (Array.isArray(circlesStore.myCircles)) {
      for (const circle of circlesStore.myCircles) {
        options.push({
          label: circle.title || 'Untitled Circle',
          value: circle.name,
          account: circle,
          icon: 'mdi-account-group',
        })
      }
    }
    if (Array.isArray(usersStore.friends)) {
      for (const user of usersStore.friends) {
        let label = ''
        if (user.givenName || user.familyName) { // user full name
          label = user.givenName + ' ' + user.familyName
          label = label.trim()
        } else if (user.username) { // user username
          label = user.username
        } 
        options.push({
          label: label,
          value: user.name,
          account: user,
          icon: 'mdi-account-circle',
        })
      }
    }
  }
  return options
})

// Set default selectedAccount to user on mount or when user changes
watch(
  () => authStore.user?.name,
  (newUserName) => {
    if (newUserName && (!selectedAccount.value || selectedAccount.value.value !== newUserName)) {
      selectedAccount.value = accountOptions.value[0]
    }
  },
  { immediate: true }
)

// When selectedAccount changes, reload the current tab
watch(
  selectedAccount,
  () => {
    // Only reload if tabsPage is ready
    nextTick(() => {
      tabsPage.value?.reloadActiveTab()
    })
  }
)

function expandSearch() {
  searchExpanded.value = true
  nextTick(() => {
    if (searchInput.value && searchInput.value.focus) {
      searchInput.value.focus()
    }
  })
}
function onSearchFocus() {
  searchExpanded.value = true
}
function onSearchBlur() {
  if (!searchQuery.value) {
    searchExpanded.value = false
  }
}

const searchBarStyle = computed(() => {
  return searchExpanded.value || searchQuery.value
    ? { maxWidth: '350px', width: '100%', transition: 'max-width 0.3s cubic-bezier(0.4,0,0.2,1)' }
    : { maxWidth: '44px', width: '44px', transition: 'max-width 0.3s cubic-bezier(0.4,0,0.2,1)' }
})

const showFilterModal = ref(false)
const selectedVisibility = ref<string[]>([])

function clearFilters() {
  selectedVisibility.value = []
}

const allVisibility = computed(() => {
  // Gather all unique visibility levels from all loaded lists in all tabs
  const lists = [
    ...listStore.myLists,
    ...listStore.sharedAcceptedLists,
    ...listStore.sharedPendingLists,
    ...listStore.publicLists,
  ]
  const set = new Set<string>()
  for (const list of lists) {
    if (list.visibility) {
      set.add(list.visibility)
    }
  }
  return Array.from(set).sort()
})

function getFilteredLists(items: List[]) {
  let filtered = items
  // Fuzzy search
  if (searchQuery.value) {
    const fuse = new Fuse(filtered, { keys: ['title'], threshold: 0.4 })
    filtered = fuse.search(searchQuery.value).map(result => result.item)
  }
  // Visibility filter
  if (selectedVisibility.value.length > 0) {
    filtered = filtered.filter(list =>
      list.visibility && selectedVisibility.value.includes(list.visibility)
    )
  }
  return filtered
}

const tabs = [
  {
    label: 'My Lists',
    value: 'my',
    icon: 'mdi-format-list-bulleted',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      // Use selectedAccount for context
      const account = selectedAccount.value?.account
      await listStore.loadMyLists(account.name)
      return [...listStore.myLists]
    },
  },
  {
    label: 'Pending',
    value: 'pending',
    icon: 'mdi-email-arrow-left-outline',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      const account = selectedAccount.value?.account
      await listStore.loadPendingLists(account.name)
      return [...listStore.sharedPendingLists]
    },
  },
  {
    label: 'Explore',
    value: 'explore',
    icon: 'mdi-card-search-outline',
    loader: async () => {
      if (!authStore.user || !authStore.user.name) throw new Error('User not authenticated')
      const account = selectedAccount.value?.account
      await listStore.loadPublicLists(account.name)
      return [...listStore.publicLists]
    },
  },
]

async function onAcceptList(list: List) {
  if (!list.listAccess?.name) return
  acceptingListId.value = list.listAccess.name
  try {
    await listStore.acceptList(list.listAccess.name)
    tabsPage.value?.reloadTab('pending')
  } catch (err) {
    alertStore.addAlert(err instanceof Error ? "Unable to accept list\n" + err.message : String(err), 'error')
  } finally {
    acceptingListId.value = null
  }
}

async function onDeclineList(list: List) {
  if (!list.listAccess?.name) return
  try {
    await listStore.deleteListAccess(list.listAccess.name)
    // Reload only the pending subtab
    tabsPage.value?.reloadTab('pending')
  } catch (err) {
    alertStore.addAlert(err instanceof Error ? "Unable to decline list\n" + err.message : String(err), 'error')
  }
}

</script>

<style scoped>
.v-tabs {
  margin-bottom: 24px;
}

.search-bar {
  transition: max-width 0.3s cubic-bezier(0.4,0,0.2,1), width 0.3s cubic-bezier(0.4,0,0.2,1);
  min-width: 44px;
}
.search-bar.collapsed {
  max-width: 44px !important;
  width: 44px !important;
  padding-left: 0 !important;
}
.search-bar.expanded {
  max-width: 350px !important;
  width: 100% !important;
}
.search-icon-btn {
  min-width: 44px;
  width: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.account-title-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
  display: inline-block;
  vertical-align: middle;
}
</style>