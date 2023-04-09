<template>
  <div class="category">
    <div class="filters">
      <div class="filter-container filter-section-dash">
        <el-row :gutter="10">
          <el-col :span="2" style="font-size: larger">
            {{ filter.filter.product_types.title }}</el-col
          >
          <el-button
            text
            style="font-size: small"
            v-for="item in filter.filter.product_types.options"
            >{{ item.name }}</el-button
          >
        </el-row>
        <el-divider border-style="dashed" />
        <el-row :gutter="20">
          <el-col :span="2" style="font-size: larger">
            {{ filter.filter.navigations.title }}</el-col
          >
          <el-button
            text
            style="font-size: small"
            v-for="item in filter.filter.navigations.options"
            >{{ item.name }}</el-button
          >
        </el-row>
      </div>
      <!-- <div class="filter-container filter-section-dash">{{ filter }}</div> -->
      <div class="filter-container"></div>
    </div>

    <div class="category-source-container">
      <p style="text-align: left">已为你找到{{ filter.total }}个内容</p>
      <div class="sort-filter"></div>
      <el-divider />

      <div class="category-source-list">
        <el-row :gutter="20">
        <el-col :span="8" v-for="(item, index) in product.product_list" style="padding: 5px;">
          <el-card :body-style="{ padding: '10px', display: 'flex',}">
            <img :src="ossProcess(item.index_img)" class="source-container-cover" />
            <div >
              <div style="width: 290px; height: 120px; text-align: left">
                <span style="display: block">{{
                  productTitle(item.name, 15)
                }}</span>
                <span style="font-size: small;">
                  {{ productTitle(item.intro, 40) }}
                </span>
                <span style="font-size: small;display: block" v-if="item.product_type==66">
                  {{ item.lecturer_name_and_title }}
                </span>
                <span style="font-size: small;display: block" v-if="item.product_type==2">
                  {{ authorList(item.author_list) }}
                </span>
                <p v-if="item.product_type==66" style="font-size: small; line-height: 0;">
                  共{{ item.phase_num }}{{ item.price_desc }}
                </p>
                <el-rate
                  :model-value="handleScore(item.score)"
                  disabled
                  show-score
                  allow-half
                  size="small"
                  text-color="#ff6b00"
                  style="display: block;"
                />
                <span v-if="item.product_type==66" style="text-align: right;font-size: x-small;display: block;">
                  {{ item.learn_user_count }}人加入学习
                </span>
               
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      </div>
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

let filter = reactive(new services.AlgoFilterResp());
let product = reactive(new services.AlgoProductResp());

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
      param.classfc_name = name.value != "" ? name.value : "全部";
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

          getAlgoProduct(param)
            .then((list) => {
              console.log(list);
              Object.assign(product, list);
            })
            .catch((error) => {
              console.log(error);
            });
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
const getAlgoProduct = async (param: services.AlgoFilterParam) => {
  await AlgoProduct(param)
    .then((list) => {
      console.log(list);
      Object.assign(product, list);
    })
    .catch((error) => {
      console.log(error);
    });
};

const ossProcess = (url: string) => {
  return url + "?x-oss-process=image/resize,m_fill,h_224,w_168";
};

const handleScore = (score: string) => {
  return score != "" ? parseFloat(score) : 0;
};

const productTitle = (title: string, len:number) => {
  if (title.length <= len) {
    return title;
  } else {
    return title.substring(0, len).concat("...");
  }
};

const authorList = (authors: string[]) => {
  if (authors.length > 1) {
    return authors[0] + " 等";
  } else {
    return authors.toString();
  }
};
</script>

<style scoped>
.category-source-container .source-container-cover {
    height: 112px;
    min-width: 84px;
    width: 84px;
    position: relative;
    margin-right: 10px;
    background: var(--default-background);
    background-size: contain;
}
</style>
