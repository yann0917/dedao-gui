<template>
  <el-dialog v-model="dialogVisible" width="60%" :before-close="closeDialog">
    <template #header="{titleId, titleClass }">
      <div class="my-header">
        <h4 :id="titleId" :class="titleClass">{{ebookInfo.operating_title}}</h4>
        <el-alert :title="ebookInfo.author_list.join(' ')+' 著/'+ebookInfo.press.name" type="info" :closable="false" center />
      </div>
    </template>
    <el-space wrap>
      <el-card v-for="i in 1" :key="i" class="box-card" width="100%">
        <template #header>
          <div class="card-header">
            <el-row :gutter="5" align="middle">
              <el-col :span="6">
                <el-row justify="center">
                  <el-image :src="ebookInfo.cover" fit="cover" />
                </el-row>
                <el-row justify="center">
                  <el-tag type="success">{{ebookInfo.book_author}}</el-tag>
                  <el-tag type="success" v-if="ebookInfo.classify_name">{{ebookInfo.classify_name}}</el-tag>
                </el-row>
              </el-col>
              <el-col :span="18" style="text-align:left">
                <p class="author-info" v-html="ebookInfo.author_info?.replaceAll('\n','<br/>')"></p><br />
                <el-tag class="ml-2" type="warning" v-if="ebookInfo.is_vip_book==1" round>
                  <el-icon><HotWater /></el-icon>会员免费</el-tag>
                <el-tag class="ml-2" type="danger" v-if="ebookInfo.is_tts_switch==true" round>
                  <el-icon><Microphone /></el-icon>可朗读
                </el-tag>
                <el-tag class="ml-2" type="danger" v-else round>
                  <el-icon><Mute /></el-icon>不可朗读
                </el-tag>
              </el-col>
            </el-row>
          </div>
        </template>
        <div class="book-tag" style="text-align:left">
          <el-tag  class="ml-2" type="info">出版日期：{{ebookInfo.publish_time}}</el-tag>
          <el-tag class="ml-2" type="success" v-if="ebookInfo.douban_score">豆瓣评分：{{ebookInfo.douban_score}}</el-tag>
          <el-tag class="ml-2" type="warning" v-if="ebookInfo.count">字数：{{Math.ceil(ebookInfo.count/1000)}}千字</el-tag>
          <el-tag class="ml-2" type="info" v-if="ebookInfo.read_time">阅读总时长：{{secondToHour(ebookInfo.read_time)}}</el-tag>
        </div>
        <el-divider content-position="left">{{ebookInfo.press.name}}</el-divider>
        <el-alert v-html="ebookInfo.press.brief?.replaceAll('\n','<br/>')" type="info" :closable="false" align="left" />

        <el-divider content-position="left">简 介</el-divider>
        <div class="book-intro" style="text-align:left">
          <p>{{ebookInfo.other_share_summary}}</p>
          <p v-html="ebookInfo.book_intro?.replaceAll('\n','<br/>')"></p>
        </div>
        <el-divider content-position="left">目 录</el-divider>
        <div class="catalog">
          <el-scrollbar height="320px" align="left">
            <span v-for="item in ebookInfo.catalog_list" :key="item.playOrder" class="scrollbar-catalog-item">
            {{repeat("&nbsp;&nbsp;", item.level)+ item.text }}<br/></span>
          </el-scrollbar>
        </div>
      </el-card>
    </el-space>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, defineEmits, defineProps,inject } from 'vue'
import { ElTable, ElMessage } from 'element-plus'
import { EbookInfo } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { repeat } from 'lodash'
import { secondToHour } from '../utils/utils'
import { useRouter } from 'vue-router'

const router = useRouter()
const dialogVisible = ref(true)

let ebookInfo = reactive(new services.EbookDetail)

const emits  = defineEmits(['showInfo']);

const props = defineProps({
        enid: {
            type: String,
            default: ""
        },
        dialogVisible:{
            type: Boolean,
            default:false
        }
    })

onMounted(() => {
getEbookInfo(props.enid);
})

const getEbookInfo = async (enid: string) => {
  await EbookInfo(enid).then((info) => {
    console.log(info)
    Object.assign(ebookInfo, info)
    openDialog()
  }).catch((error) => {
    ElMessage({
      message: error,
      type: 'warning'
    })
  })
  return
}

const openDialog = () => {
    dialogVisible.value = props.dialogVisible
    console.log(props.dialogVisible)
}
const closeDialog = () => {
  dialogVisible.value = !props.dialogVisible
}

const gotoCommentList = (row: any) => {
    router.push({
        path: `/ebook/comment`,
        query: {
            id: row.id,
            enid:row.enid,
            total: row.publish_num,
            title: row.title
        }
    })
}

</script>

<style scoped>
.card-header.el-row {
  height: 100%;
}

.el-row {
  margin-bottom: 10px;
}

.el-row:last-child {
  margin-bottom: 0;
}

.el-col {
  border-radius: 4px;
}

.divider-text {
  font-size: larger;
}

.catalog,
.author-info,
.book-intro {
  line-height: 1.6;
}
.el-tag {
  margin-right: 5px;
    /* height: auto; */
  text-align: center;
}
</style>