<template>
  <v-container max-width="600">
    <div v-for="(direction, i) in recipe.directions" :key="i" class="d-flex gap-2 mb-4">
      <move-buttons
        :show-up="i > 0"
        :show-down="i < (recipe.directions?.length || 0) - 1"
        @move-up="moveArrayItem(recipe.directions, i, 'up')"
        @move-down="moveArrayItem(recipe.directions, i, 'down')"
      />
      <v-card class="flex-grow-1">
        <v-card-title class="d-flex align-center">
          <v-text-field
            v-model="direction.title"
            class="flex-grow-1"
            placeholder="Direction Group Title"
          ></v-text-field>
          <v-btn icon="mdi-delete" variant="text" color="error" @click="removeDirection(i)"></v-btn>
        </v-card-title>
        <v-card-text>
          <v-list>
            <div v-for="(step, n) in direction.steps || []" :key="n" class="d-flex gap-2 mb-2">
              <move-buttons
                size="x-small"
                :show-up="n > 0"
                :show-down="n < (direction.steps?.length || 0) - 1"
                @move-up="moveNestedArrayItem(recipe.directions, i, 'steps', n, 'up')"
                @move-down="moveNestedArrayItem(recipe.directions, i, 'steps', n, 'down')"
              />
              <v-list-item class="flex-grow-1">
                <div class="d-flex align-center gap-2">
                  <div class="font-weight-bold">Step {{ n + 1 }}</div>
                  <v-spacer></v-spacer>
                  <v-btn
                    icon="mdi-delete"
                    size="small"
                    variant="text"
                    color="error"
                    @click="removeStep(i, n)"
                  ></v-btn>
                </div>
                <v-textarea
                  v-model="direction.steps![n]"
                  placeholder="Enter step instructions"
                ></v-textarea>
              </v-list-item>
            </div>
          </v-list>
          <v-btn block variant="text" prepend-icon="mdi-plus" class="mt-2" @click="addStep(i)"
            >Add Step</v-btn
          >
        </v-card-text>
      </v-card>
    </div>
    <v-btn block variant="outlined" prepend-icon="mdi-plus" class="mt-4" @click="addDirection"
      >Add Direction</v-btn
    >
  </v-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Recipe, Recipe_Direction } from '@/genapi/api/meals/recipe/v1alpha1'
import MoveButtons from '@/components/common/MoveButtons.vue'
import { moveArrayItem, moveNestedArrayItem } from '@/utils/array'

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

function addDirection() {
  if (!recipe.value) return
  if (!recipe.value.directions) {
    recipe.value.directions = []
  }
  recipe.value.directions.push({
    title: '',
    steps: [],
  } as Recipe_Direction)
}

function removeDirection(index: number) {
  recipe.value?.directions?.splice(index, 1)
}

function addStep(directionIndex: number) {
  const direction = recipe.value?.directions?.[directionIndex]
  if (!direction) return
  if (!direction.steps) {
    direction.steps = []
  }
  direction.steps.push('')
}

function removeStep(directionIndex: number, stepIndex: number) {
  recipe.value?.directions?.[directionIndex].steps?.splice(stepIndex, 1)
}
</script>

<style scoped>
.gap-1 {
  gap: 4px;
}

.gap-2 {
  gap: 8px;
}
</style>
