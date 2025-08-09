import { ref } from 'vue'
import { defineStore } from 'pinia'
import { calendarService } from '@/api/api'
import type {
  Calendar,
  ListCalendarsRequest,
  ListCalendarsResponse,
} from '@/genapi/api/calendars/calendar/v1alpha1'
import { calendarAccessService } from '@/api/api'

export const useCalendarsStore = defineStore('calendars', () => {
  // const calendars = ref<Calendar[]>([])
  const calendar = ref<Calendar | undefined>()
  const myCalendars = ref<Calendar[]>([])
  const sharedAcceptedCalendars = ref<Calendar[]>([])
  const sharedPendingCalendars = ref<Calendar[]>([])
  const publicCalendars = ref<Calendar[]>([])

  async function loadCalendars(parent: string, filter?: string): Promise<Calendar[]> {
    try {
      const request = {
        parent,
        pageSize: undefined,
        pageToken: undefined,
        filter,
      }
      const response = (await calendarService.ListCalendars(
        request as ListCalendarsRequest,
      )) as ListCalendarsResponse
      return response.calendars ?? []
    } catch (error) {
      console.error('Failed to load calendars:', error)
      return []
    }
  }

  // Load my calendars (calendars where I have admin permission)
  async function loadMyCalendars(parent: string) {
    const calendars = await loadCalendars(parent, 'state = 200')
    myCalendars.value = calendars
  }

  // Load shared calendars (calendars shared with me - read or write permission)
  async function loadPendingCalendars(parent: string) {
    const calendars = await loadCalendars(parent, 'state = 100')
    sharedPendingCalendars.value = calendars
  }

  // Load public calendars (calendars with public visibility)
  async function loadPublicCalendars(parent: string) {
    const calendars = await loadCalendars(parent, 'visibility = 1')
    publicCalendars.value = calendars
  }

  async function loadCalendar(calendarName: string) {
    try {
      const result = await calendarService.GetCalendar({ name: calendarName })
      calendar.value = result
    } catch (error) {
      console.error('Failed to load calendar:', error)
      calendar.value = undefined
    }
  }

  function initEmptyCalendar() {
    calendar.value = {
      name: '',
      title: '',
      visibility: 'VISIBILITY_LEVEL_PRIVATE',
      description: '',
      calendarAccess: undefined,
    }
  }

  async function createCalendar(parent: string) {
    if (!calendar.value) {
      throw new Error('No calendar to create')
    }
    if (calendar.value.name) {
      throw new Error('Calendar already has a name and cannot be created')
    }
    console.log('Creating calendar with data:', calendar.value)
    console.log('Parent path:', parent)
    try {
      const created = await calendarService.CreateCalendar({
        parent,
        calendar: calendar.value,
      })
      calendar.value = created
      return created
    } catch (error) {
      console.error('Failed to create calendar:', error)
      throw error
    }
  }

  async function updateCalendar() {
    if (!calendar.value) {
      throw new Error('No calendar to update')
    }
    if (!calendar.value.name) {
      throw new Error('Calendar must have a name to be updated')
    }
    try {
      const updated = await calendarService.UpdateCalendar({
        calendar: calendar.value,
        updateMask: undefined, // Optionally specify fields to update
      })
      calendar.value = updated
      return updated
    } catch (error) {
      console.error('Failed to update calendar:', error)
      throw error
    }
  }

  async function acceptCalendar(accessName: string) {
    try {
      await calendarAccessService.AcceptAccess({ name: accessName })
    } catch (error) {
      console.error('Failed to accept calendar access:', error)
      throw error
    }
  }

  async function deleteCalendarAccess(accessName: string) {
    try {
      await calendarAccessService.DeleteAccess({ name: accessName })
    } catch (error) {
      console.error('Failed to decline calendar access:', error)
      throw error
    }
  }

  return {
    loadMyCalendars,
    loadPendingCalendars,
    loadPublicCalendars,
    loadCalendar,
    initEmptyCalendar,
    createCalendar,
    updateCalendar,
    acceptCalendar,
    deleteCalendarAccess,
    calendar,
    myCalendars,
    sharedAcceptedCalendars,
    sharedPendingCalendars,
    publicCalendars,
  }
})
