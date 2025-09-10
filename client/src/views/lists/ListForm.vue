<template>
  <v-form ref="form" v-model="valid" @submit.prevent="handleSave">
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


        <!-- List Sections -->
        <v-col cols="12">
          <div class="text-h6 mb-4">List Sections</div>
          
          <!-- Existing Sections -->
          <div v-for="(section, index) in listStore.list.sections" :key="index" class="mb-4">
            <v-card variant="outlined">
              <v-card-text>
                <v-text-field
                  v-model="section.title"
                  label="Section Title"
                  variant="outlined"
                  density="compact"
                  :rules="sectionTitleRules"
                  prepend-inner-icon="mdi-folder-outline"
                  counter="100"
                  maxlength="100"
                />
              </v-card-text>
              <v-card-actions>
                <v-btn
                  color="error"
                  variant="text"
                  size="small"
                  @click="removeSection(index)"
                >
                  <v-icon>mdi-delete</v-icon>
                  Remove
                </v-btn>
              </v-card-actions>
            </v-card>
          </div>

          <!-- Add Section Button -->
          <v-btn
            color="primary"
            variant="outlined"
            prepend-icon="mdi-plus"
            @click="addSection"
            class="mb-4"
          >
            Add Section
          </v-btn>
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
      style="position: fixed; bottom: 46px; left: 16px"
      @click="handleCancel"
    >
      <v-icon>mdi-close</v-icon>
      Cancel
    </v-btn>

    <!-- Save FAB -->
    <v-btn
      color="success"
      density="compact"
      style="position: fixed; bottom: 46px; right: 16px"
      :loading="saving"
      @click="handleSave"
    >
      <v-icon>mdi-content-save</v-icon>
      {{ isEditMode ? 'Save Changes' : 'Create List' }}
    </v-btn>
  </v-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useListStore } from '@/stores/list'
import type { List_ListSection } from '@/genapi/api/lists/list/v1alpha1'

interface Props {
  isEditMode?: boolean
  saving?: boolean
}

interface Emits {
  (e: 'save'): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  isEditMode: false,
  saving: false
})

const emit = defineEmits<Emits>()

const listStore = useListStore()

const form = ref()
const valid = ref(false)

// Form validation rules
const titleRules = [
  (v: string) => !!v || 'Title is required',
  (v: string) => (v && v.length >= 3) || 'Title must be at least 3 characters',
  (v: string) => (v && v.length <= 100) || 'Title must be less than 100 characters',
]

const sectionTitleRules = [
  (v: string) => !!v || 'Section title is required',
  (v: string) => (v && v.length >= 1) || 'Section title must be at least 1 character',
  (v: string) => (v && v.length <= 100) || 'Section title must be less than 100 characters',
]

// Visibility options
const visibilityOptions = [
  { title: 'Public', value: 'VISIBILITY_LEVEL_PUBLIC', subtitle: 'Anyone can see this list' },
  { title: 'Restricted', value: 'VISIBILITY_LEVEL_RESTRICTED', subtitle: 'Only people with access can see this list' },
  { title: 'Private', value: 'VISIBILITY_LEVEL_PRIVATE', subtitle: 'Only you can see this list' },
  { title: 'Hidden', value: 'VISIBILITY_LEVEL_HIDDEN', subtitle: 'List is hidden from everyone' },
]

// Section management
function addSection() {
  if (!listStore.list) return
  
  const newSection: List_ListSection = {
    name: undefined, // Will be set by server
    title: 'New Section',
  }
  
  if (!listStore.list.sections) {
    listStore.list.sections = []
  }
  
  listStore.list.sections.push(newSection)
}

function removeSection(index: number) {
  if (!listStore.list?.sections) return
  
  listStore.list.sections.splice(index, 1)
}

// Event handlers
function handleSave() {
  emit('save')
}

function handleCancel() {
  emit('cancel')
}
</script>
