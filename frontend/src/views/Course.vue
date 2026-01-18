<template>
    <div v-if="groupMode.active" style="display:flex; align-items:center; justify-content:space-between; margin-bottom: 8px;">
        <div style="display:flex; align-items:center; gap: 8px;">
            <el-button type="primary" link @click="exitGroup">返回</el-button>
            <span>{{ groupMode.title }}</span>
        </div>
    </div>
    <el-table :data="tableData.list" v-loading="loading" height="97%" style="width: 100%"
        :cell-style="{ textAlign: 'left' }" :header-cell-style="{ textAlign: 'left' }"
        table-layout="auto">
        <!-- <el-table-column prop="class_id" label="ID" width="80" /> -->
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
                    :src="scope.row.icon"
                    :preview-teleported="true"
                    :preview-src-list="[scope.row.icon]"
                    style="width: 32px;"
                />
                <span v-else>-</span>
            </template>
        </el-table-column>
        <el-table-column prop="intro" label="简介" width="280" />

        <el-table-column prop="price" label="价格" width="100" />

        <el-table-column prop="publish_num" label="已更新" width="100" >
            <template #default="scope">
                <span v-if="!scope.row.is_group">{{ scope.row.publish_num || 0 }}/{{ scope.row.course_num || 0 }}</span>
                <span v-else>-</span>
            </template>
        </el-table-column>
        <el-table-column prop="progress" label="已学%" width="80" />
        <el-table-column fixed="right" label="操作" width="240">
            <template #default="scope">
                <template v-if="!scope.row.is_group">
                    <el-button icon="list" size="small" type="primary" link @click="gotoArticleList(scope.row)">章节列表
                    </el-button>
                    <el-button icon="view" size="small" type="primary" link @click="handleProd(scope.row.enid)">详情
                    </el-button>
                    <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">下载
                    </el-button>
                </template>

            </template>
        </el-table-column>
    </el-table>
    
    <pagination :key="paginationKey" :total="total" @pageChange="handleChangePage"></pagination>
    <course-info v-if="dialogVisible" :enid= "prodEnid" :dialog-visible="dialogVisible" @close="closeDialog"></course-info>
    <download-dialog
        v-if="dialogDownloadVisible"
        :dialog-visible="dialogDownloadVisible"
        :download-type-options="downloadTypeOptions"
        :prod-type="66"
        :download-id="downloadId"
        :en-id="downloadEnId"
        @close="closeDownloadDialog">
    </download-dialog>
</template>
  
<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { ElTable, ElMessage } from 'element-plus'
import {CourseList, CourseCategory, CourseGroupList, SetDir} from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { userStore } from '../stores/user';
import { settingStore } from '../stores/setting';
import { useAppRouter } from '../composables/useRouter';
import Pagination from '../components/Pagination.vue'
import CourseInfo from '../components/CourseInfo.vue'
import DownloadDialog from "../components/DownloadDialog.vue";
import { Local } from '../utils/storage';

const store = userStore()
const setStore = settingStore()
const { pushLogin, pushCourseDetail, pushSetting } = useAppRouter()

const loading = ref(true)
const page = ref(1)
const total = ref(0)
const outerTotal = ref(0)
const pageSize = ref(15)
const paginationKey = ref(0)
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
    { value: 1, label: "MP3" }, { value: 2, label: "PDF" }, { value: 3, label: "Markdown" }
]

let tableData = reactive(new services.CourseList)
let courseInfo = reactive(new services.CourseInfo)
courseInfo.class_info = new services.ClassInfo
courseInfo.intro_article = new services.ArticleIntro

onMounted(() => {
    CourseCategory().then(result => {
        result.forEach((item, key) => {
            if (item.category == "bauhinia") {
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
        ? CourseGroupList("bauhinia", "study", groupMode.groupId, page.value, pageSize.value)
        : CourseList("bauhinia", "study", page.value, pageSize.value)

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

getTableData()

const handleProd = (enid:string)=>{
  prodEnid.value = enid
  dialogVisible.value = true
}

const gotoArticleList = (row: any) => {
    pushCourseDetail(row.id, {
        enid: row.enid,
        total: row.publish_num,
        title: row.title
    })
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

const openDialog = () => {
    dialogVisible.value = true
}
const closeDialog = () => {
    //   initForm()
    dialogVisible.value = false
}

const openDownloadDialog = (row: any) => {
    downloadId.value = row.class_id
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
</script>
  
<style scoped>
.card-header.el-row {
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
.el-tag {
  margin-right: 5px;
    /* height: auto; */
  text-align: center;
}
</style>
