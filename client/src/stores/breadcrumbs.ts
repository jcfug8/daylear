import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { RouteLocationRaw } from 'vue-router'

// the crumb used in the breadcrumb
export type Breadcrumb = {
  // the string used to identify the breadcrumb
  title: string
  // the href or route object navigated to when the breadcrumb is clicked
  to: RouteLocationRaw
}

export const useBreadcrumbStore = defineStore('breadCrumbs', () => {
  const breadcrumbs = ref<Breadcrumb[]>([])

  function setBreadcrumbs(crumbs: Breadcrumb[]) {
    breadcrumbs.value = crumbs
  }

  return {
    breadcrumbs,
    setBreadcrumbs,
  }
})
