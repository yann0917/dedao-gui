<template>
  <el-row :gutter="30" class="home-banner">
    <el-col :span="3">
      <el-menu
        default-active=""
        class="classification"
        :collapse="true"
        :router="true"
        active-text-color="#ffd04b"
        :popper-class="'menu-popper'"
        @open="handleOpen"
        @close="handleClose"
        @mouseleave="mouseleave"
      >
        <el-sub-menu
          :index="item.enid"
          v-for="(item, index) in initial.homeData.categoryList"
          v-show="index < moreCategory"
          popper-append-to-body
        >
          <template #title>{{ item.name }}</template>
          <el-menu-item
            :index="i.enid"
            v-for="i in item.labelList"
            @click="gotoCategory(item, i.enid)"
          >
            <template #title>{{ i.name }}</template>
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu
          index="more"
          @mouseenter="mouseenter"
          popper-append-to-body
          v-show="9 == moreCategory"
        >
          <template #title>更多</template>
        </el-sub-menu>
      </el-menu>
    </el-col>

    <el-col :span="16" class="banner">
      <el-carousel :interval="5000" arrow="always">
        <el-carousel-item v-for="item in initial.homeData.banner" :key="item">
          <el-image :src="item.img" @click="BrowserOpenURL(item.url)">
          </el-image>
        </el-carousel-item>
      </el-carousel>
    </el-col>
    <el-col :span="5" class="user">
      <div :class="Local.get('cookies')==null ?'not-login' : 'logged'">
        <div v-if="Local.get('cookies')==null">
          <div class="receive"></div>
          <el-button class="login-btn" @click="openLoginDialog()">
            登录
          </el-button>
        </div>
        <div v-else>
          <div class="personal">
            <el-avatar :size="72" :src="user.avatar" fit="fill" />
            <h3>{{ user.nickname }}</h3>
          </div>
          <div class="data">
            <p class="time">
              <span>今日学习</span><em style="font-size: 22px;color:#333;">{{user.today_study_time > 0 ? (user.today_study_time / 60).toFixed(0):''}}</em><span>分钟</span>
            </p>
            <el-divider border-style="dotted" />
            <p class="time">
              <span>连续学习</span><em style="font-size: 22px;color: #333;">{{user.study_serial_days}}</em><span>天</span>
            </p>
          </div>
          <!-- <el-button class="button" @click="logout()"> 退出 </el-button> -->
        </div>
      </div>
    </el-col>
  </el-row>

  <div class="module">
    <div class="iget-ui-container" v-for="item in initial.homeData.moduleList">
      <div
        v-if="
          item.isShow == 3 &&
          (item.type == 'free_class' ||
            item.type == 'class' ||
            item.type == 'ebook')
        "
      >
        <div class="module-title-wrap">
          <h1>{{ item.title }}</h1>
          <p>{{ item.description }}</p>
        </div>

        <div class="home-container">
          <el-scrollbar>
            <div class="scrollbar-flex-content" v-if="item.type == 'class'">
              <el-button
                v-for="(item, index) in courseLabelList.list"
                :key="item.id"
                :class="
                  idxCourseLabel === index ? 'active-btn' : 'scrollbar-item'
                "
                text
                @click="handleLabel(item, index, 4)"
              >
                {{ item.name }}
              </el-button>
            </div>
            <div class="scrollbar-flex-content" v-if="item.type == 'ebook'">
              <el-button
                v-for="(item, index) in ebookLabelList.list"
                :key="item.id"
                :class="
                  idxEbookLabel === index ? 'active-btn' : 'scrollbar-item'
                "
                text
                @click="handleLabel(item, index, 2)"
              >
                {{ item.name }}
              </el-button>
            </div>
          </el-scrollbar>

          <div class="module-cards-wrap">
            <el-scrollbar v-if="item.type == 'free_class'">
              <div class="cards-cover">
                <el-card
                  :body-style="{ padding: '0px' }"
                  shadow="hover"
                  v-for="(item, index) in freeResourceList.list"
                >
                <div @click="handleProd(item.enid, item.class_type)">
                  <img :src="ossProcess(item.logo)" :alt="item.name" />
                  <div style="padding: 16px; text-align: left">
                    <span style="display: block">{{ item.name }}</span>
                    <p style="font-size: small">{{ item.intro }}</p>
                    <el-rate
                      v-model.number="item.score"
                      disabled
                      show-score
                      allow-half
                      size="small"
                      text-color="#ff6b00"
                    />
                  </div>
                  <div class="bottom">
                    <el-button text class="button">立即学习</el-button>
                  </div>
                </div>
                </el-card>
              </div>
            </el-scrollbar>

            <el-scrollbar v-if="item.type == 'class'">
              <div class="cards-cover">
                <el-card
                  :body-style="{ padding: '0px' }"
                  shadow="hover"
                  v-for="item in courseContentList.product_list"
                >
                <div @click="handleProd(item.product_enid, item.product_type)">
                  <img
                    :src="ossProcess(item.horizontal_image)"
                    :alt="item.title"
                  />
                  <div style="padding: 6px; text-align: left">
                    <span style="display: block">{{ item.title }}</span>
                    <p style="font-size: small; line-height: 0">
                      {{ item.intro }}
                    </p>
                    <el-rate
                      :model-value="handleScore(item.score)"
                      disabled
                      show-score
                      allow-half
                      size="small"
                      text-color="#ff6b00"
                    />
                    <p style="font-size: small; line-height: 0">
                      {{ item.learn_user_count }}人加入学习
                    </p>
                  </div>
                  <div class="bottom">
                    <el-button text class="button">立即学习</el-button>
                  </div>
                </div>
                </el-card>
              </div>
            </el-scrollbar>

            <el-scrollbar v-if="item.type == 'ebook'">
              <div class="ebook-cards-cover">
                <el-card
                  :body-style="{ padding: '0px' }"
                  shadow="hover"
                  v-for="(item, index) in ebookContentList.product_list"
                >
                  <div @click="handleProd(item.product_enid, item.product_type)">
                    <img
                      :src="ossEbookProcess(item.index_image)"
                      :alt="item.title"
                    />

                    <el-popover
                      placement="top-end"
                      :title="item.title"
                      :width="200"
                      trigger="hover"
                      :content="
                        ebookPopoverContent(item.intro, item.introduction)
                      "
                    >
                      <template #reference>
                        <div style="padding: 6px; text-align: left">
                          <span>{{ ebookTitle(item.title) }}</span>
                          <p style="font-size: small; line-height: 0">
                            {{ authorList(item.author_list) }}
                          </p>
                          <span
                            style="font-size: small"
                            v-if="item.user_score_count > 0"
                            >{{ item.user_score_count }}人评分</span
                          >
                          <span
                            style="font-size: small"
                            v-if="item.score.length == 0"
                            >暂无评分</span
                          >
                          <el-rate
                            :model-value="handleScore(item.score)"
                            disabled
                            show-score
                            allow-half
                            size="small"
                            text-color="#ff6b00"
                            v-if="item.score.length > 0"
                          />
                        </div>
                      </template>
                    </el-popover>
                  </div>

                </el-card>
              </div>
            </el-scrollbar>
            <el-button
              size="large"
              type="info"
              text
              v-if="item.type == 'class'"
              @click="gotoCategory(currentCourse, '')"
              >更多 {{ currentCourse.name }} 课程 ></el-button
            >
            <el-button
              size="large"
              type="info"
              text
              v-if="item.type == 'ebook'"
              @click="gotoCategory(currentEbook, '')"
              >更多 {{ currentEbook.name }} 电子书 ></el-button
            >
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
import { ElTable, ElMessage, ElDivider } from "element-plus";
import {
  GetHomeInitialState,
  SearchHot,
  SunflowerLabelList,
  SunflowerLabelContent,
  SunflowerResourceList,
  UserInfo,
  Logout,
} from "../../wailsjs/go/backend/App";
import { services } from "../../wailsjs/go/models";
import { BrowserOpenURL, WindowReloadApp } from "../../wailsjs/runtime";
import QrLogin from "../components/QrLogin.vue";
import EbookInfo from "../components/EbookInfo.vue";
import CourseInfo from "../components/CourseInfo.vue";
import { useRouter } from "vue-router";
import { Local } from "../utils/storage";

const router = useRouter();

const loading = ref(true);
const page = ref(0);
const total = ref(0);
const pageSize = ref(4);
const dialogVisible = ref(false);
const loginVisible = ref(false);
const ebookVisible = ref(false);
const courseVisible = ref(false);
const prodEnid = ref("");

const moreCategory = ref(9);
const idxEbookLabel = ref(0);
const idxCourseLabel = ref(0);

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
      SunflowerLabelContent("", 2, 0, 6)
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

const logout = async () => {
  await Logout()
    .then((result) => {
      console.log(result);
      WindowReloadApp();
      Local.remove("cookies")
      router.push("/home");
    })
    .catch((error) => {
      console.log(error);
    });
};

const handleProd = (enid: string, iType: number) => {
  prodEnid.value = enid;
  if (iType == 2) {
    ebookVisible.value = true;
  } else {
    // 4-专栏,36-大师课,66-class,22-course
    courseVisible.value = true;
  }
};

const handleLabel = (
  item: services.Navigation,
  index: number,
  nType: number
) => {
  if (nType == 2) {
    pageSize.value = 6;
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
  moreCategory.value = initial.homeData.categoryList.length;
};

const mouseleave = () => {
  // 添加一个小延时，避免菜单闪烁
  setTimeout(() => {
    moreCategory.value = 9;
  }, 100);
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

const gotoCategory = (item: any, label_id: string) => {
  let product_type = "0";
  if (item.navType == 2) {
    product_type = "2";
  } else if (item.navType == 4) {
    product_type = "66";
  }
  router.push({
    path: `/category`,
    query: {
        id: item.id,
        enid: item.enid,
        name: item.name,
        nav_type: item.navType,
        label_id: label_id,
        product_type: product_type,
    },
  });
};
</script>

<style scoped>
h1,
h4 {
  display: block;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0px;
  margin-inline-end: 0px;
  font-weight: bold;
}
.el-menu {
  border-right: 0;
  border-radius: 12px;
  overflow: hidden;
}

.classification {
  width: 160px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  background: #fff;
}

/* 主菜单项样式 */
.classification .el-sub-menu {
  position: relative;
  width: 100%;
  height: 48px;
  line-height: 48px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.classification .el-sub-menu:hover {
  background-color: #fff5f0;
}

/* 子菜单项样式 */
.el-menu .el-menu-item {
  position: relative;
  width: 100%;
  height: 40px;
  line-height: 40px;
  cursor: pointer;
  padding: 0 20px;
  color: #666;
  transition: all 0.3s ease;
}

.el-menu .el-menu-item:hover {
  background-color: #fff5f0 !important;
  color: #ff6b00;
}

.el-menu .el-menu-item.is-active {
  background-color: #ff6b00 !important;
  color: #fff;
}

/* 子菜单弹出层样式 */
.el-menu--popup {
  min-width: 180px;
  padding: 8px 0;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.el-menu--popup .el-menu-item {
  height: 40px;
  line-height: 40px;
  padding: 0 20px;
  margin: 0 8px;
  border-radius: 4px;
}

/* 菜单标题样式 */
.el-sub-menu__title {
  position: relative;
  padding-left: 20px !important;
  color: #333;
  font-weight: 500;
}

.el-sub-menu__title:hover {
  background-color: #fff5f0 !important;
  color: #ff6b00;
}

/* 菜单图标样式 */
.el-sub-menu__icon-arrow {
  color: #999;
  transition: all 0.3s ease;
}

.el-sub-menu:hover .el-sub-menu__icon-arrow {
  color: #ff6b00;
}

/* 更多菜单样式 */
.el-sub-menu[index="more"] .el-sub-menu__title {
  color: #666;
  font-size: 14px;
}

/* 优化菜单展开/收起动画 */
.el-menu-item,
.el-sub-menu__title {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 菜单分割线 */
.el-menu-item:not(:last-child) {
  position: relative;
}

.el-menu-item:not(:last-child)::after {
  content: '';
  position: absolute;
  left: 20px;
  right: 20px;
  bottom: 0;
  height: 1px;
  background: rgba(0, 0, 0, 0.03);
}

.el-carousel {
  height: 380px;
}
.el-carousel__item {
  height: 380px;
}
/* .el-scrollbar {
  height: 100px;
} */
.el-scrollbar__wrap {
  overflow-x: hidden;
}

/* 隐藏水平滚动条 */
.el-scrollbar__bar.is-horizontal {
  display: none !important;
}

.scrollbar-flex-content {
  display: flex;
  gap: 12px;
  padding: 4px 0;  /* 只保留上下内边距，移除左右内边距 */
}
.cards-cover {
  display: flex;
  gap: 20px;
}

.cards-cover .el-card {
  flex-shrink: 0;
  width: 290px;
  transition: all 0.3s ease;
  transform: translateZ(0);
  will-change: transform;
}

.cards-cover .el-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.cards-cover .el-card img {
  width: 100%;
  height: 163px;
  object-fit: cover;
  display: block;
}

.cards-cover .el-card .el-card__body {
  padding: 16px !important;
  height: 180px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.cards-cover .el-card .bottom {
  margin-top: 0;
  padding-top: 12px;
  border-top: 1px solid var(--el-border-color-lighter);
  display: flex;
  justify-content: center;
}

.cards-cover .el-card .button {
  padding: 8px 16px;
  transition: all 0.2s ease;
  width: 120px;
  text-align: center;
}

.cards-cover .el-card .button:hover,
.cards-cover .el-card .button:active,
.cards-cover .el-card .button:focus {
  padding: 8px 16px;
  width: 120px;
}

.cards-cover .el-card:active {
  transform: translateY(0);
}

.ebook-cards-cover {
  display: flex;
  /* border-radius: 8px; */
}
.ebook-cards-cover .el-card {
  flex-shrink: 0;
  width: 186px;
  margin-right: 20px;
  /* background: var(--el-color-info-light-9); */
  /* color: var(--el-color-info); */
}
.scrollbar-item {
  flex-shrink: 0;
  height: 32px;
  padding: 0 16px;
  border-radius: 4px;
  font-size: 14px;
  line-height: 32px;
  color: #666;
  background: #f5f5f5;
  transition: all 0.2s ease;
  min-width: 64px;
  text-align: center;
}
.scrollbar-item:hover {
  color: #ff6b00;
  background: #f7f7f7;
}
.active-btn {
  height: 32px;
  padding: 0 16px;
  border-radius: 4px;
  font-size: 14px;
  line-height: 32px;
  background: #ff6b00;
  color: #fff;
  min-width: 64px;
  text-align: center;
}
.active-btn:hover {
  color: #fff !important;
  background: #ff7b1a;
}
/* 覆盖 element-plus 的默认样式 */
.el-button.is-text {
  padding: 0 16px;
  margin: 0;
}
.el-button.is-text:not(.is-disabled):hover {
  background-color: var(--el-fill-color);
}

.module-title-wrap {
  text-align: left;
  padding-top: 20px;
  padding-bottom: 20px;
}
.module-cards-wrap {
  padding-top: 20px;
}

.bottom {
  margin-top: 13px;
  line-height: 12px;
  justify-content: space-between;
  align-items: center;
}

.el-row {
  margin-bottom: 20px;
}
.el-row:last-child {
  margin-bottom: 0;
}
.el-col {
  border-radius: 4px;
}

.home-banner {
  margin-bottom: 30px;
}

/* 左侧菜单样式 */
.el-menu {
  border-right: 0;
  border-radius: 12px;
  overflow: hidden;
}

.classification {
  width: 160px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  background: #fff;
}

.classification .el-sub-menu {
  position: relative;
  width: 100%;
  height: 48px;
  line-height: 48px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.classification .el-sub-menu:hover {
  background-color: #fff5f0;
}

/* 中间 banner 样式 */
.banner {
  .el-carousel {
    height: 380px;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    margin-left: 32px;
  }

  .el-carousel__item {
    height: 380px;
    overflow: hidden;
  }

  .el-image {
    width: 100%;
    height: 100%;
    transition: transform 0.3s ease;

    &:hover {
      transform: scale(1.02);
    }
  }

  .el-carousel__arrow {
    background-color: rgba(255, 255, 255, 0.8);
    border-radius: 50%;
    color: #666;
    
    &:hover {
      background-color: #fff;
      color: #ff6b00;
    }
  }
}

/* 右侧用户信息样式 */
.user {
  .not-login, .logged {
    background: #fff;
    border-radius: 12px;
    height: 380px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .not-login {
    .receive {
      width: 85%;
      height: 0;
      padding-bottom: calc(85% * (13 / 8));
      margin: 36px auto 0;
      background: url(https://piccdn2.umiwi.com/fe-oss/default/MTYzNTIzNTI3OTIw.png) no-repeat;
      background-size: contain;
      transition: transform 0.3s ease;

      &:hover {
        transform: scale(1.02);
      }
    }

    .login-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 160px;
      height: 44px;
      margin: 24px auto 0;
      border-radius: 22px;
      font-weight: 500;
      border: none;
      background-color: #ff6b00;
      color: #fff;
      transition: all 0.3s ease;

      &:hover {
        background-color: #ff7b1a;
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(255, 107, 0, 0.2);
      }
    }
  }

  .logged {
    padding: 30px 0;

    .personal {
      display: flex;
      flex-direction: column;
      align-items: center;
      margin-bottom: 30px;

      .el-avatar {
        border: 4px solid rgba(255, 107, 0, 0.1);
        transition: transform 0.3s ease;

        &:hover {
          transform: scale(1.05);
        }
      }

      h3 {
        margin: 16px 0 0;
        font-size: 18px;
        color: #333;
        font-weight: 500;
        text-align: center;
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .data {
      width: 180px;
      background: #fff;
      border-radius: 16px;
      padding: 20px 0;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
      background: url(https://piccdn2.umiwi.com/fe-oss/default/MTYzNTMwNDkxMjc1.png) no-repeat center bottom / 180px 40px;
      
      .time {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 16px;
        flex-direction: column;
        text-align: center;
        
        span {
          color: #666;
          font-size: 14px;
          margin-bottom: 4px;
        }
        
        em {
          font-style: normal;
          font-size: 28px;
          color: #ff6b00;
          font-weight: 500;
          line-height: 1.2;
        }

        &:last-child {
          margin-bottom: 0;
        }
      }

      .el-divider {
        margin: 20px 0;
        border-color: rgba(0, 0, 0, 0.06);
      }
    }
  }
}

/* 确保三列布局间距合适 */
.el-row {
  --el-row-gutter: 24px;
}

.el-col {
  padding: 0 calc(var(--el-row-gutter) / 2);
}
</style>
