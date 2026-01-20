<template>
  <el-dialog
    v-model="dialogVisible"
    width="400px"
    :before-close="closeDialog"
    :show-close="true"
    destroy-on-close
    center
    class="login-dialog"
  >
    <div class="login-container">
      <div class="login-header">
        <h3 class="title">{{ data.title }}</h3>
        <p class="subtitle">{{ data.tips }}</p>
      </div>
      
      <div class="qr-container">
        <div class="qr-wrapper">
            <el-image
              v-if="data.qrCode"
              class="qr-code-img"
              :src="data.qrCode"
              fit="fill"
            />
            <div v-else class="qr-loading">
                <el-icon class="is-loading"><Loading /></el-icon>
            </div>
        </div>
      </div>

      <div class="login-footer">
        <el-image
          src="https://piccdn2.umiwi.com/fe-oss/default/MTYzNzMwNzUyMzQy.png"
          class="app-logo"
        />
        <span class="footer-text">得到 App 扫码登录</span>
      </div>
    </div>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onBeforeUnmount } from "vue";
import { ElMessage } from "element-plus";
import { Loading } from '@element-plus/icons-vue';
import { GetQrcode, CheckLogin } from "../../wailsjs/go/backend/App";
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
  tips: "使用得到App或微信扫码登录",
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

timeState.timer = window.setInterval(() => {
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
    // ElMessage({
    //   message: error,
    //   type: 'warning'
    // })
    // Silent fail on check login errors to avoid spam
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
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.login-header {
  text-align: center;
  margin-bottom: 24px;
}

.title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

.qr-container {
  background: var(--bg-color);
  padding: 20px;
  border-radius: 12px;
  box-shadow: var(--shadow-inner);
  margin-bottom: 24px;
}

.qr-wrapper {
  width: 180px;
  height: 180px;
  background: white;
  padding: 10px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.qr-code {
  width: 100%;
  height: 100%;
}

.qr-loading {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  font-size: 24px;
}

.login-footer {
  width: 100%;
  display: flex;
  justify-content: center;
}

.logo-image {
  height: 40px;
  opacity: 0.8;
}

:deep(.el-dialog) {
  border-radius: 16px;
  overflow: hidden;
  background: var(--card-bg);
  box-shadow: var(--shadow-heavy);
}

:deep(.el-dialog__header) {
  margin: 0;
  padding: 0;
}

:deep(.el-dialog__body) {
  padding: 30px;
}
</style>
