import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { User, UserSettings } from '@/genapi/api/users/user/v1alpha1'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'
import { userService, authService, circleService, userSettingsService } from '@/api/api'
import type { PermissionLevel  } from '@/genapi/api/types'

export enum AccountType {
  USER = 'user',
  CIRCLE = 'circle',
}

const JWT_STORAGE_KEY = 'jwt'

export const useAuthStore = defineStore('auth', () => {
  const isLoggedIn = ref(false)
  const authInitialized = ref(false)
  let authInitPromise: Promise<void> | null = null
  const userId = ref<number>(0)
  const user = ref<User>({
    name: '',
    username: '',
    givenName: '',
    familyName: '',
    imageUri: '',
    access: undefined,
    bio: '',
  })
  const userSettings = ref<UserSettings>({
    name: '',
    email: '',
  })

  /**
   * Logs the user out by clearing authentication data and removing the JWT from both storages.
   */
  function logOut() {
    console.log('logOut')
    sessionStorage.removeItem(JWT_STORAGE_KEY) // Remove the JWT from sessionStorage
    localStorage.removeItem(JWT_STORAGE_KEY) // Remove the JWT from localStorage
    _clearAuthData()
    isLoggedIn.value = false
  }

  /**
   * loadAuthUser
   */
  async function loadAuthUser() {
    try {
      user.value = await userService.GetUser({
        name: `users/${userId.value}`,
      })
      userSettings.value = await userSettingsService.GetUserSettings({
        name: `users/${userId.value}/settings`,
      })
    } catch (error) {
      console.error('Error:', error)
      throw error
    }
  }

  /**
   * Checks if a JWT exists in localStorage or sessionStorage and sets up authentication data if it does.
   * Clears authentication data if no JWT is found.
   *
   * @private
   */
  async function _checkAuth() {
    const url = new URLSearchParams(window.location.search)
    const value = url.get('token_key')

    try {
    if (value) {
        let res = await authService.ExchangeToken({
          tokenKey: value,
        })
        if (res.token) {
          // Check if rememberMe is set in localStorage
          const rememberMe = localStorage.getItem('rememberMe') === 'true'
          if (rememberMe) {
            localStorage.setItem(JWT_STORAGE_KEY, res.token)
          } else {
            sessionStorage.setItem(JWT_STORAGE_KEY, res.token)
          }
          // Remove rememberMe flag after use
          localStorage.removeItem('rememberMe')
        } else {
          throw new Error('No token returned from auth service')
        }
      }
      
      // Check for JWT in localStorage first, then sessionStorage
      const token = localStorage.getItem(JWT_STORAGE_KEY) || sessionStorage.getItem(JWT_STORAGE_KEY)
      if (token) {
        console.log('JWT found in storage')
        await _setupAuthData()
      } else {
        console.log('No JWT found in storage')
        _clearAuthData()
      }
    } catch (error) {
      console.error('Error:', error)
      _clearAuthData()
    }
      
    authInitialized.value = true
  }

  /**
   * Returns a promise that resolves when auth initialization is complete
   */
  function waitForAuthInit(): Promise<void> {
    if (authInitialized.value) {
      return Promise.resolve()
    }
    if (!authInitPromise) {
      authInitPromise = _checkAuth()
    }
    return authInitPromise
  }

  /**
   * Sets up authentication data for the user
   * Marks the user as logged in.
   *
   * @private
   */
  async function _setupAuthData() {
    try {
      let res = await authService.CheckToken({})
      if (res.userId) {
        userId.value = res.userId
      } else {
        throw new Error('No user id returned from auth service')
      }
      await loadAuthUser()
    } catch (error) {
      console.error('Error:', error)
      throw error
    }
    
    isLoggedIn.value = true
  }

  /**
   * Clears all authentication data.
   *
   * @private
   */
  function _clearAuthData() {
    isLoggedIn.value = false
    user.value = {
      name: '',
      username: '',
      givenName: '',
      familyName: '',
      imageUri: '',
      access: undefined,
      bio: '',
    }
  }

  // Run authentication check on store initialization
  if (!authInitPromise) {
    authInitPromise = _checkAuth()
  }

  return {
    user,
    userSettings,
    isLoggedIn,
    authInitialized,
    waitForAuthInit,
    loadAuthUser,
    logOut,
  }
})
