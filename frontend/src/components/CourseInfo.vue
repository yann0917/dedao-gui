<template>
    <el-dialog 
        v-model="dialogVisible" 
        width="800px" 
        :before-close="closeDialog"
        class="course-info-dialog"
        destroy-on-close
    >
        <template #header>
            <div class="dialog-header-content">
                <h3>{{ courseInfo.class_info.name }}</h3>
                <p class="course-intro" v-if="courseInfo.class_info.intro">{{ courseInfo.class_info.intro }}</p>
            </div>
        </template>
        
        <div class="course-detail-container">
            <!-- 讲师信息卡片 -->
            <div class="lecturer-section">
                <div class="lecturer-card">
                    <div class="lecturer-avatar">
                        <el-avatar :size="80" :src="courseInfo.class_info.lecturer_avatar" />
                        <el-tag class="lecturer-tag" effect="dark" type="success" size="small">主讲人</el-tag>
                    </div>
                    <div class="lecturer-info">
                        <div class="lecturer-name">{{ courseInfo.class_info.lecturer_name }}</div>
                        <div class="lecturer-desc" v-html="formattedLecturerIntro"></div>
                    </div>
                </div>
            </div>

            <!-- 数据统计标签 -->
            <div class="stats-container">
                <div class="stat-item">
                    <el-icon class="stat-icon"><Collection /></el-icon>
                    <span class="stat-text">{{ updateStatusText }}</span>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item">
                    <el-icon class="stat-icon"><User /></el-icon>
                    <span class="stat-text">{{ courseInfo.class_info.learn_user_count }}人加入学习</span>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item clickable" @click="toggleCollection">
                    <el-icon class="stat-icon" :class="{ active: courseInfo.class_info?.collection?.is_collected }">
                        <StarFilled v-if="courseInfo.class_info?.collection?.is_collected" />
                        <Star v-else />
                    </el-icon>
                    <span class="stat-text">
                        {{ courseInfo.class_info?.collection?.is_collected ? '已收藏' : '收藏' }}
                        <span class="stat-count" v-if="courseInfo.class_info?.collection?.collection_count">
                            ({{ courseInfo.class_info?.collection?.collection_count }})
                        </span>
                    </span>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item">
                    <el-icon class="stat-icon"><ChatLineSquare /></el-icon>
                    <span class="stat-text">{{ courseInfo.class_comment_info?.count || 0 }}条评价</span>
                </div>
            </div>

            <!-- 评论轮播 -->
            <div class="comments-section" v-if="courseInfo.class_comment_info?.comment_list?.length > 0">
                <div class="section-title">
                    <span class="label">学员评价</span>
                    <el-rate 
                        v-model="averageScore" 
                        disabled 
                        show-score 
                        text-color="#ff9900" 
                        score-template="{value}分"
                    />
                </div>
                <div class="comments-carousel-wrapper">
                    <el-carousel 
                        height="160px" 
                        direction="vertical" 
                        :interval="3000" 
                        :autoplay="true"
                        indicator-position="none"
                        pause-on-hover
                    >
                        <el-carousel-item v-for="item in courseInfo.class_comment_info.comment_list" :key="item.id">
                            <div class="comment-card">
                                <div class="comment-header">
                                    <div class="user-info">
                                        <el-avatar :size="24" :src="item.avatar_s" />
                                        <span class="username">{{ item.nickname }}</span>
                                    </div>
                                    <el-rate v-model.number="item.score" disabled size="small" />
                                </div>
                                <div class="comment-content">
                                    <h5 v-if="item.title">{{ item.title }}</h5>
                                    <p class="comment-text" :title="item.no_style_content">{{ item.no_style_content }}</p>
                                </div>
                            </div>
                        </el-carousel-item>
                    </el-carousel>
                </div>
            </div>

            <!-- 课程亮点 -->
            <div class="highlights-section" v-if="courseInfo.class_info.highlight">
                <div class="section-title">课程亮点</div>
                <div class="highlight-content" v-html="formattedHighlight"></div>
            </div>

            <!-- 课程大纲 -->
            <div class="outline-section" v-if="courseInfo.class_info.outline_img">
                <div class="section-title">课程大纲</div>
                <div class="outline-image-wrapper">
                    <el-image 
                        :src="courseInfo.class_info.outline_img" 
                        :preview-src-list="[courseInfo.class_info.outline_img]"
                        fit="contain"
                        loading="lazy"
                    >
                        <template #placeholder>
                            <div class="image-placeholder">加载中...</div>
                        </template>
                    </el-image>
                </div>
            </div>
        </div>
    </el-dialog>
</template>


<script lang="ts" setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { CourseInfo as GetCourseInfo } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { Collection, User, Star, StarFilled, ChatLineSquare } from '@element-plus/icons-vue'

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

const emits = defineEmits(['close']);

const dialogVisible = ref(props.dialogVisible)
let averageScore = ref(0)

let courseInfo = reactive(new services.CourseInfo)
courseInfo.class_info = new services.ClassInfo
courseInfo.class_info.collection = new services.Collection
courseInfo.intro_article = new services.ArticleIntro

// 计算属性：格式化讲师简介
const formattedLecturerIntro = computed(() => {
    return courseInfo.class_info?.lecturer_intro?.replaceAll('\n', '<br/>') || ''
})

// 计算属性：格式化课程亮点
const formattedHighlight = computed(() => {
    return courseInfo.class_info?.highlight?.replaceAll('\n', '<br/>') || ''
})

// 计算属性：更新状态文本
const updateStatusText = computed(() => {
    const info = courseInfo.class_info
    // Ensure info exists and has properties before accessing
    if (!info) return ''
    
    const phaseNum = info.phase_num || 0
    const priceDesc = info.price_desc || ''
    const currentArticleCount = info.current_article_count || 0
    
    if (info.is_finished) {
        return `共${phaseNum}${priceDesc}`
    } else {
        return `共${phaseNum}${priceDesc} · 已更新${currentArticleCount}${priceDesc}`
    }
})

const resetCourseInfo = () => {
  courseInfo.class_info = new services.ClassInfo;
  courseInfo.class_info.collection = new services.Collection;
  courseInfo.intro_article = new services.ArticleIntro;
  averageScore.value = 0;
}

const closeDialog = () => {
  dialogVisible.value = false;
  // resetCourseInfo(); // 关闭时不需要立即重置，等下次打开覆盖即可，避免动画时内容闪烁
}

const toggleCollection = () => {
    // 暂时仅支持显示状态，不处理实际收藏逻辑
    // 如果需要实现收藏功能，需要调用后端API
    ElMessage.info('收藏功能开发中')
}

const getCourseInfo = async (enid: string) => {
  try {
    const info = await GetCourseInfo(enid);
    Object.assign(courseInfo, info);
    averageScore.value = Number(courseInfo.class_comment_info.average_score);
  } catch (error) {
    console.error('获取课程信息失败:', error);
    ElMessage({
      message: typeof error === 'string' ? error : '获取课程信息失败',
      type: 'error'
    });
    closeDialog(); 
  }
}

// 监听 props.dialogVisible 变化以同步内部状态
watch(() => props.dialogVisible, (newVal) => {
    dialogVisible.value = newVal
    if (newVal && props.enid) {
        // 每次打开时重新获取，或者你可以选择只在 enid 变化时获取
        getCourseInfo(props.enid)
    }
}, { immediate: true })

// 监听内部 dialogVisible 变化通知父组件
watch(dialogVisible, (newVal) => {
    if (!newVal) {
        emits('close')
    }
})
</script>

<style scoped lang="scss">
.course-info-dialog {
    :deep(.el-dialog__header) {
        margin-right: 0;
        padding: 20px 24px;
        border-bottom: 1px solid var(--border-soft);
    }

    :deep(.el-dialog__body) {
        padding: 0;
    }
}

.dialog-header-content {
    h3 {
        margin: 0 0 8px 0;
        font-size: 20px;
        font-weight: 600;
        color: var(--text-primary);
        line-height: 1.4;
    }
    
    .course-intro {
        margin: 0;
        font-size: 14px;
        color: var(--text-secondary);
        line-height: 1.5;
    }
}

.course-detail-container {
    padding: 24px;
    max-height: 70vh;
    overflow-y: auto;
    
    // 滚动条样式
    &::-webkit-scrollbar {
        width: 6px;
    }
    &::-webkit-scrollbar-thumb {
        background-color: var(--border-soft);
        border-radius: 3px;
    }
}

/* 讲师区域 */
.lecturer-section {
    margin-bottom: 24px;
    
    .lecturer-card {
        display: flex;
        gap: 20px;
        padding: 20px;
        background: var(--fill-color);
        border-radius: 12px;
        border: 1px solid var(--border-soft);
        
        .lecturer-avatar {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 8px;
            flex-shrink: 0;
            
            .lecturer-tag {
                margin-top: -12px;
                z-index: 1;
                box-shadow: 0 2px 6px rgba(0,0,0,0.1);
            }
        }
        
        .lecturer-info {
            flex: 1;
            
            .lecturer-name {
                font-size: 18px;
                font-weight: 600;
                color: var(--text-primary);
                margin-bottom: 8px;
            }
            
            .lecturer-desc {
                font-size: 14px;
                color: var(--text-secondary);
                line-height: 1.6;
                text-align: justify;
            }
        }
    }
}

/* 统计标签区域 */
.stats-container {
    display: flex;
    align-items: center;
    background: var(--fill-color);
    padding: 12px 20px;
    border-radius: 8px;
    margin-bottom: 24px;
    border: 1px solid var(--border-soft);
    
    .stat-item {
        display: flex;
        align-items: center;
        gap: 6px;
        color: var(--text-secondary);
        font-size: 13px;
        
        &.clickable {
            cursor: pointer;
            transition: color 0.2s;
            
            &:hover {
                color: var(--text-primary);
            }
        }
        
        .stat-icon {
            font-size: 16px;
            
            &.active {
                color: #e6a23c;
            }
        }
        
        .stat-text {
            display: flex;
            align-items: center;
            gap: 4px;
        }
        
        .stat-count {
            font-size: 12px;
        }
    }
    
    .stat-divider {
        width: 1px;
        height: 14px;
        background-color: var(--border-soft);
        margin: 0 16px;
    }
}

/* 标题样式 */
.section-title {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 16px;
    padding-left: 12px;
    border-left: 4px solid var(--el-color-primary);
    
    .label {
        margin-right: 4px;
    }
}

/* 评论区域 */
.comments-section {
    margin-bottom: 32px;
    
    .comments-carousel-wrapper {
        background: var(--fill-color);
        border-radius: 12px;
        padding: 16px;
        
        .comment-card {
            height: 100%;
            display: flex;
            flex-direction: column;
            gap: 8px;
            
            .comment-header {
                display: flex;
                align-items: center;
                justify-content: space-between;
                
                .user-info {
                    display: flex;
                    align-items: center;
                    gap: 8px;
                    
                    .username {
                        font-size: 13px;
                        color: var(--text-secondary);
                        font-weight: 500;
                    }
                }
            }
            
            .comment-content {
                flex: 1;
                overflow: hidden;
                
                h5 {
                    margin: 0 0 4px 0;
                    font-size: 14px;
                    font-weight: 600;
                    color: var(--text-primary);
                }
                
                .comment-text {
                    margin: 0;
                    font-size: 13px;
                    color: var(--text-regular);
                    line-height: 1.5;
                    display: -webkit-box;
                    -webkit-line-clamp: 3;
                    -webkit-box-orient: vertical;
                    overflow: hidden;
                    text-overflow: ellipsis;
                }
            }
        }
    }
}

/* 课程亮点 */
.highlights-section {
    margin-bottom: 32px;
    
    .highlight-content {
        background: linear-gradient(to bottom right, var(--fill-color), var(--bg-color));
        padding: 20px;
        border-radius: 12px;
        border: 1px solid var(--border-soft);
        font-size: 15px;
        line-height: 1.8;
        color: var(--text-regular);
        text-align: left;
        
        :deep(br) {
            display: block;
            margin: 8px 0;
            content: "";
        }
    }
}

/* 课程大纲 */
.outline-section {
    .outline-image-wrapper {
        border-radius: 12px;
        overflow: hidden;
        border: 1px solid var(--border-soft);
        background: var(--fill-color);
        min-height: 200px;
        display: flex;
        align-items: center;
        justify-content: center;
        
        .el-image {
            width: 100%;
            display: block;
        }
        
        .image-placeholder {
            color: var(--text-secondary);
            font-size: 14px;
        }
    }
}

/* 暗色模式微调 (大部分由 CSS 变量自动处理，这里只需处理特殊情况) */
// 如果有特定的暗色模式覆盖，可以在这里写，但优先使用变量
</style>