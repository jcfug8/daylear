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
const ACTIVE_ACCOUNT_STORAGE_KEY = 'daylear_active_account_name'

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
    visibility: undefined,
    imageUri: '',
    access: undefined,
    bio: '',
  })
  const userSettings = ref<UserSettings>({
    name: '',
    email: '',
  })
  const circles = ref<Circle[]>([])
  const activeAccount = ref<User | Circle>()
  const activeAccountName = computed(() => {
    if (activeAccountType.value === AccountType.USER) {
      return `` // the user name is not needed when getting resrouces
    } else {
      return activeAccount.value?.name ?? ''
    }
  })
  const activeAccountTitle = computed(() => {
    if (activeAccount.value && 'username' in activeAccount.value) {
      return activeAccount.value.username
    } else {
      return activeAccount.value?.title ?? ''
    }
  })
  const activeAccountType = computed(() => {
    if (activeAccount.value && 'username' in activeAccount.value) {
      return AccountType.USER
    } else if (activeAccount.value) {
      return AccountType.CIRCLE
    }
  })
  const activeAccountPermissionLevel = computed<PermissionLevel>(() => {
    if (activeAccount.value && 'username' in activeAccount.value) {
      return "PERMISSION_LEVEL_ADMIN"
    } else if (activeAccount.value) {
      return activeAccount.value.circleAccess?.permissionLevel as PermissionLevel ?? "PERMISSION_LEVEL_UNSPECIFIED"
    }
    return "PERMISSION_LEVEL_UNSPECIFIED"
  })
  /**
   * Logs the user out by clearing authentication data and removing the JWT from both storages.
   */
  function logOut() {
    console.log('logOut')
    sessionStorage.removeItem(JWT_STORAGE_KEY) // Remove the JWT from sessionStorage
    localStorage.removeItem(JWT_STORAGE_KEY) // Remove the JWT from localStorage
    sessionStorage.removeItem(ACTIVE_ACCOUNT_STORAGE_KEY) // Remove the active account from sessionStorage
    localStorage.removeItem(ACTIVE_ACCOUNT_STORAGE_KEY) // Remove the active account from localStorage (if ever used)
    _clearAuthData()
    isLoggedIn.value = false
  }

  /**
   * Sets the active account to the specified user or circle.
   *
   * @param {User | Circle} account - The account to set as active.
   */
  function setActiveAccount(account: User | Circle) {
    console.log('setActiveAccount', account)
    activeAccount.value = account
    _saveActiveAccountToStorage(account)
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
   * loadAuthCircles
   */
  async function loadAuthCircles() {
    try {
      let res = await circleService.ListCircles({
        pageSize: 100,
        pageToken: '',
        filter: 'permission > 1',
      })
      circles.value = res.circles ?? []
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
   * Sets up authentication data for the user and circles.
   * Sets the active account to the user.
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
  
      await loadAuthCircles()
    } catch (error) {
      console.error('Error:', error)
      throw error
    }

    // Try to restore the active account from session storage
    const restoredAccount = _restoreActiveAccountFromStorage()
    if (restoredAccount) {
      activeAccount.value = restoredAccount
    } else {
      // Fall back to user account if no valid stored account
      activeAccount.value = user.value
    }
    
    isLoggedIn.value = true
  }

  /**
   * Clears all authentication data, including user, circles, and active account.
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
      visibility: undefined,
      imageUri: '',
      access: undefined,
      bio: '',
    }
    circles.value = []
    activeAccount.value = undefined
  }

  /**
   * Saves the active account name to session storage.
   *
   * @private
   * @param {User | Circle} account - The account to save.
   */
  function _saveActiveAccountToStorage(account: User | Circle) {
    try {
      if (account?.name) {
        sessionStorage.setItem(ACTIVE_ACCOUNT_STORAGE_KEY, account.name)
      }
    } catch (error) {
      console.warn('Failed to save active account to session storage:', error)
    }
  }

  /**
   * Gets the stored active account name from session storage.
   *
   * @private
   * @returns {string | null} The stored account name or null if not found.
   */
  function _getStoredActiveAccountName(): string | null {
    try {
      return sessionStorage.getItem(ACTIVE_ACCOUNT_STORAGE_KEY)
    } catch (error) {
      console.warn('Failed to retrieve active account from session storage:', error)
      return null
    }
  }

  /**
   * Finds an account by name in the user or circles data.
   *
   * @private
   * @param {string} name - The account name to find.
   * @returns {User | Circle | null} The found account or null.
   */
  function _findAccountByName(name: string): User | Circle | null {
    // Check if it matches the user
    if (user.value?.name === name) {
      return user.value
    }
    
    // Check if it matches any circle
    const circle = circles.value.find(c => c.name === name)
    if (circle) {
      return circle
    }
    
    return null
  }

  /**
   * Restores the active account from session storage if the user still has access.
   *
   * @private
   * @returns {User | Circle | null} The restored account or null if not found/accessible.
   */
  function _restoreActiveAccountFromStorage(): User | Circle | null {
    const storedAccountName = _getStoredActiveAccountName()
    if (!storedAccountName) {
      return null
    }

    const account = _findAccountByName(storedAccountName)
    if (account) {
      console.log('Restored active account from session storage:', storedAccountName)
      return account
    } else {
      console.log('Stored active account no longer accessible, clearing from storage:', storedAccountName)
      // Clear the invalid stored account
      try {
        sessionStorage.removeItem(ACTIVE_ACCOUNT_STORAGE_KEY)
      } catch (error) {
        console.warn('Failed to clear invalid active account from session storage:', error)
      }
      return null
    }
  }

  // Run authentication check on store initialization
  if (!authInitPromise) {
    authInitPromise = _checkAuth()
  }

  return {
    user,
    userSettings,
    circles,
    activeAccount,
    activeAccountName,
    activeAccountPermissionLevel,
    activeAccountType,
    activeAccountTitle,
    isLoggedIn,
    authInitialized,
    waitForAuthInit,
    loadAuthUser,
    loadAuthCircles,
    logOut,
    setActiveAccount,
  }
})
