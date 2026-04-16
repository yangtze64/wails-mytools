<script setup lang="ts">
import { computed, ref } from 'vue'

const sourceText = ref('')
const trimLines = ref(true)
const removeEmptyLines = ref(true)
const caseSensitive = ref(true)
const copyState = ref('')

const lines = computed(() => sourceText.value.split(/\r?\n/))

const dedupeResult = computed(() => {
  const seen = new Set<string>()
  const output: string[] = []
  let duplicateCount = 0
  let emptyCount = 0

  for (const rawLine of lines.value) {
    const normalizedLine = trimLines.value ? rawLine.trim() : rawLine

    if (removeEmptyLines.value && normalizedLine.length === 0) {
      emptyCount += 1
      continue
    }

    const key = caseSensitive.value ? normalizedLine : normalizedLine.toLocaleLowerCase()

    if (seen.has(key)) {
      duplicateCount += 1
      continue
    }

    seen.add(key)
    output.push(normalizedLine)
  }

  return {
    text: output.join('\n'),
    outputCount: output.length,
    duplicateCount,
    emptyCount,
  }
})

const inputCount = computed(() => (sourceText.value.length === 0 ? 0 : lines.value.length))

function useSample() {
  sourceText.value = ['apple', 'banana', 'apple', ' Orange ', 'orange', '', 'banana'].join('\n')
}

function clearText() {
  sourceText.value = ''
  copyState.value = ''
}

async function copyOutput() {
  if (!dedupeResult.value.text) {
    copyState.value = '没有可复制的结果'
    return
  }

  try {
    await navigator.clipboard.writeText(dedupeResult.value.text)
    copyState.value = '已复制'
  } catch {
    copyState.value = '复制失败'
  }
}
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Text</p>
        <h1>文本去重</h1>
        <p class="summary">按行保留第一次出现的内容，快速清理重复文本。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="useSample">示例</el-button>
        <el-button @click="clearText">清空</el-button>
        <el-button type="primary" @click="copyOutput">复制结果</el-button>
      </div>
    </header>

    <el-card class="option-card" shadow="never">
      <el-checkbox v-model="trimLines">去除首尾空格</el-checkbox>
      <el-checkbox v-model="removeEmptyLines">移除空行</el-checkbox>
      <el-checkbox v-model="caseSensitive">区分大小写</el-checkbox>
    </el-card>

    <el-row :gutter="12" class="dedupe-stats">
      <el-col :xs="12" :sm="6">
        <el-card shadow="never">
          <el-statistic title="输入行数" :value="inputCount" />
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="never">
          <el-statistic title="输出行数" :value="dedupeResult.outputCount" />
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="never">
          <el-statistic title="重复行" :value="dedupeResult.duplicateCount" />
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="never">
          <el-statistic title="空行" :value="dedupeResult.emptyCount" />
        </el-card>
      </el-col>
    </el-row>

    <section class="dedupe-editor">
      <el-card class="text-panel" shadow="never">
        <template #header>原始文本</template>
        <el-input
          v-model="sourceText"
          type="textarea"
          resize="none"
          clearable
          spellcheck="false"
          placeholder="每行输入一条内容"
          @clear="clearText"
        />
      </el-card>

      <el-card class="text-panel" shadow="never">
        <template #header>去重结果</template>
        <el-input
          :model-value="dedupeResult.text"
          type="textarea"
          readonly
          resize="none"
          spellcheck="false"
          placeholder="结果会自动显示在这里"
        />
      </el-card>
    </section>

    <el-alert v-if="copyState" class="copy-state" :title="copyState" type="success" :closable="false" show-icon />
  </section>
</template>
