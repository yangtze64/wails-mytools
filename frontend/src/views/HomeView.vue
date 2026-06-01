<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, ref, type Component } from 'vue'
import { getAppInfo } from '../api/app'
import ToolTabBar from '../components/ToolTabBar.vue'
import { groupTools, toolRegistry, type ToolKey } from '../config/toolRegistry'
import { useToolTabs } from '../composables/useToolTabs'
import { useTheme } from '../composables/useTheme'
import type { AppInfo } from '../types/app'
import {
  Clock,
  CollectionTag,
  Connection,
  CopyDocument,
  DataLine,
  Document,
  EditPen,
  Grid,
  Key,
  Link,
  Lock,
  MagicStick,
  Monitor,
  Picture,
  Share,
  Sort,
  Stamp,
  Timer,
  View,
} from '@element-plus/icons-vue'
import BookmarkManagerView from './tools/BookmarkManagerView.vue'
import DevToolsView from './tools/DevToolsView.vue'
import ImageCompressView from './tools/ImageCompressView.vue'
import JsonToolView from './tools/JsonToolView.vue'
import MarkdownEditorView from './tools/MarkdownEditorView.vue'
import QrDecodeView from './tools/QrDecodeView.vue'
import QrCodeView from './tools/QrCodeView.vue'
import RandomStringView from './tools/RandomStringView.vue'
import SecretManagerView from './tools/SecretManagerView.vue'
import TextDiffView from './tools/TextDiffView.vue'
import TextDedupeView from './tools/TextDedupeView.vue'
import TimestampView from './tools/TimestampView.vue'
import Base64ToolView from './tools/Base64ToolView.vue'
import MarkmapView from './tools/MarkmapView.vue'
import CronGeneratorView from './tools/CronGeneratorView.vue'
import JwtParserView from './tools/JwtParserView.vue'
import UrlToolView from './tools/UrlToolView.vue'

const iconMap: Record<string, Component> = {
  CopyDocument,
  MagicStick,
  Sort,
  Document,
  EditPen,
  Share,
  Clock,
  CollectionTag,
  Lock,
  Monitor,
  Grid,
  View,
  Picture,
  Key,
  Connection,
  DataLine,
  Timer,
  Stamp,
  Link,
}

const appInfo = ref<AppInfo | null>(null)
const keyword = ref('')
const sidebarCollapsed = ref(false)
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
const { currentTheme, toggleTheme } = useTheme()

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
  'markmap-tool': MarkmapView,
  'image-compress': ImageCompressView,
  'base64-tool': Base64ToolView,
  'cron-generator': CronGeneratorView,
  'jwt-parser': JwtParserView,
  'url-tool': UrlToolView,
}

const activeToolComponent = computed(() => toolComponents[activeToolInfo.value.key])
const groupedTools = computed(() => {
  const normalizedKeyword = keyword.value.trim().toLocaleLowerCase()
  const filteredTools = normalizedKeyword
    ? toolRegistry.filter((tool) => `${tool.title} ${tool.description}`.toLocaleLowerCase().includes(normalizedKeyword))
    : toolRegistry

  return groupTools(filteredTools)
})

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

onMounted(async () => {
  appInfo.value = await getAppInfo()
})
</script>

<template>
  <el-container class="app-shell">
    <el-aside class="sidebar" :class="{ 'is-collapsed': sidebarCollapsed }">
      <el-input v-model="keyword" class="tool-search" clearable placeholder="搜索工具" />

      <el-menu class="tool-menu" :default-active="activeTool" :collapse="sidebarCollapsed" @select="(key: string) => openToolTab(key as ToolKey)">
        <template v-for="group in groupedTools" :key="group.category">
          <div class="menu-group" v-show="!sidebarCollapsed">{{ group.category }}</div>
          <el-menu-item
            v-for="tool in group.tools"
            :key="tool.key"
            :index="tool.key"
          >
            <el-icon class="menu-item__icon"><component :is="iconMap[tool.icon]" /></el-icon>
            <template #title>
              <div class="menu-item__text">
                <strong>{{ tool.title }}</strong>
                <span>{{ tool.description }}</span>
              </div>
            </template>
          </el-menu-item>
        </template>
      </el-menu>

      <div class="sidebar-bottom">
        <button class="sidebar-toggle" type="button" @click="toggleSidebar" :title="sidebarCollapsed ? '展开侧栏' : '收起侧栏'">
          <svg class="sidebar-toggle__icon" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10 3L5 8L10 13" />
          </svg>
        </button>
        <div class="brand">
          <div class="brand__logo">
            <img src="/logo.png" alt="MyTools" />
          </div>
          <div class="brand__text">
            <strong>MyTools</strong>
            <span>Version {{ appInfo?.version ?? '0.1.0' }}</span>
          </div>
        </div>
      </div>
    </el-aside>

    <el-main class="workspace">
      <ToolTabBar
        :tabs="openTabItems"
        :active-tool="activeTool"
        :theme="currentTheme"
        @select="openToolTab"
        @close="closeToolTab"
        @close-others="closeOtherToolTabs"
        @toggle-theme="toggleTheme"
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