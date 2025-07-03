import type { Recipe_MeasurementType } from '@/genapi/api/meals/recipe/v1alpha1'

export const MEASUREMENT_TYPES = [
  { value: 'MEASUREMENT_TYPE_UNSPECIFIED', title: '' },
  { value: 'MEASUREMENT_TYPE_CUP', title: 'cups' },
  { value: 'MEASUREMENT_TYPE_TABLESPOON', title: 'tablespoons' },
  { value: 'MEASUREMENT_TYPE_TEASPOON', title: 'teaspoons' },
  { value: 'MEASUREMENT_TYPE_OUNCE', title: 'ounces' },
  { value: 'MEASUREMENT_TYPE_POUND', title: 'pounds' },
  { value: 'MEASUREMENT_TYPE_GRAM', title: 'grams' },
  { value: 'MEASUREMENT_TYPE_MILLILITER', title: 'milliliters' },
  { value: 'MEASUREMENT_TYPE_LITER', title: 'liters' },
] as const

export const MEASUREMENT_TYPE_TO_STRING: Record<Recipe_MeasurementType, string> = {
  MEASUREMENT_TYPE_UNSPECIFIED: '',
  MEASUREMENT_TYPE_TABLESPOON: 'tablespoons',
  MEASUREMENT_TYPE_TEASPOON: 'teaspoons',
  MEASUREMENT_TYPE_OUNCE: 'ounces',
  MEASUREMENT_TYPE_POUND: 'pounds',
  MEASUREMENT_TYPE_GRAM: 'grams',
  MEASUREMENT_TYPE_MILLILITER: 'milliliters',
  MEASUREMENT_TYPE_LITER: 'liters',
  MEASUREMENT_TYPE_CUP: 'cups',
}
