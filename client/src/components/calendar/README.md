# Calendar Components

This directory contains calendar-related components for the Daylear application.

## Components

### RecurrenceRuleForm

A comprehensive form component for creating and editing recurrence rules for calendar events. It provides a user-friendly interface similar to Apple Calendar's repeat options.

**Features:**
- Toggle to enable/disable recurrence
- Predefined frequency options (Daily, Weekly, Bi-weekly, Monthly, Yearly)
- Custom frequency options with detailed configuration
- End conditions (Never, After X occurrences, Until date)
- Real-time validation of generated RRULE strings
- Preview of generated recurrence rule

**Usage:**
```vue
<RecurrenceRuleForm
  v-model="recurrenceData"
  :disabled="false"
/>
```

**Props:**
- `modelValue`: RecurrenceRuleData object
- `disabled`: Boolean to disable the form

**Events:**
- `update:modelValue`: Emitted when recurrence data changes

### RecurrenceRuleDisplay

A display component that shows recurrence rules in a human-readable format with additional information about upcoming occurrences.

**Features:**
- Human-readable description of recurrence rules
- Display of next few occurrences (if start date is provided)
- Clean, card-based design with visual indicators

**Usage:**
```vue
<RecurrenceRuleDisplay
  :recurrence-rule="event.recurrenceRule"
  :start-date="event.startTime"
/>
```

**Props:**
- `recurrenceRule`: String containing the RRULE
- `startDate`: Optional start date for calculating next occurrences

## Data Types

### RecurrenceRuleData

```typescript
interface RecurrenceRuleData {
  isRepeating: boolean
  frequency: string
  customFrequency: string
  dailyInterval: number
  weeklyInterval: number
  selectedWeekDays: string[]
  monthlyInterval: number
  monthlyType: string
  monthlyPatternOccurrence: string
  monthlyPatternWeekday: string
  yearlyInterval: number
  yearlyType: string
  yearlyPatternOccurrence: string
  yearlyPatternWeekday: string
  yearlyPatternMonth: string
  endType: string
  occurrenceCount: number
  endDate: string
  recurrenceRule: string
}
```

## Utilities

### recurrence.ts

Utility functions for working with recurrence rules:

- `parseRecurrenceRule(rule: string)`: Parse RRULE strings into structured data
- `validateRecurrenceRule(rule: string)`: Validate RRULE syntax
- `getHumanReadableDescription(rule: string)`: Convert RRULE to human-readable text
- `getNextOccurrences(rule: string, startDate: Date, count: number)`: Calculate next occurrences

## RRULE Format

The components generate and parse RRULE strings following the iCalendar specification:

- `FREQ=DAILY`: Daily recurrence
- `FREQ=WEEKLY;INTERVAL=2;BYDAY=MO,WE,FR`: Every 2 weeks on Monday, Wednesday, Friday
- `FREQ=MONTHLY;BYDAY=1MO`: First Monday of each month
- `FREQ=YEARLY;BYDAY=1TH;BYMONTH=JAN`: First Thursday of January each year

## Integration

These components are integrated into:

- `EventForm.vue`: For creating/editing events with recurrence
- `EventCreateDialog.vue`: For creating new events
- `EventEditDialog.vue`: For editing existing events
- `EventViewDialog.vue`: For displaying event details with recurrence information

## Dependencies

- Vue 3 with Composition API
- Vuetify 3 for UI components
- `rrule` library for RRULE parsing and generation
