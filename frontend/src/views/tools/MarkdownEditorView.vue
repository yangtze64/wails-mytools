<script setup lang="ts">
import { computed, ref } from 'vue'
import MarkdownIt from 'markdown-it'

const markdown = ref('')
const statusMessage = ref('')
const errorMessage = ref('')

const parser = new MarkdownIt({
  html: false,
  linkify: true,
  breaks: true,
})

const html = computed(() => parser.render(markdown.value))
const stats = computed(() => {
  const text = markdown.value
  return {
    chars: text.length,
    words: text.trim() ? text.trim().split(/\s+/).length : 0,
    lines: text ? text.split(/\r?\n/).length : 0,
  }
})

function useSample() {
  markdown.value = [
    '# MyTools',
    '',
    '一个基于 **Wails 3** 和 **Vue 3** 的桌面工具箱。',
    '',
    '## 工具',
    '',
    '- JSON 格式化',
    '- 文本差异对比',
    '- 二维码生成和解析',
    '- 专业思维导图',
    '',
    '```json',
    '{ "name": "MyTools", "version": "0.1.0" }',
    '```',
  ].join('\n')
  statusMessage.value = ''
  errorMessage.value = ''
}

function clearMarkdown() {
  markdown.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

async function copyMarkdown() {
  await copyText(markdown.value, '已复制 Markdown')
}

async function copyHtml() {
  await copyText(html.value, '已复制 HTML')
}

async function copyText(value: string, message: string) {
  if (!value) {
    statusMessage.value = '没有可复制的内容'
    return
  }

  try {
    await navigator.clipboard.writeText(value)
    statusMessage.value = message
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

useSample()
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Text</p>
        <h1>Markdown 编辑器</h1>
        <p class="summary">实时编辑和预览 Markdown，支持复制源码或 HTML。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="useSample">示例</el-button>
        <el-button @click="clearMarkdown">清空</el-button>
        <el-button @click="copyMarkdown">复制 Markdown</el-button>
        <el-button type="primary" @click="copyHtml">复制 HTML</el-button>
      </div>
    </header>

    <el-row :gutter="12" class="markdown-stats">
      <el-col :xs="12" :sm="8">
        <el-card shadow="never">
          <el-statistic title="字符数" :value="stats.chars" />
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="8">
        <el-card shadow="never">
          <el-statistic title="词数" :value="stats.words" />
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="8">
        <el-card shadow="never">
          <el-statistic title="行数" :value="stats.lines" />
        </el-card>
      </el-col>
    </el-row>

    <section class="markdown-layout">
      <el-card class="text-panel markdown-editor" shadow="never">
        <template #header>Markdown</template>
        <el-input
          v-model="markdown"
          type="textarea"
          resize="none"
          clearable
          spellcheck="false"
          placeholder="输入 Markdown 内容"
          @clear="clearMarkdown"
        />
      </el-card>

      <el-card class="markdown-preview-card" shadow="never">
        <template #header>预览</template>
        <article class="markdown-preview" v-html="html"></article>
      </el-card>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
