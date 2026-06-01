<script setup lang="ts">
import mermaid from 'mermaid'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import CanvasToolbar from '../../components/CanvasToolbar.vue'
import { useCanvasZoom } from '../../composables/useCanvasZoom'

type MermaidTheme = 'default' | 'neutral' | 'dark' | 'forest' | 'base'
type MermaidSecurityLevel = 'strict' | 'loose'
type PanelKey = 'code' | 'config' | 'examples'

interface DiagramExample {
  title: string
  description: string
  code: string
}

const source = ref('')
const previewSvg = ref('')
const statusMessage = ref('')
const errorMessage = ref('')
const isRendering = ref(false)
const activePanel = ref<PanelKey>('code')
const mermaidTheme = ref<MermaidTheme>('default')
const securityLevel = ref<MermaidSecurityLevel>('strict')
const previewViewport = ref<HTMLElement | null>(null)
const isFullscreen = ref(false)
let renderTimer: number | undefined
let renderVersion = 0

const appTheme = ref<'dark' | 'light'>(document.documentElement.getAttribute('data-theme') as 'dark' | 'light' || 'dark')

const themeObserver = new MutationObserver(() => {
  const newTheme = document.documentElement.getAttribute('data-theme') as 'dark' | 'light' || 'dark'
  if (newTheme !== appTheme.value) {
    appTheme.value = newTheme
    renderMermaid()
  }
})
themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })

const {
  zoomText,
  transform: previewTransform,
  isDragging: isPreviewDragging,
  reset: resetPreviewTransform,
  zoomIn,
  zoomOut,
  handleWheel: handlePreviewWheel,
  startDrag: startPreviewDrag,
  moveDrag: movePreviewDrag,
  stopDrag: stopPreviewDrag,
} = useCanvasZoom()

const effectiveMermaidTheme = computed<MermaidTheme>(() => {
  if (mermaidTheme.value !== 'default') return mermaidTheme.value
  return appTheme.value === 'dark' ? 'dark' : 'default'
})

const themeOptions: Array<{ label: string; value: MermaidTheme }> = [
  { label: '默认', value: 'default' },
  { label: '中性', value: 'neutral' },
  { label: '深色', value: 'dark' },
  { label: '森林', value: 'forest' },
  { label: '基础', value: 'base' },
]

const securityOptions: Array<{ label: string; value: MermaidSecurityLevel }> = [
  { label: '严格', value: 'strict' },
  { label: '宽松', value: 'loose' },
]

const panelItems: Array<{ key: PanelKey; title: string; hint: string }> = [
  { key: 'code', title: '代码', hint: '编辑 Mermaid' },
  { key: 'config', title: '配置', hint: '主题和安全' },
  { key: 'examples', title: '示例', hint: '常用图表' },
]

const diagramExamples: DiagramExample[] = [
  {
    title: '流程图',
    description: '适合业务流程和决策路径',
    code: [
      'flowchart TD',
      '  A[开始] --> B{选择工具}',
      '  B --> C[文本处理]',
      '  B --> D[二维码]',
      '  B --> E[账号密钥]',
      '  E --> F[本地加密保存]',
    ].join('\n'),
  },
  {
    title: '时序图',
    description: '适合接口调用和链路说明',
    code: [
      'sequenceDiagram',
      '  participant U as 用户',
      '  participant A as MyTools',
      '  participant S as 本地存储',
      '  U->>A: 保存账号密钥',
      '  A->>A: AES-GCM 加密',
      '  A->>S: 写入本地数据',
      '  S-->>A: 保存成功',
      '  A-->>U: 返回结果',
    ].join('\n'),
  },
  {
    title: '状态图',
    description: '适合状态流转和任务阶段',
    code: [
      'stateDiagram-v2',
      '  [*] --> 草稿',
      '  草稿 --> 校验中',
      '  校验中 --> 已发布',
      '  校验中 --> 草稿: 修改',
      '  已发布 --> [*]',
    ].join('\n'),
  },
  {
    title: '甘特图',
    description: '适合简单计划和时间线',
    code: [
      'gantt',
      '  title 工具开发计划',
      '  dateFormat  YYYY-MM-DD',
      '  section 基础',
      '  项目结构 :done, 2026-04-10, 2d',
      '  Mermaid 编辑器 :active, 2026-04-17, 2d',
      '  section 优化',
      '  导出和配置 :2026-04-19, 2d',
    ].join('\n'),
  },
]

const stats = computed(() => {
  const text = source.value
  return {
    chars: text.length,
    lines: text ? text.split(/\r?\n/).length : 0,
  }
})

mermaid.initialize({
  startOnLoad: false,
  securityLevel: securityLevel.value,
  theme: effectiveMermaidTheme.value,
})

function useSample() {
  applyExample(diagramExamples[0])
}

function applyExample(example: DiagramExample) {
  source.value = example.code
  activePanel.value = 'code'
  statusMessage.value = ''
  errorMessage.value = ''
}

function clearSource() {
  source.value = ''
  previewSvg.value = ''
  statusMessage.value = ''
  errorMessage.value = ''
}

async function copySource() {
  if (!source.value.trim()) {
    statusMessage.value = '没有可复制的内容'
    return
  }

  try {
    await navigator.clipboard.writeText(source.value)
    statusMessage.value = '已复制 Mermaid 源码'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

async function copySvg() {
  if (!previewSvg.value) {
    statusMessage.value = '没有可复制的 SVG'
    return
  }

  try {
    await navigator.clipboard.writeText(previewSvg.value)
    statusMessage.value = '已复制 SVG'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

function downloadSvg() {
  if (!previewSvg.value) {
    statusMessage.value = '没有可下载的 SVG'
    return
  }

  const blob = new Blob([previewSvg.value], { type: 'image/svg+xml;charset=utf-8' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = 'mermaid-diagram.svg'
  link.click()
  URL.revokeObjectURL(link.href)
  statusMessage.value = '已下载 SVG'
  errorMessage.value = ''
}

async function downloadPng() {
  if (!previewSvg.value) {
    statusMessage.value = '没有可下载的 PNG'
    return
  }

  try {
    let svgContent = previewSvg.value
    if (!svgContent.includes('xmlns="http://www.w3.org/2000/svg"')) {
      svgContent = svgContent.replace('<svg', '<svg xmlns="http://www.w3.org/2000/svg"')
    }

    const svgDocument = new DOMParser().parseFromString(svgContent, 'image/svg+xml')
    const svgElement = svgDocument.querySelector('svg')
    const viewBox = svgElement?.getAttribute('viewBox')?.split(/\s+/).map(Number)
    let width = Number(svgElement?.getAttribute('width')) || (viewBox?.[2] || 1200)
    let height = Number(svgElement?.getAttribute('height')) || (viewBox?.[3] || 800)

    if (!svgElement?.getAttribute('width')) {
      svgElement?.setAttribute('width', String(width))
    }
    if (!svgElement?.getAttribute('height')) {
      svgElement?.setAttribute('height', String(height))
    }

    const updatedSvgContent = new XMLSerializer().serializeToString(svgDocument)

    const scale = Math.min(window.devicePixelRatio || 1, 2)
    const canvas = document.createElement('canvas')
    const context = canvas.getContext('2d')

    if (!context) {
      throw new Error('无法创建图片画布')
    }

    canvas.width = Math.ceil(width * scale)
    canvas.height = Math.ceil(height * scale)
    context.scale(scale, scale)
    context.fillStyle = '#ffffff'
    context.fillRect(0, 0, width, height)

    const blob = new Blob([updatedSvgContent], { type: 'image/svg+xml;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const image = new Image()

    await new Promise<void>((resolve, reject) => {
      image.onload = () => resolve()
      image.onerror = () => reject(new Error('PNG 生成失败'))
      image.src = url
    })

    context.drawImage(image, 0, 0, width, height)
    URL.revokeObjectURL(url)

    canvas.toBlob((pngBlob) => {
      if (!pngBlob) {
        errorMessage.value = 'PNG 生成失败'
        return
      }
      downloadBlob(pngBlob, 'mermaid-diagram.png')
      statusMessage.value = '已下载 PNG'
      errorMessage.value = ''
    }, 'image/png', 0.95)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'PNG 生成失败'
  }
}

function downloadBlob(blob: Blob, filename: string) {
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = filename
  link.click()
  URL.revokeObjectURL(link.href)
}

function toggleFullscreen() {
  isFullscreen.value = !isFullscreen.value
}

async function renderMermaid() {
  const code = source.value.trim()
  const currentVersion = ++renderVersion

  if (!code) {
    previewSvg.value = ''
    errorMessage.value = ''
    isRendering.value = false
    return
  }

  isRendering.value = true

  try {
    await nextTick()
    mermaid.initialize({
      startOnLoad: false,
      securityLevel: securityLevel.value,
      theme: effectiveMermaidTheme.value,
    })
    const id = `mermaid-${Date.now()}-${currentVersion}`
    const { svg } = await mermaid.render(id, code)

    if (currentVersion === renderVersion) {
      previewSvg.value = svg
      errorMessage.value = ''
      resetPreviewTransform()
    }
  } catch (error) {
    if (currentVersion === renderVersion) {
      previewSvg.value = ''
      errorMessage.value = error instanceof Error ? error.message : 'Mermaid 语法解析失败'
    }
  } finally {
    if (currentVersion === renderVersion) {
      isRendering.value = false
    }
  }
}

watch(source, () => {
  if (renderTimer) {
    window.clearTimeout(renderTimer)
  }
  renderTimer = window.setTimeout(() => {
    renderMermaid()
  }, 260)
})

watch([mermaidTheme, securityLevel], () => {
  renderMermaid()
})

onMounted(() => {
  useSample()
  renderMermaid()
})

onBeforeUnmount(() => {
  if (renderTimer) {
    window.clearTimeout(renderTimer)
  }
  themeObserver.disconnect()
})
</script>

<template>
  <section class="tool-page" :class="{ 'is-canvas-fullscreen': isFullscreen }">
    <header v-if="!isFullscreen" class="tool-header">
      <div>
        <p class="eyebrow">Diagram</p>
        <h1>Mermaid 编辑器</h1>
        <p class="summary">左侧编辑图表源码和配置，右侧实时预览 Mermaid 图形。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="copySource">复制源码</el-button>
        <el-button @click="copySvg">复制 SVG</el-button>
        <el-button @click="downloadSvg">下载 SVG</el-button>
        <el-button type="primary" @click="downloadPng">下载 PNG</el-button>
      </div>
    </header>

    <section class="mermaid-layout" :class="{ 'is-fullscreen': isFullscreen }">
      <aside v-if="!isFullscreen" class="mermaid-sidebar">
        <div class="mermaid-sidebar__tabs">
          <button
            v-for="panel in panelItems"
            :key="panel.key"
            class="mermaid-tab"
            :class="{ 'is-active': activePanel === panel.key }"
            type="button"
            @click="activePanel = panel.key"
          >
            <span class="mermaid-tab__icon">{{ panel.title.slice(0, 1) }}</span>
            <span>
              <strong>{{ panel.title }}</strong>
              <small>{{ panel.hint }}</small>
            </span>
          </button>
        </div>

        <div class="mermaid-sidebar__body">
          <div v-if="activePanel === 'code'" class="mermaid-panel mermaid-panel--code">
            <div class="mermaid-panel__head">
              <div>
                <strong>Mermaid 源码</strong>
                <span>{{ stats.lines }} 行 · {{ stats.chars }} 字符</span>
              </div>
              <el-button @click="clearSource">清空</el-button>
            </div>
            <el-input
              v-model="source"
              type="textarea"
              resize="none"
              clearable
              spellcheck="false"
              placeholder="输入 Mermaid 内容，例如 flowchart TD"
              @clear="clearSource"
            />
          </div>

          <div v-else-if="activePanel === 'config'" class="mermaid-panel mermaid-config">
            <div class="mermaid-panel__head">
              <div>
                <strong>渲染配置</strong>
                <span>{{ isRendering ? '正在渲染' : '实时生效' }}</span>
              </div>
            </div>

            <label class="mermaid-field">
              <span>主题</span>
              <el-select v-model="mermaidTheme" placeholder="选择主题">
                <el-option
                  v-for="option in themeOptions"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </el-select>
            </label>

            <label class="mermaid-field">
              <span>安全级别</span>
              <el-select v-model="securityLevel" placeholder="选择安全级别">
                <el-option
                  v-for="option in securityOptions"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </el-select>
            </label>

            <el-alert
              title="严格模式会过滤不安全内容，日常使用建议保持严格。"
              type="info"
              :closable="false"
              show-icon
            />
          </div>

          <div v-else class="mermaid-panel mermaid-examples">
            <button
              v-for="example in diagramExamples"
              :key="example.title"
              class="mermaid-example"
              type="button"
              @click="applyExample(example)"
            >
              <strong>{{ example.title }}</strong>
              <span>{{ example.description }}</span>
            </button>
          </div>
        </div>
      </aside>

      <section class="mermaid-canvas-wrap">
        <div
          ref="previewViewport"
          class="mermaid-canvas"
          :class="{ 'is-dragging': isPreviewDragging }"
          @wheel="handlePreviewWheel"
          @pointerdown="startPreviewDrag"
          @pointermove="movePreviewDrag"
          @pointerup="stopPreviewDrag"
          @pointercancel="stopPreviewDrag"
          @dblclick="resetPreviewTransform()"
        >
          <div
            v-if="previewSvg"
            class="mermaid-canvas__content"
            :style="previewTransform"
          >
            <div class="mermaid-canvas__svg" v-html="previewSvg"></div>
          </div>
          <div v-else class="mermaid-canvas__empty">
            <p>输入 Mermaid 后预览图形</p>
          </div>
        </div>
        <div class="mermaid-canvas-toolbar">
          <CanvasToolbar
            :zoom-text="isRendering ? '渲染中' : zoomText"
            @zoom-in="zoomIn()"
            @zoom-out="zoomOut()"
            @reset="resetPreviewTransform()"
            @fullscreen="toggleFullscreen()"
          />
        </div>
      </section>
    </section>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
