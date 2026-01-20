import { createApp } from 'vue'
import { createPinia } from "pinia";
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import App from './App.vue'
import router from "./router";
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { themeStore } from './stores/theme'

import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css' //样式
import 'element-plus/dist/index.css'
import './assets/css/font.css' // 字体样式
import './assets/css/global.css' // 全局自定义样式

const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)

//创建v-highlight全局指令
app.directive('highlight', function (el) {
  let blocks = el.querySelectorAll('pre code');
  blocks.forEach((block: any) => {
    hljs.highlightBlock(block)
  })
})

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.config.errorHandler = (err, vm, info) => {
  console.log('❌ 提示:', err);
  if (err instanceof Error) {
    console.log('❌ 提示:', err.message);
    if (err.message == '401 Unauthorized') {
      router.push("/user/login")
    }
  }
}

// 生产环境下它会被忽略
app.config.warnHandler = function(msg, vm, trace) {
  // `trace` 是组件的继承关系追踪
  console.warn('warnHandler msg: %o, vm: %o, trace: %o', msg, vm, trace)
}

// 在应用挂载后初始化主题
app.mount('#app')

// 初始化主题 - 必须在应用挂载后
const theme = themeStore()
theme.initTheme()
