import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/genapi/api/users/user/v1alpha1'
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'

export const useAuthStore = defineStore('auth', () => {
  const isLoggedIn = ref(false)
  const user = ref<User>({
    name: '',
    email: '',
    givenName: '',
    familyName: '',
  })
  const circles = ref<Circle[]>([])
  const activeAccount = ref<User | Circle | undefined>()

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
        let res = await fetch(`http://localhost:8080/auth/token/${value}`)
        let data = await res.json()
        sessionStorage.setItem('jwt', data.token)
      } catch (error) {
        console.error('Error:', error)
      }
    }

    const token = sessionStorage.getItem('jwt') // Retrieve the JWT from sessionStorage
    if (token) {
      console.log('JWT found in sessionStorage')
      _setupAuthData()
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
  function _setupAuthData() {
    user.value = {
      name: 'users/1',
      email: 'jcfug8@gmail.com',
      givenName: 'Jace',
      familyName: 'Fugate',
    }
    circles.value = [
      {
        name: 'circles/1',
        title: 'Circle 1',
      },
      {
        name: 'circles/2',
        title: 'Circle 2',
      },
    ]
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
    isLoggedIn,
    logOut,
    setActiveAccount,
  }
})
