<template>
  <el-breadcrumb :separator-icon="ArrowRight">
    <el-breadcrumb-item :to="{ path: '/ebook' }">电子书架</el-breadcrumb-item>
    <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
  </el-breadcrumb>
 
  <el-space wrap alignment="flex-start" class="card">
    <el-card v-for="item in tableData.list" :key="item.note_id" class="box-card" shadow="hover" style="width: 400px;">
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
              <p>{{ item.note_title }}</p>
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
  </el-space>

  <Pagination :total="total" @pageChange="handleChangePage"></Pagination>

</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowRight } from '@element-plus/icons-vue'
import { EbookCommentList } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import Pagination from '../components/Pagination.vue'
import { useRoute, useRouter } from 'vue-router'
import { userStore } from '../stores/user';

const store = userStore()
const router = useRouter()
const route = useRoute()
const loading = ref(true)
const page = ref(1)
const total = ref(0)
const pageSize = ref(15)

let averageScore = ref(0)
let scoreInfo = ref()

let tableData = reactive(new services.EbookCommentList)


let id = ref()
let enid = ref()
let breadcrumbTitle = ref()

onMounted(() => {
  watch(() => {
    id.value = route.query.id
    enid.value = route.query.enid
    breadcrumbTitle.value = route.query.title
  },
    () => getTableData(),
    { immediate: true }
  )
})

// 分页
const handleChangePage = (item: any) => {
  page.value = item.page
  pageSize.value = item.pageSize
  getTableData()
}


const getTableData = async () => {
  await EbookCommentList(enid.value, page.value, pageSize.value).then((table) => {
    loading.value = false
    console.log(table)
    Object.assign(tableData, table)
    total.value = tableData.total
    averageScore.value = Number(tableData.ebook_score.average_score)
    scoreInfo.value =tableData.ebook_score.score_info

  }).catch((error) => {
    loading.value = false
    ElMessage({
      message: error,
      type: 'warning'
    })
  })
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
.el-card{
  color:#606266;
  .card-header,
  .card-content {
    text-align: left ;
  }
}

.el-tag {
  margin-right: 5px;
  text-align: center;
}

.card {
  display:flex;
  column-gap: 5px;
}

.card .box-card {
  margin-bottom: 10px;
}

</style>