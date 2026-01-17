<script lang="ts" setup>
import { onMounted, onUnmounted, computed } from 'vue'
import Menu from './components/Menu.vue'
import GlobalAudioPlayer from './components/GlobalAudioPlayer.vue'
import { ElConfigProvider } from 'element-plus'
import 'element-plus/es/components/message/style/css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import { themeStore } from './stores/theme'
import { playerStore } from './stores/player'
import { AudioDetailAlias } from '../wailsjs/go/backend/App'

// 初始化主题
const store = themeStore()
onMounted(() => {
  store.initTheme()
})

const pStore = playerStore()
const odobCache = new Map<string, { src: string; poster?: string }>()
const odobPending = new Map<string, Promise<{ src: string; poster?: string }>>()

const globalPlayerHeight = computed(() => {
  if (!pStore.hasTrack) return 0
  const h = Number(pStore.barHeight) || 0
  if (h > 0) return h
  return pStore.collapsed ? 76 : 120
})

const mainStyle = computed(() => {
  return { '--global-player-height': `${globalPlayerHeight.value}px` } as any
})

const resolveOdobSrc = async (aliasId: string) => {
  const key = String(aliasId || '').trim()
  if (!key) return { src: '' }
  const cached = odobCache.get(key)
  if (cached) return cached
  const pending = odobPending.get(key)
  if (pending) return pending
  const p = AudioDetailAlias(key)
    .then((detail) => {
      const src = String(detail?.mp3_play_url ?? '').trim()
      const poster = String(detail?.icon ?? '').trim() || undefined
      const val = { src, poster }
      if (src) odobCache.set(key, val)
      return val
    })
    .finally(() => {
      odobPending.delete(key)
    })
  odobPending.set(key, p)
  return p
}

const onResolveTrack = async (ev: any) => {
  const detail = ev?.detail || {}
  const contextKey = String(detail?.contextKey ?? '')
  if (contextKey !== 'odob:study') return
  const trackId = String(detail?.trackId ?? '')
  if (!trackId) return
  const aliasId = trackId.startsWith('odob:') ? trackId.slice(5) : trackId
  if (!aliasId) return
  try {
    const { src, poster } = await resolveOdobSrc(aliasId)
    if (!src) return
    pStore.updateTrackSource(trackId, src, poster)
  } catch {
  }
}

onMounted(() => {
  window.addEventListener('player:resolveTrack', onResolveTrack as any)
})

onUnmounted(() => {
  window.removeEventListener('player:resolveTrack', onResolveTrack as any)
})

// Element Plus 主题配置
const elementTheme = computed(() => store.isDark ? 'dark' : 'light')
</script>

<template>
  <el-config-provider :locale="zhCn" :theme="elementTheme">
    <el-container>
      <el-header>
        <Menu />
      </el-header>
      <el-main :style="mainStyle">
        <router-view></router-view>
      </el-main>
      <GlobalAudioPlayer />
      <!-- <el-footer>Footer</el-footer> -->
    </el-container>
  </el-config-provider>
</template>

<style lang="scss">
@import url("./assets/css/font.css");

body {

  // background-color: transparent;
  img {
    max-width: 100%;
  }
}


// #app {
//   // position: relative;
//   // width: 900px;
//   // height: 520px;
// }

.el-container {
  height: 100%;
  background-color: var(--bg-color);
  transition: background-color 0.3s ease;
}

.el-header {
  position: relative;
  flex-direction: row;
  flex-wrap: nowrap;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 0;
  background-color: var(--bg-color);
  border-bottom: 1px solid var(--border-color-lighter);
  transition: background-color 0.3s ease, border-color 0.3s ease;

  .el-menu {
    height: auto;
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;
    justify-content: flex-start;
    background-color: transparent;
    border-bottom-width: 0px;
    flex: 1;

    .el-menu-item:hover,
    .el-menu-item:active,
    .el-menu-item:focus{
      background:none;
    }

    ul {
      a,
      a:hover,
      a:active,
      a:visited,
      a:link,
      a:focus {
        display: inline-block;
        padding: 0 5px;
        margin-right: 8px;
        text-align: center;
        text-decoration: none;
        white-space: nowrap;
        text-decoration: none;
        color: var(--text-color);
        transition: color 0.3s ease;
      }
    }
  }

}

.el-main {
  overflow-y: scroll;
  color: var(--text-color-secondary);
  width: 100%;
  height: 100%;
  padding-bottom: calc(var(--global-player-height, 0px) + 10px);
  background-color: var(--bg-color);
  transition: background-color 0.3s ease, color 0.3s ease;

  .el-pagination {
    margin-top: 10px;
    margin-bottom: 10px;
  }

  // 『只要在el-table元素中定义了height属性，即可实现固定表头的表格』不生效解决办法。
  .el-table {
    .el-table__body-wrapper {
      height: calc(100% - 5px) !important; // 表格高度减去表头的高度
    }
  }

  .el-breadcrumb {
    margin-bottom: 15px;
  }


}
</style>
