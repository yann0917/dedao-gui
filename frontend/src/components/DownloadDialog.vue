<template>
    <el-dialog 
        v-model="dialogVisible" 
        title="下载选项" 
        align-center 
        center 
        width="420px" 
        :before-close="closeDialog"
        class="custom-download-dialog"
    >
        <div class="download-container">
            <div class="format-selector">
                <div class="section-label">选择导出格式</div>
                <div class="format-options">
                    <div 
                        v-for="item in props.downloadTypeOptions" 
                        :key="item.value" 
                        class="format-option"
                        :class="downloadType === item.value ? 'active' : ''"
                        @click="downloadType = item.value"
                    >
                        <span class="format-text">{{ item.label }}</span>
                        <el-icon v-if="downloadType === item.value" class="selected-icon"><Check /></el-icon>
                    </div>
                </div>
            </div>

            <div class="download-status">
                <div class="status-header">
                    <span class="status-text">{{ content }}</span>
                </div>
            </div>
        </div>

        <template #footer>
            <div class="dialog-footer">
                <el-button @click="closeDialog" :disabled="isDownloading">取消</el-button>
                <el-button type="primary" @click="download()" :loading="isDownloading">
                    {{ isDownloading ? '提交中...' : '开始下载' }}
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script lang="ts" setup>
import {onMounted, ref, PropType} from "vue";
import {EbookDownload, CourseDownload, OdobDownload} from "../../wailsjs/go/backend/App";
import {ElMessage} from "element-plus";
import { Check } from '@element-plus/icons-vue'

let content=ref('任务将提交到下载队列，可在“下载任务”中查看进度')
const isDownloading = ref(false)

const dialogVisible = ref(false)
const downloadType = ref(1)

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
        type: Array as PropType<Array<{value: number, label: string}>>,
        required:true,
        default:() => []
    },
    downloadData: {
        type: Object,
        default: () => ({})
    }
});
const emits = defineEmits(["close"]);

onMounted(() => {
    openDialog();
});

const openDialog = () => {
    dialogVisible.value = props.dialogVisible
    // Set default download type to the first option if available
    if (props.downloadTypeOptions && props.downloadTypeOptions.length > 0) {
        downloadType.value = props.downloadTypeOptions[0].value
    }
}

const closeDialog = () => {
    content.value = '任务将提交到下载队列，可在“下载任务”中查看进度'
    emits("close")
}

const download = async () => {
    isDownloading.value = true
    content.value = '正在提交任务...'
    
    try {
        switch (props.prodType) {
            case 2: // Ebook
                await EbookDownload(props.downloadId, downloadType.value, props.enId)
                break;
            case 66: // Course
                await CourseDownload(props.downloadId, props.articleId, downloadType.value, props.enId)
                break;
            case 3: // Odob
                await OdobDownload(props.downloadId, downloadType.value, props.downloadData as any)
                break;
        }
        ElMessage({
            message: '任务已加入下载队列',
            type: 'success'
        })
        window.dispatchEvent(new CustomEvent('download-task:open'))
    } catch (error) {
        ElMessage({
            message: String(error),
            type: 'warning'
        })
    } finally {
        isDownloading.value = false
        closeDialog()
    }
}
</script>

<style scoped>
.download-container {
    padding: 10px 20px;
}

.format-selector {
    margin-bottom: 24px;
}

.section-label {
    font-size: 14px;
    color: var(--text-secondary, #606266);
    margin-bottom: 12px;
    font-weight: 500;
}

.format-options {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
}

.format-option {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    height: 40px;
    border: 1px solid var(--border-color, #dcdfe6);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    background-color: var(--fill-color-light, #f5f7fa);
    font-size: 14px;
    color: var(--text-primary, #303133);
}

.format-option:hover {
    border-color: var(--primary-color, #409eff);
    color: var(--primary-color, #409eff);
    background-color: var(--primary-color-light-9, #ecf5ff);
}

.format-option.active {
    border-color: var(--primary-color, #409eff);
    background-color: var(--primary-color, #409eff);
    color: white;
    font-weight: 500;
}

.selected-icon {
    position: absolute;
    right: 4px;
    top: 4px;
    font-size: 12px;
}

.download-status {
    margin-top: 20px;
    background: var(--fill-color-lighter, #fafafa);
    padding: 16px;
    border-radius: 8px;
    border: 1px solid var(--border-color-lighter, #ebeef5);
}

.status-header {
    display: flex;
    justify-content: center;
    margin-bottom: 8px;
    font-size: 13px;
}

.status-text {
    color: var(--text-regular, #606266);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
}

.dialog-footer {
    display: flex;
    justify-content: center;
    gap: 16px;
    padding-bottom: 8px;
}
</style>
