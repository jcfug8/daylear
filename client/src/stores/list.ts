import { ref } from 'vue'
import { defineStore } from 'pinia'
import { listService, listItemService } from '@/api/api'
import type {
  List,
  ListListsRequest,
  ListListsResponse,
  apitypes_VisibilityLevel,
  ListItem,
  CreateListItemRequest,
  ListListItemsRequest,
  ListListItemsResponse,
  UpdateListItemRequest,
  DeleteListItemRequest,
  GetListItemRequest,
} from '@/genapi/api/lists/list/v1alpha1'
import { listAccessService } from '@/api/api'

export const useListStore = defineStore('list', () => {
  // const lists = ref<List[]>([])
  const list = ref<List | undefined>()
  const myLists = ref<List[]>([])
  const sharedAcceptedLists = ref<List[]>([])
  const sharedPendingLists = ref<List[]>([])
  const publicLists = ref<List[]>([])
  
  // List items state
  const listItems = ref<Record<string, ListItem[]>>({})

  async function loadLists(parent: string, filter?: string): Promise<List[]> {
    try {
      const request = {
        parent,
        pageSize: undefined,
        pageToken: undefined,
        filter,
      }
      const response = (await listService.ListLists(
        request as ListListsRequest,
      )) as ListListsResponse
      return response.lists ?? []
    } catch (error) {
      console.error('Failed to load lists:', error)
      return []
    }
  }

  // Load my lists (lists where I have admin permission)
  async function loadMyLists(parent: string) {
    myLists.value = []
    const lists = await loadLists(parent, 'state = 200 OR favorited = true')
    myLists.value = lists
  }

  // Load shared lists (lists shared with me - read or write permission)
  async function loadPendingLists(parent: string) {
    sharedPendingLists.value = []
    const lists = await loadLists(parent, 'state = 100')
    sharedPendingLists.value = lists
  }

  // Load public lists (lists with public visibility)
  async function loadPublicLists(parent: string) {
    publicLists.value = []
    const lists = await loadLists(parent, 'visibility = 1')
    publicLists.value = lists
  }

  async function loadList(listName: string) {
    try {
      list.value = undefined
      const result = await listService.GetList({ name: listName })
      list.value = result
    } catch (error) {
      console.error('Failed to load list:', error)
      list.value = undefined
    }
  }

  function initEmptyList() {
    list.value = {
      name: undefined,
      title: undefined,
      description: undefined,
      showCompleted: false,
      visibility: 'VISIBILITY_LEVEL_PUBLIC' as apitypes_VisibilityLevel,
      listAccess: undefined,
      createTime: undefined,
      updateTime: undefined,
      favorited: false,
      sections: [],
    }
  }

  async function createList(parent: string) {
    if (!list.value) {
      throw new Error('No list to create')
    }
    if (list.value.name) {
      throw new Error('List already has a name and cannot be created')
    }
    console.log('Creating list with data:', list.value)
    console.log('Parent path:', parent)

    cleanList()

    try {
      const created = await listService.CreateList({
        parent,
        list: list.value,
        listId: crypto.randomUUID(),
      })
      list.value = created
      return created
    } catch (error) {
      console.error('Failed to create list:', error)
      throw error
    }
  }

  async function updateList() {
    if (!list.value) {
      throw new Error('No list to update')
    }
    if (!list.value.name) {
      throw new Error('List must have a name to be updated')
    }

    cleanList()

    try {
      const updated = await listService.UpdateList({
        list: list.value,
        updateMask: undefined, // Optionally specify fields to update
      })
      list.value = updated
      return updated
    } catch (error) {
      console.error('Failed to update list:', error)
      throw error
    }
  }

  function cleanList() {
    if (!list.value) {
      return
    }
    // Ensure sections array exists
    if (!list.value.sections) {
      list.value.sections = []
    }
    // Clean up sections - ensure each section has required fields
    list.value.sections.forEach((section) => {
      if (!section.title) {
        section.title = 'Untitled Section'
      }
    })
  }

  async function acceptList(accessName: string) {
    try {
      await listAccessService.AcceptListAccess({ name: accessName })
    } catch (error) {
      console.error('Failed to accept list access:', error)
      throw error
    }
  }

  async function deleteListAccess(accessName: string) {
    try {
      await listAccessService.DeleteAccess({ name: accessName })
    } catch (error) {
      console.error('Failed to decline list access:', error)
      throw error
    }
  }

  async function favoriteList(listName: string) {
    try {
      await listService.FavoriteList({ name: listName })
    } catch (error) {
      console.error('Failed to favorite list:', error)
      throw error
    }
  }

  async function unfavoriteList(listName: string) {
    try {
      await listService.UnfavoriteList({ name: listName })
    } catch (error) {
      console.error('Failed to unfavorite list:', error)
      throw error
    }
  }

  // *** List Item Methods ***

  // Get list items for a specific list
  function getListItems(listName: string): ListItem[] {
    return listItems.value[listName] || []
  }

  // Load list items from the server
  async function loadListItems(listName: string): Promise<ListItem[]> {
    try {
      const request: ListListItemsRequest = {
        parent: listName,
        pageSize: undefined,
        pageToken: undefined,
        filter: undefined,
      }
      
      const response: ListListItemsResponse = await listItemService.ListListItems(request)
      const items = response.listItems || []
      
      // Store items in state
      listItems.value[listName] = items
      
      return items
    } catch (error) {
      console.error('Failed to load list items:', error)
      return []
    }
  }

  // Create a new list item
  async function createListItem(
    listName: string, 
    title: string, 
    listSection?: string, 
    points: number = 0,
    recurrenceRule?: string
  ): Promise<ListItem> {
    try {
      const request: CreateListItemRequest = {
        parent: listName,
        listItem: {
          name: undefined, // Will be set by server
          title,
          listSection,
          points,
          recurrenceRule,
          createTime: undefined, // OUTPUT_ONLY
          updateTime: undefined, // OUTPUT_ONLY
        },
      }
      
      const newItem = await listItemService.CreateListItem(request)
      
      // Add to local state
      if (!listItems.value[listName]) {
        listItems.value[listName] = []
      }
      listItems.value[listName].push(newItem)
      
      return newItem
    } catch (error) {
      console.error('Failed to create list item:', error)
      throw error
    }
  }

  // Update an existing list item
  async function updateListItem(
    itemName: string, 
    updates: Partial<Pick<ListItem, 'title' | 'listSection' | 'points' | 'recurrenceRule'>>
  ): Promise<ListItem> {
    try {
      const request: UpdateListItemRequest = {
        listItem: {
          name: itemName,
          title: updates.title,
          listSection: updates.listSection,
          points: updates.points,
          recurrenceRule: updates.recurrenceRule,
          createTime: undefined, // OUTPUT_ONLY
          updateTime: undefined, // OUTPUT_ONLY
        },
        updateMask: Object.keys(updates).join(','),
      }
      
      const updatedItem = await listItemService.UpdateListItem(request)
      
      // Update local state
      const listName = itemName.split('/').slice(0, -1).join('/')
      if (listItems.value[listName]) {
        const index = listItems.value[listName].findIndex(item => item.name === itemName)
        if (index !== -1) {
          listItems.value[listName][index] = updatedItem
        }
      }
      
      return updatedItem
    } catch (error) {
      console.error('Failed to update list item:', error)
      throw error
    }
  }

  // Delete a list item
  async function deleteListItem(itemName: string): Promise<void> {
    try {
      const request: DeleteListItemRequest = {
        name: itemName,
      }
      
      await listItemService.DeleteListItem(request)
      
      // Remove from local state
      const listName = itemName.split('/').slice(0, -1).join('/')
      if (listItems.value[listName]) {
        listItems.value[listName] = listItems.value[listName].filter(item => item.name !== itemName)
      }
    } catch (error) {
      console.error('Failed to delete list item:', error)
      throw error
    }
  }

  // Get a specific list item
  async function getListItem(itemName: string): Promise<ListItem | undefined> {
    try {
      const request: GetListItemRequest = {
        name: itemName,
      }
      
      return await listItemService.GetListItem(request)
    } catch (error) {
      console.error('Failed to get list item:', error)
      return undefined
    }
  }

  // Clear list items for a specific list (useful when switching lists)
  function clearListItems(listName: string) {
    delete listItems.value[listName]
  }

  return {
    loadMyLists,
    loadPendingLists,
    loadPublicLists,
    loadList,
    initEmptyList,
    createList,
    updateList,
    acceptList,
    deleteListAccess,
    favoriteList,
    unfavoriteList,
    // List item methods
    getListItems,
    loadListItems,
    createListItem,
    updateListItem,
    deleteListItem,
    getListItem,
    clearListItems,
    // State
    list,
    myLists,
    sharedAcceptedLists,
    sharedPendingLists,
    publicLists,
    listItems,
  }
})
