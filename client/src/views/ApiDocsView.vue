<template>
  <v-container class="py-8">
    <v-row justify="center">
      <v-col cols="12" md="8">
        <h1 class="mb-4">API Documentation</h1>
        <v-form v-if="specs.length > 0">
          <v-select
            v-model="selectedUrl"
            :items="specs"
            item-title="name"
            item-value="url"
            label="Select API Spec"
            outlined
            dense
            class="mb-6"
          />
        </v-form>
        <v-alert
          v-if="error"
          type="error"
          class="mb-4"
        >
          {{ error }}
        </v-alert>
        <v-progress-circular
          v-if="loading"
          indeterminate
          color="primary"
          class="mb-4"
        />
        <SwaggerUI v-if="selectedUrl" :url="API_BASE_URL + selectedUrl" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { API_BASE_URL } from '@/constants/api'
import { fetchApiSpecs, type ApiSpec } from '@/api/api'

import { ref, onMounted } from 'vue';
import SwaggerUI from '../components/common/SwaggerUI.vue';

const specs = ref<ApiSpec[]>([]);
const selectedUrl = ref<string>('');
const loading = ref(true);
const error = ref<string>('');

onMounted(async () => {
  try {
    loading.value = true;
    error.value = '';
    const fetchedSpecs = await fetchApiSpecs();
    specs.value = fetchedSpecs;
    if (fetchedSpecs.length > 0) {
      selectedUrl.value = fetchedSpecs[0].url;
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load API specs';
  } finally {
    loading.value = false;
  }
});
</script> 