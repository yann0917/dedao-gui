<template>
  <div class="compass-container">
    <div class="compass-grid" v-loading="loading">
      <div v-for="item in tableData.list" :key="item.id" class="compass-card">
        <div class="card-inner">
          <div class="card-cover">
            <el-image :src="item.icon" loading="lazy" fit="cover">
              <template #placeholder>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div class="card-overlay">
              <el-button type="primary" circle @click="openDownloadDialog(item)">
                <el-icon><DownloadIcon /></el-icon>
              </el-button>
            </div>
          </div>
          <div class="card-info">
            <h3 class="card-title" :title="item.title">{{ item.title }}</h3>
            <div class="card-meta">
              <span class="replier" v-if="item.ext_info && item.ext_info[0]">
                <el-icon><User /></el-icon>
                {{ item.ext_info[0].replier_name }}
              </span>
            </div>
            <div class="card-intro" v-if="item.intro">
              <p>{{ item.intro }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="pagination-container">
      <Pagination :total="total" @pageChange="handleChangePage"></Pagination>
    </div>

    <!-- Download Dialog -->
    <el-dialog v-model="dialogDownloadVisible" title="下载内容" width="30%" align-center class="custom-dialog">
      <div class="dialog-content">
        <el-form label-position="top">
          <el-form-item label="选择格式">
            <el-select v-model="downloadType" placeholder="请选择下载格式" style="width: 100%">
              <el-option
                v-for="item in downloadTypeOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeDownloadDialog">取消</el-button>
          <el-button type="primary" @click="download(downloadId, downloadType)">
            开始下载
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { CourseList, CourseCategory, CourseDownload } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { useRouter } from 'vue-router'
import { userStore } from '../stores/user';
import Pagination from '../components/Pagination.vue'
import { Local } from '../utils/storage';
import { Download as DownloadIcon, User, Picture } from '@element-plus/icons-vue'

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
let tableData = reactive(new services.CourseList)

onMounted(() => {
    CourseCategory().then(result => {
        result.forEach((item, key) => {
            if (item.category == "compass") {
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
    await CourseList("compass", "study", page.value, pageSize.value).then((table) => {
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

getTableData()


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
    await CourseDownload(id, 0, dType,'').then((info) => {
        console.log(info)
        ElMessage({
            message: '任务已添加到下载队列',
            type: 'success'
        })
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
.compass-container {
  height: calc(100vh - 20px);
  display: flex;
  flex-direction: column;
  background-color: var(--fill-color-light);
  padding: 20px;
  box-sizing: border-box;
}

.compass-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
  overflow-y: auto;
  padding: 4px;
  align-content: start;
  
  /* 隐藏滚动条但保留功能 - 清新风格 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

.compass-grid::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

.compass-card {
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: var(--shadow-soft);
  transition: all 0.3s ease;
  overflow: hidden;
  height: 100%;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border-soft);
}

.compass-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-medium);
  border-color: var(--primary-color-light);
}

.card-inner {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.card-cover {
  position: relative;
  height: 160px;
  background: var(--fill-color);
  overflow: hidden;
}

.card-cover .el-image {
  width: 100%;
  height: 100%;
  transition: transform 0.5s ease;
}

.compass-card:hover .card-cover .el-image {
  transform: scale(1.05);
}

.image-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: var(--text-secondary);
  font-size: 24px;
}

.card-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  backdrop-filter: blur(2px);
}

.compass-card:hover .card-overlay {
  opacity: 1;
}

.card-info {
  padding: 16px;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-secondary);
}

.replier {
  display: flex;
  align-items: center;
  gap: 4px;
}

.card-intro {
  margin-top: auto;
  font-size: 13px;
  color: var(--text-tertiary);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  background: var(--card-bg);
  padding: 12px;
  border-radius: 8px;
  box-shadow: var(--shadow-soft);
}

/* Dark mode adjustments */
.theme-dark .compass-card {
  background: var(--card-bg);
  border-color: var(--border-soft);
}

.theme-dark .card-title {
  color: var(--text-primary);
}

.theme-dark .card-meta {
  color: var(--text-secondary);
}

.theme-dark .card-intro {
  color: var(--text-tertiary);
}

/* Custom Dialog */
.custom-dialog {
  border-radius: 12px;
  overflow: hidden;
}

.dialog-content {
  padding: 10px 0;
}
</style>