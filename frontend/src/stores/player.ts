import { defineStore } from 'pinia'

export type PlayerTrack = {
  id: string
  title: string
  src: string
  poster?: string
}

export type PlayerContext = {
  key: string
  title?: string
}

export const playerStore = defineStore('player', {
  state: () => {
    return {
      queue: [] as PlayerTrack[],
      currentIndex: -1,
      collapsed: false,
      playlistVisible: false,
      context: null as PlayerContext | null,
      pendingAutoNext: false,
      barHeight: 0,
    }
  },
  getters: {
    hasTrack: (state) => state.currentIndex >= 0 && state.currentIndex < state.queue.length,
    currentTrack: (state) => {
      if (state.currentIndex < 0 || state.currentIndex >= state.queue.length) return null
      return state.queue[state.currentIndex]
    },
    isLastTrack: (state) => state.queue.length > 0 && state.currentIndex === state.queue.length - 1,
  },
  actions: {
    setContext(context: PlayerContext | null) {
      this.context = context
    },
    setQueue(queue: PlayerTrack[], startIndex = 0) {
      this.queue = queue
      this.currentIndex = queue.length > 0 ? Math.min(Math.max(startIndex, 0), queue.length - 1) : -1
      this.collapsed = false
      this.pendingAutoNext = false
    },
    appendQueue(tracks: PlayerTrack[]) {
      if (!tracks || tracks.length === 0) return
      const existing = new Set(this.queue.map((t) => t.id))
      const appended = tracks.filter((t) => t && !existing.has(t.id))
      if (appended.length === 0) return
      this.queue = this.queue.concat(appended)
    },
    updateTrackSource(trackId: string, src: string, poster?: string) {
      if (!trackId) return
      const idx = this.queue.findIndex((t) => t.id === trackId)
      if (idx < 0) return
      const t = this.queue[idx]
      t.src = src
      if (poster) t.poster = poster
    },
    setCurrentIndex(index: number) {
      if (index < 0 || index >= this.queue.length) return
      this.currentIndex = index
      this.pendingAutoNext = false
    },
    next() {
      if (this.queue.length === 0) return
      const nextIndex = this.currentIndex + 1
      if (nextIndex >= this.queue.length) return
      this.currentIndex = nextIndex
      this.pendingAutoNext = false
    },
    prev() {
      if (this.queue.length === 0) return
      const prevIndex = this.currentIndex - 1
      if (prevIndex < 0) return
      this.currentIndex = prevIndex
      this.pendingAutoNext = false
    },
    requestAutoNext() {
      this.pendingAutoNext = true
    },
    consumeAutoNext() {
      if (!this.pendingAutoNext) return false
      this.pendingAutoNext = false
      this.next()
      return true
    },
    toggleCollapsed() {
      this.collapsed = !this.collapsed
    },
    openPlaylist() {
      this.playlistVisible = true
    },
    closePlaylist() {
      this.playlistVisible = false
    },
    setBarHeight(height: number) {
      const h = Number(height) || 0
      this.barHeight = h > 0 ? Math.floor(h) : 0
    },
    closePlayer() {
      this.queue = []
      this.currentIndex = -1
      this.collapsed = false
      this.playlistVisible = false
      this.context = null
      this.pendingAutoNext = false
      this.barHeight = 0
    },
  },
})
