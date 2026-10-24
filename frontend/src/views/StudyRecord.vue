<template>
    <div class="record-container">
        <div class="page-header">
            <h2>学习记录</h2>
            <p class="page-desc">当前登录账号的最近学习内容，按学习时间倒序。</p>
        </div>

        <div v-if="statsError" class="stats-fallback">统计加载失败：{{ statsError }}</div>
        <div v-else-if="stats && stats.total > 0" class="stats-panel">
            <div class="stats-head">
                <span class="stats-title">近 {{ stats.window_days }} 天统计</span>
                <span v-if="stats.truncated" class="stats-note">窗口内记录较多，仅统计最近 500 条</span>
            </div>

            <div class="stat-cards">
                <div class="stat-card">
                    <span class="stat-value">{{ stats.total }}</span>
                    <span class="stat-label">记录数</span>
                </div>
                <div class="stat-card">
                    <span class="stat-value">{{ stats.active_days }}<i class="stat-unit">/{{ stats.window_days }}</i></span>
                    <span class="stat-label">学习天数</span>
                </div>
                <div class="stat-card">
                    <span class="stat-value">{{ stats.finish_rate }}<i class="stat-unit">%</i></span>
                    <span class="stat-label">完成率</span>
                </div>
                <div class="stat-card">
                    <span class="stat-value">{{ stats.top_type || '—' }}</span>
                    <span class="stat-label">主力类型</span>
                </div>
            </div>

            <div class="stats-section">
                <p class="section-title">内容构成</p>
                <div class="compose-bar">
                    <div
                        v-for="t in stats.types"
                        :key="t.name"
                        class="compose-seg"
                        :style="{ flexGrow: t.count, background: typeColor(t.name) }"
                        :title="`${t.name} ${t.count} 条`"
                    />
                </div>
                <div class="legend">
                    <span v-for="t in stats.types" :key="t.name" class="legend-item">
                        <i class="legend-dot" :style="{ background: typeColor(t.name) }" />{{ t.name }} {{ t.count }}
                    </span>
                </div>
            </div>

            <div class="charts-grid">
                <div class="stats-section">
                    <p class="section-title">每日学习条数（近 {{ stats.daily.length }} 天）</p>
                    <div class="daily-chart">
                        <div v-for="day in stats.daily" :key="day.date" class="daily-col" :title="dailyTip(day)">
                            <span v-if="day.count > 0" class="daily-count">{{ day.count }}</span>
                            <div v-if="day.count > 0" class="daily-bar" :style="{ height: dailyHeight(day.count) }">
                                <div
                                    v-for="seg in day.types"
                                    :key="seg.name"
                                    class="daily-seg"
                                    :style="{ flexGrow: seg.count, background: typeColor(seg.name) }"
                                />
                            </div>
                            <div v-else class="daily-empty" />
                            <span class="daily-date">{{ day.date }}</span>
                        </div>
                    </div>
                </div>

                <div class="stats-section">
                    <div class="hours-head">
                        <p class="section-title">学习时段分布</p>
                        <span v-if="peakHour.count > 0" class="stats-note">
                            单小时高峰 {{ padHour(peakHour.hour) }}:00（{{ peakHour.count }} 条）
                        </span>
                    </div>
                    <div class="hours-chart">
                        <div v-for="b in hourBuckets" :key="b.label" class="hour-row" :title="bucketTip(b)">
                            <span class="hour-label">{{ b.label }}<i>{{ b.range }}</i></span>
                            <div class="hour-track">
                                <div class="hour-fill" :style="{ width: hourWidth(b.count) }">
                                    <div
                                        v-for="seg in b.types"
                                        :key="seg.name"
                                        class="hour-seg"
                                        :style="{ flexGrow: seg.count, background: typeColor(seg.name) }"
                                    />
                                </div>
                            </div>
                            <span class="hour-count">{{ b.count }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <p class="list-title">全部记录</p>

        <el-alert v-if="error" :title="error" type="error" :closable="false" class="list-error" />

        <div v-loading="loading" class="record-list">
            <div v-for="item in items" :key="`${item.product_type}-${item.product_id}-${item.timestamp}`" class="record-card">
                <div class="record-cover">
                    <el-image v-if="item.square_img || item.index_img" :src="item.square_img || item.index_img" fit="cover" loading="lazy">
                        <template #error>
                            <div class="cover-placeholder"><el-icon><Picture /></el-icon></div>
                        </template>
                    </el-image>
                    <div v-else class="cover-placeholder"><el-icon><Reading /></el-icon></div>
                </div>
                <div class="record-body">
                    <div class="record-title-row">
                        <el-tag size="small" type="info" effect="plain" class="record-type">{{ item.type_name || '内容' }}</el-tag>
                        <span class="record-title" :title="stripHtml(item.title)">{{ stripHtml(item.title) }}</span>
                    </div>
                    <div class="record-meta">
                        <span v-if="item.author" class="meta-author">{{ item.author }}</span>
                        <span v-if="progressPercent(item) !== null" class="meta-progress">
                            <span class="prog-track">
                                <span class="prog-fill" :style="{ width: clampPercent(progressPercent(item)!) + '%' }" />
                            </span>
                            {{ progressPercent(item) }}%
                        </span>
                        <span v-else class="meta-progress-text">{{ progressText(item) }}</span>
                        <span v-if="studyTimeOf(item)">{{ studyTimeOf(item) }}</span>
                    </div>
                </div>
            </div>
            <el-empty v-if="!loading && !error && items.length === 0" description="还没有学习记录" />
        </div>

        <div v-if="hasMore && !error" class="load-more">
            <el-button :loading="loadingMore" @click="load(true)">加载更多</el-button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'
import { Picture, Reading } from '@element-plus/icons-vue'
import { RecentList, RecentStats } from '@backend-bindings/app'
import * as services from '@backend-bindings/services/models';
import { userStore } from '../stores/user';
import { useAppRouter } from '../composables/useRouter';
import { isUnauthorized } from '../utils/errors';

const PAGE_SIZE = 20

const store = userStore()
const { pushLogin } = useAppRouter()

const stats = ref<services.RecentStats | null>(null)
const statsError = ref('')
const items = ref<services.RecentItem[]>([])
const hasMore = ref(false)
const loading = ref(true)
const loadingMore = ref(false)
const error = ref('')

// 上游多个文本字段（title/last_info/intro）内嵌 <hl>、<h1> 等高亮标签，展示前统一剥掉
const stripHtml = (raw: string) =>
    (raw || '').replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim()

// 记录时间为毫秒时间戳；last_info 为空时（听书条目常见）直接格式化时间戳
const studyTimeOf = (item: services.RecentItem) => {
    if (item.last_info) return stripHtml(item.last_info)
    if (!item.timestamp) return ''
    return new Date(item.timestamp).toLocaleString('zh-CN')
}

// 进度展示优先级：max_progress 百分比 > intro 文本 > 原始 progress 数值
// 特例：上游 recent 接口对课程类条目的 progress 恒为 0，真实学习进度百分比放在 max_progress
const progressPercent = (item: services.RecentItem): number | null => {
    const intro = item.progress_intro
    if (!intro) return 0
    if (item.type_name === '课程' && intro.progress === 0 && intro.max_progress > 0) {
        return intro.max_progress
    }
    if (intro.max_progress > 0) {
        return Math.round((intro.progress / intro.max_progress) * 100)
    }
    return null
}

const progressText = (item: services.RecentItem) => {
    const intro = item.progress_intro
    if (intro?.intro) return stripHtml(intro.intro)
    return String(intro?.progress ?? 0)
}

const clampPercent = (n: number) => Math.max(0, Math.min(n, 100))

// 内容类型固定配色（与全局主题色对应），图例与各图表堆叠分段共用；未知类型落中性灰
const TYPE_COLORS: Record<string, string> = {
    听书: 'var(--danger-color)',
    课程: 'var(--success-color)',
    电子书: 'var(--warning-color)',
    锦囊: 'var(--info-color)',
}
const typeColor = (name: string) => TYPE_COLORS[name] || 'var(--text-placeholder)'

const breakdownText = (types: services.RecentNameCount[] | undefined) =>
    (types || []).map((t) => `${t.name} ${t.count}`).join('、')

// 时段分桶横向条：与近14天柱状图共用类型配色
const HOUR_BUCKETS = [
    { label: '凌晨', range: '00–06', from: 0, to: 6 },
    { label: '清晨', range: '06–09', from: 6, to: 9 },
    { label: '上午', range: '09–12', from: 9, to: 12 },
    { label: '午间', range: '12–14', from: 12, to: 14 },
    { label: '下午', range: '14–18', from: 14, to: 18 },
    { label: '晚间', range: '18–21', from: 18, to: 21 },
    { label: '深夜', range: '21–24', from: 21, to: 24 },
]

const hourBuckets = computed(() => {
    const s = stats.value
    if (!s) return []
    const byHour = new Map(s.hours.map((h) => [h.hour, h]))
    return HOUR_BUCKETS.map((bucket) => {
        const typeMap = new Map<string, number>()
        let count = 0
        for (let hour = bucket.from; hour < bucket.to; hour++) {
            const hourStat = byHour.get(hour)
            if (!hourStat) continue
            count += hourStat.count
            for (const seg of hourStat.types) {
                typeMap.set(seg.name, (typeMap.get(seg.name) ?? 0) + seg.count)
            }
        }
        // 按 stats.types 的全局顺序还原，保证各图分段顺序一致
        const types = s.types
            .filter((t) => typeMap.has(t.name))
            .map((t) => ({ name: t.name, count: typeMap.get(t.name) ?? 0 }))
        return { ...bucket, count, types }
    })
})

const peakHour = computed(() =>
    (stats.value?.hours || []).reduce(
        (top, item) => (item.count > top.count ? item : top),
        { hour: 0, count: 0 } as services.RecentHourCount,
    ),
)

const dailyMax = computed(() => Math.max(...(stats.value?.daily || []).map((d) => d.count), 1))
const dailyHeight = (count: number) => `${Math.max((count / dailyMax.value) * 100, 6)}%`
const dailyTip = (day: services.RecentDayCount) =>
    day.count > 0 ? `${day.date} · ${day.count} 条：${breakdownText(day.types)}` : `${day.date} · 无记录`

const hourMax = computed(() => Math.max(...hourBuckets.value.map((b) => b.count), 1))
const hourWidth = (count: number) => `${(count / hourMax.value) * 100}%`
const bucketTip = (b: { label: string; range: string; count: number; types: services.RecentNameCount[] }) =>
    b.count > 0 ? `${b.label} ${b.range} · ${b.count} 条：${breakdownText(b.types)}` : `${b.label} ${b.range} · 无记录`

const padHour = (h: number) => String(h).padStart(2, '0')

// 游标翻页：下一页 maxID 传已加载最后一条的 timestamp（毫秒）
const load = async (append: boolean) => {
    if (append) {
        loadingMore.value = true
    } else {
        loading.value = true
    }
    try {
        const cursor = append && items.value.length > 0 ? items.value[items.value.length - 1].timestamp : 0
        const result = await RecentList(cursor, PAGE_SIZE)
        const list = result?.list || []
        if (append) {
            items.value.push(...list)
        } else {
            items.value = list
        }
        hasMore.value = result?.has_more || false
        error.value = ''
    } catch (e) {
        if (isUnauthorized(e)) {
            store.sessionExpired()
            pushLogin()
            return
        }
        error.value = typeof e === 'string' ? e : String(e)
    } finally {
        loading.value = false
        loadingMore.value = false
    }
}

const loadStats = async () => {
    try {
        stats.value = await RecentStats()
    } catch (e) {
        if (isUnauthorized(e)) {
            store.sessionExpired()
            pushLogin()
            return
        }
        statsError.value = typeof e === 'string' ? e : String(e)
    }
}

onMounted(() => {
    loadStats()
    load(false)
})
</script>

<style lang="scss" scoped>
.record-container {
    padding: 20px 24px;
    overflow-y: auto;
    height: 100%;
}

.page-header {
    margin-bottom: 20px;

    h2 {
        margin: 0;
        font-size: 22px;
        color: var(--text-primary);
    }

    .page-desc {
        margin: 6px 0 0;
        font-size: 13px;
        color: var(--text-secondary);
    }
}

.stats-panel {
    background: var(--card-bg-light);
    border-radius: 10px;
    padding: 20px;
    margin-bottom: 24px;
}

.stats-head {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 8px;
}

.stats-title {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
}

.stats-note {
    font-size: 12px;
    color: var(--text-placeholder);
}

.stats-fallback {
    font-size: 13px;
    color: var(--text-secondary);
    margin-bottom: 16px;
}

.stat-cards {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;
    margin-top: 16px;
}

.stat-card {
    background: var(--fill-color-light);
    border-radius: 8px;
    padding: 14px 16px;
    display: flex;
    flex-direction: column;
    gap: 4px;

    .stat-value {
        font-size: 24px;
        font-weight: 700;
        color: var(--text-primary);
        line-height: 1.2;

        .stat-unit {
            font-style: normal;
            font-size: 13px;
            font-weight: 400;
            color: var(--text-secondary);
            margin-left: 2px;
        }
    }

    .stat-label {
        font-size: 12px;
        color: var(--text-secondary);
    }
}

.stats-section {
    margin-top: 22px;
}

.section-title {
    margin: 0 0 12px;
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
}

.compose-bar {
    display: flex;
    gap: 2px;
    height: 12px;
    border-radius: 3px;
    overflow: hidden;
}

.legend {
    display: flex;
    flex-wrap: wrap;
    gap: 8px 20px;
    margin-top: 10px;

    .legend-item {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        font-size: 12px;
        color: var(--text-regular);
    }

    .legend-dot {
        width: 10px;
        height: 10px;
        border-radius: 3px;
        display: inline-block;
    }
}

.charts-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;

    @media (max-width: 1100px) {
        grid-template-columns: 1fr;
    }
}

.daily-chart {
    display: flex;
    align-items: flex-end;
    gap: 6px;
    height: 150px;
}

.daily-col {
    flex: 1;
    min-width: 0;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-end;
    gap: 4px;

    .daily-count {
        font-size: 10px;
        line-height: 1;
        color: var(--text-secondary);
    }

    .daily-bar {
        width: 100%;
        border-radius: 3px;
        overflow: hidden;
        display: flex;
        flex-direction: column-reverse;
    }

    .daily-seg {
        width: 100%;
    }

    .daily-empty {
        width: 100%;
        height: 2px;
        border-radius: 1px;
        background: var(--fill-color-dark);
    }

    .daily-date {
        font-size: 10px;
        line-height: 1;
        color: var(--text-placeholder);
    }
}

.hours-head {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 8px;

    .section-title {
        margin-bottom: 0;
    }
}

.hours-chart {
    margin-top: 12px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.hour-row {
    display: flex;
    align-items: center;
    gap: 12px;

    .hour-label {
        width: 88px;
        flex-shrink: 0;
        font-size: 12px;
        color: var(--text-regular);

        i {
            font-style: normal;
            margin-left: 4px;
            color: var(--text-placeholder);
        }
    }

    .hour-track {
        flex: 1;
        height: 16px;
        border-radius: 3px;
        background: var(--fill-color-light);
        overflow: hidden;
    }

    .hour-fill {
        display: flex;
        height: 100%;
        border-radius: 3px;
        overflow: hidden;
    }

    .hour-seg {
        height: 100%;
    }

    .hour-count {
        width: 32px;
        flex-shrink: 0;
        text-align: right;
        font-size: 12px;
        font-variant-numeric: tabular-nums;
        color: var(--text-regular);
    }
}

.list-title {
    margin: 0 0 12px;
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
}

.list-error {
    margin-bottom: 12px;
}

.record-list {
    min-height: 120px;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.record-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 14px 16px;
    background: var(--card-bg-light);
    border-radius: 10px;
    cursor: default;

    .record-cover {
        width: 56px;
        height: 56px;
        flex-shrink: 0;
        border-radius: 6px;
        overflow: hidden;
        background: var(--fill-color);

        :deep(.el-image) {
            width: 100%;
            height: 100%;
            display: block;
        }
    }

    .cover-placeholder {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--text-placeholder);
    }

    .record-body {
        flex: 1;
        min-width: 0;
    }

    .record-title-row {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .record-type {
        flex-shrink: 0;
    }

    .record-title {
        font-size: 15px;
        font-weight: 600;
        color: var(--text-primary);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .record-meta {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 4px 16px;
        margin-top: 8px;
        font-size: 12px;
        color: var(--text-secondary);

        .meta-author {
            max-width: 200px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .meta-progress {
            display: inline-flex;
            align-items: center;
            gap: 8px;
        }

        .prog-track {
            width: 96px;
            height: 6px;
            border-radius: 3px;
            background: var(--fill-color-dark);
            overflow: hidden;
            display: inline-block;
        }

        .prog-fill {
            display: block;
            height: 100%;
            border-radius: 3px;
            background: var(--accent-color);
        }
    }
}

.load-more {
    display: flex;
    justify-content: center;
    margin-top: 16px;
}
</style>
