<template>
  <el-breadcrumb :separator-icon="ArrowRight">
    <el-breadcrumb-item :to="{ path: '/bought/ebook' }">电子书架</el-breadcrumb-item>
    <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
  </el-breadcrumb>

  <!-- 简单的数据概览 -->
  <div class="overview">
    <el-tag type="info">总评论: {{ tableData.total }}</el-tag>
    <el-tag type="success">平均评分: {{ averageScore.toFixed(1) }}</el-tag>
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
            <el-row :gutter="5" align="top">
              <el-col :span="6">
                <el-avatar :size="72" :src="item.notes_owner.avatar" fit="fill" />
                <el-tag class="ml-2" type="info">{{ item.notes_owner.name }}</el-tag>
              </el-col>
              <el-col :span="18">
                <el-alert v-if="item.notes_owner.slogan" :title="item.notes_owner.slogan" type="success"
                  :closable="false" />
                <h3 class="note-title">{{ item.note_title }}</h3>
                <el-tag class="ml-2" type="success" v-if="item.notes_count?.like_count" round>
                  👍🏻 {{ item.notes_count?.like_count }}
                </el-tag>
              </el-col>
            </el-row>
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
.overview {
  margin: 20px;
  .el-tag {
    margin-right: 10px;
  }
}

.waterfall-container {
  display: flex;
  gap: 20px;
  padding: 20px;
  min-height: 100vh;
  
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
  
  .card-header,
  .card-content {
    text-align: left;
  }
}

.note-title {
  margin: 10px 0;
  font-size: 16px;
  color: #303133;
}

.el-tag {
  margin-right: 5px;
  text-align: center;
}

.card-content {
  color: #606266;
  line-height: 1.6;
  
  p {
    margin: 8px 0;
  }
  
  ul, ol {
    padding-left: 20px;
    margin: 8px 0;
  }
}

.el-avatar {
  margin-bottom: 8px;
  display: block;
  margin-left: auto;
  margin-right: auto;
}

.el-alert {
  margin-bottom: 10px;
}

.el-card {
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
  transition: box-shadow 0.3s ease;
  
  &:hover {
    box-shadow: 0 4px 16px 0 rgba(0,0,0,0.2);
  }
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: #909399;
  
  .el-icon {
    margin-right: 8px;
    font-size: 20px;
  }
}

.no-more {
  text-align: center;
  padding: 20px;
  color: #909399;
  font-size: 14px;
}
</style>