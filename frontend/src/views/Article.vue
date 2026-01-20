<template>
    <div class="article-container">
        <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="parentPath">{{breadcrumbItem1}}</el-breadcrumb-item>
            <el-breadcrumb-item>文稿</el-breadcrumb-item>
        </el-breadcrumb>
        <div class="markdown-body">
            <div v-highlight v-html="htmlStr" id="content"></div>
        </div>
    </div>
</template>


<script lang="ts" setup>
import { ref, reactive, onMounted, watch } from 'vue'
import 'element-plus/es/components/message/style/css'
import { ElMessage } from 'element-plus'
import { ArticleDetail } from '../../wailsjs/go/backend/App'
import { useRoute } from 'vue-router'
import { ROUTE_NAMES } from '../router/routes'

import { marked } from 'marked'

const route = useRoute()

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
let parentPath = reactive({name: '', params: {}, query: {}})


const PDFFile = {
    A4:[592.28, 841.89],
};


onMounted(() => {
    watch(() => {
        from.value = route.query.from
        id.value = route.params.id
        classId.value = route.query.class_id

        if (from.value == "course") {
            aType = 1
            parentPath.name = ROUTE_NAMES.ARTICLE_LIST
            parentPath.params = { id: String(classId.value) }
            parentPath.query={total: route.query.total, enid: route.query.class_enid, title: route.query.parentTitle}
            breadcrumbItem1.value = route.query.title
        }
        if (from.value == "odob") {
            aType = 2
            parentPath.name = ROUTE_NAMES.ODOB
            breadcrumbItem1.value = '听书书架'
        }
        if (from.value == "aiChannel") {
            aType = 1
            parentPath.name = ROUTE_NAMES.AI_CHANNEL
            breadcrumbItem1.value = route.query.parentTitle || 'AI学习圈'
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
<style scoped lang="scss">
.article-container {
    padding: 24px;
    background: var(--fill-color-light);
    height: calc(100vh - 60px);
    overflow-y: auto;
    box-sizing: border-box;

    /* 隐藏滚动条但保留功能 - 清新风格 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE 10+ */
}

.article-container::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}

.markdown-body {
    color: var(--text-primary);
    text-align: left;
    line-height: 1.8;
    background: var(--card-bg);
    padding: 32px;
    border-radius: 12px;
    box-shadow: var(--shadow-soft);
    margin-top: 20px;
    max-width: 900px; /* Improve readability */
    margin-left: auto;
    margin-right: auto;
}

:deep(#content h2>code){
    background-color: var(--accent-color);
    padding: 2px 6px;
    border-radius: 4px;
    color: white;
}

:deep(#content p>em){
    font-style: normal;
    color: var(--accent-color);
}

/* Markdown content styles */
:deep(#content) {
    h1, h2, h3, h4, h5, h6 {
        color: var(--text-primary);
        margin-top: 24px;
        margin-bottom: 16px;
        font-weight: 600;
    }

    p {
        color: var(--text-primary);
        margin-bottom: 16px;
        font-size: 16px;
    }

    blockquote {
        border-left: 4px solid var(--accent-color);
        background: var(--fill-color);
        padding: 16px;
        margin: 16px 0;
        color: var(--text-secondary);
        border-radius: 0 4px 4px 0;
    }

    ul, ol {
        padding-left: 24px;
        margin-bottom: 16px;
        color: var(--text-primary);
    }

    li {
        margin-bottom: 8px;
    }

    img {
        max-width: 100%;
        border-radius: 8px;
        margin: 16px 0;
    }

    hr {
        border: none;
        border-top: 1px solid var(--border-soft);
        margin: 24px 0;
    }
}
</style>
