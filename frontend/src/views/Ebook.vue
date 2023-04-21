<template>
    <el-table :data="tableData.list" v-loading="loading" height="97%" width="100%" :cell-style="{ textAlign: 'left' }"
              :header-cell-style="{textAlign: 'left'}" :row-style="{height: '50px'}" table-layout="auto"
              @sort-change="handleChange">
        <el-table-column prop="id" label="ID" width="100"/>
        <el-table-column prop="title" label="标题" width="280"/>
        <el-table-column prop="icon" label="封面" width="80">
            <template #default="scope">
                <el-image :src="scope.row.icon" :preview-src-list="[scope.row.icon]" :initial-index="0"
                          style="width: 32px;"/>
            </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="100"/>
        <el-table-column prop="intro" label="简介" width="320">
            <template #default="scope">
                <el-popover title="简介" trigger="hover" placement="left" :width="480"
                            :disabled="scope.row.intro.length <= 30">
                    <p v-html="scope.row.intro?.replaceAll('\n','<br/>')"></p>
                    <template #reference>
                        <span slot="reference"
                              v-if="scope.row.intro && scope.row.intro.length <= 30">{{ scope.row.intro }}</span>
                        <span slot="reference" v-if="scope.row.intro && scope.row.intro.length > 30">{{
                            scope.row.intro.substring(0, 30)
                            + "..."
                            }}</span>
                    </template>
                </el-popover>
            </template>
        </el-table-column>
        <el-table-column prop="progress" label="已读%" width="100">
        </el-table-column>

        <el-table-column fixed="right" label="操作" width="200">
            <template #default="scope">
                <el-button icon="ChatDotRound" size="small" type="primary" link @click="gotoCommentList(scope.row)">书评
                </el-button>
                <el-button icon="view" size="small" type="primary" link @click="handleProd(scope.row.enid)">详情
                </el-button>
                <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">下载
                </el-button>

            </template>
        </el-table-column>
    </el-table>
    <Pagination :total="total" @pageChange="handleChangePage"></Pagination>
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
import {ElMessage, ElTable} from 'element-plus'
import {CourseCategory, CourseList, OpenDirectoryDialog, SetDir} from '../../wailsjs/go/backend/App'
import {services} from '../../wailsjs/go/models'
import Pagination from '../components/Pagination.vue'
import EbookInfo from '../components/EbookInfo.vue'
import DownloadDialog from "../components/DownloadDialog.vue";

import {useRouter} from 'vue-router'
import {userStore} from '../stores/user'
import {settingStore} from "../stores/setting";
import {Local} from '../utils/storage'

const store = userStore()
const setStore = settingStore()
const router = useRouter()
const loading = ref(true)
const page = ref(1)
const total = ref(0)
const pageSize = ref(15)
const searchInfo = ref({})
const dialogVisible = ref(false)
const prodEnid = ref("")

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
    await CourseList("ebook", "study", page.value, pageSize.value).then((table) => {
        loading.value = false
        Object.assign(tableData, table)
        console.log(table)
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


function handleChange() {
}

getTableData()

const dialogTitle = ref('detail')
const openDialog = () => {
    dialogVisible.value = true
}
const closeDialog = () => {
    //   initForm()
    dialogVisible.value = false
}


const openDownloadDialog = (row: any) => {
    // openDialogDir("Select download directory")
    downloadId.value = row.id
    downloadEnId.value = row.enid
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

const openDialogDir = async (title: string) => {
    await OpenDirectoryDialog(title).then((result) => {
        console.log(result)
    }).catch((error) => {
        console.log(error)
    })
}

const gotoCommentList = (row: any) => {
    router.push({
        path: `/ebook/comment`,
        query: {
            id: row.id,
            enid: row.enid,
            total: row.publish_num,
            title: row.title
        }
    })
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