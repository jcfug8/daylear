import { createRecipeServiceClient, createRecipeAccessServiceClient } from '@/genapi/api/meals/recipe/v1alpha1'
import { createListServiceClient, createListAccessServiceClient, createListItemServiceClient, createListItemCompletionServiceClient } from '@/genapi/api/lists/list/v1alpha1'
import { createUserServiceClient, createUserAccessServiceClient, createUserSettingsServiceClient, createAccessKeyServiceClient } from '@/genapi/api/users/user/v1alpha1'
import { createCircleServiceClient, createCircleAccessServiceClient } from '@/genapi/api/circles/circle/v1alpha1'
import { createCalendarServiceClient, createCalendarAccessServiceClient, createEventServiceClient, createEventRecipeServiceClient } from '@/genapi/api/calendars/calendar/v1alpha1'
import { createAuthServiceClient } from './auth'
import { createFileServiceClient } from './files'
import { API_BASE_URL } from '@/constants/api'

// Type for API spec
export type ApiSpec = {
  name: string;
  url: string;
};

// Generic fetch handler for the generated API client
export const authenticatedFetchHandler = function(contentType: string) {
  return async function <T = string | null | FormData>(
    request: { path: string; method: string; body: T, signal?: AbortSignal, responseType?: string },
    _meta?: { service: string; method: string },
  ) {

    const token = localStorage.getItem('jwt') || sessionStorage.getItem('jwt')
    const headers: Record<string, string> = {}

    if (contentType) {
      headers['Content-Type'] = contentType
    }

    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    const fetchOptions: any = {
      method: request.method,
      headers,
      body: request.body as any,
    }
    if (request.signal) {
      fetchOptions.signal = request.signal
    }
    const res = await fetch(API_BASE_URL + request.path, fetchOptions)
    if (!res.ok) {
      const error = await res.json()
      let message = `API error: ${res.status} ${res.statusText}`
      if (error.message) {
        message = `\n${error.message}`
      }
      throw new Error(message)
    }
    if (res.status === 204) return undefined
    if (request.responseType === "blob") {
      return await res.blob()
    }
    return await res.json()
  }
}

// Function to fetch API specs
export async function fetchApiSpecs(): Promise<ApiSpec[]> {
  const res = await fetch(API_BASE_URL + 'openapi/specs.json')
  if (!res.ok) {
    throw new Error(`Failed to fetch API specs: ${res.status} ${res.statusText}`)
  }
  return await res.json()
}

export const recipeService = createRecipeServiceClient(authenticatedFetchHandler('application/json'))
export const recipeAccessService = createRecipeAccessServiceClient(authenticatedFetchHandler('application/json'))
export const listService = createListServiceClient(authenticatedFetchHandler('application/json'))
export const listAccessService = createListAccessServiceClient(authenticatedFetchHandler('application/json'))
export const listItemService = createListItemServiceClient(authenticatedFetchHandler('application/json'))
export const listItemCompletionService = createListItemCompletionServiceClient(authenticatedFetchHandler('application/json'))
export const userService = createUserServiceClient(authenticatedFetchHandler('application/json'))
export const userSettingsService = createUserSettingsServiceClient(authenticatedFetchHandler('application/json'))
export const userAccessService = createUserAccessServiceClient(authenticatedFetchHandler('application/json'))
export const accessKeyService = createAccessKeyServiceClient(authenticatedFetchHandler('application/json'))
export const authService = createAuthServiceClient(authenticatedFetchHandler('application/json'))
export const circleService = createCircleServiceClient(authenticatedFetchHandler('application/json'))
export const circleAccessService = createCircleAccessServiceClient(authenticatedFetchHandler('application/json'))
export const fileService = createFileServiceClient(authenticatedFetchHandler(''))
export const calendarService = createCalendarServiceClient(authenticatedFetchHandler('application/json'))
export const calendarAccessService = createCalendarAccessServiceClient(authenticatedFetchHandler('application/json'))
export const eventService = createEventServiceClient(authenticatedFetchHandler('application/json'))
export const eventRecipeService = createEventRecipeServiceClient(authenticatedFetchHandler('application/json'))