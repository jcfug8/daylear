<template>
  <v-container class="d-flex justify-center align-center" style="height: 75vh">
    <v-sheet elevation="2" max-width="400" class="pa-4 ma-4">
      <v-form>
        <h1 class="font-weight-thin text-center">Login</h1>
        <v-divider class="mb-2"></v-divider>
        <v-btn prepend-icon="mdi-google" @click="withGooge">Login With Google</v-btn>
      </v-form>
    </v-sheet>
  </v-container>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { watch } from 'vue'
import { API_BASE_URL } from '@/constants/api'

const router = useRouter()
const authStore = useAuthStore()

watch(
  () => authStore.isLoggedIn,
  (newValue: boolean) => {
    if (newValue) {
      router.push({ name: 'recipes' })
    }
  },
)

const withGooge = () => {
  window.location.href = API_BASE_URL + 'auth/google'
}
</script>

<style></style>
