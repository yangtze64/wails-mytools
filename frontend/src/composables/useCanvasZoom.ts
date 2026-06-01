import { computed, ref, type CSSProperties } from 'vue'

export function useCanvasZoom(options?: { min?: number; max?: number }) {
  const scaleMin = options?.min ?? 0.2
  const scaleMax = options?.max ?? 4

  const scale = ref(1)
  const offsetX = ref(0)
  const offsetY = ref(0)
  const isDragging = ref(false)

  let dragStart = { x: 0, y: 0 }
  let dragOffsetStart = { x: 0, y: 0 }

  const zoomText = computed(() => `${Math.round(scale.value * 100)}%`)

  const transform = computed<CSSProperties>(() => ({
    transform: `translate(${offsetX.value}px, ${offsetY.value}px) scale(${scale.value})`,
    transformOrigin: 'center center',
  }))

  function clamp(value: number, min: number, max: number) {
    return Math.min(Math.max(value, min), max)
  }

  function reset() {
    scale.value = 1
    offsetX.value = 0
    offsetY.value = 0
  }

  function zoomIn(delta = 0.1) {
    scale.value = clamp(Number((scale.value + delta).toFixed(2)), scaleMin, scaleMax)
  }

  function zoomOut(delta = 0.1) {
    scale.value = clamp(Number((scale.value - delta).toFixed(2)), scaleMin, scaleMax)
  }

  function handleWheel(event: WheelEvent) {
    event.preventDefault()

    const currentScale = scale.value
    const nextScale = clamp(
      Number((currentScale * (event.deltaY > 0 ? 0.9 : 1.1)).toFixed(2)),
      scaleMin,
      scaleMax,
    )

    const viewport = event.currentTarget as HTMLElement
    if (!viewport || nextScale === currentScale) {
      scale.value = nextScale
      return
    }

    const rect = viewport.getBoundingClientRect()
    const pointerX = event.clientX - rect.left - rect.width / 2
    const pointerY = event.clientY - rect.top - rect.height / 2
    const scaleRatio = nextScale / currentScale

    offsetX.value = pointerX - (pointerX - offsetX.value) * scaleRatio
    offsetY.value = pointerY - (pointerY - offsetY.value) * scaleRatio
    scale.value = nextScale
  }

  function startDrag(event: PointerEvent) {
    if (event.button !== 0) return
    isDragging.value = true
    dragStart = { x: event.clientX, y: event.clientY }
    dragOffsetStart = { x: offsetX.value, y: offsetY.value }
    ;(event.currentTarget as HTMLElement).setPointerCapture(event.pointerId)
  }

  function moveDrag(event: PointerEvent) {
    if (!isDragging.value) return
    offsetX.value = dragOffsetStart.x + event.clientX - dragStart.x
    offsetY.value = dragOffsetStart.y + event.clientY - dragStart.y
  }

  function stopDrag(event: PointerEvent) {
    if (!isDragging.value) return
    isDragging.value = false
    ;(event.currentTarget as HTMLElement).releasePointerCapture(event.pointerId)
  }

  return {
    scale,
    offsetX,
    offsetY,
    isDragging,
    zoomText,
    transform,
    reset,
    zoomIn,
    zoomOut,
    handleWheel,
    startDrag,
    moveDrag,
    stopDrag,
  }
}
