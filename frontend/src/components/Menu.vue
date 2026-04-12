<template>
    <div class="menu-container" style="--wails-draggable:drag" @dblclick="handleTitlebarDblClick">
        <!-- 窗口控制按钮 -->
        <div class="window-controls" :class="isWindows ? 'windows-style' : 'mac-style'">
            <template v-if="isWindows">
                <!-- Windows 风格 -->
                <button class="window-btn windows-btn" @click="windowMinimize" title="最小化">
                    <svg viewBox="0 0 10 10" class="icon-minimize"><path d="M0 5h10" stroke="currentColor" stroke-width="1"/></svg>
                </button>
                <button class="window-btn windows-btn" @click="windowMaximize" :title="isMaximized ? '还原' : '最大化'">
                    <svg v-if="!isMaximized" viewBox="0 0 10 10" class="icon-maximize"><rect x="0.5" y="0.5" width="9" height="9" stroke="currentColor" fill="none" stroke-width="1"/></svg>
                    <svg v-else viewBox="0 0 10 10" class="icon-restore">
                        <path d="M3 1.5H8.5V7H7.5V2.5H3V1.5Z" fill="currentColor"/>
                        <path d="M1.5 3H7V8.5H1.5V3ZM2.5 4V7.5H6V4H2.5Z" fill="currentColor"/>
                    </svg>
                </button>
                <button class="window-btn windows-btn close-btn windows-close" @click="windowClose" title="关闭">
                    <svg viewBox="0 0 10 10" class="icon-close"><path d="M0 0L10 10M10 0L0 10" stroke="currentColor" stroke-width="1"/></svg>
                </button>
            </template>
            <template v-else>
                <!-- macOS 风格 -->
                <button class="window-btn mac-btn close-btn" @click="windowClose" title="关闭"></button>
                <button class="window-btn mac-btn minimize-btn" @click="windowMinimize" title="最小化"></button>
                <button class="window-btn mac-btn maximize-btn" @click="windowMaximize" title="最大化"></button>
            </template>
        </div>
        
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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Moon, Sunny } from '@element-plus/icons-vue'
import MenuItem from './MenuItem.vue'
import { themeStore } from '../stores/theme'
import { Environment, Quit, WindowIsMaximised, WindowMinimise, WindowToggleMaximise } from '../../wailsjs/runtime'

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

const isWindows = ref(false)
const isMaximized = ref(false)

const initPlatform = async () => {
  try {
    const env = await Environment()
    isWindows.value = env.platform === 'windows'
    if (isWindows.value) {
      await syncMaximizedState()
    }
  } catch {
    isWindows.value = false
  }
}

const syncMaximizedState = async () => {
  try {
    isMaximized.value = await WindowIsMaximised()
  } catch {
    isMaximized.value = false
  }
}

onMounted(() => {
  initPlatform()
})

const windowClose = () => {
  Quit()
}

const windowMinimize = () => {
  WindowMinimise()
}

const windowMaximize = async () => {
  WindowToggleMaximise()
  await syncMaximizedState()
}

const handleTitlebarDblClick = async (event: MouseEvent) => {
  if (!isWindows.value) {
    return
  }
  const target = event.target as HTMLElement | null
  const ignore = target?.closest('.window-controls, .theme-switch-container, .el-menu')
  if (ignore) {
    return
  }
  await windowMaximize()
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

/* 窗口控制按钮通用样式 */
.window-controls {
  display: flex;
  align-items: center;
  padding: 0 12px;
  margin-right: 16px;
  -webkit-app-region: no-drag;
}

.window-controls.mac-style {
  gap: 8px;
}

.window-controls.windows-style {
  gap: 4px;
}

.window-btn {
  border: none;
  cursor: pointer;
  transition: opacity 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* macOS 风格按钮 */
.mac-btn {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.mac-btn:hover {
  opacity: 0.8;
}

.mac-btn.close-btn {
  background-color: #ff5f57;
}

.mac-btn.minimize-btn {
  background-color: #febc2e;
}

.mac-btn.maximize-btn {
  background-color: #28c840;
}

/* Windows 风格按钮 */
.windows-btn {
  width: 28px;
  height: 28px;
  background: transparent;
  border-radius: 4px;
}

.windows-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.windows-btn .icon-minimize,
.windows-btn .icon-maximize,
.windows-btn .icon-close {
  width: 10px;
  height: 10px;
  color: var(--text-primary);
}

.windows-btn.windows-close:hover {
  background-color: #e81123;
}

.windows-btn.windows-close .icon-close {
  color: var(--text-primary);
}

.windows-btn.windows-close:hover .icon-close {
  color: white;
}
</style>
