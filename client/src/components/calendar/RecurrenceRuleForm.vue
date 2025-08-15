<template>
  <div class="recurrence-rule-form">
    <!-- Repeat Toggle -->
    <div class="d-flex align-center mb-4">
      <v-checkbox
        v-model="isRepeating"
        label="Repeat"
        hide-details
        density="compact"
        @update:model-value="onRepeatToggle"
      />
    </div>

    <!-- Recurrence Options (only shown when repeating) -->
    <div v-if="isRepeating" class="recurrence-options">
      <!-- Frequency Selection -->
      <div class="mb-4">
        <v-select
          v-model="frequency"
          :items="frequencyOptions"
          item-title="label"
          item-value="value"
          label="Repeat"
          density="compact"
          hide-details
          @update:model-value="updateRule"
        />
      </div>

      <!-- Interval -->
      <div v-if="frequency !== 'never'" class="mb-4">
        <v-text-field
          v-model="interval"
          type="number"
          min="1"
          max="99"
          :label="`Every ${interval} ${getIntervalLabel()}`"
          density="compact"
          hide-details
          @update:model-value="updateRule"
        />
      </div>

      <!-- Weekly Options -->
      <div v-if="frequency === 'WEEKLY'" class="mb-4">
        <div class="text-caption text-medium-emphasis mb-2">Repeat on</div>
        <div class="d-flex flex-wrap gap-1">
          <v-chip
            v-for="day in weekDays"
            :key="day.value"
            :color="selectedWeekDays.includes(day.value) ? 'primary' : 'default'"
            variant="outlined"
            size="small"
            @click="toggleWeekDay(day.value)"
          >
            {{ day.label }}
          </v-chip>
        </div>
      </div>

      <!-- Monthly Options -->
      <div v-if="frequency === 'MONTHLY'" class="mb-4">
        <v-radio-group
          v-model="monthlyType"
          density="compact"
          hide-details
          @update:model-value="updateRule"
        >
          <v-radio value="day" label="On the same day of the month" />
          <v-radio value="weekday" label="On the same weekday of the month" />
        </v-radio-group>
        
        <!-- Pattern options for monthly -->
        <div v-if="monthlyType === 'weekday'" class="mt-2">
          <div class="d-flex gap-2 align-center">
            <v-select
              v-model="monthlyPatternOccurrence"
              :items="monthlyPatternOccurrences"
              item-title="label"
              item-value="value"
              density="compact"
              hide-details
              style="max-width: 120px"
              @update:model-value="updateRule"
            />
            <v-select
              v-model="monthlyPatternWeekday"
              :items="weekDays"
              item-title="label"
              item-value="value"
              density="compact"
              hide-details
              style="max-width: 120px"
              @update:model-value="updateRule"
            />
            <span class="text-caption">of each month</span>
          </div>
        </div>
      </div>

      <!-- Yearly Options -->
      <div v-if="frequency === 'YEARLY'" class="mb-4">
        <v-radio-group
          v-model="yearlyType"
          density="compact"
          hide-details
          @update:model-value="updateRule"
        >
          <v-radio value="date" label="On the same date each year" />
          <v-radio value="pattern" label="On a specific pattern" />
        </v-radio-group>
        
        <!-- Pattern options for yearly -->
        <div v-if="yearlyType === 'pattern'" class="mt-2">
          <div class="d-flex gap-2 align-center">
            <v-select
              v-model="yearlyPatternOccurrence"
              :items="monthlyPatternOccurrences"
              item-title="label"
              item-value="value"
              density="compact"
              hide-details
              style="max-width: 120px"
              @update:model-value="updateRule"
            />
            <v-select
              v-model="yearlyPatternWeekday"
              :items="weekDays"
              item-title="label"
              item-value="value"
              density="compact"
              hide-details
              style="max-width: 120px"
              @update:model-value="updateRule"
            />
            <v-select
              v-model="yearlyPatternMonth"
              :items="months"
              item-title="label"
              item-value="value"
              density="compact"
              hide-details
              style="max-width: 120px"
              @update:model-value="updateRule"
            />
          </div>
        </div>
      </div>

      <!-- End Options -->
      <div class="mb-4">
        <div class="text-caption text-medium-emphasis mb-2">End</div>
        <v-radio-group
          v-model="endType"
          density="compact"
          hide-details
          @update:model-value="onEndTypeChange"
        >
          <v-radio value="never" label="Never" />
          <v-radio value="after" label="After X occurrences" />
          <v-radio value="until" label="Until date" />
        </v-radio-group>
      </div>

      <!-- End Options Details -->
      <div v-if="endType === 'after'" class="mb-4">
        <v-text-field
          v-model="occurrenceCount"
          type="number"
          min="1"
          max="999"
          label="Number of occurrences"
          density="compact"
          hide-details
          @update:model-value="updateRule"
        />
      </div>

      <div v-if="endType === 'until'" class="mb-4">
        <v-text-field
          v-model="endDate"
          type="date"
          label="End date"
          density="compact"
          hide-details
          @update:model-value="updateRule"
        />
      </div>

      <!-- Preview of Generated Rule -->
      <div v-if="isRepeating" class="mt-4">
        <div class="text-caption text-medium-emphasis mb-2">Generated Rule:</div>
        <v-chip
          :color="isValidRule ? 'success' : 'error'"
          variant="outlined"
          size="small"
          class="font-mono"
        >
          {{ recurrenceRule }}
        </v-chip>
        
        <!-- Human readable preview -->
        <div v-if="humanReadableRule" class="mt-2 text-caption text-medium-emphasis">
          {{ humanReadableRule }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, nextTick } from 'vue'
import { RRule, type ByWeekday } from 'rrule'

const props = withDefaults(defineProps<{
  modelValue: string
  disabled?: boolean
}>(), {
  disabled: false,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

// Local state
const isRepeating = ref(false)
const frequency = ref('DAILY')
const interval = ref(1)
const selectedWeekDays = ref<string[]>(['MO'])
const monthlyType = ref('day')
const monthlyPatternOccurrence = ref(1)
const monthlyPatternWeekday = ref('MO')
const yearlyType = ref('date')
const yearlyPatternOccurrence = ref(1)
const yearlyPatternWeekday = ref('MO')
const yearlyPatternMonth = ref(1)
const endType = ref('never')
const occurrenceCount = ref(10)
const endDate = ref('')
const recurrenceRule = ref('')
const isInternalUpdate = ref(false)

// Computed properties
const isValidRule = computed(() => {
  if (!recurrenceRule.value) return true
  try {
    RRule.fromString(recurrenceRule.value)
    return true
  } catch {
    return false
  }
})

const humanReadableRule = computed(() => {
  if (!recurrenceRule.value) return ''
  try {
    return RRule.fromString(recurrenceRule.value).toText()
  } catch {
    return ''
  }
})

// Options for dropdowns
const frequencyOptions = [
  { label: 'Daily', value: 'DAILY' },
  { label: 'Weekly', value: 'WEEKLY' },
  { label: 'Monthly', value: 'MONTHLY' },
  { label: 'Yearly', value: 'YEARLY' },
]

const weekDays = [
  { label: 'Mon', value: 'MO' },
  { label: 'Tue', value: 'TU' },
  { label: 'Wed', value: 'WE' },
  { label: 'Thu', value: 'TH' },
  { label: 'Fri', value: 'FR' },
  { label: 'Sat', value: 'SA' },
  { label: 'Sun', value: 'SU' },
]

const months = [
  { label: 'January', value: 1 },
  { label: 'February', value: 2 },
  { label: 'March', value: 3 },
  { label: 'April', value: 4 },
  { label: 'May', value: 5 },
  { label: 'June', value: 6 },
  { label: 'July', value: 7 },
  { label: 'August', value: 8 },
  { label: 'September', value: 9 },
  { label: 'October', value: 10 },
  { label: 'November', value: 11 },
  { label: 'December', value: 12 },
]

const monthlyPatternOccurrences = [
  { label: '1st', value: 1 },
  { label: '2nd', value: 2 },
  { label: '3rd', value: 3 },
  { label: '4th', value: 4 },
  { label: '5th', value: 5 },
  { label: 'Last', value: -1 },
]

// Helper functions
function getIntervalLabel(): string {
  switch (frequency.value) {
    case 'DAILY': return interval.value === 1 ? 'day' : 'days'
    case 'WEEKLY': return interval.value === 1 ? 'week' : 'weeks'
    case 'MONTHLY': return interval.value === 1 ? 'month' : 'months'
    case 'YEARLY': return interval.value === 1 ? 'year' : 'years'
    default: return 'period'
  }
}

function getFrequencyValue(freq: string): number {
  switch (freq) {
    case 'DAILY': return RRule.DAILY
    case 'WEEKLY': return RRule.WEEKLY
    case 'MONTHLY': return RRule.MONTHLY
    case 'YEARLY': return RRule.YEARLY
    default: return RRule.DAILY
  }
}

// Watch for prop changes and update local state
watch(() => props.modelValue, (newValue) => {
  if (isInternalUpdate.value) {
    // Skip parsing if this is an internal update
    return
  }
  if (newValue) {
    parseRecurrenceRule(newValue)
  } else {
    resetForm()
  }
}, { immediate: true })

function resetForm() {
  isRepeating.value = false
  frequency.value = 'DAILY'
  interval.value = 1
  selectedWeekDays.value = ['MO']
  monthlyType.value = 'day'
  monthlyPatternOccurrence.value = 1
  monthlyPatternWeekday.value = 'MO'
  yearlyType.value = 'date'
  yearlyPatternOccurrence.value = 1
  yearlyPatternWeekday.value = 'MO'
  yearlyPatternMonth.value = 1
  endType.value = 'never'
  occurrenceCount.value = 10
  endDate.value = ''
  recurrenceRule.value = ''
}

function onRepeatToggle(value: boolean | null) {
  isRepeating.value = value || false
  if (!value) {
    resetForm()
    emit('update:modelValue', '')
  } else {
    updateRule()
  }
}

function onEndTypeChange(value: string | null) {
  endType.value = value || 'never'
  updateRule()
}

function toggleWeekDay(day: string) {
  const index = selectedWeekDays.value.indexOf(day)
  if (index > -1) {
    selectedWeekDays.value.splice(index, 1)
  } else {
    selectedWeekDays.value.push(day)
  }
  
  // Ensure at least one day is selected
  if (selectedWeekDays.value.length === 0) {
    selectedWeekDays.value = ['MO']
  }
  
  updateRule()
}

watch(selectedWeekDays, (newValue) => {
  console.log('selectedWeekDays', newValue)
}, { deep: true })

function updateRule() {
  if (!isRepeating.value) {
    recurrenceRule.value = ''
    emit('update:modelValue', '')
    return
  }

     try {
     const options: {
       freq: number
       interval: number
       byweekday?: number[]
       bysetpos?: number[]
       bymonth?: number[]
       count?: number
       until?: Date
     } = {
       freq: getFrequencyValue(frequency.value),
       interval: interval.value,
     }

    // Add weekly options
    if (frequency.value === 'WEEKLY' && selectedWeekDays.value.length > 0) {
      options.byweekday = selectedWeekDays.value.map(day => {
        const dayMap: { [key: string]: number } = { MO: 0, TU: 1, WE: 2, TH: 3, FR: 4, SA: 5, SU: 6 }
        return dayMap[day] || 0
      })
    }

    // Add monthly options
    if (frequency.value === 'MONTHLY' && monthlyType.value === 'weekday') {
      options.bysetpos = [monthlyPatternOccurrence.value]
      const dayMap: { [key: string]: number } = { MO: 0, TU: 1, WE: 2, TH: 3, FR: 4, SA: 5, SU: 6 }
      options.byweekday = [dayMap[monthlyPatternWeekday.value] || 0]
    }

    // Add yearly options
    if (frequency.value === 'YEARLY' && yearlyType.value === 'pattern') {
      options.bysetpos = [yearlyPatternOccurrence.value]
      const dayMap: { [key: string]: number } = { MO: 0, TU: 1, WE: 2, TH: 3, FR: 4, SA: 5, SU: 6 }
      options.byweekday = [dayMap[yearlyPatternWeekday.value] || 0]
      options.bymonth = [yearlyPatternMonth.value]
    }

    // Add end conditions
    if (endType.value === 'after') {
      options.count = occurrenceCount.value
    } else if (endType.value === 'until' && endDate.value) {
      options.until = new Date(endDate.value)
    }

    const rrule = new RRule(options)
    recurrenceRule.value = rrule.toString()
    
    // Set flag to prevent watch from triggering
    isInternalUpdate.value = true
    emit('update:modelValue', recurrenceRule.value)
    // Reset flag after emit
    nextTick(() => {
      isInternalUpdate.value = false
    })
  } catch (error) {
    console.error('Error creating recurrence rule:', error)
    recurrenceRule.value = ''
    emit('update:modelValue', '')
  }
}

function parseRecurrenceRule(rule: string) {
  if (!rule) {
    resetForm()
    return
  }

  try {
    const rrule = RRule.fromString(rule)
    const options = rrule.origOptions
    
    isRepeating.value = true
    frequency.value = Object.keys(RRule).find(key => RRule[key as keyof typeof RRule] === options.freq) || 'DAILY'
    interval.value = options.interval || 1

    // Parse weekly options
    if (options.byweekday && Array.isArray(options.byweekday) && options.byweekday.length > 0) {
      const dayNames = ['MO', 'TU', 'WE', 'TH', 'FR', 'SA', 'SU']
      const uniqueDays = new Set<string>()
      options.byweekday.forEach((day: ByWeekday) => {
        if (typeof day === 'number') {
          const dayName = dayNames[day] || 'MO'
          uniqueDays.add(dayName)
        }
      })
      selectedWeekDays.value = Array.from(uniqueDays)
    }

    // Parse monthly options
    if (options.bysetpos && Array.isArray(options.bysetpos) && options.bysetpos.length > 0) {
      monthlyType.value = 'weekday'
      monthlyPatternOccurrence.value = options.bysetpos[0]
      if (options.byweekday && Array.isArray(options.byweekday) && options.byweekday.length > 0) {
        const dayNames = ['MO', 'TU', 'WE', 'TH', 'FR', 'SA', 'SU']
        const firstDay = options.byweekday[0]
        if (typeof firstDay === 'number') {
          monthlyPatternWeekday.value = dayNames[firstDay] || 'MO'
        }
      }
    }

    // Parse yearly options
    if (options.bymonth && Array.isArray(options.bymonth) && options.bymonth.length > 0) {
      yearlyType.value = 'pattern'
      yearlyPatternMonth.value = options.bymonth[0]
      if (options.bysetpos && Array.isArray(options.bysetpos) && options.bysetpos.length > 0) {
        yearlyPatternOccurrence.value = options.bysetpos[0]
      }
      if (options.byweekday && Array.isArray(options.byweekday) && options.byweekday.length > 0) {
        const dayNames = ['MO', 'TU', 'WE', 'TH', 'FR', 'SA', 'SU']
        const firstDay = options.byweekday[0]
        if (typeof firstDay === 'number') {
          yearlyPatternWeekday.value = dayNames[firstDay] || 'MO'
        }
      }
    }

    // Parse end conditions
    if (options.count) {
      endType.value = 'after'
      occurrenceCount.value = options.count
    } else if (options.until) {
      endType.value = 'until'
      endDate.value = options.until.toISOString().split('T')[0]
    } else {
      endType.value = 'never'
    }

    recurrenceRule.value = rule
  } catch (error) {
    console.error('Error parsing recurrence rule:', error)
    resetForm()
  }
}
</script>

<style scoped>
.recurrence-rule-form {
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 8px;
  padding: 16px;
}

.recurrence-options {
  margin-top: 8px;
  padding-top: 16px;
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
</style>
