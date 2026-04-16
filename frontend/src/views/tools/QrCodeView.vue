<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import QRCode, { type QRCodeErrorCorrectionLevel } from 'qrcode'

type CorrectionLevel = Extract<QRCodeErrorCorrectionLevel, 'L' | 'M' | 'Q' | 'H'>

const content = ref('https://github.com')
const size = ref(260)
const margin = ref(2)
const darkColor = ref('#111827')
const lightColor = ref('#ffffff')
const errorCorrectionLevel = ref<CorrectionLevel>('M')
const logoDataUrl = ref('')
const logoSizePercent = ref(22)
const qrDataUrl = ref('')
const errorMessage = ref('')
const logoInput = ref<HTMLInputElement | null>(null)
let renderVersion = 0

const correctionOptions: Array<{ label: string; value: CorrectionLevel }> = [
  { label: '低 L', value: 'L' },
  { label: '中 M', value: 'M' },
  { label: '较高 Q', value: 'Q' },
  { label: '高 H', value: 'H' },
]

const hasContent = computed(() => content.value.trim().length > 0)

async function generateQrCode() {
  const currentVersion = ++renderVersion
  errorMessage.value = ''

  if (!hasContent.value) {
    qrDataUrl.value = ''
    return
  }

  try {
    const canvas = document.createElement('canvas')
    await QRCode.toCanvas(canvas, content.value, {
      width: size.value,
      margin: margin.value,
      errorCorrectionLevel: errorCorrectionLevel.value,
      color: {
        dark: toQrColor(darkColor.value),
        light: toQrColor(lightColor.value),
      },
    })

    if (logoDataUrl.value) {
      await drawLogo(canvas, logoDataUrl.value)
    }

    if (currentVersion !== renderVersion) {
      return
    }

    qrDataUrl.value = canvas.toDataURL('image/png')
  } catch (error) {
    if (currentVersion !== renderVersion) {
      return
    }

    qrDataUrl.value = ''
    errorMessage.value = error instanceof Error ? error.message : '二维码生成失败'
  }
}

function toQrColor(color: string) {
  const normalizedColor = color.trim()

  if (/^#[0-9a-fA-F]{3}$/.test(normalizedColor)) {
    return `${normalizedColor}f`
  }

  if (/^#[0-9a-fA-F]{4}$/.test(normalizedColor) || /^#[0-9a-fA-F]{8}$/.test(normalizedColor)) {
    return normalizedColor
  }

  if (/^#[0-9a-fA-F]{6}$/.test(normalizedColor)) {
    return `${normalizedColor}ff`
  }

  return '#000000ff'
}

function drawLogo(canvas: HTMLCanvasElement, logoSource: string) {
  return new Promise<void>((resolve, reject) => {
    const image = new Image()
    image.onload = () => {
      const context = canvas.getContext('2d')
      if (!context) {
        reject(new Error('无法创建图片上下文'))
        return
      }

      const logoSize = Math.round(canvas.width * (logoSizePercent.value / 100))
      const padding = Math.round(logoSize * 0.12)
      const backgroundSize = logoSize + padding * 2
      const x = Math.round((canvas.width - backgroundSize) / 2)
      const y = Math.round((canvas.height - backgroundSize) / 2)
      const radius = Math.round(backgroundSize * 0.14)

      drawRoundedRect(context, x, y, backgroundSize, backgroundSize, radius)
      context.fillStyle = lightColor.value || '#ffffff'
      context.fill()

      const logoX = x + padding
      const logoY = y + padding
      context.drawImage(image, logoX, logoY, logoSize, logoSize)
      resolve()
    }
    image.onerror = () => reject(new Error('Logo 图片读取失败'))
    image.src = logoSource
  })
}

function drawRoundedRect(
  context: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  height: number,
  radius: number,
) {
  context.beginPath()
  context.moveTo(x + radius, y)
  context.lineTo(x + width - radius, y)
  context.quadraticCurveTo(x + width, y, x + width, y + radius)
  context.lineTo(x + width, y + height - radius)
  context.quadraticCurveTo(x + width, y + height, x + width - radius, y + height)
  context.lineTo(x + radius, y + height)
  context.quadraticCurveTo(x, y + height, x, y + height - radius)
  context.lineTo(x, y + radius)
  context.quadraticCurveTo(x, y, x + radius, y)
  context.closePath()
}

function useSample() {
  content.value = 'https://github.com'
}

function clearContent() {
  content.value = ''
}

function handleLogoUpload(event: Event) {
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
  reader.onload = () => {
    logoDataUrl.value = String(reader.result)
    input.value = ''
  }
  reader.onerror = () => {
    errorMessage.value = 'Logo 图片读取失败'
    input.value = ''
  }
  reader.readAsDataURL(file)
}

function triggerLogoUpload() {
  logoInput.value?.click()
}

function clearLogo() {
  logoDataUrl.value = ''
}

function downloadQrCode() {
  if (!qrDataUrl.value) {
    errorMessage.value = '没有可下载的二维码'
    return
  }

  const link = document.createElement('a')
  link.href = qrDataUrl.value
  link.download = 'qrcode.png'
  link.click()
}

watch(
  [content, size, margin, darkColor, lightColor, errorCorrectionLevel, logoDataUrl, logoSizePercent],
  generateQrCode,
  { immediate: true },
)
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Image</p>
        <h1>生成二维码</h1>
        <p class="summary">输入文本、链接或配置内容，生成可下载的 PNG 二维码。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="useSample">示例</el-button>
        <el-button @click="clearContent">清空</el-button>
        <el-button type="primary" :disabled="!qrDataUrl" @click="downloadQrCode">下载 PNG</el-button>
      </div>
    </header>

    <section class="qr-layout">
      <el-card class="qr-form" shadow="never">
        <template #header>内容</template>
        <el-input
          v-model="content"
          type="textarea"
          :autosize="{ minRows: 8, maxRows: 12 }"
          resize="none"
          clearable
          placeholder="输入需要编码到二维码里的内容"
          @clear="clearContent"
        />

        <div class="form-grid">
          <label>
            <span>尺寸</span>
            <el-input-number v-model="size" :min="160" :max="720" :step="20" controls-position="right" />
          </label>
          <label>
            <span>边距</span>
            <el-input-number v-model="margin" :min="0" :max="8" controls-position="right" />
          </label>
          <label>
            <span>纠错级别</span>
            <el-select v-model="errorCorrectionLevel">
              <el-option
                v-for="option in correctionOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </label>
        </div>

        <div class="color-grid">
          <label class="color-field">
            <span>前景色</span>
            <div>
              <input v-model="darkColor" class="native-color-picker" type="color" aria-label="前景色" />
              <el-input v-model="darkColor" />
            </div>
          </label>
          <label class="color-field">
            <span>背景色</span>
            <div>
              <input v-model="lightColor" class="native-color-picker" type="color" aria-label="背景色" />
              <el-input v-model="lightColor" />
            </div>
          </label>
        </div>

        <div class="logo-panel">
          <div class="logo-panel__header">
            <span>中间 Logo</span>
            <div>
              <input ref="logoInput" class="file-input" type="file" accept="image/*" @change="handleLogoUpload" />
              <el-button @click="triggerLogoUpload">上传 Logo</el-button>
              <el-button :disabled="!logoDataUrl" @click="clearLogo">移除</el-button>
            </div>
          </div>

          <div class="logo-panel__body">
            <div class="logo-preview">
              <img v-if="logoDataUrl" :src="logoDataUrl" alt="Logo 预览" />
              <span v-else>未上传</span>
            </div>
            <label class="logo-size">
              <span>Logo 比例</span>
              <el-input-number
                v-model="logoSizePercent"
                :min="12"
                :max="32"
                :step="1"
                controls-position="right"
              />
            </label>
          </div>
        </div>
      </el-card>

      <el-card class="qr-preview" shadow="never">
        <template #header>预览</template>
        <div class="qr-preview__stage">
          <img v-if="qrDataUrl" :src="qrDataUrl" alt="二维码预览" />
          <el-empty v-else description="输入内容后生成二维码" />
        </div>
      </el-card>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
  </section>
</template>
