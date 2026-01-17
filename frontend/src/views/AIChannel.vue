<template>
  <div class="ai-channel">
    <!-- 频道信息头部 -->
    <div class="channel-header" v-if="channelInfo">
      <!-- 模糊背景层 -->
      <div class="header-bg-blur" :style="{ backgroundImage: `url(${channelInfo.logo})` }"></div>
      <div class="header-bg-overlay"></div>
      
      <div class="header-content">
        <!-- 左侧 Logo -->
        <div class="channel-logo" v-if="channelInfo.logo">
          <img :src="channelInfo.logo" :alt="channelInfo.title" />
        </div>

        <!-- 中间主要信息 -->
        <div class="channel-main">
          <div class="title-row">
            <h1>{{ channelInfo.title }}</h1>
            <el-tag v-if="channelInfo.is_vip" class="vip-tag" effect="dark">VIP</el-tag>
          </div>
          <p class="subtitle">{{ channelInfo.description }}</p>
          
          <!-- 人员信息区域 -->
          <div class="people-section">
            <div class="person-group" v-if="channelInfo.host">
              <span class="group-label">主持人</span>
              <div class="person-content">
                <el-avatar :size="24" :src="channelInfo.host.avatar" v-if="channelInfo.host.avatar" />
                <span class="person-name">{{ channelInfo.host.name }}</span>
              </div>
            </div>

            <el-divider direction="vertical" class="section-divider" v-if="channelInfo.host && channelInfo.guests && channelInfo.guests.length > 0" />

            <div class="person-group" v-if="channelInfo.guests && channelInfo.guests.length > 0">
              <span class="group-label">嘉宾</span>
              <div class="guest-list">
                <el-tooltip
                  v-for="guest in channelInfo.guests"
                  :key="guest.uid"
                  effect="dark"
                  :content="guest.bio"
                  placement="top"
                  :disabled="!guest.bio"
                >
                  <div class="guest-item">
                    <el-avatar :size="24" :src="guest.avatar" />
                    <span class="person-name">{{ guest.name }}</span>
                  </div>
                </el-tooltip>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧数据与状态 -->
        <div class="channel-side">
           <div class="stat-card" v-if="channelInfo.statistics?.total_subscribers">
              <div class="stat-value">{{ formatNumber(channelInfo.statistics.total_subscribers) }}</div>
              <div class="stat-label">人已加入</div>
           </div>

           <template v-if="vipInfo && vipInfo.is_ever_subscribed">
              <div class="status-card" :class="{ 'is-expired': !vipInfo.is_vip }">
                <div class="status-icon">
                  <el-icon v-if="vipInfo.is_vip"><Trophy /></el-icon>
                  <el-icon v-else><Warning /></el-icon>
                </div>
                <div class="status-info">
                   <template v-if="vipInfo.is_vip">
                      <div class="status-value">{{ vipInfo.surplus_days }}天</div>
                      <div class="status-label">VIP剩余有效期</div>
                   </template>
                   <template v-else>
                      <div class="status-value">{{ getExpiredDays(vipInfo) }}天</div>
                      <div class="status-label">VIP已过期</div>
                   </template>
                </div>
              </div>
           </template>
        </div>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="error-state">
      <el-result
        icon="error"
        title="加载失败"
        :sub-title="error"
      >
        <template #extra>
          <el-button type="primary" @click="fetchChannelHomepage">重新加载</el-button>
        </template>
      </el-result>
    </div>

    <!-- 骨架屏加载状态 (仅初始化) -->
    <div v-if="loading && categories.length === 0" class="loading-state">
      <div class="main-content">
         <div class="category-tabs">
            <el-skeleton :rows="1" animated />
         </div>
         <div class="content-area">
            <div class="subcategory-list">
              <el-skeleton :count="5" animated />
            </div>
            <div class="items-container">
               <el-skeleton :count="3" animated style="margin-bottom: 20px" />
            </div>
         </div>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="main-content" v-else-if="categories.length > 0">
      <!-- 一级分类 Tab 导航 -->
      <div class="category-tabs">
        <el-tabs
          v-model="activeCategoryId"
          @tab-click="handleTabClick"
          class="category-tabs-container"
        >
          <el-tab-pane
            v-for="cat in categories"
            :key="cat.category_id"
            :name="cat.category_id"
          >
            <template #label>
              <div class="tab-label">
                <div class="tab-icon">
                  <img v-if="cat.category_icon" :src="cat.category_icon" class="category-icon" />
                  <el-icon v-else><Document /></el-icon>
                </div>
                <span class="tab-name">
                  {{ cat.category_name }}
                  <span class="tab-count" v-if="getCategoryItemCount(cat) > 0">
                    {{ getCategoryItemCount(cat) }}
                  </span>
                </span>
              </div>
            </template>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 内容区域 -->
      <div class="content-area">
        <!-- 二级分类列表 -->
        <div class="subcategory-wrapper" v-if="currentSubcategories.length > 0">
           <div class="subcategory-list">
              <div
                v-for="subcat in currentSubcategories"
                :key="subcat.id"
                class="subcategory-item"
                :class="{ 'is-active': activeSubcategoryId === subcat.id }"
                @click="handleSubcategoryClick(subcat)"
              >
                <div class="subcategory-indicator"></div>
                <div class="subcategory-icon" v-if="subcat.icon">
                  <img :src="subcat.icon" :alt="subcat.title" />
                </div>
                <div class="subcategory-content">
                  <h4 class="subcategory-title">{{ subcat.title }}</h4>
                  <p v-if="subcat.intro" class="subcategory-intro">{{ subcat.intro }}</p>
                </div>
                <div class="subcategory-count" v-if="subcat.length">
                  <span class="count-badge">{{ subcat.length }}</span>
                </div>
              </div>
           </div>
        </div>

        <!-- 内容项列表 -->
        <div class="items-container" v-if="activeSubcategoryId">
          <div class="items-header">
            <h3>
              {{ currentSubcategories.find(s => s.id === activeSubcategoryId)?.title }}
              <span class="items-total" v-if="!listLoading">{{ currentItems.length }} 项内容</span>
            </h3>
          </div>
          
          <!-- 列表加载中 -->
          <div v-if="listLoading" class="items-loading">
             <el-skeleton :rows="3" animated />
             <br/>
             <el-skeleton :rows="3" animated />
          </div>

          <!-- 列表内容 -->
          <transition name="el-fade-in-linear">
            <div v-if="!listLoading && currentItems.length > 0" class="items-grid">
              <div
                v-for="(item, index) in currentItems"
                :key="item.product_id || index"
                class="item-card"
                @click="handleItemClick(item)"
              >
                <div class="item-cover-wrapper">
                  <div class="item-cover">
                    <img v-if="item.cover" :src="item.cover" loading="lazy" />
                    <div class="item-difficulty-badge" v-if="getDifficultyLabel(item.difficulty_level)">
                      {{ getDifficultyLabel(item.difficulty_level) }}
                    </div>
                    <div class="item-type-badge">
                      {{ getProductTypeName(item.product_type) }}
                    </div>
                  </div>
                </div>
                
                <div class="item-content">
                  <h4 class="item-title" :title="item.title">{{ item.title }}</h4>
                  <p v-if="item.summary" class="item-description">{{ item.summary }}</p>
                  
                  <div class="item-footer">
                    <div class="meta-left">
                        <div v-if="item.duration" class="meta-item">
                          <el-icon><Clock /></el-icon>
                          {{ formatDuration(item.duration) }}
                        </div>
                        <div v-if="item.learn_count" class="meta-item">
                          <el-icon><VideoPlay /></el-icon>
                          {{ formatNumber(item.learn_count) }}
                        </div>
                    </div>
                    
                    <div class="item-actions" v-if="item.en_id">
                      <el-button
                        class="download-btn"
                        circle
                        size="small"
                        text
                        @click.stop="openItemDetail(item)"
                      >
                        <el-icon><Document /></el-icon>
                      </el-button>
                      <el-button
                        class="download-btn"
                        circle
                        size="small"
                        text
                        @click.stop="openOfficialArticle(item)"
                      >
                        <el-icon><Link /></el-icon>
                      </el-button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </transition>

          <!-- 无数据状态 -->
          <div v-if="!listLoading && currentItems.length === 0" class="empty-data">
            <el-empty description="该分类下暂无内容" />
          </div>
        </div>

        <!-- 空状态 (未选择二级分类) -->
        <div v-else-if="!loading" class="empty-state">
          <el-empty
            :image-size="200"
            description="请选择二级分类"
          />
        </div>
      </div>
    </div>

    <!-- 课程详情弹窗 -->
    <CourseInfo
      v-if="showCourseInfo"
      :enid="currentClassEnId"
      :dialogVisible="showCourseInfo"
      @close="showCourseInfo = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document,
  VideoPlay,
  Clock,
  Link,
  Trophy,
  Warning
} from '@element-plus/icons-vue'
import type { services } from '../../wailsjs/go/models'
import {
  ChannelInfo as getChannelInfo,
  ChannelHomepage as getChannelHomepage,
  ChannelVipInfo as getVipInfo,
  ChannelTopicDetail,
} from '../../wailsjs/go/backend/App'
import CourseInfo from '../components/CourseInfo.vue'
import { useAppRouter } from '../composables/useRouter'

// 状态变量
const loading = ref(false)
const listLoading = ref(false)
const channelInfo = ref<services.ChannelInfo | null>(null)
const vipInfo = ref<services.ChannelVipInfo | null>(null)
const categories = ref<services.ChannelHomepageCategory[]>([])
const activeCategoryId = ref<number | null>(null)
const activeSubcategoryId = ref<number | null>(null)
const error = ref<string | null>(null)
const showCourseInfo = ref(false)
const currentClassEnId = ref('')
const { pushArticleDetail, openDedaoArticle } = useAppRouter()

// 存储完整的二级分类数据
const subcategoryDataMap = ref<Map<number, services.ChannelTopicCategory>>(new Map())

// 计算属性
const currentSubcategories = computed(() => {
  if (!activeCategoryId.value) return []
  const category = categories.value.find((cat: services.ChannelHomepageCategory) => cat.category_id === activeCategoryId.value)
  return category?.list || []
})

const currentItems = computed(() => {
  if (!activeSubcategoryId.value) return []
  const subcategoryData = subcategoryDataMap.value.get(activeSubcategoryId.value)
  return subcategoryData?.items || []
})

const openItemDetail = (item: any) => {
  const enid = String(item?.en_id ?? '').trim()
  if (!enid) {
    ElMessage.warning('缺少 enid，无法查看详情')
    return
  }
  pushArticleDetail(enid, 'aiChannel', {
    parentTitle: channelInfo.value?.title ?? 'AI学习圈',
  })
}

const openOfficialArticle = (item: any) => {
  const enid = String(item?.en_id ?? '').trim()
  if (!enid) {
    ElMessage.warning('缺少 enid，无法打开官网页面')
    return
  }
  openDedaoArticle(enid)
}

// 获取难度标签
const getDifficultyLabel = (level: number): string => {
  const map: Record<number, string> = {
    1: '入门',
    2: '进阶',
    3: '高阶'
  }
  return map[level] || ''
}

// 获取分类下的总项目数
const getCategoryItemCount = (category: services.ChannelHomepageCategory): number => {
  return category.list?.reduce((total: number, subcat: services.ChannelTopicCategory) => total + (subcat.length || 0), 0) || 0
}

// 获取频道信息
const fetchChannelInfo = async () => {
  try {
    const info = await getChannelInfo(props.channelId)
    channelInfo.value = info
    console.log('频道信息:', info)
    error.value = null
  } catch (err) {
    console.error('获取频道信息失败:', err)
    error.value = '获取频道信息失败'
    ElMessage.error('获取频道信息失败')
  }
}

// 获取频道首页分类
const fetchChannelHomepage = async () => {
  try {
    loading.value = true
    const data = await getChannelHomepage(props.channelId)
    categories.value = data
    error.value = null
    // 如果有分类，尝试恢复状态
    if (data.length > 0) {
      const savedState = restoreState()
      
      // 检查保存的一级分类是否存在于当前数据中
      const categoryExists = savedState?.categoryId && data.some(c => c.category_id === savedState.categoryId)
      
      if (categoryExists) {
        activeCategoryId.value = savedState.categoryId
        
        // 恢复二级分类选中状态
        if (savedState.subcategoryId) {
           // 找到对应的一级分类
           const currentCategory = data.find(c => c.category_id === savedState.categoryId)
           // 检查二级分类是否存在
           const subcategory = currentCategory?.list?.find(s => s.id === savedState.subcategoryId)
           
           if (subcategory) {
             // 触发点击事件逻辑以加载数据
             handleSubcategoryClick(subcategory)
           }
        }
      } else if (!activeCategoryId.value) {
        // 如果没有保存状态或保存的状态无效，且当前未选中，则默认选中第一个
        activeCategoryId.value = data[0].category_id
      }
    }
  } catch (err) {
    console.error('获取频道分类失败:', err)
    error.value = '获取频道分类失败'
    ElMessage.error('获取频道分类失败')
  } finally {
    loading.value = false
  }
}

// 获取VIP信息
const fetchVipInfo = async () => {
  try {
    const info = await getVipInfo(props.channelId)
    vipInfo.value = info
    console.log('VIP信息:', info)
  } catch (error) {
    console.error('获取VIP信息失败:', error)
  }
}

// 存储 key 前缀
const STORAGE_KEY_PREFIX = 'ai_channel_state_'

// 保存状态到 localStorage
const saveState = () => {
  if (props.channelId) {
    const state = {
      categoryId: activeCategoryId.value,
      subcategoryId: activeSubcategoryId.value
    }
    localStorage.setItem(`${STORAGE_KEY_PREFIX}${props.channelId}`, JSON.stringify(state))
  }
}

// 恢复状态
const restoreState = () => {
  if (props.channelId) {
    try {
      const savedState = localStorage.getItem(`${STORAGE_KEY_PREFIX}${props.channelId}`)
      if (savedState) {
        const state = JSON.parse(savedState)
        return state
      }
    } catch (e) {
      console.error('恢复状态失败:', e)
    }
  }
  return null
}

// 处理一级分类切换
const handleTabClick = (tab: any) => {
  const categoryId = Number(tab.props.name)
  activeCategoryId.value = categoryId
  activeSubcategoryId.value = null
  // 清空当前的数据
  subcategoryDataMap.value.clear()
  saveState()
}

// 处理二级分类点击
const handleSubcategoryClick = async (subcategory: services.ChannelTopicCategory) => {
  activeSubcategoryId.value = subcategory.id
  saveState()

  // 如果已经缓存了数据，直接使用
  if (subcategoryDataMap.value.has(subcategory.id)) {
    return
  }

  // 否则获取完整数据
  listLoading.value = true
  try {

    const fullData = await ChannelTopicDetail(subcategory.id)
    if (fullData && fullData.items) {
      subcategoryDataMap.value.set(subcategory.id, fullData)
    } else {
      // 即使返回空结构体，也缓存它避免重复请求
      // 创建一个默认的结构体
      const defaultData = {
        id: subcategory.id,
        title: subcategory.title,
        items: subcategory.items || [], // 如果原始数据有items，使用它
        length: subcategory.length || 0
      }
      subcategoryDataMap.value.set(subcategory.id, defaultData as services.ChannelTopicCategory)
    }
  } catch (error) {
    ElMessage.error('获取分类详情失败')
  } finally {
    listLoading.value = false
  }
}

// 处理内容项点击
const handleItemClick = (item: services.ChannelItem) => {
  console.log('点击内容项:', item)

  if (canPlayItem(item)) {
    openItemDetail(item)
    return
  }

  if (item.class_en_id) {
    currentClassEnId.value = item.class_en_id
    showCourseInfo.value = true
    return
  }

  ElMessage.warning('该内容暂无课程信息')
}

const canPlayItem = (item: services.ChannelItem) => {
  return item.product_type === 65 && !!item.en_id
}

// 获取产品类型名称
const getProductTypeName = (productType: number): string => {
  const typeMap: Record<number, string> = {
    65: '视频',
    66: '课程',
  }
  return typeMap[productType] || `类型${productType}`
}

// 格式化时长
const formatDuration = (seconds: number): string => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)

  if (hours > 0) {
    return `${hours}时${minutes}分`
  }
  return `${minutes}分`
}

// 格式化数字
const formatNumber = (num: number): string => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// 计算过期天数
const getExpiredDays = (info: services.ChannelVipInfo): number => {
  if (!info.current_time || !info.end_time) return 0
  const diffSeconds = info.current_time - info.end_time
  if (diffSeconds <= 0) return 0
  return Math.floor(diffSeconds / (24 * 3600))
}

// Props
const props = defineProps<{
  channelId: number
}>()

// 初始化
onMounted(async () => {
  await Promise.all([
    fetchChannelInfo(),
    fetchChannelHomepage(),
    fetchVipInfo()
  ])
})
</script>

<style scoped lang="scss">
.ai-channel {
  min-height: 100vh;
  /* background-color: var(--bg-color); */
  /* color: var(--text-color); */
  /* transition: background-color 0.3s ease, color 0.3s ease; */
  padding-bottom: 40px;
}

/* 头部样式优化 */
.channel-header {
  position: relative;
  overflow: hidden;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border-soft);
  
  .header-bg-blur {
    position: absolute;
    top: -20px;
    left: -20px;
    right: -20px;
    bottom: -20px;
    background-size: cover;
    background-position: center;
    filter: blur(60px);
    opacity: 0.15;
    z-index: 0;
  }
  
  .header-bg-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(to bottom, transparent, var(--bg-color));
    z-index: 1;
  }

  .header-content {
    position: relative;
    z-index: 2;
    max-width: 1400px;
    margin: 0 auto;
    padding: 48px 32px;
    display: flex;
    align-items: flex-start;
    gap: 32px;

    .channel-logo {
      width: 120px;
      height: 120px;
      flex-shrink: 0;
      border-radius: 16px;
      overflow: hidden;
      background: var(--fill-color);
      border: 1px solid var(--border-soft);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .channel-main {
      flex: 1;
      padding-top: 8px;

      .title-row {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 12px;

        h1 {
          margin: 0;
          font-size: 32px;
          font-weight: 700;
          color: var(--text-primary);
          line-height: 1.2;
        }

        .vip-tag {
           background: linear-gradient(135deg, #ffd700, #ffa500);
           border: none;
           color: #fff;
           font-weight: bold;
           letter-spacing: 1px;
        }
      }

      .subtitle {
        margin: 0 0 24px 0;
        font-size: 16px;
        color: var(--text-secondary);
        line-height: 1.6;
        max-width: 800px;
      }

      .people-section {
        display: flex;
        align-items: flex-start;
        gap: 32px;
        padding-top: 8px;

        .section-divider {
           height: 40px;
           margin: 0;
           align-self: center;
           border-color: var(--border-soft);
        }

        .person-group {
          display: flex;
          flex-direction: column;
          gap: 12px;

          .group-label {
            font-size: 13px;
            color: var(--text-tertiary);
            font-weight: 500;
          }

          .person-content {
             display: flex;
             align-items: center;
             gap: 10px;
             background: var(--fill-color);
             padding: 6px 16px 6px 6px;
             border-radius: 50px;
             border: 1px solid var(--border-soft);

             .person-name {
               font-size: 14px;
               font-weight: 500;
               color: var(--text-primary);
             }
          }

          .guest-list {
             display: flex;
             align-items: center;
             gap: 12px;
             flex-wrap: wrap;
             
             .guest-item {
               display: flex;
               align-items: center;
               gap: 8px;
               background: var(--fill-color);
               padding: 4px 12px 4px 4px;
               border-radius: 50px;
               border: 1px solid var(--border-soft);
               
               .person-name {
                 font-size: 13px;
                 color: var(--text-primary);
               }
             }
          }
        }
      }
    }

    .channel-side {
      width: 260px;
      flex-shrink: 0;
      display: flex;
      flex-direction: column;
      gap: 16px;
      padding-top: 8px;

      .stat-card {
        background: rgba(255, 255, 255, 0.5);
        backdrop-filter: blur(10px);
        padding: 20px;
        border-radius: 16px;
        border: 1px solid var(--border-soft);
        text-align: center;
        transition: all 0.3s ease;
        
        &:hover {
          background: var(--card-bg);
          box-shadow: var(--shadow-soft);
        }

        .stat-value {
           font-size: 24px;
           font-weight: 700;
           color: var(--text-primary);
           line-height: 1.2;
           margin-bottom: 4px;
           font-family: var(--font-family-mono);
        }

        .stat-label {
           font-size: 13px;
           color: var(--text-tertiary);
        }
      }

      .status-card {
        background: linear-gradient(135deg, #2c3e50, #000000);
        color: #ffd700;
        padding: 16px;
        border-radius: 16px;
        display: flex;
        align-items: center;
        gap: 16px;
        box-shadow: 0 8px 20px rgba(0,0,0,0.1);
        border: 1px solid rgba(255,255,255,0.1);

        &.is-expired {
           background: var(--fill-color);
           color: var(--text-secondary);
           border: 1px solid var(--border-soft);
           box-shadow: none;
        }

        .status-icon {
           font-size: 28px;
           display: flex;
           align-items: center;
           justify-content: center;
           width: 48px;
           height: 48px;
           background: rgba(255,255,255,0.1);
           border-radius: 12px;
        }

        .status-info {
           .status-value {
              font-size: 20px;
              font-weight: 700;
              line-height: 1.2;
              margin-bottom: 2px;
           }
           .status-label {
              font-size: 12px;
              opacity: 0.8;
           }
        }
      }
    }
  }
}

.main-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px;
}

/* 分类 Tab 优化 */
.category-tabs {
  margin-bottom: 32px;
  /* position: sticky; */
  /* top: 0; */
  /* z-index: 10; */
  background: var(--bg-color);
  padding-top: 16px;
  padding-bottom: 8px;

  .category-tabs-container {
    background: var(--card-bg);
    border-radius: 16px;
    padding: 6px;
    box-shadow: var(--shadow-soft);
    border: 1px solid var(--border-soft);

    :deep(.el-tabs__header) {
      margin-bottom: 0;
    }

    :deep(.el-tabs__nav-wrap::after) {
      display: none;
    }

    :deep(.el-tabs__item) {
      padding: 0 24px;
      height: 48px;
      line-height: 48px;
      border-radius: 10px;
      transition: all 0.3s ease;
      margin: 4px;

      &:hover {
        background: var(--fill-color);
        color: var(--text-primary);
      }

      &.is-active {
        background: var(--fill-color);
        color: var(--el-color-primary);
        font-weight: 600;
        box-shadow: none;
      }
    }
  }

  .tab-label {
    display: flex;
    align-items: center;
    gap: 8px;

    .tab-icon {
      width: 18px;
      height: 18px;
      display: flex;
      align-items: center;
      justify-content: center;

      .category-icon {
        width: 100%;
        height: 100%;
        object-fit: contain;
      }
    }

    .tab-name {
      font-size: 15px;
      font-weight: 500;

      .tab-count {
        font-size: 12px;
        opacity: 0.8;
        margin-left: 2px;
        background: rgba(255,255,255,0.2);
        padding: 0 6px;
        border-radius: 10px;
      }
    }
  }
}

.content-area {
  display: flex;
  gap: 32px;
  align-items: flex-start;
}

/* 左侧侧边栏优化 */
.subcategory-wrapper {
  flex: 0 0 280px;
  position: sticky;
  top: 20px; /* 调整为较小的顶部间距，因为 tab 不再固定 */
  max-height: calc(100vh - 40px);
  overflow-y: auto;
  
  /* 隐藏滚动条但保留功能 */
  &::-webkit-scrollbar {
    width: 4px;
  }
  &::-webkit-scrollbar-thumb {
    background: transparent;
  }
  &:hover::-webkit-scrollbar-thumb {
    background: var(--border-soft);
  }
}

.subcategory-list {
  background: var(--card-bg);
  border-radius: 16px;
  padding: 12px;
  border: 1px solid var(--border-soft);

  .subcategory-item {
    position: relative;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.2s ease;
    margin-bottom: 4px;
    overflow: hidden;

    &:hover {
      background: var(--fill-color);
    }

    &.is-active {
      background: var(--fill-color);
      
      .subcategory-title {
        color: var(--el-color-primary);
        font-weight: 600;
      }
      
      .subcategory-indicator {
        opacity: 1;
        transform: scaleY(1);
      }
    }
    
    .subcategory-indicator {
      position: absolute;
      left: 0;
      top: 12px;
      bottom: 12px;
      width: 4px;
      background: var(--el-color-primary);
      border-radius: 0 4px 4px 0;
      opacity: 0;
      transform: scaleY(0.5);
      transition: all 0.2s ease;
    }

    .subcategory-icon {
      width: 36px;
      height: 36px;
      flex-shrink: 0;
      border-radius: 8px;
      overflow: hidden;
      background: var(--fill-color);

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .subcategory-content {
      flex: 1;
      min-width: 0;

      .subcategory-title {
        margin: 0;
        font-size: 14px;
        font-weight: 500;
        color: var(--text-primary);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        transition: color 0.2s;
      }

      .subcategory-intro {
        margin: 4px 0 0 0;
        font-size: 12px;
        color: var(--text-secondary);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        opacity: 0.8;
      }
    }

    .subcategory-count {
      flex-shrink: 0;
      
      .count-badge {
         font-size: 11px;
         background: var(--fill-color);
         color: var(--text-secondary);
         padding: 2px 8px;
         border-radius: 10px;
      }
    }
  }
}

/* 内容区域优化 */
.items-container {
  flex: 1;
  min-width: 0; /* 防止 grid 溢出 */

  .items-loading {
    padding: 20px 0;
  }

  .items-header {
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--border-soft);

    h3 {
      margin: 0;
      font-size: 20px;
      font-weight: 600;
      color: var(--text-primary);
      display: flex;
      align-items: baseline;
      gap: 12px;

      .items-total {
        font-size: 14px;
        color: var(--text-secondary);
        font-weight: normal;
      }
    }
  }

  .items-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 24px;
  }

  .item-card {
    background: var(--card-bg);
    border-radius: 16px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    border: 1px solid var(--border-soft);
    display: flex;
    flex-direction: column;
    height: 100%;

    &:hover {
      transform: translateY(-6px);
      box-shadow: 0 12px 32px rgba(0, 0, 0, 0.1);
      border-color: transparent;
      
      .item-cover img {
        transform: scale(1.05);
      }
      
      .hover-overlay {
        opacity: 1;
      }
      
      .item-title {
        color: var(--el-color-primary);
      }
    }

    .item-cover-wrapper {
      position: relative;
      padding-top: 56.25%; /* 16:9 比例 */
      overflow: hidden;
    }

    .item-cover {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: var(--fill-color);

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.5s ease;
      }
      
      .item-type-badge {
         position: absolute;
         top: 8px;
         right: 8px;
         background: rgba(0,0,0,0.6);
         backdrop-filter: blur(4px);
         color: white;
         font-size: 11px;
         padding: 4px 8px;
         border-radius: 6px;
         z-index: 2;
      }

      .item-difficulty-badge {
         position: absolute;
         top: 8px;
         left: 8px;
         background: rgba(0,0,0,0.6);
         backdrop-filter: blur(4px);
         color: white;
         font-size: 11px;
         padding: 4px 8px;
         border-radius: 6px;
         z-index: 2;
      }
      
      .hover-overlay {
         position: absolute;
         top: 0;
         left: 0;
         right: 0;
         bottom: 0;
         background: rgba(0,0,0,0.3);
         display: flex;
         align-items: center;
         justify-content: center;
         opacity: 0;
         transition: opacity 0.3s ease;
         z-index: 1;
         
         .play-icon {
            font-size: 48px;
            color: white;
            filter: drop-shadow(0 4px 8px rgba(0,0,0,0.3));
         }
      }
    }

    .item-content {
      padding: 16px;
      flex: 1;
      display: flex;
      flex-direction: column;

      .item-title {
        margin: 0 0 8px 0;
        font-size: 16px;
        font-weight: 600;
        color: var(--text-primary);
        line-height: 1.4;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        transition: color 0.2s;
      }

      .item-description {
        margin: 0 0 16px 0;
        font-size: 13px;
        color: var(--text-secondary);
        line-height: 1.6;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        flex: 1; /* 撑开高度，保持底部对齐 */
      }

      .item-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-top: 12px;
        border-top: 1px solid var(--border-soft);
        
        .meta-left {
           display: flex;
           gap: 12px;
        }

        .meta-item {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 12px;
          color: var(--text-tertiary);

          .el-icon {
            font-size: 14px;
          }
        }
        
        .download-btn {
           color: var(--text-tertiary);
           opacity: 0.8;
           transition: all 0.2s;
           
           &:hover {
              color: var(--el-color-primary);
              background: var(--fill-color);
              opacity: 1;
              transform: scale(1.1);
           }
        }

        .item-actions {
          display: flex;
          gap: 6px;
        }
      }
    }
  }
}

.empty-state, .loading-state {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 60px 0;
  
  :deep(.el-empty__description) {
    color: var(--text-secondary);
  }
}

@media (max-width: 900px) {
  .content-area {
    flex-direction: column;
  }
  
  .subcategory-wrapper {
     flex: none;
     width: 100%;
     position: static;
     max-height: none;
     margin-bottom: 24px;
  }

  .subcategory-list {
    display: flex;
    overflow-x: auto;
    padding: 8px;
    gap: 8px;
    
    .subcategory-item {
       flex: 0 0 auto;
       width: 200px;
       margin-bottom: 0;
       
       .subcategory-intro {
          display: none;
       }
    }
  }

  .items-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }
}
</style>
