import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const isLoggedIn = ref(false)
  function logIn() {
    console.log('logIn')
    isLoggedIn.value = true
  }
  function logOut() {
    console.log('logOut')
    isLoggedIn.value = false
  }
  return { isLoggedIn, logIn, logOut }
})
