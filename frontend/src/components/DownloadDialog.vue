<template>
    <el-dialog v-model="dialogVisible" title="请选择下载格式" align-center center width="30%" :before-close="closeDialog">
        <el-form >
            <el-form-item label="下载格式" label-width="80px">
                <el-select v-model="downloadType" placeholder="请选择下载格式">
                    <el-option v-for="item in props.downloadTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-form-item>
            <el-progress v-show="percentage"
                    :percentage="percentage"
                    status="success"
                    :stroke-width="20"
                    :text-inside="true"
            ><span>{{content}}</span></el-progress>
        </el-form>

        <template #footer>
      <span class="dialog-footer">
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" @click="download()">
          确认
        </el-button>
      </span>
        </template>
    </el-dialog>

</template>

<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {EbookDownload, CourseDownload} from "../../wailsjs/go/backend/App";
import {ElMessage} from "element-plus";
import { EventsOn, EventsOff} from "../../wailsjs/runtime/runtime";

let percentage=ref(0)
let content=ref('')

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
const closeDialog = () => {
    //   initForm()
    // downloadType.value = 1
    // dialogVisible.value = false
    EventsOff("courseDownload", "ebookDownload")
    percentage.value = 0
    content.value = ''
    emits("close")
}

const download = async () => {
    // percentage.value = 66
    content.value = '下载中'
    switch (props.prodType) {
        case 2:
        EventsOn("ebookDownload",  data=>{
                if (data) {
                    console.log(data)
                    percentage.value = data.pct
                    content.value = data.value + '下载中...'
                }
            })
            await EbookDownload(props.downloadId, downloadType.value, props.enId).then((info) => {
                // console.log(info)
            }).catch((error) => {
                ElMessage({
                    message: error,
                    type: 'warning'
                })
            })
            break;
        case 66:
            EventsOn("courseDownload",  data=>{
                if (data) {
                    console.log(data)
                    percentage.value = data.pct
                    content.value = data.value + '下载中...'
                }
            })
            await CourseDownload(props.downloadId, props.articleId, downloadType.value, props.enId).then((info) => {
                console.log(info)
            }).catch((error) => {
                ElMessage({
                    message: error,
                    type: 'warning'
                })
            })
            break;
    }
    closeDialog()
    return
}

</script>

<style scoped>

</style>