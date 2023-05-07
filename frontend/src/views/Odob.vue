<template>
    <el-table :data="tableData.list" v-loading="loading" height="97%" width="100%" :cell-style="{ textAlign: 'left' }"
              :header-cell-style="{ textAlign: 'left' }" :row-style="{ height: '50px' }" table-layout="auto" stripe>
        <el-table-column prop="id" label="ID" width="100"/>
        <el-table-column prop="title" label="标题" width="280"/>
        <el-table-column prop="icon" label="封面" width="80">
            <template #default="scope">
                <el-image :src="scope.row.icon" :preview-src-list="[scope.row.icon]" :initial-index="0"
                          style="width: 32px;"/>
            </template>
        </el-table-column>
        <el-table-column prop="duration" label="时长" width="100">
            <template #default="scope">
                {{ secondToHour(scope.row.duration) }}
            </template>
        </el-table-column>
        <el-table-column prop="intro" label="简介" width="300">
            <template #default="scope">
                <el-popover title="简介" trigger="hover" placement="top" :width="480"
                            :disabled="scope.row.intro.length <= 30"
                            :content="scope.row.intro">
                    <template #reference>
                        <span slot="reference" v-if="scope.row.intro && scope.row.intro.length <= 30">{{
                            scope.row.intro
                            }}</span>
                        <span slot="reference" v-if="scope.row.intro && scope.row.intro.length > 30">{{
                            scope.row.intro.substring(0,
                                30)
                            + "..."
                            }}</span>
                    </template>
                </el-popover>
            </template>
        </el-table-column>
        <el-table-column prop="progress" label="已学%" width="100"/>

        <el-table-column fixed="right" label="操作" width="200">
            <template #default="scope">
                <el-button icon="VideoPlay" size="small" type="primary" link @click="handlePlay(scope.row)">播放
                </el-button>
                <el-button icon="Memo" size="small" type="primary" link @click="gotoArticleDetail(scope.row)">文稿
                </el-button>
                <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">下载
                </el-button>

            </template>
        </el-table-column>
    </el-table>
    <Pagination :total="total" @pageChange="handleChangePage"></Pagination>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="60%" :before-close="closeDialog">
        <el-space wrap>
            <el-card v-for="i in 1" :key="i" class="box-card" width="100%">
                <template #header>
                    <div class="card-header">
                        <el-row :gutter="5">
                            <el-col :span="4">
                                <el-row>
                                    <el-avatar :size="84" :src="courseInfo.class_info.lecturer_avatar" left fit="fill"/>
                                </el-row>
                                <el-row>
                                    <el-tag type="success">{{ courseInfo.class_info.lecturer_name }}</el-tag>
                                </el-row>
                            </el-col>
                            <el-col :span="18">
                                <p style="text-align:left">{{ courseInfo.class_info.lecturer_intro }}</p>
                            </el-col>
                        </el-row>
                    </div>
                </template>
                <div style="text-align:left">
                    <span v-html="courseInfo.class_info.highlight?.replaceAll('\n', '<br/>')"></span>
                </div>

            </el-card>
        </el-space>

    </el-dialog>

    <el-dialog v-model="dialogDownloadVisible" title="请选择下载格式" align-center center width="30%">
        <el-form >
            <el-form-item label="下载格式" label-width="80px">
                <el-select v-model="downloadType" placeholder="请选择下载格式">
                    <el-option v-for="item in downloadTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
            <el-progress v-show="percentage"
                :percentage="percentage"
                status="success"
                :indeterminate="true"
                :duration="3"
            ><span>{{content}}</span></el-progress>
        </el-form>
        <template #footer>
      <span class="dialog-footer">
        <el-button @click="closeDownloadDialog">取消</el-button>
        <el-button type="primary" @click="download(downloadId, downloadType)">
          确认
        </el-button>
      </span>
        </template>

    </el-dialog>

    <el-drawer direction="btt" :title="media?.title" v-model="videoVisible" size="30%" @close="closeVideo" @open="open"
               @opened="openVideo(media)">
        <div style="position:relative;" v-html="videohtml"></div>
    </el-drawer>


</template>

<script lang="ts" setup>
import {nextTick, onMounted, onUnmounted, reactive, ref} from 'vue'
import 'element-plus/es/components/message/style/css'
import {ElMessage, ElTable} from 'element-plus'
import {CourseCategory, CourseList, OdobDownload, SetDir} from '../../wailsjs/go/backend/App'
import {services} from '../../wailsjs/go/models'
import Pagination from '../components/Pagination.vue'
import {useRouter} from 'vue-router'
import {userStore} from '../stores/user';
import {settingStore} from "../stores/setting";
import {Local} from '../utils/storage';

import {secondToHour} from '../utils/utils'
import videojs from 'video.js'
import "video.js/dist/video-js.css"

const videoPlayer = ref()
let media = ref()
const videoVisible = ref(false)
const videohtml = ref('')

const store = userStore()
const setStore = settingStore()
const router = useRouter()
const loading = ref(true)
const page = ref(1)
const total = ref(0)
const pageSize = ref(15)
const searchInfo = ref({})
const dialogVisible = ref(false)

const dialogDownloadVisible = ref(false)
const downloadType = ref(1)
const downloadId = ref(0)
let downloadData = reactive(new services.Course)
const downloadTypeOptions = [
    {value: 1, label: "MP3"}, {value: 2, label: "PDF"}, {value: 3, label: "Markdown"}
]

let percentage=ref(0)
let content=ref('')

const closeVideo = () => {
    // if (videoPlayer.value) {
    //   videoPlayer.value.dispose();
    //   videohtml.value = '';
    // }
    videoVisible.value = false;
    // media.value = ''
}

const handlePlay = (row: any) => {
    videoVisible.value = true;
    media.value = row
}

const open = () => {
    videohtml.value = '<audio id=' + media.value.log_type + media.value.enid + ' controls class="video-js vjs-big-play-centered vjs-default-skin" style="width:100%;height:80px" muted></audio>';
}

const openVideo = async (row: any) => {
    console.log(row)

    setTimeout(() => {
        nextTick(() => {

                videoPlayer.value = videojs(row.log_type + row.enid, {
                    language: 'zh-Hans',
                    poster: row.audio_detail.icon,
                    controls: true,
                    sources: [
                        {
                            src: row.audio_detail.mp3_play_url,
                            type: 'application/x-mpegURL',
                        }
                    ],
                    controlBar: {
                        remainingTimeDisplay: {
                            displayNegative: false
                        },
                        fullscreenToggle: false
                    },
                    playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2],
                    // aspectRatio:'1:1',
                    // audioOnlyMode:true,
                    audioPosterMode: true,
                    userActions: {
                        hotkeys: function (event: any) {

                            // `x` key = pause
                            if (event.which === 88) {
                                videoPlayer.value.pause();
                            }
                            // `y` key = play
                            if (event.which === 89) {
                                videoPlayer.value.play();
                            }
                        }
                    }
                }, () => {
                    videoPlayer.value.log("play.....")
                })
            }
        )
    }, 300)
}

let tableData = reactive(new services.CourseList)
let courseInfo = reactive(new services.CourseInfo)

onMounted(() => {
    CourseCategory().then(result => {
        result.forEach((item, key) => {
            if (item.category == "odob") {
                total.value = item.count
            }
        })

    }).catch((error) => {
        if (error == '401 Unauthorized') {
            store.user = null
            router.push("/user/login")
        }
        Local.remove("cookies")
        Local.remove("userStore")
    })

})

onUnmounted(() => {
    if (videoPlayer.value) {
        videoPlayer.value.dispose()
    }
})
// 分页
const handleChangePage = (item: any) => {
    page.value = item.page
    pageSize.value = item.pageSize
    getTableData()
}

const getTableData = async () => {
    await CourseList("odob", "study", page.value, pageSize.value).then((table) => {
        loading.value = false
        Object.assign(tableData, table)
        console.log(tableData)
    }).catch((error) => {
        loading.value = false
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
}
getTableData()


const openDownloadDialog = (row: any) => {
    downloadId.value = row.id
    dialogDownloadVisible.value = true

    Object.assign(downloadData, row)
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
    percentage.value = 0
    content.value = ''
}

const download = async (id: number, dType: number) => {
    percentage.value = 66
    content.value = '下载中'
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
    let id = row.audio_detail.alias_id
    router.push({path: `/odob/${id}`, query: {from: "odob"}})
}

</script>
  