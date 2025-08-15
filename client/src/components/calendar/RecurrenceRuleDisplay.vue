<template>
  <div v-if="recurrenceRule" class="recurrence-rule-display">
    <div class="d-flex align-center mb-2">
      <v-icon icon="mdi-repeat" size="small" class="mr-2" />
      <span class="text-caption text-medium-emphasis">Repeats</span>
    </div>
    <div class="text-body-2">{{ humanReadableRule }}</div>
    
    <!-- Show next occurrences if available -->
    <div v-if="nextOccurrences.length > 0" class="mt-3">
      <div class="text-caption text-medium-emphasis mb-1">Next occurrences:</div>
      <div class="text-body-2">
        <div v-for="(date, index) in nextOccurrences" :key="index" class="mb-1">
          {{ formatDate(date) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RRule } from 'rrule'

const props = defineProps<{
  recurrenceRule: string
  startDate?: string
}>()

const nextOccurrences = ref<Date[]>([])

const humanReadableRule = computed(() => {
  if (!props.recurrenceRule) return ''
  try {
    return RRule.fromString(props.recurrenceRule).toText()
  } catch {
    return 'Invalid recurrence rule'
  }
})

function formatDate(date: Date): string {
  return date.toLocaleDateString(undefined, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

function getNextOccurrences(rule: string, startDate: Date, count: number = 5): Date[] {
  try {
    const rrule = RRule.fromString(rule)
    rrule.options.dtstart = startDate
    
    // Get the next occurrences within the next year
    const endDate = new Date(startDate.getTime() + 365 * 24 * 60 * 60 * 1000)
    return rrule.between(startDate, endDate, true, (date, i) => i < count)
  } catch {
    return []
  }
}

onMounted(() => {
  if (props.recurrenceRule && props.startDate) {
    try {
      const startDate = new Date(props.startDate)
      if (!isNaN(startDate.getTime())) {
        nextOccurrences.value = getNextOccurrences(props.recurrenceRule, startDate, 3)
      }
    } catch {
      // If we can't parse the start date, just skip showing next occurrences
    }
  }
})
</script>

<style scoped>
.recurrence-rule-display {
  padding: 8px 12px;
  background-color: rgba(var(--v-theme-surface-variant), 0.1);
  border-radius: 6px;
  border-left: 3px solid var(--v-theme-primary);
}
</style>
