<template>
  <el-table :data="tableData.list" v-loading="loading" height="97%" width="100%" :cell-style="{ textAlign: 'left' }"
    :header-cell-style="{textAlign: 'left'}" :row-style="{height: '50px'}" table-layout="auto" @sort-change="handleChange">
    <el-table-column prop="id" label="ID" width="100" />
    <el-table-column prop="title" label="标题" width="280" />
    <el-table-column prop="icon" label="封面" width="80">
      <template #default="scope">
        <el-image :src="scope.row.icon" :preview-src-list="[scope.row.icon]" :initial-index="0" style="width: 32px;" />
      </template>
    </el-table-column>
    <el-table-column prop="price" label="价格" width="100" />
    <el-table-column prop="intro" label="简介"  width="300">
      <template #default="scope">
        <el-popover title="简介" trigger="hover" placement="left" :width="480" :disabled="scope.row.intro.length <= 30">
          <p v-html="scope.row.intro?.replaceAll('\n','<br/>')"></p>
          <template #reference>
            <span slot="reference" v-if="scope.row.intro && scope.row.intro.length <= 30">{{scope.row.intro}}</span>
            <span slot="reference" v-if="scope.row.intro && scope.row.intro.length > 30">{{scope.row.intro.substr(0, 30)
            + "..."}}</span>
          </template>
        </el-popover>
      </template>
    </el-table-column>
    <el-table-column prop="progress" label="已读%" width="100" >
     
    </el-table-column>

    <el-table-column fixed="right" label="操作" width="200">
      <template #default="scope">
        <el-button icon="ChatDotRound" size="small" type="primary" link @click="gotoCommentList(scope.row)">书评
                </el-button>
        <el-button icon="view" size="small" type="primary" link @click="getEbookInfo(scope.row.enid)">详情</el-button>
        <el-button icon="download" size="small" type="primary" link @click="openDownloadDialog(scope.row)">下载
        </el-button>

      </template>
    </el-table-column>
  </el-table>
  <Pagination :total="total" @pageChange="handleChangePage"></Pagination>

  <el-dialog v-model="dialogVisible" width="60%" :before-close="closeDialog">
    <template #header="{titleId, titleClass }">
      <div class="my-header">
        <h4 :id="titleId" :class="titleClass">{{ebookInfo.operating_title}}</h4>
        <el-alert :title="ebookInfo.author_list.join(' ')+' 著/'+ebookInfo.press.name" type="info" :closable="false" center />
      </div>
    </template>
    <el-space wrap>
      <el-card v-for="i in 1" :key="i" class="box-card" width="100%">
        <template #header>
          <div class="card-header">
            <el-row :gutter="5" align="middle">
              <el-col :span="6">
                <el-row justify="center">
                  <el-image :src="ebookInfo.cover" fit="cover" />
                </el-row>
                <el-row justify="center">
                  <el-tag type="success">{{ebookInfo.book_author}}</el-tag>
                  <el-tag type="success" v-if="ebookInfo.classify_name">{{ebookInfo.classify_name}}</el-tag>
                </el-row>
              </el-col>
              <el-col :span="18" style="text-align:left">
                <p class="author-info" v-html="ebookInfo.author_info?.replaceAll('\n','<br/>')"></p><br />
                <el-tag class="ml-2" type="warning" v-if="ebookInfo.is_vip_book==1" round>
                  <el-icon><HotWater /></el-icon>会员免费</el-tag>
                <el-tag class="ml-2" type="danger" v-if="ebookInfo.is_tts_switch==true" round>
                  <el-icon><Microphone /></el-icon>可朗读
                </el-tag>
                <el-tag class="ml-2" type="danger" v-else round>
                  <el-icon><Mute /></el-icon>不可朗读
                </el-tag>
              </el-col>
            </el-row>
          </div>
        </template>
        <div class="book-tag" style="text-align:left">
          <el-tag  class="ml-2" type="info">出版日期：{{ebookInfo.publish_time}}</el-tag>
          <el-tag class="ml-2" type="success" v-if="ebookInfo.douban_score">豆瓣评分：{{ebookInfo.douban_score}}</el-tag>
          <el-tag class="ml-2" type="warning" v-if="ebookInfo.count">字数：{{Math.ceil(ebookInfo.count/1000)}}千字</el-tag>
          <el-tag class="ml-2" type="info" v-if="ebookInfo.read_time">阅读总时长：{{secondToHour(ebookInfo.read_time)}}</el-tag>
        </div>
        <el-divider content-position="left">{{ebookInfo.press.name}}</el-divider>
        <el-alert v-html="ebookInfo.press.brief?.replaceAll('\n','<br/>')" type="info" :closable="false" align="left" />

        <el-divider content-position="left">简 介</el-divider>
        <div class="book-intro" style="text-align:left">
          <p>{{ebookInfo.other_share_summary}}</p>
          <p v-html="ebookInfo.book_intro?.replaceAll('\n','<br/>')"></p>
        </div>
        <el-divider content-position="left">目 录</el-divider>
        <div class="catalog">
          <el-scrollbar height="320px" align="left">
            <span v-for="item in ebookInfo.catalog_list" :key="item.playOrder" class="scrollbar-catalog-item">
            {{repeat("&nbsp;&nbsp;", item.level)+ item.text }}<br/></span>
          </el-scrollbar>
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
        <el-button type="primary" @click="download(downloadId,downloadType)">
          确认
        </el-button>
      </span>
    </template>
  </el-dialog>

</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { ElTable, ElMessage } from 'element-plus'
import { CourseList, CourseCategory, EbookInfo, EbookDownload } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import Pagination from '../components/Pagination.vue'
import { repeat } from 'lodash'
import { secondToHour } from '../utils/utils'
import { useRouter } from 'vue-router'
import { userStore } from '../stores/user';


const store = userStore()
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
const downloadTypeOptions = [
  { value: 1, label: "HTML" }, { value: 2, label: "PDF" }, { value: 3, label: "EPUB" }
]

let tableData = reactive(new services.CourseList)
let ebookInfo = reactive(new services.EbookDetail)

const fontWeight = reactive({
  'font-weight': 'bold', 
})
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
    loading.value=false
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

const getEbookInfo = async (enid: string) => {
  await EbookInfo(enid).then((info) => {
    console.log(info)
    Object.assign(ebookInfo, info)
    openDialog()
  }).catch((error) => {
    ElMessage({
      message: error,
      type: 'warning'
    })
  })
  return
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
  downloadId.value = row.id
  dialogDownloadVisible.value = true
}
const closeDownloadDialog = () => {
  //   initForm()
  downloadType.value = 1
  dialogDownloadVisible.value = false
}

const download = async (id: number, dType: number) => {
  await EbookDownload(id, dType).then((info) => {
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

const gotoCommentList = (row: any) => {
    router.push({
        path: `/ebook/comment`,
        query: {
            id: row.id,
            enid:row.enid,
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