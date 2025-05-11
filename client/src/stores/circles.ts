import { ref } from 'vue'
import { defineStore } from 'pinia'
import { circleService, publicCircleService } from '@/api/api'
import type {
  Circle,
  PublicCircle,
  ListCirclesRequest,
  ListCirclesResponse,
  ListPublicCirclesRequest,
  ListPublicCirclesResponse,
} from '@/genapi/api/circles/circle/v1alpha1'

export const useCirclesStore = defineStore('circles', () => {
  const circles = ref<Circle[]>([])
  const publicCircles = ref<PublicCircle[]>([])
  const circle = ref<Circle | undefined>()
  const publicCircle = ref<PublicCircle | undefined>()

  async function loadCircles(parent: string = 'users/1', filter?: string) {
    try {
      const request: ListCirclesRequest = {
        parent,
        pageSize: undefined,
        pageToken: undefined,
        filter: filter,
      }
      const response = (await circleService.ListCircles(request)) as ListCirclesResponse
      circles.value = response.circles ?? []
    } catch (error) {
      console.error('Failed to load circles:', error)
      circles.value = []
    }
  }

  async function loadPublicCircles(filter?: string) {
    try {
      const request: ListPublicCirclesRequest = {
        pageSize: undefined,
        pageToken: undefined,
        filter: filter,
      }
      const response = (await publicCircleService.ListPublicCircles(request)) as ListPublicCirclesResponse
      publicCircles.value = response.publicCircles ?? []
    } catch (error) {
      console.error('Failed to load circles:', error)
      circles.value = []
    }
  }

  async function loadCircle(circleName: string) {
    try {
      const result = await circleService.GetCircle({ name: circleName })
      circle.value = result
    } catch (error) {
      console.error('Failed to load circle:', error)
      circle.value = undefined
    }
  }

  async function loadPublicCircle(circleName: string) {
    try {
      const result = await publicCircleService.GetPublicCircle({ name: circleName })
      publicCircle.value = result
    } catch (error) {
      console.error('Failed to load public circle:', error)
      publicCircle.value = undefined
    }
  }

  function initEmptyCircle() {
    circle.value = {
      name: undefined,
      title: '',
      isPublic: false,
    }
  }

  async function createCircle() {
    if (!circle.value) {
      throw new Error('No circle to create')
    }
    if (circle.value.name) {
      throw new Error('Circle already has a name and cannot be created')
    }
    try {
      const created = await circleService.CreateCircle({
        parent: 'users/1',
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

  return {
    loadCircle,
    loadCircles,
    loadPublicCircles,
    loadPublicCircle,
    initEmptyCircle,
    createCircle,
    updateCircle,
    circles,
    circle,
    publicCircles,
    publicCircle,
  }
}) 