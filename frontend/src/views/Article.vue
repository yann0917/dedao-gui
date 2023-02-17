<template>
    <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="parentPath">{{breadcrumbItem1}}</el-breadcrumb-item>
        <el-breadcrumb-item>文稿</el-breadcrumb-item>
    </el-breadcrumb>
    <div class="markdown-body">
        <div v-highlight v-html="htmlStr" id="content"></div>
    </div>
</template>


<script lang="ts" setup>
import { ref, reactive, onMounted, watch } from 'vue'
import 'element-plus/es/components/message/style/css'
import { ElMessage } from 'element-plus'
import { ArticleDetail } from '../../wailsjs/go/backend/App'
import { useRoute, useRouter } from 'vue-router'

import { marked } from 'marked'

const route = useRoute()
const router = useRouter()

marked.setOptions({
    pedantic: false,
    gfm: true,
    breaks: false,
    sanitize: false,
    smartLists: true,
    smartypants: false,
    xhtml: false,
})

let id = ref()
let classId = ref()
let from = ref()
let breadcrumbItem1 = ref()
let htmlStr = ref()
let aType = 0
let parentPath = reactive({path:"",query:{}})

onMounted(() => {
    watch(() => {
        from.value = route.query.from
        id.value = route.params.id
        classId.value = route.query.class_id

        if (from.value == "course") {
            aType = 1
            parentPath.path = '/'+from.value+'/'+classId.value
            parentPath.query={total: route.query.total, enid: route.query.class_enid, title: route.query.parentTitle}
            breadcrumbItem1.value = route.query.title
        }
        if (from.value == "odob") {
            aType = 2
            parentPath.path= '/'+from.value
            breadcrumbItem1.value = '听书书架'
        }
    },
        () => articleDetail(aType, id.value),
        { immediate: true }
    )
})

const articleDetail = async (aType: number, enid: string) => {
    await ArticleDetail(aType, enid).then((info) => {
        // console.log(info)
        htmlStr.value = marked(info)
    }).catch((error) => {
        ElMessage({
            message: error,
            type: 'warning'
        })
    })
    return
}

</script>
<style scoped>
.markdown-body {
    text-align: left;
    line-height: 1.6;
}
</style>