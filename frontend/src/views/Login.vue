<script lang="ts" setup>
import { reactive, onBeforeMount, onMounted, onBeforeUnmount } from 'vue'
import { GetQrcode, CheckLogin, UserInfo } from '../../wailsjs/go/backend/App'
import { useRouter } from 'vue-router'
import { userStore } from '../stores/user';
import { services } from '../../wailsjs/go/models'
import { Local } from '../utils/storage';

// const route = useRoute()
const store = userStore()
const router = useRouter()
const data = reactive({
  qrCode: "",
  qrCodeString: "",
  token: "",
  title: "扫码登录",
  tips: "同时支持「得到App」和「微信」扫码👇",
})

const timeState = reactive({
  time: 600,  // 600s倒计时
  timer: 0,  // 定时器对象
})
onBeforeMount(() => {
  GetQrcode().then(result => {
    data.qrCode = result.qrCode
    data.token = result.token
    data.qrCodeString = result.qrCodeString
    console.log(result)
  })
})

onMounted(() => {
  console.log(data)
  timeState.timer = setInterval(() => {
    timeState.time--
    if (!timeState.time) {
      timeState.time = 600
      clearInterval(timeState.timer)
      timeState.timer = 0
    }

    CheckLogin(data.token, data.qrCodeString).then(loginResult => {

      if (loginResult.status == 1 || loginResult.status == 2) {
        clearInterval(timeState.timer)
        timeState.timer = 0;
        if (loginResult.status == 1) {
          UserInfo().then(result => {
            let user = reactive(new services.User)
            Object.assign(user, result)
            store.user = user
            
            Local.set("cookies",loginResult.cookie)

            if ( store.userList.length == 0) {
              store.userList.push(user)
            } else {
              store.userList.forEach(item=>{
                if (item.uid_hazy != user.uid_hazy) {
                  store.userList.push(user)
                }
              })
            }
          
            console.log(store)
            router.push("/user/profile")
          }).catch((error)=>{
            console.log(error)
            Local.remove("cookies")
            Local.remove("userStore")
          })
        }else if (loginResult.status == 2) {
          router.push("/user/login")
        } else {
          Local.remove("cookies")
          Local.remove("userStore")
        }
      }

      console.log(loginResult)
    })

  }, 2000)
})

onBeforeUnmount(() => {
  clearInterval(timeState.timer)
  timeState.timer = 0
  router.push("/user/profile")
})

</script>

<template>
  <el-space wrap>
    <el-card v-for="i in 1" :key="i" class="box-card" style="width: 340px">
      <template #header>
        <div class="card-header">
          <span>{{ data.title }}</span>
          <el-alert :title="data.tips" type="success" center :closable="false" />
        </div>
      </template>
      <div id="qrcode-code" class="qrcode-code">
        <el-image style="width: 124px; height: 124px" :src="data.qrCode" fit="fill" />
        <!-- <button class="btn" @click="getQrCode">登录</button> -->
      </div>
      <div id="bottom">
        <el-image src="https://piccdn2.umiwi.com/fe-oss/default/MTYzNzMwNzUyMzQy.png" class="qrcode-login-example" />
      </div>
    </el-card>
  </el-space>
</template>

<style scoped>
.qrcode-login-example {
  width: 247px;
  height: 80px;
  line-height: 20px;
  margin: 1.5rem auto;
}
</style>
