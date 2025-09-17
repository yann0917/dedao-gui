<script lang="ts" setup>
import { onMounted, computed } from 'vue'
import Menu from './components/Menu.vue'
import { ElConfigProvider } from 'element-plus'
import 'element-plus/es/components/message/style/css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import { themeStore } from './stores/theme'

// 初始化主题
const store = themeStore()
onMounted(() => {
  store.initTheme()
})

// Element Plus 主题配置
const elementTheme = computed(() => store.isDark ? 'dark' : 'light')
</script>

<template>
  <el-config-provider :locale="zhCn" :theme="elementTheme">
    <el-container>
      <el-header>
        <Menu />
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
      <!-- <el-footer>Footer</el-footer> -->
    </el-container>
  </el-config-provider>
</template>

<style lang="scss">
@import url("./assets/css/font.css");

body {

  // background-color: transparent;
  img {
    max-width: 100%;
  }
}


// #app {
//   // position: relative;
//   // width: 900px;
//   // height: 520px;
// }

.el-container {
  height: 100%;
  background-color: var(--bg-color);
  transition: background-color 0.3s ease;
}

.el-header {
  position: relative;
  flex-direction: row;
  flex-wrap: nowrap;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 0;
  background-color: var(--bg-color);
  border-bottom: 1px solid var(--border-color-lighter);
  transition: background-color 0.3s ease, border-color 0.3s ease;

  .el-menu {
    height: auto;
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;
    justify-content: flex-start;
    background-color: transparent;
    border-bottom-width: 0px;
    flex: 1;

    .el-menu-item:hover,
    .el-menu-item:active,
    .el-menu-item:focus{
      background:none;
    }

    ul {
      a,
      a:hover,
      a:active,
      a:visited,
      a:link,
      a:focus {
        display: inline-block;
        padding: 0 5px;
        margin-right: 8px;
        text-align: center;
        text-decoration: none;
        white-space: nowrap;
        text-decoration: none;
        color: var(--text-color);
        transition: color 0.3s ease;
      }
    }
  }

}

.el-main {
  overflow-y: scroll;
  color: var(--text-color-secondary);
  width: 100%;
  height: 100%;
  background-color: var(--bg-color);
  transition: background-color 0.3s ease, color 0.3s ease;

  .el-pagination {
    margin-top: 10px;
    margin-bottom: 10px;
  }

  // 『只要在el-table元素中定义了height属性，即可实现固定表头的表格』不生效解决办法。
  .el-table {
    .el-table__body-wrapper {
      height: calc(100% - 5px) !important; // 表格高度减去表头的高度
    }
  }

  .el-breadcrumb {
    margin-bottom: 15px;
  }


}
</style>