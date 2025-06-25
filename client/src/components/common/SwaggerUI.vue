<template>
  <div ref="swaggerContainer" style="width: 100%; min-height: 80vh;"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import 'swagger-ui-dist/swagger-ui.css';

const props = defineProps<{
  url: string;
}>();

const swaggerContainer = ref<HTMLElement | null>(null);
let ui: any = null;

async function loadSwaggerUI() {
  const SwaggerUIBundle = (await import('swagger-ui-dist/swagger-ui-es-bundle.js')).default;
  if (swaggerContainer.value) {
    // Clear previous UI if any
    swaggerContainer.value.innerHTML = '';
    ui = SwaggerUIBundle({
      url: props.url,
      domNode: swaggerContainer.value,
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis],
      layout: 'BaseLayout',
    });
  }
}

onMounted(() => {
  loadSwaggerUI();
});

watch(() => props.url, () => {
  loadSwaggerUI();
});
</script> 