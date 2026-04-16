<script setup lang="ts">
import { computed, ref } from 'vue'
import { diffLines, type Change } from 'diff'

const originalText = ref('')
const changedText = ref('')
const ignoreWhitespace = ref(false)
const statusMessage = ref('')
const errorMessage = ref('')

const changes = computed(() => diffLines(originalText.value, changedText.value, { ignoreWhitespace: ignoreWhitespace.value }))
const stats = computed(() => {
  let added = 0
  let removed = 0
  let unchanged = 0

  for (const change of changes.value) {
    const lineCount = countLines(change)
    if (change.added) {
      added += lineCount
    } else if (change.removed) {
      removed += lineCount
    } else {
      unchanged += lineCount
    }
  }

  return { added, removed, unchanged }
})
const diffText = computed(() => {
  return changes.value
    .map((change) => {
      const prefix = change.added ? '+ ' : change.removed ? '- ' : '  '
      return change.value
        .split(/\r?\n/)
        .filter((line, index, lines) => line.length > 0 || index < lines.length - 1)
        .map((line) => `${prefix}${line}`)
        .join('\n')
    })
    .filter(Boolean)
    .join('\n')
})

function countLines(change: Change) {
  return change.value.split(/\r?\n/).filter((line, index, lines) => line.length > 0 || index < lines.length - 1).length
}

function useSample() {
  originalText.value = ['项目计划', '1. 文本工具', '2. 二维码工具', '3. 简易脑图', '4. 发布测试'].join('\n')
  changedText.value = ['项目计划', '1. 文本工具', '2. 二维码工具', '3. 专业思维导图', '4. 文本差异对比', '5. 发布测试'].join('\n')
  statusMessage.value = ''
  errorMessage.value = ''
}

function clearText() {
  originalText.value = ''
  changedText.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

async function copyDiff() {
  if (!diffText.value) {
    statusMessage.value = '没有可复制的差异'
    return
  }

  try {
    await navigator.clipboard.writeText(diffText.value)
    statusMessage.value = '已复制差异'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Text</p>
        <h1>文本差异对比</h1>
        <p class="summary">对比两段文本，快速查看新增、删除和未变化的行。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="useSample">示例</el-button>
        <el-button @click="clearText">清空</el-button>
        <el-button type="primary" @click="copyDiff">复制差异</el-button>
      </div>
    </header>

    <el-card class="option-card" shadow="never">
      <el-checkbox v-model="ignoreWhitespace">忽略空白差异</el-checkbox>
    </el-card>

    <el-row :gutter="12" class="diff-stats">
      <el-col :xs="12" :sm="8">
        <el-card shadow="never">
          <el-statistic title="新增行" :value="stats.added" />
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="8">
        <el-card shadow="never">
          <el-statistic title="删除行" :value="stats.removed" />
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="8">
        <el-card shadow="never">
          <el-statistic title="未变化" :value="stats.unchanged" />
        </el-card>
      </el-col>
    </el-row>

    <section class="diff-inputs">
      <el-card class="text-panel" shadow="never">
        <template #header>原始文本</template>
        <el-input
          v-model="originalText"
          type="textarea"
          resize="none"
          clearable
          spellcheck="false"
          placeholder="粘贴原始文本"
        />
      </el-card>

      <el-card class="text-panel" shadow="never">
        <template #header>新文本</template>
        <el-input
          v-model="changedText"
          type="textarea"
          resize="none"
          clearable
          spellcheck="false"
          placeholder="粘贴新文本"
        />
      </el-card>
    </section>

    <el-card class="diff-result" shadow="never">
      <template #header>差异结果</template>
      <div class="diff-lines">
        <template v-for="(change, index) in changes" :key="index">
          <div
            v-for="(line, lineIndex) in change.value.split(/\r?\n/).filter((item, itemIndex, items) => item.length > 0 || itemIndex < items.length - 1)"
            :key="`${index}-${lineIndex}`"
            class="diff-line"
            :class="{ 'diff-line--added': change.added, 'diff-line--removed': change.removed }"
          >
            <span>{{ change.added ? '+' : change.removed ? '-' : ' ' }}</span>
            <code>{{ line }}</code>
          </div>
        </template>
      </div>
    </el-card>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
