<template>
  <el-row class="tac">
    <el-col :span="4">
      <el-menu
        default-active=""
        class="classification"
        collapse="true"
        router="true"
        active-text-color="#ffd04b"
        @open="handleOpen"
        @close="handleClose"
      >
        <el-sub-menu
          :index="item.enid"
          v-for="item in initial.homeData.categoryList"
        >
          <template #title>{{ item.name }}</template>
          <el-menu-item :index="i.enid" v-for="i in item.labelList">
            <template #title>{{ i.name }}</template>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-col>

    <el-col :span="16" class="banner">
      <el-carousel :interval="5000" arrow="always">
        <el-carousel-item v-for="item in initial.homeData.banner" :key="item">
          <el-image :src="item.img">
            <a :href="item.url"></a>
          </el-image>
        </el-carousel-item>
      </el-carousel>
    </el-col>
    <el-col :span="4" class="user"></el-col>
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
          <el-scrollbar v-if="item.type == 'class'">
            <div class="scrollbar-flex-content">
              <el-button
                v-for="item in courseLabelList.list"
                :key="item.id"
                class="scrollbar-item"
                text
                @click="handleLabel(item.enid,4)"
              >
                {{ item.name }}
              </el-button>
            </div>
          </el-scrollbar>

          <el-scrollbar v-if="item.type == 'ebook'">
            <div class="scrollbar-flex-content">
              <el-button
                v-for="item in ebookLabelList.list"
                :key="item.id"
                class="scrollbar-item"
                text
                @click="handleLabel(item.enid,2)"
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
                  <img :src="ossProcess(item.logo)" :alt="item.name" />
                  <div style="padding: 14px">
                    <span style="display: block; text-align: left;">{{ item.name }}</span>
                    <el-alert type="info" :closable="false">
                      <p>{{ item.intro }}</p>
                      <el-tag class="ml-2" type="warning">{{
                        item.score
                      }}</el-tag></el-alert>
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
                  v-for="(item, index) in courseContentList.product_list"
                >
                  <img
                    :src="ossProcess(item.horizontal_image)"
                    :alt="item.title"
                  />
                  <div style="padding: 10px">
                    <span style="display: block; text-align: left;">{{ item.title }}</span>
                    <el-alert type="info" :closable="false">
                      <span>{{ item.intro }}<br/></span>
                      <span>{{ item.learn_user_count }}人加入学习</span>
                      <el-tag class="ml-2" type="warning">{{
                        item.score
                      }}</el-tag></el-alert
                    >
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
                  <img :src="ossProcess(item.index_image)" :alt="item.title" />

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
                      <div style="padding: 10px; text-align: left">
                        <span style=" ;">{{ ebookTitle(item.title) }}</span>
                        <el-alert type="info" :closable="false">
                          <p>{{ item.author_list.toString() }}</p>
                          <span v-if="item.user_score_count>0">{{ item.user_score_count }}人评分</span>
                          <span v-if="item.score.length==0">暂无评分</span>
                          <el-tag class="ml-2" type="warning" v-if="item.score.length>0">{{
                            item.score
                          }}</el-tag>
                          </el-alert
                        >
                      </div>
                    </template>
                  </el-popover>
                </el-card>
              </div>
            </el-scrollbar>
          </div>

          <el-button
            size="large"
            key="course"
            type=""
            text
            v-if="item.type == 'class'"
            >更多 视频课 课程 ></el-button
          >
          <el-button
            size="large"
            key="ebook"
            type=""
            text
            v-if="item.type == 'ebook'"
            >更多 小说 电子书 ></el-button
          >
        </div>
      </div>
    </div>
  </div>

  <el-backtop :right="100" :bottom="100" />
</template>

<script lang="ts" setup>
import { ref, reactive, onBeforeMount, onMounted } from "vue";
import { ElTable, ElMessage } from "element-plus";
import {
  GetHomeInitialState,
  SearchHot,
  SunflowerLabelList,
  SunflowerLabelContent,
  SunflowerResourceList,
} from "../../wailsjs/go/backend/App";
import { services } from "../../wailsjs/go/models";
import { useRouter } from "vue-router";

const router = useRouter();

const loading = ref(true);
const page = ref(0);
const total = ref(0);
const pageSize = ref(4);
const dialogVisible = ref(false);

let initial = reactive(new services.HomeInitState());
let ebookLabelList = reactive(new services.SunflowerLabelList());
let courseLabelList = reactive(new services.SunflowerLabelList());

let freeResourceList = reactive(new services.SunflowerResourceList());
let ebookContentList = reactive(new services.SunflowerContent());
let courseContentList = reactive(new services.SunflowerContent());

onBeforeMount(() => {
  // 分类
  GetHomeInitialState()
    .then((state) => {
      Object.assign(initial, state);
      console.log(state);
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
      SunflowerLabelContent("", 2, 0, 6)
        .then((list) => {
          console.log(list);
          Object.assign(ebookContentList, list);
        })
        .catch((error) => {
          console.log(error);
        });
      // console.log(result);
    })
    .catch((error) => {
      console.log(error);
    }),
    // 精选课程
    SunflowerLabelList(4)
      .then((result) => {
        Object.assign(courseLabelList, result);
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

const handleLabel = (enid:string,nType: number) => {
  if(nType==2) {
    pageSize.value = 6
  }
  if(nType==4) {
    pageSize.value = 4
  }
  sunflowerLabelContent(enid, nType)
};

const sunflowerLabelContent = async (enid:string  ,nType:number) => {
  await SunflowerLabelContent(enid, nType, page.value, pageSize.value)
    .then((list) => {
      if(nType==2) {
        Object.assign(ebookContentList, list);
    } else if (nType==4) {
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

const ebookTitle = (title: string) => {
  if (title.length <= 9) {
    return title;
  } else {
    return title.substring(0, 9).concat("...");
  }
};

const ebookPopoverContent = (intro: string, introduction: string) => {
  if (intro.concat(introduction).length <= 146) {
    return intro.concat(introduction);
  } else {
    return intro.concat(introduction).substring(0, 146).concat("...");
  }
};

const openDialog = () => {
  dialogVisible.value = true;
};
const closeDialog = () => {
  //   initForm()
  dialogVisible.value = false;
};
const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath);
};
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath);
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
  width: 180px;
  /* height: 376px; */
  box-shadow: 0 5px 10px rgba(51, 51, 51, 0.06);
}

.classification .el-sub-menu {
  position: relative;
  width: 100%;
  height: 38px;
  line-height: 38px;
  cursor: pointer;
  z-index: 300;
}
.el-carousel {
  height: 600px;
}
/* .el-scrollbar {
  height: 100px;
} */
.scrollbar-flex-content {
  display: flex;
  border-radius: 8px;
}
.cards-cover {
  display: flex;
  /* border-radius: 8px; */
}

.cards-cover .el-card {
  flex-shrink: 0;
  width: 310px;
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
  width: 200px;
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
  height: 50px;
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
</style>
