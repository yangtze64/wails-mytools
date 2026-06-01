<script setup lang="ts">
import { computed, ref, watch } from 'vue'

interface ImageFile {
  file: File
  name: string
  size: number
  type: string
}

const quality = ref(85)
const originalImage = ref<string | null>(null)
const compressedImage = ref<string | null>(null)
const originalFile = ref<ImageFile | null>(null)
const compressedBlob = ref<Blob | null>(null)
const originalSize = ref(0)
const compressedSize = ref(0)
const isCompressing = ref(false)
const isDragging = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const compressionRatio = computed(() => {
  if (!originalSize.value || !compressedSize.value) return 0
  return Math.round((1 - compressedSize.value / originalSize.value) * 100)
})

const sizeDisplay = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

function triggerUpload() {
  fileInput.value?.click()
}

function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement
  if (input.files && input.files[0]) {
    processFile(input.files[0])
  }
  input.value = ''
}

function handleDrop(event: DragEvent) {
  event.preventDefault()
  isDragging.value = false
  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    processFile(event.dataTransfer.files[0])
  }
}

function handleDragOver(event: DragEvent) {
  event.preventDefault()
  isDragging.value = true
}

function handleDragLeave() {
  isDragging.value = false
}

async function processFile(file: File) {
  if (!file.type.startsWith('image/')) {
    return
  }

  originalFile.value = {
    file,
    name: file.name,
    size: file.size,
    type: file.type,
  }
  originalSize.value = file.size

  const reader = new FileReader()
  reader.onload = (e) => {
    originalImage.value = e.target?.result as string
    compressImage()
  }
  reader.readAsDataURL(file)
}

async function compressImage() {
  if (!originalImage.value) return

  isCompressing.value = true

  try {
    const img = new Image()
    img.onload = () => {
      const canvas = document.createElement('canvas')
      canvas.width = img.width
      canvas.height = img.height

      const ctx = canvas.getContext('2d')
      if (!ctx) return

      ctx.drawImage(img, 0, 0)

      canvas.toBlob(
        (blob) => {
          if (blob) {
            if (compressedImage.value) {
              URL.revokeObjectURL(compressedImage.value)
            }
            compressedBlob.value = blob
            compressedSize.value = blob.size
            compressedImage.value = URL.createObjectURL(blob)
          }
          isCompressing.value = false
        },
        originalFile.value?.type || 'image/jpeg',
        quality.value / 100
      )
    }
    img.src = originalImage.value
  } catch (error) {
    isCompressing.value = false
  }
}

function downloadCompressed() {
  if (!compressedBlob.value || !originalFile.value) return

  const link = document.createElement('a')
  const originalName = originalFile.value.name
  const lastDot = originalName.lastIndexOf('.')
  const baseName = lastDot > 0 ? originalName.substring(0, lastDot) : originalName
  const ext = lastDot > 0 ? originalName.substring(lastDot) : '.jpg'

  link.href = URL.createObjectURL(compressedBlob.value)
  link.download = `${baseName}_compressed_${quality.value}${ext}`
  link.click()
  URL.revokeObjectURL(link.href)
}

function clearAll() {
  if (compressedImage.value) {
    URL.revokeObjectURL(compressedImage.value)
  }
  originalImage.value = null
  compressedImage.value = null
  originalFile.value = null
  compressedBlob.value = null
  originalSize.value = 0
  compressedSize.value = 0
}

watch(quality, () => {
  if (originalImage.value) {
    compressImage()
  }
})
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Image</p>
        <h1>图片压缩</h1>
        <p class="summary">上传图片并调整压缩质量，左侧原图，右侧压缩效果预览。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="triggerUpload">上传图片</el-button>
        <el-button @click="clearAll">清空</el-button>
        <el-button type="primary" :disabled="!compressedBlob" @click="downloadCompressed">
          下载压缩图
        </el-button>
        <input ref="fileInput" class="file-input" type="file" accept="image/*" @change="handleFileSelect" />
      </div>
    </header>

    <div class="image-compress-controls">
      <div class="quality-control">
        <span class="quality-label">压缩质量</span>
        <el-slider v-model="quality" :min="10" :max="100" :step="1" />
        <span class="quality-value">{{ quality }}%</span>
      </div>
      <div v-if="originalFile" class="image-compress-stats">
        <div class="stat-item">
          <span class="stat-label">原始</span>
          <strong>{{ sizeDisplay(originalSize) }}</strong>
        </div>
        <div class="stat-item">
          <span class="stat-label">压缩后</span>
          <strong>{{ sizeDisplay(compressedSize) }}</strong>
        </div>
        <div class="stat-item stat-accent">
          <span class="stat-label">节省</span>
          <strong>{{ compressionRatio }}%</strong>
        </div>
      </div>
    </div>

    <div
      class="image-compress-preview"
      :class="{ 'is-dragging': isDragging }"
      @drop="handleDrop"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
    >
      <div v-if="!originalImage" class="preview-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <rect x="3" y="3" width="18" height="18" rx="2" />
          <circle cx="8.5" cy="8.5" r="1.5" />
          <path d="M21 15l-5-5L5 21" />
        </svg>
        <p>拖拽图片到这里</p>
        <span>或点击上方"上传图片"按钮</span>
      </div>

      <template v-else>
        <div class="preview-panel">
          <div class="preview-header">
            <strong>原图</strong>
            <span>{{ sizeDisplay(originalSize) }}</span>
          </div>
          <div class="preview-image">
            <img :src="originalImage" alt="原始图片" />
          </div>
        </div>

        <div class="preview-panel">
          <div class="preview-header">
            <strong>压缩后</strong>
            <span>{{ sizeDisplay(compressedSize) }}</span>
          </div>
          <div class="preview-image">
            <img v-if="compressedImage && !isCompressing" :src="compressedImage" alt="压缩后的图片" />
            <div v-else-if="isCompressing" class="preview-loading">
              <p>压缩中...</p>
            </div>
          </div>
        </div>
      </template>
    </div>
  </section>
</template>
