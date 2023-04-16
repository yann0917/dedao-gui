<template>
  <el-dialog v-model="dialogVisible" width="50%" :before-close="closeDialog">
    <el-space>
      <el-card v-for="i in 1" :key="i" class="box-card" style="width: 340px">
        <template #header>
          <div class="card-header">
            <span>{{ data.title }}</span>
            <el-alert
              :title="data.tips"
              type="success"
              center
              :closable="false"
            />
          </div>
        </template>
        <div id="qrcode-code" class="qrcode-code">
          <el-image
            style="width: 124px; height: 124px"
            :src="data.qrCode"
            fit="fill"
          />
        </div>
        <div id="bottom">
          <el-image
            src="https://piccdn2.umiwi.com/fe-oss/default/MTYzNzMwNzUyMzQy.png"
            class="qrcode-login"
          />
        </div>
      </el-card>
    </el-space>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, onBeforeMount, onMounted, onBeforeUnmount } from "vue";
import { ElTable, ElMessage } from "element-plus";
import { GetQrcode, CheckLogin, UserInfo } from "../../wailsjs/go/backend/App";
import { useRouter } from "vue-router";
import { userStore } from "../stores/user";
import { services } from "../../wailsjs/go/models";
import { Local } from "../utils/storage";

// const route = useRoute()
const store = userStore();
const router = useRouter();
const data = reactive({
  qrCode: "",
  qrCodeString: "",
  token: "",
  title: "扫码登录",
  tips: "同时支持「得到App」和「微信」扫码👇",
});

const timeState = reactive({
  time: 600, // 600s倒计时
  timer: 0, // 定时器对象
});

const dialogVisible = ref(false);
const emits = defineEmits(["close"]);
const props = defineProps({
  dialogVisible: {
    type: Boolean,
    default: false,
  },
});

onMounted(() => {
  getQrcode();
  openDialog();
});

const getQrcode = () => {
  GetQrcode()
    .then((result) => {
      data.qrCode = result.qrCode;
      data.token = result.token;
      data.qrCodeString = result.qrCodeString;
      console.log(result);
    })
    .catch((error) => {
      ElMessage({
        message: error,
        type: "warning",
      });
    });
};

timeState.timer = setInterval(() => {
  timeState.time--;
  if (!timeState.time) {
    timeState.time = 600;
    clearInterval(timeState.timer);
    timeState.timer = 0;
  }

  CheckLogin(data.token, data.qrCodeString).then((loginResult) => {
    if (loginResult.status == 1 || loginResult.status == 2) {
      clearInterval(timeState.timer);
      timeState.timer = 0;
      if (loginResult.status == 1) {
        let user = reactive(new services.User());
        Object.assign(user, loginResult.user);
        store.user = user;

        Local.set("cookies", loginResult.cookie);

        if (store.userList.length == 0) {
          store.userList.push(user);
        } else {
          store.userList.forEach((item) => {
            if (item.uid_hazy != user.uid_hazy) {
              store.userList.push(user);
            }
          });
        }

        console.log(store);
        router.push("/user/profile");
      } else if (loginResult.status == 2) {
        router.push("/user/login");
      } else {
        Local.remove("cookies");
        Local.remove("userStore");
      }
    }

    console.log(loginResult);
  }).catch((error)=>{
    ElMessage({
      message: error,
      type: 'warning'
    })
    clearInterval(timeState.timer);
    timeState.timer = 0;
  });
}, 2000);

onBeforeUnmount(() => {
  clearInterval(timeState.timer);
  timeState.timer = 0;
//   router.push("/user/profile");
});

const openDialog = () => {
  dialogVisible.value = props.dialogVisible;
};

const closeDialog = () => {
  clearInterval(timeState.timer);
  timeState.timer = 0;
  emits("close");
};
</script>

<style scoped>
.qrcode-login {
  width: 247px;
  height: 80px;
  line-height: 20px;
  margin: 1.5rem auto;
}
</style>
