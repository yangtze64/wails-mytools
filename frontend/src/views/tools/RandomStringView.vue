<script setup lang="ts">
import { computed, ref } from 'vue'

const length = ref(32)
const count = ref(5)
const useUppercase = ref(true)
const useLowercase = ref(true)
const useNumbers = ref(true)
const useSymbols = ref(false)
const excludeAmbiguous = ref(true)
const customChars = ref('')
const results = ref<string[]>([])
const statusMessage = ref('')
const errorMessage = ref('')

const uppercaseChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
const lowercaseChars = 'abcdefghijklmnopqrstuvwxyz'
const numberChars = '0123456789'
const symbolChars = '!@#$%^&*_-+=?'
const ambiguousChars = new Set(['0', 'O', 'o', '1', 'I', 'l'])

const charPool = computed(() => {
  let pool = ''

  if (useUppercase.value) {
    pool += uppercaseChars
  }
  if (useLowercase.value) {
    pool += lowercaseChars
  }
  if (useNumbers.value) {
    pool += numberChars
  }
  if (useSymbols.value) {
    pool += symbolChars
  }
  if (customChars.value) {
    pool += customChars.value
  }

  const uniqueChars = Array.from(new Set(pool.split('')))
  const filteredChars = excludeAmbiguous.value
    ? uniqueChars.filter((char) => !ambiguousChars.has(char))
    : uniqueChars

  return filteredChars.join('')
})

function generateRandomStrings() {
  statusMessage.value = ''
  errorMessage.value = ''

  if (!charPool.value) {
    results.value = []
    errorMessage.value = '请至少选择一种字符来源'
    return
  }

  const generated: string[] = []
  const randomValues = new Uint32Array(length.value * count.value)
  crypto.getRandomValues(randomValues)

  for (let rowIndex = 0; rowIndex < count.value; rowIndex += 1) {
    let value = ''
    for (let charIndex = 0; charIndex < length.value; charIndex += 1) {
      const randomIndex = randomValues[rowIndex * length.value + charIndex] % charPool.value.length
      value += charPool.value[randomIndex]
    }
    generated.push(value)
  }

  results.value = generated
  statusMessage.value = `已生成 ${generated.length} 条`
}

async function copyResults() {
  if (!results.value.length) {
    statusMessage.value = '没有可复制的内容'
    return
  }

  try {
    await navigator.clipboard.writeText(results.value.join('\n'))
    statusMessage.value = '已复制'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

function clearResults() {
  results.value = []
  statusMessage.value = ''
  errorMessage.value = ''
}

generateRandomStrings()
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Text</p>
        <h1>随机字符串</h1>
        <p class="summary">生成密码、Token、测试数据和一次性随机标识。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="clearResults">清空</el-button>
        <el-button :disabled="!results.length" @click="copyResults">复制结果</el-button>
        <el-button type="primary" @click="generateRandomStrings">生成</el-button>
      </div>
    </header>

    <section class="random-layout">
      <el-card class="random-options" shadow="never">
        <template #header>生成配置</template>

        <div class="form-grid">
          <label>
            <span>字符串长度</span>
            <el-input-number v-model="length" :min="1" :max="256" controls-position="right" />
          </label>
          <label>
            <span>生成数量</span>
            <el-input-number v-model="count" :min="1" :max="100" controls-position="right" />
          </label>
        </div>

        <div class="random-checkboxes">
          <el-checkbox v-model="useUppercase">大写字母</el-checkbox>
          <el-checkbox v-model="useLowercase">小写字母</el-checkbox>
          <el-checkbox v-model="useNumbers">数字</el-checkbox>
          <el-checkbox v-model="useSymbols">符号</el-checkbox>
          <el-checkbox v-model="excludeAmbiguous">排除易混淆字符</el-checkbox>
        </div>

        <label class="custom-chars">
          <span>自定义字符</span>
          <el-input v-model="customChars" clearable placeholder="可追加自己的字符集" />
        </label>

        <el-alert
          class="pool-preview"
          :title="`当前字符池：${charPool.length} 个字符`"
          :description="charPool || '未选择字符'"
          type="info"
          :closable="false"
        />
      </el-card>

      <el-card class="random-result" shadow="never">
        <template #header>生成结果</template>
        <el-input
          :model-value="results.join('\n')"
          type="textarea"
          readonly
          resize="none"
          placeholder="点击生成后显示结果"
        />
      </el-card>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
