<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(defineProps<{
  zoomText?: string
  showFullscreen?: boolean
}>(), {
  zoomText: '100%',
  showFullscreen: true,
})

const emit = defineEmits<{
  zoomIn: []
  zoomOut: []
  reset: []
  fullscreen: []
}>()

const isFullscreen = ref(false)

function toggleFullscreen() {
  isFullscreen.value = !isFullscreen.value
  emit('fullscreen')
}
</script>

<template>
  <div class="canvas-toolbar">
    <button type="button" class="canvas-toolbar__btn" title="缩小" @click="emit('zoomOut')">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="5" y1="12" x2="19" y2="12" /></svg>
    </button>
    <span class="canvas-toolbar__zoom">{{ zoomText }}</span>
    <button type="button" class="canvas-toolbar__btn" title="放大" @click="emit('zoomIn')">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" /></svg>
    </button>
    <div class="canvas-toolbar__divider" />
    <button type="button" class="canvas-toolbar__btn" title="重置" @click="emit('reset')">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" /><path d="M3 3v5h5" /></svg>
    </button>
    <div v-if="showFullscreen" class="canvas-toolbar__divider" />
    <button v-if="showFullscreen" type="button" class="canvas-toolbar__btn" :class="{ 'is-active': isFullscreen }" title="全屏" @click="toggleFullscreen">
      <svg v-if="!isFullscreen" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 3 21 3 21 9" /><polyline points="9 21 3 21 3 15" /><line x1="21" y1="3" x2="14" y2="10" /><line x1="3" y1="21" x2="10" y2="14" /></svg>
      <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="4 14 10 14 10 20" /><polyline points="20 10 14 10 14 4" /><line x1="14" y1="10" x2="21" y2="3" /><line x1="3" y1="21" x2="10" y2="14" /></svg>
    </button>
  </div>
</template>

<style scoped>
.canvas-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: var(--radius-md);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-border);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.canvas-toolbar__btn {
  display: grid;
  width: 28px;
  height: 28px;
  place-items: center;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background-color var(--transition-fast), color var(--transition-fast);
}

.canvas-toolbar__btn svg {
  width: 16px;
  height: 16px;
}

.canvas-toolbar__btn:hover {
  background: var(--color-bg-hover);
  color: var(--color-text-primary);
}

.canvas-toolbar__btn.is-active {
  background: var(--el-color-primary);
  color: #FFFFFF;
}

.canvas-toolbar__zoom {
  min-width: 44px;
  padding: 0 4px;
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 600;
  text-align: center;
  user-select: none;
}

.canvas-toolbar__divider {
  width: 1px;
  height: 18px;
  margin: 0 2px;
  background: var(--color-border);
}
</style>
