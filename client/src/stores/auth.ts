import { ref, computed } from 'vue'
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

  function checkAuth() {
    const token = sessionStorage.getItem('jwt') // Retrieve the JWT from sessionStorage
    if (token) {
      console.log('JWT found in sessionStorage')
      setupAuthData()
      isLoggedIn.value = true
    } else {
      console.log('No JWT found in sessionStorage')
      isLoggedIn.value = false
    }
  }

  function setAccounts() {
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
  }

  function clearAccounts() {
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

  function setupAuthData() {
    setAccounts()
    isLoggedIn.value = true
  }

  function logIn() {
    console.log('logIn')
    sessionStorage.setItem('jwt', 'fake-jwt') // Save a fake JWT to sessionStorage
    setupAuthData()
  }
  function logOut() {
    console.log('logOut')
    sessionStorage.removeItem('jwt') // Remove the JWT from sessionStorage
    clearAccounts()
    isLoggedIn.value = false
  }
  function setActiveAccount(account: User | Circle) {
    console.log('setActiveAccount', account)
    activeAccount.value = account
  }

  checkAuth()

  return {
    user,
    circles,
    activeAccount,
    setActiveAccount,
    isLoggedIn,
    logIn,
    logOut,
  }
})
