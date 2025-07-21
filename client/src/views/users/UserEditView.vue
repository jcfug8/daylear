<template>
  <v-container>
    <v-card class="mx-auto" max-width="600">
      <v-card-title>Edit User Settings</v-card-title>
      <v-card-text>
        <div class="image-container mb-4" style="position: relative;">
          <v-img
            class="mt-1"
            style="background-color: lightgray"
            :src="imageCleared ? '' : (previewImage || editedUser.imageUri)"
            cover
            height="300"
          >
            <template #placeholder>
              <v-row class="fill-height ma-0" align="center" justify="center">
                <v-icon size="64" color="grey-darken-1">mdi-image-outline</v-icon>
              </v-row>
            </template>
            <template #error>
              <v-row class="fill-height ma-0" align="center" justify="center">
                <v-icon size="64" color="grey-darken-1">mdi-image-outline</v-icon>
              </v-row>
            </template>
          </v-img>
          <v-btn
            icon="mdi-camera"
            color="primary"
            class="image-upload-btn"
            @click="showImageDialog = true"
          ></v-btn>
          <v-btn
            v-if="(previewImage || editedUser.imageUri) && !imageCleared"
            icon="mdi-close"
            color="warning"
            class="image-x-btn"
            @click="clearImage"
            title="Remove Image"
          ></v-btn>
          <v-btn
            v-if="imageCleared"
            icon="mdi-arrow-u-left-top"
            color="info"
            class="image-undo-btn"
            @click="undoClearImage"
            title="Undo Remove Image"
          ></v-btn>
        </div>
        <v-textarea
          v-model="editedUser.bio"
          label="Bio"
          placeholder="Tell us about yourself..."
          rows="3"
          auto-grow
          class="mb-4"
        ></v-textarea>
        <v-dialog v-model="showImageDialog" max-width="500">
          <v-card>
            <v-card-title>Add User Image</v-card-title>
            <v-card-text>
              <v-file-input
                v-model="imageFile"
                label="Choose Image"
                accept="image/jpeg, image/jpg, image/png, image/gif, image/webp, image/bmp, image/svg, image/heic, image/heif"
                class="mt-4"
                prepend-icon="mdi-camera"
                @update:model-value="handleFileSelect"
              ></v-file-input>
              <v-img
                v-if="previewImage"
                :src="previewImage"
                max-height="200"
                contain
                class="mt-4"
              ></v-img>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="error" @click="cancelImageDialog">Cancel</v-btn>
              <v-btn color="primary" @click="handleImageSubmit">OK</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <v-form @submit.prevent="saveSettings">
          <v-text-field
            v-model="editedUser.givenName"
            label="Given Name"
          ></v-text-field>

          <v-text-field
            v-model="editedUser.familyName"
            label="Family Name"
          ></v-text-field>

          <v-text-field
            disabled
            v-model="editedUser.email"
            label="Email"
            type="email"
            required
          ></v-text-field>

          <v-text-field
            v-model="editedUser.username"
            label="Username"
            required
          ></v-text-field>

          <v-select
            v-model="editedUser.visibility"
            :items="visibilityOptions"
            item-title="label"
            item-value="value"
            label="Profile Visibility"
            density="compact"
            variant="outlined"
          >
            <template #selection="{ item }">
              <div class="d-flex align-center">
                <v-icon :icon="item.raw.icon" class="me-2" size="small"></v-icon>
                {{ item.raw.label }}
              </div>
            </template>
            <template #item="{ props, item }">
              <v-list-item v-bind="props">
                <template #prepend>
                  <v-icon :icon="item.raw.icon" size="small"></v-icon>
                </template>
                <v-list-item-subtitle class="text-wrap">
                  {{ item.raw.description }}
                </v-list-item-subtitle>
              </v-list-item>
            </template>
          </v-select>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="secondary"
          @click="navigateBack"
        >
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          @click="saveSettings"
        >
          Save Changes
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { onMounted, ref, computed } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumbs'
import { useRouter } from 'vue-router'
import type { User, UserSettings, apitypes_VisibilityLevel } from '@/genapi/api/users/user/v1alpha1'
import { useAlertStore } from '@/stores/alerts'
import { fileService } from '@/api/api'

const router = useRouter()
const authStore = useAuthStore()
const { user, userSettings } = storeToRefs(authStore)
const breadcrumbStore = useBreadcrumbStore()
const alertStore = useAlertStore()

// Create a copy of the user object for editing
const editedUser = ref<User & UserSettings>({
  name: user.value.name,
  email: userSettings.value.email,
  username: user.value.username,
  givenName: user.value.givenName,
  familyName: user.value.familyName,
  visibility: (user.value.visibility || 'VISIBILITY_LEVEL_PUBLIC') as apitypes_VisibilityLevel,
  imageUri: user.value.imageUri,
  access: user.value.access,
  bio: user.value.bio,
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

const showImageDialog = ref(false)
const imageFile = ref<File | null>(null)
const previewImage = ref<string | null>(null)
const imageCleared = ref(false)
const originalImageUri = ref<string | null>(null)

function handleFileSelect(files: File | File[] | null) {
  if (files instanceof File) {
    const reader = new FileReader()
    reader.onload = (e) => {
      previewImage.value = e.target?.result as string
    }
    reader.readAsDataURL(files)
  } else {
    previewImage.value = null
  }
}

function cancelImageDialog() {
  showImageDialog.value = false
  imageFile.value = null
  previewImage.value = null
}

function handleImageSubmit() {
  showImageDialog.value = false
}

function clearImage() {
  if (!imageCleared.value) {
    originalImageUri.value = previewImage.value || editedUser.value.imageUri || null
    imageCleared.value = true
    previewImage.value = null
    editedUser.value.imageUri = ''
  }
}

function undoClearImage() {
  if (imageCleared.value && originalImageUri.value) {
    editedUser.value.imageUri = originalImageUri.value
    previewImage.value = originalImageUri.value
    imageCleared.value = false
  }
}

function navigateBack() {
  router.push({ name: 'user', params: { userId: user.value.name } })
}

async function saveSettings() {
  try {
    await authStore.updateAuthUser(editedUser.value)
    // Upload image if there's a pending file
    if (imageFile.value && editedUser.value.name) {
      const response = await fileService.UploadUserImage({
        name: editedUser.value.name,
        file: imageFile.value,
      })
      editedUser.value.imageUri = response.imageUri
    }
    navigateBack()
  } catch (err) {
    console.log('Error saving settings:', err)
    alertStore.addAlert(err instanceof Error ? "Unable to save settings\n" + err.message : String(err), 'error')
  }
}

onMounted(async () => {
  await authStore.loadAuthUser()

  breadcrumbStore.setBreadcrumbs([
    { title: 'User Settings', to: { name: 'user', params: { userId: user.value.name } } },
    { title: 'Edit', to: { name: 'user-edit', params: { userId: user.value.name } } },
  ])
})
</script>

<style>
.image-container {
  position: relative;
}
.image-upload-btn {
  position: absolute;
  bottom: 16px;
  right: 16px;
}
.image-x-btn, .image-undo-btn {
  position: absolute;
  top: 16px;
  right: 16px;
}
</style>
