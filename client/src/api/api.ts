import { createRecipeServiceClient } from '@/genapi/api/meals/recipe/v1alpha1'
import { createPublicUserServiceClient, createUserServiceClient } from '@/genapi/api/users/user/v1alpha1'
import { createCircleServiceClient, createPublicCircleServiceClient } from '@/genapi/api/circles/circle/v1alpha1'
import { createAuthServiceClient } from './auth'
import { createFileServiceClient } from './files'

const API_BASE_URL = 'http://localhost:8080/'

// Generic fetch handler for the generated API client
export const authenticatedFetchHandler = function(contentType: string = 'application/json') {
  return async function <T = string | null | FormData>(
    request: { path: string; method: string; body: T },
    _meta?: { service: string; method: string },
  ) {
    const token = sessionStorage.getItem('jwt')
    const headers: Record<string, string> = {
      'Content-Type': contentType,
    }
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    const res = await fetch(API_BASE_URL + request.path, {
      method: request.method,
      headers,
      body: request.body as any,
    })
    if (!res.ok) {
      throw new Error(`API error: ${res.status} ${res.statusText}`)
    }
    if (res.status === 204) return undefined
    return await res.json()
  }
}

export const recipeService = createRecipeServiceClient(authenticatedFetchHandler())
export const userService = createUserServiceClient(authenticatedFetchHandler())
export const publicUserService = createPublicUserServiceClient(authenticatedFetchHandler())
export const authService = createAuthServiceClient(authenticatedFetchHandler())
export const circleService = createCircleServiceClient(authenticatedFetchHandler())
export const publicCircleService = createPublicCircleServiceClient(authenticatedFetchHandler())
export const fileService = createFileServiceClient(authenticatedFetchHandler("application/octet-stream"))
