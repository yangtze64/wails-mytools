<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, ref, type Component } from 'vue'
import { getAppInfo } from '../api/app'
import ToolTabBar from '../components/ToolTabBar.vue'
import { groupTools, toolRegistry, type ToolKey } from '../config/toolRegistry'
import { useToolTabs } from '../composables/useToolTabs'
import type { AppInfo } from '../types/app'
import BookmarkManagerView from './tools/BookmarkManagerView.vue'
import DevToolsView from './tools/DevToolsView.vue'
import JsonToolView from './tools/JsonToolView.vue'
import MarkdownEditorView from './tools/MarkdownEditorView.vue'
import QrDecodeView from './tools/QrDecodeView.vue'
import QrCodeView from './tools/QrCodeView.vue'
import RandomStringView from './tools/RandomStringView.vue'
import SecretManagerView from './tools/SecretManagerView.vue'
import TextDiffView from './tools/TextDiffView.vue'
import TextDedupeView from './tools/TextDedupeView.vue'
import TimestampView from './tools/TimestampView.vue'

const appInfo = ref<AppInfo | null>(null)
const keyword = ref('')
const MindMapView = defineAsyncComponent(() => import('./tools/MindMapView.vue'))
const MermaidEditorView = defineAsyncComponent(() => import('./tools/MermaidEditorView.vue'))
const {
  activeTool,
  activeToolInfo,
  openTabItems,
  openToolTab,
  closeToolTab,
  closeOtherToolTabs,
} = useToolTabs()

const toolComponents: Record<ToolKey, Component> = {
  'text-dedupe': TextDedupeView,
  'random-string': RandomStringView,
  'text-diff': TextDiffView,
  'json-tool': JsonToolView,
  'markdown-editor': MarkdownEditorView,
  'mermaid-editor': MermaidEditorView,
  timestamp: TimestampView,
  bookmarks: BookmarkManagerView,
  'secret-manager': SecretManagerView,
  devtools: DevToolsView,
  'qr-code': QrCodeView,
  'qr-decode': QrDecodeView,
  'mind-map': MindMapView,
}

const activeToolComponent = computed(() => toolComponents[activeToolInfo.value.key])
const groupedTools = computed(() => {
  const normalizedKeyword = keyword.value.trim().toLocaleLowerCase()
  const filteredTools = normalizedKeyword
    ? toolRegistry.filter((tool) => `${tool.title} ${tool.description}`.toLocaleLowerCase().includes(normalizedKeyword))
    : toolRegistry

  return groupTools(filteredTools)
})

onMounted(async () => {
  appInfo.value = await getAppInfo()
})
</script>

<template>
  <el-container class="app-shell">
    <el-aside class="sidebar" width="260px">
      <el-input v-model="keyword" class="tool-search" clearable placeholder="搜索工具" />

      <el-menu class="tool-menu" :default-active="activeTool" @select="(key: string) => openToolTab(key as ToolKey)">
        <template v-for="group in groupedTools" :key="group.category">
          <div class="menu-group">{{ group.category }}</div>
          <el-menu-item
            v-for="tool in group.tools"
            :key="tool.key"
            :index="tool.key"
          >
            <div class="menu-item">
              <strong>{{ tool.title }}</strong>
              <span>{{ tool.description }}</span>
            </div>
          </el-menu-item>
        </template>
      </el-menu>

      <div class="brand brand--footer">
        <div class="brand__logo">
          <img src="/logo.png" alt="MyTools" />
        </div>
        <div class="brand__text">
          <strong>MyTools</strong>
          <span>Version {{ appInfo?.version ?? '0.1.0' }}</span>
        </div>
      </div>
    </el-aside>

    <el-main class="workspace">
      <ToolTabBar
        :tabs="openTabItems"
        :active-tool="activeTool"
        @select="openToolTab"
        @close="closeToolTab"
        @close-others="closeOtherToolTabs"
      />

      <div class="workspace__body">
        <KeepAlive>
          <component :is="activeToolComponent" :key="activeToolInfo.key" />
        </KeepAlive>

        <footer class="app-footer">
          <span>{{ appInfo?.name ?? 'MyTools' }} {{ appInfo?.version ?? '' }}</span>
          <span>{{ appInfo ? `${appInfo.os}/${appInfo.arch}` : '' }}</span>
        </footer>
      </div>
    </el-main>
  </el-container>
</template>
