<template>
    <div class="breadcrumb-container">
        <el-breadcrumb :separator-icon="ArrowRight">
            <el-breadcrumb-item :to="{ name: ROUTE_NAMES.COURSE }">课程列表</el-breadcrumb-item>
            <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
        </el-breadcrumb>
        <el-button 
            type="primary" 
            link 
            @click="showCourseDetail"
            class="detail-btn"
        >
            课程详情
            <el-icon class="el-icon--right"><InfoFilled /></el-icon>
        </el-button>
        <el-button 
            type="primary" 
            link 
            @click="toggleSort"
            class="sort-btn"
        >
            {{ isReverse ? '正序' : '倒序' }}
            <el-icon class="el-icon--right"><Sort /></el-icon>
        </el-button>
    </div>
    <div
        class="infinite-list-wrapper"
        v-infinite-scroll="loadMoreArticles"
        :infinite-scroll-disabled="loading || finished"
        infinite-scroll-distance="10"
        @scroll="onScroll"
    >
        <ul class="article-list">
            <li v-for="item in tableData.article_list" :key="item.id" class="article-card">
                <el-card shadow="hover" class="article-el-card">
                    <div class="card-content">
                        <div class="card-header">
                            <el-tooltip
                                effect="dark"
                                :content="item.title"
                                placement="top"
                                :show-after="500"
                            >
                                <span class="card-title">{{ item.title }}</span>
                            </el-tooltip>
                            <el-tag v-if="isLearned(item)" type="success" size="small">已学习</el-tag>
                        </div>
                        <el-image 
                            :src="item.logo" 
                            class="card-image"
                            fit="cover"
                            :preview-teleported="true"
                            :preview-src-list="[item.logo]"
                        />
                        <div class="card-meta-row">
                            <span class="card-meta">
                                <el-icon><Clock /></el-icon>
                                {{ item.video_status == 1 ? secondToHour(item.video?.[0]?.duration !== undefined ? item.video[0].duration : 0) : (item.audio?.duration ? secondToHour(item.audio.duration) : '') }}
                            </span>
                            <span class="card-meta">
                                <el-icon><User /></el-icon>
                                {{ item.cur_learn_count }}
                            </span>
                            <span class="card-meta">
                                <el-icon><Calendar /></el-icon>
                                {{ item.publish_time ? timestampToDate(item.publish_time) : '' }}
                            </span>
                        </div>
                        <div class="card-actions">
                            <el-button icon="VideoPlay" size="small" type="primary" link @click="handlePlay(item)"
                                       v-if="item.audio_alias_ids.length || item.video_status">播放
                            </el-button>
                            <el-button icon="Memo" size="small" type="primary" link @click="gotoArticleDetail(item)">文稿</el-button>
                            <el-button
                              icon="download"
                              size="small"
                              type="primary"
                              link
                              v-if="canDownload(item)"
                              @click="openDownloadDialog(item)"
                            >下载</el-button>
                        </div>
                    </div>
                </el-card>
            </li>
        </ul>
        <div v-if="loading" class="loading" style="text-align:center;padding:12px;">加载中...</div>
        <div v-if="finished && hasScrolled" class="finished" style="text-align:center;padding:12px;color:#999;">没有更多了</div>
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
import {ArrowRight} from '@element-plus/icons-vue'
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

onMounted(() => {
    window.addEventListener('player:needMore', onNeedMoreForPlayer as any)
})

onUnmounted(() => {
    window.removeEventListener('player:needMore', onNeedMoreForPlayer as any)
})

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
    margin-bottom: 20px;
    padding-right: 20px;
}
.sort-btn {
    font-size: 14px;
}
.infinite-list-wrapper {
    height: calc(100vh - 120px);
    overflow-y: auto;
}
.article-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 24px;
  padding: 0;
  margin: 0;
  list-style: none;
}
.article-card {
  /* 让卡片高度自适应内容 */
}
.article-el-card {
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  transition: box-shadow 0.2s;
}
.article-el-card:hover {
  box-shadow: 0 6px 24px rgba(0,0,0,0.12);
}
.card-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 160px;
}
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}
.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  line-height: 1.5;
  flex: 1;
  
  /* Multi-line truncation */
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 8px;
  height: 48px; /* Ensure consistent height for alignment */
}
.card-image {
  width: 100%;
  height: 192px; /* 16:9 aspect ratio approx or fixed height */
  border-radius: 8px;
  display: block;
}
.card-meta-row {
  display: flex;
  justify-content: space-between;
  color: #6b7280;
  font-size: 13px;
  align-items: center;
  margin-top: 8px;
}
.card-meta {
  display: flex;
  align-items: center;
  gap: 4px;
}
.card-meta i {
  font-size: 16px;
  vertical-align: middle;
}
.card-actions {
  margin-top: auto;
  display: flex;
  gap: 16px;
  justify-content: flex-end;
}
@media (max-width: 600px) {
  .article-list {
    grid-template-columns: 1fr;
  }
  .card-title {
    font-size: 17px;
  }
  .card-meta-row {
    font-size: 13px;
  }
}
</style>
