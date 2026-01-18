<template>
    <div v-if="groupMode.active" style="display:flex; align-items:center; justify-content:space-between; margin-bottom: 8px;">
        <div style="display:flex; align-items:center; gap: 8px;">
            <el-button type="primary" link @click="exitGroup">返回</el-button>
            <span>{{ groupMode.title }}</span>
        </div>
    </div>
    <el-table 
        :data="tableData.list" 
        v-loading="loading" 
        height="97%" 
        width="100%" 
        class="custom-table"
        table-layout="auto"
    >
        <!-- <el-table-column prop="id" label="ID" width="100"/> -->
        <el-table-column prop="title" label="标题" width="280">
            <template #default="scope">
                <div style="display: flex; align-items: center; gap: 8px;">
                    <span>{{ scope.row.title }}</span>
                    <el-tag v-if="scope.row.type === 1013" type="warning" size="small">名家讲书</el-tag>
                    <el-tag v-if="scope.row.is_group" type="info" size="small">分组</el-tag>
                    <span v-if="scope.row.is_group">共{{ scope.row.course_num || 0 }}本</span>
                    <el-button v-if="scope.row.is_group" type="primary" link @click.stop="enterGroup(scope.row)">进入</el-button>
                </div>
            </template>
        </el-table-column>
        <el-table-column prop="icon" label="封面" width="80">
            <template #default="scope">
                <el-image
                    v-if="scope.row.icon"
                    :src="scope.row.icon"
                    :preview-teleported="true"
                    :preview-src-list="[scope.row.icon]"
                    style="width: 32px;"
                />
                <span v-else>-</span>
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
                            :disabled="(scope.row.intro || '').length <= 30"
                            :content="scope.row.intro || ''">
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

        <el-table-column fixed="right" label="操作" width="240">
            <template #default="scope">
                <template v-if="!scope.row.is_group">
                    <el-tooltip content="播放">
                    <el-button icon="VideoPlay" size="small" type="primary" link @click="handlePlay(scope.row)">
                    </el-button>
                    </el-tooltip>
                    <el-tooltip content="文稿">
                    <el-button icon="Memo" size="small" type="primary" link @click="gotoArticleDetail(scope.row)">
                    </el-button>
                    </el-tooltip>
                    <el-tooltip content="详情">
                    <el-button icon="view" size="small" type="primary" link @click="handleProd(scope.row)">
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="下载">
                    <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">
                    </el-button>
                    </el-tooltip>
                </template>

            </template>
        </el-table-column>
    </el-table>
    <Pagination :key="paginationKey" :total="total" @pageChange="handleChangePage"></Pagination>
    <audio-info v-if="dialogVisible" :enid="prodEnid" :dialog-visible="dialogVisible" @close="closeDialog"></audio-info>
    <outside-info v-if="outsideVisible" :enid="prodEnid" :dialog-visible="outsideVisible" @close="closeDialog"></outside-info>

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
                :stroke-width="20"
                :text-inside="true"
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
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import 'element-plus/es/components/message/style/css'
import { ElMessage } from 'element-plus'
import { AudioDetailAlias, CourseCategory, CourseGroupList, CourseList, OdobDownload, SetDir } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import Pagination from '../components/Pagination.vue'
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
const loading = ref(true)
const page = ref(1)
const total = ref(0)
const outerTotal = ref(0)
const pageSize = ref(15)
const paginationKey = ref(0)
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
// 分页
const handleChangePage = (item: any) => {
    page.value = item.page
    pageSize.value = item.pageSize
    getTableData()
}

const getTableData = async () => {
    loading.value = true
    const fetcher = groupMode.active
        ? CourseGroupList("odob", "study", groupMode.groupId, page.value, pageSize.value)
        : CourseList("odob", "study", page.value, pageSize.value)
    await fetcher.then((table) => {
        loading.value = false
        Object.assign(tableData, table)
        total.value = groupMode.active ? (table.total || 0) : outerTotal.value
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
    paginationKey.value += 1
    getTableData()
}

const exitGroup = () => {
    groupMode.active = false
    groupMode.groupId = 0
    groupMode.title = ''
    page.value = 1
    total.value = outerTotal.value
    paginationKey.value += 1
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
/* 自定义表格样式 */
.custom-table {
  text-align: left;
}

.custom-table :deep(.el-table__cell) {
  text-align: left;
  height: 50px;
}

.custom-table :deep(.el-table__header .el-table__cell) {
  text-align: left;
  font-weight: 500;
  background-color: var(--fill-color);
  color: var(--text-primary);
}

.custom-table :deep(.el-table__body .el-table__cell) {
  text-align: left;
  background-color: var(--card-bg);
  color: var(--text-primary);
  border-color: var(--border-soft);
}

.custom-table :deep(.el-table__row:hover .el-table__cell) {
  background-color: var(--fill-color);
  color: var(--text-primary);
}

/* 暗色模式下的表格样式 */
.theme-dark .custom-table :deep(.el-table__header .el-table__cell) {
  background-color: var(--fill-color) !important;
  color: var(--text-primary) !important;
  border-color: var(--border-soft) !important;
}

.theme-dark .custom-table :deep(.el-table__body .el-table__cell) {
  background-color: var(--card-bg) !important;
  color: var(--text-primary) !important;
  border-color: var(--border-soft) !important;
}

.theme-dark .custom-table :deep(.el-table__row:hover .el-table__cell) {
  background-color: var(--fill-color) !important;
  color: var(--text-primary) !important;
}


/* 表格标签样式 */
.custom-table :deep(.el-tag) {
  background-color: var(--fill-color);
  border-color: var(--border-soft);
  color: var(--text-primary);
}

/* 表格图片样式 */
.custom-table :deep(.el-image) {
  border-radius: 4px;
  border: 1px solid var(--border-soft);
}

/* 表格弹出层样式 */
.custom-table :deep(.el-popover) {
  background-color: var(--card-bg);
  border-color: var(--border-soft);
  color: var(--text-primary);
}

.custom-table :deep(.el-popover .el-popover__title) {
  color: var(--text-primary);
}

.custom-table :deep(.el-popover .el-popover__content) {
  color: var(--text-secondary);
}
</style>
  
