<script setup lang="ts">
import { computed, ref } from 'vue'

const sourceText = ref('')
const resultText = ref('')
const indentSize = ref(2)
const statusMessage = ref('')
const errorMessage = ref('')

const canCopy = computed(() => resultText.value.length > 0)

function parseJson() {
  if (!sourceText.value.trim()) {
    throw new Error('请输入 JSON 内容')
  }

  return JSON.parse(sourceText.value)
}

function formatJson() {
  runJsonAction(() => {
    resultText.value = JSON.stringify(parseJson(), null, indentSize.value)
    statusMessage.value = '格式化成功'
  })
}

function minifyJson() {
  runJsonAction(() => {
    resultText.value = JSON.stringify(parseJson())
    statusMessage.value = '压缩成功'
  })
}

function validateJson() {
  runJsonAction(() => {
    parseJson()
    resultText.value = ''
    statusMessage.value = 'JSON 校验通过'
  })
}

function runJsonAction(action: () => void) {
  errorMessage.value = ''
  statusMessage.value = ''

  try {
    action()
  } catch (error) {
    resultText.value = ''
    errorMessage.value = error instanceof Error ? error.message : 'JSON 处理失败'
  }
}

function useSample() {
  sourceText.value = '{"name":"MyTools","version":"0.1.0","tools":["JSON格式化","二维码","文本差异"],"enabled":true}'
  resultText.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

function clearText() {
  sourceText.value = ''
  resultText.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

async function copyResult() {
  if (!resultText.value) {
    statusMessage.value = '没有可复制的结果'
    return
  }

  try {
    await navigator.clipboard.writeText(resultText.value)
    statusMessage.value = '已复制'
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
        <h1>JSON 工具</h1>
        <p class="summary">格式化、压缩和校验 JSON，快速定位语法错误。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="useSample">示例</el-button>
        <el-button @click="clearText">清空</el-button>
        <el-button :disabled="!canCopy" @click="copyResult">复制结果</el-button>
        <el-button @click="validateJson">校验</el-button>
        <el-button @click="minifyJson">压缩</el-button>
        <el-button type="primary" @click="formatJson">格式化</el-button>
      </div>
    </header>

    <el-card class="option-card" shadow="never">
      <label class="json-indent">
        <span>缩进空格</span>
        <el-input-number v-model="indentSize" :min="0" :max="8" controls-position="right" />
      </label>
    </el-card>

    <section class="json-editor">
      <el-card class="text-panel" shadow="never">
        <template #header>JSON 输入</template>
        <el-input
          v-model="sourceText"
          type="textarea"
          resize="none"
          clearable
          spellcheck="false"
          placeholder="粘贴 JSON 内容"
          @clear="clearText"
        />
      </el-card>

      <el-card class="text-panel" shadow="never">
        <template #header>处理结果</template>
        <el-input
          :model-value="resultText"
          type="textarea"
          readonly
          resize="none"
          spellcheck="false"
          placeholder="格式化或压缩后的结果会显示在这里"
        />
      </el-card>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
