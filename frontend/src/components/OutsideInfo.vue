<template>
    <el-dialog v-model="dialogVisible" width="75%" :before-close="closeDialog">
        <template #header="{ titleId, titleClass }">
            <div class="my-header">
                <h4 :id="titleId" :class="titleClass">{{ outsideInfo.spu.title }}</h4>
            </div>
        </template>
        <el-space wrap>
            <el-card v-for="i in 1" :key="i" class="box-card" width="100%">
                <template #header>
                    <div class="card-header">
                        <el-row :gutter="5" align="middle">
                            <el-col :span="4">
                                <el-row justify="center">
                                    <el-avatar :size="84" :src="outsideInfo.spu.extra.teacher_avatar" fit="fill" />
                                </el-row>
                                <el-row justify="center">
                                    <el-tag type="success">{{ outsideInfo.spu.extra.teacher_name }}</el-tag>
                                </el-row>
                            </el-col>
                            <el-col :span="18">
                                <p style="text-align:left"
                                    v-html="outsideInfo.spu.extra.teacher_intro?.replaceAll('\n', '<br/>')"></p>
                            </el-col>
                        </el-row>
                    </div>
                </template>

                <div class="book-tag" style="text-align:left">
                    <el-tag >
                        共{{ outsideInfo.count }}期
                    </el-tag>
                   
                    <el-tag class="ml-2" type="success" v-if="outsideInfo.spu.extra?.odob_consumer_num">{{ outsideInfo.spu.extra.odob_consumer_num }}人加入学习</el-tag>
                    <el-tag class="ml-2" type="info" v-if="outsideInfo.spu.extra?.rn_learn_count">{{ outsideInfo.spu.extra.rn_learn_count }}人学习</el-tag>
                    <el-tag class="ml-2" type="warning" v-if="outsideInfo.spu.extra?.book_shelf_status">已收藏</el-tag>
                    <el-tag class="ml-2" type="warning" v-else>未收藏</el-tag>
                </div>

                <el-divider content-position="left">课程亮点</el-divider>
                <div style="text-align:left">
                    <span v-html="outsideInfo.spu.extra?.intro_text?.replaceAll('\n', '<br/>')"></span>
                </div>

                <el-divider content-position="left">课程介绍</el-divider>
                <div class="course-intro" style="text-align:left">
                    <div v-for="(paragraph, index) in parsedIntro" :key="index" class="intro-paragraph">
                        <p v-if="paragraph.text && paragraph.text.trim()" v-html="paragraph.text.replaceAll('\n', '<br/>')"></p>
                    </div>
                </div>

                <el-divider content-position="left">课程列表</el-divider>
                <div class="course-list" style="text-align:left">
                    <div v-for="(item, index) in outsideInfo.items" :key="index" class="course-item">
                        <el-row :gutter="16" align="middle" style="margin-bottom: 12px; padding: 8px; border: 1px solid #f0f0f0; border-radius: 8px;">
                            <el-col :span="4">
                                <el-image :src="item.icon" fit="cover" style="width: 60px; height: 80px; border-radius: 4px;" />
                            </el-col>
                            <el-col :span="16">
                                <h5 style="margin: 0 0 8px 0; font-size: 14px; color: #333;">{{ item.title }}</h5>
                                <p style="margin: 0; font-size: 12px; color: #666; line-height: 1.4;">{{ item.summary }}</p>
                                <div style="margin-top: 8px; display: flex; gap: 12px; font-size: 11px; color: #999;">
                                    <span v-if="item.extra?.publish_time">
                                        <el-icon><Clock /></el-icon>
                                        {{ formatPublishTime(item.extra.publish_time) }}
                                    </span>
                                    <span v-if="item.extra?.duration">
                                        <el-icon><VideoPlay /></el-icon>
                                        {{ formatDuration(item.extra.duration) }}
                                    </span>
                                </div>
                            </el-col>
                            <el-col :span="4" style="text-align: right;">
                                <el-tag v-if="item.extra?.progress_learned" type="success" size="small">已学习</el-tag>
                                <el-tag v-else type="info" size="small">未学习</el-tag>
                            </el-col>
                        </el-row>
                    </div>
                </div>
            </el-card>
        </el-space>

    </el-dialog>
</template>


<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Clock, VideoPlay } from '@element-plus/icons-vue'
import { OutsideDetail } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'

const dialogVisible = ref(false)

let outsideInfo = reactive(new services.OutsideDetail)

// 解析课程介绍JSON
const parsedIntro = computed(() => {
  if (!outsideInfo.spu.intro) return [];
  
  try {
    const introData = JSON.parse(outsideInfo.spu.intro);
    if (Array.isArray(introData)) {
      return introData.filter(item => item.type === 'paragraph' && item.text && item.text.trim());
    }
  } catch (error) {
    console.warn('解析课程介绍JSON失败:', error);
  }
  
  return [];
})

const emits = defineEmits(['close']);

const props = defineProps({
  enid: {
    type: String,
    default: ""
  },
  dialogVisible: {
    type: Boolean,
    default: false
  }
})

onMounted(() => {
    getOutsideInfo(props.enid);
})

const resetOutsideInfo = () => {
  Object.assign(outsideInfo, new services.OutsideDetail);
}

const closeDialog = () => {
  dialogVisible.value = false;
  resetOutsideInfo();
  emits('close');
}

const getOutsideInfo = async (enid: string) => {
  try {
    const info = await OutsideDetail(enid);
    Object.assign(outsideInfo, info);
    dialogVisible.value = true;
    console.log(outsideInfo);
  } catch (error) {
    console.error('获取名家讲书信息失败:', error);
    ElMessage({
      message: typeof error === 'string' ? error : '获取名家讲书信息失败',
      type: 'error'
    });
    closeDialog(); // 发生错误时关闭对话框并重置状态
  }
}

// 格式化发布时间
const formatPublishTime = (publishTime: number) => {
  const date = new Date(publishTime * 1000);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  
  if (days === 0) {
    return '今天发布';
  } else if (days === 1) {
    return '昨天发布';
  } else if (days < 7) {
    return `${days}天前发布`;
  } else if (days < 30) {
    const weeks = Math.floor(days / 7);
    return `${weeks}周前发布`;
  } else {
    return date.toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    });
  }
}

// 格式化音频时长（秒转换为分秒）
const formatDuration = (duration: number) => {
  const minutes = Math.floor(duration / 60);
  const seconds = duration % 60;
  return `${minutes}分${seconds}秒`;
}

</script>

<style scoped>
.my-header {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.my-header h4 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.book-tag {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.course-list {
  max-height: 400px;
  overflow-y: auto;
}

.course-item {
  margin-bottom: 8px;
}

.course-item:last-child {
  margin-bottom: 0;
}

:deep(.el-dialog__header) {
  padding: 20px 20px 0;
}

:deep(.el-alert) {
  margin-bottom: 8px;
}

.course-item .el-icon {
  margin-right: 4px;
  font-size: 12px;
}

.course-intro {
  line-height: 1.6;
  color: #333;
}

.intro-paragraph {
  margin-bottom: 12px;
}

.intro-paragraph p {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
}

.intro-paragraph p:last-child {
  margin-bottom: 0;
}
</style>
