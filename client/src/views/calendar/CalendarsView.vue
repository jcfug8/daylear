<template>
  <ListTabsPage
    ref="tabsPage"
    :tabs="tabs"
    v-model:activeTab="currentTab"
  >
    <template #filter>
      <v-row class="align-center" style="max-width: 900px; margin: 0 auto;">
        <!-- Left: Search (grid mode, or any non-my tab) -->
        <v-col cols="3" class="pa-0 d-flex align-center justify-start">
          <template v-if="(isMyTab && viewMode === 'grid') || !isMyTab">
            <template v-if="searchExpanded || searchQuery">
              <v-text-field
                v-model="searchQuery"
                label="Search calendars"
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
              <v-btn icon variant="text" class="search-icon-btn" @click="expandSearch">
                <v-icon>mdi-magnify</v-icon>
              </v-btn>
            </template>
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
            <span class="account-title-ellipsis">{{ selectedAccount?.label || 'My Calendars' }}</span>
          </div>
        </v-col>
        <!-- Right: Grid -> Filter button; Schedule -> Calendar menu -->
        <v-col cols="3" class="pa-0 d-flex align-center justify-end">
          <template v-if="(isMyTab && viewMode === 'schedule')">
            <v-menu v-model="settingsOpen" :close-on-content-click="false">
              <template #activator="{ props }">
                <v-btn density="compact" icon="mdi-cog" v-bind="props" />
              </template>
              <v-list>
                <v-list-item
                  v-for="calendar in calendarsStore.myCalendars" :key="calendar.name"
                  prepend-icon="mdi-calendar"
                  :title="calendar.title"
                />
              </v-list>
            </v-menu>
          </template>
          <template v-else>
            <v-btn class="filter-button mr-2" color="white" variant="flat" @click="showFilterModal = true" title="Filter calendars">
              <v-icon>mdi-filter-variant</v-icon>
            </v-btn>
          </template>
        </v-col>
      </v-row>
    </template>
    <template #my="{ items, loading }">
      <CalendarGrid v-if="viewMode === 'grid'" :calendars="getFilteredCalendars(items as Calendar[])" :loading="(loading as boolean)" />
      <template v-else>
        <ScheduleCal 
          v-if="!loading" 
          :events="events" :calendars="(items as Calendar[])" 
          :show-create-button="true" 
          @created="onEventCreated" 
          @updated="onEventUpdated"
          @deleted="onEventDeleted"
        />
      </template>
      <!-- View mode toggle FAB -->
      <v-btn
        color="primary"
        density="compact"
        class="text-none"
        style="position: fixed; bottom: 16px; left: 16px; z-index: 10;"
        @click="toggleViewMode"
      >
        <v-icon class="mr-1">{{ viewMode === 'grid' ? 'mdi-calendar-month' : 'mdi-view-grid' }}</v-icon>
        <span>{{ viewMode === 'grid' ? 'Schedule' : 'Grid' }}</span>
      </v-btn>
    </template>
    <template #pending="{ items, loading }">
      <CalendarGrid :calendars="(items as Calendar[])" @accept="acceptCalendarAccess" @decline="onDeclinetCalendar" :acceptingCalendarId="acceptingCalendarId" :loading="(loading as boolean)" />
    </template>
    <template #explore="{ items, loading }">
      <div class="d-flex justify-space-between align-center mb-4">
      </div>
      <CalendarGrid :calendars="(items as Calendar[])" :loading="(loading as boolean)" />
    </template>
    <template #fab>
        <!-- Create Calendar (only in grid view) -->
        <v-btn
          v-if="selectedAccount?.value === authStore.user.name && viewMode === 'grid'"
          color="primary"
          density="compact"
          style="position: fixed; bottom: 16px; right: 16px"
          :to="{ name: 'calendarCreate' }"
        >
          <v-icon>mdi-plus</v-icon>
          <span>Create Calendar</span>
        </v-btn>
      </template>
  </ListTabsPage>

  <!-- Simple filter dialog for account selection -->
  <v-dialog v-model="showFilterModal" max-width="500">
    <v-card>
      <v-card-title>Filter Calendars</v-card-title>
      <v-card-text>
        <div class="mb-4">
          <div class="font-weight-bold mb-2">Account</div>
          <v-autocomplete
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
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text @click="showFilterModal = false">Close</v-btn>
        <v-btn color="primary" @click="showFilterModal = false">Ok</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useCalendarsStore } from '@/stores/calendar'
import { useAuthStore } from '@/stores/auth'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import CalendarGrid from '@/components/CalendarGrid.vue'
import { useAlertStore } from '@/stores/alerts'
import { useCirclesStore } from '@/stores/circles'
import { useUsersStore } from '@/stores/users'
import type { Calendar } from '@/genapi/api/calendars/calendar/v1alpha1'
import ScheduleCal from '@/views/calendar/event/ScheduleCal.vue'
import type { Event } from '@/genapi/api/calendars/calendar/v1alpha1'

const calendarsStore = useCalendarsStore()
const authStore = useAuthStore()
const circlesStore = useCirclesStore()
const usersStore = useUsersStore()

const acceptingCalendarId = ref<string | null>(null)
const tabsPage = ref()
const alertsStore = useAlertStore()
const viewMode = ref<'grid' | 'schedule'>('schedule')
const events = ref<Event[]>([])
const currentTab = ref<string>('my')
const isMyTab = computed(() => currentTab.value === 'my')
const loading = ref<boolean>(true)

// Account selection like RecipesView
type NamedEntity = { name: string }
type AccountOption = { label: string; value: string; account: NamedEntity; icon: string }

const selectedAccount = ref<AccountOption | null>(null)
const showFilterModal = ref(false)
const settingsOpen = ref(false)

const accountOptions = computed<AccountOption[]>(() => {
  const options: AccountOption[] = []
  if (authStore.user?.name) {
    const userName = authStore.user.name as string
    options.push({ label: 'My Calendars', value: userName, account: { name: userName }, icon: 'mdi-account-circle' })
    if (Array.isArray(circlesStore.myCircles)) {
      for (const circle of circlesStore.myCircles) {
        if (!circle?.name) continue
        options.push({ label: circle.title || 'Untitled Circle', value: circle.name as string, account: { name: circle.name as string }, icon: 'mdi-account-group' })
      }
    }
    if (Array.isArray(usersStore.friends)) {
      for (const user of usersStore.friends) {
        if (!user?.name) continue
        let label = ''
        if (user.givenName || user.familyName) {
          label = (user.givenName + ' ' + user.familyName).trim()
        } else if (user.username) {
          label = user.username
        }
        options.push({ label, value: user.name as string, account: { name: user.name as string }, icon: 'mdi-account-circle' })
      }
    }
  }
  return options
})

watch(
  () => authStore.user?.name,
  (newUserName) => {
    if (newUserName && (!selectedAccount.value || selectedAccount.value.value !== newUserName)) {
      selectedAccount.value = accountOptions.value[0]
    }
  },
  { immediate: true }
)

watch(selectedAccount, async () => {
  await nextTick()
  // Reload current tab with selected account context if needed later
  if (isMyTab.value) {
      tabsPage.value?.reloadActiveTab()
  }
})

// Search like RecipesView (grid or non-my)
const searchQuery = ref('')
const searchExpanded = ref(false)
const searchInput = ref<HTMLInputElement | null>(null)

function onSearchEnter() { searchExpanded.value = false }
function expandSearch() { searchExpanded.value = true; nextTick(() => { searchInput.value?.focus?.() }) }
function onSearchFocus() { searchExpanded.value = true }
function onSearchBlur() { if (!searchQuery.value) searchExpanded.value = false }

const searchBarStyle = computed(() => {
  return searchExpanded.value || searchQuery.value
    ? { maxWidth: '350px', width: '100%', transition: 'max-width 0.3s cubic-bezier(0.4,0,0.2,1)' }
    : { maxWidth: '44px', width: '44px', transition: 'max-width 0.3s cubic-bezier(0.4,0,0.2,1)' }
})

function getFilteredCalendars(items: Calendar[]) {
  if (!searchQuery.value) return items
  const q = searchQuery.value.toLowerCase()
  return items.filter(c => (c.title || '').toLowerCase().includes(q))
}

const tabs = [
  { 
    label: 'My Calendars', 
    value: 'my', 
    icon: 'mdi-calendar-account',
    loader: async () => {
      const account = selectedAccount.value?.account
      await calendarsStore.loadMyCalendars(account?.name || '')
      return [...calendarsStore.myCalendars]
    }
  },
  { 
    label: 'Pending',  
    value: 'pending', 
    icon: 'mdi-calendar-clock',
    loader: async () => {
      await calendarsStore.loadPendingCalendars(authStore.user?.name || '')
      return [...calendarsStore.sharedPendingCalendars]
    }
  },
  { 
    label: 'Explore',  
    value: 'explore', 
    icon: 'mdi-compass',
    loader: async () => {
      await calendarsStore.loadPublicCalendars(authStore.user?.name || '')
      return [...calendarsStore.publicCalendars]
    }
  }
]

async function acceptCalendarAccess(calendar: Calendar) {
  if (!calendar.name) return
  
  acceptingCalendarId.value = calendar.name
  try {
    await calendarsStore.acceptCalendar(calendar.calendarAccess?.name || '')
    tabsPage.value?.reloadTab('pending')
  } catch (error) {
    alertsStore.addAlert(`Failed to accept calendar access: ${error}`)
  } finally {
    acceptingCalendarId.value = null
  }
}

async function onDeclinetCalendar(calendar: Calendar) {
  if (!calendar.name) return
  
  try {
    await calendarsStore.deleteCalendarAccess(calendar.name)
    // Refresh the calendars list
    tabsPage.value?.reloadTab('pending')
  } catch (error) {
    alertsStore.addAlert(`Failed to decline calendar access: ${error}`)
  }
}

function toggleViewMode() {
  viewMode.value = viewMode.value === 'grid' ? 'schedule' : 'grid'
  if (viewMode.value === 'schedule' && isMyTab.value) {
    void loadEventsForMyCalendars()
  }
}

// Load events for all calendars currently in the "My" tab
async function loadEventsForMyCalendars() {
  loading.value = true
  events.value = []
  for (const calendar of calendarsStore.myCalendars) {
    if (!calendar?.name) continue
    const es = await calendarsStore.loadEvents(calendar.name)
    events.value.push(...es)
  }
  loading.value = false
}

async function onEventUpdated() {
  if (viewMode.value === 'schedule' && isMyTab.value) {
    await loadEventsForMyCalendars()
  }
}

async function onEventDeleted() {
  if (viewMode.value === 'schedule' && isMyTab.value) {
    await loadEventsForMyCalendars()
  }
}

async function onEventCreated() {
  if (viewMode.value === 'schedule' && isMyTab.value) {
    await loadEventsForMyCalendars()
  }
}

// Keep events in sync when the calendars list changes or when switching to schedule mode
watch(
  () => calendarsStore.myCalendars,
  async () => {
    if (viewMode.value === 'schedule' && isMyTab.value) {
      await loadEventsForMyCalendars()
    }
  },
  { deep: true }
)

// When switching tabs, ensure events are loaded if entering My tab in schedule mode
watch(
  () => tabsPage.value?.activeTab?.value,
  async (tab) => {
    if (tab === 'my' && viewMode.value === 'schedule') {
      await loadEventsForMyCalendars()
    }
  }
)

</script>

<style scoped>
</style>