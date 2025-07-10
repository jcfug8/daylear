<template>
  <v-container max-width="600">
    <div v-for="(ingredientGroup, i) in recipe.ingredientGroups" :key="i" class="d-flex gap-2 mb-4">
      <move-buttons
        :show-up="i > 0"
        :show-down="i < (recipe.ingredientGroups?.length || 0) - 1"
        @move-up="moveArrayItem(recipe.ingredientGroups, i, 'up')"
        @move-down="moveArrayItem(recipe.ingredientGroups, i, 'down')"
      />
      <v-card class="flex-grow-1">
        <v-card-title class="d-flex align-center">
          <v-text-field
            v-model="ingredientGroup.title"
            class="flex-grow-1"
            placeholder="Ingredient Group Title"
          ></v-text-field>
          <v-btn
            icon="mdi-delete"
            variant="text"
            color="error"
            @click="removeIngredientGroup(i)"
          ></v-btn>
        </v-card-title>
        <v-card-text>
          <div
            v-for="(ingredient, j) in ingredientGroup.ingredients"
            :key="j"
            class="d-flex gap-2 mb-2"
          >
            <move-buttons
              size="x-small"
              :show-up="j > 0"
              :show-down="j < (ingredientGroup.ingredients?.length || 0) - 1"
              @move-up="moveNestedArrayItem(recipe.ingredientGroups, i, 'ingredients', j, 'up')"
              @move-down="moveNestedArrayItem(recipe.ingredientGroups, i, 'ingredients', j, 'down')"
            />
            <div class="flex-grow-1">
              <v-row dense>
                <v-col cols="4" sm="2">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    hide-details
                    v-model.number="ingredient.measurementAmount"
                    type="number"
                    min="0"
                    placeholder="Amount"
                    class="mt-0"
                  ></v-text-field>
                </v-col>
                <v-col cols="8" sm="4">
                  <v-select
                    density="compact"
                    variant="outlined"
                    hide-details
                    v-model="ingredient.measurementType"
                    :items="MEASUREMENT_TYPES"
                    item-title="title"
                    item-value="value"
                    class="mt-0"
                  ></v-select>
                </v-col>
                <v-col cols="6" sm="5">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    hide-details
                    v-model="ingredient.title"
                    placeholder="Ingredient Name"
                    class="mt-0"
                  ></v-text-field>
                </v-col>
                <v-col cols="6" sm="1" class="d-flex justify-end align-center">
                  <v-btn
                    icon="mdi-delete"
                    size="small"
                    variant="text"
                    color="error"
                    @click="removeIngredient(i, j)"
                  ></v-btn>
                </v-col>
                <v-col cols="12" sm="12" class="d-flex align-center">
                  <v-btn
                    size="small"
                    variant="text"
                    color="primary"
                    @click="toggleSecondMeasurement(i, j)"
                  >
                    {{ showSecondMeasurementMap[i]?.[j] ? 'Remove second measurement' : 'Add another measurement' }}
                  </v-btn>
                </v-col>
                <template v-if="showSecondMeasurementMap[i]?.[j]">
                  <v-col cols="12" sm="3">
                    <v-select
                      density="compact"
                      variant="outlined"
                      hide-details
                      v-model="ingredient.measurementConjunction"
                      :items="MEASUREMENT_CONJUNCTIONS"
                      item-title="title"
                      item-value="value"
                      label="Conjunction"
                      class="mt-0"
                    ></v-select>
                  </v-col>
                  <v-col cols="4" sm="2">
                    <v-text-field
                      density="compact"
                      variant="outlined"
                      hide-details
                      v-model.number="ingredient.secondMeasurementAmount"
                      type="number"
                      min="0"
                      placeholder="Second Amount"
                      class="mt-0"
                    ></v-text-field>
                  </v-col>
                  <v-col cols="8" sm="4">
                    <v-select
                      density="compact"
                      variant="outlined"
                      hide-details
                      v-model="ingredient.secondMeasurementType"
                      :items="MEASUREMENT_TYPES"
                      item-title="title"
                      item-value="value"
                      placeholder="Second Type"
                      class="mt-0"
                    ></v-select>
                  </v-col>
                </template>
                <v-col cols="6" sm="12" class="d-flex justify-end">
                  <v-checkbox
                    density="compact"
                    hide-details
                    v-model="ingredient.optional"
                    label="Optional"
                    class="mt-0"
                  ></v-checkbox>
                </v-col>
              </v-row>
            </div>
          </div>
          <v-btn block variant="text" prepend-icon="mdi-plus" class="mt-2" @click="addIngredient(i)"
            >Add Ingredient</v-btn
          >
        </v-card-text>
      </v-card>
    </div>
    <v-btn block variant="outlined" prepend-icon="mdi-plus" class="mt-4" @click="addIngredientGroup"
      >Add Ingredient Group</v-btn
    >
  </v-container>
</template>

<script setup lang="ts">
import { computed, reactive, onMounted, watch } from 'vue'
import type {
  Recipe,
  Recipe_IngredientGroup,
  Recipe_Ingredient,
} from '@/genapi/api/meals/recipe/v1alpha1'
import { MEASUREMENT_TYPES } from '@/constants/measurements'
import MoveButtons from '@/components/common/MoveButtons.vue'
import { moveArrayItem, moveNestedArrayItem } from '@/utils/array'

onMounted(() => {
  initializeShowSecondMeasurementMap()
})

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
 
watch(recipe, () => {
  initializeShowSecondMeasurementMap()
})

// Track showSecondMeasurement state for each ingredient by group and index
const showSecondMeasurementMap = reactive<{ [groupIdx: number]: { [ingredientIdx: number]: boolean } }>({})
function initializeShowSecondMeasurementMap() {
  recipe.value.ingredientGroups?.forEach((group, groupIdx) => {
    if (!showSecondMeasurementMap[groupIdx]) showSecondMeasurementMap[groupIdx] = {}
    group.ingredients?.forEach((ingredient, ingredientIdx) => {
      const hasSecond =
        ingredient.measurementConjunction !== 'MEASUREMENT_CONJUNCTION_UNSPECIFIED' ||
        (ingredient.secondMeasurementAmount !== undefined && ingredient.secondMeasurementAmount !== 0) ||
        ingredient.secondMeasurementType !== 'MEASUREMENT_TYPE_UNSPECIFIED'
      showSecondMeasurementMap[groupIdx][ingredientIdx] = hasSecond
    })
  })
}

function toggleSecondMeasurement(groupIdx: number, ingredientIdx: number) {
  if (!showSecondMeasurementMap[groupIdx]) showSecondMeasurementMap[groupIdx] = {}
  const current = !!showSecondMeasurementMap[groupIdx][ingredientIdx]
  showSecondMeasurementMap[groupIdx][ingredientIdx] = !current
  if (!showSecondMeasurementMap[groupIdx][ingredientIdx]) {
    const ingredient = recipe.value.ingredientGroups?.[groupIdx]?.ingredients?.[ingredientIdx]
    if (ingredient) {
      ingredient.measurementConjunction = undefined
      ingredient.secondMeasurementAmount = undefined
      ingredient.secondMeasurementType = undefined
    }
  }
}

const MEASUREMENT_CONJUNCTIONS = [
  { title: '', value: 'MEASUREMENT_CONJUNCTION_UNSPECIFIED' },
  { title: 'and', value: 'MEASUREMENT_CONJUNCTION_AND' },
  { title: 'to', value: 'MEASUREMENT_CONJUNCTION_TO' },
  { title: 'or', value: 'MEASUREMENT_CONJUNCTION_OR' },
]

function addIngredientGroup() {
  if (!recipe.value) return
  if (!recipe.value.ingredientGroups) {
    recipe.value.ingredientGroups = []
  }
  recipe.value.ingredientGroups.push({
    title: '',
    ingredients: [],
  } as Recipe_IngredientGroup)
}

function removeIngredientGroup(index: number) {
  recipe.value?.ingredientGroups?.splice(index, 1)
}

function addIngredient(groupIndex: number) {
  const group = recipe.value?.ingredientGroups?.[groupIndex]
  if (!group) return
  if (!group.ingredients) {
    group.ingredients = []
  }
  group.ingredients.push({
    title: '',
    measurementAmount: 0,
    measurementType: 'MEASUREMENT_TYPE_UNSPECIFIED',
    optional: false,
    showSecondMeasurement: false,
  } as any)
}

function removeIngredient(groupIndex: number, ingredientIndex: number) {
  recipe.value?.ingredientGroups?.[groupIndex].ingredients?.splice(ingredientIndex, 1)
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
