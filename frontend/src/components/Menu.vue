<template>
    <div class="menu-container">
        <el-menu router :default-active="activeIndex" class="el-menu" mode="horizontal" 
             text-color="var(--text-primary)" 
             active-text-color="var(--accent-color)"
             @select="handleSelect" :collapse-transition="false">
            <menu-item v-for="(menu, key) in allRoutes" :key="key" :menu="menu" :path="menu.path" />
        </el-menu>
        
        <!-- 主题切换按钮 -->
        <div class="theme-switch-container">
            <el-switch 
                v-model="isDark" 
                :active-action-icon="Moon" 
                :inactive-action-icon="Sunny" 
                inline-prompt 
                @change="toggleDark"
                class="theme-switch"
            />
            <span style="margin-left: 8px; color: var(--text-primary); font-size: 12px;">
                {{ isDark ? '暗色' : '亮色' }}
            </span>
        </div>
    </div>
</template>
  
<script lang="ts" setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Moon, Sunny } from '@element-plus/icons-vue'
import MenuItem from './MenuItem.vue'
import { themeStore } from '../stores/theme'

const route = useRoute()
const router = useRouter()

// 主题相关
const store = themeStore()
const isDark = computed(() => store.isDark)

const toggleDark = () => {
  console.log('切换主题，当前状态:', store.isDark)
  store.toggleTheme()
  console.log('切换后状态:', store.isDark)
}

const allRoutes = router.options.routes.filter(route => 
  route.meta?.name !== '主题' // 过滤掉主题页面
)
const activeIndex = computed(() => {
    return route.path
})
const handleSelect = (key: string, keyPath: string[]) => {
    // console.log(key, keyPath)
}
</script>

<style scoped>
.menu-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.theme-switch-container {
  display: flex;
  align-items: center;
  padding: 0 16px;
  margin-left: auto;
  min-width: 120px;
  justify-content: flex-end;
}

.theme-switch {
  --el-switch-on-color: var(--accent-color);
  --el-switch-off-color: var(--border-soft);
  --el-switch-border-color: var(--border-soft);
}

.theme-switch :deep(.el-switch__core) {
  background-color: var(--card-bg);
  border-color: var(--border-soft);
  transition: all 0.3s ease;
}

.theme-switch :deep(.el-switch__core:hover) {
  border-color: var(--accent-color);
}

.theme-switch :deep(.el-switch__action) {
  background-color: var(--card-bg);
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.theme-switch :deep(.el-switch__action:hover) {
  color: var(--accent-color);
}

/* 暗色模式下的主题切换按钮 */
.theme-dark .theme-switch :deep(.el-switch__core) {
  background-color: var(--card-bg) !important;
  border-color: var(--border-soft) !important;
}

.theme-dark .theme-switch :deep(.el-switch__core:hover) {
  border-color: var(--accent-color) !important;
}

.theme-dark .theme-switch :deep(.el-switch__action) {
  background-color: var(--card-bg) !important;
  color: var(--text-primary) !important;
}

.theme-dark .theme-switch :deep(.el-switch__action:hover) {
  color: var(--accent-color) !important;
}

/* 菜单样式适配 */
.el-menu {
  background-color: transparent !important;
  border-bottom: none !important;
}

.el-menu :deep(.el-menu-item) {
  color: var(--text-primary) !important;
  background-color: transparent !important;
  border-bottom: 2px solid transparent !important;
  transition: all 0.3s ease !important;
}

.el-menu :deep(.el-menu-item:hover) {
  color: var(--accent-color) !important;
  background-color: var(--card-hover-bg) !important;
  border-bottom-color: var(--accent-color) !important;
}

.el-menu :deep(.el-menu-item.is-active) {
  color: var(--accent-color) !important;
  background-color: var(--card-hover-bg) !important;
  border-bottom-color: var(--accent-color) !important;
  font-weight: 500 !important;
}

.el-menu :deep(.el-sub-menu__title) {
  color: var(--text-primary) !important;
  background-color: transparent !important;
  border-bottom: 2px solid transparent !important;
  transition: all 0.3s ease !important;
}

.el-menu :deep(.el-sub-menu__title:hover) {
  color: var(--accent-color) !important;
  background-color: var(--card-hover-bg) !important;
  border-bottom-color: var(--accent-color) !important;
}

/* 暗色主题下的菜单样式 */
.theme-dark .el-menu {
  background-color: transparent !important;
  border-bottom: none !important;
}

.theme-dark .el-menu :deep(.el-menu-item) {
  color: var(--text-primary) !important;
  background-color: transparent !important;
  border-bottom: 2px solid transparent !important;
}

.theme-dark .el-menu :deep(.el-menu-item:hover) {
  color: var(--accent-color) !important;
  background-color: var(--card-hover-bg) !important;
  border-bottom-color: var(--accent-color) !important;
}

.theme-dark .el-menu :deep(.el-menu-item.is-active) {
  color: var(--accent-color) !important;
  background-color: var(--card-hover-bg) !important;
  border-bottom-color: var(--accent-color) !important;
  font-weight: 500 !important;
}

.theme-dark .el-menu :deep(.el-sub-menu__title) {
  color: var(--text-primary) !important;
  background-color: transparent !important;
  border-bottom: 2px solid transparent !important;
}

.theme-dark .el-menu :deep(.el-sub-menu__title:hover) {
  color: var(--accent-color) !important;
  background-color: var(--card-hover-bg) !important;
  border-bottom-color: var(--accent-color) !important;
}
</style>