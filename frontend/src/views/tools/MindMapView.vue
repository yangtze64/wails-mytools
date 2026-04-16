<script setup lang="ts">
import { markRaw, nextTick, onBeforeUnmount, onMounted, ref, shallowRef } from 'vue'
import MindElixir, { type MindElixirData, type MindElixirInstance, type NodeObj } from 'mind-elixir'
import 'mind-elixir/style.css'

const container = ref<HTMLElement | null>(null)
const mind = shallowRef<MindElixirInstance | null>(null)
const statusMessage = ref('')
const errorMessage = ref('')

function createSampleData(): MindElixirData {
  return {
    nodeData: {
      id: 'root',
      topic: '产品规划',
      children: [
        {
          id: 'users',
          topic: '用户',
          children: [
            { id: 'developer', topic: '开发者' },
            { id: 'designer', topic: '设计师' },
          ],
        },
        {
          id: 'features',
          topic: '功能',
          children: [
            { id: 'text-tools', topic: '文本工具' },
            { id: 'qr-tools', topic: '二维码工具' },
            { id: 'mind-map-tool', topic: '思维导图' },
          ],
        },
        {
          id: 'release',
          topic: '发布',
          children: [
            { id: 'package', topic: '打包' },
            { id: 'test', topic: '测试' },
            { id: 'docs', topic: '文档' },
          ],
        },
      ],
    },
    direction: MindElixir.SIDE,
  }
}

async function initMindMap(data = createSampleData()) {
  await nextTick()

  if (!container.value) {
    return
  }

  mind.value?.destroy()
  mind.value = markRaw(new MindElixir({
    el: container.value,
    direction: MindElixir.SIDE,
    editable: true,
    contextMenu: true,
    toolBar: true,
    keypress: true,
    mouseSelectionButton: 2,
    draggable: true,
    newTopicName: '新节点',
    allowUndo: true,
    overflowHidden: false,
    alignment: 'nodes',
    scaleMin: 0.45,
    scaleMax: 2.4,
    theme: {
      name: 'mytools',
      palette: ['#256f5a', '#1d4ed8', '#9333ea', '#dc2626', '#ca8a04', '#0891b2'],
      cssVar: {
        '--node-gap-x': '42px',
        '--node-gap-y': '18px',
        '--main-gap-x': '72px',
        '--main-gap-y': '42px',
        '--main-color': '#1f2937',
        '--main-bgcolor': '#ffffff',
        '--main-bgcolor-transparent': 'rgba(255, 255, 255, 0.72)',
        '--color': '#1f2937',
        '--bgcolor': '#ffffff',
        '--selected': '#dff2e9',
        '--accent-color': '#256f5a',
        '--root-color': '#ffffff',
        '--root-bgcolor': '#256f5a',
        '--root-border-color': '#256f5a',
        '--root-radius': '8px',
        '--main-radius': '8px',
        '--topic-padding': '8px 12px',
        '--panel-color': '#374151',
        '--panel-bgcolor': '#ffffff',
        '--panel-border-color': '#e5e7eb',
        '--map-padding': '80px',
      },
    },
  }))

  mind.value.init(data)
  mind.value.scaleFit()
}

function resetSample() {
  initMindMap(createSampleData())
  statusMessage.value = '已重置示例'
  errorMessage.value = ''
}

function createBlank() {
  initMindMap(MindElixir.new('我的思维导图'))
  statusMessage.value = '已创建空白脑图'
  errorMessage.value = ''
}

async function exportPng() {
  if (!mind.value) {
    return
  }

  const blob = await mind.value.exportPng()
  if (!blob) {
    errorMessage.value = '导出 PNG 失败'
    return
  }

  downloadBlob(blob, 'mind-map.png')
  statusMessage.value = '已导出 PNG'
  errorMessage.value = ''
}

function exportSvg() {
  if (!mind.value) {
    return
  }

  const blob = mind.value.exportSvg()
  downloadBlob(blob, 'mind-map.svg')
  statusMessage.value = '已导出 SVG'
  errorMessage.value = ''
}

async function copyJson() {
  if (!mind.value) {
    return
  }

  try {
    await navigator.clipboard.writeText(JSON.stringify(cleanData(mind.value.getData()), null, 2))
    statusMessage.value = '已复制 JSON'
    errorMessage.value = ''
  } catch {
    errorMessage.value = '复制失败'
  }
}

function centerMap() {
  mind.value?.scaleFit()
}

function cleanData(data: MindElixirData): MindElixirData {
  return {
    ...data,
    nodeData: cleanNode(data.nodeData),
  }
}

function cleanNode(node: NodeObj): NodeObj {
  const { parent, ...nodeData } = node
  return {
    ...nodeData,
    children: node.children?.map(cleanNode),
  }
}

function downloadBlob(blob: Blob, filename: string) {
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = filename
  link.click()
  URL.revokeObjectURL(link.href)
}

onMounted(() => {
  initMindMap()
})

onBeforeUnmount(() => {
  mind.value?.destroy()
})
</script>

<template>
  <section class="tool-page">
    <header class="tool-header">
      <div>
        <p class="eyebrow">Diagram</p>
        <h1>思维导图</h1>
        <p class="summary">专业脑图编辑器，支持拖拽、右键菜单、快捷键、缩放和导出。</p>
      </div>
      <div class="tool-actions">
        <el-button @click="createBlank">空白</el-button>
        <el-button @click="resetSample">示例</el-button>
        <el-button @click="centerMap">适应画布</el-button>
        <el-button @click="copyJson">复制 JSON</el-button>
        <el-button @click="exportSvg">导出 SVG</el-button>
        <el-button type="primary" @click="exportPng">导出 PNG</el-button>
      </div>
    </header>

    <el-card class="mind-map-pro-card" shadow="never">
      <div ref="container" class="mind-map-pro"></div>
    </el-card>

    <el-alert v-if="errorMessage" class="copy-state" :title="errorMessage" type="error" :closable="false" show-icon />
    <el-alert v-else-if="statusMessage" class="copy-state" :title="statusMessage" type="success" :closable="false" show-icon />
  </section>
</template>
