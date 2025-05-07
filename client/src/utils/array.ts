export function moveArrayItem<T>(
  array: T[] | undefined,
  index: number,
  direction: 'up' | 'down',
): void {
  if (!array) return

  const newIndex = direction === 'up' ? index - 1 : index + 1
  if (newIndex < 0 || newIndex >= array.length) return

  const temp = array[index]
  array[index] = array[newIndex]
  array[newIndex] = temp
}

export function moveNestedArrayItem<T>(
  parentArray: T[] | undefined,
  parentIndex: number,
  childArrayKey: keyof T,
  childIndex: number,
  direction: 'up' | 'down',
): void {
  if (!parentArray?.[parentIndex]) return

  const childArray = parentArray[parentIndex][childArrayKey] as unknown as any[]
  if (!childArray) return

  moveArrayItem(childArray, childIndex, direction)
}
