<script setup lang="ts">
import { ref } from 'vue'

const jwtInput = ref<string>('')
const statusMessage = ref('')
const errorMessage = ref('')
const decodedHeader = ref<any>(null)
const decodedPayload = ref<any>(null)
const signature = ref<string>('')

const sampleJwt = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE4MTYyMzkwMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c'

function urlBase64Decode(str: string): string {
  let output = str.replace(/-/g, '+').replace(/_/g, '/')
  switch (output.length % 4) {
    case 0:
      break
    case 2:
      output += '=='
      break
    case 3:
      output += '='
      break
    default:
      throw new Error('Invalid base64 string')
  }
  return decodeURIComponent(escape(atob(output)))
}

function parseJwt() {
  statusMessage.value = ''
  errorMessage.value = ''
  decodedHeader.value = null
  decodedPayload.value = null
  signature.value = ''

  if (!jwtInput.value.trim()) {
    errorMessage.value = '请输入 JWT 令牌'
    return
  }

  try {
    const parts = jwtInput.value.split('.')
    if (parts.length !== 3) {
      throw new Error('无效的 JWT 格式，需要三个部分')
    }

    const headerPart = parts[0]
    const payloadPart = parts[1]
    const signaturePart = parts[2]

    decodedHeader.value = JSON.parse(urlBase64Decode(headerPart))
    decodedPayload.value = JSON.parse(urlBase64Decode(payloadPart))
    signature.value = signaturePart

    statusMessage.value = '解析成功'
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : '解析失败'
  }
}

function useSample() {
  jwtInput.value = sampleJwt
  parseJwt()
}

function clearAll() {
  jwtInput.value = ''
  decodedHeader.value = null
  decodedPayload.value = null
  signature.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

async function copyJwt() {
  if (!jwtInput.value.trim()) {
    errorMessage.value = '没有可复制的内容'
    return
  }

  try {
    await navigator.clipboard.writeText(jwtInput.value)
    statusMessage.value = '已复制 JWT'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

function formatDate(timestamp: number): string {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleString()
}
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">JWT</p>
        <h1>JWT 解析器</h1>
        <p class="summary">解析 JSON Web Token，查看 Header、Payload 和签名</p>
      </div>
    </header>

    <section class="jwt-layout">
      <aside class="jwt-sidebar">
        <div class="jwt-sidebar-section">
          <h3>JWT 令牌</h3>
          <el-input
            v-model="jwtInput"
            type="textarea"
            :rows="10"
            placeholder="请输入 JWT 令牌..."
            @input="parseJwt"
          />
        </div>
        <div class="jwt-sidebar-actions">
          <el-button @click="useSample">示例</el-button>
          <el-button @click="clearAll">清空</el-button>
          <el-button type="primary" @click="copyJwt">复制</el-button>
        </div>
      </aside>

      <section class="jwt-content">
        <div class="jwt-panels">
          <div class="jwt-panel">
            <div class="jwt-panel-header">
              <span class="jwt-panel-title">Header</span>
            </div>
            <div class="jwt-panel-content">
              <pre v-if="decodedHeader">{{ JSON.stringify(decodedHeader, null, 2) }}</pre>
              <div v-else class="jwt-empty">等待解析...</div>
            </div>
          </div>

          <div class="jwt-panel">
            <div class="jwt-panel-header">
              <span class="jwt-panel-title">Payload</span>
            </div>
            <div class="jwt-panel-content">
              <template v-if="decodedPayload">
                <pre>{{ JSON.stringify(decodedPayload, null, 2) }}</pre>
                <div v-if="decodedPayload.iat || decodedPayload.exp" class="jwt-payload-info">
                  <div v-if="decodedPayload.iat" class="jwt-info-item">
                    <span class="jwt-info-label">签发时间 (iat):</span>
                    <span class="jwt-info-value">{{ formatDate(decodedPayload.iat) }}</span>
                  </div>
                  <div v-if="decodedPayload.exp" class="jwt-info-item">
                    <span class="jwt-info-label">过期时间 (exp):</span>
                    <span class="jwt-info-value">{{ formatDate(decodedPayload.exp) }}</span>
                  </div>
                </div>
              </template>
              <div v-else class="jwt-empty">等待解析...</div>
            </div>
          </div>

          <div class="jwt-panel">
            <div class="jwt-panel-header">
              <span class="jwt-panel-title">Signature</span>
            </div>
            <div class="jwt-panel-content">
              <pre v-if="signature">{{ signature }}</pre>
              <div v-else class="jwt-empty">等待解析...</div>
            </div>
          </div>
        </div>
      </section>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>