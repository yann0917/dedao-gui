<template>
    <div class="knowledge-container">
        <div class="knowledge-layout">
            <div class="content-area">
                <notes-item :topic-detail="topicDetail" :key="timer()"/>
            </div>
            <div class="sidebar-area">
                <topic-item :topic-detail="topicDetail" @send-detail="getTopicDetail"/>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import NotesItem from "../components/NotesItem.vue";
import TopicItem from "../components/TopicItem.vue";
import {reactive} from "vue";
import {services} from "../../wailsjs/go/models";

const topicDetail = reactive(new services.TopicIntro)

const getTopicDetail = (row:any)=> {
    Object.assign(topicDetail, row)
    timer()
}

const timer = ()=> {
    return new Date().getTime()
}
</script>

<style scoped>
.knowledge-container {
    height: 100%;
    padding: 32px;
    max-width: 1600px;
    margin: 0 auto;
    box-sizing: border-box;
}

.knowledge-layout {
    display: grid;
    grid-template-columns: 320px 1fr;
    gap: 32px;
    height: 100%;
    align-items: start;
}

.content-area {
    height: 100%;
    overflow-y: auto;
    border-radius: 16px;
    /* Hide scrollbar for cleaner look */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE/Edge */
}

.content-area::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}

.sidebar-area {
    height: 100%;
    overflow-y: auto;
    border-radius: 12px;
    order: -1; /* Move sidebar to left on desktop */
}

@media (max-width: 1024px) {
    .knowledge-layout {
        grid-template-columns: 280px 1fr;
        gap: 24px;
    }
    
    .knowledge-container {
        padding: 24px;
    }
}

@media (max-width: 768px) {
    .knowledge-layout {
        grid-template-columns: 1fr;
        grid-template-rows: 1fr auto;
    }

    .sidebar-area {
        height: auto;
        max-height: 300px;
        order: 2; /* Put sidebar at bottom on mobile */
    }

    .content-area {
        order: 1;
    }
    
    .knowledge-container {
        padding: 16px;
    }
}
</style>
