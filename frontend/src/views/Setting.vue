<template>
  <div class="setting-container">
    <div class="setting-header">
      <h1>系统设置</h1>
      <p class="subtitle">管理您的下载路径、工具配置和界面偏好</p>
    </div>

    <div class="setting-content">
      <!-- 基本设置卡片 -->
      <el-card class="setting-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon><Folder /></el-icon>
            <span>基本设置</span>
          </div>
        </template>
        
        <el-form :model="form" label-position="top">
          <el-form-item label="文件保存目录">
            <div class="dir-input-wrapper">
              <el-input 
                v-model="form.downloadDir"
                :prefix-icon="Folder"
                placeholder="选择文件保存位置"
                clearable 
                @clear="clearDir"
              />
              <el-button type="primary" @click="openDialogDir('选择文件保存目录')">
                <el-icon><FolderOpened /></el-icon>
                浏览
              </el-button>
            </div>
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              <span>设置下载文件的默认保存位置</span>
            </div>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 下载工具设置卡片 -->
      <el-card class="setting-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon><Download /></el-icon>
            <span>下载工具设置</span>
          </div>
        </template>
        
        <el-form :model="form" label-position="top">
          <el-form-item>
            <div class="info-card">
              <div class="info-header">
                <el-icon><InfoFilled /></el-icon>
                <span>下载功能说明</span>
              </div>
              <ul class="info-list">
                <li>
                  <div class="info-item">
                    <el-icon><CircleCheck /></el-icon>
                    <span>只有已购买的课程、加入书架的听书和电子书支持下载</span>
                  </div>
                </li>
                <li>
                  <div class="info-item">
                    <el-icon><CircleCheck /></el-icon>
                    <span>每天听本书和电子书均不支持批量下载</span>
                  </div>
                </li>
                <li>
                  <div class="info-item">
                    <el-icon><Connection /></el-icon>
                    <span>如需下载音频，需安装 <el-link href="https://ffmpeg.org" target="_blank" type="primary">ffmpeg</el-link> 并设置路径</span>
                  </div>
                </li>
                <li>
                  <div class="info-item">
                    <el-icon><Document /></el-icon>
                    <span>如需下载PDF格式电子书，需安装 <el-link href="https://wkhtmltopdf.org/downloads.html" target="_blank" type="primary">wkhtmltopdf</el-link> 并设置路径</span>
                  </div>
                </li>
              </ul>
            </div>
          </el-form-item>

          <el-form-item label="ffmpeg 可执行文件路径">
            <div class="dir-input-wrapper">
              <el-input 
                v-model="form.ffmpegDir"
                :prefix-icon="Connection"
                placeholder="例如: /usr/local/bin/ffmpeg"
                clearable
              />
              <el-button @click="openToolPathDialog('选择ffmpeg路径', 'ffmpeg')">
                <el-icon><FolderOpened /></el-icon>
                浏览
              </el-button>
            </div>
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              <span>音频需要借助 <el-link href="https://ffmpeg.org" target="_blank" type="primary">ffmpeg</el-link> 合成，安装后终端执行 whereis ffmpeg 查找路径</span>
            </div>
          </el-form-item>

          <el-form-item label="wkhtmltopdf 可执行文件路径">
            <div class="dir-input-wrapper">
              <el-input 
                v-model="form.wkhtmltopdfDir"
                :prefix-icon="Document"
                placeholder="例如: /usr/local/bin/wkhtmltopdf"
                clearable
              />
              <el-button @click="openToolPathDialog('选择wkhtmltopdf路径', 'wkhtmltopdf')">
                <el-icon><FolderOpened /></el-icon>
                浏览
              </el-button>
            </div>
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              <span>电子书转 PDF 需要借助 <el-link href="https://wkhtmltopdf.org/downloads.html" target="_blank" type="primary">wkhtmltopdf</el-link>，安装后终端执行 whereis wkhtmltopdf 查找路径</span>
            </div>
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              <span>macOS 里可用 ⌘+⇧+. 显示隐藏的路径，也可在输入框手动粘贴完整路径</span>
            </div>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 界面设置卡片 -->
      <el-card class="setting-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon><Brush /></el-icon>
            <span>界面设置</span>
          </div>
        </template>
        
        <el-form :model="form" label-position="top">
          <el-form-item label="主题色">
            <div class="theme-color-container">
              <el-color-picker 
                v-model="form.systemColor"
                :predefine="predefineColors"
                show-alpha
                size="large"
                @change="setThemeColor(form.systemColor)"
              />
              <div class="theme-preview">
                <div class="color-preview" :style="{ backgroundColor: form.systemColor }"></div>
                <span class="color-value">{{ form.systemColor }}</span>
              </div>
            </div>
          </el-form-item>

          <el-form-item label="字体设置">
            <el-radio-group v-model="form.fontFamily" @change="setFontFamily(form.fontFamily)">
              <el-radio label="default" border>默认 (PingFang SC)</el-radio>
              <el-radio label="jetbrains" border>JetBrains Mono</el-radio>
              <el-radio label="wenkai" border>霞鹜文楷 (LXGW WenKai)</el-radio>
            </el-radio-group>
            <div class="form-tip">
              <el-icon><InfoFilled /></el-icon>
              <span>选择您喜欢的界面字体风格</span>
            </div>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button @click="resetForm" :icon="RefreshLeft" size="large">重置</el-button>
        <el-button type="primary" @click="onSubmit" :icon="Check" size="large">保存设置</el-button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'
import { 
  Folder, 
  FolderOpened, 
  Download, 
  Connection, 
  Document, 
  Brush,
  Check,
  RefreshLeft,
  InfoFilled,
  CircleCheck
} from '@element-plus/icons-vue'
import { settingStore } from "../stores/setting"
import { themeStore } from "../stores/theme"
import { OpenDirectoryDialog, OpenFileDialog } from "../../wailsjs/go/backend/App"
import { setThemeColor, setFontFamily } from "../utils/utils"
import { ElMessage } from 'element-plus'

const store = settingStore()
const themeStoreInstance = themeStore()

const predefineColors = ref([
  '#ff6b00',
  '#ff4500',
  '#ff8c00',
  '#ffd700',
  '#90ee90',
  '#00ced1',
  '#1e90ff',
  '#409EFF',
  '#c71585',
  'rgba(255, 69, 0, 0.68)',
  'hsv(51, 100, 98)',
  'hsva(120, 40, 94, 0.5)',
  'hsl(181, 100%, 37%)',
  'hsla(209, 100%, 56%, 0.73)',
  '#c7158577',
])

const form = reactive({
  downloadDir: store.getDownloadDir,
  systemColor: themeStoreInstance.color || store.getColor || '#ff6b00',
  ffmpegDir: store.getFfmpegDirDir,
  wkhtmltopdfDir: store.getWkDir,
  fontFamily: store.setting.fontFamily || 'default',
})

// 保存原始设置值用于重置
const originalSettings = {
  downloadDir: store.getDownloadDir,
  systemColor: themeStoreInstance.color || store.getColor || '#ff6b00',
  ffmpegDir: store.getFfmpegDirDir,
  wkhtmltopdfDir: store.getWkDir,
  fontFamily: store.setting.fontFamily || 'default',
}

const openDialogDir = async (title: string) => {
  try {
    const result = await OpenDirectoryDialog(title)
    if (result) {
      form.downloadDir = result
      console.log('选择的目录:', result)
    }
  } catch (error) {
    console.error('打开目录选择器失败:', error)
  }
}

const openToolPathDialog = async (title: string, toolType: string) => {
  try {
    const result = await OpenFileDialog(title)
    if (result) {
      if (toolType === 'ffmpeg') {
        form.ffmpegDir = result
      } else if (toolType === 'wkhtmltopdf') {
        form.wkhtmltopdfDir = result
      }
      console.log(`选择的${toolType}路径:`, result)
    }
  } catch (error) {
    console.error(`打开${toolType}路径选择器失败:`, error)
  }
}

const clearDir = () => {
  form.downloadDir = ''
}

const resetForm = () => {
  form.downloadDir = originalSettings.downloadDir
  form.systemColor = originalSettings.systemColor
  form.ffmpegDir = originalSettings.ffmpegDir
  form.wkhtmltopdfDir = originalSettings.wkhtmltopdfDir
  
  // 重新应用主题色
  setThemeColor(form.systemColor)
  // 重新应用字体
  setFontFamily(form.fontFamily)
  
  ElMessage({
    message: '已重置为上次保存的设置',
    type: 'info',
  })
}

const onSubmit = () => {
  console.log('保存设置:', form)
  
  // 更新设置存储
  store.setting.downloadDir = form.downloadDir
  store.setting.theme = form.systemColor
  store.setting.ffmpegDir = form.ffmpegDir
  store.setting.wkhtmltopdfDir = form.wkhtmltopdfDir
  store.setting.fontFamily = form.fontFamily
  
  // 更新主题存储
  themeStoreInstance.setThemeColor(form.systemColor)
  setFontFamily(form.fontFamily)
  
  // 更新原始设置，用于重置
  Object.assign(originalSettings, {
    downloadDir: form.downloadDir,
    systemColor: form.systemColor,
    ffmpegDir: form.ffmpegDir,
    wkhtmltopdfDir: form.wkhtmltopdfDir,
    fontFamily: form.fontFamily,
  })
  
  ElMessage({
    message: '设置已保存',
    type: 'success',
  })
}
</script>

<style scoped>
.setting-container {
  padding: 32px;
  max-width: 900px;
  margin: 0 auto;
}

.setting-header {
  margin-bottom: 32px;
  text-align: left;
}

.setting-header h1 {
  font-size: 28px;
  margin: 0 0 8px 0;
  font-weight: 600;
  color: var(--text-primary);
}

.subtitle {
  color: var(--text-secondary);
  font-size: 14px;
  margin: 0;
}

.setting-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.setting-card {
  border-radius: 16px;
  border: 1px solid var(--border-soft);
  background-color: var(--card-bg);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.setting-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-medium);
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.card-header .el-icon {
  margin-right: 8px;
  font-size: 18px;
}

.dir-input-wrapper {
  display: flex;
  gap: 12px;
}

.dir-input-wrapper .el-input {
  flex: 1 1 auto;
  min-width: 420px;
}

.theme-color-container {
  display: flex;
  align-items: center;
  gap: 20px;
}

.theme-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background-color: var(--fill-color);
  border-radius: 8px;
}

.color-preview {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.color-value {
  font-family: monospace;
  color: var(--text-secondary);
  font-size: 14px;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 16px;
  padding-bottom: 32px;
}

.form-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-tertiary);
  font-size: 13px;
  margin-top: 8px;
}

.info-card {
  background-color: var(--fill-color);
  border-radius: 8px;
  padding: 16px;
}

.info-header {
  display: flex;
  align-items: center;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.info-header .el-icon {
  margin-right: 8px;
  color: var(--primary-color);
}

.info-list {
  padding: 0;
  margin: 0;
  list-style-type: none;
}

.info-list li {
  margin-bottom: 12px;
}

.info-list li:last-child {
  margin-bottom: 0;
}

.info-item {
  display: flex;
  align-items: flex-start;
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.5;
}

.info-item .el-icon {
  margin-right: 10px;
  margin-top: 3px;
  color: var(--success-color);
  flex-shrink: 0;
}

.info-item .el-icon.connection-icon,
.info-item .el-icon.document-icon {
  color: var(--primary-color);
}

@media (max-width: 600px) {
  .setting-container {
    padding: 20px;
  }
  
  .dir-input-wrapper {
    flex-direction: column;
  }
  
  .action-buttons {
    flex-direction: column-reverse;
  }
  
  .action-buttons .el-button {
    width: 100%;
  }
}
</style>
