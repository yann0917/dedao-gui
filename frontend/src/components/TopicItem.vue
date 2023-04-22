<template>
    <p style="text-align: left; font-size: 20px; font-weight: bold"><span style="color:#ff6b00"># </span>推荐话题</p>
      <div class="topic">
          <ul v-infinite-scroll="loadTopic"
              class="infinite-list"
              style="overflow-y: auto; max-height: 800px;"
              :infinite-scroll-disabled="infLoadingTopic"
              infinite-scroll-distance="10">
              <li v-for="i in topicList.list"  class="infinite-list-item">
                  <el-button text class="topic-item" @click="sendTopicDetail(i)">
                      <div class="title">
                          <span style="color:#ff6b00">#</span>{{i.name}}
                          <el-badge :value="i.tag==1?'新':i.tag==2?'热':''" class="badge-item" :type="i.tag==2?'warning':'danger'">
                          </el-badge>
                      </div>
                      <div class="intro">{{i.intro}}</div>
                  </el-button>
              </li>
          </ul>
      </div>
</template>

<script setup lang="ts">
import {ref, reactive, onMounted, defineProps} from 'vue'
import { ElMessage } from 'element-plus'
import { TopicAll } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { useRoute, useRouter } from 'vue-router'
import { userStore } from '../stores/user';

const store = userStore()
const router = useRouter()
const route = useRoute()
const page = ref(1)
const total = ref(0)
const pageSize = ref(15)
const infLoadingTopic = ref(true)

let topicAll = reactive(new services.TopicAll)
const topicList = reactive(new services.TopicAll)
topicList.list = []

const emits = defineEmits(["sendDetail"]);

// let id = ref()
// let enid = ref()

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
    // console.log(row)
    emits("sendDetail", row)
}
</script>

<style scoped lang="scss">
.topic {
    height: 800px;
    border: 1px solid var(--el-border-color);
    border-radius:8px;
}
.topic-item {
    height: 110px;
}
/* 深度选择器（样式穿透）*/
:deep .el-button>span {
    display: inline-block;
    align-items: center;
}
.topic-item .title {
    font-size: large;
    font-weight: bold;
    text-align: left;
    margin-bottom: 10px;
}
.topic-item .intro{
    white-space: pre-wrap;
    font-weight:normal;
    text-align: left;
    line-height: 1.5;
}
ul {
    padding: 0;
    margin: 0;
    text-align: left;
    list-style: none;
}

.infinite-list {
    height: 90%;
    padding: 0;
    margin: 0;
    text-align: left;
    list-style: none;
}
.infinite-list .infinite-list-item {
    margin: 1%;
    /*background: var(--el-color-primary-light-9);*/
    /*color: var(--el-color-primary);*/
    color: #606266;
    padding-bottom: 1%;
    border-bottom: 1px dashed;
}
.infinite-list .infinite-list-item:last-child {
    margin-bottom: 0;
}

.infinite-list .infinite-list-item+.list-item {
    margin-top: 10px;
}

</style>