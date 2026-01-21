<template>
    <div class="course-container">
        <div v-if="filterOptions.length > 0 && !groupMode.active" class="filter-container">
            <el-radio-group v-model="currentFilter" @change="handleFilterChange" size="small">
                <el-radio-button 
                    v-for="item in filterOptions" 
                    :key="item.filter" 
                    :label="item.filter"
                >
                    {{ item.name }}
                    <span v-if="item.show_count">({{ item.count }})</span>
                </el-radio-button>
            </el-radio-group>
        </div>
        <div v-if="groupMode.active" class="group-header">
            <el-button type="primary" link @click="exitGroup">
                <el-icon class="el-icon--left"><ArrowLeft /></el-icon>返回
            </el-button>
            <span class="group-title">{{ groupMode.title }}</span>
        </div>

        <div v-loading="initLoading" class="course-grid-container" v-infinite-scroll="loadMore" :infinite-scroll-disabled="disabled" :infinite-scroll-immediate="false">
            <div v-if="tableData.list && tableData.list.length > 0" class="course-grid">
                <div v-for="item in tableData.list" :key="item.id" class="course-card" @click="item.is_group ? enterGroup(item) : gotoArticleList(item)">
                    <div class="card-cover">
                        <!-- 分组封面拼图 -->
                        <div v-if="item.is_group && item.group_books && item.group_books.length > 0" class="group-cover-grid">
                            <div v-for="(book, index) in item.group_books.slice(0, 4)" :key="book.id || index" class="group-grid-item">
                                <el-image :src="book.icon" fit="cover" loading="lazy" class="grid-image">
                                    <template #error>
                                        <div class="grid-placeholder">
                                            <el-icon><Picture /></el-icon>
                                        </div>
                                    </template>
                                </el-image>
                            </div>
                            <div v-for="n in (4 - Math.min(item.group_books.length, 4))" :key="'ph-'+n" class="group-grid-item">
                                <div class="grid-placeholder bg-gray">
                                    <el-icon><Picture /></el-icon>
                                </div>
                            </div>
                        </div>

                        <el-image 
                            v-else-if="item.icon" 
                            :src="item.icon" 
                            fit="cover"
                            loading="lazy"
                        >
                            <template #placeholder>
                                <div class="image-placeholder">
                                    <el-icon><Picture /></el-icon>
                                </div>
                            </template>
                        </el-image>
                        <div v-else class="no-cover">
                            <el-icon v-if="item.is_group" :size="40"><Folder /></el-icon>
                            <span v-else>无封面</span>
                        </div>
                        
                        <!-- 悬停遮罩层 (仅非分组显示操作) -->
                        <div v-if="!item.is_group" class="card-overlay" @click.stop>
                            <div class="overlay-actions">
                                <el-button circle type="primary" :icon="View" @click="handleProd(item.enid)" title="详情" />
                                <el-button circle type="success" :icon="Download" @click="openDownloadDialog(item)" title="下载" />
                            </div>
                        </div>
                        
                        <!-- 分组标识 -->
                        <div v-if="item.is_group" class="group-badge">
                            <el-icon><Folder /></el-icon>
                            <span>分组</span>
                        </div>
                    </div>
                    
                    <div class="card-content">
                        <h3 class="card-title" :title="item.title">{{ item.title }}</h3>
                        <div class="card-meta">
                            <span v-if="item.is_group" class="meta-info">{{ item.course_num || 0 }} 本课程</span>
                            <span v-else class="meta-info">已更 {{ item.publish_num || 0 }}/{{ item.course_num || 0 }}</span>
                            
                            <span v-if="!item.is_group" class="meta-progress">{{ item.progress || 0 }}%</span>
                        </div>
                        <div v-if="!item.is_group" class="progress-bar">
                             <div class="progress-fill" :style="{ width: (item.progress || 0) + '%' }"></div>
                        </div>
                    </div>
                </div>
            </div>
            <el-empty v-else description="暂无课程" />
        </div>
    
    </div>

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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, View, Download, Picture, Folder } from '@element-plus/icons-vue'
import {CourseList, CourseCategory, CourseGroupList, SetDir, GetNavbar} from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { userStore } from '../stores/user';
import { settingStore } from '../stores/setting';
import { useAppRouter } from '../composables/useRouter';
import CourseInfo from '../components/CourseInfo.vue'
import DownloadDialog from "../components/DownloadDialog.vue";
import { Local } from '../utils/storage';

const store = userStore()
const setStore = settingStore()
const { pushLogin, pushCourseDetail, pushSetting } = useAppRouter()

const loading = ref(false)
const initLoading = ref(true)
const page = ref(1)
const total = ref(0)
const outerTotal = ref(0)
const pageSize = ref(20) // Increase page size for better scrolling experience
const lastPageSize = ref(20)
const searchInfo = ref({})
const currentFilter = ref('all')
const filterOptions = ref<any[]>([])

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
    GetNavbar().then((res: any) => {
        if (res && res.list) {
            const opts: any[] = []
            res.list.forEach((item: any) => {
                if (item.category === "bauhinia" && item.children) {
                    opts.push(...item.children)
                }
            })
            filterOptions.value = opts
        }
    })

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

const noMore = computed(() => {
    const currentCount = tableData.list ? tableData.list.length : 0
    if (groupMode.active || currentFilter.value !== 'all') {
        const total = tableData.total || 0
        if (total > 0) {
            return currentCount >= total
        }
        if (tableData.is_more === 0) {
            return true
        }
        return lastPageSize.value < pageSize.value
    }
    if (outerTotal.value > 0) {
        return currentCount >= outerTotal.value
    }
    return lastPageSize.value < pageSize.value
})

const disabled = computed(() => loading.value || noMore.value)

const loadMore = () => {
    if (disabled.value) return
    page.value += 1
    getTableData(true)
}

const handleFilterChange = () => {
    console.log('Course filter changed:', currentFilter.value)
    if (currentFilter.value === "all") {
        groupMode.active = false
        groupMode.groupId = 0
        groupMode.title = ''
    }
    page.value = 1
    tableData.list = [] // Clear list to avoid confusion
    getTableData()
}

const getTableData = async (append = false) => {
    console.log('Course getTableData:', { append, currentFilter: currentFilter.value, groupMode: groupMode.active })
    loading.value = true
    if (!append) initLoading.value = true
    
    let fetcher;
    if (groupMode.active) {
        fetcher = CourseGroupList("bauhinia", "study", currentFilter.value,groupMode.groupId, page.value, pageSize.value)
    } else {
        fetcher = CourseList("bauhinia", "study", currentFilter.value, page.value, pageSize.value)
    }

    await fetcher.then((table) => {
        loading.value = false
        initLoading.value = false
        
        const fetchedList = table.list || []
        lastPageSize.value = fetchedList.length
        
        if (append) {
            if (fetchedList.length > 0) {
                tableData.list.push(...fetchedList)
            }
        } else {
            Object.assign(tableData, table)
        }
        
        if (groupMode.active || currentFilter.value !== 'all') {
            total.value = table.total || 0
        } else {
            total.value = outerTotal.value
        }
    }).catch((error) => {
        loading.value = false
        initLoading.value = false
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
    const groupId = Number(row?.group_id || row?.id || 0)

    if (!groupId) return

    groupMode.active = true
    groupMode.groupId = groupId
    groupMode.title = String(row?.title || '')
    page.value = 1
    tableData.list = [] 
    getTableData()
}

const exitGroup = () => {
    groupMode.active = false
    groupMode.groupId = 0
    groupMode.title = ''
    page.value = 1
    total.value = outerTotal.value
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
.course-container {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 20px;
    box-sizing: border-box;
}

.group-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    padding-top: 10px;
}

.group-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
}

.filter-container {
    margin-bottom: 20px;
    overflow-x: auto;
    white-space: nowrap;
    padding-bottom: 4px;
}
.filter-container::-webkit-scrollbar {
    height: 4px;
}
.filter-container::-webkit-scrollbar-thumb {
    background: var(--border-color);
    border-radius: 2px;
}

.course-grid-container {
    flex: 1;
    overflow-y: auto;
    padding-bottom: 20px;
    /* 隐藏滚动条但保留功能 - 清新风格 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE 10+ */
}

.course-grid-container::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}

.course-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 20px;
    padding: 4px; /* 防止阴影被切 */
}

.course-card {
    background: var(--card-bg, #fff);
    border-radius: 12px;
    box-shadow: var(--shadow-soft, 0 2px 12px rgba(0, 0, 0, 0.08));
    overflow: hidden;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    position: relative;
    border: 1px solid var(--border-soft, #ebeef5);
    display: flex;
    flex-direction: column;
}

.course-card:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-medium, 0 8px 24px rgba(0, 0, 0, 0.12));
    border-color: transparent;
}

.card-cover {
    position: relative;
    width: 100%;
    aspect-ratio: 1;
    background-color: var(--fill-color-light, #f5f7fa);
    overflow: hidden;
}

.card-cover .el-image {
    width: 100%;
    height: 100%;
    display: block;
    transition: transform 0.5s ease;
}

.course-card:hover .card-cover .el-image {
    transform: scale(1.05);
}

.image-placeholder, .no-cover {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--text-tertiary);
    gap: 8px;
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

.course-card:hover .card-overlay {
    opacity: 1;
}

.overlay-actions {
    display: flex;
    gap: 16px;
    transform: translateY(10px);
    transition: transform 0.3s ease;
}

.course-card:hover .overlay-actions {
    transform: translateY(0);
}

.group-badge {
    position: absolute;
    top: 8px;
    right: 8px;
    background: rgba(0, 0, 0, 0.6);
    color: #fff;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 12px;
    display: flex;
    align-items: center;
    gap: 4px;
    backdrop-filter: blur(4px);
}

.group-cover-grid {
    width: 100%;
    height: 100%;
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr 1fr;
    gap: 2px;
    background: var(--fill-color-light, #f5f7fa);
}

.group-grid-item {
    position: relative;
    overflow: hidden;
    background: var(--fill-color-light, #f5f7fa);
    width: 100%;
    height: 100%;
}

.grid-image {
    width: 100%;
    height: 100%;
    display: block;
}

.grid-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-tertiary, #909399);
    background: var(--fill-color-light, #f5f7fa);
}

.bg-gray {
    background: var(--fill-color, #f0f2f5);
}

.card-content {
    padding: 12px 16px;
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

.card-title {
    margin: 0 0 8px;
    font-size: 15px;
    line-height: 1.4;
    color: var(--text-primary);
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    font-weight: 600;
}

.card-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 12px;
    color: var(--text-secondary);
    margin-bottom: 8px;
}

.meta-progress {
    color: var(--accent-color, #ff6b00);
    font-weight: 500;
}

.progress-bar {
    height: 4px;
    background: var(--fill-color, #f0f2f5);
    border-radius: 2px;
    overflow: hidden;
}

.progress-fill {
    height: 100%;
    background: var(--accent-color, #ff6b00);
    border-radius: 2px;
    transition: width 0.3s ease;
}
</style>
