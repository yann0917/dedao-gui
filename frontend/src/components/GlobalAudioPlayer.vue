<template>
  <div v-if="store.hasTrack" ref="rootRef" class="global-audio-player" :class="store.collapsed ? 'collapsed' : ''">
    <div class="player-bar" @dblclick="store.toggleCollapsed">
      <div class="cover" @click="store.openPlaylist">
        <el-image
          v-if="store.currentTrack?.poster"
          :src="store.currentTrack.poster"
          fit="cover"
          class="cover-img"
          :preview-teleported="true"
          :preview-src-list="[store.currentTrack.poster]"
        />
        <div v-else class="cover-fallback">
          <el-icon><Headset /></el-icon>
        </div>
      </div>

      <div class="main">
        <div class="top-row">
          <div class="title" :title="store.currentTrack?.title ?? ''">{{ store.currentTrack?.title }}</div>
          <div class="meta">
            <span class="idx">{{ store.currentIndex + 1 }}/{{ store.queue.length }}</span>
            <el-dropdown @command="setRate">
              <span class="rate">{{ playbackRate }}x</span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-for="r in rates" :key="r" :command="r">{{ r }}x</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button text size="small" class="collapse-btn" @click="store.toggleCollapsed">
              <el-icon v-if="store.collapsed"><ArrowUp /></el-icon>
              <el-icon v-else><ArrowDown /></el-icon>
            </el-button>
            <el-button text size="small" class="close-btn" @click="handleClose">
              <el-icon><Close /></el-icon>
            </el-button>
          </div>
        </div>

        <div class="progress-row">
          <span class="time">{{ formatTime(displayTime) }}</span>
          <el-slider
            class="progress"
            :min="0"
            :max="Math.max(duration, 0)"
            :step="0.25"
            :show-tooltip="false"
            :model-value="displayTime"
            @update:modelValue="handleSeekInput"
            @change="handleSeekChange"
          />
          <span class="time">{{ formatTime(duration) }}</span>
        </div>

        <div class="controls-row">
          <div class="left-controls">
            <el-button text size="small" @click="store.prev" :disabled="store.currentIndex <= 0">
              <el-icon><DArrowLeft /></el-icon>
            </el-button>
            <el-button text size="small" @click="togglePlay">
              <el-icon v-if="isPlaying"><VideoPause /></el-icon>
              <el-icon v-else><VideoPlay /></el-icon>
            </el-button>
            <el-button text size="small" @click="store.next" :disabled="store.currentIndex >= store.queue.length - 1">
              <el-icon><DArrowRight /></el-icon>
            </el-button>
          </div>

          <div class="right-controls">
            <el-button text size="small" @click="store.openPlaylist">
              <el-icon><List /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <audio id="global-audio" class="video-js vjs-default-skin vjs-big-play-centered"></audio>

    <el-drawer
      v-model="store.playlistVisible"
      direction="btt"
      size="45%"
      :title="store.context?.title ? `播放列表 - ${store.context.title}` : '播放列表'"
      :append-to-body="true"
    >
      <div class="playlist">
        <div
          v-for="(item, idx) in store.queue"
          :key="item.id"
          class="playlist-item"
          :class="idx === store.currentIndex ? 'active' : ''"
          @click="handlePick(idx)"
        >
          <div class="playlist-title">{{ item.title }}</div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'
import { playerStore } from '../stores/player'
import { ArrowDown, ArrowUp, Close, DArrowLeft, DArrowRight, Headset, List, VideoPause, VideoPlay } from '@element-plus/icons-vue'

const store = playerStore()

let player: any = null
let resizeObserver: ResizeObserver | null = null
const rootRef = ref<HTMLElement | null>(null)
const duration = ref(0)
const currentTime = ref(0)
const seeking = ref(false)
const seekTime = ref(0)
const playbackRate = ref(1)
const rates = [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2]
const playing = ref(false)

const formatTime = (sec: number) => {
  const s = Math.max(0, Math.floor(sec || 0))
  const h = Math.floor(s / 3600)
  const m = Math.floor((s % 3600) / 60)
  const r = s % 60
  const pad = (n: number) => String(n).padStart(2, '0')
  if (h > 0) return `${h}:${pad(m)}:${pad(r)}`
  return `${m}:${pad(r)}`
}

const displayTime = computed(() => (seeking.value ? seekTime.value : currentTime.value))

const isPlaying = computed(() => playing.value)

const inferSourceType = (src: string) => {
  const raw = String(src || '').trim()
  const s = raw.toLowerCase()
  if (!s) return 'audio/mpeg'
  if (s.includes('.m3u8')) return 'application/x-mpegURL'
  const noHash = s.split('#')[0]
  const noQuery = noHash.split('?')[0]
  if (noQuery.endsWith('.mp3')) return 'audio/mpeg'
  if (noQuery.endsWith('.m4a')) return 'audio/mp4'
  if (noQuery.endsWith('.aac')) return 'audio/aac'
  return 'audio/mpeg'
}

const ensurePlayer = async () => {
  if (player) return
  await nextTick()
  player = videojs('global-audio', {
    language: 'zh-Hans',
    controls: false,
    autoplay: true,
    preload: 'auto',
    html5: {
      vhs: {
        withCredentials: false,
      },
    },
    controlBar: {
      remainingTimeDisplay: { displayNegative: false },
      fullscreenToggle: false,
      pictureInPictureToggle: false,
    },
    playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2],
    audioOnlyMode: true,
    audioPosterMode: true,
    volume: 0.6,
  })

  const syncPlaying = () => {
    try {
      playing.value = !!player && !player.paused()
    } catch {
      playing.value = false
    }
  }

  player.on('play', syncPlaying)
  player.on('pause', syncPlaying)
  player.on('ended', syncPlaying)
  player.on('error', syncPlaying)
  syncPlaying()

  player.on('timeupdate', () => {
    if (seeking.value) return
    try {
      currentTime.value = player.currentTime() || 0
    } catch {
    }
  })

  const refreshDuration = () => {
    try {
      duration.value = player.duration() || 0
    } catch {
      duration.value = 0
    }
  }

  player.on('loadedmetadata', refreshDuration)
  player.on('durationchange', refreshDuration)
  player.on('ratechange', () => {
    try {
      playbackRate.value = player.playbackRate() || 1
    } catch {
      playbackRate.value = 1
    }
  })

  player.on('ended', async () => {
    if (store.currentIndex < store.queue.length - 1) {
      store.next()
      return
    }
    if (store.context?.key) {
      store.requestAutoNext()
      window.dispatchEvent(new CustomEvent('player:needMore', { detail: { contextKey: store.context.key } }))
    }
  })
}

const handleClose = () => {
  store.closePlayer()
}

const loadCurrent = async () => {
  const track = store.currentTrack
  if (!track) return
  await ensurePlayer()
  const src = String(track.src || '').trim()
  if (!src) {
    if (store.context?.key) {
      window.dispatchEvent(
        new CustomEvent('player:resolveTrack', {
          detail: { contextKey: store.context.key, trackId: track.id },
        }),
      )
    }
    return
  }
  player.poster(track.poster || '')
  player.src({ src, type: inferSourceType(src) })
  try {
    await player.play()
  } catch {
  }
}

const setRate = async (r: number) => {
  await ensurePlayer()
  const rate = Number(r)
  if (!Number.isFinite(rate)) return
  playbackRate.value = rate
  try {
    player.playbackRate(rate)
  } catch {
  }
}

const togglePlay = async () => {
  await ensurePlayer()
  if (player.paused()) {
    try {
      await player.play()
    } catch {
    }
    return
  }
  player.pause()
}

const handlePick = (index: number) => {
  store.setCurrentIndex(index)
}

const handleSeekInput = (v: number) => {
  seeking.value = true
  seekTime.value = Number(v) || 0
}

const handleSeekChange = async (v: number) => {
  await ensurePlayer()
  const t = Number(v) || 0
  try {
    player.currentTime(t)
  } catch {
  }
  seeking.value = false
}

const maybePrefetchMore = () => {
  if (!store.context?.key) return
  if (store.queue.length === 0) return
  if (store.currentIndex < store.queue.length - 2) return
  window.dispatchEvent(new CustomEvent('player:needMore', { detail: { contextKey: store.context.key } }))
}

watch(
  () => store.currentIndex,
  () => {
    seeking.value = false
    seekTime.value = 0
    currentTime.value = 0
    duration.value = 0
    loadCurrent()
    maybePrefetchMore()
  },
  { immediate: true },
)

watch(
  () => `${store.currentTrack?.id ?? ''}|${store.currentTrack?.src ?? ''}`,
  () => {
    if (!store.currentTrack) return
    loadCurrent()
  },
)

const refreshBarHeight = () => {
  const el = rootRef.value
  if (!el) return
  store.setBarHeight(el.getBoundingClientRect().height)
}

onMounted(() => {
  if (store.hasTrack) ensurePlayer()
})

watch(
  () => store.hasTrack,
  async (hasTrack) => {
    if (!hasTrack) {
      seeking.value = false
      seekTime.value = 0
      currentTime.value = 0
      duration.value = 0
      playing.value = false
      store.setBarHeight(0)

      if (resizeObserver) {
        try {
          resizeObserver.disconnect()
        } catch {
        }
        resizeObserver = null
      }

      if (player) {
        try {
          player.pause()
        } catch {
        }
        try {
          player.dispose()
        } catch {
        }
        player = null
      }
      return
    }

    await nextTick()
    refreshBarHeight()
    if (!resizeObserver && rootRef.value) {
      resizeObserver = new ResizeObserver(() => {
        refreshBarHeight()
      })
      resizeObserver.observe(rootRef.value)
    }
  },
  { immediate: true },
)

watch(
  () => store.collapsed,
  async () => {
    if (!store.hasTrack) return
    await nextTick()
    refreshBarHeight()
  },
)

onUnmounted(() => {
  if (resizeObserver) {
    try {
      resizeObserver.disconnect()
    } catch {
    }
    resizeObserver = null
  }
  if (player) {
    try {
      player.dispose()
    } catch {
    }
    player = null
  }
})
</script>

<style scoped>
.global-audio-player {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background: var(--bg-color);
  border-top: 1px solid var(--border-color-lighter);
  padding: 10px 10px calc(10px + env(safe-area-inset-bottom));
}

#global-audio {
  position: absolute;
  left: -9999px;
  width: 1px;
  height: 1px;
  opacity: 0;
}

.player-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cover {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  overflow: hidden;
  flex: 0 0 auto;
  background: var(--border-color-lighter);
  cursor: pointer;
}

.cover-img {
  width: 56px;
  height: 56px;
}

.cover-fallback {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-color-secondary);
}

.main {
  flex: 1;
  min-width: 0;
}

.top-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title {
  flex: 1;
  min-width: 0;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 0 0 auto;
}

.collapse-btn {
  margin-left: -6px;
}

.close-btn {
  margin-left: -8px;
}

.idx {
  font-size: 12px;
  color: var(--text-color-secondary);
}

.rate {
  font-size: 12px;
  color: var(--text-color-secondary);
  cursor: pointer;
  user-select: none;
}

.progress-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 6px;
}

.time {
  width: 44px;
  font-size: 12px;
  color: var(--text-color-secondary);
  text-align: center;
  flex: 0 0 auto;
}

.progress {
  flex: 1;
  min-width: 0;
}

.controls-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
}

.left-controls,
.right-controls {
  display: flex;
  align-items: center;
  gap: 2px;
}

.collapsed .progress-row,
.collapsed .controls-row {
  display: none;
}

.playlist {
  padding: 8px 0;
}

.playlist-item {
  padding: 10px 12px;
  border-radius: 10px;
  cursor: pointer;
}

.playlist-item:hover {
  background: var(--border-color-lighter);
}

.playlist-item.active {
  background: rgba(255, 107, 0, 0.12);
}

.playlist-title {
  color: var(--text-color);
  font-size: 14px;
  line-height: 20px;
}
</style>
