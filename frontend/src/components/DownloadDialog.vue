<template>
    <el-dialog v-model="dialogVisible" title="请选择下载格式" align-center center width="30%" :before-close="handleCloseDialog">
        <el-form >
            <el-form-item label="下载格式" label-width="80px">
                <el-select v-model="downloadType" placeholder="请选择下载格式" :disabled="isDownloading">
                    <el-option v-for="item in props.downloadTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
            <div v-if="isDownloading" class="download-status">
                <el-progress
                    :percentage="percentage"
                    :status="percentage === 100 ? 'success' : 'warning'"
                    :stroke-width="20"
                    :text-inside="true"
                >
                    <span>{{content}}</span>
                </el-progress>
                <div class="status-text">
                    <p>当前下载: <strong>{{ currentFile }}</strong></p>
                    <p>总进度: {{ currentCount }}/{{ totalCount }} ({{ percentage }}%)</p>
                </div>
            </div>
        </el-form>

        <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCloseDialog">{{ isDownloading ? '后台继续下载' : '取消' }}</el-button>
        <el-button type="primary" @click="download()" :disabled="isDownloading">
          确认下载
        </el-button>
      </span>
        </template>
    </el-dialog>

</template>

<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {EbookDownload, CourseDownload} from "../../wailsjs/go/backend/App";
import {ElMessage, ElNotification} from "element-plus";
import { EventsOn, EventsOff} from "../../wailsjs/runtime/runtime";

let percentage=ref(0)
let content=ref('')
let isDownloading=ref(false)
let currentFile=ref('')
let currentCount=ref(0)
let totalCount=ref(0)
let downloadCompleted=ref(false)

const dialogVisible = ref(false)
const downloadType = ref(1)
// const downloadId = ref(0)
// const prodType = ref(2) //2-电子书,66-课程
// const downloadTypeOptions = [
//     { value: 1, label: "HTML" }, { value: 2, label: "PDF" }, { value: 3, label: "EPUB" }
// ]

const props = defineProps({
    downloadId:{
        type:Number,
        required:true,
        default:0,
    },
    enId:{
        type:String,
        default:'',
    },
    prodType:{
        type:Number,
        required:true,
        default:0,
    },
    articleId:{
        type:Number,
        default:0,
    },
    dialogVisible: {
        type: Boolean,
        default: false,
    },
    downloadTypeOptions:{
        type: Object,
        required:true,
        default:{}
    }
});
const emits = defineEmits(["close"]);

onMounted(() => {
    openDialog();
});

const openDialog = () => {
    // openDialogDir("Select download directory")
    // downloadId.value = row.id
    dialogVisible.value = props.dialogVisible
}

const handleCloseDialog = () => {
    if (isDownloading.value) {
        ElMessage({
            type: 'info',
            message: '下载将在后台继续进行，完成后会通知您'
        })
        dialogVisible.value = false
        emits("close")
    } else {
        closeDialog()
    }
}

const closeDialog = () => {
    //   initForm()
    // downloadType.value = 1
    // dialogVisible.value = false
    EventsOff("courseDownload")
    EventsOff("ebookDownload")
    percentage.value = 0
    content.value = ''
    isDownloading.value = false
    currentFile.value = ''
    currentCount.value = 0
    totalCount.value = 0
    downloadCompleted.value = false
    emits("close")
}

const showDownloadCompleteNotification = () => {
    ElNotification({
        title: '下载完成',
        message: '所有文件已下载完成！',
        type: 'success',
        duration: 5000,
        position: 'bottom-right'
    })
}

const handleDownloadProgress = (data) => {
    if (data) {
        console.log(data)
        percentage.value = data.pct
        content.value = data.value ? data.value + ' 下载中...' : '正在下载...'
        currentFile.value = data.value || '未知文件'
        currentCount.value = data.current || 0
        totalCount.value = data.total || 0
        
        // 检测下载是否完成
        if (data.pct === 100 && currentCount.value === totalCount.value) {
            downloadCompleted.value = true
            content.value = '下载完成！'
            showDownloadCompleteNotification()
            // 延迟关闭对话框，让用户看到下载完成状态
            setTimeout(() => {
                closeDialog()
            }, 3000)
        }
    }
}

const download = async () => {
    isDownloading.value = true
    content.value = '准备下载中...'
    
    switch (props.prodType) {
        case 2:
            EventsOn("ebookDownload", handleDownloadProgress)
            await EbookDownload(props.downloadId, downloadType.value, props.enId).then((info) => {
                // console.log(info)
                if (!downloadCompleted.value) {
                    downloadCompleted.value = true
                    showDownloadCompleteNotification()
                }
            }).catch((error) => {
                isDownloading.value = false
                ElMessage({
                    message: error,
                    type: 'error'
                })
            })
            break;
        case 66:
            EventsOn("courseDownload", handleDownloadProgress)
            await CourseDownload(props.downloadId, props.articleId, downloadType.value, props.enId).then((info) => {
                console.log(info)
                if (!downloadCompleted.value) {
                    downloadCompleted.value = true
                    showDownloadCompleteNotification()
                }
            }).catch((error) => {
                isDownloading.value = false
                ElMessage({
                    message: error,
                    type: 'error'
                })
            })
            break;
    }
    
    return
}

</script>

<style scoped>
.download-status {
    margin-top: 20px;
    padding: 10px;
    border-radius: 4px;
    background-color: #f5f7fa;
}

.status-text {
    margin-top: 10px;
    font-size: 14px;
    color: #606266;
}

.status-text p {
    margin: 5px 0;
}
</style>