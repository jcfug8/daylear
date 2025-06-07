import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/genapi/api/users/user/v1alpha1'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'
import { userService, authService, circleService } from '@/api/api'

export enum AccountType {
  USER = 'user',
  CIRCLE = 'circle',
}

export const useAuthStore = defineStore('auth', () => {
  const isLoggedIn = ref(false)
  const userId = ref<number>(0)
  const user = ref<User>({
    name: '',
    email: '',
    username: '',
    givenName: '',
    familyName: '',
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
    } else {
      return AccountType.CIRCLE
    }
  })
  /**
   * Logs the user out by clearing authentication data and removing the JWT from sessionStorage.
   */
  function logOut() {
    console.log('logOut')
    sessionStorage.removeItem('jwt') // Remove the JWT from sessionStorage
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
  }

  /**
   * loadAuthUser
   */
  async function loadAuthUser() {
    try {
      user.value = await userService.GetUser({
        name: `users/${userId.value}`,
      })
    } catch (error) {
      console.error('Error:', error)
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
        filter: '',
      })
      circles.value = res.circles ?? []
    } catch (error) {
      console.error('Error:', error)
    }
  }

  /**
   * updateAuthUser
   */
  async function updateAuthUser(editUser: User) {
    try {
      user.value = await userService.UpdateUser({
        user: editUser,
        updateMask: undefined,
      })
    } catch (error) {
      console.error('Error:', error)
    }
  }

  /**
   * Checks if a JWT exists in sessionStorage and sets up authentication data if it does.
   * Clears authentication data if no JWT is found.
   *
   * @private
   */
  async function _checkAuth() {
    const url = new URLSearchParams(window.location.search)
    const value = url.get('token_key')

    if (value) {
      try {
        let res = await authService.ExchangeToken({
          tokenKey: value,
        })
        if (res.token) {
          sessionStorage.setItem('jwt', res.token)
        } else {
          throw new Error('No token returned from auth service')
        }
      } catch (error) {
        console.error('Error:', error)
      }
    }

    const token = sessionStorage.getItem('jwt') // Retrieve the JWT from sessionStorage
    if (token) {
      console.log('JWT found in sessionStorage')
      await _setupAuthData()
    } else {
      console.log('No JWT found in sessionStorage')
      _clearAuthData()
    }
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
    } catch (error) {
      console.error('Error:', error)
    }

    await loadAuthUser()

    await loadAuthCircles()

    if (activeAccount.value === undefined) {
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
      email: '',
      username: '',
      givenName: '',
      familyName: '',
    }
    circles.value = []
    activeAccount.value = undefined
  }

  // Run authentication check on store initialization
  _checkAuth()

  return {
    user,
    circles,
    activeAccount,
    activeAccountName,
    activeAccountType,
    activeAccountTitle,
    isLoggedIn,
    loadAuthUser,
    loadAuthCircles,
    updateAuthUser,
    logOut,
    setActiveAccount,
  }
})
