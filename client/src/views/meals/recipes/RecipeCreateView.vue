<template>
  <!-- Initial Choice Dialog -->
  <v-dialog v-model="showChoiceDialog" max-width="500" persistent>
    <v-card>
      <v-card-title class="text-h5 text-center pa-6">
        Create New Recipe
      </v-card-title>
      <v-card-text class="text-center pa-6">
        <p class="text-body-1 mb-6">
          Choose how you'd like to create your recipe:
        </p>
        <v-row>
          <v-col cols="12">
            <v-btn
              block
              size="large"
              color="primary"
              variant="outlined"
              class="mb-3"
              prepend-icon="mdi-pencil"
              @click="handleManualEntry"
            >
              Manually Enter Recipe
            </v-btn>
          </v-col>
          <v-col cols="12">
            <v-btn
              block
              size="large"
              color="secondary"
              variant="outlined"
              class="mb-3"
              prepend-icon="mdi-camera"
              @click="handleImageImport"
            >
              Import from Image
            </v-btn>
          </v-col>
          <v-col cols="12">
            <v-btn
              block
              size="large"
              color="info"
              variant="outlined"
              prepend-icon="mdi-link"
              @click="handleUrlImport"
            >
              Import from URL
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-dialog>

  <!-- Recipe Form with Import Modal -->
  <recipe-form
    v-if="recipesStore.recipe"
    v-model="recipesStore.recipe"
    :show-import-dialog="showImportDialog"
    :import-tab="importTab"
    @save="saveRecipe"
    @close="navigateBack"
    @close-import-dialog="handleCloseImportDialog"
  />
</template>

<script setup lang="ts">
import { useRecipesStore } from '@/stores/recipes'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import RecipeForm from '@/views/meals/recipes/forms/RecipeForm.vue'
import { useAuthStore } from '@/stores/auth'
import { useAlertStore } from '@/stores/alerts'
import { computed } from 'vue'

const router = useRouter()
const recipesStore = useRecipesStore()
const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()
const alertStore = useAlertStore()
const route = useRoute()

// Dialog state
const showChoiceDialog = ref(true)
const showImportDialog = ref(false)
const importTab = ref('scrape') // 'scrape' for URL, 'ocr' for image

function navigateBack() {
  if (route.params.circleId) {
    router.push({ name: 'circle', params: { circleId: route.params.circleId } })
  } else {
    router.push({ name: 'recipes' })
  }
}

const circleName = computed(() => {
  return route.path.replace('/recipes/create', '').slice(1)
})

async function saveRecipe() {
  if (!authStore.user.name && !route.params.circleId) {
    throw new Error('User not authenticated')
  }
  try {
    const recipe = await recipesStore.createRecipe(circleName.value ? circleName.value : authStore.user.name)
    router.push('/'+recipe.name)
  } catch (err) {
    alertStore.addAlert(err instanceof Error ? "Unable to create recipe\n" + err.message : String(err), 'error')
  }
}

// Handle user choices
function handleManualEntry() {
  showChoiceDialog.value = false
  // Just close the dialog and show the form normally
}

function handleImageImport() {
  showChoiceDialog.value = false
  showImportDialog.value = true
  importTab.value = 'ocr'
}

function handleUrlImport() {
  showChoiceDialog.value = false
  showImportDialog.value = true
  importTab.value = 'scrape'
}

function handleCloseImportDialog() {
  showImportDialog.value = false
}

onMounted(() => {
  // Initialize an empty recipe
  recipesStore.initEmptyRecipe()

  breadcrumbStore.setBreadcrumbs([
    { title: 'Recipes', to: { name: 'recipes' } },
    { title: 'Create New Recipe', to: { name: 'recipeCreate' } },
  ])
})
</script>
