<template>
  <v-container v-if="user">
    <ListTabsPage :tabs="tabs" ref="tabsPage">
      <template #general>
        <v-card class="mx-auto" max-width="600">
          <div class="image-container">
            <v-img
              v-if="user.imageUri"
              class="mt-1"
              style="background-color: lightgray"
              :src="user.imageUri"
              cover
              height="300"
            ></v-img>
            <div v-else class="mt-1 d-flex align-center justify-center" style="background-color: lightgray; height: 300px; border-radius: 4px;">
              <div class="text-center">
                <v-icon size="64" color="grey-darken-1">mdi-image-outline</v-icon>
                <div class="text-grey-darken-1 mt-2">No image available</div>
              </div>
            </div>
          </div>
          <div class="bio-section pa-4">
            <div class="text-subtitle-1 font-weight-bold mb-1">Bio</div>
            <div class="text-body-1" style="white-space: pre-line;">
              {{ user.bio || 'No bio set.' }}
            </div>
          </div>
          <v-card-title>User Settings</v-card-title>
          <v-card-text>
            <v-list>
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon icon="mdi-account"></v-icon>
                </template>
                <v-list-item-title>Given Name</v-list-item-title>
                <v-list-item-subtitle>{{ user.givenName || 'Not set' }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon icon="mdi-account"></v-icon>
                </template>
                <v-list-item-title>Family Name</v-list-item-title>
                <v-list-item-subtitle>{{ user.familyName || 'Not set' }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon icon="mdi-account-circle"></v-icon>
                </template>
                <v-list-item-title>Username</v-list-item-title>
                <v-list-item-subtitle>{{ user.username }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon :icon="selectedVisibilityIcon"></v-icon>
                </template>
                <v-list-item-title>Visibility</v-list-item-title>
                <v-list-item-subtitle>
                  <strong>{{ selectedVisibilityLabel }}</strong>: {{ selectedVisibilityDescription }}
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <!-- Removed Edit Settings button -->
          </v-card-actions>
        </v-card>
      </template>
      <template #recipes="{ items, loading }">
        <RecipeGrid :recipes="items" :loading="loading" />
      </template>
    </ListTabsPage>
  </v-container>
  <template v-if="user && authStore.user.name && user.name === authStore.user.name">
    <v-fab
      location="bottom right"
      color="primary"
      icon
      style="position: fixed; bottom: 24px; right: 24px; z-index: 10;"
      @click="router.push({ name: 'user-edit', params: { userId: user.name } })"
    >
      <v-icon>mdi-pencil</v-icon>
    </v-fab>
  </template>
</template>

<script setup lang="ts">
import { useUsersStore } from '@/stores/users'
import { useRecipesStore } from '@/stores/recipes'
import { storeToRefs } from 'pinia'
import { onMounted, ref, computed } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRoute } from 'vue-router'
import ListTabsPage from '@/components/common/ListTabsPage.vue'
import RecipeGrid from '@/components/RecipeGrid.vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const usersStore = useUsersStore()
const recipesStore = useRecipesStore()
const { currentUser: user } = storeToRefs(usersStore)
const breadcrumbStore = useBreadcrumbStore()
const route = useRoute()
const tabsPage = ref()

const userId = computed(() => String(route.params.userId || ''))

onMounted(async () => {
  await Promise.all([
    usersStore.loadUser(userId.value),
  ])
  breadcrumbStore.setBreadcrumbs([
    { title: 'User Settings', to: { name: 'user', params: { userId: userId.value } } },
  ])
})

const visibilityOptions = [
  {
    value: 'VISIBILITY_LEVEL_PUBLIC',
    label: 'Public',
    icon: 'mdi-earth',
    color: 'success',
    description: 'Everyone can see your profile.'
  },
  {
    value: 'VISIBILITY_LEVEL_RESTRICTED',
    label: 'Restricted',
    icon: 'mdi-account-group',
    color: 'warning',
    description: 'Only shared users and their connections can see your profile.'
  },
  {
    value: 'VISIBILITY_LEVEL_PRIVATE',
    label: 'Private',
    icon: 'mdi-lock',
    color: 'info',
    description: 'Only specifically shared users can see your profile.'
  },
  {
    value: 'VISIBILITY_LEVEL_HIDDEN',
    label: 'Hidden',
    icon: 'mdi-eye-off',
    color: 'secondary',
    description: 'Only you can see your profile.'
  }
]

const selectedVisibility = computed(() => {
  return visibilityOptions.find(option => option.value === user.value?.visibility)
})
const selectedVisibilityLabel = computed(() => selectedVisibility.value?.label || '')
const selectedVisibilityDescription = computed(() => selectedVisibility.value?.description || '')
const selectedVisibilityIcon = computed(() => selectedVisibility.value?.icon || 'mdi-help-circle')

const tabs = [
  {
    label: 'General',
    value: 'general',
  },
  {
    label: 'Recipes',
    value: 'recipes',
    loader: async () => {
      if (!user.value?.name) return []
      // The parent resource for user recipes is the user's resource name
      await recipesStore.loadMyRecipes(user.value.name)
      return [...recipesStore.myRecipes]
    },
  },
]

const authStore = useAuthStore()
const router = useRouter()
</script>

<style></style>
