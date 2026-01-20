<template>
  <el-dialog
    v-model="dialogVisible"
    width="75%"
    :before-close="closeDialog"
    class="info-dialog ebook-info-dialog"
    destroy-on-close
  >
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
              <el-col :span="18" class="card-main">
                <p class="author-info" v-html="ebookInfo.author_info?.replaceAll('\n','<br/>')"></p><br />
                <span class="price">{{ ebookInfo.price }}元</span>
                <el-tag class="ml-2" type="warning" v-if="ebookInfo.is_vip_book==1" round>
                  会员免费</el-tag>
                <el-tag class="ml-2" type="info" v-if="ebookInfo.read_time">阅读总时长：{{secondToHour(ebookInfo.read_time)}}</el-tag>
                <el-button v-if="ebookInfo.is_on_bookshelf==false" class="primary-action" @click="ebookShelfAdd(ebookInfo.enid)">
                  <el-icon><Plus /></el-icon>加入书架
                </el-button>
                <el-button v-if="ebookInfo.is_on_bookshelf==true" @click="ebookShelfRemove(ebookInfo.enid)" type="danger">
                  <el-icon><Delete /></el-icon>移出书架
                </el-button>
              </el-col>
            </el-row>
          </div>
        </template>
        <div class="section-info ebook college" >
          <div class="item-content" v-if="ebookInfo.product_score !='0.0'">
            <span class="info-value">{{ ebookInfo.product_score }}</span>
            <span class="info-text">用户推荐指数</span>
          </div>
          <div class="item-content" v-if="ebookInfo.classify_name !=''">
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
            <span class="info-value">{{Math.floor(ebookInfo.count/1000)}}千字</span>
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
              <span :style="[item.level == 0 ? 'font-size:large;font-weight:bold':'']">{{repeat("&nbsp;&nbsp;", item.level)+ item.text }}<br/></span>
            </span>
          </el-scrollbar>
        </div>
      </el-card>
    </el-space>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { EbookInfo, EbookShelfAdd, EbookShelfRemove } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { repeat } from 'lodash'
import { secondToHour } from '../utils/utils'
import { useAppRouter } from '../composables/useRouter'
import { Plus, Delete } from '@element-plus/icons-vue'

const { pushEbookComment } = useAppRouter()
const dialogVisible = ref(false)

let ebookInfo = reactive(new services.EbookDetail)

const emits = defineEmits(['close', 'bookshelf-changed']);

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

const ebookShelfAdd = async (enid: string) => {
  await EbookShelfAdd([enid]).then((info) => {
    console.log(info)
    ElMessage({
      message: '已加入书架',
      type: 'success'
    })
    getEbookInfo(enid)
    // 发射事件通知父组件书架状态已改变
    emits('bookshelf-changed')
  }).catch((error) => {
    ElMessage({
      message: error,
      type: 'warning'
    })
  })
  return
}

const ebookShelfRemove = async (enid: string) => {
  ElMessageBox.confirm(
    '是否移出书架?',
    '',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(() => {
      EbookShelfRemove([enid]).then((info) => {
        console.log(info)
        ElMessage({
          message: '已移出书架',
          type: 'success'
        })
        getEbookInfo(enid)
        // 发射事件通知父组件书架状态已改变
        emits('bookshelf-changed')
      }).catch((error) => {
        ElMessage({
          message: error,
          type: 'warning'
        })
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消',
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
    pushEbookComment({
        id: row.id,
        enid:row.enid,
        total: row.publish_num,
        title: row.title
    })
}
</script>

<style scoped lang="scss">
.ebook-info-dialog {
  :deep(.el-dialog__header) {
    margin-right: 0;
    padding: 20px 24px;
    border-bottom: 1px solid var(--border-soft);
  }

  :deep(.el-dialog__body) {
    padding: 0;
  }
}

.my-header {
  display: flex;
  flex-direction: column;
  gap: 12px;

  h4 {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    color: var(--text-primary);
    line-height: 1.4;
  }
}

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
  text-align: center;
}

.card-main {
  text-align: left;
  padding-left: 15px;
}

.price {
  color: var(--accent-color) !important;
  font-weight: 600;
}

.section-info {
  display: flex;
  flex-wrap: wrap;
  margin: 20px 0;
}

.section-info .item-content {
  position: relative;
  text-align: center;
  min-width: 24%;
  margin-bottom: 15px;
  border-left: 1px solid var(--border-soft);

  &:first-child {
    border-left: none;
  }
}

.item-content span {
  display: block;
}

.item-content .info-value {
  font-size: 16px;
  color: var(--text-primary);
  letter-spacing: 0;
  line-height: 28px;
  font-weight: 500;
}

.item-content .info-text {
  color: var(--text-secondary);
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

.card-header .primary-action {
  border-radius: 8px;
  font-weight: bold;
  border-color: var(--accent-color);
  background-color: var(--accent-color);
  color: #fff !important;
  transition: all 0.3s ease;
  
  &:hover {
    background-color: var(--accent-hover);
    border-color: var(--accent-hover);
    transform: translateY(-2px);
  }
}

.book-intro {
  h1 {
    font-size: 18px;
    color: var(--text-primary) !important;
    margin: 16px 0 8px;
    
    &:first-child {
      margin-top: 0;
    }
  }
  
  p {
    color: var(--text-secondary) !important;
    margin: 0 0 16px;
    line-height: 1.8;
  }
}

:deep(.el-alert) {
  margin: 16px 0;
}

:deep(.el-scrollbar) {
  .scrollbar-catalog-item {
    display: block;
    padding: 4px 0;
    color: var(--text-secondary);
    
    &:hover {
      background-color: var(--fill-color);
    }
  }
}
</style>
