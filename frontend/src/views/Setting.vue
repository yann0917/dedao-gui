<template>
    <el-form :model="form" label-width="120px" style="max-width: 640px">
        <el-form-item label="文件保存目录">
            <el-input v-model="form.downloadDir"
                      :prefix-icon="Folder"
                      @blur="openDialogDir('')"
                      clearable @clear="clearDir "/>
        </el-form-item>
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
})

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
}
</script>

<style scoped>

</style>