<template>
  <div class="category">
    <div class="filters">
        <div class="filter-container filter-section-dash">
            <el-row :gutter="10">
                <el-col :span="2" style="font-size: larger;"> {{ filter.filter.product_types.title  }}</el-col>
                <el-col :span="2" style="font-size: small;" v-for="item in filter.filter.product_types.options">{{ item.name }}</el-col>
            </el-row>
            <el-divider border-style="dashed"/>
            <el-row :gutter="10">
                <el-col :span="2" style="font-size: larger;"> {{ filter.filter.navigations.title  }}</el-col>
                <el-col :span="2"  style="font-size: small;" v-for="item in filter.filter.navigations.options">{{ item.name }}</el-col>
            </el-row>
        </div>
      <!-- <div class="filter-container filter-section-dash">{{ filter }}</div> -->
      <div class="filter-container"></div>
    </div>

    <div class="category-source-container">
        <p style="text-align: left;"> 已为你找到{{ filter.total }}个内容</p>
      <div class="sort-filter"></div>
      <div class="category-source-list"></div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onBeforeMount, onMounted, watch } from "vue";
import { ElTable, ElMessage } from "element-plus";
import {
  SearchHot,
  AlgoFilter,
  AlgoProduct,
} from "../../wailsjs/go/backend/App";
import { services } from "../../wailsjs/go/models";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const page = ref(0);
const total = ref(0);
const pageSize = ref(4);
const dialogVisible = ref(false);

const moreCategory = ref(9);
const currentCourse = ref("");
const currentEbook = ref("");

let ebookLabelList = reactive(new services.SunflowerLabelList());
let courseLabelList = reactive(new services.SunflowerLabelList());

let freeResourceList = reactive(new services.SunflowerResourceList());
let ebookContentList = reactive(new services.SunflowerContent());
let courseContentList = reactive(new services.SunflowerContent());
let filter = reactive(new services.AlgoFilterResp());

let enid = ref();
let name = ref();
let navType = ref();
let labelId = ref();

onMounted(() => {
  console.log(route.query);
  console.log(route.params);
  let param = new services.AlgoFilterParam();
  watch(
    () => {
      enid.value = route.query.enid;
      name.value = route.query.name;
      navType.value = route.query.nav_type;
      labelId.value = route.query.label_id;

      param.navigation_id = enid.value;
      param.classfc_name = name.value;
      param.label_id = labelId.value;
      param.page_size = 18;
      param.product_types = "0";

      if (navType.value == 2) {
        param.product_types = "2";
      } else if (navType.value == 4) {
        param.product_types = "66";
      }
      param.sort_strategy = "HOT";
    },
    () =>
      getAlgoFilter(param)
        .then((list) => {
          console.log(list);
          Object.assign(filter, list);
        })
        .catch((error) => {
          console.log(error);
        }),
    { immediate: true }
  );
  // 热搜词
  // SearchHot()
  //   .then((result) => {
  //     console.log(result);
  //   })
  //   .catch((error) => {
  //     console.log(error);
  //   })
});

const getAlgoFilter = async (param: services.AlgoFilterParam) => {
  await AlgoFilter(param)
    .then((list) => {
      console.log(list);
      Object.assign(filter, list);
    })
    .catch((error) => {
      console.log(error);
    });
};
</script>

<style scoped></style>
