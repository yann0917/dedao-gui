<template>
    <el-dialog v-model="dialogVisible" width="75%" :before-close="closeDialog">
        <template #header="{ titleId, titleClass }">
            <div class="my-header">
                <h4 :id="titleId" :class="titleClass">{{ courseInfo.class_info.name }}</h4>
                <el-alert :title="courseInfo.class_info.intro" type="info" :closable="false" center />
            </div>
        </template>
        <el-space wrap>
            <el-card v-for="i in 1" :key="i" class="box-card" width="100%">
                <template #header>
                    <div class="card-header">
                        <el-row :gutter="5" align="middle">
                            <el-col :span="4">
                                <el-row justify="center">
                                    <el-avatar :size="84" :src="courseInfo.class_info.lecturer_avatar" fit="fill" />
                                </el-row>
                                <el-row justify="center">
                                    <el-tag type="success">{{ courseInfo.class_info.lecturer_name }}</el-tag>
                                </el-row>
                            </el-col>
                            <el-col :span="18">
                                <p style="text-align:left"
                                    v-html="courseInfo.class_info?.lecturer_intro.replaceAll('\n', '<br/>')"></p>
                            </el-col>
                        </el-row>
                    </div>
                </template>

                <div class="book-tag" style="text-align:left">
                    <el-tag v-if="courseInfo.class_info.is_finished">
                        共{{ courseInfo.class_info.phase_num }}{{ courseInfo.class_info.price_desc }}
                    </el-tag>
                    <el-tag v-else>
                        共{{ courseInfo.class_info.phase_num }}{{ courseInfo.class_info.price_desc }}·已更新{{
                                courseInfo.class_info.current_article_count
                        }}{{ courseInfo.class_info.price_desc }}
                    </el-tag>
                    <el-tag class="ml-2" type="success">{{ courseInfo.class_info.learn_user_count }}人加入学习</el-tag>
                    <el-tag class="ml-2" type="info">{{ courseInfo.class_info.collection.collection_count }}人收藏</el-tag>
                    <el-tag class="ml-2" type="warning" v-if="courseInfo.class_info.collection.is_collected">已收藏
                    </el-tag>
                    <el-tag class="ml-2" type="warning" v-else>未收藏</el-tag>
                    <el-tag class="ml-2" type="danger">{{ courseInfo.class_comment_info?.count }}条评价</el-tag>
                </div>

                <div class="comment-info flex items-left" style="text-align: left;" v-if="courseInfo.class_comment_info?.comment_list.length >0">
                    <el-divider content-position="left">
                        <el-rate v-model="averageScore" disabled show-score :allow-half="true" text-color="var(--accent-color)" />
                    </el-divider>
                    <el-carousel height="200px" direction="vertical" loop>
                        <el-carousel-item v-for="item in courseInfo.class_comment_info.comment_list" :key="item.id">
                            <p class="text-large font-600 mr-3" v-if="item.title"> {{ item.title }} </p>
                            <p class="text-sm mr-2" style="color: var(--text-secondary)"
                                v-if="item.no_style_content">
                                <el-popover trigger="hover" placement="right" :width="300"
                                    :disabled="item.no_style_content.length <= 240" :content="item.no_style_content">
                                    <template #reference>
                                        <span slot="reference" v-if="item.no_style_content.length <= 240">{{
                                                item.no_style_content
                                        }}</span>
                                        <span slot="reference" v-if="item.no_style_content.length > 240">{{
                                                item.no_style_content.slice(0, 240)
                                                + "..."
                                        }}</span>
                                    </template>
                                </el-popover>
                            </p>
                            <el-avatar class="mr-3" :size="20" :src="item.avatar_s" /><span> &emsp;{{ item.nickname }}
                                评价了
                                <el-rate v-model.number="item.score" disabled show-score :allow-half="true"
                                    text-color="var(--accent-color)" />
                            </span>
                        </el-carousel-item>
                    </el-carousel>
                </div>

                <el-divider content-position="left">课程亮点</el-divider>
                <div style="text-align:left">
                    <span v-html="courseInfo.class_info.highlight?.replaceAll('\n', '<br/>')"></span>
                </div>

                <el-divider content-position="left">课程表</el-divider>
                <div class="outline_img" style="text-align:left">
                    <el-image :src="courseInfo.class_info.outline_img"></el-image>
                </div>
            </el-card>
        </el-space>

    </el-dialog>
</template>


<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { CourseInfo } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { useRouter } from 'vue-router'

const dialogVisible = ref(false)
let averageScore = ref(0)

let courseInfo = reactive(new services.CourseInfo)
courseInfo.class_info = new services.ClassInfo
courseInfo.intro_article = new services.ArticleIntro

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
    getCourseInfo(props.enid);
})

// // 监听 props 变化
// watch(() => props.dialogVisible, (newVal) => {
//   dialogVisible.value = newVal;
//   if (newVal && props.enid) {
//     getCourseInfo(props.enid);
//   }
// }, { immediate: true })

// // 监听 enid 变化
// watch(() => props.enid, (newVal) => {
//   if (newVal && props.dialogVisible) {
//     getCourseInfo(newVal);
//   }
// })

const resetCourseInfo = () => {
  courseInfo.class_info = new services.ClassInfo;
  courseInfo.intro_article = new services.ArticleIntro;
  averageScore.value = 0;
}

const closeDialog = () => {
  dialogVisible.value = false;
  resetCourseInfo();
  emits('close');
}

const getCourseInfo = async (enid: string) => {
  try {
    const info = await CourseInfo(enid);
    Object.assign(courseInfo, info);
    averageScore.value = Number(courseInfo.class_comment_info.average_score);
    dialogVisible.value = true;
  } catch (error) {
    console.error('获取课程信息失败:', error);
    ElMessage({
      message: typeof error === 'string' ? error : '获取课程信息失败',
      type: 'error'
    });
    closeDialog(); // 发生错误时关闭对话框并重置状态
  }
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
  color: var(--text-primary);
}

.book-tag {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.outline_img {
  width: 100%;
  overflow: hidden;
  
  .el-image {
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
  }
}

:deep(.el-dialog__header) {
  padding: 20px 20px 0;
}

:deep(.el-alert) {
  margin-bottom: 8px;
}

/* 暗色主题适配 */
.theme-dark .my-header h4 {
  color: var(--text-primary) !important;
}

.theme-dark .text-sm {
  color: var(--text-secondary) !important;
}

.theme-dark :deep(.el-rate__text) {
  color: var(--accent-color) !important;
}

.theme-dark :deep(.el-dialog) {
  background-color: var(--card-bg) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-dialog__header) {
  background-color: var(--card-bg) !important;
  color: var(--text-primary) !important;
  border-bottom-color: var(--border-soft) !important;
}

.theme-dark :deep(.el-dialog__body) {
  background-color: var(--card-bg) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-card) {
  background-color: var(--card-bg) !important;
  border-color: var(--border-soft) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-card__header) {
  background-color: var(--card-bg) !important;
  border-bottom-color: var(--border-soft) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-card__body) {
  background-color: var(--card-bg) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-divider__text) {
  background-color: var(--card-bg) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-carousel__item) {
  background-color: var(--card-bg) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-popover) {
  background-color: var(--card-bg) !important;
  border-color: var(--border-soft) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-popover__content) {
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-tag) {
  background-color: var(--fill-color-light) !important;
  border-color: var(--border-soft) !important;
  color: var(--text-primary) !important;
}

.theme-dark :deep(.el-tag--success) {
  background-color: var(--success-color-light) !important;
  color: var(--success-color) !important;
  border-color: var(--success-color) !important;
}

.theme-dark :deep(.el-tag--info) {
  background-color: var(--info-color-light) !important;
  color: var(--info-color) !important;
  border-color: var(--info-color) !important;
}

.theme-dark :deep(.el-tag--warning) {
  background-color: var(--warning-color-light) !important;
  color: var(--warning-color) !important;
  border-color: var(--warning-color) !important;
}

.theme-dark :deep(.el-tag--danger) {
  background-color: var(--danger-color-light) !important;
  color: var(--danger-color) !important;
  border-color: var(--danger-color) !important;
}
</style>