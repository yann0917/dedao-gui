<template>
  <div class="category">
    <div class="filters">
      <div class="filter-container">
        <el-row class="filter-row" :gutter="24">
          <el-col :span="2" class="filter-label">
            {{ productTypes.title }}
          </el-col>
          <el-col :span="22" class="filter-options">
            <el-button
              v-for="(item, index) in productTypes.options"
              :key="index"
              :class="['filter-btn', idxProd === index ? 'active-btn' : '']"
              text
              @click="handleFilter(item, index, 1)"
            >
              {{ item.name }}
            </el-button>
          </el-col>
        </el-row>

        <el-divider class="filter-divider" />

        <el-row class="filter-row" :gutter="24">
          <el-col :span="2" class="filter-label">
            {{ navigations.title }}
          </el-col>
          <el-col :span="22" class="filter-options">
            <el-button
              v-for="(item, index) in navigations.options"
              :key="index"
              :class="['filter-btn', idxLabel === index ? 'active-btn' : '']"
              text
              @click="handleFilter(item, index, 2)"
            >
              {{ item.name }}
            </el-button>
          </el-col>
        </el-row>

        <el-row v-if="subOptions.length > 0" class="filter-row" :gutter="24">
          <el-col :span="2" class="filter-label"></el-col>
          <el-col :span="22">
            <div class="sub-options">
              <el-button
                v-for="(item, index) in subOptions"
                :key="index"
                :class="['filter-btn', idxSubLabel === index ? 'active-btn' : '']"
                text
                @click="handleFilter(item, index, 3)"
              >
                {{ item.name }}
              </el-button>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <div class="content-container">
      <div class="result-header">
        <div class="result-count">
          已为你找到 <span class="highlight">{{ filter.total }}</span> 个内容
        </div>
        <div class="sort-options">
          <el-button
            v-for="item in sort.options"
            :key="item.value"
            text
            class="sort-btn"
            :class="{ 'active-sort': param.sort_strategy === item.value }"
            @click="handleSort(item)"
          >
            {{ item.name }}
          </el-button>
        </div>
      </div>

      <el-divider class="content-divider" />

      <div class="content-list">
        <ul
          v-infinite-scroll="loadProduction"
          class="infinite-list"
          :infinite-scroll-disabled="infLoadingProd"
          infinite-scroll-distance="10"
        >
          <li
            v-for="(item, index) in product_list.product_list"
            :key="index"
            class="content-item"
          >
            <el-card
              class="content-card"
              shadow="hover"
              @click="handleProd(item.id_out, item.item_type, item.product_type)"
            >
              <div class="card-content">
                <img
                  :src="ossProcess(item.index_img)"
                  class="content-cover"
                  :alt="item.name"
                />
                <div class="content-info">
                  <h3 class="content-title">{{ productTitle(item.name, 16) }}</h3>
                  <p class="content-intro">{{ productTitle(item.intro, item.item_type === 66 ? 20 : 40) }}</p>
                  
                  <div v-if="item.item_type === 66" class="meta-info">
                    <span class="lecturer">{{ productTitle(item.lecturer_name_and_title, 20) }}</span>
                    <span class="episode-count">共{{ item.phase_num }}{{ item.price_desc }}</span>
                  </div>
                  
                  <div v-if="item.item_type === 2" class="meta-info">
                    <span class="author">{{ authorList(item.author_list) }}</span>
                  </div>

                  <div class="rating-container">
                    <div class="rating">
                      <el-rate
                        v-if="item.score !== ''"
                        :model-value="handleScore(item.score)"
                        disabled
                        show-score
                        allow-half
                        size="small"
                        text-color="var(--accent-color)"
                      />
                      <span v-else class="no-rating">暂无评分</span>
                    </div>
                    <span v-if="item.learn_user_count > 0" class="learner-count">
                      {{ item.learn_user_count }}人加入学习
                    </span>
                  </div>
                </div>
              </div>
              <div class="card-actions">
                <!-- 电子书加入书架按钮 -->
                <el-button v-if="item.item_type === 2 && !item.is_on_bookshelf" 
                           @click="ebookShelfAdd(item.id_out, $event)" 
                           size="small" 
                           type="primary">
                  <el-icon><Plus /></el-icon>加入书架
                </el-button>
                <el-button v-if="item.item_type === 2 && item.is_on_bookshelf" 
                           @click="ebookShelfRemove(item.id_out, $event)" 
                           size="small" 
                           type="danger">
                  <el-icon><Delete /></el-icon>移出书架
                </el-button>
              </div>
            </el-card>
          </li>
        </ul>
      </div>
    </div>

    <EbookInfo v-if="ebookVisible" :enid="prodEnid" :dialog-visible="ebookVisible" @close="closeDialog" @bookshelf-changed="refreshAlgoData" />
    <CourseInfo v-if="courseVisible" :enid="prodEnid" :dialog-visible="courseVisible" @close="closeDialog" />
    <AudioInfo v-if="audioVisible" :enid="prodEnid" :dialog-visible="audioVisible" @close="closeDialog" />
    <OutsideInfo v-if="outsideVisible" :enid="prodEnid" :dialog-visible="outsideVisible" @close="closeDialog" />
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onBeforeMount, onMounted, watch} from "vue";
import { ElTable, ElMessage, ElMessageBox } from "element-plus";
import {
  SearchHot,
  AlgoFilter,
  AlgoProduct,
  EbookShelfAdd,
  EbookShelfRemove,
  OutsideDetail,
} from "../../wailsjs/go/backend/App";
import { services } from "../../wailsjs/go/models";
import { useRoute, useRouter } from "vue-router";
import EbookInfo  from '../components/EbookInfo.vue';
import CourseInfo from "../components/CourseInfo.vue";
import AudioInfo from "../components/AudioInfo.vue";
import OutsideInfo from "../components/OutsideInfo.vue";
import { Plus, Delete } from '@element-plus/icons-vue';

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const page = ref(0);
const total = ref(0);
const pageSize = ref(4);

const ebookVisible = ref(false);
const courseVisible = ref(false);
const audioVisible = ref(false);
const outsideVisible = ref(false);
const prodEnid = ref("");

const idxProd = ref(0);
const idxLabel = ref(0);
const idxSubLabel = ref(0);
const infLoadingProd = ref(false);

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
let param = new services.AlgoFilterParam();

onMounted(() => {
    enid.value = route.query.enid;
    name.value = route.query.name;
    navType.value = route.query.nav_type;
    labelId.value = route.query.label_id;
    productType.value = route.query.product_type;
    param.request_id = "";
    param.tags_ids = [];
    param.page = 0;

    param.nav_type = navType.value?Number.parseInt(navType.value):0;
    param.navigation_id = enid.value;
    param.classfc_name = name.value != "" ? name.value : "全部";
    param.label_id = labelId.value;
    param.page_size = 18;
    param.product_types = productType.value;
    param.sort_strategy = "HOT";
    getAlgoFilter(param)
        .then(() => {
            productTypes.options?.forEach((item,index)=> {
                if (item.value == productType.value) {
                    idxProd.value = index;
                }
            })
            navigations.options?.forEach((item, index) => {
                if (item.value == enid.value) {
                    idxLabel.value = index;
                    if (item.sub_options != undefined && item.sub_options?.length > 0) {
                        let subOption = JSON.parse(JSON.stringify(item.sub_options));
                        subOption.forEach((item,index)=>{
                            if (item.value == labelId.value) {
                                idxSubLabel.value = index;
                            }
                            subOptions.push(item)
                        })
                    }
                }
            });
            getAlgoProduct(param);
        })
        .catch((error) => {
            console.log(error);
        })

  // 热搜词
  // SearchHot()
  //   .then((result) => {
  //     console.log(result);
  //   })
  //   .catch((error) => {
  //     console.log(error);
  //   })
});
const handleFilter = (item: services.Option, idx: number, nType: number) => {
    if (nType == 1) {
        productType.value = item.value;
        idxProd.value = idx;

        param.product_types = item.value;
        param.classfc_name = item.name;
        param.navigation_id = "";
        param.label_id = "";
        subOptions.splice(0, subOptions.length)
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
    product_list.product_list = [];
    getAlgoFilter(param);
    getAlgoProduct(param);
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
};

const getAlgoProduct = async (param: services.AlgoFilterParam) => {
  loading.value = true;
  try {
    const list = await AlgoProduct(param);
    Object.assign(product, list);
    product_list.is_more = product.is_more;
    product_list.total = product.total;
    product_list.product_list.push(...product.product_list);
    console.log(product_list.product_list);
  } catch (error) {
    console.log(error);
    ElMessage({
      message: '获取数据失败',
      type: 'error'
    });
  } finally {
    loading.value = false;
    infLoadingProd.value = false;
  }
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

const handleProd = (enid:string, iType:number, pType:number)=>{
  prodEnid.value = enid
  if (iType == 2) {
    ebookVisible.value = true
  } else if(iType == 66) {
    courseVisible.value = true
  } else if(iType == 13 && pType == 1013) {
    // 13-单本, 1013-名家讲书合集
    outsideVisible.value = true
  } else if(iType == 13 && pType == 13) {
    audioVisible.value = true
  }
}

const openDialog = () => {
    ebookVisible.value = true
}
const closeDialog = () => {
  ebookVisible.value = false
  courseVisible.value = false
  audioVisible.value = false
  outsideVisible.value = false
}

const handleSort = (item: services.Option) => {
  param.sort_strategy = item.value;
  param.page = 0;
  page.value = 0;
  product_list.product_list = [];
  getAlgoProduct(param);
};

const ebookShelfAdd = async (enid: string, event?: Event) => {
  // 阻止事件冒泡，避免触发父级 el-card 的点击事件
  if (event) {
    event.stopPropagation()
  }
  
  await EbookShelfAdd([enid]).then((info) => {
    console.log(info)
    ElMessage({
      message: '已加入书架',
      type: 'success'
    })
    // 清空现有数据并重新加载
    product_list.product_list = []
    page.value = 0
    param.page = 0
    getAlgoProduct(param)
  }).catch((error) => {
    ElMessage({
      message: error,
      type: 'warning'
    })
  })
}

const ebookShelfRemove = async (enid: string, event?: Event) => {
  // 阻止事件冒泡，避免触发父级 el-card 的点击事件
  if (event) {
    event.stopPropagation()
  }
  
  ElMessageBox.confirm(
    '是否移出书架?',
    '',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(() => {
      EbookShelfRemove([enid]).then((info) => {
        console.log(info)
        ElMessage({
          message: '已移出书架',
          type: 'success'
        })
        // 清空现有数据并重新加载
        product_list.product_list = []
        page.value = 0
        param.page = 0
        getAlgoProduct(param)
      }).catch((error) => {
        ElMessage({
          message: error,
          type: 'warning'
        })
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消',
      })
    })
}

// 刷新 Algo 页面数据的函数
const refreshAlgoData = () => {
  // 清空现有数据并重新加载
  product_list.product_list = []
  page.value = 0
  param.page = 0
  getAlgoProduct(param)
}

</script>

<style scoped>
.category {
  padding: 24px;
  background: var(--fill-color-light);
  min-height: calc(100vh - 60px);
}

.filters {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 24px;
  box-shadow: var(--shadow-soft);
}

.filter-container {
  .filter-row {
    margin-bottom: 16px;
    
    &:last-child {
      margin-bottom: 0;
    }
  }

  .filter-label {
    color: var(--text-primary);
    font-size: 16px;
    font-weight: 500;
    line-height: 32px;
  }

  .filter-options {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .filter-btn {
    font-size: 14px;
    padding: 6px 16px;
    border-radius: 16px;
    transition: all 0.3s ease;
    color: var(--text-secondary);

    &:hover {
      color: var(--accent-color);
      background: rgba(255, 107, 0, 0.05);
    }
  }

  .active-btn {
    color: var(--accent-color);
    font-weight: 500;
    background: rgba(255, 107, 0, 0.1);
  }
}

.sub-options {
  background: var(--fill-color);
  padding: 12px;
  border-radius: 8px;
}

.filter-divider {
  margin: 16px 0;
  border-color: var(--border-soft);
}

.content-container {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 24px;
  box-shadow: var(--shadow-soft);
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  .result-count {
    font-size: 16px;
    color: var(--text-primary);

    .highlight {
      color: var(--accent-color);
      font-weight: 500;
    }
  }

  .sort-options {
    display: flex;
    gap: 16px;

    .sort-btn {
      color: var(--text-secondary);
      font-size: 14px;
      transition: all 0.3s ease;

      &:hover {
        color: var(--accent-color);
      }

      &.active-sort {
        color: var(--accent-color);
        font-weight: 500;
      }
    }
  }
}

.content-divider {
  margin: 16px 0;
  border-color: var(--border-soft);
}

.content-list {
  .infinite-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
    gap: 24px;
    padding: 0;
    margin: 0;
    list-style: none;
  }
}

.content-item {
  .content-card {
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-4px);
    }
  }
}

.card-content {
  display: flex;
  gap: 16px;
}

.content-cover {
  width: 112px;
  height: 150px;
  object-fit: cover;
  border-radius: 8px;
  background: var(--fill-color-light);
}

.content-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  text-align: left;
}

.content-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.4;
  text-align: left;
}

.content-intro {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-align: left;
}

.meta-info {
  font-size: 13px;
  color: var(--text-tertiary);
  margin-top: auto;
  text-align: left;

  .lecturer, .author {
    display: block;
    margin-bottom: 4px;
    text-align: left;
  }

  .episode-count {
    color: var(--text-secondary);
    text-align: left;
  }
}

.rating-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  text-align: left;

  .rating {
    text-align: left;
    .no-rating {
      font-size: 13px;
      color: var(--text-tertiary);
    }
  }

  .learner-count {
    font-size: 12px;
    color: var(--text-tertiary);
    text-align: right;
  }
}

:deep(.el-card) {
  border-radius: 12px;
  border: none;
  overflow: hidden;

  .el-card__body {
    padding: 16px;
    text-align: left;
  }
}

:deep(.el-rate) {
  display: inline-flex;
  align-items: center;
  
  .el-rate__text {
    color: var(--accent-color);
  }
}

.card-actions .el-button {
  border-radius: 8px;
  font-weight: bold;
  border-color: var(--accent-color);
  background-color: var(--accent-color);
  color: #fff;
  transition: all 0.3s ease;
}

.card-actions .el-button:hover {
  background-color: var(--accent-hover);
  border-color: var(--accent-hover);
  transform: translateY(-2px);
}

.card-actions .el-button:not(.is-disabled) {
  border-color: var(--accent-color);
  background-color: var(--accent-color);
}

/* 暗色主题适配 */
.theme-dark .result-count {
  color: var(--text-primary) !important;
}

.theme-dark .result-count .highlight {
  color: var(--accent-color) !important;
}

.theme-dark .sort-btn {
  color: var(--text-secondary) !important;
}

.theme-dark .sort-btn:hover {
  color: var(--accent-color) !important;
}

.theme-dark .sort-btn.active-sort {
  color: var(--accent-color) !important;
  font-weight: 500;
}

.theme-dark .content-divider {
  border-color: var(--border-soft) !important;
}

.theme-dark .content-cover {
  background: var(--fill-color-light) !important;
}

.theme-dark .content-title {
  color: var(--text-primary) !important;
}

.theme-dark .content-intro {
  color: var(--text-secondary) !important;
}

.theme-dark .meta-info {
  color: var(--text-tertiary) !important;
}

.theme-dark .meta-info .episode-count {
  color: var(--text-secondary) !important;
}

.theme-dark .rating .no-rating {
  color: var(--text-tertiary) !important;
}

.theme-dark .learner-count {
  color: var(--text-tertiary) !important;
}

.theme-dark :deep(.el-rate__text) {
  color: var(--accent-color) !important;
}

.theme-dark .card-actions .el-button {
  border-color: var(--accent-color) !important;
  background-color: var(--accent-color) !important;
  color: #fff !important;
}

.theme-dark .card-actions .el-button:hover {
  background-color: var(--accent-hover) !important;
  border-color: var(--accent-hover) !important;
}

/* 确保选中按钮在暗色主题下可见 */
.theme-dark .filter-btn.active-btn {
  color: var(--accent-color) !important;
  background: rgba(255, 107, 0, 0.15) !important;
  border: 1px solid var(--accent-color) !important;
}

.theme-dark .filter-btn {
  color: var(--text-secondary) !important;
  background: transparent !important;
  border: 1px solid transparent !important;
}

.theme-dark .filter-btn:hover {
  color: var(--accent-color) !important;
  background: rgba(255, 107, 0, 0.05) !important;
  border: 1px solid var(--accent-color) !important;
}
</style>
