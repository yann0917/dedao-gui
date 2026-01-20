<template>
    <div class="topic-container" v-if="topicList.list.length>0">
        <div class="topic-header">
            <h3 class="section-title"><span class="highlight">#</span> 推荐话题</h3>
        </div>
        <div class="topic-list-wrapper">
            <ul v-infinite-scroll="loadTopic"
                class="topic-list"
                :infinite-scroll-disabled="infLoadingTopic"
                infinite-scroll-distance="10">
                <li v-for="i in topicList.list" :key="i.notes_topic_id" class="topic-list-item" @click="sendTopicDetail(i)">
                    <div class="item-content">
                        <div class="item-title">
                            <span class="hash">#</span>
                            <span class="text">{{i.name}}</span>
                            <el-tag v-if="i.tag==1" size="small" type="danger" effect="plain" class="badge">新</el-tag>
                            <el-tag v-if="i.tag==2" size="small" type="warning" effect="plain" class="badge">热</el-tag>
                        </div>
                        <div class="item-intro">{{i.intro}}</div>
                    </div>
                </li>
            </ul>
        </div>
    </div>
</template>

<script setup lang="ts">
import {ref, reactive, onMounted} from 'vue'
import { ElMessage } from 'element-plus'
import { TopicAll } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'

const infLoadingTopic = ref(true)

let topicAll = reactive(new services.TopicAll)
const topicList = reactive(new services.TopicAll)
topicList.list = []

const emits = defineEmits(["sendDetail"]);

const topicPage = ref(0)
const topicHasMore = ref(false)
const topicPageSize = ref(20)

const props = defineProps({
    topicDetail: {
        type: Object,
        default: services.TopicIntro
    }
})

onMounted(() => {
    getTopicAll()
})

const loadTopic = () => {
    infLoadingTopic.value=true
    if (topicHasMore.value === true) {
        topicPage.value += 1
        getTopicAll()
    }
}

const getTopicAll =async () => {
    await TopicAll(topicPage.value, topicPageSize.value).then((result)=>{
        Object.assign(topicAll, result)
        topicHasMore.value= topicAll.has_more
        topicList.has_more = topicAll.has_more
        topicList.list.push(...topicAll.list)
    }).catch((error)=>{
        ElMessage({
            message: error,
            type: 'warning'
        })
    })
    infLoadingTopic.value = false
}

const sendTopicDetail = (row: any) => {
    emits("sendDetail", row)
}
</script>

<style scoped>
.topic-container {
    height: 100%;
    display: flex;
    flex-direction: column;
    background: var(--card-bg);
    border-radius: 12px;
    box-shadow: var(--shadow-soft);
    border: 1px solid var(--border-soft);
    overflow: hidden;
}

.topic-header {
    padding: 16px 20px;
    border-bottom: 1px solid var(--border-soft);
    background: var(--card-bg);
}

.section-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
    display: flex;
    align-items: center;
    gap: 4px;
}

.highlight {
    color: var(--accent-color);
}

.topic-list-wrapper {
    flex: 1;
    overflow-y: auto;
}

.topic-list {
    padding: 0;
    margin: 0;
    list-style: none;
}

.topic-list-item {
    padding: 16px 20px;
    cursor: pointer;
    transition: all 0.2s ease;
    border-bottom: 1px solid var(--border-lighter);
}

.topic-list-item:last-child {
    border-bottom: none;
}

.topic-list-item:hover {
    background-color: var(--fill-color-light);
}

.item-content {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.item-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
    display: flex;
    align-items: center;
    gap: 4px;
    line-height: 1.4;
}

.hash {
    color: var(--accent-color);
    font-weight: normal;
}

.badge {
    margin-left: 4px;
    height: 18px;
    padding: 0 6px;
    font-size: 11px;
}

.item-intro {
    font-size: 13px;
    color: var(--text-secondary);
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
</style>