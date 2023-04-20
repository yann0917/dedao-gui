<template>
    <el-form :model="form" label-width="140px" style="max-width: 640px">
        <el-form-item label="文件保存目录">
            <el-input v-model="form.downloadDir"
                      :prefix-icon="Folder"
                      @blur="openDialogDir('')"
                      clearable @clear="clearDir "/>
        </el-form-item>
        <el-tooltip
                effect="light"
                content="下载 mp3 格式需要，安装后终端执行 whereis ffmpeg 查找路径"
                placement="right"
        >
            <el-form-item label="ffmpeg目录">
                <el-input v-model="form.ffmpegDir"
                          :prefix-icon="Folder"
                          clearable/>
            </el-form-item>
        </el-tooltip>
        <el-tooltip
                effect="light"
                content="电子书下载 pdf 格式需要，安装后终端执行 whereis wkhtmltopdf 查找路径"
                placement="right"
        >
            <el-form-item label="wkhtmltopdf目录">
                <el-input v-model="form.wkhtmltopdfDir"
                          :prefix-icon="Folder"
                          clearable/>
            </el-form-item>
        </el-tooltip>
        <el-form-item label="主题色">
            <el-color-picker class="color-block"
                             v-model="form.systemColor"
                             :predefine="predefineColors"
                             show-alpha
                             @change="setThemeColor(form.systemColor)"/>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="onSubmit">保存</el-button>
            <el-button>取消</el-button>
        </el-form-item>
    </el-form>
    <el-divider content-position="left">说明</el-divider>
    <div class="notes">
        <ul style="text-align: left">
            <li v-for="item in notes"> {{ item.content }} {{ item.link }}</li>
        </ul>
    </div>
</template>

<script lang="ts" setup>
import {reactive, ref} from 'vue'
import {Folder} from '@element-plus/icons-vue'
import {settingStore} from "../stores/setting";
import {OpenDirectoryDialog} from "../../wailsjs/go/backend/App";
import {setThemeColor} from "../utils/utils";

const store = settingStore();
const predefineColors = ref([
    '#ff6b00',
    '#ff4500',
    '#ff8c00',
    '#ffd700',
    '#90ee90',
    '#00ced1',
    '#1e90ff',
    '#409EFF',
    '#c71585',
    'rgba(255, 69, 0, 0.68)',
    'hsv(51, 100, 98)',
    'hsva(120, 40, 94, 0.5)',
    'hsl(181, 100%, 37%)',
    'hsla(209, 100%, 56%, 0.73)',
    '#c7158577',
])

const form = reactive({
    downloadDir: store.getDownloadDir,
    systemColor: store.getColor,
    ffmpegDir: store.getFfmpegDirDir,
    wkhtmltopdfDir: store.getWkDir,
})

const notes = reactive([{
        content: "音频需要借助[ffmpeg]合成",
        link: "https://ffmpeg.org",
    }, {
        content: "电子书转 PDF 需要借助[wkhtmltopdf]",
        link: "https://wkhtmltopdf.org/downloads.html",
    }, {
        content: "如需下载音频，电脑上安装ffmpeg后，找到可执行文件的路径，并设置好目录",
        link: "",
    },{
        content: "如需下载 PDF 格式电子书，电脑上安装wkhtmltopdf，找到可执行文件的路径，并设置好目录",
        link: "",
    },{
        content: "每天听本书和电子书均不支持批量下载",
        link: "",
    }
])

const openDialogDir = async (title: string) => {
    if (store.getDownloadDir == "") {
        await OpenDirectoryDialog(title).then((result) => {
            form.downloadDir = result
            console.log(result)
        }).catch((error) => {
            console.log(error)
        })
    }
}

const clearDir = () => {
    store.setting.downloadDir = ''
}

const onSubmit = () => {
    console.log(form)
    store.setting.downloadDir = form.downloadDir
    store.setting.theme = form.systemColor
    store.setting.ffmpegDir = form.ffmpegDir
    store.setting.wkhtmltopdfDir = form.wkhtmltopdfDir
}
</script>

<style scoped>

</style>