import { ref } from 'vue'
import { defineStore } from 'pinia'
import { circleService, circleAccessService } from '@/api/api'
import type {
  Circle,
  apitypes_VisibilityLevel,
  apitypes_PermissionLevel,
  apitypes_AccessState,
} from '@/genapi/api/circles/circle/v1alpha1'
import { useAuthStore } from '@/stores/auth'

export const useCirclesStore = defineStore('circles', () => {
  const circles = ref<Circle[]>([])
  const circle = ref<Circle | undefined>()
  const myCircles = ref<Circle[]>([])
  const sharedAcceptedCircles = ref<Circle[]>([])
  const sharedPendingCircles = ref<Circle[]>([])
  const publicCircles = ref<Circle[]>([])

  async function loadCircle(circleName: string) {
    try {
      const result = await circleService.GetCircle({ name: circleName })
      circle.value = result
    } catch (error) {
      console.error('Failed to load circle:', error)
      circle.value = undefined
    }
  }

  async function loadCircles(filter?: string) {
    try {
      const result = await circleService.ListCircles({ 
        filter: filter || undefined,
        pageSize: 50,
        pageToken: undefined
      })
      circles.value = result.circles || []
    } catch (error) {
      console.error('Failed to load public circles:', error)
      circles.value = []
    }
  }

  function initEmptyCircle() {
    circle.value = {
      name: undefined,
      title: '',
      visibility: 'VISIBILITY_LEVEL_PRIVATE' as apitypes_VisibilityLevel,
      permission: 'PERMISSION_LEVEL_UNSPECIFIED' as apitypes_PermissionLevel,
      state: 'ACCESS_STATE_UNSPECIFIED' as apitypes_AccessState,
    }
  }

  async function createCircle(parent: string) {
    if (!circle.value) {
      throw new Error('No circle to create')
    }
    if (circle.value.name) {
      throw new Error('Circle already has a name and cannot be created')
    }
    try {
      const created = await circleService.CreateCircle({
        circle: circle.value,
        circleId: crypto.randomUUID(),
      })
      circle.value = created
      return created
    } catch (error) {
      console.error('Failed to create circle:', error)
      throw error
    }
  }

  async function updateCircle() {
    if (!circle.value) {
      throw new Error('No circle to update')
    }
    if (!circle.value.name) {
      throw new Error('Circle must have a name to be updated')
    }
    try {
      const updated = await circleService.UpdateCircle({
        circle: circle.value,
        updateMask: undefined, // Optionally specify fields to update
      })
      circle.value = updated
      return updated
    } catch (error) {
      console.error('Failed to update circle:', error)
      throw error
    }
  }

  // Load my circles (admin permission)
  async function loadMyCircles() {
    await loadCircles('permission = 300')
  }

  // Load shared circles (accepted or pending)
  async function loadSharedCircles(state?: number) {
    let filter = 'permission = 100 OR permission = 200'
    if (state) {
      filter += ` AND state = ${state}`
    }
    await loadCircles(filter)
  }

  // Load public circles
  async function loadPublicCircles() {
    await loadCircles('visibility = 1')
  }

  return {
    loadCircle,
    loadCircles,
    initEmptyCircle,
    createCircle,
    updateCircle,
    circles,
    circle,
    myCircles,
    sharedAcceptedCircles,
    sharedPendingCircles,
    publicCircles,
    loadMyCircles,
    loadSharedCircles,
    loadPublicCircles,
  }
}) 