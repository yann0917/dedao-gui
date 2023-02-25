<template>
  <div>
    <el-space wrap>
    <el-row :gutter="20">
    <el-col :span="16" class="notes">
      <ul v-infinite-scroll="loadNotes" 
      class="infinite-note-list" 
      style="overflow-y: auto; max-height: 800px;"  
      :infinite-scroll-disabled="infLoadingNote"
      infinite-scroll-distance="10"	>
      <li class="infinite-note-list-item">
      <el-card v-for="item in noteList.notes" :key="item.f_part.note_id_hazy" class="box-card" shadow="hover" style="width: 800px;">
        <template #header>
          <div class="card-header">

            <el-row :gutter="20" align="top">
              <el-col :span="24" v-if="item.comb?.length">
                <svg t="1677244742093" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="3684" width="14" height="14"><path d="M922.026667 439.04l-267.221334-282.368v185.258667h-85.333333c-250.325333 0-455.253333 181.504-481.322667 413.44 115.157333-133.674667 291.925333-218.922667 480.682667-218.922667h85.077333l0.554667 185.173333 267.52-282.581333z m-438.528 189.44C268.8 662.570667 87.04 821.76 47.232 1024A529.664 529.664 0 0 1 0 805.461333C0 502.272 254.976 256.597333 569.472 256.597333V40.96a35.242667 35.242667 0 0 1 10.368-30.208 39.68 39.68 0 0 1 54.4 0l378.794667 400.298667a35.413333 35.413333 0 0 1 10.88 27.861333 35.498667 35.498667 0 0 1-10.88 27.904l-376.704 398.08a37.973333 37.973333 0 0 1-56.448 2.133333 35.114667 35.114667 0 0 1-10.410667-30.208l-0.64-215.082666c-26.88 0-53.546667 2.005333-79.701333 5.845333l-5.632 0.853333z" fill="#8a8a8a" p-id="3685"></path></svg>
                  
                <span style="font-size: small;" v-if="item.comb.length<=2">  
                    {{ ' '+item.comb.map((v,i)=>{
                  return v.name
                }).join("、")+" 转发了" }}</span>
                 <span style="font-size: small;" v-else="item.comb.length>2">  
                    {{ ' '+item.comb.map((v,i)=>{
                      if(i<2) {
                        return v.name
                      }
                }).join("、")+"和其他"+(item.comb.length-2)+"转发了" }}</span>
              </el-col>
            </el-row>

            <el-row :gutter="20" align="top">
              <el-col :span="4">
                
                <el-avatar :size="72" :src="item.f_part.avatar" fit="fill" />
                <el-icon v-if="item.f_part.is_v==2"><svg t="1677248742045" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="18550" width="14" height="14"><path d="M512 0C229.376 0 0 229.376 0 512s229.376 512 512 512 512-229.376 512-512S794.624 0 512 0z m55.296 832L194.56 311.296l174.592-0.512 198.656 313.856 141.824-313.856 112.64 0.512-254.976 520.704z" fill="#ff6b00" p-id="18551"></path></svg></el-icon>
                <el-icon v-if="item.f_part.is_v==4"><svg t="1677248508098" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6456" width="14" height="14"><path d="M511.953455 1002.146909c-142.987636 0-408.901818-218.763636-408.901818-425.634909L103.051636 164.421818l40.657455-0.674909c0.861091 0 91.624727-1.931636 185.274182-39.936 96.046545-39.028364 157.998545-83.828364 158.580364-84.247273l24.273455-17.687273 24.482909 17.687273c0.581818 0.442182 62.533818 45.218909 158.580364 84.247273 93.649455 38.004364 184.413091 39.936 185.367273 39.936l40.471273 0.674909 0.186182 412.090182C920.948364 783.36 655.034182 1002.146909 511.953455 1002.146909L511.953455 1002.146909zM185.623273 243.409455l0 333.079273c0 159.953455 231.633455 343.063273 326.330182 343.063273 94.72 0 326.330182-183.109818 326.330182-343.063273L838.283636 243.409455c-40.471273-4.375273-106.170182-15.429818-174.405818-43.124364-69.934545-28.439273-123.042909-59.345455-151.947636-77.754182-28.811636 18.408727-81.989818 49.314909-151.854545 77.754182C291.793455 228.002909 226.071273 239.034182 185.623273 243.409455L185.623273 243.409455zM490.077091 731.345455l-173.614545-147.898182 53.387636-62.813091 111.383273 94.813091 211.386182-243.525818 62.417455 54.155636L490.077091 731.345455 490.077091 731.345455zM490.077091 731.345455" fill="#ff6b00" p-id="6457"></path></svg></el-icon>
              </el-col>
              <el-col :span="4" style="text-align: left; line-height: 10px;">
                <h3>{{ item.f_part.nick_name }}</h3>
                <el-tag class="ml-2" type="info">{{ item.f_part.time_desc }}</el-tag>
              </el-col>
              <el-col :span="16">
                <el-tag class="ml-2" type="success" v-if="item.f_part.slogan" >{{ item.f_part.slogan }}</el-tag>
                <el-tag class="ml-2" type="warning" v-if="item.f_part.v_info" >{{ item.f_part.v_info }}</el-tag>
              </el-col>
            </el-row>
          </div>
        </template>
        <div class="card-content">
          <span style="white-space: pre-wrap;">{{item.f_part.note }}</span>
          <div class="box" v-if="item.f_part.images?.length>0">
            <div class="imageBox">
              <el-image v-for="i in item.f_part.images" :src="i" :preview-src-list="item.f_part.images.map((v,i)=>{
                return v
              })" :initial-index="0"  fit="cover" style="width: 32%;"/>
            </div>
          </div>

          <div v-if="item.notes_type == 7">
              给这本书评了{{ item.f_part.note_score  }}分
          </div>

          <div class="note-line" v-if="item.f_part.base_source.title !=''">
            <div v-if="item.notes_type == 4">
              {{item.f_part.note_line }}
            <el-divider />
            </div>
            <div class="base-source">
              <el-row :gutter="2" align="top">
              <el-col :span="2" style="padding-left: 1%;">
                <el-image :src="item.f_part.base_source.img" fit="cover" style="width:90%;"/>
              </el-col>
              <el-col :span="16" style="text-align: left;line-height: 2.2em;">
                <div style="font-weight: bold;">{{ item.f_part.base_source.title }}</div>
                <div style="font-size: small;">{{ item.f_part.base_source.sub_title }}</div>
              </el-col>
            </el-row>
            </div>
          </div>

          <div class="s-part" v-if="item.s_part">
            <el-divider />

            <el-row :gutter="5" align="top">
              <el-col :span="3">
                <el-avatar :size="48" :src="item.s_part.avatar" fit="fill" />
                <el-icon v-if="item.s_part.is_v==2"><svg t="1677248742045" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="18550" width="14" height="14"><path d="M512 0C229.376 0 0 229.376 0 512s229.376 512 512 512 512-229.376 512-512S794.624 0 512 0z m55.296 832L194.56 311.296l174.592-0.512 198.656 313.856 141.824-313.856 112.64 0.512-254.976 520.704z" fill="#ff6b00" p-id="18551"></path></svg></el-icon>
                <el-icon v-if="item.s_part.is_v==4"><svg t="1677248508098" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6456" width="14" height="14"><path d="M511.953455 1002.146909c-142.987636 0-408.901818-218.763636-408.901818-425.634909L103.051636 164.421818l40.657455-0.674909c0.861091 0 91.624727-1.931636 185.274182-39.936 96.046545-39.028364 157.998545-83.828364 158.580364-84.247273l24.273455-17.687273 24.482909 17.687273c0.581818 0.442182 62.533818 45.218909 158.580364 84.247273 93.649455 38.004364 184.413091 39.936 185.367273 39.936l40.471273 0.674909 0.186182 412.090182C920.948364 783.36 655.034182 1002.146909 511.953455 1002.146909L511.953455 1002.146909zM185.623273 243.409455l0 333.079273c0 159.953455 231.633455 343.063273 326.330182 343.063273 94.72 0 326.330182-183.109818 326.330182-343.063273L838.283636 243.409455c-40.471273-4.375273-106.170182-15.429818-174.405818-43.124364-69.934545-28.439273-123.042909-59.345455-151.947636-77.754182-28.811636 18.408727-81.989818 49.314909-151.854545 77.754182C291.793455 228.002909 226.071273 239.034182 185.623273 243.409455L185.623273 243.409455zM490.077091 731.345455l-173.614545-147.898182 53.387636-62.813091 111.383273 94.813091 211.386182-243.525818 62.417455 54.155636L490.077091 731.345455 490.077091 731.345455zM490.077091 731.345455" fill="#ff6b00" p-id="6457"></path></svg></el-icon>
              </el-col>
              <el-col :span="6">
                <h4>{{ item.s_part.nick_name }}</h4>
              </el-col>
            </el-row>
            <el-row :gutter="5" align="top">
              <el-col :span="24">
                <span style="white-space: pre-wrap;">{{item.s_part.note }}</span>
                <div class="note-line" v-if="item.s_part.base_source.title !=''">
                  <div v-if="item.notes_type == 4">
                    {{item.s_part.note_line }}
                  <el-divider />
                  </div>
                <div class="base-source">
                  <el-row :gutter="2" align="top">
                  <el-col :span="2" style="padding-left: 1%;">
                    <el-image :src="item.s_part.base_source.img" fit="cover" style="width:90%;"/>
                  </el-col>
                  <el-col :span="16" style="text-align: left;line-height: 2.2em;">
                    <div style="font-weight: bold;">{{ item.s_part.base_source.title }}</div>
                    <div style="font-size: small;">{{ item.s_part.base_source.sub_title }}</div>
                  </el-col>
                </el-row>
                </div>
                </div>
              </el-col>
             
            </el-row>
          </div>

          <div class="topic-name" style="padding-top: 1em;">
            <el-tag class="ml-2" type="info" v-if="item.topic?.topic_name" round>
              {{ '#'+item.topic?.topic_name }}
            </el-tag>
          </div>
          
          <div class="note-count" style="padding-top: 1em;">
            <el-row :gutter="20">
              <el-col :span="4">
                <el-icon>
                <svg t="1677237782377" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="4603" width="16" height="16"><path d="M922.026667 439.04l-267.221334-282.368v185.258667h-85.333333c-250.325333 0-455.253333 181.504-481.322667 413.44 115.157333-133.674667 291.925333-218.922667 480.682667-218.922667h85.077333l0.554667 185.173333 267.52-282.581333z m-438.528 189.44C268.8 662.570667 87.04 821.76 47.232 1024A529.664 529.664 0 0 1 0 805.461333C0 502.272 254.976 256.597333 569.472 256.597333V40.96a35.242667 35.242667 0 0 1 10.368-30.208 39.68 39.68 0 0 1 54.4 0l378.794667 400.298667a35.413333 35.413333 0 0 1 10.88 27.861333 35.498667 35.498667 0 0 1-10.88 27.904l-376.704 398.08a37.973333 37.973333 0 0 1-56.448 2.133333 35.114667 35.114667 0 0 1-10.410667-30.208l-0.64-215.082666c-26.88 0-53.546667 2.005333-79.701333 5.845333l-5.632 0.853333z" fill="#191919" p-id="4604"></path></svg>
            </el-icon>
            {{ item.note_count?.repost_count }}
              </el-col>
              <el-col :span="4">
                <el-icon>
              <svg t="1677237746696" class="icon" viewBox="0 0 1058 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="3564" width="16" height="16"><path d="M330.744242 885.372121l194.779798-129.861818 16.665859-11.106263h383.844848c36.486465 0 66.19798-29.659798 66.19798-66.146262v-529.19596c0-36.434747-29.711515-66.107475-66.19798-66.107475H132.305455c-36.486465 0-66.146263 29.659798-66.146263 66.107475v529.19596c0 36.486465 29.659798 66.146263 66.146263 66.146262h198.438787v140.968081m-66.146262 123.578182V810.550303H132.305455c-73.024646 0-132.305455-59.216162-132.305455-132.292525v-529.19596C0 76.024242 59.267879 16.808081 132.305455 16.808081h793.742222c73.076364 0 132.357172 59.216162 132.357171 132.240808v529.195959c0 73.076364-59.267879 132.292525-132.357171 132.292526h-363.830303L264.59798 1008.950303z m0 0" p-id="3565"></path></svg>
            </el-icon>
            {{ item.note_count?.comment_count }}

              </el-col>
              <el-col :span="4">
                <el-icon>
              <svg t="1677237515155" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="1474" width="16" height="16"><path d="M857.28 344.992h-264.832c12.576-44.256 18.944-83.584 18.944-118.208 0-78.56-71.808-153.792-140.544-143.808-60.608 8.8-89.536 59.904-89.536 125.536v59.296c0 76.064-58.208 140.928-132.224 148.064l-117.728-0.192A67.36 67.36 0 0 0 64 483.04V872c0 37.216 30.144 67.36 67.36 67.36h652.192a102.72 102.72 0 0 0 100.928-83.584l73.728-388.96a102.72 102.72 0 0 0-100.928-121.824zM128 872V483.04c0-1.856 1.504-3.36 3.36-3.36H208v395.68H131.36A3.36 3.36 0 0 1 128 872z m767.328-417.088l-73.728 388.96a38.72 38.72 0 0 1-38.048 31.488H272V476.864a213.312 213.312 0 0 0 173.312-209.088V208.512c0-37.568 12.064-58.912 34.72-62.176 27.04-3.936 67.36 38.336 67.36 80.48 0 37.312-9.504 84-28.864 139.712a32 32 0 0 0 30.24 42.496h308.512a38.72 38.72 0 0 1 38.048 45.888z" p-id="1475"></path></svg>
            </el-icon>
            {{ item.note_count?.like_count }}
              </el-col>
            </el-row>
          
          </div>
        </div>
      </el-card>
    </li>
  </ul>

    </el-col>
   
    <el-col :span="8" justify="center" class="topic"  v-if="topicList.list.length>0">
      <h2 style="text-align: left;"><span style="color:#ff6b00">#</span>推荐话题</h2>
      <ul v-infinite-scroll="loadTopic" 
      class="infinite-list" 
      style="overflow-y: auto; max-height: 800px;" 
      :infinite-scroll-disabled="infLoadingTopic"
      infinite-scroll-distance="10"
      >
        <li v-for="i in topicList.list"  class="infinite-list-item">
          <el-button text class="topic-item">
          <el-badge :value="i.tag==1?'新':''" class="badge-item">
            <div style="font-size: large; font-weight: bold; text-align: left;padding-bottom: 10px;">
              <span style="color:#ff6b00">#</span>{{i.name}} 
            </div>
            <div style="white-space: pre-wrap; font-weight:light;text-align: left;">{{i.intro}}</div>
          </el-badge>
          </el-button>
        </li>
      </ul>
    </el-col>
  </el-row>
    </el-space>
  </div>
       
  </template>
  
  <script lang="ts" setup>
  import { ref, reactive, onMounted, watch } from 'vue'
  import { ElMessage } from 'element-plus'
  import { NotesTimeline, TopicAll } from '../../wailsjs/go/backend/App'
  import { services } from '../../wailsjs/go/models'
  import { useRoute, useRouter } from 'vue-router'
  import { userStore } from '../stores/user';
import { it } from 'element-plus/es/locale'
  
  const store = userStore()
  const router = useRouter()
  const route = useRoute()
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(15)
  const infLoadingNote = ref(true)
  const infLoadingTopic = ref(true)

  let tableData = reactive(new services.NotesTimeline)
  let topicAll = reactive(new services.TopicAll)
  const topicList = reactive(new services.TopicAll)
  topicList.list = []

  const noteList = reactive(new services.NotesTimeline)
  noteList.notes = []

  let id = ref()
  let isMore= ref(false)
  let maxId= ref("")
  let enid = ref()
  let breadcrumbTitle = ref()
  
  const topicPage = ref(0)
  const topicHasMore = ref(false)
  const topicPageSize = ref(20)
 

  onMounted(() => {
    getTableData()
    getTopicAll()
  })
  
  const loadTopic = () => {
    infLoadingTopic.value=true
    if (topicHasMore.value === true) {
      topicPage.value += 1
      getTopicAll()
    }
  }
  
  const loadNotes = () => {
    infLoadingNote.value = true
    if (noteList.is_more === true) {
      getTableData()
    }
  }
  
  const getTableData = async () => {
    await NotesTimeline(maxId.value).then((table) => {
      // console.log(table)
      Object.assign(tableData, table)
      isMore.value = tableData.is_more
      maxId.value = tableData.max_id
      tableData.notes.forEach((item,index,array)=>{
        if(item.f_part.images?.length>0) {
          item.f_part.images.forEach((item1, index1, array1) => {
            let imgs = JSON.parse(item1)
            array1[index1]=imgs.url
        })
        }
      })

      noteList.is_more= tableData.is_more
      noteList.max_id = tableData.max_id
      noteList.notes.push(...tableData.notes)
     
    }).catch((error) => {
      ElMessage({
        message: error,
        type: 'warning'
      })
    })
    infLoadingNote.value = false
  }
  
  const getTopicAll =async () => {
    await TopicAll(topicPage.value, topicPageSize.value).then((result)=>{
      Object.assign(topicAll, result)
      // console.log(topicAll)
      topicHasMore.value= topicAll.has_more
      topicList.has_more = topicAll.has_more
      topicList.list.push(...topicAll.list)
    }).catch((error)=>{
      ElMessage({
        message: error,
        type: 'warning'
      })
    })
    infLoadingTopic.value = false
  }

  </script>
  

  <style scoped lang="scss">
  .el-card{
    color:#606266;
    .card-header,
    .card-content {
      text-align: left ;
  }
  }
  .topic-item {
    height: 120px;
}
  .el-tag {
    margin-right: 5px;
      // height: auto;
    text-align: center;
  }

  .el-row {
  margin-bottom: 5px;
}
.el-row:last-child {
  margin-bottom: 0;
}

.topic {
  height: 800px;
  border: 1px solid var(--el-border-color);
  border-radius:15px;
}

ul {
  padding: 0;
  margin: 0;
  text-align: left;
  list-style: none;
}

  .infinite-list {
  height: 90%;
  padding: 0;
  margin: 0;
  text-align: left;
  list-style: none;

}
.infinite-list .infinite-list-item {
  margin: 1%;
  // background: var(--el-color-primary-light-9);
  // color: var(--el-color-primary);
  color: #606266;
  padding-bottom: 1%;
  border-bottom: 1px dashed;

}
.infinite-list .infinite-list-item:last-child {
  margin-bottom: 0;
}

.infinite-list .infinite-list-item + .list-item {
  margin-top: 10px;
}

.note-line{
  border-radius: 0.4em;
  white-space: pre-wrap; 
  background:rgb(242, 241, 241); 
  opacity: 1; 
  padding-top: 1em;
  text-align: left;
}
.box {
    display: flex;
    flex-wrap: wrap;
}
.imageBox .el-image {
    position: relative;
    overflow: hidden;
    margin: 0.5%;
    width: 245px;
}
.imageBox .el-image img:nth-child(2):nth-last-child(1),
.imageBox .el-image img:nth-child(1):nth-last-child(2) {
    width: 49%;
}
    /* 2/3 */
.imageBox .el-image img:nth-child(1):nth-last-child(3),
.imageBox .el-image img:nth-child(2):nth-last-child(2),
.imageBox .el-image img:nth-child(3):nth-last-child(1) {
    width: 32%;
}

/* 4 or 1/2 */
.imageBox .el-image img:nth-child(1):nth-last-child(4),
.imageBox .el-image img:nth-child(2):nth-last-child(3),
.imageBox .el-image img:nth-child(3):nth-last-child(2),
.imageBox .el-image img:nth-child(4):nth-last-child(1) {
    width: 49%;
}

/*  5张以上图片  */
.imageBox .el-image img:nth-child(1):nth-last-child(n + 5),
.imageBox .el-image img:nth-child(1):nth-last-child(n + 5)~img {
    width: 32%;
}

</style>