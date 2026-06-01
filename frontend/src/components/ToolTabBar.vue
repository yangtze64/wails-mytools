<script setup lang="ts">
import type { ToolKey, ToolMenuItem } from '../config/toolRegistry'

defineProps<{
  tabs: ToolMenuItem[]
  activeTool: ToolKey
  theme: 'dark' | 'light'
}>()

const emit = defineEmits<{
  select: [key: ToolKey]
  close: [key: ToolKey]
  closeOthers: [key: ToolKey]
  toggleTheme: []
}>()
</script>

<template>
  <div class="tool-tabs" role="tablist" aria-label="已打开工具">
    <div class="tool-tabs__scroller">
      <div
        v-for="tab in tabs"
        :key="tab.key"
        class="tool-tab"
        :class="{ 'is-active': tab.key === activeTool }"
        role="tab"
        :aria-selected="tab.key === activeTool"
        :title="tab.description"
        @dblclick="emit('closeOthers', tab.key)"
      >
        <button
          class="tool-tab__trigger"
          type="button"
          @click="emit('select', tab.key)"
          @click.middle.prevent="emit('close', tab.key)"
        >
          <span class="tool-tab__title">{{ tab.title }}</span>
        </button>
        <button
          v-if="tabs.length > 1"
          class="tool-tab__close"
          type="button"
          aria-label="关闭标签"
          @click.stop="emit('close', tab.key)"
        >
          x
        </button>
      </div>
    </div>

    <div class="tool-tabs__actions">
      <button
        class="theme-toggle"
        type="button"
        :title="theme === 'dark' ? '切换到亮色模式' : '切换到暗色模式'"
        @click="emit('toggleTheme')"
      >
        <svg v-if="theme === 'dark'" class="theme-toggle__icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="5" />
          <line x1="12" y1="1" x2="12" y2="3" />
          <line x1="12" y1="21" x2="12" y2="23" />
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
          <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
          <line x1="1" y1="12" x2="3" y2="12" />
          <line x1="21" y1="12" x2="23" y2="12" />
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
          <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
        </svg>
        <svg v-else class="theme-toggle__icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
        </svg>
      </button>
    </div>
  </div>
</template>
