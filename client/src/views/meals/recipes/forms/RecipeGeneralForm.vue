<template>
  <v-container max-width="600" class="pa-1">
    <v-row>
      <v-col class="pt-5">
        <div class="text-h4">
          <v-text-field
            density="compact"
            v-model="recipe.title"
            placeholder="Recipe Title"
          ></v-text-field>
        </div>
        <div class="text-body-1">
          <v-textarea
            density="compact"
            v-model="recipe.description"
            placeholder="Recipe Description"
          ></v-textarea>
        </div>
      </v-col>
    </v-row>
    <v-row>
      <v-spacer></v-spacer>
      <v-col align-self="auto" cols="12" sm="8">
        <v-img
          class="mt-1"
          style="background-color: lightgray"
          :src="recipe.imageUri"
          cover
        ></v-img>
      </v-col>
      <v-spacer></v-spacer>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Recipe } from '@/genapi/api/meals/recipe/v1alpha1'

const props = defineProps<{
  modelValue: Recipe
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Recipe): void
}>()

const recipe = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})
</script>
