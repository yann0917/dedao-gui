<template>
  <el-breadcrumb :separator-icon="ArrowRight">
    <el-breadcrumb-item :to="{ path: '/course' }">课程列表</el-breadcrumb-item>
    <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
  </el-breadcrumb>
  <el-table :data="tableData.article_list" v-loading="loading" height="97%" style="width: 100%"
    :cell-style="{ textAlign: 'left' }" :header-cell-style="{ textAlign: 'left' }" :row-style="{ height: '50px' }"
    table-layout="auto" stripe>
    <el-table-column prop="id" label="ID" width="100" />
    <el-table-column prop="title" label="标题" width="380" />
    <el-table-column prop="audio.duration" label="时长" width="100">
      <template #default="scope">
          <span v-if="scope.row.video_status == 1">{{ secondToHour(scope.row.video[0]?.duration) }}</span>
          <span v-else-if="scope.row.audio?.duration">{{ secondToHour(scope.row.audio?.duration) }}</span>
      </template>
    </el-table-column>
    <el-table-column prop="summary" label="简介">
      <template #default="scope">
        <el-popover title="简介" trigger="hover" placement="left" :width="480" :disabled="scope.row.summary.length <= 30">
          <p v-html="scope.row.summary?.replaceAll('\n', '<br/>')"></p>
          <template #reference>
            <span slot="reference" v-if="scope.row.summary && scope.row.summary.length <= 30">{{ scope.row.summary
            }}</span>
            <span slot="reference" v-if="scope.row.summary && scope.row.summary.length > 30">{{
                scope.row.summary.substring(0, 30)
                + "..."
            }}</span>
          </template>
        </el-popover>
      </template>
    </el-table-column>
    <el-table-column prop="cur_learn_count" label="学习人数" width="100" />
    <el-table-column fixed="right" label="操作" width="200">
      <template #default="scope">
        <el-button icon="VideoPlay" size="small" type="primary" link @click="handlePlay(scope.row)"
          v-if="scope.row.audio_alias_ids.length || scope.row.video_status">播放</el-button>
        <el-button icon="Memo" size="small" type="primary" link @click="gotoArticleDetail(scope.row)">文稿</el-button>
        <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">下载
        </el-button>

      </template>
    </el-table-column>
  </el-table>
  <el-pagination background v-model:currentPage="page" v-model:page-size="pageSize" :page-sizes="pageSizes"
    :total="total" :layout="layout" @current-change="handleChangePage" @size-change="handleSizeChange" />

  <download-dialog
      v-if="dialogDownloadVisible"
      :dialog-visible="dialogDownloadVisible"
      :download-type-options="downloadTypeOptions"
      :prod-type="66"
      :download-id="downloadId"
      :article-id="downloadaId"
      @close="closeDownloadDialog">
  </download-dialog>

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
import { ref, reactive, onMounted, onUnmounted,watch, nextTick } from 'vue'
import { ElTable, ElMessage } from 'element-plus'
import { ArrowRight } from '@element-plus/icons-vue'
import { ArticleList, SetDir} from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { useRoute, useRouter } from 'vue-router'
import { secondToHour } from '../utils/utils'
import DownloadDialog from "../components/DownloadDialog.vue";

import videojs from 'video.js'
import "video.js/dist/video-js.css"
import {settingStore} from "../stores/setting";

const audioPlayer = ref()
let media = ref()
const audioVisible = ref(false)
const audiohtml = ref('')

const vePlayer = ref()
const vedioVisible = ref(false)
// const audiohtml = ref('')


const loading = ref(true)
const page = ref(1)
let total = ref(0)
const pageSize = ref(30)
const dialogVisible = ref(false)
const layout = ref('total, sizes, next')
const pageSizes = ref([10, 15, 20, 30, 50]);
const route = useRoute()
const setStore = settingStore()
const router = useRouter()

const dialogDownloadVisible = ref(false)
const downloadType = ref(1)
const downloadId = ref(0)
const downloadaId = ref(0)
const downloadTypeOptions = [
  { value: 1, label: "MP3" }, { value: 2, label: "PDF" }, { value: 3, label: "Markdown" }
]

let id = ref()
let enid = ref()
let maxId = ref()
let tableData = reactive(new services.ArticleList)
let breadcrumbTitle = ref()

onMounted(() => {
  watch(() => {
    id.value = route.params.id
    enid.value = route.query.enid
    total.value = Number(route.query.total)
    breadcrumbTitle.value = route.query.title
  },
    () => getTableData(),
    { immediate: true }
  )

})

onUnmounted(() => {
  if (audioPlayer.value) {
    audioPlayer.value.dispose()
  }
})

const closeAudio = () => {
  // if (audioPlayer.value) {
  //   audioPlayer.value.dispose();
  //   audiohtml.value = '';
  // }
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
  audiohtml.value = '<audio id='+media.value.log_type+media.value.enid+' controls class="video-js vjs-big-play-centered  vjs-default-skin " style="width:100%;height:100px"></audio>';
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
        audioPlayer.value = videojs(row.log_type+row.enid, {
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

// 分页
const handleChangePage = (item: any) => {
  page.value = item.page
  getTableData()
}

const handleSizeChange = (item: any) => {
  pageSize.value = item
  maxId.value = 0
  getTableData()
}

const getTableData = async () => {
  await ArticleList(enid.value, "", pageSize.value, maxId.value).then((table) => {
    loading.value = false
    Object.assign(tableData, table)
    console.log(table)
    maxId.value = tableData.article_list[tableData.article_list.length - 1]?.id

  }).catch((error) => {
    loading.value = false
    ElMessage({
      message: error,
      type: 'warning'
    })
  })
}


const openDialog = () => {
  dialogVisible.value = true
}
const closeDialog = () => {
  //   initForm()
  dialogVisible.value = false
}

const openDownloadDialog = (row: any) => {
  downloadId.value = row.class_id
  downloadaId.value = row.id
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
  //   initForm()
  downloadType.value = 1
  dialogDownloadVisible.value = false
}

const gotoArticleDetail = (row: any) => {
  let total = route.query.total
  router.push({
    path: `/article/${row.enid}`,
    query: {
      from: "course",
      class_id: row.class_id,
      class_enid: row.class_enid,
      total: total,
      title: row.title,
      parentTitle: breadcrumbTitle.value
    }
  })
}

const gotoArticleVideo = (row: any) => {
  router.push({
    path: `/video`,
    query: {
      from: "course",
      media_id: row.media_base_info![0].source_id,
      security_token: row.media_base_info![0].security_token,
      title: row.title,
      parentTitle: breadcrumbTitle.value
    }
  })
}
</script>