<template>
    <div>
        <el-space wrap>
            <el-row :gutter="10">
                <el-col :span="16">
                    <notes-item :topic-detail="topicDetail" :key="timer"/>
                </el-col>
                <el-col :span="8" justify="center">
                    <topic-item :topic-detail="topicDetail" @send-detail="getTopicDetail"/>
                </el-col>
            </el-row>
        </el-space>
    </div>
</template>

<script lang="ts" setup>
import NotesItem from "../components/NotesItem.vue";
import TopicItem from "../components/TopicItem.vue";
import {useRoute} from "vue-router";
import {onMounted, reactive} from "vue";
import {services} from "../../wailsjs/go/models";

const route = useRoute()
const topicDetail = reactive(new services.TopicIntro)

onMounted(() => {
    // console.log(route.query)
    // Object.assign(topicDetail, route.query)
    // watch(() => {
    //         id.value = route.params.id
    //         enid.value = route.query.enid
    //         total.value = Number(route.query.total)
    //         breadcrumbTitle.value = route.query.title
    //     },
    //     () => getTableData(),
    //     {immediate: true}
    // )

})
const getTopicDetail = (row:any)=> {
    Object.assign(topicDetail, row)
    console.log('topic子组件向父组件传值', row)
    timer()
}

const timer = ()=> {
    return new Date().getTime()
}
</script>

<style scoped lang="scss">
</style>
