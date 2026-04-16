<script setup lang="ts">
import { computed, ref } from 'vue'
import jsQR from 'jsqr'

const imageDataUrl = ref('')
const decodedText = ref('')
const statusMessage = ref('')
const errorMessage = ref('')
const imageInput = ref<HTMLInputElement | null>(null)

const hasResult = computed(() => decodedText.value.length > 0)

function triggerImageUpload() {
  imageInput.value?.click()
}

function handleImageUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) {
    return
  }

  if (!file.type.startsWith('image/')) {
    errorMessage.value = '请选择图片文件'
    input.value = ''
    return
  }

  const reader = new FileReader()
  reader.onload = async () => {
    imageDataUrl.value = String(reader.result)
    input.value = ''
    await decodeImage(imageDataUrl.value)
  }
  reader.onerror = () => {
    errorMessage.value = '图片读取失败'
    input.value = ''
  }
  reader.readAsDataURL(file)
}

async function decodeImage(source: string) {
  decodedText.value = ''
  statusMessage.value = ''
  errorMessage.value = ''

  try {
    const image = await loadImage(source)
    const canvas = document.createElement('canvas')
    const context = canvas.getContext('2d', { willReadFrequently: true })

    if (!context) {
      errorMessage.value = '无法创建图片上下文'
      return
    }

    canvas.width = image.naturalWidth
    canvas.height = image.naturalHeight
    context.drawImage(image, 0, 0)

    const pixels = context.getImageData(0, 0, canvas.width, canvas.height)
    const result = jsQR(pixels.data, pixels.width, pixels.height, {
      inversionAttempts: 'attemptBoth',
    })

    if (!result) {
      statusMessage.value = '未识别到二维码内容'
      return
    }

    decodedText.value = result.data
    statusMessage.value = `识别成功，二维码版本 ${result.version}`
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '二维码解析失败'
  }
}

function loadImage(source: string) {
  return new Promise<HTMLImageElement>((resolve, reject) => {
    const image = new Image()
    image.onload = () => resolve(image)
    image.onerror = () => reject(new Error('图片加载失败'))
    image.src = source
  })
}

async function copyResult() {
  if (!decodedText.value) {
    statusMessage.value = '没有可复制的内容'
    return
  }

  try {
    await navigator.clipboard.writeText(decodedText.value)
    statusMessage.value = '已复制'
  } catch {
    errorMessage.value = '复制失败'
  }
}

function clearImage() {
  imageDataUrl.value = ''
  decodedText.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Image</p>
        <h1>解析二维码</h1>
        <p class="summary">上传包含二维码的图片，提取其中的文本、链接或配置内容。</p>
      </div>
      <div class="tool-actions">
        <input ref="imageInput" class="file-input" type="file" accept="image/*" @change="handleImageUpload" />
        <el-button type="primary" @click="triggerImageUpload">上传图片</el-button>
        <el-button :disabled="!imageDataUrl" @click="clearImage">清空</el-button>
        <el-button :disabled="!hasResult" @click="copyResult">复制结果</el-button>
      </div>
    </header>

    <section class="qr-decode-layout">
      <el-card class="decode-preview" shadow="never">
        <template #header>图片预览</template>
        <div class="decode-preview__stage">
          <img v-if="imageDataUrl" :src="imageDataUrl" alt="二维码图片预览" />
          <el-empty v-else description="上传二维码图片" />
        </div>
      </el-card>

      <el-card class="decode-result" shadow="never">
        <template #header>解析结果</template>
        <el-input
          :model-value="decodedText"
          type="textarea"
          readonly
          resize="none"
          placeholder="识别结果会显示在这里"
        />
      </el-card>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
