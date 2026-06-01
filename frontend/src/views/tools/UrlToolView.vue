<script setup lang="ts">
import { computed, ref } from 'vue'

type Action = 'encode' | 'decode' | 'parse'

const action = ref<Action>('encode')
const inputText = ref('')
const statusMessage = ref('')
const errorMessage = ref('')

const encodedResult = computed(() => {
  if (!inputText.value.trim()) return ''
  if (action.value !== 'encode') return ''
  try {
    return encodeURIComponent(inputText.value)
  } catch {
    return ''
  }
})

const decodedResult = computed(() => {
  if (!inputText.value.trim()) return ''
  if (action.value !== 'decode') return ''
  try {
    return decodeURIComponent(inputText.value)
  } catch {
    return ''
  }
})

interface UrlPart {
  label: string
  value: string
}

interface QueryParam {
  key: string
  value: string
}

const parsedParts = computed<UrlPart[]>(() => {
  if (!inputText.value.trim()) return []
  if (action.value !== 'parse') return []

  try {
    const url = new URL(inputText.value)
    const parts: UrlPart[] = []

    if (url.protocol) parts.push({ label: '协议', value: url.protocol })
    if (url.hostname) parts.push({ label: '主机名', value: url.hostname })
    if (url.port) parts.push({ label: '端口', value: url.port })
    if (url.pathname && url.pathname !== '/') parts.push({ label: '路径', value: url.pathname })
    if (url.hash) parts.push({ label: '锚点', value: url.hash })
    if (url.username) parts.push({ label: '用户名', value: url.username })
    if (url.password) parts.push({ label: '密码', value: url.password })
    if (url.origin) parts.push({ label: 'Origin', value: url.origin })
    if (url.href) parts.push({ label: '完整 URL', value: url.href })

    return parts
  } catch {
    return []
  }
})

const queryParams = computed<QueryParam[]>(() => {
  if (!inputText.value.trim()) return []
  if (action.value !== 'parse') return []

  try {
    const url = new URL(inputText.value)
    const params: QueryParam[] = []
    url.searchParams.forEach((value, key) => {
      params.push({ key, value })
    })
    return params
  } catch {
    return []
  }
})

const queryString = computed(() => {
  if (!inputText.value.trim()) return ''
  if (action.value !== 'parse') return ''
  try {
    const url = new URL(inputText.value)
    return url.search
  } catch {
    return ''
  }
})

const outputText = computed(() => {
  if (action.value === 'encode') return encodedResult.value
  if (action.value === 'decode') return decodedResult.value
  return ''
})

function switchAction(newAction: Action) {
  action.value = newAction
  statusMessage.value = ''
  errorMessage.value = ''
}

function clearAll() {
  inputText.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

async function copyOutput() {
  if (!outputText.value) {
    errorMessage.value = '没有可复制的内容'
    return
  }

  try {
    await navigator.clipboard.writeText(outputText.value)
    statusMessage.value = '已复制结果'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

async function copyText(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    statusMessage.value = '已复制'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

function useSample() {
  if (action.value === 'encode') {
    inputText.value = 'https://example.com/api?name=张三&age=25&city=北京'
  } else if (action.value === 'decode') {
    inputText.value = 'https%3A%2F%2Fexample.com%2Fapi%3Fname%3D%E5%BC%A0%E4%B8%89%26age%3D25%26city%3D%E5%8C%97%E4%BA%AC'
  } else {
    inputText.value = 'https://example.com:8080/api/users?name=张三&age=25&city=北京&tags=vue,react#section1'
  }
  statusMessage.value = '已加载示例'
  errorMessage.value = ''
}
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">URL</p>
        <h1>URL 编解码工具</h1>
        <p class="summary">URL 编码、解码和参数解析</p>
      </div>
    </header>

    <section class="url-layout">
      <aside class="url-sidebar">
        <div class="url-action-switch">
          <button
            v-for="a in (['encode', 'decode', 'parse'] as Action[])"
            :key="a"
            class="action-btn"
            :class="{ 'is-active': action === a }"
            type="button"
            @click="switchAction(a)"
          >
            {{ a === 'encode' ? '编码' : a === 'decode' ? '解码' : '解析' }}
          </button>
        </div>

        <div class="url-sidebar-section">
          <h3>输入</h3>
          <el-input
            v-model="inputText"
            type="textarea"
            :rows="12"
            :placeholder="action === 'encode' ? '请输入要编码的 URL 或文本...' : action === 'decode' ? '请输入要解码的 URL 编码...' : '请输入要解析的 URL...'"
            resize="none"
          />
        </div>

        <div class="url-sidebar-footer">
          <el-button @click="useSample">示例</el-button>
          <el-button @click="clearAll">清空</el-button>
          <el-button
            v-if="action !== 'parse'"
            type="primary"
            @click="copyOutput"
            :disabled="!outputText"
          >
            复制结果
          </el-button>
        </div>
      </aside>

      <section class="url-content">
        <template v-if="action === 'parse'">
          <div class="url-parse-section">
            <h3>URL 组成部分</h3>
            <div v-if="parsedParts.length > 0" class="url-parts">
              <div v-for="part in parsedParts" :key="part.label" class="url-part-item">
                <span class="url-part-label">{{ part.label }}</span>
                <span class="url-part-value">{{ part.value }}</span>
                <button class="url-part-copy" @click="copyText(part.value)">复制</button>
              </div>
            </div>
            <div v-else class="url-empty">输入 URL 后将自动解析各组成部分</div>
          </div>

          <div class="url-parse-section">
            <h3>Query 参数</h3>
            <div v-if="queryParams.length > 0" class="url-params">
              <div class="url-params-header">
                <span>键</span>
                <span>值</span>
                <span></span>
              </div>
              <div v-for="param in queryParams" :key="param.key + param.value" class="url-param-item">
                <span class="url-param-key">{{ param.key }}</span>
                <span class="url-param-value">{{ param.value }}</span>
                <button class="url-part-copy" @click="copyText(`${param.key}=${param.value}`)">复制</button>
              </div>
            </div>
            <div v-else class="url-empty">没有查询参数</div>
          </div>

          <div v-if="queryString" class="url-parse-section">
            <h3>Query String</h3>
            <div class="url-query-string">
              <code>{{ queryString }}</code>
              <button class="url-part-copy" @click="copyText(queryString)">复制</button>
            </div>
          </div>
        </template>

        <template v-else>
          <div class="url-result-section">
            <div class="url-result-header">
              <h3>{{ action === 'encode' ? '编码结果' : '解码结果' }}</h3>
            </div>
            <div class="url-result-body">
              <el-input
                :model-value="outputText"
                type="textarea"
                :rows="18"
                :placeholder="action === 'encode' ? '编码后的结果...' : '解码后的结果...'"
                resize="none"
                readonly
              />
            </div>
          </div>
        </template>
      </section>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>