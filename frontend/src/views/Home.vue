<template>
  <div class="home-container-wrapper">
    <!-- 顶部区域：分类菜单 + 轮播图 + 用户信息 -->
    <div class="top-section">
      <!-- 左侧分类菜单 -->
      <div class="menu-wrapper">
        <el-menu
          ref="menuRef"
          default-active=""
          class="classification-menu"
          :collapse="false"
          unique-opened
          :router="true"
          active-text-color="var(--accent-color)"
          :popper-class="'menu-popper'"
          @open="handleOpen"
          @close="handleClose"
          @mouseenter="onMenuEnter"
          @mouseleave="onMenuLeave"
        >
          <el-sub-menu
            :index="item.enid"
            v-for="(item, index) in initial.homeData.categoryList"
            :key="item.enid"
            v-show="index < moreCategory"
            @mouseenter="handleMenuEnter(item.enid)"
          >
            <template #title>
              <span class="menu-title">{{ item.name }}</span>
            </template>
            <el-menu-item
              :index="i.enid"
              v-for="i in item.labelList"
              :key="i.enid"
              @click="gotoCategory(item, i.enid)"
            >
              <template #title>{{ i.name }}</template>
            </el-menu-item>
          </el-sub-menu>

          <el-sub-menu
            index="more"
            @mouseenter="mouseenter"
            v-show="9 == moreCategory"
          >
            <template #title><span class="menu-title">更多</span></template>
          </el-sub-menu>
        </el-menu>
      </div>

      <!-- 中间轮播图 -->
      <div class="banner-wrapper">
        <el-carousel :interval="5000" arrow="hover" height="380px" class="custom-carousel">
          <el-carousel-item v-for="item in initial.homeData.banner" :key="item">
            <el-image :src="item.img" fit="cover" class="banner-image" @click="BrowserOpenURL(item.url)" />
          </el-carousel-item>
        </el-carousel>
      </div>

      <!-- 右侧用户信息卡片 -->
      <div class="user-card-wrapper">
        <div class="user-card" :class="!Local.get('cookies') ? 'not-login' : ''">
          <div v-if="Local.get('cookies')==null" class="login-prompt">
            <div class="login-placeholder">
               <img src="../assets/images/logo-universal.png" alt="Logo" class="login-logo" />
               <p>登录开启学习之旅</p>
            </div>
            <el-button type="primary" class="login-btn" @click="openLoginDialog()" round>
              立即登录
            </el-button>
          </div>
          <div v-else class="user-info">
            <div class="avatar-section">
              <el-avatar
                :size="80"
                :src="user.avatar"
                class="user-avatar"
                @click="goToUserCenter"
              />
              <h3 class="user-name" @click="goToUserCenter">{{ user.nickname }}</h3>
            </div>
            <div class="user-stats">
              <div class="stat-item">
                <span class="stat-label">今日学习</span>
                <div class="stat-value">
                  <span class="number">{{ user.today_study_time > 0 ? (user.today_study_time / 60).toFixed(0) : '0' }}</span>
                  <span class="unit">分钟</span>
                </div>
              </div>
              <div class="stat-divider"></div>
              <div class="stat-item">
                <span class="stat-label">连续学习</span>
                <div class="stat-value">
                  <span class="number">{{ user.study_serial_days }}</span>
                  <span class="unit">天</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 内容模块列表 -->
    <div class="modules-container">
      <div class="module-section" v-for="item in initial.homeData.moduleList" :key="item.title">
        <div
          v-if="
            item.isShow == 3 &&
            (item.type == 'free_class' ||
              item.type == 'class' ||
              item.type == 'ebook')
          "
        >
          <div class="module-header">
            <div class="header-left">
              <h2 class="module-title">{{ item.title }}</h2>
              <span class="module-desc">{{ item.description }}</span>
            </div>
            
            <!-- 标签筛选 -->
            <div class="module-filters" v-if="item.type == 'class' || item.type == 'ebook'">
               <el-scrollbar>
                  <div class="filter-scroll">
                    <template v-if="item.type == 'class'">
                        <span 
                            v-for="(label, index) in courseLabelList.list"
                            :key="label.id"
                            class="filter-tag"
                            :class="idxCourseLabel === index ? 'active' : ''"
                            @click="handleLabel(label, index, 4)"
                        >
                            {{ label.name }}
                        </span>
                    </template>
                     <template v-if="item.type == 'ebook'">
                        <span 
                            v-for="(label, index) in ebookLabelList.list"
                            :key="label.id"
                            class="filter-tag"
                            :class="idxEbookLabel === index ? 'active' : ''"
                            @click="handleLabel(label, index, 2)"
                        >
                            {{ label.name }}
                        </span>
                    </template>
                  </div>
               </el-scrollbar>
            </div>
          </div>

          <div class="module-content">
            <!-- 免费专区 -->
             <div v-if="item.type == 'free_class'" class="cards-grid free-class-grid">
                <div 
                    v-for="(prod, index) in freeResourceList.list" 
                    :key="index"
                    class="content-card free-card"
                    @click="handleFreeProd(prod)"
                >
                    <div class="card-cover">
                        <img :src="ossProcess(prod.logo)" :alt="prod.name" loading="lazy" />
                        <div class="play-overlay"><el-icon><VideoPlay /></el-icon></div>
                    </div>
                    <div class="card-info">
                        <h4 class="card-title">{{ prod.name }}</h4>
                        <p class="card-intro">{{ prod.intro }}</p>
                        <div class="card-footer">
                             <el-rate
                                v-model.number="prod.score"
                                disabled
                                allow-half
                                size="small"
                            />
                        </div>
                    </div>
                </div>
             </div>

             <!-- 课程列表 -->
             <div v-if="item.type == 'class'" class="cards-grid course-grid">
                 <div 
                    v-for="(prod, index) in courseContentList.product_list" 
                    :key="prod.product_enid || index"
                    class="content-card course-card"
                    @click="handleClassProd(prod)"
                >
                    <div class="card-cover">
                         <img :src="ossProcess(prod.horizontal_image)" :alt="prod.title" loading="lazy" />
                    </div>
                     <div class="card-info">
                        <h4 class="card-title">{{ prod.title }}</h4>
                        <p class="card-intro">{{ prod.intro }}</p>
                        <div class="card-footer">
                            <span class="learn-count">{{ prod.learn_user_count }}人加入</span>
                        </div>
                    </div>
                </div>
                 <div class="more-btn-wrapper">
                    <el-button round @click="gotoCategory(currentCourse, '')">查看更多 {{ currentCourse.name }}</el-button>
                </div>
             </div>

             <!-- 电子书列表 -->
              <div v-if="item.type == 'ebook'" class="cards-grid ebook-grid">
                 <div 
                    v-for="(prod, index) in ebookContentList.product_list" 
                    :key="index"
                    class="content-card ebook-card"
                    @click="handleProd(prod.product_enid, prod.product_type)"
                >
                    <div class="ebook-cover-wrapper">
                         <img :src="ossEbookProcess(prod.index_image)" :alt="prod.title" loading="lazy" class="ebook-cover" />
                    </div>
                    <div class="card-info">
                        <h4 class="card-title">{{ ebookTitle(prod.title) }}</h4>
                        <p class="card-author">{{ authorList(prod.author_list) }}</p>
                         <div class="card-footer">
                            <el-rate
                                :model-value="handleScore(prod.score)"
                                disabled
                                allow-half
                                size="small"
                                v-if="prod.score.length > 0"
                            />
                             <span v-else class="no-score">暂无评分</span>
                        </div>
                    </div>
                </div>
                 <div class="more-btn-wrapper">
                    <el-button round @click="gotoCategory(currentEbook, '')">查看更多 {{ currentEbook.name }}</el-button>
                </div>
             </div>

          </div>
        </div>
      </div>
    </div>
  </div>

  <QrLogin
    v-if="loginVisible"
    :dialog-visible="loginVisible"
    @close="closeDialog"
  ></QrLogin>
  <EbookInfo
    v-if="ebookVisible"
    :enid="prodEnid"
    :dialog-visible="ebookVisible"
    @close="closeDialog"
  ></EbookInfo>
  <CourseInfo
    v-if="courseVisible"
    :enid="prodEnid"
    :dialog-visible="courseVisible"
    @close="closeDialog"
  ></CourseInfo>
</template>

<script lang="ts" setup>
import { ref, reactive, onBeforeMount, onMounted } from "vue";
import { ElMessage, ElDivider } from "element-plus";
import { VideoPlay } from '@element-plus/icons-vue'
import {
  GetHomeInitialState,
  SearchHot,
  SunflowerLabelList,
  SunflowerLabelContent,
  SunflowerResourceList,
  ArticleList,
  UserInfo,
} from "../../wailsjs/go/backend/App";
import { services } from "../../wailsjs/go/models";
import { BrowserOpenURL } from "../../wailsjs/runtime";
import QrLogin from "../components/QrLogin.vue";
import EbookInfo from "../components/EbookInfo.vue";
import CourseInfo from "../components/CourseInfo.vue";
import { useAppRouter } from "../composables/useRouter";
import { ROUTE_NAMES } from "../router/routes";
import { Local } from "../utils/storage";

const { pushByName, replace, pushCourseDetail } = useAppRouter();

const loading = ref(true);
const page = ref(0);
const total = ref(0);
const pageSize = ref(4);
const dialogVisible = ref(false);
const loginVisible = ref(false);
const ebookVisible = ref(false);
const courseVisible = ref(false);
const prodEnid = ref("");
const menuRef = ref();

const moreCategory = ref(9);
const idxEbookLabel = ref(0);
const idxCourseLabel = ref(0);
const isMouseInMenu = ref(false);
const currentOpenedMenu = ref<string>("");

let initial = reactive(new services.HomeInitState());
initial.homeData= reactive(new services.HomeData)

let ebookLabelList = reactive(new services.SunflowerLabelList());
let courseLabelList = reactive(new services.SunflowerLabelList());

let freeResourceList = reactive(new services.SunflowerResourceList());
let ebookContentList = reactive(new services.SunflowerContent());
let courseContentList = reactive(new services.SunflowerContent());

let currentCourse = reactive(new services.Navigation());
let currentEbook = reactive(new services.Navigation());
let user = reactive(new services.User());

onBeforeMount(() => {
  // 分类
  GetHomeInitialState()
    .then((state) => {
      Object.assign(initial, state);
      console.log(state);
      if (initial.isLogin) {
        getUserInfo();
      }
    })
    .catch((error) => {
      console.log(error);
    });
});

onMounted(() => {
  // 热搜词
  // SearchHot()
  //   .then((result) => {
  //     console.log(result);
  //   })
  //   .catch((error) => {
  //     console.log(error);
  //   }),
  // 电子书
  SunflowerLabelList(2)
    .then((result) => {
      Object.assign(ebookLabelList, result);
      Object.assign(currentEbook, ebookLabelList.list[0]);
      SunflowerLabelContent("", 2, 0, 10)
        .then((list) => {
          console.log(list);
          Object.assign(ebookContentList, list);
        })
        .catch((error) => {
          console.log(error);
        });
    })
    .catch((error) => {
      console.log(error);
    }),
    // 精选课程
    SunflowerLabelList(4)
      .then((result) => {
        Object.assign(courseLabelList, result);
        Object.assign(currentCourse, courseLabelList.list[0]);
        SunflowerLabelContent("", 4, 0, 4)
          .then((list) => {
            console.log(list);
            Object.assign(courseContentList, list);
          })
          .catch((error) => {
            console.log(error);
          });
        // console.log(result)
      })
      .catch((error) => {
        console.log(error);
      });
});

const getFreeResourceList = async () => {
  await SunflowerResourceList()
    .then((list) => {
      Object.assign(freeResourceList, list);
      console.log(freeResourceList);
    })
    .catch((error) => {
      console.log(error);
    });
};
getFreeResourceList();

const getUserInfo = async () => {
  await UserInfo()
    .then((result) => {
      Object.assign(user, result);
      console.log(user);
    })
    .catch((error) => {
      console.log(error);
    });
};


const handleProd = (enid: string, iType: number) => {
  // 2-ebook,4-专栏,36-大师课,66-class,22-course,65-学习圈
  // console.log('handleProd:', enid, iType);
  prodEnid.value = enid;
  if (iType == 2) {
    ebookVisible.value = true;
  } else if (iType == 65) {
    courseVisible.value = false;
  } else {
    courseVisible.value = true;
  }
};

const handleFreeProd = (item: any) => {
  // 统一跳转到详情页，不再弹窗
  pushByName(ROUTE_NAMES.ARTICLE_LIST, { id: item.enid }, { 
    enid: item.enid, 
    title: item.name 
  });
};

const handleClassProd = async (item: services.ProductSimple) => {
  const enid = String(item?.product_enid ?? "").trim();
  if (!enid) return;
  if (item?.product_type != 66) {
    courseVisible.value = false;
    return;
  };
  pushByName(ROUTE_NAMES.ARTICLE_LIST, { id: enid }, { enid, title: item.title });
};

const handleLabel = (
  item: services.Navigation,
  index: number,
  nType: number
) => {
  if (nType == 2) {
    pageSize.value = 10;
    Object.assign(currentEbook, item);
    idxEbookLabel.value = index;
  }
  if (nType == 4) {
    pageSize.value = 4;
    Object.assign(currentCourse, item);
    idxCourseLabel.value = index;
  }
  sunflowerLabelContent(item.enid, nType);
};

const sunflowerLabelContent = async (enid: string, nType: number) => {
  await SunflowerLabelContent(enid, nType, page.value, pageSize.value)
    .then((list) => {
      if (nType == 2) {
        Object.assign(ebookContentList, list);
      } else if (nType == 4) {
        Object.assign(courseContentList, list);
      }
      console.log(list);
    })
    .catch((error) => {
      console.log(error);
    });
};

const ossProcess = (url: string) => {
  return url + "?x-oss-process=image/crop,h_608/resize,w_1080,h_608,m_fill";
};

const ossAvatarProcess = (url: string) => {
  return url + "?x-oss-process=image/resize,w_96,m_lfit";
};

const ossEbookProcess = (url: string) => {
  return url + "?x-oss-process=image/crop,h_790/resize,w_600,h_790,m_fill";
};

const ebookTitle = (title: string) => {
  if (title.length <= 9) {
    return title;
  } else {
    return title.substring(0, 9).concat("...");
  }
};

const authorList = (authors: string[]) => {
  if (authors.length > 1) {
    return authors[0] + " 等";
  } else {
    return authors.toString();
  }
};

const ebookPopoverContent = (intro: string, introduction: string) => {
  if (intro.concat(introduction).length <= 146) {
    return intro.concat(introduction);
  } else {
    return intro.concat(introduction).substring(0, 146).concat("...");
  }
};

const mouseenter = () => {
  // 显示所有分类
  moreCategory.value = initial.homeData.categoryList.length;
};

const onMenuEnter = () => {
  // 鼠标进入菜单区域
  isMouseInMenu.value = true;
};

const onMenuLeave = (event: MouseEvent) => {
  // 检查鼠标是否还在菜单相关的元素上
  const menuElement = (event.currentTarget as HTMLElement);
  const relatedTarget = event.relatedTarget as Element | null;

  // 如果鼠标还在菜单内或二级菜单上，不隐藏
  if (relatedTarget &&
      (menuElement.contains(relatedTarget) ||
       (relatedTarget.closest && relatedTarget.closest('.el-menu--popup')))) {
    return;
  }

  // 鼠标真正离开菜单区域，隐藏多余的分类
  isMouseInMenu.value = false;
  moreCategory.value = 9;
  
  // 收起已展开的菜单
  if (menuRef.value && currentOpenedMenu.value) {
    menuRef.value.close(currentOpenedMenu.value);
    currentOpenedMenu.value = "";
  }
};

const handleMenuEnter = (index: string) => {
  // 鼠标悬浮展开二级菜单
  if (menuRef.value) {
    menuRef.value.open(index);
    currentOpenedMenu.value = index;
  }
};

const handleOpen = (key: string, keyPath: string[]) => {
  // 不需要在这里处理 more 的展开
};

const handleClose = (key: string, keyPath: string[]) => {
  // 不需要在这里处理 more 的收起
};

const handleScore = (score: string) => {
  return score != "" ? parseFloat(score) : 0;
};

const openLoginDialog = () => {
  loginVisible.value = true;
};

const openDialog = () => {
  dialogVisible.value = true;
};
const closeDialog = () => {
  //   initForm()
  dialogVisible.value = false;
  loginVisible.value = false;
  ebookVisible.value = false;
  courseVisible.value = false;
};

const goToUserCenter = () => {
  pushByName(ROUTE_NAMES.PROFILE);
};

const gotoCategory = (item: any, label_id: string) => {
  if (item.nav_type == 4 && item.enid && item.enid.startsWith("aiSphereGroupType:")) {
    pushByName("aiChannel");
    return;
  }

  let product_type = "0"; // 默认设置为课程类型
  if (item.nav_type == 2) {
    product_type = "2"; // 电子书
  } else if (item.nav_type == 4) {
    product_type = "66"; // 课程/专栏
  }
  pushByName(ROUTE_NAMES.CATEGORY, {}, {
    id: item.id,
    enid: item.enid,
    name: item.name,
    nav_type: item.nav_type,
    label_id: label_id,
    product_type: product_type,
  });
};
</script>

<style scoped>
.home-container-wrapper {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
  overflow-y: auto;
  box-sizing: border-box;
  /* 隐藏滚动条但保留功能 - 清新风格 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

.home-container-wrapper::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

/* --- Top Section --- */
.top-section {
  display: flex;
  gap: 24px;
  margin-bottom: 40px;
  height: 380px;
}

.menu-wrapper {
  width: 200px;
  flex-shrink: 0;
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: var(--shadow-soft);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.classification-menu {
  border: none;
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
  background-color: transparent !important; /* 确保背景透明，使用外层容器背景 */
  
  /* 隐藏滚动条但保留功能 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

/* 隐藏滚动条 */
.classification-menu::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

/* 强制覆盖 Element Plus 菜单选中和 Hover 样式，实现清新风格 */
:deep(.el-menu) {
  background-color: transparent;
}

:deep(.el-sub-menu__title) {
  height: 40px;
  line-height: 40px;
  border-radius: 8px;
  margin: 0 8px;
  color: var(--text-primary);
}

:deep(.el-sub-menu__title:hover) {
  background-color: var(--fill-color-light) !important;
  color: var(--accent-color);
}

:deep(.el-menu-item) {
  height: 36px;
  line-height: 36px;
  border-radius: 8px;
  margin: 4px 16px;
  color: var(--text-secondary);
}

:deep(.el-menu-item:hover) {
  background-color: var(--fill-color-light);
  color: var(--accent-color);
}

:deep(.el-menu-item.is-active) {
  background-color: var(--fill-color);
  color: var(--accent-color);
  font-weight: 500;
}


.banner-wrapper {
  flex: 1;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: var(--shadow-medium);
  min-width: 0; /* 防止 Flex 子元素溢出 */
}

.banner-image {
  width: 100%;
  height: 100%;
  cursor: pointer;
  transition: transform 0.5s ease;
}

.banner-image:hover {
  transform: scale(1.02);
}

.user-card-wrapper {
  width: 260px;
  flex-shrink: 0;
}

.user-card {
  height: 100%;
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: var(--shadow-soft);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 24px;
  box-sizing: border-box;
  text-align: center;
}

.user-card.not-login {
  background: linear-gradient(135deg, var(--card-bg) 0%, var(--fill-color-light) 100%);
}

.login-placeholder {
  margin-bottom: 24px;
}

.login-logo {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
  border-radius: 12px;
}

.login-placeholder p {
  color: var(--text-secondary);
  font-size: 14px;
  margin: 0;
}

.login-btn {
  width: 100%;
  font-weight: 600;
}

.user-info {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 20px;
}

.user-avatar {
  border: 4px solid var(--fill-color-light);
  cursor: pointer;
  transition: transform 0.3s ease;
}

.user-avatar:hover {
  transform: scale(1.05);
  border-color: var(--accent-color);
}

.user-name {
  margin: 16px 0 0;
  font-size: 18px;
  color: var(--text-primary);
  cursor: pointer;
}

.user-stats {
  display: flex;
  justify-content: space-around;
  align-items: center;
  background: var(--fill-color-light);
  border-radius: 8px;
  padding: 16px 8px;
  margin-bottom: 10px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.stat-value {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.stat-value .number {
  font-size: 20px;
  font-weight: bold;
  color: var(--accent-color);
  font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace;
}

.stat-value .unit {
  font-size: 12px;
  color: var(--text-tertiary);
}

.stat-divider {
  width: 1px;
  height: 24px;
  background-color: var(--border-color);
}

/* --- Modules Section --- */
.module-section {
  margin-bottom: 48px;
}

.module-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-soft);
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 16px;
}

.module-title {
  font-size: 24px;
  color: var(--text-primary);
  margin: 0;
}

.module-desc {
  font-size: 14px;
  color: var(--text-tertiary);
}

.module-filters {
  max-width: 50%;
}

.filter-scroll {
  display: flex;
  gap: 8px;
  padding-bottom: 4px;
}

.filter-tag {
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 16px;
  background: var(--fill-color);
  transition: all 0.2s ease;
  white-space: nowrap;
}

.filter-tag:hover {
  color: var(--accent-color);
  background: var(--fill-color-light);
}

.filter-tag.active {
  color: #fff;
  background: var(--accent-color);
}

/* --- Cards Grid --- */
.cards-grid {
  display: grid;
  gap: 24px;
}

.free-class-grid {
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
}

.course-grid {
  grid-template-columns: repeat(4, 1fr); /* 强制4列，与旧版逻辑保持一致，或者用 auto-fill */
}

.ebook-grid {
  grid-template-columns: repeat(5, 1fr);
}

.content-card {
  background: var(--card-bg);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: var(--shadow-soft);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  display: flex;
  flex-direction: column;
}

.content-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-medium);
}

.card-cover {
  width: 100%;
  aspect-ratio: 16/9;
  position: relative;
  overflow: hidden;
}

.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.content-card:hover .card-cover img {
  transform: scale(1.05);
}

.play-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.play-overlay .el-icon {
  font-size: 48px;
  color: #fff;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.5));
}

.content-card:hover .play-overlay {
  opacity: 1;
}

.ebook-cover-wrapper {
  width: 100%;
  aspect-ratio: 3/4;
  overflow: hidden;
  background: var(--fill-color-light);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 12px;
  box-sizing: border-box;
}

.ebook-cover {
  width: auto;
  height: 100%;
  max-width: 100%;
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
  transition: transform 0.3s ease;
}

.content-card:hover .ebook-cover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 8px 16px rgba(0,0,0,0.2);
}

.card-info {
  padding: 12px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 6px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-intro, .card-author {
  font-size: 12px;
  color: var(--text-secondary);
  margin: 0 0 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.card-footer {
  margin-top: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.learn-count {
  font-size: 12px;
  color: var(--text-tertiary);
}

.no-score {
  font-size: 12px;
  color: var(--text-placeholder);
}

.more-btn-wrapper {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    margin-top: 24px;
}

.more-btn-wrapper .el-button {
    width: 240px;
    height: 40px;
    font-size: 14px;
    border: 1px solid var(--border-soft);
    background-color: var(--card-bg);
    color: var(--text-secondary);
    transition: all 0.3s ease;
}

.more-btn-wrapper .el-button:hover {
    background-color: var(--fill-color-light);
    border-color: var(--accent-color);
    color: var(--accent-color);
    transform: translateY(-2px);
    box-shadow: var(--shadow-soft);
}

/* Responsive adjustments */
@media (max-width: 1200px) {
  .top-section {
    height: auto;
    flex-wrap: wrap;
  }
  
  .banner-wrapper {
    min-width: 60%;
  }
  
  .course-grid {
      grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 992px) {
    .course-grid {
      grid-template-columns: repeat(2, 1fr);
  }
}
</style>

