import { computed, ref, watch } from 'vue'
import { toolRegistry, type ToolKey, type ToolMenuItem } from '../config/toolRegistry'

const STORAGE_KEY = 'mytools.toolTabs.v1'
const DEFAULT_TOOL: ToolKey = 'text-dedupe'
const validToolKeys = new Set<ToolKey>(toolRegistry.map((tool) => tool.key))

interface StoredToolTabs {
  activeTool?: ToolKey
  openTools?: ToolKey[]
}

function isToolKey(value: unknown): value is ToolKey {
  return typeof value === 'string' && validToolKeys.has(value as ToolKey)
}

function normalizeOpenTools(value: unknown): ToolKey[] {
  if (!Array.isArray(value)) {
    return [DEFAULT_TOOL]
  }

  const tools = value.filter(isToolKey)
  const uniqueTools = Array.from(new Set(tools))
  return uniqueTools.length > 0 ? uniqueTools : [DEFAULT_TOOL]
}

function readStoredTabs(): Required<StoredToolTabs> {
  if (typeof window === 'undefined') {
    return { activeTool: DEFAULT_TOOL, openTools: [DEFAULT_TOOL] }
  }

  try {
    const rawValue = window.localStorage.getItem(STORAGE_KEY)
    const parsedValue = rawValue ? (JSON.parse(rawValue) as StoredToolTabs) : null
    const openTools = normalizeOpenTools(parsedValue?.openTools)
    const activeTool = isToolKey(parsedValue?.activeTool) && openTools.includes(parsedValue.activeTool)
      ? parsedValue.activeTool
      : openTools[0]

    return { activeTool, openTools }
  } catch {
    return { activeTool: DEFAULT_TOOL, openTools: [DEFAULT_TOOL] }
  }
}

function persistTabs(openTools: ToolKey[], activeTool: ToolKey) {
  if (typeof window === 'undefined') {
    return
  }

  window.localStorage.setItem(STORAGE_KEY, JSON.stringify({ activeTool, openTools }))
}

export function useToolTabs() {
  const storedTabs = readStoredTabs()
  const openTools = ref<ToolKey[]>(storedTabs.openTools)
  const activeTool = ref<ToolKey>(storedTabs.activeTool)

  const openTabItems = computed<ToolMenuItem[]>(() =>
    openTools.value
      .map((key) => toolRegistry.find((tool) => tool.key === key))
      .filter((tool): tool is ToolMenuItem => Boolean(tool)),
  )

  const activeToolInfo = computed(() =>
    toolRegistry.find((tool) => tool.key === activeTool.value) ?? toolRegistry[0],
  )

  function openToolTab(key: ToolKey) {
    if (!validToolKeys.has(key)) {
      return
    }

    if (!openTools.value.includes(key)) {
      openTools.value = [...openTools.value, key]
    }
    activeTool.value = key
  }

  function closeToolTab(key: ToolKey) {
    if (openTools.value.length <= 1) {
      return
    }

    const currentIndex = openTools.value.indexOf(key)
    if (currentIndex === -1) {
      return
    }

    const nextOpenTools = openTools.value.filter((toolKey) => toolKey !== key)
    openTools.value = nextOpenTools

    if (activeTool.value === key) {
      activeTool.value = nextOpenTools[Math.max(0, currentIndex - 1)]
    }
  }

  function closeOtherToolTabs(key: ToolKey) {
    if (!openTools.value.includes(key)) {
      return
    }

    openTools.value = [key]
    activeTool.value = key
  }

  watch(
    [openTools, activeTool],
    ([nextOpenTools, nextActiveTool]) => persistTabs(nextOpenTools, nextActiveTool),
    { deep: true },
  )

  return {
    activeTool,
    activeToolInfo,
    openTabItems,
    openToolTab,
    closeToolTab,
    closeOtherToolTabs,
  }
}
