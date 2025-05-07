<template>
  <div
    class="button-column d-flex flex-column gap-1"
    :class="{ 'mt-4': !size || size === 'small' }"
  >
    <v-btn
      icon="mdi-arrow-up"
      :size="size || 'small'"
      @click="$emit('move-up')"
      variant="text"
      :style="{ visibility: showUp ? 'visible' : 'hidden' }"
    ></v-btn>
    <v-btn
      icon="mdi-arrow-down"
      :size="size || 'small'"
      @click="$emit('move-down')"
      variant="text"
      :style="{ visibility: showDown ? 'visible' : 'hidden' }"
    ></v-btn>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  showUp: boolean
  showDown: boolean
  size?: 'x-small' | 'small'
}>()

defineEmits<{
  (e: 'move-up'): void
  (e: 'move-down'): void
}>()
</script>

<style scoped>
.button-column {
  width: 10px;
  align-items: center;
}
</style>
