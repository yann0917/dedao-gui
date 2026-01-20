<template>
  <div class="veplayer-page">
    <div class="veplayer-header">
      <el-button text @click="goBack">返回</el-button>
      <div class="veplayer-title" :title="title">{{ title }}</div>
      <el-button text :loading="loading" @click="reload">重试</el-button>
    </div>

    <div class="veplayer-body">
      <div ref="playerRoot" id="veplayer" class="veplayer-container"></div>

      <div v-if="loading" class="veplayer-loading">
        <el-skeleton :rows="3" animated style="width: 360px" />
      </div>

      <div v-if="errorText" class="veplayer-error">
        <el-result icon="error" title="播放失败" :sub-title="errorText">
          <template #extra>
            <el-button type="primary" @click="reload">重试</el-button>
          </template>
        </el-result>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, onUnmounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'
import { GetVolcPlayAuthToken } from '../../wailsjs/go/backend/App'
import { userStore } from '../stores/user'

type VePlayerCtor = new (options: Record<string, any>) => {
  dispose?: () => void
  on?: (event: string, cb: (...args: any[]) => void) => void
}

type MediaVolcLike = {
  tracks?: Array<{
    formats?: Array<{
      volc_play_auth_token?: string
    }>
  }>
}

const route = useRoute()
const router = useRouter()
const store = userStore()

const playerRoot = ref<HTMLDivElement | null>(null)
const playerSdk = ref<InstanceType<VePlayerCtor> | null>(null)
const loading = ref(false)
const errorText = ref('')

const title = computed(() => String(route.query.title ?? '视频'))

const getLineAppId = () => {
  const raw = (route.query.line_app_id ?? '') as any
  const n = Number(raw)
  if (Number.isFinite(n) && n > 0) return n
  return 233260
}

const VEPLAYER_SDK_URLS = [
  {
    css: 'https://lf-unpkg.volccdn.com/obj/vcloudfe/sdk/@volcengine/veplayer/1.15.1/index.min.css',
    js: 'https://lf-unpkg.volccdn.com/obj/vcloudfe/sdk/@volcengine/veplayer/1.15.1/index.min.js',
  },
  {
    css: 'https://sf-unpkg.bytepluscdn.com/obj/byteplusfe-sg/sdk/@volcengine/veplayer/1.15.1/index.min.css',
    js: 'https://sf-unpkg.bytepluscdn.com/obj/byteplusfe-sg/sdk/@volcengine/veplayer/1.15.1/index.min.js',
  },
]

const ensureCssLoaded = (href: string) => {
  const existing = document.querySelector(`link[data-veplayer-css="true"][href="${href}"]`)
  if (existing) return Promise.resolve()

  return new Promise<void>((resolve, reject) => {
    const link = document.createElement('link')
    link.rel = 'stylesheet'
    link.href = href
    link.dataset.veplayerCss = 'true'
    link.onload = () => resolve()
    link.onerror = () => reject(new Error(`VePlayer 样式加载失败：${href}`))
    document.head.appendChild(link)
  })
}

const ensureScriptLoaded = (src: string) => {
  const existing = document.querySelector(`script[data-veplayer-js="true"][src="${src}"]`)
  if (existing) return Promise.resolve()

  return new Promise<void>((resolve, reject) => {
    const script = document.createElement('script')
    script.src = src
    script.async = true
    script.dataset.veplayerJs = 'true'
    script.onload = () => resolve()
    script.onerror = () => reject(new Error(`VePlayer 脚本加载失败：${src}`))
    document.head.appendChild(script)
  })
}

const hasRequiredApi = (VePlayerLike: any) => {
  return typeof VePlayerLike === 'function'
}

const ensureVePlayer = async () => {
  const existing = (window as any).VePlayer
  if (hasRequiredApi(existing)) return existing as VePlayerCtor

  for (const url of VEPLAYER_SDK_URLS) {
    try {
      await ensureCssLoaded(url.css)
      await ensureScriptLoaded(url.js)
      const loaded = (window as any).VePlayer
      if (hasRequiredApi(loaded)) return loaded as VePlayerCtor
    } catch {
      continue
    }
  }

  throw new Error('VePlayer SDK 加载失败，请检查网络连接')
}

const pickPlayAuthToken = (volc: MediaVolcLike | null | undefined) => {
  const formats = volc?.tracks?.flatMap((t) => t.formats ?? []) ?? []
  const token = formats.map((f) => f?.volc_play_auth_token).find((t) => typeof t === 'string' && t.length > 0)
  return token ?? ''
}

const getErrorMessage = (err: unknown) => {
  if (err instanceof Error) return err.message
  if (typeof err === 'string') return err
  try {
    return JSON.stringify(err)
  } catch {
    return String(err)
  }
}

const destroyPlayer = () => {
  try {
    playerSdk.value?.dispose?.()
  } finally {
    playerSdk.value = null
    if (playerRoot.value) playerRoot.value.innerHTML = ''
  }
}

const fetchPlayAuthToken = async () => {
  const mediaId = String(route.query.media_id ?? '')
  const securityToken = String(route.query.security_token ?? '')

  if (!mediaId || !securityToken) {
    throw new Error('缺少 media_id 或 security_token')
  }

  const volc = (await GetVolcPlayAuthToken(mediaId, securityToken)) as unknown as MediaVolcLike
  const token = pickPlayAuthToken(volc)
  if (!token) {
    throw new Error('未获取到火山点播 playAuthToken')
  }
  return token
}

const createPlayer = async (playAuthToken: string) => {
  const VePlayer = await ensureVePlayer()

  await nextTick()
  if (!playerRoot.value) {
    throw new Error('播放器容器未就绪')
  }

  destroyPlayer()

  const instance = new (VePlayer as any)({
    root: playerRoot.value,
    getVideoByToken: {
      playAuthToken,
      definitionMap: {
        original: { definition: 'ori', definitionTextKey: 'ORI' },
        '360p': { definition: 'ld', definitionTextKey: 'LD' },
        '480p': { definition: 'sd', definitionTextKey: 'SD' },
        '720p': { definition: 'hd', definitionTextKey: 'HD' },
      },
    },
    vodLogOpts: {
      vtype: "FLV",
      tag: "直播",
      codec_type: "h264",
      drm_type: 1,
      line_app_id: getLineAppId(),
      line_user_id: String(store.user?.uid_hazy ?? 'unknown'),
      playerCoreVersion: "2.16.2"
    },
    languages: {
      zh: {
        ORI: '原画',
        LD: '流畅',
        SD: '标清',
        HD: '高清',
      },
    },
    lang: 'zh',
    autoplay: true,
    onTokenExpired: async () => {
      const newToken = await fetchPlayAuthToken()
      return { playAuthToken: newToken }
    },
  })

  instance?.on?.('error', (...args: any[]) => {
    const msg = args.map(getErrorMessage).join(' ')
    if (!msg) return
    errorText.value = msg
  })

  playerSdk.value = instance
}

const reload = async () => {
  loading.value = true
  errorText.value = ''
  try {
    const token = await fetchPlayAuthToken()
    await createPlayer(token)
  } catch (err) {
    const msg = getErrorMessage(err)
    errorText.value = msg
    ElMessage({ message: msg, type: 'warning' })
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

const routeKey = computed(() => `${String(route.query.media_id ?? '')}|${String(route.query.security_token ?? '')}`)
watch(routeKey, () => reload(), { immediate: true })

onUnmounted(() => {
  destroyPlayer()
})
</script>

<style scoped>
.veplayer-page {
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
  padding: 12px 0;
  box-sizing: border-box;
  overflow-y: auto;
  /* 隐藏滚动条但保留功能 - 清新风格 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

.veplayer-page::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

.veplayer-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 12px 12px;
}

.veplayer-title {
  flex: 1;
  text-align: left;
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--text-primary);
}

.veplayer-body {
  position: relative;
  flex: 1;
  min-height: 360px;
  padding: 0 12px 12px;
  box-sizing: border-box;
}

.veplayer-container {
  width: 100%;
  height: 100%;
  border-radius: 12px;
  overflow: hidden;
  background: #000;
}

.veplayer-loading {
  position: absolute;
  inset: 0 12px 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--bg-color);
  border-radius: 12px;
  z-index: 2;
}

.veplayer-error {
  position: absolute;
  inset: 0 12px 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--bg-color);
  border-radius: 12px;
  z-index: 3;
}
</style>
