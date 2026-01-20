<template>
  <div class="ebook-comment-page">
    <div class="page-header">
      <el-breadcrumb :separator-icon="ArrowRight">
        <el-breadcrumb-item :to="{ path: '/bought/ebook' }">电子书架</el-breadcrumb-item>
        <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
      </el-breadcrumb>
      
      <!-- 简单的数据概览 -->
      <div class="overview">
        <el-tag type="info" effect="plain" round>总评论: {{ tableData.total }}</el-tag>
        <el-tag type="success" effect="plain" round>平均评分: {{ averageScore.toFixed(1) }}</el-tag>
      </div>
    </div>
 
    <div 
      class="waterfall-container" 
      v-infinite-scroll="loadMore"
      :infinite-scroll-disabled="loading || !hasMore"
      :infinite-scroll-distance="10"
      :infinite-scroll-immediate="false"
    >
      <div class="waterfall-column" v-for="colIndex in 3" :key="colIndex">
        <el-card 
          v-for="(item, index) in columns[`col${colIndex}`]" 
          :key="item.note_id" 
          class="waterfall-item" 
          shadow="hover"
        >
          <template #header>
            <div class="card-header">
              <div class="user-info">
                <el-avatar :size="40" :src="item.notes_owner.avatar" />
                <div class="user-meta">
                  <span class="user-name">{{ item.notes_owner.name }}</span>
                  <el-tag v-if="item.notes_owner.slogan" size="small" type="info" class="user-slogan" effect="plain">
                    {{ item.notes_owner.slogan }}
                  </el-tag>
                </div>
              </div>
              <div class="note-meta">
                <h3 class="note-title">{{ item.note_title }}</h3>
                <div class="likes" v-if="item.notes_count?.like_count">
                  👍 {{ item.notes_count?.like_count }}
                </div>
              </div>
            </div>
          </template>
          <div class="card-content">
            <span v-html="handleStyleNote(item.style_note_line)"></span>
          </div>
        </el-card>
      </div>
    </div>

    <div v-if="loading" class="loading-container">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    
    <div v-if="!hasMore && tableData.list.length > 0" class="no-more">
      没有更多数据了
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowRight } from '@element-plus/icons-vue'
import { EbookCommentList } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { useRoute, useRouter } from 'vue-router'
import { userStore } from '../stores/user';

const store = userStore()
const router = useRouter()
const route = useRoute()
const loading = ref(false)
const page = ref(1)
const hasMore = ref(true)
const pageSize = ref(15)

let averageScore = ref(0)
let scoreInfo = ref()

// 将数据分为三列
const columns = reactive({
  col1: [] as any[],
  col2: [] as any[],
  col3: [] as any[]
})

let tableData = reactive({
  list: [] as any[],
  total: 0,
  ebook_score: {
    average_score: 0,
    score_info: {}
  }
})

let id = ref()
let enid = ref()
let breadcrumbTitle = ref()

onMounted(() => {
  watch(() => {
    id.value = route.query.id
    enid.value = route.query.enid
    breadcrumbTitle.value = route.query.title
  },
    () => {
      page.value = 1
      tableData.list = []
      getTableData()
    },
    { immediate: true }
  )
})

const loadMore = async () => {
  if (loading.value || !hasMore.value) return
  page.value++
  await getTableData()
}

// 将数据分配到三列中
const distributeData = (newData: any[]) => {
  newData.forEach((item, index) => {
    const colIndex = index % 3
    switch(colIndex) {
      case 0:
        columns.col1.push(item)
        break
      case 1:
        columns.col2.push(item)
        break
      case 2:
        columns.col3.push(item)
        break
    }
  })
}

const getTableData = async () => {
  if (!hasMore.value) return
  
  loading.value = true
  try {
    const result = await EbookCommentList(enid.value, page.value, pageSize.value)
    
    if (page.value === 1) {
      // 重置所有列
      columns.col1 = []
      columns.col2 = []
      columns.col3 = []
      tableData.list = result.list
    } else {
      tableData.list.push(...result.list)
    }
    
    // 分配新数据到列中
    distributeData(result.list)
    
    tableData.total = result.total
    averageScore.value = Number(result.ebook_score.average_score)
    scoreInfo.value = result.ebook_score.score_info
    
    hasMore.value = tableData.list.length < result.total
    
  } catch (error) {
    ElMessage({
      message: error as string,
      type: 'warning'
    })
  } finally {
    loading.value = false
  }
}

const handleStyleNote = (style_note_line: string) => {
  let notes = JSON.parse(style_note_line)
  let htmlString = ''
  for (let item of notes) {
    switch (item.type) {
      case "paragraph":
        htmlString += "<p>"
        for (let ele of item.contents) {
          switch (ele.type) {
            case "text":
              let temp = ele.text.content.trim()
              if (ele.text.bold) {
                htmlString += "<span style='font-weight:bold;'>" + temp + "</span>"
              } else {
                htmlString += temp
              }
              break
          }
        }
        htmlString += "</p>"
        break

      case "list":
        if (item.ordered) {
          htmlString += "<ol>"
        } else {
          htmlString += "<ul>"
        }

        for (let ele1 of item.contents) {
          let subHtml = ""
          for (let el of ele1) {
            if (el.type === "text") {
              let temp = el.text.content.trim()
              if (temp.length > 0) {
                if (el.text.bold) {
                  subHtml += "<li style='font-weight:bold;'>" + temp + "</li>"
                } else {
                  subHtml += "<li>" + temp + "</li>"
                }
              }
            }
          }
          htmlString += subHtml
        }
        if (item.ordered) {
          htmlString += "</ol>"
        } else {
          htmlString += "</ul>"
        }
        break
    }
  }

  return htmlString
}

</script>

<style scoped lang="scss">
.ebook-comment-page {
  height: calc(100vh - 60px);
  overflow-y: auto;
  box-sizing: border-box;
  padding: 24px;
  background: var(--fill-color-light);
  
  /* 隐藏滚动条但保留功能 */
  scrollbar-width: none;
  -ms-overflow-style: none;
  &::-webkit-scrollbar {
    display: none;
  }
}

.page-header {
  margin-bottom: 24px;
}

.overview {
  margin-top: 16px;
  .el-tag {
    margin-right: 12px;
  }
}

.waterfall-container {
  display: flex;
  gap: 20px;
  align-items: flex-start;
  
  @media screen and (max-width: 1200px) {
    gap: 15px;
  }
  
  @media screen and (max-width: 768px) {
    flex-direction: column;
  }
}

.waterfall-column {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  
  @media screen and (max-width: 768px) {
    width: 100%;
  }
}

.waterfall-item {
  width: 100%;
  border: none;
  border-radius: 12px;
  background: var(--card-bg);
  box-shadow: var(--shadow-soft);
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-medium);
  }
  
  :deep(.el-card__header) {
    padding: 16px;
    border-bottom: 1px solid var(--border-color-light);
  }
  
  :deep(.el-card__body) {
    padding: 16px;
  }
}

.card-header {
  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
    
    .user-meta {
      display: flex;
      flex-direction: column;
      gap: 4px;
      align-items: flex-start;
    }
    
    .user-name {
      font-size: 14px;
      font-weight: 500;
      color: var(--text-primary);
    }
    
    .user-slogan {
      max-width: 200px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
  
  .note-meta {
    .note-title {
      margin: 0 0 8px;
      font-size: 16px;
      font-weight: 600;
      color: var(--text-primary);
      line-height: 1.4;
    }
    
    .likes {
      font-size: 12px;
      color: var(--text-secondary);
    }
  }
}

.card-content {
  color: var(--text-regular);
  line-height: 1.6;
  font-size: 14px;
  
  :deep(p) {
    margin: 8px 0;
  }
  
  :deep(ul), :deep(ol) {
    padding-left: 20px;
    margin: 8px 0;
  }
  
  :deep(li) {
    margin: 4px 0;
  }
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: var(--text-secondary);
  
  .el-icon {
    margin-right: 8px;
    font-size: 20px;
  }
}

.no-more {
  text-align: center;
  padding: 20px;
  color: var(--text-secondary);
  font-size: 14px;
}
</style>