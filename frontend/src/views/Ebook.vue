<template>
  <div class="ebook-container">
    <div class="header-actions" v-if="groupMode.active">
      <el-button type="primary" link @click="exitGroup">
        <el-icon><ArrowLeft /></el-icon> 返回
      </el-button>
      <span class="group-title">{{ groupMode.title }}</span>
    </div>

    <div class="ebook-grid-container" v-loading="initLoading" v-infinite-scroll="loadMore" :infinite-scroll-disabled="disabled" :infinite-scroll-immediate="false">
      <div class="ebook-grid">
        <div v-for="item in tableData.list" :key="item.enid" class="ebook-card" @click="handleCardClick(item)">
          <div class="card-cover">
            <el-image 
              v-if="item.icon" 
              :src="item.icon" 
              fit="cover" 
              loading="lazy"
              class="cover-image"
            >
              <template #placeholder>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div v-else class="no-cover">
              <el-icon><Notebook /></el-icon>
            </div>
            
            <!-- Group Indicator -->
            <div v-if="item.is_group" class="group-badge">
              <el-icon><Folder /></el-icon>
              <span>{{ item.course_num || 0 }}本</span>
            </div>

            <!-- Progress Badge -->
            <div v-if="!item.is_group && item.progress" class="progress-badge">
              {{ item.progress }}%
            </div>

            <!-- Hover Overlay -->
            <div class="card-overlay">
              <div class="overlay-actions">
                <template v-if="!item.is_group">
                  <el-tooltip content="详情" :show-after="500">
                    <el-button circle size="small" @click.stop="handleProd(item.enid)">
                      <el-icon><View /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="下载" :show-after="500">
                    <el-button circle size="small" type="primary" @click.stop="openDownloadDialog(item)">
                      <el-icon><DownloadIcon /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="书评" :show-after="500">
                    <el-button circle size="small" @click.stop="gotoCommentList(item)">
                      <el-icon><ChatDotRound /></el-icon>
                    </el-button>
                  </el-tooltip>
                  <el-tooltip content="移出书架" :show-after="500">
                    <el-button circle size="small" type="danger" @click.stop="ebookShelfRemove(item.enid)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </el-tooltip>
                </template>
                <template v-else>
                   <el-button type="primary" round size="small" @click.stop="enterGroup(item)">进入分组</el-button>
                </template>
              </div>
            </div>
          </div>

          <div class="card-info">
            <h3 class="card-title" :title="item.title">{{ item.title }}</h3>
            <p class="card-intro" :title="item.intro">{{ stripHtml(item.intro) }}</p>
            <div class="card-meta">
               <span v-if="item.price" class="price">¥{{ item.price }}</span>
               <el-tag v-if="item.is_group" size="small" effect="plain">分组</el-tag>
            </div>
          </div>
        </div>
      </div>
    </div>

    
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
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref, computed} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {
  ChatDotRound,
  View,
  Download as DownloadIcon,
  Delete,
  Picture,
  Notebook,
  Folder,
  ArrowLeft
} from '@element-plus/icons-vue'
import {
  CourseCategory,
  CourseList,
  CourseGroupList,
  EbookShelfRemove,
  SetDir
} from '../../wailsjs/go/backend/App'
import {services} from '../../wailsjs/go/models'
import EbookInfo from '../components/EbookInfo.vue'
import DownloadDialog from "../components/DownloadDialog.vue";

import {useAppRouter} from '../composables/useRouter'
import {userStore} from '../stores/user'
import {settingStore} from "../stores/setting";
import {Local} from '../utils/storage'

const store = userStore()
const setStore = settingStore()
const { pushEbookComment, pushSetting, pushLogin } = useAppRouter()
const loading = ref(false)
const initLoading = ref(true)
const page = ref(1)
const total = ref(0)
const outerTotal = ref(0)
const pageSize = ref(20) // Increased for grid view
const lastPageSize = ref(20)
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
    
    // Initial load
    getTableData()
})

const noMore = computed(() => {
    const currentCount = tableData.list ? tableData.list.length : 0
    if (groupMode.active) {
        return currentCount >= (tableData.total || 0)
    }
    if (outerTotal.value > 0) {
        return currentCount >= outerTotal.value
    }
    // Fallback if outerTotal is not yet available or reliable
    return lastPageSize.value < pageSize.value
})

const disabled = computed(() => loading.value || noMore.value)

const loadMore = () => {
    if (disabled.value) return
    page.value += 1
    getTableData(true)
}

const getTableData = async (append = false) => {
    loading.value = true
    if (!append) initLoading.value = true
    const fetcher = groupMode.active
        ? CourseGroupList("ebook", "study", groupMode.groupId, page.value, pageSize.value)
        : CourseList("ebook", "study", page.value, pageSize.value)
    await fetcher.then((table) => {
        loading.value = false
        initLoading.value = false
        
        // Update lastPageSize based on the fetched list length
        const fetchedList = table.list || []
        lastPageSize.value = fetchedList.length
        
        if (append) {
            if (fetchedList.length > 0) {
                tableData.list.push(...fetchedList)
            }
        } else {
            Object.assign(tableData, table)
        }
        total.value = groupMode.active ? (table.total || 0) : outerTotal.value
    }).catch((error) => {
        loading.value = false
        initLoading.value = false
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
  getTableData()
}

const exitGroup = () => {
    groupMode.active = false
    groupMode.groupId = 0
    groupMode.title = ''
    page.value = 1
    getTableData()
}

const handleCardClick = (item: any) => {
    if (item.is_group) {
        enterGroup(item)
    } else {
        handleProd(item.enid)
    }
}

const openDownloadDialog = (row: any) => {
    downloadId.value = row.id
    downloadEnId.value = row.enid
    
    if (setStore.getDownloadDir == "") {
        ElMessage({
            message: '请设置文件保存目录',
            type: 'warning'
        })
        pushSetting()
    } else {
        // Ensure directories are set in backend
        SetDir([setStore.getDownloadDir,
            setStore.getFfmpegDirDir,
            setStore.getWkDir]).then(() => {
            dialogDownloadVisible.value = true
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
}

const gotoCommentList = (row: any) => {
    pushEbookComment({
        id: row.id,
        enid: row.enid,
        total: row.publish_num,
        title: row.title
    })
}

const ebookShelfRemove = (enid: string) => {
    ElMessageBox.confirm(
        '确定要将该电子书移出书架吗?',
        '提示',
        {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        // backend expects array
        EbookShelfRemove([enid]).then(() => {
            ElMessage({
                type: 'success',
                message: '移除成功',
            })
            getTableData()
        }).catch((err) => {
            ElMessage({
                type: 'error',
                message: err,
            })
        })
    })
}

const closeDialog = () => {
    dialogVisible.value = false
}

const stripHtml = (html: string) => {
    if (!html) return ''
    const tmp = document.createElement("DIV")
    tmp.innerHTML = html
    return tmp.textContent || tmp.innerText || ""
}
</script>

<style scoped>
.ebook-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-color);
  padding: 20px;
  overflow: hidden;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.group-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.ebook-grid-container {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 20px;
  /* 隐藏滚动条但保留功能 */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.ebook-grid-container::-webkit-scrollbar {
  display: none;
}

.ebook-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 20px;
  padding: 4px;
  align-content: start;
}

.ebook-card {
  background-color: var(--card-bg);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  cursor: pointer;
  position: relative;
}

.ebook-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-medium);
  border-color: var(--primary-color-light);
}

.card-cover {
  position: relative;
  aspect-ratio: 2/3;
  overflow: hidden;
  background-color: var(--fill-color);
}

.cover-image {
  width: 100%;
  height: 100%;
  transition: transform 0.5s ease;
}

.ebook-card:hover .cover-image {
  transform: scale(1.05);
}

.no-cover {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  color: var(--text-secondary);
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--fill-color);
  color: var(--text-secondary);
}

.group-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  background-color: rgba(0, 0, 0, 0.6);
  color: #fff;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  backdrop-filter: blur(4px);
}

.progress-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  background-color: var(--primary-color);
  color: #fff;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.card-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  backdrop-filter: blur(2px);
}

.ebook-card:hover .card-overlay {
  opacity: 1;
}

.overlay-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
  padding: 0 12px;
}

.card-info {
  padding: 12px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-intro {
  font-size: 12px;
  color: var(--text-secondary);
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
}

.price {
  font-size: 13px;
  color: var(--accent-color);
  font-weight: 600;
  font-family: var(--font-family-mono);
}

/* Dark Mode Adaptation */
.theme-dark .ebook-card {
  border-color: var(--border-soft);
}

.theme-dark .group-badge {
  background-color: rgba(255, 255, 255, 0.2);
}

/* Scrollbar Styling */
.ebook-grid {
  /* 隐藏滚动条但保留功能 - 清新风格 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

.ebook-grid::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}
</style>