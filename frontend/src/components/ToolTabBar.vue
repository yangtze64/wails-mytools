<script setup lang="ts">
import type { ToolKey, ToolMenuItem } from '../config/toolRegistry'

defineProps<{
  tabs: ToolMenuItem[]
  activeTool: ToolKey
}>()

const emit = defineEmits<{
  select: [key: ToolKey]
  close: [key: ToolKey]
  closeOthers: [key: ToolKey]
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
  </div>
</template>
