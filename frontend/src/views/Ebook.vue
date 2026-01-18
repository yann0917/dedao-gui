<template>
    <div v-if="groupMode.active" style="display:flex; align-items:center; justify-content:space-between; margin-bottom: 8px;">
        <div style="display:flex; align-items:center; gap: 8px;">
            <el-button type="primary" link @click="exitGroup">返回</el-button>
            <span>{{ groupMode.title }}</span>
        </div>
    </div>
    <el-table :data="tableData.list" v-loading="loading" height="97%" width="100%" :cell-style="{ textAlign: 'left' }"
              :header-cell-style="{textAlign: 'left'}" :row-style="{height: '50px'}" table-layout="auto"
              @sort-change="handleChange">
        <!-- <el-table-column prop="id" label="ID" width="100"/> -->
        <el-table-column prop="title" label="标题" width="280">
            <template #default="scope">
                <div style="display:flex; align-items:center; gap: 8px;">
                    <span>{{ scope.row.title }}</span>
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
                    preview-teleported
                    :src="scope.row.icon"
                    :preview-src-list="[scope.row.icon]"
                    :initial-index="0"
                    style="width: 32px;"
                />
                <span v-else>-</span>
            </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="100"/>
        <el-table-column prop="intro" label="简介" width="400">
            <template #default="scope">
                <el-popover title="简介" trigger="hover" placement="left" :width="480"
                            :disabled="(scope.row.intro || '').length <= 40">
                    <p v-html="(scope.row.intro || '').replaceAll('\n','<br/>')"></p>
                    <template #reference>
                        <span slot="reference"
                              v-if="scope.row.intro && scope.row.intro.length <= 40">{{ scope.row.intro }}</span>
                        <span slot="reference" v-if="scope.row.intro && scope.row.intro.length > 40">{{
                            scope.row.intro.substring(0, 40)
                            + "..."
                            }}</span>
                    </template>
                </el-popover>
            </template>
        </el-table-column>
        <el-table-column prop="progress" label="已读%" width="80" align="center">
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="220">
            <template #default="scope">
              <template v-if="!scope.row.is_group">
                <el-tooltip content="书评">
                  <el-button icon="ChatDotRound" size="small" type="primary" link @click="gotoCommentList(scope.row)">
                  </el-button>
                </el-tooltip>
                <el-tooltip content="详情">
                  <el-button icon="view" size="small" type="primary" link @click="handleProd(scope.row.enid)">
                  </el-button>
                </el-tooltip>
                <el-tooltip content="下载">
                  <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">
                  </el-button>
                </el-tooltip>
                <el-tooltip content="移出书架">
                  <el-button icon="delete" size="small" type="primary" link @click="ebookShelfRemove(scope.row.enid)">
                  </el-button>
                </el-tooltip>
              </template>
            </template>
        </el-table-column>
    </el-table>
    <Pagination :key="paginationKey" :total="total" @pageChange="handleChangePage"></Pagination>
    <ebook-info v-if="dialogVisible" :enid="prodEnid" :dialog-visible="dialogVisible" @close="closeDialog"></ebook-info>

    <download-dialog
            v-if="dialogDownloadVisible"
            :dialog-visible="dialogDownloadVisible"
            :download-type-options="downloadTypeOptions"
            :prod-type="2"
            :download-id="downloadId"
            :en-id="downloadEnId"
            @close="closeDownloadDialog">
    </download-dialog>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElTable, ElMessageBox} from 'element-plus'
import {
  CourseCategory,
  CourseList,
  CourseGroupList,
  EbookShelfRemove,
  SetDir
} from '../../wailsjs/go/backend/App'
import {services} from '../../wailsjs/go/models'
import Pagination from '../components/Pagination.vue'
import EbookInfo from '../components/EbookInfo.vue'
import DownloadDialog from "../components/DownloadDialog.vue";

import {useAppRouter} from '../composables/useRouter'
import {userStore} from '../stores/user'
import {settingStore} from "../stores/setting";
import {Local} from '../utils/storage'

const store = userStore()
const setStore = settingStore()
const { pushEbookComment, pushSetting, pushLogin } = useAppRouter()
const loading = ref(true)
const page = ref(1)
const total = ref(0)
const outerTotal = ref(0)
const pageSize = ref(15)
const paginationKey = ref(0)
const searchInfo = ref({})
const dialogVisible = ref(false)
const prodEnid = ref("")

const groupMode = reactive({
  active: false,
  groupId: 0,
  title: '',
})

const dialogDownloadVisible = ref(false)
const downloadType = ref(1)
const downloadId = ref(0)
const downloadEnId = ref('')
const downloadTypeOptions = [
    {value: 1, label: "HTML"}, {value: 2, label: "PDF"}, {value: 3, label: "EPUB"}
]

let tableData = reactive(new services.CourseList)

onMounted(() => {
    CourseCategory().then(result => {
        result.forEach((item, key) => {
            if (item.category == "ebook") {
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
        ? CourseGroupList("ebook", "study", groupMode.groupId, page.value, pageSize.value)
        : CourseList("ebook", "study", page.value, pageSize.value)
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

const handleProd = (enid: string) => {
    prodEnid.value = enid
    dialogVisible.value = true
}

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


function handleChange() {
}

getTableData()

const openDialog = () => {
    dialogVisible.value = true
}
const closeDialog = () => {
    //   initForm()
    dialogVisible.value = false
}

const ebookShelfRemove = async (enid: string) => {
  ElMessageBox.confirm(
      '是否移出书架?',
      '',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        EbookShelfRemove([enid]).then((info) => {
          console.log(info)
          getTableData()
        }).catch((error) => {
          ElMessage({
            message: error,
            type: 'warning'
          })
        })
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: '取消',
        })
      })
  return
}

const openDownloadDialog = (row: any) => {
    downloadId.value = row.id
    downloadEnId.value = row.enid
    dialogDownloadVisible.value = true
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
    //   initForm()
    downloadType.value = 1
    dialogDownloadVisible.value = false
}

const gotoCommentList = (row: any) => {
    pushEbookComment({
        id: row.id,
        enid: row.enid,
        total: row.publish_num,
        title: row.title
    })
}

</script>


<style scoped>
.card-header .el-row {
    height: 100%;
}

.el-row {
    margin-bottom: 10px;
}

.el-row:last-child {
    margin-bottom: 0;
}

.el-col {
    border-radius: 4px;
}

.divider-text {
    font-size: larger;
}

.catalog,
.author-info,
.book-intro {
    line-height: 1.6;
}

.el-tag {
    margin-right: 5px;
    /* height: auto; */
    text-align: center;
}
</style>
