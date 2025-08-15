import { RRule } from 'rrule'

/**
 * Validate if a recurrence rule is syntactically correct
 */
export function validateRecurrenceRule(rule: string): boolean {
  try {
    RRule.fromString(rule)
    return true
  } catch {
    return false
  }
}

/**
 * Generate the next few occurrences of a recurring event
 */
export function getNextOccurrences(rule: string, startDate: Date, count: number = 5): Date[] {
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
