<template>
    <el-table stripe :data="tableData.list" v-loading="loading" height="97%" style="width: 100%"
        :cell-style="{ textAlign: 'left' }" :header-cell-style="{ textAlign: 'left' }" :row-style="{ height: '30px' }"
        table-layout="auto">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" width="280" />
        <el-table-column prop="icon" label="封面" width="80">
            <template #default="scope">
                <el-image :src="scope.row.icon" :preview-src-list="[scope.row.icon]" :initial-index="0"
                    style="width: 32px;" />
            </template>
        </el-table-column>
        <el-table-column prop="intro" label="简介" width="280" />

        <el-table-column prop="price" label="价格" width="100" />

        <el-table-column prop="publish_num" label="已更新" width="100" />
        <el-table-column prop="progress" label="已学%" width="100" />
        <el-table-column fixed="right" label="操作" width="200">
            <template #default="scope">
                <el-button icon="list" size="small" type="primary" link @click="gotoArticleList(scope.row)">章节列表
                </el-button>
                <el-button icon="view" size="small" type="primary" link @click="handleProd(scope.row.enid)">详情
                </el-button>
                <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">下载
                </el-button>

            </template>
        </el-table-column>
    </el-table>
    <pagination :total="total" @pageChange="handleChangePage"></pagination>
    <course-info v-if="dialogVisible" :enid= "prodEnid" :dialog-visible="dialogVisible" @close="closeDialog"></course-info>
    <download-dialog
        v-if="dialogDownloadVisible"
        :dialog-visible="dialogDownloadVisible"
        :download-type-options="downloadTypeOptions"
        :prod-type="66"
        :download-id="downloadId"
        @close="closeDownloadDialog">
    </download-dialog>
</template>
  
<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { ElTable, ElMessage } from 'element-plus'
import { CourseList, CourseCategory } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { useRouter } from 'vue-router'
import { userStore } from '../stores/user';
import Pagination from '../components/Pagination.vue'
import CourseInfo from '../components/CourseInfo.vue'
import DownloadDialog from "../components/DownloadDialog.vue";
import { Local } from '../utils/storage';

const store = userStore()
const router = useRouter()

const loading = ref(true)
const page = ref(1)
const total = ref(0)
const pageSize = ref(15)
const dialogVisible = ref(false)
const prodEnid = ref("")


const dialogDownloadVisible = ref(false)
const downloadType = ref(1)
const downloadId = ref(0)

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

// 分页
const handleChangePage = (item: any) => {
    page.value = item.page
    pageSize.value = item.pageSize
    getTableData()
}


const getTableData = async () => {
    await CourseList("bauhinia", "study", page.value, pageSize.value).then((table) => {
        loading.value = false
        Object.assign(tableData, table)
        // console.log(tableData)
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
    router.push({
        path: `/course/${row.id}`,
        query: {
            enid: row.enid,
            total: row.publish_num,
            title: row.title
        }
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
    downloadId.value = row.id
    dialogDownloadVisible.value = true
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