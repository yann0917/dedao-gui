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
            <el-row :gutter="0" align="middle">
              <el-col :span="6">
                <el-image :src="ebookInfo.cover" fit="cover" />
              </el-col>
              <el-col :span="18" style="text-align:left; padding-left: 15px;">
                <p class="author-info" v-html="ebookInfo.author_info?.replaceAll('\n','<br/>')"></p><br />
                <span style="color:#ff6b00">{{ ebookInfo.price }}元</span>
                <el-tag class="ml-2" type="warning" v-if="ebookInfo.is_vip_book==1" round>
                  <el-icon><HotWater /></el-icon>会员免费</el-tag>
                  <el-tag class="ml-2" type="info" v-if="ebookInfo.read_time">阅读总时长：{{secondToHour(ebookInfo.read_time)}}</el-tag>
              </el-col>
            </el-row>
          </div>
        </template>
        <div class="section-info ebook college" >
          <div class="item-content" v-if="ebookInfo.product_score !='0.0'">
            <span class="info-value">{{ ebookInfo.product_score }}</span>
            <span class="info-text">用户推荐指数</span>
          </div>
          <div class="item-content">
            <span class="info-value">{{ ebookInfo.classify_name }}</span>
            <span class="info-text">类型</span>
          </div>
          <div class="item-content" v-if="ebookInfo.douban_score !=''">
            <span class="info-value">{{ebookInfo.douban_score}}</span>
            <span class="info-text douban">豆瓣评分</span>
          </div>
          <div class="item-content">
            <span class="info-value"> {{ebookInfo.is_tts_switch==true ? '可以朗读':'不可朗读'}} </span>
            <span class="info-text">语音朗读</span>
          </div>
          <div class="item-content">
            <span class="info-value">{{Math.ceil(ebookInfo.count/1000)}}千字</span>
            <span class="info-text">字数</span>
          </div>
          <div class="item-content" v-if="ebookInfo.rank_num !=0">
            <span class="info-value"> No.{{ ebookInfo.rank_num }}</span>
            <span class="info-text">{{ ebookInfo.rank_name }}</span>
          </div>
          <div class="item-content">
            <span class="info-value">{{ebookInfo.publish_time}}</span>
            <span class="info-text">发行日期</span>
          </div>
        </div>
        <el-divider content-position="left">{{ebookInfo.press.name}}</el-divider>
        <el-alert v-html="ebookInfo.press.brief?.replaceAll('\n','<br/>')" type="info" :closable="false" align="left" />

        <el-divider content-position="left">简 介</el-divider>
        <div class="book-intro" style="text-align:left">
          <h1>主编推荐语</h1>
          <p>{{ebookInfo.other_share_summary}}</p>
          <h1>内容简介</h1>
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
const dialogVisible = ref(false)

let ebookInfo = reactive(new services.EbookDetail)

const emits  = defineEmits(['close']);

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
}
const closeDialog = () => {
  ebookInfo = reactive(new services.EbookDetail)
  emits('close')
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
.section-info{
  display: flex;
  flex-wrap: wrap;
}
.section-info .item-content {
  position: relative;
  text-align: center;
  min-width: 24%;
  margin-bottom: 15px;
  border-left: 1px solid #d8d8d8;
}
.item-content span {
  display: block;
}

.item-content .info-value {
  font-size: 16px;
  color: #333;
  letter-spacing: 0;
  line-height: 28px;
}

.item-content .info-text {
  color: #666;
  line-height: 16px;
  font-size: 12px;
  letter-spacing: 0;
}

.douban {
  background: url(https://piccdn2.umiwi.com/fe-oss/default/doubanlogo.png) no-repeat;
  background-size: 20px;
  padding-left: 20px;
  background-position: 35%;
}
</style>