<template>
  <el-dialog v-model="dialogVisible" width="75%" :before-close="closeDialog">
    <template #header="{titleId, titleClass }">
      <div class="my-header">
        <h4 :id="titleId" :class="titleClass">{{audioInfo.detail.title}}</h4>
        <el-alert :title="audioInfo.detail.audio_summary" type="info" :closable="false" center />
      </div>
    </template>
    <el-space wrap>
      <el-card v-for="i in 1" :key="i" class="box-card" width="100%">
        <template #header>
          <div class="card-header">
            <el-row :gutter="0" align="middle">
              <el-col :span="6">
                <el-image :src="audioInfo.detail.icon" fit="cover" />
              </el-col>
              <el-col :span="18" style="text-align:left; padding-left: 15px;">
                <el-row :gutter="0" align="middle">
                  <el-col :span="2">
                    <el-avatar :src="audioInfo.detail.agency_detail.qcg_member_avatar"></el-avatar>
                  </el-col>
                  <el-col :span="16">
                    <p style="margin: 0;font-size: medium;">{{ audioInfo.detail.agency_detail.qcg_member_name }}</p>
                    <span style="color: #909399;">解读书本{{ audioInfo.detail.agency_detail.book_count }}本 | 被学习{{ audioInfo.detail.agency_detail.uv }}次</span>
                  </el-col>
                </el-row>
                <el-row :gutter="0" >
                  <p class="author-info" v-html="audioInfo.detail.agency_detail.intro?.replaceAll('\n','<br/>')"></p><br />
                  <span style="color:#ff6b00">{{ audioInfo.detail.audio_price }}元</span>
                  <el-tag class="ml-2" type="warning" v-if="audioInfo.detail.is_vip == true" round>
                    <el-icon><HotWater /></el-icon>会员免费</el-tag>
                  <el-tag class="ml-2" type="info" v-if="audioInfo.detail.in_bookrack == true" round>
                    <el-icon><HotWater /></el-icon>已加入书架</el-tag>
                  <el-tag class="ml-2" type="success" v-if="audioInfo.detail.tag?.length >0" round>
                    {{ audioInfo.detail.tag.join(",") }}</el-tag>
                </el-row>

              </el-col>
            </el-row>
          </div>
        </template>
        <div class="section-info ebook college" >
          <!-- <div class="item-content" v-if="audioInfo.detail.rank.product_score !='0.0'">
            <span class="info-value">{{ audioInfo.detail.product_score }}</span>
            <span class="info-text">用户推荐指数</span>
          </div> -->
          <div class="item-content" v-if="audioInfo.detail.rank.rank_number !=0">
            <span class="info-value"> No.{{ audioInfo.detail.rank.rank_number }}</span>
            <span class="info-text">{{ audioInfo.detail.rank.rank_name }}</span>
          </div>
          <div class="item-content">
            <span class="info-value">{{ audioInfo.detail.learn_count_desc }}</span>
            <span class="info-text">学习次数</span>
          </div>
          <div class="item-content" v-if="audioInfo.detail.duration !=0">
            <span class="info-value">{{ secondToHour(audioInfo.detail.duration) }}</span>
            <span class="info-text">音频时长</span>
          </div>
          <div class="item-content">
            <span class="info-value">{{timestampToDate(audioInfo.detail.publish_time)}}</span>
            <span class="info-text">上架时间</span>
          </div>
        </div>
        <el-divider content-position="left">音频介绍</el-divider>
        <div class="book-intro" style="text-align:left">
          <div v-for="item in audioInfo.detail.topic_summary" :key="item.id">
              <h3>{{ item.title}}</h3>
              <p v-html="item.sub_title?.replaceAll('\n','<br/>')" style="color: #909399;"></p>
          </div>
        </div>
      </el-card>
    </el-space>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { AudioDetail } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { repeat } from 'lodash'
import { secondToHour,timestampToDate } from '../utils/utils'
import { useRouter } from 'vue-router'

const router = useRouter()
const dialogVisible = ref(false)

let audioInfo = reactive(new services.AudioInfoResp)

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
  getAudioInfo(props.enid);
})

const getAudioInfo = async (enid: string) => {
  await AudioDetail(enid).then((info) => {
    console.log(info)
    Object.assign(audioInfo, info)
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
  audioInfo = reactive(new services.AudioInfoResp)
  emits('close')
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

.card-header .el-button {
    /* height: 32px; */
    border-radius: 8px;
    font-weight: bold;
    border-color: #ff6b00;
    background-color: #ff6b00;
    color: #fff;
}
</style>