<template>
  <el-row :gutter="30" class="home-banner">
    <el-col :span="4">
      <el-menu
        default-active=""
        class="classification"
        :collapse="true"
        :router="true"
        active-text-color="#ffd04b"
        @open="handleOpen"
        @close="handleClose"
        @mouseleave="mouseleave"
      >
        <el-sub-menu
          :index="item.enid"
          v-for="(item, index) in initial.homeData.categoryList"
          v-show="index < moreCategory"
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
    <el-col :span="4" class="user">
      <div :class="!initial.isLogin ? 'not-login' : 'logged'">
        <div v-if="!initial.isLogin">
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
              <span>今日学习</span
              ><span
                ><em style="font-size: 22px">{{
                  (user.today_study_time / 60).toFixed(0)
                }}</em>
                分钟</span
              >
            </p>
            <el-divider border-style="dotted" />
            <p class="time">
              <span>连续学习</span
              ><span
                ><em style="font-size: 22px">{{ user.study_serial_days }}</em>
                天</span
              >
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
  moreCategory.value = 9;
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
const handleOpen = (key: string, keyPath: string[]) => {
  // console.log(key, keyPath);
};
const handleClose = (key: string, keyPath: string[]) => {
  // console.log(key, keyPath);
};

const gotoCategory = (item: any, label_id: string) => {
  // services.Navigation
  console.log(item);
  let product_type = "0";
  if (item.nav_type == 2) {
    product_type = "2";
  } else if (item.nav_type == 4) {
    product_type = "66";
  }
  router.push({
    path: `/category`,
    query: {
        id: item.id,
        enid: item.enid,
        name: item.name,
        nav_type: item.nav_type,
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
}
.classification {
  width: 190px;
  box-shadow: 0 5px 10px rgba(51, 51, 51, 0.06);
}

.classification .el-sub-menu {
  position: relative;
  width: 100%;
  height: 38px;
  line-height: 38px;
  cursor: pointer;
}

/* .el-menu .el-sub-menu:hover {
  background-color: rgba(255, 107, 0, 0.06) !important;
} */
.el-menu .el-menu-item {
  position: relative;
  width: 100%;
  height: 38px;
  line-height: 38px;
  cursor: pointer;
}

.el-menu-item:hover {
  background-color: rgba(255, 107, 0, 0.06) !important;
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
.scrollbar-flex-content {
  display: flex;
  /* border-radius: 8px; */
  /* background: var(--el-color-info-light-9); */
}
.cards-cover {
  display: flex;
  /* border-radius: 8px; */
}

.cards-cover .el-card {
  flex-shrink: 0;
  width: 290px;
  margin-right: 20px;
  /* background: var(--el-color-info-light-9); */
  /* color: var(--el-color-info); */
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
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 40px;
  text-align: center;
  background: var(--el-color-info-light-9);
  color: var(--el-color-info);
}
.el-scrollbar__wrap {
  overflow-x: hidden;
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
.active-btn {
  font-weight: 500;
  height: 40px;
  color: #fff;
  border-radius: 8px;
  background: #ff6b00;
  box-shadow: 0px 6px 10px rgba(251, 72, 16, 0.2);
}
.el-button.is-text:not(.is-disabled):hover {
  background-color: var(--el-fill-color);
}

.home-banner .user {
  box-shadow: 0 5px 10px rgba(51, 51, 51, 0.06);
  height: 380px;
}

.home-banner .user .not-login .receive {
  width: 89%;
  height: 0;
  padding-bottom: calc(89% * (13 / 8));
  margin: 36px auto 0;
  background: url(https://piccdn2.umiwi.com/fe-oss/default/MTYzNTIzNTI3OTIw.png)
    no-repeat;
  background-size: contain;
}

.home-banner .user .not-login .login-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 148px;
  margin: 14px auto 0;
  border-radius: 8px;
  line-height: 28px;
  font-weight: bold;
  border-color: #ff6b00;
  background-color: #ff6b00;
  color: #fff;
}

.home-banner .user .not-login .login-btn .line {
  display: inline-block;
  width: 1px;
  height: 14px;
  margin: 0 8px;
}

.home-banner .user .logged .personal {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.home-banner .user .logged .personal .nickname {
  width: 100%;
  height: 24px;
  margin-top: 4px;
  line-height: 24px;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  font-size: 16px;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

.home-banner .user .logged .data {
  width: 148px;
  height: 144px;
  margin-top: 8px;
  padding: 6px 12px 0;
  border-radius: 10px;
  background: url(https://piccdn2.umiwi.com/fe-oss/default/MTYzNTMwNDkxMjc1.png)
    no-repeat center bottom / 148px 40px;
  box-shadow: 0 5px 10px rgba(255, 107, 0, 0.1);
}

.home-banner .user .logged .data .time {
  display: flex;
  padding: 2px 0;
  line-height: 22px;
  justify-content: space-between;
}

.home-banner .user .logged .button {
  height: 32px;
  margin-top: 16px;
  border-radius: 8px;
  font-weight: bold;
  border-color: #ff6b00;
  background-color: #ff6b00;
  color: #fff;
}

.home-banner .banner .el-carousel {
    box-shadow: 0 5px 10px rgba(51, 51, 51, 0.06);
    /*height: 380px;*/
    border-radius: 8px;
}

</style>
