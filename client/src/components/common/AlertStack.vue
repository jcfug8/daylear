<template>
  <div class="alert-stack">
    <div
      v-for="alert in alerts"
      :key="alert.id"
      class="alert"
      :class="alert.severity"
    >
      <button class="close" @click="remove(alert.id)">&times;</button>
      <span class="message">{{ alert.message }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAlertStore } from '@/stores/alerts';
import { storeToRefs } from 'pinia';
import { onMounted } from 'vue';

const alertStore = useAlertStore();
const { alerts } = storeToRefs(alertStore);
const remove = alertStore.removeAlert;

onMounted(() => {
//   alertStore.addAlert('This is a test alert', 'info')
})
</script>

<style scoped>
.alert-stack {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 10000;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.alert {
  min-width: 250px;
  padding: 1rem 2.5rem 1rem 1rem;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  color: #fff;
  position: relative;
  font-size: 1rem;
  animation: fadein 0.2s;
}
.alert.info { background: #2196f3; }
.alert.success { background: #43a047; }
.alert.warning { background: #ffa000; }
.alert.error { background: #e53935; }
.close {
  position: absolute;
  top: 0.5rem;
  right: 0.7rem;
  background: none;
  border: none;
  color: #fff;
  font-size: 1.2rem;
  cursor: pointer;
  line-height: 1;
}
.message {
  display: block;
  word-break: break-word;
}
@keyframes fadein {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style> 