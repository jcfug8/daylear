import { createRecipeServiceClient, createRecipeAccessServiceClient } from '@/genapi/api/meals/recipe/v1alpha1'
import { createUserServiceClient, createUserAccessServiceClient, createUserSettingsServiceClient } from '@/genapi/api/users/user/v1alpha1'
import { createCircleServiceClient, createCircleAccessServiceClient } from '@/genapi/api/circles/circle/v1alpha1'
import { createAuthServiceClient } from './auth'
import { createFileServiceClient } from './files'
import { AccountType, useAuthStore } from '@/stores/auth'
import { API_BASE_URL } from '@/constants/api'

// Generic fetch handler for the generated API client
export const authenticatedFetchHandler = function(contentType: string) {
  return async function <T = string | null | FormData>(
    request: { path: string; method: string; body: T, signal?: AbortSignal, responseType?: string },
    _meta?: { service: string; method: string },
  ) {
    const authStore = useAuthStore()

    const token = sessionStorage.getItem('jwt')
    const headers: Record<string, string> = {}

    if (contentType) {
      headers['Content-Type'] = contentType
    }

    if (authStore.activeAccountType === AccountType.CIRCLE) {
      headers["X-Daylear-Circle"] = authStore.activeAccount?.name ? authStore.activeAccount?.name : ""
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

export const recipeService = createRecipeServiceClient(authenticatedFetchHandler('application/json'))
export const recipeAccessService = createRecipeAccessServiceClient(authenticatedFetchHandler('application/json'))
export const userService = createUserServiceClient(authenticatedFetchHandler('application/json'))
export const userSettingsService = createUserSettingsServiceClient(authenticatedFetchHandler('application/json'))
export const userAccessService = createUserAccessServiceClient(authenticatedFetchHandler('application/json'))
export const authService = createAuthServiceClient(authenticatedFetchHandler('application/json'))
export const circleService = createCircleServiceClient(authenticatedFetchHandler('application/json'))
export const circleAccessService = createCircleAccessServiceClient(authenticatedFetchHandler('application/json'))
export const fileService = createFileServiceClient(authenticatedFetchHandler(''))
