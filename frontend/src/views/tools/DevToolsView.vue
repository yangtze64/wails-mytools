<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getDevEnvironment, openDevTools, type DevEnvironment } from '../../api/devtools'

const environment = ref<DevEnvironment | null>(null)
const statusMessage = ref('')
const errorMessage = ref('')

onMounted(loadEnvironment)

async function loadEnvironment() {
  errorMessage.value = ''

  try {
    environment.value = await getDevEnvironment()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '读取调试信息失败'
  }
}

async function handleOpenDevTools() {
  errorMessage.value = ''
  statusMessage.value = ''

  try {
    await openDevTools()
    statusMessage.value = '已打开开发者工具'
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '打开开发者工具失败'
  }
}

async function copyEnvironment() {
  if (!environment.value) {
    statusMessage.value = '没有可复制的信息'
    return
  }

  try {
    await navigator.clipboard.writeText(JSON.stringify(environment.value, null, 2))
    statusMessage.value = '已复制环境信息'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}
</script>

<template>
  <section class="tool-page devtools-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Developer</p>
        <h1>开发者调试</h1>
        <p class="summary">开发环境专用入口，用于打开 WebView DevTools 和查看运行信息。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="loadEnvironment">刷新信息</el-button>
        <el-button @click="copyEnvironment">复制信息</el-button>
        <el-button type="primary" @click="handleOpenDevTools">打开 DevTools</el-button>
      </div>
    </header>

    <el-card class="devtools-card" shadow="never">
      <template #header>运行环境</template>
      <dl class="devtools-info">
        <div>
          <dt>Go 版本</dt>
          <dd>{{ environment?.goVersion || '-' }}</dd>
        </div>
        <div>
          <dt>系统</dt>
          <dd>{{ environment ? `${environment.os}/${environment.arch}` : '-' }}</dd>
        </div>
        <div>
          <dt>前端 Dev Server</dt>
          <dd>{{ environment?.frontendDevUrl || '-' }}</dd>
        </div>
        <div>
          <dt>工作目录</dt>
          <dd>{{ environment?.workingDirectory || '-' }}</dd>
        </div>
      </dl>
    </el-card>

    <el-card class="devtools-card" shadow="never">
      <template #header>常用提示</template>
      <ul class="devtools-tips">
        <li>页面样式、控制台日志、网络请求可以在 DevTools 中查看。</li>
        <li>如果窗口没有反应，先确认当前运行的是 `task dev`。</li>
        <li>生产构建不会显示这个工具入口。</li>
      </ul>
    </el-card>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
