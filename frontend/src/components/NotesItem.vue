<template>
  <div class="notes-container">
    <div class="topic-header" v-if="hasTopicDetail">
      <el-image class="topic-cover" :src="ossProcess(props.topicDetail.img)" fit="cover" />
      <div class="topic-info">
        <h2 class="topic-title"><span class="hash">#</span>{{ props.topicDetail.name }}</h2>
        <div class="topic-stats">
          <span class="stat-item">阅读 {{ props.topicDetail.view_count }}</span>
          <span class="stat-item">讨论 {{ props.topicDetail.notes_count }}</span>
        </div>
        <div class="topic-intro">{{ props.topicDetail.intro }}</div>
      </div>
    </div>

    <div class="topic-tabs" v-if="hasTopicDetail">
      <div 
        class="tab-item" 
        :class="isElected ? 'active' : ''" 
        @click="handleElected(true)"
      >精选</div>
      <div 
        class="tab-item" 
        :class="!isElected ? 'active' : ''" 
        @click="handleElected(false)"
      >最新</div>
    </div>

    <div class="notes-list-container" v-infinite-scroll="loadNotes" :infinite-scroll-disabled="infLoadingNote" infinite-scroll-distance="10">
      <div class="note-item" v-for="item in noteList.notes" :key="item.f_part.note_id_hazy">
        <!-- Repost Header -->
        <div class="repost-info" v-if="item.comb?.length">
          <el-icon><Connection /></el-icon>
          <span>{{ handleComb(item.comb) }}</span>
        </div>

        <!-- User Info -->
        <div class="user-info">
          <div class="avatar-wrapper">
            <el-avatar :size="40" :src="item.f_part.avatar" />
            <div class="vip-badge" v-if="item.f_part.is_v">
               <el-icon v-if="item.f_part.is_v==2" color="#ff6b00"><Medal /></el-icon>
               <el-icon v-if="item.f_part.is_v==3||item.f_part.is_v==4" color="#ff6b00"><Trophy /></el-icon>
            </div>
          </div>
          <div class="user-meta">
            <div class="name-line">
              <span class="username">{{ item.f_part.nick_name }}</span>
              <el-tag v-if="item.f_part.slogan" size="small" type="success" effect="plain">{{ item.f_part.slogan }}</el-tag>
            </div>
            <div class="time-line">{{ item.f_part.time_desc }}</div>
          </div>
        </div>

        <!-- Content -->
        <div class="note-content">
          <div class="text">{{ item.f_part.note }}</div>
          
          <!-- Images -->
          <div class="image-grid" v-if="item.f_part.images?.length">
            <el-image 
              v-for="(img, idx) in item.f_part.images" 
              :key="idx"
              :src="img" 
              :preview-src-list="item.f_part.images"
              :initial-index="idx"
              fit="cover" 
              class="note-image"
            />
          </div>

          <!-- Book Rating -->
          <div class="rating-info" v-if="item.notes_type == 7">
            给这本书评了 <span class="score">{{ item.f_part.note_score }}</span> 分
          </div>

          <!-- Quote/Source -->
          <div class="source-card" v-if="item.f_part.base_source.title">
            <div class="quote-text" v-if="item.notes_type == 4">
              {{ item.f_part.note_line }}
            </div>
            <div class="source-link">
              <el-image :src="item.f_part.base_source.img" class="source-cover" fit="cover" />
              <div class="source-details">
                <div class="source-title">{{ item.f_part.base_source.title }}</div>
                <div class="source-subtitle">{{ item.f_part.base_source.sub_title }}</div>
              </div>
            </div>
          </div>

          <!-- Second Part (Reposted Content) -->
          <div class="reposted-content" v-if="item.s_part">
            <div class="repost-user">
              <el-avatar :size="24" :src="item.s_part.avatar" />
              <span class="repost-username">{{ item.s_part.nick_name }}</span>
            </div>
            <div class="repost-text">{{ item.s_part.note }}</div>
             <div class="image-grid" v-if="item.s_part.images?.length">
                <el-image 
                  v-for="(img, idx) in item.s_part.images" 
                  :key="idx"
                  :src="img" 
                  :preview-src-list="item.s_part.images"
                  :initial-index="idx"
                  fit="cover" 
                  class="note-image"
                />
            </div>
             <div class="source-card" v-if="item.s_part.base_source.title">
                <div class="quote-text" v-if="item.notes_type == 4">
                  {{ item.s_part.note_line }}
                </div>
                <div class="source-link">
                  <el-image :src="item.s_part.base_source.img" class="source-cover" fit="cover" />
                  <div class="source-details">
                    <div class="source-title">{{ item.s_part.base_source.title }}</div>
                    <div class="source-subtitle">{{ item.s_part.base_source.sub_title }}</div>
                  </div>
                </div>
              </div>
          </div>
          
           <!-- Topic Tag -->
           <div class="note-topic-tag" v-if="item.topic?.topic_name">
             <el-tag round size="small">#{{ item.topic.topic_name }}</el-tag>
           </div>
        </div>

        <!-- Actions -->
        <div class="note-actions">
          <div class="action-item">
            <el-icon><Share /></el-icon>
            <span>{{ item.note_count?.repost_count || '转发' }}</span>
          </div>
          <div class="action-item">
            <el-icon><ChatDotRound /></el-icon>
             <span>{{ item.note_count?.comment_count || '评论' }}</span>
          </div>
          <div class="action-item">
            <el-icon><Star /></el-icon>
             <span>{{ item.note_count?.like_count || '点赞' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Connection, Medal, Trophy, Share, ChatDotRound, Star } from '@element-plus/icons-vue'
import { NotesTimeline, TopicNotesList } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'

const page = ref(0)
const total = ref(0)
const pageSize = ref(20)
const infLoadingNote = ref(false)

const tableData = reactive(new services.NotesTimeline())
const noteList = reactive(new services.NotesTimeline())
noteList.notes = []

const maxId = ref("")
const topic_id_hazy = ref("")
const isElected = ref(true)

const props = defineProps({
    topicDetail: {
        type: Object,
        default: () => new services.TopicIntro()
    }
})

const hasTopicDetail = computed(() => {
    return props.topicDetail && Object.keys(props.topicDetail).length > 0 && props.topicDetail.topic_id_hazy
})

const ossProcess = (url: string) => {
    if (!url) return ''
    return url + "?x-oss-process=image/resize,m_fill,h_200,w_200"
}

onMounted(() => {
    watch(() => props.topicDetail?.topic_id_hazy,
        (newVal) => {
            topic_id_hazy.value = newVal || ''
            if (topic_id_hazy.value) {
                noteList.notes = []
                isElected.value = true
                page.value = 0
                getTableData()
            } else {
                 // Default timeline if no topic selected
                 if (noteList.notes.length === 0) {
                     getTableData()
                 }
            }
        },
        { immediate: true }
    )
})

const loadNotes = () => {
    if (infLoadingNote.value) return
    infLoadingNote.value = true
    if (noteList.is_more) {
        if (topic_id_hazy.value) {
            page.value += 1
        }
        getTableData()
    } else {
        infLoadingNote.value = false
    }
}

const getTableData = async () => {
    try {
        if (topic_id_hazy.value) {
            const table = await TopicNotesList(topic_id_hazy.value, isElected.value, page.value, pageSize.value)
            const list = new services.NotesList()
            Object.assign(list, table)
            
            // Process images
            list.note_detail_list.forEach(processImages)
            
            tableData.is_more = list.has_more
            tableData.notes = list.note_detail_list
            noteList.is_more = tableData.is_more
            noteList.notes.push(...tableData.notes)
        } else {
            const table = await NotesTimeline(maxId.value)
            Object.assign(tableData, table)
            maxId.value = tableData.max_id
            
             // Process images
            tableData.notes.forEach(processImages)
            
            noteList.is_more = tableData.is_more
            noteList.max_id = tableData.max_id
            noteList.notes.push(...tableData.notes)
        }
    } catch (error) {
        ElMessage({
            message: '加载失败: ' + error,
            type: 'warning'
        })
    } finally {
        infLoadingNote.value = false
    }
}

const processImages = (item: any) => {
    if (item.f_part.images?.length > 0) {
        item.f_part.images.forEach((imgStr: string, idx: number, arr: string[]) => {
            try {
                // Check if it's a JSON string first
                if (imgStr.startsWith('{')) {
                    const imgObj = JSON.parse(imgStr)
                    arr[idx] = imgObj.url
                }
            } catch (e) {
                // If parse fails, assume it's already a URL or handle error
            }
        })
    }
    if (item.s_part?.images?.length > 0) {
         item.s_part.images.forEach((imgStr: string, idx: number, arr: string[]) => {
            try {
                if (imgStr.startsWith('{')) {
                    const imgObj = JSON.parse(imgStr)
                    arr[idx] = imgObj.url
                }
            } catch (e) {}
        })
    }
}

const handleElected = (val: boolean) => {
    if (isElected.value === val) return
    isElected.value = val
    noteList.notes = []
    page.value = 0
    getTableData()
}

const handleComb = (list: services.Comb[]) => {
    if (!list) return ''
    if (list.length <= 2) {
        return list.map(v => v.name).join("、") + "转发了"
    } else {
        return list.slice(0, 2).map(v => v.name).join("、") + "和其他" + (list.length - 2) + "转发了"
    }
}
</script>

<style scoped>
.notes-container {
    max-width: 800px;
    margin: 0 auto;
    padding-bottom: 40px;
}

.topic-header {
    display: flex;
    gap: 20px;
    margin-bottom: 24px;
    background: var(--card-bg);
    padding: 24px;
    border-radius: 12px;
    box-shadow: var(--shadow-soft);
}

.topic-cover {
    width: 100px;
    height: 100px;
    border-radius: 8px;
    flex-shrink: 0;
}

.topic-info {
    flex: 1;
}

.topic-title {
    margin: 0 0 12px 0;
    font-size: 20px;
    color: var(--text-primary);
}

.hash {
    color: var(--el-color-primary);
    margin-right: 4px;
}

.topic-stats {
    font-size: 13px;
    color: var(--text-secondary);
    margin-bottom: 12px;
}

.stat-item {
    margin-right: 16px;
}

.topic-intro {
    font-size: 14px;
    color: var(--text-regular);
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

.topic-tabs {
    display: flex;
    gap: 24px;
    margin-bottom: 20px;
    border-bottom: 1px solid var(--border-soft);
    padding: 0 12px;
}

.tab-item {
    padding: 12px 0;
    cursor: pointer;
    font-size: 15px;
    color: var(--text-secondary);
    position: relative;
    transition: all 0.3s;
}

.tab-item.active {
    color: var(--el-color-primary);
    font-weight: 600;
}

.tab-item.active::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 0;
    width: 100%;
    height: 2px;
    background: var(--el-color-primary);
}

.note-item {
    background: var(--card-bg);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 16px;
    box-shadow: var(--shadow-soft);
    border: 1px solid var(--border-soft);
    transition: transform 0.2s;
}

.note-item:hover {
    box-shadow: var(--shadow-medium);
}

.repost-info {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: var(--text-secondary);
    margin-bottom: 12px;
}

.user-info {
    display: flex;
    gap: 12px;
    margin-bottom: 12px;
}

.avatar-wrapper {
    position: relative;
}

.vip-badge {
    position: absolute;
    bottom: -4px;
    right: -4px;
    background: #fff;
    border-radius: 50%;
    padding: 2px;
    display: flex;
}

.user-meta {
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.name-line {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
}

.username {
    font-weight: 600;
    color: var(--text-primary);
    font-size: 15px;
}

.time-line {
    font-size: 12px;
    color: var(--text-secondary);
}

.note-content {
    margin-left: 52px; /* Align with avatar text */
}

.text {
    font-size: 15px;
    color: var(--text-regular);
    line-height: 1.6;
    white-space: pre-wrap;
    margin-bottom: 12px;
}

.image-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 8px;
    margin-bottom: 12px;
}

.note-image {
    border-radius: 8px;
    aspect-ratio: 1;
    width: 100%;
    cursor: zoom-in;
}

.rating-info {
    color: var(--el-color-warning);
    font-size: 14px;
    margin-bottom: 12px;
}

.score {
    font-weight: bold;
    font-size: 16px;
}

.source-card {
    background: var(--bg-color);
    border-radius: 8px;
    padding: 12px;
    margin-bottom: 12px;
}

.quote-text {
    font-style: italic;
    color: var(--text-secondary);
    margin-bottom: 12px;
    padding-left: 12px;
    border-left: 3px solid var(--border-soft);
}

.source-link {
    display: flex;
    gap: 12px;
    align-items: center;
}

.source-cover {
    width: 40px;
    height: 60px;
    border-radius: 4px;
    flex-shrink: 0;
}

.source-details {
    overflow: hidden;
}

.source-title {
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 4px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.source-subtitle {
    font-size: 12px;
    color: var(--text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.reposted-content {
    background: var(--bg-color);
    border-radius: 8px;
    padding: 12px;
    margin-bottom: 12px;
}

.repost-user {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
}

.repost-username {
    font-weight: 600;
    font-size: 13px;
}

.repost-text {
    font-size: 14px;
    color: var(--text-regular);
    margin-bottom: 8px;
}

.note-actions {
    margin-left: 52px;
    display: flex;
    gap: 32px;
    margin-top: 16px;
    padding-top: 12px;
    border-top: 1px solid var(--border-soft);
}

.action-item {
    display: flex;
    align-items: center;
    gap: 6px;
    color: var(--text-secondary);
    cursor: pointer;
    font-size: 13px;
    transition: color 0.2s;
}

.action-item:hover {
    color: var(--el-color-primary);
}

.note-topic-tag {
    margin-top: 8px;
}

@media (max-width: 768px) {
    .note-content, .note-actions {
        margin-left: 0;
    }
    .topic-header {
        flex-direction: column;
    }
    .topic-cover {
        width: 100%;
        height: 150px;
    }
}
</style>
