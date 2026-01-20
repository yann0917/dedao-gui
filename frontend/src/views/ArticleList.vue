<template>
    <div class="breadcrumb-container">
        <el-breadcrumb :separator-icon="ArrowRight">
            <el-breadcrumb-item :to="{ name: ROUTE_NAMES.COURSE }">课程列表</el-breadcrumb-item>
            <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
        </el-breadcrumb>
        <div class="breadcrumb-actions">
            <el-button 
                type="primary" 
                link 
                @click="showCourseDetail"
                class="action-btn"
            >
                <el-icon><InfoFilled /></el-icon>
                课程详情
            </el-button>
            <el-button 
                type="primary" 
                link 
                @click="toggleSort"
                class="action-btn"
            >
                <el-icon><Sort /></el-icon>
                {{ isReverse ? '正序' : '倒序' }}
            </el-button>
        </div>
    </div>
    <div
        class="infinite-list-wrapper"
        v-infinite-scroll="loadMoreArticles"
        :infinite-scroll-disabled="loading || finished"
        infinite-scroll-distance="10"
        @scroll="onScroll"
    >
        <ul class="article-list">
            <li v-for="item in tableData.article_list" :key="item.id" class="article-card-wrapper">
                <div class="article-card" @click="handlePlay(item)">
                    <div class="card-cover">
                        <el-image 
                            :src="item.logo" 
                            class="card-image"
                            fit="cover"
                            loading="lazy"
                        >
                            <template #placeholder>
                                <div class="image-placeholder">
                                    <el-icon><Picture /></el-icon>
                                </div>
                            </template>
                        </el-image>
                        <div class="card-overlay" @click.stop>
                            <div class="overlay-actions">
                                <el-tooltip content="播放" :show-after="500">
                                    <el-button circle type="primary" :icon="VideoPlay" @click="handlePlay(item)" v-if="item.audio_alias_ids.length || item.video_status"/>
                                </el-tooltip>
                                <el-tooltip content="文稿" :show-after="500">
                                    <el-button circle type="success" :icon="Memo" @click="gotoArticleDetail(item)" />
                                </el-tooltip>
                                <el-tooltip content="下载" :show-after="500" v-if="canDownload(item)">
                                    <el-button circle type="warning" :icon="Download" @click="openDownloadDialog(item)" />
                                </el-tooltip>
                            </div>
                        </div>
                        <div class="card-badges">
                            <el-tag v-if="isLearned(item)" type="success" size="small" effect="dark">已学</el-tag>
                        </div>
                    </div>
                    
                    <div class="card-content">
                        <div class="card-header">
                            <h3 class="card-title" :title="item.title">{{ item.title }}</h3>
                        </div>
                        
                        <div class="card-meta">
                            <span class="meta-item">
                                <el-icon><Clock /></el-icon>
                                {{ item.video_status == 1 ? secondToHour(item.video?.[0]?.duration !== undefined ? item.video[0].duration : 0) : (item.audio?.duration ? secondToHour(item.audio.duration) : '') }}
                            </span>
                            <span class="meta-item">
                                <el-icon><User /></el-icon>
                                {{ item.cur_learn_count }}
                            </span>
                        </div>
                        
                        <div class="card-footer">
                            <span class="publish-time">
                                <el-icon><Calendar /></el-icon>
                                {{ item.publish_time ? timestampToDate(item.publish_time) : '' }}
                            </span>
                        </div>
                    </div>
                </div>
            </li>
        </ul>
        <div v-if="loading" class="loading-state">
            <el-icon class="is-loading"><Loading /></el-icon> 加载中...
        </div>
        <div v-if="finished && hasScrolled" class="finished-state">没有更多了</div>
    </div>
    <download-dialog
            v-if="dialogDownloadVisible"
            :dialog-visible="dialogDownloadVisible"
            :download-type-options="downloadTypeOptions"
            :prod-type="66"
            :download-id="downloadId"
            :article-id="downloadaId"
            :en-id="downloadEnId"
            @close="closeDownloadDialog">
    </download-dialog>

    <course-info
            v-if="courseDetailVisible"
            :dialog-visible="courseDetailVisible"
            :enid="enid"
            @close="closeCourseDetail">
    </course-info>

</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowRight, VideoPlay, Memo, Download, Picture, Loading } from '@element-plus/icons-vue'
import {ArticleList, SetDir} from '../../wailsjs/go/backend/App'
import {services} from '../../wailsjs/go/models'
import {useRoute} from 'vue-router'
import {secondToHour} from '../utils/utils'
import DownloadDialog from "../components/DownloadDialog.vue";
import CourseInfo from "../components/CourseInfo.vue";

import {settingStore} from "../stores/setting";
import { User, Clock, Calendar, Sort, InfoFilled } from '@element-plus/icons-vue'
import { timestampToDate } from '../utils/utils'
import {useAppRouter} from '../composables/useRouter'
import {ROUTE_NAMES} from '../router/routes'
import { playerStore, type PlayerTrack } from '../stores/player'

const pStore = playerStore()

const loading = ref(false) // 懒加载 loading 状态
const finished = ref(false) // 懒加载是否全部加载完
const page = ref(1) // 当前页码
const pageSize = ref(30) // 每页条数
const isReverse = ref(false) // 是否倒序
let total = ref(0)
const dialogVisible = ref(false)
const layout = ref('total, sizes, next')
const pageSizes = ref([10, 15, 20, 30, 50]);
const route = useRoute()
const setStore = settingStore()
const { router, pushArticleDetail, pushVideo } = useAppRouter()

const dialogDownloadVisible = ref(false)
const downloadType = ref(1)
const downloadId = ref(0)
const downloadaId = ref(0)
const downloadEnId = ref('')
const downloadTypeOptions = [
    {value: 1, label: "MP3"}, {value: 2, label: "PDF"}, {value: 3, label: "Markdown"}
]

const courseDetailVisible = ref(false)

const showCourseDetail = () => {
    courseDetailVisible.value = true
}

const closeCourseDetail = () => {
    courseDetailVisible.value = false
}

let id = ref()
let enid = ref()
let maxId = ref()
let tableData = reactive(new services.ArticleList)
let breadcrumbTitle = ref()
const hasScrolled = ref(false)

const onScroll = () => {
    if (!hasScrolled.value) hasScrolled.value = true
}

const toggleSort = () => {
    isReverse.value = !isReverse.value
    // 重置列表
    tableData.article_list = []
    page.value = 1
    finished.value = false
    maxId.value = 0
    loading.value = false // 确保不处于 loading 状态
    
    // 重新加载
    loadMoreArticles()
}

onMounted(() => {
    watch(() => {
            id.value = route.params.id
            enid.value = route.query.enid
            total.value = Number(route.query.total)
            breadcrumbTitle.value = route.query.title
        },
        () => {
            // 初始化时清空数据，重置分页
            tableData.article_list = []
            page.value = 1
            finished.value = false
            maxId.value = 0 // 初始化 maxId
            
            loadMoreArticles()
        },
        {immediate: true}
    )
    window.addEventListener('player:needMore', onNeedMoreForPlayer as any)
})

onUnmounted(() => {
    window.removeEventListener('player:needMore', onNeedMoreForPlayer as any)
})

const buildTrack = (row: any): PlayerTrack | null => {
    const src = String(row?.audio?.mp3_play_url ?? '').trim()
    if (!src) return null

    return {
        id: `${String(enid.value ?? '')}|${String(row?.enid ?? row?.id ?? '')}`,
        title: String(row?.title ?? ''),
        src,
        poster: row?.logo || row?.audio?.icon || row?.video?.[0]?.cover_img,
    }
}

const isLearned = (row: services.ArticleIntro) => {
    const isBuy = !!(row as any)?.is_buy
    const isUserFreeTry = !!(row as any)?.is_user_free_try
    const isRead = !!(row as any)?.is_read
    return (!isBuy && isUserFreeTry) || (isBuy && isRead)
}

const canDownload = (row: services.ArticleIntro) => {
    const isBuy = !!(row as any)?.is_buy
    const isUserFreeTry = !!(row as any)?.is_user_free_try
    return isBuy || (!isBuy && isUserFreeTry)
}

const handlePlay = (row: any) => {
    pStore.setContext({ key: `courseArticle:${String(enid.value ?? '')}`, title: String(breadcrumbTitle.value ?? '') })
    const queue = (tableData.article_list || []).map(buildTrack).filter((t): t is PlayerTrack => !!t)
    const target = buildTrack(row)

    if (!target) {
        ElMessage({ message: '该条目没有可播放的音频地址', type: 'warning' })
        return
    }

    const startIndex = queue.findIndex((t) => t.id === target.id)
    pStore.setQueue(queue, startIndex >= 0 ? startIndex : 0)
}

const appendToPlayerQueueIfActive = (items: any[]) => {
    const ctxKey = `courseArticle:${String(enid.value ?? '')}`
    if (pStore.context?.key !== ctxKey) return
    if (pStore.queue.length === 0) return
    const tracks = (items || []).map(buildTrack).filter((t): t is PlayerTrack => !!t)
    pStore.appendQueue(tracks)
    if (pStore.consumeAutoNext()) {
        return
    }
}

// 懒加载：滚动到底部时加载更多
const loadMoreArticles = async () => {
    if (loading.value || finished.value) return
    
    // 如果 total > 0 且已加载条数 >= total，则 finished
    if (total.value > 0 && tableData.article_list.length >= total.value) {
        finished.value = true
        return
    }
    
    loading.value = true
    try {
        // 使用 maxId 分页
        const res = await ArticleList(enid.value, "", pageSize.value, maxId.value, isReverse.value)
        console.log("xxxxxxxxxxxxxxxx")
        console.log(res)
        // 如果返回数据为空，直接 finished
        if (!res.article_list || res.article_list.length === 0) {
            finished.value = true
            loading.value = false
            return
        }
        tableData.article_list.push(...res.article_list)
        appendToPlayerQueueIfActive(res.article_list)
        // 更新 maxId 为最新一条的 id
        maxId.value = res.article_list[res.article_list.length - 1]?.id
        // 如果本次返回数据不足 pageSize，或总数已达 total，finished
        if (res.article_list.length < pageSize.value || (total.value > 0 && tableData.article_list.length >= total.value)) {
            finished.value = true
        }
    } catch (error) {
        ElMessage({ message: '获取数据失败', type: 'error' })
        finished.value = true // 防止异常时死循环
    } finally {
        loading.value = false
    }
}

const onNeedMoreForPlayer = async (ev: any) => {
    const detail = ev?.detail || {}
    const ctxKey = `courseArticle:${String(enid.value ?? '')}`
    if (detail?.contextKey !== ctxKey) return
    if (finished.value) return
    await loadMoreArticles()
}

const openDialog = () => {
    dialogVisible.value = true
}
const closeDialog = () => {
    dialogVisible.value = false
}

const openDownloadDialog = (row: any) => {
    downloadId.value = row.class_id
    downloadaId.value = row.id
    downloadEnId.value = row.class_enid
    dialogDownloadVisible.value = true
    if (setStore.getDownloadDir == "") {
        ElMessage({
            message: '请设置文件保存目录',
            type: 'warning'
        })
        router.push('/setting')
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
}

const gotoArticleDetail = (row: any) => {
    const total = route.query.total
    pushArticleDetail(row.enid, 'course', {
        class_id: row.class_id,
        class_enid: row.class_enid,
        total: total,
        title: row.title,
        parentTitle: breadcrumbTitle.value
    })
}

const gotoArticleVideo = (row: any) => {
    const pickVideoMediaBaseInfo = (list: any[] | undefined) => {
        if (!list || list.length === 0) return null
        return list.find((m) => m?.media_type === 2) ?? list[0]
    }

    const mediaBase = pickVideoMediaBaseInfo(row?.media_base_info)
    const mediaId = String(mediaBase?.source_id ?? '')
    const securityToken = String(mediaBase?.security_token ?? '')

    if (!mediaId || !securityToken) {
        ElMessage({ message: '未获取到可播放的鉴权信息', type: 'warning' })
        return
    }

    pushVideo({
        from: "course",
        media_id: mediaId,
        security_token: securityToken,
        title: row.title,
        parentTitle: breadcrumbTitle.value
    })
}

</script>

<style scoped>
.breadcrumb-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 24px 32px 16px 32px;
    background: var(--bg-color);
}

.breadcrumb-actions {
    display: flex;
    gap: 16px;
}

.action-btn {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 14px;
    color: var(--text-secondary);
    transition: color 0.2s;
}

.action-btn:hover {
    color: var(--primary-color);
}

.infinite-list-wrapper {
    height: calc(100vh - 140px);
    overflow-y: auto;
    padding: 0 32px 32px 32px;
    
    /* 隐藏滚动条但保留功能 - 清新风格 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE 10+ */
}

.infinite-list-wrapper::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}

.article-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
    padding: 0;
    margin: 0;
    list-style: none;
}

.article-card-wrapper {
    /* Wrapper for layout control if needed */
}

.article-card {
    background: var(--card-bg);
    border-radius: 16px;
    overflow: hidden;
    box-shadow: var(--shadow-soft);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    border: 1px solid var(--border-soft);
    display: flex;
    flex-direction: column;
    height: 100%;
    cursor: pointer;
}

.article-card:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-medium);
    border-color: var(--accent-color);
}

.card-cover {
    position: relative;
    width: 100%;
    padding-top: 56.25%; /* 16:9 Aspect Ratio */
    background-color: var(--fill-color);
    overflow: hidden;
}

.card-image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    transition: transform 0.5s ease;
}

.article-card:hover .card-image {
    transform: scale(1.05);
}

.image-placeholder {
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

.article-card:hover .card-overlay {
    opacity: 1;
}

.overlay-actions {
    display: flex;
    gap: 12px;
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
    z-index: 2;
}

.card-content {
    padding: 16px;
    flex: 1;
    display: flex;
    flex-direction: column;
}

.card-header {
    margin-bottom: 12px;
}

.card-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
    line-height: 1.5;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    overflow: hidden;
    height: 48px;
}

.card-meta {
    display: flex;
    justify-content: space-between;
    color: var(--text-secondary);
    font-size: 13px;
    margin-bottom: 12px;
}

.meta-item {
    display: flex;
    align-items: center;
    gap: 4px;
}

.card-footer {
    margin-top: auto;
    display: flex;
    justify-content: flex-end;
    border-top: 1px solid var(--border-soft);
    padding-top: 12px;
}

.publish-time {
    font-size: 12px;
    color: var(--text-tertiary);
    display: flex;
    align-items: center;
    gap: 4px;
}

.loading-state, .finished-state {
    text-align: center;
    padding: 24px;
    color: var(--text-secondary);
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

@media (max-width: 600px) {
    .article-list {
        grid-template-columns: 1fr;
    }
}
</style>
