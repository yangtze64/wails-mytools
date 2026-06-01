<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue'
import CanvasToolbar from '../../components/CanvasToolbar.vue'

const sampleMarkdown = `# 开发工具箱

## 文本工具
- 文本去重
- 随机字符串
- 文本差异对比
- JSON 工具
- Markdown 编辑器
- Base64 编解码

## 图形工具
- Mermaid 编辑器
- 脑图
- Markmap 脑图

## 开发工具
- 时间戳转换
- Cron 表达式生成器
- JWT 解析器
- URL 编解码

## 图像工具
- 生成二维码
- 解析二维码
- 图片压缩
`

const inputMarkdown = ref(sampleMarkdown)
const statusMessage = ref('')
const errorMessage = ref('')
const svgRef = ref<SVGElement | null>(null)
const isFullscreen = ref(false)
const zoomLevel = ref(100)
let markmapInstance: any = null

async function renderMarkmap() {
  if (!inputMarkdown.value.trim()) {
    if (markmapInstance) {
      markmapInstance.setData({ content: '', children: [] })
    }
    return
  }

  try {
    const { Transformer } = await import('markmap-lib')
    const { Markmap }: any = await import('markmap-view')

    const transformer = new Transformer()
    const { root } = transformer.transform(inputMarkdown.value)

    await nextTick()

    if (svgRef.value) {
      if (markmapInstance) {
        markmapInstance.setData(root)
        await nextTick()
        markmapInstance.fit()
      } else {
        markmapInstance = Markmap.create(svgRef.value, {
          zoom: true,
          pan: true,
        }, root)
      }
    }
  } catch (error) {
    console.error('Markmap load error:', error)
    errorMessage.value = '加载失败，请刷新重试'
  }
}

function useSample() {
  inputMarkdown.value = sampleMarkdown
  statusMessage.value = '已加载示例'
  errorMessage.value = ''
}

function clearAll() {
  inputMarkdown.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

async function copyMarkdown() {
  if (!inputMarkdown.value.trim()) {
    errorMessage.value = '没有可复制的内容'
    return
  }

  try {
    await navigator.clipboard.writeText(inputMarkdown.value)
    statusMessage.value = '已复制 Markdown'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

async function downloadSvg() {
  const svgElement = svgRef.value
  if (!svgElement) {
    errorMessage.value = '没有可下载的内容'
    return
  }

  try {
    const svgData = new XMLSerializer().serializeToString(svgElement)
    const blob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'mindmap.svg'
    link.click()
    URL.revokeObjectURL(url)
    statusMessage.value = '已下载 SVG'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '下载失败'
  }
}

async function downloadPng() {
  const svgElement = svgRef.value
  if (!svgElement) {
    errorMessage.value = '没有可下载的内容'
    return
  }

  try {
    const svgData = new XMLSerializer().serializeToString(svgElement)
    let svgContent = svgData
    if (!svgContent.includes('xmlns="http://www.w3.org/2000/svg"')) {
      svgContent = svgContent.replace('<svg', '<svg xmlns="http://www.w3.org/2000/svg"')
    }

    const canvas = document.createElement('canvas')
    const context = canvas.getContext('2d')
    if (!context) {
      throw new Error('无法创建画布')
    }

    const img = new Image()
    const svgBlob = new Blob([svgContent], { type: 'image/svg+xml;charset=utf-8' })
    const url = URL.createObjectURL(svgBlob)

    await new Promise<void>((resolve, reject) => {
      img.onload = () => resolve()
      img.onerror = () => reject(new Error('图片加载失败'))
      img.src = url
    })

    const padding = 40
    canvas.width = img.width + padding * 2
    canvas.height = img.height + padding * 2

    context.fillStyle = '#ffffff'
    context.fillRect(0, 0, canvas.width, canvas.height)
    context.drawImage(img, padding, padding)

    URL.revokeObjectURL(url)

    canvas.toBlob((pngBlob) => {
      if (!pngBlob) {
        errorMessage.value = 'PNG 生成失败'
        return
      }

      const link = document.createElement('a')
      link.href = URL.createObjectURL(pngBlob)
      link.download = 'mindmap.png'
      link.click()
      URL.revokeObjectURL(link.href)

      statusMessage.value = '已下载 PNG'
      errorMessage.value = ''
    }, 'image/png', 0.95)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '下载失败'
  }
}

function handleZoomIn() {
  if (markmapInstance) {
    const currentScale = markmapInstance.transform?.k ?? 1
    const newScale = currentScale * 1.2
    markmapInstance.rescale(newScale)
    zoomLevel.value = Math.round(newScale * 100)
  }
}

function handleZoomOut() {
  if (markmapInstance) {
    const currentScale = markmapInstance.transform?.k ?? 1
    const newScale = currentScale / 1.2
    markmapInstance.rescale(newScale)
    zoomLevel.value = Math.round(newScale * 100)
  }
}

function handleReset() {
  if (markmapInstance) {
    markmapInstance.fit()
    zoomLevel.value = 100
  }
}

function toggleFullscreen() {
  isFullscreen.value = !isFullscreen.value
}

watch(inputMarkdown, () => {
  renderMarkmap()
})

onMounted(() => {
  setTimeout(() => {
    renderMarkmap()
  }, 100)
})
</script>

<template>
  <section class="tool-page" :class="{ 'is-canvas-fullscreen': isFullscreen }">
    <header v-if="!isFullscreen" class="tool-header">
      <div>
        <p class="eyebrow">Mind Map</p>
        <h1>Markmap 脑图</h1>
        <p class="summary">用 Markdown 语法快速生成脑图，左右分栏，实时预览。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="useSample">示例</el-button>
        <el-button @click="clearAll">清空</el-button>
        <el-button @click="copyMarkdown">复制</el-button>
        <el-button @click="downloadSvg">下载 SVG</el-button>
        <el-button type="primary" @click="downloadPng">下载 PNG</el-button>
      </div>
    </header>

    <section class="markmap-layout" :class="{ 'is-fullscreen': isFullscreen }">
      <aside v-if="!isFullscreen" class="markmap-sidebar">
        <div class="markmap-sidebar__header">
          <strong>Markdown</strong>
          <span>{{ inputMarkdown.split('\n').length }} 行</span>
        </div>
        <div class="markmap-sidebar__body">
          <el-input
            v-model="inputMarkdown"
            type="textarea"
            resize="none"
            spellcheck="false"
            placeholder="输入 Markdown 格式的文本，层级结构将自动生成脑图"
          />
        </div>
      </aside>

      <section class="markmap-canvas-wrap">
        <div class="markmap-canvas">
          <svg ref="svgRef" class="markmap-svg" />
          <div v-if="!inputMarkdown.trim()" class="markmap-canvas__empty">
            <p>输入 Markdown 后脑图将自动生成</p>
          </div>
        </div>
        <div class="markmap-canvas-toolbar">
          <CanvasToolbar
            :zoom-text="`${zoomLevel}%`"
            @zoom-in="handleZoomIn()"
            @zoom-out="handleZoomOut()"
            @reset="handleReset()"
            @fullscreen="toggleFullscreen()"
          />
        </div>
      </section>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
