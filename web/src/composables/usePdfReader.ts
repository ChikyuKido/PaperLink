import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import * as pdfjsLib from 'pdfjs-dist/build/pdf.mjs'
import pdfWorker from 'pdfjs-dist/build/pdf.worker.mjs?worker'
import { apiFetch } from '@/auth/api'
import type { CollabClientMessage, CollabServerMessage, CollabStatus } from '@/lib/pdf_collab'

pdfjsLib.GlobalWorkerOptions.workerPort = new pdfWorker()
type CollabMessageListener = (message: CollabServerMessage) => void

export function usePdfReader() {
  const route = useRoute()
  const pdfID = computed(() => String(route.params.id ?? ''))

  const pageCount = ref(0)
  const currentPage = ref(1)
  const canvasEl = ref<HTMLCanvasElement | null>(null)
  const thumbnailScrollEl = ref<HTMLDivElement | null>(null)
  const thumbnails = ref<(string | null)[]>([])
  const THUMB_BATCH_SIZE = 25
  const thumbnailBatchLocks = reactive<Record<string, boolean>>({})
  const thumbnailBatchCache = reactive<Record<string, boolean>>({})
  const readerError = ref<string | null>(null)
  const collabStatus = ref<CollabStatus>('idle')
  const collabError = ref<string | null>(null)
  const pageRenderVersion = ref(0)

  let keydownHandler: ((e: KeyboardEvent) => void) | null = null
  let currentRenderTask: { cancel: () => void; promise: Promise<void> } | null = null
  const pageCache = reactive<Record<number, Uint8Array>>({})
  const fetchLocks = reactive<Record<string, boolean>>({})
  let initToken = 0
  let renderToken = 0
  let collabToken = 0
  let collabSocket: WebSocket | null = null
  const collabListeners = new Set<CollabMessageListener>()

  function clearThumbnailBatchState() {
    for (const k of Object.keys(thumbnailBatchLocks)) delete thumbnailBatchLocks[k]
    for (const k of Object.keys(thumbnailBatchCache)) delete thumbnailBatchCache[k]
  }

  function clearThumbnails() {
    thumbnails.value.forEach((url) => {
      if (url) URL.revokeObjectURL(url)
    })
    thumbnails.value = []
  }

  function disposeCurrentPdf() {
    renderToken++
    pageRenderVersion.value++
    if (currentRenderTask) {
      try {
        currentRenderTask.cancel()
      } catch {
      }
      currentRenderTask = null
    }
    for (const page of Object.keys(pageCache)) delete pageCache[Number(page)]
    for (const key of Object.keys(fetchLocks)) delete fetchLocks[key]
    clearThumbnails()
    clearThumbnailBatchState()
  }

  function setCollabDisconnected(message: string | null = null) {
    collabStatus.value = 'disconnected'
    collabError.value = message
  }

  function closeCollabConnection() {
    collabToken++
    collabStatus.value = 'idle'
    collabError.value = null

    if (!collabSocket) return

    const socket = collabSocket
    collabSocket = null
    socket.onopen = null
    socket.onmessage = null
    socket.onerror = null
    socket.onclose = null

    if (socket.readyState === WebSocket.CONNECTING || socket.readyState === WebSocket.OPEN) {
      socket.close()
    }
  }

  function emitCollabMessage(message: CollabServerMessage) {
    for (const listener of Array.from(collabListeners)) {
      listener(message)
    }
  }

  function subscribeCollabMessages(listener: CollabMessageListener) {
    collabListeners.add(listener)
    return () => {
      collabListeners.delete(listener)
    }
  }

  function sendCollabMessage(message: CollabClientMessage) {
    if (!collabSocket || collabSocket.readyState !== WebSocket.OPEN) {
      return false
    }

    collabSocket.send(JSON.stringify(message))
    return true
  }

  function requestPageAnnotations(page: number) {
    return sendCollabMessage({
      type: 'annotations:get',
      page,
    })
  }

  function createAnnotation(annotation: CollabClientMessage['annotation']) {
    if (!annotation) return false
    return sendCollabMessage({
      type: 'annotation:create',
      annotation,
    })
  }

  function updateAnnotation(annotation: CollabClientMessage['annotation']) {
    if (!annotation) return false
    return sendCollabMessage({
      type: 'annotation:update',
      annotation,
    })
  }

  function moveAnnotation(annotation: CollabClientMessage['annotation']) {
    if (!annotation) return false
    return sendCollabMessage({
      type: 'annotation:move',
      annotation,
    })
  }

  function deleteAnnotation(annotationId: number) {
    return sendCollabMessage({
      type: 'annotation:delete',
      annotationId,
    })
  }

  async function connectCollab() {
    const documentID = pdfID.value
    if (!documentID) {
      closeCollabConnection()
      return
    }

    closeCollabConnection()
    const token = collabToken
    collabStatus.value = 'connecting'
    collabError.value = null

    try {
      const res = await apiFetch(`/api/v1/pdfws/create/${documentID}`)
      if (!res.ok) {
        setCollabDisconnected('Live sync unavailable.')
        return
      }

      const body = await res.json().catch(() => null)
      const singleUseToken = body?.data?.token
      if (typeof singleUseToken !== 'string' || singleUseToken.length === 0) {
        setCollabDisconnected('Live sync unavailable.')
        return
      }

      if (token !== collabToken) return

      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const socketURL = `${protocol}//${window.location.host}/api/v1/pdfws/connect/${encodeURIComponent(documentID)}?token=${encodeURIComponent(singleUseToken)}`
      const socket = new WebSocket(socketURL)
      collabSocket = socket

      socket.onopen = () => {
        if (token !== collabToken || collabSocket !== socket) return
        collabStatus.value = 'connected'
        collabError.value = null
      }

      socket.onmessage = (event) => {
        if (token !== collabToken || collabSocket !== socket) return

        try {
          const message = JSON.parse(String(event.data ?? '{}')) as CollabServerMessage
          if (message?.type === 'error' && typeof message.error === 'string') {
            collabError.value = message.error
          }
          emitCollabMessage(message)
        } catch {
        }
      }

      socket.onerror = () => {
        if (token !== collabToken || collabSocket !== socket) return
        setCollabDisconnected('Live sync unavailable.')
      }

      socket.onclose = () => {
        if (token !== collabToken || collabSocket !== socket) return
        collabSocket = null
        if (collabStatus.value !== 'disconnected') {
          setCollabDisconnected('Live sync disconnected.')
        }
      }
    } catch (err) {
      if (token !== collabToken) return
      console.error('Failed to connect to collaboration websocket', err)
      setCollabDisconnected('Live sync unavailable.')
    }
  }

  async function loadDocument() {
    const res = await apiFetch(`/api/v1/document/get/${pdfID.value}`)
    if (!res.ok) {
      throw new Error('Failed to load document metadata')
    }

    const doc = await res.json().catch(() => null)
    pageCount.value = Number(doc?.file?.pages ?? 0)
    thumbnails.value = Array.from({ length: pageCount.value }, () => null)
  }

  async function loadPDFDocument() {}

  async function fetchPages(start: number, end: number) {
    start = Math.max(1, start)
    end = Math.min(pageCount.value, end)
    if (start > end) return

    const key = `${start}-${end}`
    if (fetchLocks[key]) return
    fetchLocks[key] = true

    try {
      const res = await apiFetch(`/api/v1/pdf/${pdfID.value}/${start}-${end}`)
      if (!res.ok) return

      const buf = await res.arrayBuffer()
      const bytes = new Uint8Array(buf)

      if (bytes[0] === 0) {
        pageCache[start] = bytes.slice(1)
        return
      }

      let offset = 1
      let pageNum = start
      while (offset < bytes.length && pageNum <= end) {
        const size = Number(new DataView(bytes.buffer, offset, 8).getBigUint64(0))
        offset += 8
        if (size <= 0 || offset + size > bytes.length) break
        pageCache[pageNum] = bytes.slice(offset, offset + size)
        offset += size
        pageNum++
      }
    } finally {
      fetchLocks[key] = false
    }
  }

  async function fetchThumbnailBatch(startIndex: number) {
    const localPdfID = pdfID.value
    if (pageCount.value === 0 || startIndex >= pageCount.value) return

    const start = Math.max(0, startIndex)
    const end = Math.min(pageCount.value - 1, start + THUMB_BATCH_SIZE - 1)
    const key = `${start}-${end}`
    if (thumbnailBatchCache[key] || thumbnailBatchLocks[key]) return
    thumbnailBatchLocks[key] = true

    try {
      const res = await apiFetch(`/api/v1/pdf/thumbnails/${localPdfID}/${start}-${end}`)
      if (!res.ok) return

      const buf = await res.arrayBuffer()
      const bytes = new Uint8Array(buf)
      const dv = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength)

      let offset = 0
      let pageIndex = start

      while (offset + 8 <= bytes.length && pageIndex <= end) {
        const size = Number(dv.getBigUint64(offset, true))
        offset += 8
        if (size <= 0 || offset + size > bytes.length) break

        const pngBytes = bytes.slice(offset, offset + size)
        offset += size

        if (localPdfID !== pdfID.value) break
        const url = URL.createObjectURL(new Blob([pngBytes], { type: 'image/png' }))
        if (thumbnails.value[pageIndex]) {
          URL.revokeObjectURL(thumbnails.value[pageIndex]!)
        }
        thumbnails.value[pageIndex] = url
        pageIndex++
      }

      thumbnailBatchCache[key] = true
    } finally {
      thumbnailBatchLocks[key] = false
    }
  }

  function ensureThumbnailBatchForPage(page: number) {
    const idx = Math.max(0, page - 1)
    const batchStart = Math.floor(idx / THUMB_BATCH_SIZE) * THUMB_BATCH_SIZE
    void fetchThumbnailBatch(batchStart)
  }

  function ensureThumbnailBatchesForViewport() {
    const el = thumbnailScrollEl.value
    if (!el) return

    const itemHeight = 128
    const firstVisiblePage = Math.max(1, Math.floor(el.scrollTop / itemHeight) + 1)
    const lastVisiblePage = Math.min(
      pageCount.value,
      Math.ceil((el.scrollTop + el.clientHeight) / itemHeight),
    )

    ensureThumbnailBatchForPage(firstVisiblePage)
    ensureThumbnailBatchForPage(lastVisiblePage)
  }

  function onThumbnailScroll() {
    ensureThumbnailBatchesForViewport()
  }

  function ensureSurrounding(n: number) {
    const preloadBefore = Math.max(1, n - 1)
    const preloadAfter = Math.min(pageCount.value, n + 2)
    const ranges: Array<[number, number]> = []
    let rangeStart: number | null = null

    for (let i = preloadBefore; i <= preloadAfter; i++) {
      if (!pageCache[i]) {
        if (rangeStart === null) rangeStart = i
      } else if (rangeStart !== null) {
        ranges.push([rangeStart, i - 1])
        rangeStart = null
      }
    }

    if (rangeStart !== null) {
      ranges.push([rangeStart, preloadAfter])
    }

    for (const [start, end] of ranges) {
      void fetchPages(start, end)
    }
  }

  function scrollThumbnailIntoView(page: number) {
    void nextTick(() => {
      const el = thumbnailScrollEl.value?.querySelector<HTMLElement>(`[data-page="${page}"]`)
      el?.scrollIntoView({ block: 'nearest', inline: 'nearest' })
    })
  }

  async function renderPage(n: number) {
    const token = ++renderToken
    const canvas = canvasEl.value
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return
    if (token !== renderToken) return

    if (currentRenderTask) {
      try {
        currentRenderTask.cancel()
      } catch {
      }
      currentRenderTask = null
    }

    if (!pageCache[n]) {
      await fetchPages(n, n)
    }
    if (token !== renderToken || !pageCache[n]) {
      return
    }

    const pageBytes = pageCache[n].slice()
    const pdf = await pdfjsLib.getDocument({ data: pageBytes }).promise
    if (token !== renderToken) {
      void pdf.destroy()
      return
    }
    const page = await pdf.getPage(1)
    if (token !== renderToken) {
      void pdf.destroy()
      return
    }
    const viewport = page.getViewport({ scale: 1.5 })

    canvas.width = viewport.width
    canvas.height = viewport.height

    const renderTask = page.render({ canvasContext: ctx, viewport })
    currentRenderTask = renderTask as { cancel: () => void; promise: Promise<void> }
    await renderTask.promise
    void pdf.destroy()
    if (token !== renderToken) return

    pageRenderVersion.value++
    ensureSurrounding(n)
  }

  function go(n: number) {
    if (pageCount.value === 0) return
    n = Math.min(pageCount.value, Math.max(1, n))
    currentPage.value = n
    ensureThumbnailBatchForPage(n)
    scrollThumbnailIntoView(n)
    void renderPage(n)
  }

  const goFirst = () => go(1)
  const goLast = () => go(pageCount.value)
  const prevPage = () => go(Math.max(currentPage.value - 1, 1))
  const nextPage = () => go(Math.min(currentPage.value + 1, pageCount.value))

  async function initializeReader() {
    const token = ++initToken
    readerError.value = null
    currentPage.value = 1
    pageCount.value = 0
    disposeCurrentPdf()

    try {
      await loadDocument()
      await loadPDFDocument()
      if (token !== initToken) return
      if (pageCount.value === 0) {
        readerError.value = 'This document has no pages.'
        return
      }
      await fetchThumbnailBatch(0)
      ensureThumbnailBatchesForViewport()
      await renderPage(1)
      scrollThumbnailIntoView(1)
    } catch (err) {
      if (token !== initToken) return
      console.error('Failed to initialize PDF reader', err)
      readerError.value = 'Failed to load this PDF.'
    }
  }

  onMounted(async () => {
    void initializeReader()
    void connectCollab()

    keydownHandler = (e: KeyboardEvent) => {
      if (e.defaultPrevented || e.metaKey || e.ctrlKey || e.altKey) return

      const target = e.target as HTMLElement | null
      if (target?.closest('input, textarea, select, [contenteditable="true"]')) return

      if (e.key === 'ArrowRight' || e.key === 'ArrowDown' || e.key === 'PageDown') {
        e.preventDefault()
        nextPage()
      }
      if (e.key === 'ArrowLeft' || e.key === 'ArrowUp' || e.key === 'PageUp') {
        e.preventDefault()
        prevPage()
      }
    }
    window.addEventListener('keydown', keydownHandler)
  })

  watch(
    () => route.params.id,
    async (next, prev) => {
      if (String(next ?? '') === String(prev ?? '')) return
      void initializeReader()
      void connectCollab()
    },
  )

  onBeforeUnmount(() => {
    initToken++
    closeCollabConnection()
    if (keydownHandler) {
      window.removeEventListener('keydown', keydownHandler)
      keydownHandler = null
    }
    disposeCurrentPdf()
  })

  return {
    pageCount,
    currentPage,
    canvasEl,
    thumbnailScrollEl,
    thumbnails,
    readerError,
    collabStatus,
    collabError,
    pageRenderVersion,
    subscribeCollabMessages,
    requestPageAnnotations,
    createAnnotation,
    updateAnnotation,
    moveAnnotation,
    deleteAnnotation,
    onThumbnailScroll,
    go,
    goFirst,
    goLast,
    prevPage,
    nextPage,
  }
}
