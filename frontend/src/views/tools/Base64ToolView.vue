<script setup lang="ts">
import { computed, ref } from 'vue'

type Mode = 'text' | 'file'
type Action = 'encode' | 'decode'

const mode = ref<Mode>('text')
const action = ref<Action>('encode')
const inputText = ref('')
const outputText = ref('')
const statusMessage = ref('')
const errorMessage = ref('')
const selectedFile = ref<File | null>(null)
const filePreviewBase64 = ref('')
const isProcessing = ref(false)

const outputSizeInKB = computed(() => {
  if (!outputText.value) return '0.0'
  return (new Blob([outputText.value]).size / 1024).toFixed(1)
})

function reset() {
  inputText.value = ''
  outputText.value = ''
  selectedFile.value = null
  filePreviewBase64.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

function switchMode(newMode: Mode) {
  mode.value = newMode
  reset()
}

async function encodeText() {
  if (!inputText.value.trim()) {
    errorMessage.value = '请输入要编码的文本'
    return
  }

  try {
    outputText.value = btoa(unescape(encodeURIComponent(inputText.value)))
    statusMessage.value = '编码成功'
    errorMessage.value = ''
  } catch (error) {
    errorMessage.value = '编码失败'
  }
}

async function decodeText() {
  if (!inputText.value.trim()) {
    errorMessage.value = '请输入要解码的 Base64'
    return
  }

  try {
    outputText.value = decodeURIComponent(escape(atob(inputText.value)))
    statusMessage.value = '解码成功'
    errorMessage.value = ''
  } catch (error) {
    errorMessage.value = '解码失败，请检查输入的 Base64 是否正确'
  }
}

async function handleFileSelect(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    selectedFile.value = file
    await processFile()
  }
}

async function handleFileDrop(event: DragEvent) {
  event.preventDefault()
  const file = event.dataTransfer?.files?.[0]
  if (file) {
    selectedFile.value = file
    await processFile()
  }
}

function handleDragOver(event: DragEvent) {
  event.preventDefault()
}

async function processFile() {
  if (!selectedFile.value) return

  isProcessing.value = true
  errorMessage.value = ''

  try {
    if (action.value === 'encode') {
      const base64 = await fileToBase64(selectedFile.value)
      filePreviewBase64.value = base64
      outputText.value = base64
      statusMessage.value = '文件编码成功'
    } else {
      const binaryString = atob(inputText.value)
      const bytes = new Uint8Array(binaryString.length)
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i)
      }
      const blob = new Blob([bytes], { type: 'application/octet-stream' })
      filePreviewBase64.value = URL.createObjectURL(blob)
      statusMessage.value = '文件解码成功，可以下载'
    }
  } catch (error) {
    errorMessage.value = action.value === 'encode' ? '文件编码失败' : '解码失败，请检查输入的 Base64 是否正确'
  } finally {
    isProcessing.value = false
  }
}

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const result = reader.result as string
      // 去除 data:mime/type;base64, 前缀
      const base64 = result.split(',')[1]
      resolve(base64)
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
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

async function downloadDecodedFile() {
  if (!outputText.value) {
    errorMessage.value = '没有可下载的内容'
    return
  }

  try {
    const binaryString = atob(outputText.value)
    const bytes = new Uint8Array(binaryString.length)
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }
    const blob = new Blob([bytes], { type: 'application/octet-stream' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    const originalName = selectedFile.value?.name || 'decoded-file'
    const lastDot = originalName.lastIndexOf('.')
    const baseName = lastDot > 0 ? originalName.substring(0, lastDot) : originalName
    const ext = lastDot > 0 ? originalName.substring(lastDot) : '.bin'
    link.download = `${baseName}-decoded${ext}`
    link.click()
    URL.revokeObjectURL(link.href)
    statusMessage.value = '已下载文件'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '下载失败'
  }
}

function switchAction(newAction: Action) {
  action.value = newAction
  reset()
}
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Encoding</p>
        <h1>Base64 编解码</h1>
        <p class="summary">支持文本和文件的 Base64 编码和解码，方便在不同场景下使用。</p>
      </div>
    </header>

    <section class="base64-layout">
      <aside class="base64-sidebar">
        <div class="base64-mode-switch">
          <button
            v-for="m in ['text', 'file']"
            :key="m"
            class="mode-btn"
            :class="{ 'is-active': mode === m }"
            type="button"
            @click="switchMode(m as Mode)"
          >
            {{ m === 'text' ? '文本模式' : '文件模式' }}
          </button>
        </div>

        <div class="base64-action-switch">
          <button
            v-for="a in ['encode', 'decode']"
            :key="a"
            class="action-btn"
            :class="{ 'is-active': action === a }"
            type="button"
            @click="switchAction(a as Action)"
          >
            {{ a === 'encode' ? '编码' : '解码' }}
          </button>
        </div>

        <div class="base64-sidebar__footer">
          <el-button type="primary" v-if="mode === 'text'" @click="action === 'encode' ? encodeText() : decodeText()">
            {{ action === 'encode' ? '编码' : '解码' }}
          </el-button>
          <el-button type="primary" v-else-if="mode === 'file' && action === 'decode'" @click="processFile">
            解码
          </el-button>
          <el-button @click="copyOutput" :disabled="!outputText">复制结果</el-button>
          <el-button
            v-if="mode === 'file' && action === 'decode'"
            @click="downloadDecodedFile"
            :disabled="!outputText"
          >
            下载文件
          </el-button>
          <el-button @click="reset">清空</el-button>
        </div>
      </aside>

      <section class="base64-content">
        <div class="base64-panel">
          <div class="panel-header">
            <strong>{{ action === 'encode' ? '输入' : 'Base64' }}</strong>
          </div>
          <div v-if="mode === 'text'" class="panel-content">
            <el-input
              v-model="inputText"
              type="textarea"
              :rows="18"
              :placeholder="action === 'encode' ? '请输入要编码的文本...' : '请输入要解码的 Base64...'"
              resize="none"
            />
          </div>
          <div v-else class="panel-content">
            <div
              v-if="action === 'encode'"
              class="file-drop-zone"
              @drop="handleFileDrop"
              @dragover="handleDragOver"
            >
              <input
                ref="fileInput"
                type="file"
                class="file-input"
                @change="handleFileSelect"
              />
              <template v-if="selectedFile">
                <div class="file-info">
                  <strong>{{ selectedFile.name }}</strong>
                  <span>{{ (selectedFile.size / 1024).toFixed(1) }} KB</span>
                </div>
              </template>
              <template v-else>
                <div class="drop-hint">
                  <p>拖拽文件到这里</p>
                  <p>或</p>
                  <label class="upload-btn">
                    选择文件
                    <input type="file" @change="handleFileSelect" />
                  </label>
                </div>
              </template>
            </div>
            <div v-else class="file-decode-input">
              <el-input
                v-model="inputText"
                type="textarea"
                :rows="18"
                placeholder="请输入要解码的 Base64..."
                resize="none"
              />
            </div>
          </div>
        </div>

        <div class="base64-panel">
          <div class="panel-header">
            <strong>{{ action === 'encode' ? 'Base64' : '输出' }}</strong>
            <span v-if="outputText">{{ outputSizeInKB }} KB</span>
          </div>
          <div class="panel-content">
            <el-input
              v-model="outputText"
              type="textarea"
              :rows="18"
              :placeholder="action === 'encode' ? '编码后的 Base64 会显示在这里...' : '解码后的结果会显示在这里...'"
              resize="none"
              readonly
            />
          </div>
        </div>
      </section>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>