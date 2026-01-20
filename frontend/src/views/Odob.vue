<template>
    <div class="odob-container">
        <div v-if="groupMode.active" class="group-header">
            <el-button type="primary" link @click="exitGroup">
                <el-icon class="el-icon--left"><ArrowLeft /></el-icon>返回
            </el-button>
            <span class="group-title">{{ groupMode.title }}</span>
        </div>

        <div v-loading="initLoading" class="odob-grid-container" v-infinite-scroll="loadMore" :infinite-scroll-disabled="disabled" :infinite-scroll-immediate="false">
            <div v-if="tableData.list && tableData.list.length > 0" class="odob-grid">
                <div v-for="item in tableData.list" :key="item.id" class="odob-card" @click="item.is_group ? enterGroup(item) : null">
                    <div class="card-cover">
                        <el-image 
                            v-if="item.icon" 
                            :src="item.icon" 
                            fit="cover"
                            loading="lazy"
                        >
                            <template #placeholder>
                                <div class="image-placeholder">
                                    <el-icon><Picture /></el-icon>
                                </div>
                            </template>
                        </el-image>
                        <div v-else class="no-cover">
                            <el-icon v-if="item.is_group" :size="40"><Folder /></el-icon>
                            <span v-else>无封面</span>
                        </div>
                        
                        <!-- 悬停遮罩层 (仅非分组显示操作) -->
                        <div v-if="!item.is_group" class="card-overlay" @click.stop>
                            <div class="overlay-actions">
                                <el-tooltip content="播放" :show-after="500">
                                    <el-button circle type="primary" :icon="VideoPlay" @click="handlePlay(item)" />
                                </el-tooltip>
                                <el-tooltip content="文稿" :show-after="500">
                                    <el-button circle type="success" :icon="Memo" @click="gotoArticleDetail(item)" />
                                </el-tooltip>
                                <el-tooltip content="详情" :show-after="500">
                                    <el-button circle type="info" :icon="View" @click="handleProd(item)" />
                                </el-tooltip>
                                <el-tooltip content="下载" :show-after="500">
                                    <el-button circle type="warning" :icon="DownloadIcon" @click="openDownloadDialog(item)" />
                                </el-tooltip>
                            </div>
                        </div>
                        
                        <!-- 标签标识 -->
                        <div class="card-badges">
                            <el-tag v-if="item.type === 1013" type="warning" size="small" effect="dark">名家讲书</el-tag>
                            <el-tag v-if="item.is_group" type="info" size="small" effect="dark">分组</el-tag>
                        </div>
                    </div>
                    
                    <div class="card-content">
                        <h3 class="card-title" :title="item.title">{{ item.title }}</h3>
                        <div class="card-meta">
                            <span class="meta-info">
                                <el-icon><Clock /></el-icon>
                                {{ secondToHour(item.duration) }}
                            </span>
                            <span v-if="item.is_group" class="meta-info">共 {{ item.course_num || 0 }} 本</span>
                        </div>
                        <div class="card-intro" v-if="item.intro">
                            {{ item.intro.length > 40 ? item.intro.substring(0, 40) + '...' : item.intro }}
                        </div>
                    </div>
                </div>
            </div>
            <el-empty v-else description="暂无内容" />
        </div>
    
    
    </div>

    <audio-info v-if="dialogVisible" :enid="prodEnid" :dialog-visible="dialogVisible" @close="closeDialog"></audio-info>
    <outside-info v-if="outsideVisible" :enid="prodEnid" :dialog-visible="outsideVisible" @close="closeDialog"></outside-info>

    <el-dialog v-model="dialogDownloadVisible" title="请选择下载格式" align-center center width="400px" class="download-dialog">
        <el-form label-position="top">
            <el-form-item label="下载格式">
                <el-select v-model="downloadType" placeholder="请选择下载格式" style="width: 100%">
                    <el-option v-for="item in downloadTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
            <div v-if="percentage > 0" class="download-progress">
                <div class="progress-text">{{ content }}</div>
                <el-progress 
                    :percentage="percentage" 
                    :stroke-width="10" 
                    status="success"
                    striped
                    striped-flow
                />
            </div>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDownloadDialog">取消</el-button>
                <el-button type="primary" @click="download(downloadId, downloadType)" :loading="percentage > 0 && percentage < 100">
                    {{ percentage > 0 && percentage < 100 ? '下载中' : '开始下载' }}
                </el-button>
            </span>
        </template>
    </el-dialog>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, computed } from 'vue'
import 'element-plus/es/components/message/style/css'
import { ElMessage } from 'element-plus'
import { ArrowLeft, VideoPlay, Memo, View, Download as DownloadIcon, Picture, Folder, Clock } from '@element-plus/icons-vue'
import { AudioDetailAlias, CourseCategory, CourseGroupList, CourseList, OdobDownload, SetDir } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import AudioInfo from '../components/AudioInfo.vue'
import OutsideInfo from '../components/OutsideInfo.vue'
import { userStore } from '../stores/user'
import { settingStore } from '../stores/setting'
import { Local } from '../utils/storage'
import { useAppRouter } from '../composables/useRouter'

import { secondToHour } from '../utils/utils'
import { EventsOff, EventsOn } from '../../wailsjs/runtime/runtime'
import { playerStore, type PlayerTrack } from '../stores/player'

const pStore = playerStore()

const store = userStore()
const setStore = settingStore()
const { pushOdobDetail, pushLogin, pushSetting } = useAppRouter()
const loading = ref(false)
const initLoading = ref(true)
const page = ref(1)
const total = ref(0)
const outerTotal = ref(0)
const pageSize = ref(20) // Increased for scrolling
const lastPageSize = ref(20)
const dialogVisible = ref(false)
const outsideVisible = ref(false)
const prodEnid = ref("")

const groupMode = reactive({
    active: false,
    groupId: 0,
    title: '',
})

const dialogDownloadVisible = ref(false)
const downloadType = ref(1)
const downloadId = ref(0)
let downloadData = reactive(new services.Course)
const downloadTypeOptions = [
    {value: 1, label: "MP3"}, {value: 2, label: "PDF"}, {value: 3, label: "Markdown"}
]

let percentage=ref(0)
let content=ref('')

const aliasSrcCache = new Map<string, { src: string; poster?: string }>()
const aliasPending = new Map<string, Promise<{ src: string; poster?: string }>>()

const resolveOdobSrc = async (aliasId: string) => {
    const key = String(aliasId || '').trim()
    if (!key) return { src: '' }
    const cached = aliasSrcCache.get(key)
    if (cached) return cached
    const pending = aliasPending.get(key)
    if (pending) return pending
    const p = AudioDetailAlias(key)
        .then((detail) => {
            const src = String(detail?.mp3_play_url ?? '').trim()
            const poster = String(detail?.icon ?? '').trim() || undefined
            const val = { src, poster }
            if (src) aliasSrcCache.set(key, val)
            return val
        })
        .finally(() => {
            aliasPending.delete(key)
        })
    aliasPending.set(key, p)
    return p
}

const buildTrack = (row: any): PlayerTrack | null => {
    const aliasId = String(row?.audio_detail?.alias_id ?? '').trim()
    if (!aliasId) return null
    const cached = aliasSrcCache.get(aliasId)
    return {
        id: `odob:${aliasId}`,
        title: String(row?.title ?? ''),
        src: cached?.src ?? '',
        poster: cached?.poster || row?.icon || row?.audio_detail?.icon,
    }
}

const handlePlay = async (row: any) => {
    pStore.setContext({ key: 'odob:study', title: '每日听书' })
    const queue = (tableData.list || []).map(buildTrack).filter((t): t is PlayerTrack => !!t)
    const target = buildTrack(row)
    if (!target) {
        ElMessage({ message: '该条目没有可播放的音频地址', type: 'warning' })
        return
    }
    const startIndex = queue.findIndex((t) => t.id === target.id)
    const idx = startIndex >= 0 ? startIndex : 0
    const aliasId = String(row?.audio_detail?.alias_id ?? '').trim()
    if (aliasId) {
        try {
            const { src, poster } = await resolveOdobSrc(aliasId)
            if (src) {
                queue[idx].src = src
                if (poster && !queue[idx].poster) queue[idx].poster = poster
            }
        } catch (e) {
            ElMessage({ message: '获取音频播放地址失败', type: 'warning' })
        }
    }
    pStore.setQueue(queue, idx)
}

let tableData = reactive(new services.CourseList)
let courseInfo = reactive(new services.CourseInfo)

onMounted(() => {
    CourseCategory().then(result => {
        result.forEach((item, key) => {
            if (item.category == "odob") {
                outerTotal.value = item.count
                if (!groupMode.active) total.value = item.count
            }
        })

    }).catch((error) => {
        if (error == '401 Unauthorized') {
            store.user = null
            pushLogin()
        }
        Local.remove("cookies")
        Local.remove("userStore")
    })

})

const noMore = computed(() => {
    const currentCount = tableData.list ? tableData.list.length : 0
    if (groupMode.active) {
        return currentCount >= (tableData.total || 0)
    }
    if (outerTotal.value > 0) {
        return currentCount >= outerTotal.value
    }
    return lastPageSize.value < pageSize.value
})

const disabled = computed(() => loading.value || noMore.value)

const loadMore = () => {
    if (disabled.value) return
    page.value += 1
    getTableData(true)
}

const getTableData = async (append = false) => {
    loading.value = true
    if (!append) initLoading.value = true
    const fetcher = groupMode.active
        ? CourseGroupList("odob", "study", groupMode.groupId, page.value, pageSize.value)
        : CourseList("odob", "study", page.value, pageSize.value)
    await fetcher.then((table) => {
        loading.value = false
        initLoading.value = false
        
        const fetchedList = table.list || []
        lastPageSize.value = fetchedList.length
        
        if (append) {
            if (fetchedList.length > 0) {
                tableData.list.push(...fetchedList)
            }
        } else {
            Object.assign(tableData, table)
        }
        total.value = groupMode.active ? (table.total || 0) : outerTotal.value
    }).catch((error) => {
        loading.value = false
        initLoading.value = false
        ElMessage({
            message: error,
            type: 'warning'
        })
    })
}

const dialogTitle = ref('detail')
const openDialog = () => {
    dialogVisible.value = true
}
const closeDialog = () => {
    dialogVisible.value = false
    outsideVisible.value = false
}
getTableData()

const enterGroup = (row: any) => {
    const groupId = Number(row?.group_id || 0)
    if (!groupId) return
    groupMode.active = true
    groupMode.groupId = groupId
    groupMode.title = String(row?.title || '')
    page.value = 1
    getTableData()
}

const exitGroup = () => {
    groupMode.active = false
    groupMode.groupId = 0
    groupMode.title = ''
    page.value = 1
    total.value = outerTotal.value
    getTableData()
}


const openDownloadDialog = (row: any) => {
    downloadId.value = row.id
    dialogDownloadVisible.value = true

    Object.assign(downloadData, row)
    if (setStore.getDownloadDir == "") {
        ElMessage({
            message: '请设置文件保存目录',
            type: 'warning'
        })
        pushSetting()
    } else {
        SetDir([setStore.getDownloadDir,
            setStore.getFfmpegDirDir,
            setStore.getWkDir]).then(() => {
        }).catch((error) => {
            ElMessage({
                message: error,
                type: 'warning'
            })
        })
    }
}
const closeDownloadDialog = () => {
    downloadType.value = 1
    dialogDownloadVisible.value = false
    percentage.value = 0
    content.value = ''
    EventsOff("odobDownload")
}

const download = async (id: number, dType: number) => {
    content.value = '下载中...'
    EventsOn("odobDownload",  data=>{
        if (data) {
            console.log(data)
            percentage.value = data.pct
            content.value = data.value + '下载中...'
        }
    })

    await OdobDownload(id, dType, downloadData).then((info) => {
        console.log(info)
    }).catch((error) => {
        ElMessage({
            message: error,
            type: 'warning'
        })
    })
    closeDownloadDialog()
    return
}


const gotoArticleDetail = (row: any) => {
    const id = row.audio_detail.alias_id
    pushOdobDetail(id)
}

const handleProd = (row: any) => {
    prodEnid.value = row.enid
    if (row.type === 1013) {
        // 名家讲书类型显示 OutsideInfo 组件
        outsideVisible.value = true
    } else {
        // 普通音频类型显示 AudioInfo 组件
        dialogVisible.value = true
    }
}

</script>

<style scoped>
.odob-container {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 20px;
}

.group-header {
    display: flex;
    align-items: center;
    margin-bottom: 24px;
    gap: 16px;
}

.group-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
}

.odob-grid-container {
    flex: 1;
    overflow-y: auto;
    min-height: 0; /* Important for flex child scrolling */
    padding-bottom: 20px;
    
    /* 隐藏滚动条但保留功能 - 清新风格 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE 10+ */
}

.odob-grid-container::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}

.odob-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 20px;
    padding: 4px;
}

.odob-card {
    background: var(--card-bg);
    border-radius: 16px;
    overflow: hidden;
    box-shadow: var(--shadow-soft);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    cursor: default;
    border: 1px solid var(--border-soft);
    display: flex;
    flex-direction: column;
}

.odob-card:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-medium);
    border-color: var(--accent-color);
}

.card-cover {
    position: relative;
    width: 100%;
    padding-top: 100%; /* 1:1 Aspect Ratio */
    background-color: var(--fill-color);
    overflow: hidden;
}

.card-cover .el-image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    transition: transform 0.5s ease;
}

.odob-card:hover .card-cover .el-image {
    transform: scale(1.05);
}

.image-placeholder, .no-cover {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--text-secondary);
    background-color: var(--fill-color);
}

.card-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;
    backdrop-filter: blur(4px);
}

.odob-card:hover .card-overlay {
    opacity: 1;
}

.overlay-actions {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    justify-content: center;
    padding: 10px;
}

.overlay-actions .el-button {
    margin: 0 !important;
    transform: scale(0.9);
    transition: transform 0.2s;
}

.overlay-actions .el-button:hover {
    transform: scale(1.1);
}

.card-badges {
    position: absolute;
    top: 10px;
    left: 10px;
    display: flex;
    flex-direction: column;
    gap: 6px;
    z-index: 2;
}

.card-content {
    padding: 16px;
    flex: 1;
    display: flex;
    flex-direction: column;
}

.card-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 8px 0;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    overflow: hidden;
    height: 44px; /* Fixed height for 2 lines */
}

.card-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 13px;
    color: var(--text-secondary);
    margin-bottom: 8px;
}

.meta-info {
    display: flex;
    align-items: center;
    gap: 4px;
}

.card-intro {
    font-size: 13px;
    color: var(--text-tertiary);
    line-height: 1.5;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    overflow: hidden;
    margin-top: auto; /* Push to bottom */
}

.download-progress {
    margin-top: 16px;
}

.progress-text {
    font-size: 13px;
    color: var(--text-secondary);
    margin-bottom: 6px;
    text-align: center;
}

/* Responsive adjustments removed to enforce 5 columns */
</style>
