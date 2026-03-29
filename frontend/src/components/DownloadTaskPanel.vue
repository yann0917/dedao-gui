<template>
  <div class="download-task-panel">
    <el-tooltip content="下载任务" placement="left" :show-after="200">
      <el-badge
        :value="activeCount"
        :hidden="activeCount === 0"
        :max="99"
        class="task-badge"
        :style="badgeStyle"
      >
        <el-button
          circle
          class="task-button"
          :class="{ 'task-button-active': activeCount > 0 }"
          @mousedown="onDragStart"
          @click="handleButtonClick"
        >
          <el-icon class="task-icon"><Download /></el-icon>
        </el-button>
      </el-badge>
    </el-tooltip>

    <el-drawer v-model="visible" title="下载任务" direction="rtl" size="560px">
      <template #header>
        <div class="drawer-header">
          <span>下载任务</span>
          <el-space>
            <el-button size="small" @click="refresh">刷新</el-button>
            <el-button size="small" type="warning" plain @click="clearFinished">清理终态</el-button>
            <el-button size="small" type="danger" plain @click="clearAll">清空全部</el-button>
          </el-space>
        </div>
      </template>

      <el-table :data="tasks" stripe border height="calc(100vh - 160px)" empty-text="暂无下载任务">
        <el-table-column label="类型" width="80">
          <template #default="{ row }">
            <span>{{ bizTypeLabel(row.bizType) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="任务名" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">
            <span>{{ taskTitle(row) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="160">
          <template #default="{ row }">
            <el-progress :percentage="Number(row.progress || 0)" :stroke-width="10" :show-text="false" />
            <div class="progress-text">{{ Number(row.progress || 0) }}%</div>
          </template>
        </el-table-column>
        <el-table-column label="当前项" min-width="130" show-overflow-tooltip>
          <template #default="{ row }">
            <span>{{ row.currentName || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="错误" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <span>{{ row.errorMessage || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-space>
              <el-button
                v-if="canPause(row.status)"
                size="small"
                type="warning"
                link
                @click="pause(row.id)"
              >
                暂停
              </el-button>
              <el-button
                v-if="canResume(row.status)"
                size="small"
                type="primary"
                link
                @click="resume(row.id)"
              >
                继续
              </el-button>
              <el-button
                v-if="canCancel(row.status)"
                size="small"
                type="danger"
                link
                @click="cancel(row.id)"
              >
                取消
              </el-button>
              <el-button
                v-if="canRetry(row.status)"
                size="small"
                type="primary"
                link
                @click="retry(row.id)"
              >
                重试
              </el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import { CancelDownloadTask, ClearDownloadTasks, ListDownloadTasks, PauseDownloadTask, ResumeDownloadTask, RetryDownloadTask } from '../../wailsjs/go/backend/App'
import { downloadmgr } from '../../wailsjs/go/models'
import { EventsOff, EventsOn } from '../../wailsjs/runtime/runtime'

const visible = ref(false)
const tasks = ref<downloadmgr.DownloadTask[]>([])
let timer: number | undefined
const position = ref({ x: 0, y: 0 })
const dragState = ref({
  dragging: false,
  moved: false,
  offsetX: 0,
  offsetY: 0,
})

const activeCount = computed(() => {
  return tasks.value.filter((task) => ['pending', 'queued', 'running'].includes(String(task.status))).length
})

const refresh = async () => {
  const resp = await ListDownloadTasks({
    page: 1,
    pageSize: 200,
    status: [],
    bizType: [],
  })
  tasks.value = Array.isArray(resp?.list) ? resp.list : []
}

const openPanel = async () => {
  visible.value = true
  await refresh()
}

const badgeStyle = computed(() => {
  return {
    left: `${position.value.x}px`,
    top: `${position.value.y}px`,
  }
})

const handleButtonClick = async () => {
  if (dragState.value.moved) {
    dragState.value.moved = false
    return
  }
  await openPanel()
}

const clampPosition = (x: number, y: number) => {
  const maxX = Math.max(window.innerWidth - 64, 0)
  const maxY = Math.max(window.innerHeight - 64, 0)
  return {
    x: Math.min(Math.max(0, x), maxX),
    y: Math.min(Math.max(0, y), maxY),
  }
}

const onDragMove = (event: MouseEvent) => {
  if (!dragState.value.dragging) return
  const next = clampPosition(
    event.clientX - dragState.value.offsetX,
    event.clientY - dragState.value.offsetY,
  )
  if (Math.abs(next.x - position.value.x) > 2 || Math.abs(next.y - position.value.y) > 2) {
    dragState.value.moved = true
  }
  position.value = next
}

const onDragEnd = () => {
  dragState.value.dragging = false
  window.removeEventListener('mousemove', onDragMove)
  window.removeEventListener('mouseup', onDragEnd)
}

const onDragStart = (event: MouseEvent) => {
  dragState.value.dragging = true
  dragState.value.moved = false
  dragState.value.offsetX = event.clientX - position.value.x
  dragState.value.offsetY = event.clientY - position.value.y
  window.addEventListener('mousemove', onDragMove)
  window.addEventListener('mouseup', onDragEnd)
}

const initPosition = () => {
  const startX = Math.max(window.innerWidth - 72, 0)
  const startY = Math.max(window.innerHeight - 92, 0)
  position.value = { x: startX, y: startY }
}

const canCancel = (status: string) => ['pending', 'queued', 'running'].includes(String(status))
const canRetry = (status: string) => ['failed', 'canceled'].includes(String(status))
const canPause = (status: string) => ['pending', 'queued', 'running'].includes(String(status))
const canResume = (status: string) => ['paused'].includes(String(status))

const statusLabel = (status: string) => {
  const s = String(status)
  if (s === 'pending') return '待执行'
  if (s === 'queued') return '排队中'
  if (s === 'running') return '执行中'
  if (s === 'paused') return '已暂停'
  if (s === 'success') return '已完成'
  if (s === 'failed') return '失败'
  if (s === 'canceled') return '已取消'
  return s
}

const statusType = (status: string) => {
  const s = String(status)
  if (s === 'success') return 'success'
  if (s === 'failed' || s === 'canceled') return 'danger'
  if (s === 'paused') return 'warning'
  if (s === 'running') return 'primary'
  return 'info'
}

const bizTypeLabel = (bizType: string) => {
  const t = String(bizType)
  if (t === 'course') return '课程'
  if (t === 'ebook') return '电子书'
  if (t === 'odob') return '听书'
  return t
}

const taskTitle = (task: downloadmgr.DownloadTask) => {
  const title = String(task?.title ?? '').trim()
  if (title) return title
  return `${bizTypeLabel(task.bizType)} #${task.bizId}`
}

const updateTaskByEvent = (payload: any) => {
  const id = String(payload?.taskId ?? '')
  if (!id) return
  const idx = tasks.value.findIndex((task) => task.id === id)
  if (idx < 0) {
    refresh()
    return
  }
  const task = tasks.value[idx]
  const next = {
    ...task,
    status: payload?.status ?? task.status,
    progress: typeof payload?.progress === 'number' ? payload.progress : task.progress,
    current: typeof payload?.current === 'number' ? payload.current : task.current,
    total: typeof payload?.total === 'number' ? payload.total : task.total,
    currentName: payload?.currentName ?? task.currentName,
    errorCode: payload?.errorCode ?? task.errorCode,
    errorMessage: payload?.errorMessage ?? task.errorMessage,
  }
  tasks.value.splice(idx, 1, next as any)
}

const cancel = async (id: string) => {
  try {
    await CancelDownloadTask(id)
    ElMessage({ message: '任务已取消', type: 'success' })
    await refresh()
  } catch (error) {
    ElMessage({ message: String(error), type: 'warning' })
  }
}

const retry = async (id: string) => {
  try {
    await RetryDownloadTask(id)
    ElMessage({ message: '任务已重新入队', type: 'success' })
    await refresh()
  } catch (error) {
    ElMessage({ message: String(error), type: 'warning' })
  }
}

const pause = async (id: string) => {
  try {
    await PauseDownloadTask(id)
    ElMessage({ message: '任务已暂停', type: 'success' })
    await refresh()
  } catch (error) {
    ElMessage({ message: String(error), type: 'warning' })
  }
}

const resume = async (id: string) => {
  try {
    await ResumeDownloadTask(id)
    ElMessage({ message: '任务已继续', type: 'success' })
    await refresh()
  } catch (error) {
    ElMessage({ message: String(error), type: 'warning' })
  }
}

const clearFinished = async () => {
  try {
    await ElMessageBox.confirm('将清理已完成、失败、已取消任务记录，是否继续？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await ClearDownloadTasks(false)
    ElMessage({ message: '终态任务已清理', type: 'success' })
    await refresh()
  } catch (error) {
    if (error === 'cancel') return
    ElMessage({ message: String(error), type: 'warning' })
  }
}

const clearAll = async () => {
  try {
    await ElMessageBox.confirm('将清空所有任务（包含进行中的任务），是否继续？', '高风险操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error',
    })
    await ClearDownloadTasks(true)
    ElMessage({ message: '任务数据已清空', type: 'success' })
    await refresh()
  } catch (error) {
    if (error === 'cancel') return
    ElMessage({ message: String(error), type: 'warning' })
  }
}

const handleOpenFromEvent = () => {
  openPanel()
}

onMounted(async () => {
  initPosition()
  await refresh()
  EventsOn('download:task:update', updateTaskByEvent)
  timer = window.setInterval(() => {
    refresh()
  }, 3000)
  window.addEventListener('download-task:open', handleOpenFromEvent)
})

onUnmounted(() => {
  onDragEnd()
  EventsOff('download:task:update')
  if (timer) {
    window.clearInterval(timer)
  }
  window.removeEventListener('download-task:open', handleOpenFromEvent)
})
</script>

<style scoped>
.task-badge {
  position: fixed;
  z-index: 2100;
  cursor: grab;
  user-select: none;
}

.task-button {
  width: 52px;
  height: 52px;
  padding: 0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  color: #fff;
  background: linear-gradient(135deg, #6aa9ff 0%, #3b82f6 48%, #2563eb 100%);
  box-shadow: 0 10px 24px rgba(37, 99, 235, 0.32), 0 2px 8px rgba(37, 99, 235, 0.24);
  transition: transform 0.2s ease, box-shadow 0.2s ease, filter 0.2s ease;
}

.task-button:hover {
  transform: translateY(-1px) scale(1.04);
  filter: brightness(1.05);
  box-shadow: 0 12px 30px rgba(37, 99, 235, 0.36), 0 3px 10px rgba(37, 99, 235, 0.3);
}

.task-button:active {
  transform: translateY(0) scale(0.98);
}

.task-icon {
  font-size: 20px;
}

.task-button-active {
  animation: taskPulse 1.8s ease-in-out infinite;
}

@keyframes taskPulse {
  0% {
    box-shadow: 0 10px 24px rgba(37, 99, 235, 0.32), 0 2px 8px rgba(37, 99, 235, 0.24), 0 0 0 0 rgba(59, 130, 246, 0.4);
  }
  70% {
    box-shadow: 0 10px 24px rgba(37, 99, 235, 0.32), 0 2px 8px rgba(37, 99, 235, 0.24), 0 0 0 10px rgba(59, 130, 246, 0);
  }
  100% {
    box-shadow: 0 10px 24px rgba(37, 99, 235, 0.32), 0 2px 8px rgba(37, 99, 235, 0.24), 0 0 0 0 rgba(59, 130, 246, 0);
  }
}

.drawer-header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.progress-text {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
