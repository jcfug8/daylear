import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useRecipeFormStore = defineStore('recipeForm', () => {
  const activeTab = ref('general')

  function setActiveTab(tab: string) {
    activeTab.value = tab
  }

  return {
    activeTab,
    setActiveTab,
  }
})
