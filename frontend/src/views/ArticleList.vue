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
                            <span class="card-title">{{ item.title }}</span>
                        </div>
                        <el-image 
                            :src="item.logo" 
                            class="card-image" 
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
                            <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(item)">下载</el-button>
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

    <el-drawer :title="media?.title" direction="btt" v-model="audioVisible" @close="closeAudio" @open="open"
               @opened="openVideo(media)">
        <div style="position:relative;" v-html="audiohtml"></div>
    </el-drawer>

    <el-drawer :title="media?.title" direction="btt" v-model="vedioVisible" @close="closeVideo" @open="open"
               @opened="openVideo(media)">
        <div style="position:relative;" id="video"></div>
    </el-drawer>

</template>

<script lang="ts" setup>
import {nextTick, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElTable} from 'element-plus'
import {ArrowRight} from '@element-plus/icons-vue'
import {ArticleList, SetDir} from '../../wailsjs/go/backend/App'
import {services} from '../../wailsjs/go/models'
import {useRoute} from 'vue-router'
import {secondToHour} from '../utils/utils'
import DownloadDialog from "../components/DownloadDialog.vue";
import CourseInfo from "../components/CourseInfo.vue";

import videojs from 'video.js'
import "video.js/dist/video-js.css"
import {settingStore} from "../stores/setting";
import { User, Clock, Calendar, Sort, InfoFilled } from '@element-plus/icons-vue'
import { timestampToDate } from '../utils/utils'
import {useAppRouter} from '../composables/useRouter'
import {ROUTE_NAMES} from '../router/routes'

const audioPlayer = ref()
let media = ref()
const audioVisible = ref(false)
const audiohtml = ref('')

const vePlayer = ref()
const vedioVisible = ref(false)
// const audiohtml = ref('')

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

onUnmounted(() => {
    if (audioPlayer.value) {
        audioPlayer.value.dispose()
    }
})

const closeAudio = () => {
    audioVisible.value = false;
    media.value = ''
}

const closeVideo = () => {
    if (vePlayer.value) {
        vePlayer.value.dispose();
        vePlayer.value = '';
    }
    vePlayer.value = false;
    media.value = ''
}

const handlePlay = (row: any) => {
    if (row.video_status == 0) {
        audioVisible.value = true;
        media.value = row
    } else {
        gotoArticleVideo(row)
    }
}

const open = () => {
    audiohtml.value = '<audio id=' + media.value.log_type + media.value.enid + ' controls class="video-js vjs-big-play-centered  vjs-default-skin " style="width:100%;height:100px"></audio>';
}

const openVideo = async (row: any) => {
    console.log(row)
    let videoStatus = row.video_status
    let poster = (videoStatus == 0) ? row.audio.icon : row.video[0].cover_img
    let src = (videoStatus == 0) ? row.audio.mp3_play_url : row.video[0].bitrate_720
    let sType = (videoStatus == 0) ? 'application/x-mpegURL' : 'video/mp4'

    if (videoStatus == 0) {
        // console.log(soc)
        setTimeout(() => {
            nextTick(() => {
                audioPlayer.value = videojs(row.log_type + row.enid, {
                    language: 'zh-Hans',
                    poster: poster,
                    controls: true,
                    sources: [{
                        src: src, type: sType
                    }],
                    controlBar: {
                        remainingTimeDisplay: {
                            displayNegative: false
                        },
                        fullscreenToggle: false
                    },
                    playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2],
                    audioOnlyMode: true,
                    audioPosterMode: true,
                    volume: 0.5,
                    userActions: {
                        hotkeys: function (event: any) {
                            // // `x` key = pause
                            // if (event.which === 88) {
                            //   audioPlayer.value.pause();
                            // }
                            // // `y` key = play
                            // if (event.which === 89) {
                            //   audioPlayer.value.play();
                            // }
                        }
                    }
                }, () => {
                    audioPlayer.value.log("play.....")
                })
            })
        }, 300)
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
        // 如果返回数据为空，直接 finished
        if (!res.article_list || res.article_list.length === 0) {
            finished.value = true
            loading.value = false
            return
        }
        tableData.article_list.push(...res.article_list)
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
    pushVideo({
        from: "course",
        media_id: row.media_base_info![0].source_id,
        security_token: row.media_base_info![0].security_token,
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
  font-weight: 500;
  color: #222;
  flex: 1;
  line-height: 1.4;
  word-break: break-all;
}
.card-meta-row {
  display: flex;
  gap: 24px;
  color: #888;
  font-size: 15px;
  align-items: center;
}
.card-meta i {
  margin-right: 4px;
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