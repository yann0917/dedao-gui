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
      <div v-if="item.isShow == 3 && (item.type =='free_class'||item.type =='class'||item.type =='ebook')">
        <div class="module-title-wrap">
          <h1>{{ item.title }}</h1>
          <p>{{ item.description }}</p>
        </div>

        <div class="home-container">

          <el-scrollbar v-if="item.type=='class'">
            <div class="scrollbar-flex-content">
              <el-button
                v-for="item in courseLabelList.list"
                :key="item.id"
                class="scrollbar-item"
                text
              >
                {{ item.name }}
              </el-button>
            </div>
          </el-scrollbar>

          <el-scrollbar v-if="item.type=='ebook'">
            <div class="scrollbar-flex-content">
              <el-button
                v-for="item in ebookLabelList.list"
                :key="item.id"
                class="scrollbar-item"
                text
              >
                {{ item.name }}
              </el-button>
            </div>
          </el-scrollbar>

          <el-button size="large" key="course" type="" text
            v-if="item.type=='class'" >更多 视频课 课程 ></el-button
          >
          <el-button size="large" key="ebook" type="" text
            v-if="item.type=='ebook'" >更多 小说 电子书 ></el-button
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
import Pagination from "../components/Pagination.vue";
import { Local } from "../utils/storage";

const router = useRouter();

const loading = ref(true);
const page = ref(1);
const total = ref(0);
const pageSize = ref(15);
const dialogVisible = ref(false);

let initial = reactive(new services.HomeInitState());
let ebookLabelList = reactive(new services.SunflowerLabelList());
let courseLabelList = reactive(new services.SunflowerLabelList());

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
  SearchHot()
    .then((result) => {
      console.log(result);
    })
    .catch((error) => {
      console.log(error);
    }),
    // 电子书
    SunflowerLabelList(2)
      .then((result) => {
        Object.assign(ebookLabelList, result);
        console.log(result);
      })
      .catch((error) => {
        console.log(error);
      }),
    // 精选课程
    SunflowerLabelList(4)
      .then((result) => {
        Object.assign(courseLabelList, result);
        // console.log(result)
      })
      .catch((error) => {
        console.log(error);
      });
});

// 分页
const handleChangePage = (item: any) => {
  page.value = item.page;
  pageSize.value = item.pageSize;
  // getTableData()
};

// const getTableData = async () => {
//     await CourseList("compass", "study", page.value, pageSize.value).then((table) => {
//         loading.value = false
//         Object.assign(tableData, table)
//         console.log(tableData)
//     }).catch((error) => {
//         loading.value = false
//         ElMessage({
//             message: error,
//             type: 'warning'
//         })
//     })
// }

// getTableData()

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
h1 {
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
.el-scrollbar {
  height: 100px;
}
.scrollbar-flex-content {
  display: flex;
  border-radius: 8px;
}
.scrollbar-item {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 50px;
  /* margin-right: 2px; */
  text-align: center;
  /* border-radius: 4px; */
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
</style>
