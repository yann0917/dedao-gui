<template>
  <div class="category">
    <div class="filters">
      <div class="filter-container filter-section-dash">
        <el-row :gutter="10">
          <el-col :span="2" style="font-size: larger">
            {{ productTypes.title }}</el-col
          >
          <el-button
            text
            style="font-size: small"
            v-for="(item, index) in productTypes.options"
            @click="handleFilter(item, index, 1)"
            :class="
              productType == item.value || idxProd == index ? 'active-btn' : ''
            "
            >{{ item.name }}</el-button
          >
        </el-row>
        <el-divider border-style="dashed" />
        <el-row :gutter="20">
          <el-col :span="2" style="font-size: larger">
            {{ navigations.title }}</el-col
          >
          <el-button
            text
            style="font-size: small"
            v-for="(item, index) in navigations.options"
            :class="enid == item.value || idxLabel == index ? 'active-btn' : ''"
            @click="handleFilter(item, index, 2)"
            >{{ item.name }}</el-button
          >
        </el-row>
        <el-row :gutter="20">
          <el-col :span="2" style="font-size: larger">{{}}</el-col>
          <el-col :span="22">
            <div style="background-color: #f7f7f7; text-align: left">
              <el-button
                text
                style="font-size: small; padding: 7px"
                v-for="(item, index) in subOptions"
                :class="
                  labelId == item.value || idxSubLabel == index
                    ? 'active-btn'
                    : ''
                "
                @click="handleFilter(item, index, 3)"
                >{{ item.name }}</el-button
              >
            </div>
          </el-col>
        </el-row>
      </div>
      <div class="filter-container"></div>
    </div>

    <div class="category-source-container">
      <p style="text-align: left">已为你找到 <span style="color:#ff6b00">{{ filter.total }}</span> 个内容</p>
      <div class="sort-filter">
        <el-row>
          <el-col :span="1" v-for="item,index in sort.options">
            {{ item.name }}</el-col>
        </el-row>
      </div>
      <el-divider />

      <div class="category-source-list">
            <ul
              v-infinite-scroll="loadProduction"
              class="infinite-list"
              style="overflow-y: auto; max-height: 900px"
              :infinite-scroll-disabled="infLoadingProd"
              infinite-scroll-distance="10"
            >
            <el-space wrap>
              <li
                class="infinite-list-item"
                v-for="(item, index) in product_list.product_list"
              >
                <el-card
                  class="box-card"
                  shadow="hover"
                  :body-style="{ display: 'flex' }"
                  style="width: 400px"
                  @click="handleProd(item.id_out, item.item_type)"
                >
                  <img
                    :src="ossProcess(item.index_img)"
                    class="source-container-cover"
                  />
                  <div>
                    <div style="width: 290px; height: 120px; text-align: left">
                      <span style="display: block">{{
                        productTitle(item.name, 16)
                      }}</span>
                      <span style="font-size: small">
                        {{ productTitle(item.intro, item.item_type == 66? 20:40) }}
                      </span>
                      <span
                        style="font-size: small; display: block; color: gray"
                        v-if="item.item_type == 66"
                      >
                        {{ productTitle(item.lecturer_name_and_title, 20) }}
                      </span>
                      <span
                        style="font-size: small; display: block; color: gray"
                        v-if="item.item_type == 2"
                      >
                        {{ authorList(item.author_list) }}
                      </span>
                      <p
                        v-if="item.item_type == 66"
                        style="font-size: small; line-height: 0; color: gray"
                      >
                        共{{ item.phase_num }}{{ item.price_desc }}
                      </p>
                      <el-row>
                        <el-col :span="14">
                          <span
                            v-if="item.score == ''"
                            style="
                              text-align: left;
                              font-size: small;
                              color: gray;
                           "
                          >
                            暂无评分
                          </span>

                          <el-rate
                            v-else
                            :model-value="handleScore(item.score)"
                            disabled
                            show-score
                            allow-half
                            size="small"
                            text-color="#ff6b00"
                            style="display: block"
                          />
                        </el-col>
                        <el-col :span="2"></el-col>
                        <el-col :span="8">
                          <span
                            v-if="item.learn_user_count > 0"
                            style="font-weight: 200;
                            font-size:70%;
                            text-align: right;
                            color: gray;
                            "
                          >
                            {{ item.learn_user_count }}人加入学习
                          </span>
                        </el-col>
                      </el-row>
                    </div>
                  </div>
                </el-card>
              </li>
            </el-space>
            </ul>
      </div>
    </div>
  </div>
    <EbookInfo v-if="ebookVisible" :enid= "prodEnid" :dialog-visible="ebookVisible" @close="closeDialog"></EbookInfo>
</template>

<script lang="ts" setup>
import { ref, reactive, onBeforeMount, onMounted, watch} from "vue";
import { ElTable, ElMessage } from "element-plus";
import {
  SearchHot,
  AlgoFilter,
  AlgoProduct,
} from "../../wailsjs/go/backend/App";
import { services } from "../../wailsjs/go/models";
import { useRoute, useRouter } from "vue-router";
import  EbookInfo  from '../components/EbookInfo.vue';

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const page = ref(0);
const total = ref(0);
const pageSize = ref(4);

const ebookVisible = ref(false);
const courseVisible = ref(false);
const prodEnid = ref("");

const idxProd = ref(0);
const idxLabel = ref(0);
const idxSubLabel = ref(0);

const infLoadingProd = ref(true);

let filter = reactive(new services.AlgoFilterResp());
let product = reactive(new services.AlgoProductResp());
let product_list = reactive(new services.AlgoProductResp());
product_list.product_list = [];

let productTypes = reactive(new services.Strategy());
let navigations = reactive(new services.Strategy());
let sort = reactive(new services.Strategy());
let subOptions = reactive([] as services.Option[]);

const enid = ref();
const name = ref();
const navType = ref();
const productType = ref();
const labelId = ref();

onMounted(() => {
  let param = new services.AlgoFilterParam();
  watch(
    () => {
      enid.value = route.query.enid;
      name.value = route.query.name;
      navType.value = route.query.nav_type;
      labelId.value = route.query.label_id;
      productType.value = route.query.product_type;

      param.request_id = "";
      param.tags_ids = [];
      param.page = 0;

      param.navigation_id = enid.value;
      param.classfc_name = name.value != "" ? name.value : "全部";
      param.label_id = labelId.value;
      param.page_size = 18;
      param.product_types = productType.value;
      param.sort_strategy = "HOT";
    },
    () =>
      getAlgoFilter(param)
        .then((list) => {
          console.log(list);
          Object.assign(filter, list);

          Object.assign(productTypes, filter.filter.product_types);
          Object.assign(navigations, filter.filter.navigations);
          Object.assign(sort, filter.filter.sort_strategy)
          // navigations.options.forEach((item) => {
          //   if (item.value == enid.value) {
          //     if (
          //       item.sub_options != undefined &&
          //       item.sub_options?.length > 0
          //     ) {
                
          //     }
          //   }
          // });

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
const param = new services.AlgoFilterParam();
const handleFilter = (item: services.Option, idx: number, nType: number) => {
  if (nType == 1) {
    productType.value = item.value;
    idxProd.value = idx;

    param.product_types = item.value;
    param.classfc_name = "全部";
    param.navigation_id = "";
    param.label_id = "";
  } else if (nType == 2) {
    enid.value = item.value;
    idxLabel.value = idx;
    idxSubLabel.value = 0;
    param.label_id = "";

    subOptions.splice(0, subOptions.length)
    if (item.sub_options != undefined && item.sub_options?.length > 0) {
      let subOption = JSON.parse(JSON.stringify(item.sub_options));
      subOption.forEach((item)=>{
        subOptions.push(item)
      })
    }

    param.classfc_name = item.name;
    param.navigation_id = item.value;
  } else if (nType == 3) {
    idxSubLabel.value = idx;
    param.label_id = item.value;
  }
  param.request_id = "";
  param.tags_ids = [];
  param.page = 0;
  param.nav_type = 0;
  param.page_size = 18;
  param.sort_strategy = "HOT";
  // console.log(param);
  product_list.product_list = [];
  getAlgoFilter(param);
};

const getAlgoFilter = async (param: services.AlgoFilterParam) => {
  await AlgoFilter(param)
    .then((list) => {
      console.log(list);
      Object.assign(filter, list);

      Object.assign(productTypes, filter.filter.product_types);
      Object.assign(navigations, filter.filter.navigations);
      Object.assign(sort, filter.filter.sort_strategy)
    })
    .catch((error) => {
      console.log(error);
    });
  getAlgoProduct(param);
};

const getAlgoProduct = async (param: services.AlgoFilterParam) => {
  await AlgoProduct(param)
    .then((list) => {
      console.log(list);
      Object.assign(product, list);
      product_list.is_more = product.is_more;
      product_list.total = product.total;
      product_list.product_list.push(...product.product_list);
    })
    .catch((error) => {
      console.log(error);
    });
  infLoadingProd.value = false;
};

const ossProcess = (url: string) => {
  return url + "?x-oss-process=image/resize,m_fill,h_224,w_168";
};

const handleScore = (score: string) => {
  return score != "" ? parseFloat(score) : 0;
};

const productTitle = (title: string, len: number) => {
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

const loadProduction = () => {
  infLoadingProd.value = true;
  if (product.is_more === 1) {
    page.value += 1;
    param.page = page.value;
    getAlgoProduct(param);
  }
};

const handleProd = (enid:string, iType:number)=>{
  prodEnid.value = enid
  if (iType == 2) {
    ebookVisible.value = true
  }
  console.log(prodEnid)
}

const openDialog = () => {
    ebookVisible.value = true
}
const closeDialog = () => {
  ebookVisible.value = false
}

</script>

<style scoped lang="scss">
.category-source-container .source-container-cover {
  height: 112px;
  min-width: 84px;
  width: 84px;
  position: relative;
  margin-right: 10px;
  background: var(--default-background);
  background-size: contain;
}
.active-btn {
  font-weight: 500;
  color: #ff6b00;
}

ul {
  padding: 0;
  margin: 0;
  text-align: left;
  list-style: none;
  display: flex;
}
ul li {
  flex-shrink: 0;
}
.infinite-list {
  width: 100%;
  height: 100%;
  padding: 0;
  margin: 0;
  text-align: left;
  list-style: none;
  flex-wrap: wrap;
}

.infinite-list .infinite-list-item:last-child {
  margin-bottom: 0;
}

.infinite-list .infinite-list-item + .list-item {
  margin-top: 10px;
}

.el-card {
  --el-card-padding: 10px;
}
</style>
