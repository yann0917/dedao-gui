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
                <el-button icon="view" size="small" type="primary" link @click="getCourseInfo(scope.row.enid)">详情
                </el-button>
                <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">下载
                </el-button>

            </template>
        </el-table-column>
    </el-table>
    <Pagination :total="total" @pageChange="handleChangePage"></Pagination>

    <el-dialog v-model="dialogVisible" width="60%" :before-close="closeDialog">
        <template #header="{ titleId, titleClass }">
            <div class="my-header">
                <h4 :id="titleId" :class="titleClass">{{ courseInfo.class_info.name }}</h4>
                <el-alert :title="courseInfo.class_info.intro" type="info" :closable="false" center />
            </div>
        </template>
        <el-space wrap>
            <el-card v-for="i in 1" :key="i" class="box-card" width="100%">
                <template #header>
                    <div class="card-header">
                        <el-row :gutter="5" align="middle">
                            <el-col :span="4">
                                <el-row justify="center">
                                    <el-avatar :size="84" :src="courseInfo.class_info.lecturer_avatar" fit="fill" />
                                </el-row>
                                <el-row justify="center">
                                    <el-tag type="success">{{ courseInfo.class_info.lecturer_name }}</el-tag>
                                </el-row>
                            </el-col>
                            <el-col :span="18">
                                <p style="text-align:left"
                                    v-html="courseInfo.class_info?.lecturer_intro.replaceAll('\n', '<br/>')"></p>
                            </el-col>
                        </el-row>
                    </div>
                </template>

                <div class="book-tag" style="text-align:left">
                    <el-tag v-if="courseInfo.class_info.is_finished">
                        共{{ courseInfo.class_info.phase_num }}{{ courseInfo.class_info.price_desc }}
                    </el-tag>
                    <el-tag v-else>
                        共{{ courseInfo.class_info.phase_num }}{{ courseInfo.class_info.price_desc }}·已更新{{
                                courseInfo.class_info.current_article_count
                        }}{{ courseInfo.class_info.price_desc }}
                    </el-tag>
                    <el-tag class="ml-2" type="success">{{ courseInfo.class_info.learn_user_count }}人加入学习</el-tag>
                    <el-tag class="ml-2" type="info">{{ courseInfo.class_info.collection.collection_count }}人收藏</el-tag>
                    <el-tag class="ml-2" type="warning" v-if="courseInfo.class_info.collection.is_collected">已收藏
                    </el-tag>
                    <el-tag class="ml-2" type="warning" v-else>未收藏</el-tag>
                    <el-tag class="ml-2" type="danger">{{ courseInfo.class_comment_info?.count }}条评价</el-tag>
                </div>

                <el-divider content-position="left">
                    <el-rate v-model="averageScore" disabled show-score :allow-half="true" text-color="#ff9900" />
                </el-divider>
                <div class="comment-info flex items-left" style="text-align: left;">
                    <el-carousel height="200px" direction="vertical" loop>
                        <el-carousel-item v-for="item in courseInfo.class_comment_info?.comment_list" :key="item.id">
                            <p class="text-large font-600 mr-3" v-if="item.title"> {{ item.title }} </p>
                            <p class="text-sm mr-2" style="color: var(--el-text-color-regular)"
                                v-if="item.no_style_content">
                                <el-popover trigger="hover" placement="right" :width="300"
                                    :disabled="item.no_style_content.length <= 240" :content="item.no_style_content">
                                    <template #reference>
                                        <span slot="reference" v-if="item.no_style_content.length <= 240">{{
                                                item.no_style_content
                                        }}</span>
                                        <span slot="reference" v-if="item.no_style_content.length > 240">{{
                                                item.no_style_content.slice(0, 240)
                                                + "..."
                                        }}</span>
                                    </template>
                                </el-popover>
                            </p>
                            <el-avatar class="mr-3" :size="20" :src="item.avatar_s" /><span> &emsp;{{ item.nickname }}
                                评价了
                                <el-rate v-model.number="item.score" disabled show-score :allow-half="true"
                                    text-color="#ff9900" />
                            </span>
                        </el-carousel-item>
                    </el-carousel>
                </div>

                <div class="flex items-center">

                </div>
                <el-divider content-position="left">课程亮点</el-divider>
                <div style="text-align:left">
                    <span v-html="courseInfo.class_info.highlight?.replaceAll('\n', '<br/>')"></span>
                </div>

                <el-divider content-position="left">课程表</el-divider>
                <div class="outline_img" style="text-align:left">
                    <el-image :src="courseInfo.class_info.outline_img"></el-image>
                </div>
            </el-card>
        </el-space>

    </el-dialog>

    <el-dialog v-model="dialogDownloadVisible" title="请选择下载格式" align-center center width="30%">
        <el-select v-model="downloadType" placeholder="请选择下载格式">
            <el-option v-for="item in downloadTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
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
import { ref, reactive, onMounted } from 'vue'
import { ElTable, ElMessage } from 'element-plus'
import { CourseList, CourseCategory, CourseInfo, CourseDownload } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { useRouter } from 'vue-router'
import { userStore } from '../stores/user';
import Pagination from '../components/Pagination.vue'

const store = userStore()
const router = useRouter()

const loading = ref(true)
const page = ref(1)
const total = ref(0)
const pageSize = ref(15)
const dialogVisible = ref(false)



const dialogDownloadVisible = ref(false)
const downloadType = ref(1)
const downloadId = ref(0)

const downloadTypeOptions = [
    { value: 1, label: "MP3" }, { value: 2, label: "PDF" }, { value: 3, label: "Markdown" }
]
let averageScore = ref(0)
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

const getCourseInfo = async (enid: string) => {
    await CourseInfo(enid).then((info) => {
        console.log(info)
        Object.assign(courseInfo, info)
        averageScore.value = Number(courseInfo.class_comment_info.average_score)
        openDialog()
    }).catch((error) => {
        loading.value = false
        ElMessage({
            message: error,
            type: 'warning'
        })
    })
    return
}

getTableData()

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

const download = async (id: number, dType: number) => {
    await CourseDownload(id, 0, dType).then((info) => {
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