import { nextTick, onBeforeUnmount, onMounted, ref, watch, type Ref } from 'vue'
import { Canvas as FabricCanvas, Textbox } from 'fabric'
import {
  annotationTools,
  getAllAnnotations,
  updateAnnotation,
  type AnnotationTool,
  type PDFAnnotation,
} from '@/lib/pdf_annotations'

type OverlayOptions = {
  currentPage: Ref<number>
  pdfCanvasEl: Ref<HTMLCanvasElement | null>
  pageRenderVersion: Ref<number>
}

type FabricTextboxWithId = Textbox & {
  annotationId?: string
}

export function usePdfAnnotationOverlay({ currentPage, pdfCanvasEl, pageRenderVersion }: OverlayOptions) {
  const annotationHostEl = ref<HTMLDivElement | null>(null)
  const activeTool = ref<AnnotationTool>('select')
  const annotationCount = ref(0)
  const overlayReady = ref(false)

  let fabricCanvas: FabricCanvas | null = null
  let resizeObserver: ResizeObserver | null = null
  let loadToken = 0
  let isHydrating = false
  let resizeFrame = 0

  function getOverlaySize() {
    const canvas = pdfCanvasEl.value
    if (!canvas) return { width: 0, height: 0 }

    const rect = canvas.getBoundingClientRect()
    const width = Math.round(rect.width || canvas.clientWidth || canvas.width)
    const height = Math.round(rect.height || canvas.clientHeight || canvas.height)

    return {
      width,
      height,
    }
  }

  function syncOverlaySize() {
    if (!fabricCanvas) return false

    const { width, height } = getOverlaySize()
    if (!width || !height) return false

    fabricCanvas.setDimensions({ width, height })
    fabricCanvas.requestRenderAll()
    return true
  }

  function createTextboxObject(annotation: PDFAnnotation, width: number, height: number) {
    const textbox = new Textbox(annotation.data.text, {
      left: annotation.data.left * width,
      top: annotation.data.top * height,
      width: annotation.data.width * width,
      fontSize: Math.max(12, annotation.data.fontSize * height),
      fill: annotation.data.fill,
      angle: annotation.data.angle ?? 0,
      editable: true,
      borderColor: '#0f172a',
      cornerColor: '#0f172a',
      cornerStrokeColor: '#ffffff',
      transparentCorners: false,
    }) as FabricTextboxWithId

    textbox.annotationId = annotation.id
    return textbox
  }

  function serializeTextbox(textbox: FabricTextboxWithId, width: number, height: number): PDFAnnotation {
    return {
      id: textbox.annotationId ?? crypto.randomUUID(),
      type: 'textbox',
      data: {
        text: textbox.text ?? '',
        left: (textbox.left ?? 0) / width,
        top: (textbox.top ?? 0) / height,
        width: ((textbox.width ?? 0) * (textbox.scaleX ?? 1)) / width,
        fontSize: (textbox.fontSize ?? 16) / height,
        fill: typeof textbox.fill === 'string' ? textbox.fill : '#0f172a',
        angle: textbox.angle ?? 0,
      },
    }
  }

  async function pushAnnotationUpdates() {
    if (!fabricCanvas || isHydrating) return

    const { width, height } = getOverlaySize()
    if (!width || !height) return

    const objects = fabricCanvas.getObjects()
    annotationCount.value = objects.length

    for (const object of objects) {
      if (object.type !== 'textbox') continue
      const serialized = serializeTextbox(object as FabricTextboxWithId, width, height)
      await updateAnnotation(currentPage.value, serialized)
    }
  }

  async function reloadAnnotations(page: number) {
    if (!fabricCanvas) return

    if (!syncOverlaySize()) {
      annotationCount.value = 0
      return
    }

    const localToken = ++loadToken
    const { width, height } = getOverlaySize()
    if (!width || !height) return

    isHydrating = true

    try {
      const annotations = await getAllAnnotations(page)
      if (localToken !== loadToken || !fabricCanvas) return

      fabricCanvas.clear()

      for (const annotation of annotations) {
        fabricCanvas.add(createTextboxObject(annotation, width, height))
      }

      annotationCount.value = annotations.length
      fabricCanvas.requestRenderAll()
    } finally {
      isHydrating = false
    }
  }

  function setActiveTool(tool: AnnotationTool) {
    activeTool.value = tool
    if (!fabricCanvas) return

    fabricCanvas.selection = true
    fabricCanvas.skipTargetFind = false
    fabricCanvas.requestRenderAll()
  }

  async function addTextbox() {
    if (!fabricCanvas) return
    if (!syncOverlaySize()) return

    const { width, height } = getOverlaySize()
    const textbox = new Textbox('Text', {
      left: width * 0.16,
      top: height * 0.14,
      width: width * 0.3,
      fontSize: Math.max(18, height * 0.032),
      fill: '#111827',
      editable: true,
      borderColor: '#0f172a',
      cornerColor: '#0f172a',
      cornerStrokeColor: '#ffffff',
      transparentCorners: false,
    }) as FabricTextboxWithId

    textbox.annotationId = crypto.randomUUID()
    fabricCanvas.add(textbox)
    fabricCanvas.setActiveObject(textbox)
    fabricCanvas.requestRenderAll()
    textbox.enterEditing()
    textbox.selectAll()
    setActiveTool('select')
    await pushAnnotationUpdates()
  }

  function scheduleResizeSync() {
    if (resizeFrame) cancelAnimationFrame(resizeFrame)
    resizeFrame = requestAnimationFrame(() => {
      resizeFrame = 0
      if (!syncOverlaySize()) return
      void reloadAnnotations(currentPage.value)
    })
  }

  function mountFabricCanvas() {
    const host = annotationHostEl.value
    if (!host || fabricCanvas) return

    const canvasEl = document.createElement('canvas')
    canvasEl.className = 'block h-full w-full'
    host.replaceChildren(canvasEl)

    fabricCanvas = new FabricCanvas(canvasEl, {
      preserveObjectStacking: true,
      selection: true,
      containerClass: 'paperlink-annotation-overlay',
    })

    const wrapperEl = (fabricCanvas as unknown as { wrapperEl?: HTMLDivElement }).wrapperEl
    if (wrapperEl) {
      wrapperEl.style.position = 'absolute'
      wrapperEl.style.inset = '0'
      wrapperEl.style.width = '100%'
      wrapperEl.style.height = '100%'
    }

    fabricCanvas.on('object:added', () => {
      if (isHydrating) return
      void pushAnnotationUpdates()
    })
    fabricCanvas.on('object:modified', () => {
      if (isHydrating) return
      void pushAnnotationUpdates()
    })
    fabricCanvas.on('object:removed', () => {
      if (isHydrating) return
      void pushAnnotationUpdates()
    })
    fabricCanvas.on('text:changed', () => {
      if (isHydrating) return
      void pushAnnotationUpdates()
    })

    overlayReady.value = true
    setActiveTool('select')
    syncOverlaySize()
    void reloadAnnotations(currentPage.value)
  }

  onMounted(async () => {
    await nextTick()
    mountFabricCanvas()

    if (typeof ResizeObserver !== 'undefined') {
      resizeObserver = new ResizeObserver(() => {
        scheduleResizeSync()
      })

      if (pdfCanvasEl.value) {
        resizeObserver.observe(pdfCanvasEl.value)
      }
    }
  })

  watch(pdfCanvasEl, (canvas, prevCanvas) => {
    if (prevCanvas && resizeObserver) {
      resizeObserver.unobserve(prevCanvas)
    }
    if (canvas && resizeObserver) {
      resizeObserver.observe(canvas)
    }
    mountFabricCanvas()
    scheduleResizeSync()
  })

  watch(currentPage, (page) => {
    void reloadAnnotations(page)
  })

  watch(pageRenderVersion, () => {
    scheduleResizeSync()
  })

  onBeforeUnmount(() => {
    if (resizeFrame) cancelAnimationFrame(resizeFrame)
    resizeObserver?.disconnect()
    resizeObserver = null
    overlayReady.value = false

    if (fabricCanvas) {
      fabricCanvas.dispose()
      fabricCanvas = null
    }
  })

  return {
    annotationHostEl,
    annotationCount,
    annotationTools,
    activeTool,
    overlayReady,
    setActiveTool,
    addTextbox,
  }
}
