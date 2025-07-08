<template>
  <v-container class="py-8">
    <v-row justify="center">
      <v-col cols="12" md="8">
        <h1 class="mb-4">API Documentation</h1>
        <v-form>
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
        <SwaggerUI :url="API_BASE_URL + selectedUrl" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { API_BASE_URL } from '@/constants/api'

import { ref } from 'vue';
import SwaggerUI from '../components/common/SwaggerUI.vue';

// openapi/api/circles/circle/v1alpha1/circle.swagger.json
// openapi/api/circles/circle/v1alpha1/public_circle.swagger.json
// openapi/api/meals/recipe/v1alpha1/recipe.swagger.json
// openapi/api/meals/recipe/v1alpha1/recipe_access.swagger.json
// openapi/api/meals/recipe/v1alpha1/recipe_recipient_list.swagger.json
// openapi/api/types/permission_level.swagger.json
// openapi/api/users/user/v1alpha1/public_user.swagger.json
// openapi/api/users/user/v1alpha1/user.swagger.json
const specs = [
  {
    name: 'Circle',
    url: 'api/openapi/api/circles/circle/v1alpha1/circle.swagger.json',
  },
  {
    name: 'Public Circle',
    url: 'api/openapi/api/circles/circle/v1alpha1/public_circle.swagger.json',
  },
  {
    name: 'Recipe',
    url: 'api/openapi/api/meals/recipe/v1alpha1/recipe.swagger.json',
  },
  {
    name: 'Recipe Access',
    url: 'api/openapi/api/meals/recipe/v1alpha1/recipe_access.swagger.json',
  },
  {
    name: 'Recipe Recipient List',
    url: 'api/openapi/api/meals/recipe/v1alpha1/recipe_recipient_list.swagger.json',
  },
  {
    name: 'Permission Level',
    url: 'api/openapi/api/types/permission_level.swagger.json',
  },
  {
    name: 'Public User',
    url: 'api/openapi/api/users/user/v1alpha1/public_user.swagger.json',
  },
];

const selectedUrl = ref(specs[0].url);
</script> 