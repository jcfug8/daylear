import { createRecipeServiceClient } from '@/genapi/api/meals/recipe/v1alpha1'
import { createUserServiceClient } from '@/genapi/api/users/user/v1alpha1'
import { createAuthServiceClient } from './auth'

const API_BASE_URL = 'http://localhost:8080/'

// Generic fetch handler for the generated API client
export const authenticatedFetchHandler = async (
  request: { path: string; method: string; body: string | null },
  _meta?: { service: string; method: string },
) => {
  const token = sessionStorage.getItem('jwt')
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  const res = await fetch(API_BASE_URL + request.path, {
    method: request.method,
    headers,
    body: request.body,
  })
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`)
  }
  if (res.status === 204) return undefined
  return await res.json()
}

export const recipeService = createRecipeServiceClient(authenticatedFetchHandler)
export const userService = createUserServiceClient(authenticatedFetchHandler)
export const authService = createAuthServiceClient(authenticatedFetchHandler)
