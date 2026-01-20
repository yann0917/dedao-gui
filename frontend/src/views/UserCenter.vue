<template>
    <div class="user-center">
        <el-card class="user-card" shadow="hover">
            <!-- 用户信息头部 -->
            <div class="user-header">
                <div class="user-info">
                    <div class="avatar-container">
                        <el-avatar :size="80" :src="user.avatar" fit="cover" class="user-avatar" />
                        <div v-if="user.is_teacher" class="teacher-badge">
                            <el-icon><School /></el-icon>
                            <span>教师</span>
                        </div>
                    </div>
                    <div class="user-meta">
                        <h2 class="nickname">{{ user.nickname }}</h2>
                        <p v-if="ebookUser.slogan" class="slogan">{{ ebookUser.slogan }}</p>
                        <div class="membership-badges">
                            <div v-if="odobUser.user?.is_vip" class="vip-badge odob">
                                <div class="badge-icon">
                                    <el-icon><Headset /></el-icon>
                                </div>
                                <div class="badge-content">
                                    <div class="badge-title">听书会员</div>
                                    <div v-if="odobUser.user.surplus_time" class="badge-desc">
                                        剩余 {{ odobUser.user.surplus_time }} 天
                                    </div>
                                </div>
                            </div>
                            <div v-if="ebookUser.is_vip" class="vip-badge ebook">
                                <div class="badge-icon">
                                    <el-icon><Reading /></el-icon>
                                </div>
                                <div class="badge-content">
                                    <div class="badge-title">电子书会员</div>
                                    <div v-if="ebookUser.surplus_time" class="badge-desc">
                                        剩余 {{ ebookUser.surplus_time }} 天
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 学习数据 -->
            <div class="study-stats">
                <div class="stat-item">
                    <span class="stat-value">{{ (user.today_study_time / 60).toFixed(0) }}</span>
                    <span class="stat-label">今日学习(分钟)</span>
                </div>
                <div class="stat-item">
                    <span class="stat-value">{{ user.study_serial_days }}</span>
                    <span class="stat-label">连续学习(天)</span>
                </div>
            </div>

            <!-- 听书会员信息 -->
            <div v-if="odobUser.user.is_vip" class="membership-section">
                <div class="section-header">
                    <h3>听书会员</h3>
                </div>
                <el-alert
                    v-if="odobUser.user.err_tips"
                    :title="odobUser.user.err_tips"
                    type="warning"
                    :closable="false"
                    class="membership-alert"
                />
                <div class="membership-info">
                    <div class="info-row">
                        <div class="expire-info">
                            <div class="expire-date">
                                <el-icon><Clock /></el-icon>
                                <span class="expire-label">到期时间</span>
                                <span class="expire-value">{{ timestampToTime(odobUser.user.expire_time) }}</span>
                            </div>
                            <el-tag type="danger" effect="plain" class="surplus-tag">剩 {{ odobUser.user.surplus_time }} 天</el-tag>
                        </div>
                    </div>
                    <div class="info-row">
                        <div class="stats-grid">
                            <div class="stat-card">
                                <div class="stat-number">{{ odobUser.user.week_count }}</div>
                                <div class="stat-text">本周听书</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-number">{{ odobUser.user.total_count }}</div>
                                <div class="stat-text">累计听书</div>
                            </div>
                        </div>
                    </div>
                    <div class="info-row">
                        <div class="savings-highlight">
                            <el-icon><Present /></el-icon>
                            <span class="savings-text">累计为你节省了</span>
                            <span class="savings-amount">{{ odobUser.user.save_price }}{{ odobUser.user.price_desc }}</span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 电子书会员信息 -->
            <div v-if="ebookUser.is_vip" class="membership-section">
                <div class="section-header">
                    <h3>电子书会员</h3>
                </div>
                <el-alert
                    v-if="ebookUser.err_tips"
                    :title="ebookUser.err_tips"
                    type="warning"
                    :closable="false"
                    class="membership-alert"
                />
                <div class="membership-info">
                    <div class="info-row">
                        <div class="expire-info">
                            <div class="expire-date">
                                <el-icon><Clock /></el-icon>
                                <span class="expire-label">到期时间</span>
                                <span class="expire-value">{{ timestampToTime(ebookUser.expire_time) }}</span>
                            </div>
                            <el-tag type="danger" effect="plain" class="surplus-tag">剩 {{ ebookUser.surplus_time }} 天</el-tag>
                        </div>
                    </div>
                    <div class="info-row">
                        <div class="stats-grid">
                            <div class="stat-card">
                                <div class="stat-number">{{ ebookUser.month_count }}</div>
                                <div class="stat-text">本月读书</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-number">{{ ebookUser.total_count }}</div>
                                <div class="stat-text">累计读书</div>
                            </div>
                        </div>
                    </div>
                    <div class="info-row">
                        <div class="savings-highlight">
                            <el-icon><Present /></el-icon>
                            <span class="savings-text">累计为你节省了</span>
                            <span class="savings-amount">{{ ebookUser.save_price }}{{ ebookUser.price_desc }}</span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 退出按钮 -->
            <div class="logout-section">
                <el-button
                    type="danger"
                    size="large"
                    plain
                    @click="handleLogout"
                    class="logout-btn"
                >
                    <el-icon><SwitchButton /></el-icon>
                    <span>退出登录</span>
                </el-button>
            </div>
        </el-card>
    </div>
</template>

<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { SwitchButton, Clock, Present, School, Headset, Reading } from '@element-plus/icons-vue'

import { UserInfo, EbookUserInfo, OdobUserInfo } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { userStore } from '../stores/user';
import { timestampToTime } from '../utils/utils'
const store = userStore()
let user = reactive(new services.User)
let ebookUser = reactive(new services.EbookVIPInfo)
let odobUser = reactive(new services.OdobVip)
odobUser.user = new services.OdobUser

const handleLogout = async () => {
    try {
        await store.logout()
        ElMessage.success('已退出登录')
    } catch (error) {
        ElMessage.error('退出失败，请重试')
        console.error('Logout error:', error)
    }
}

onMounted(() => {

})

const getUserInfo = async () => {
    UserInfo().then(result => {
        Object.assign(user, result)
        store.user = user
        console.log(store)
    })
}
getUserInfo()

const getEbookUserInfo = async () => {
    await EbookUserInfo().then(result => {
        console.log(result)
        Object.assign(ebookUser, result)
    }).catch((error) => {
        ElMessage({
            message: error,
            type: 'warning'
        })
    })
}
getEbookUserInfo()

const getOdobUserInfo = async () => {
    await OdobUserInfo().then(result => {
        console.log(result)
        Object.assign(odobUser, result)
    }).catch((error) => {
        ElMessage({
            message: error,
            type: 'warning'
        })
    })
}

getOdobUserInfo()

</script>
<style scoped>
.user-center {
    padding: 24px;
    background: var(--fill-color-light);
    height: calc(100vh - 60px);
    overflow-y: auto;
    box-sizing: border-box;
    
    /* 隐藏滚动条但保留功能 - 清新风格 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE 10+ */
}

.user-center::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}

.user-card {
    max-width: 800px;
    margin: 0 auto;
    border-radius: 12px;
}

.user-header {
    margin-bottom: 30px;
}

.user-info {
    display: flex;
    align-items: center;
    gap: 24px;
}

.user-avatar {
    border: 4px solid rgba(255, 107, 0, 0.1);
    transition: transform 0.3s ease;
}

.user-avatar:hover {
    transform: scale(1.05);
}

.avatar-container {
    position: relative;
}

.teacher-badge {
    position: absolute;
    top: -8px;
    right: -8px;
    background: var(--accent-color);
    color: white;
    padding: 4px 8px;
    border-radius: 16px;
    font-size: 12px;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 4px;
    box-shadow: var(--shadow-light);
    z-index: 10;
}

.user-meta {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.nickname {
    font-size: 28px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
    letter-spacing: -0.5px;
}

.slogan {
    color: var(--text-secondary);
    font-size: 15px;
    margin: 0;
    line-height: 1.6;
}

.membership-badges {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
}

.vip-badge {
    display: flex;
    align-items: center;
    padding: 8px 16px;
    border-radius: 24px;
    gap: 8px;
    transition: all 0.3s ease;
    cursor: default;
    color: white;
    font-size: 14px;
    font-weight: 500;
}

.vip-badge:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-light);
}

.vip-badge.odob {
    background: linear-gradient(135deg, #ff6b35 0%, #ff9800 100%);
}

.vip-badge.ebook {
    background: linear-gradient(135deg, #e91e63 0%, #f44336 100%);
}

.badge-icon {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
}

.badge-icon .el-icon {
    font-size: 14px;
}

.badge-content {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.badge-title {
    font-size: 14px;
    font-weight: 600;
}

.badge-desc {
    font-size: 11px;
    opacity: 0.9;
    font-weight: 400;
}

.study-stats {
    display: flex;
    align-items: center;
    gap: 40px;
    margin-bottom: 30px;
    padding: 20px;
    background: var(--fill-color-light);
    border-radius: 12px;
    border: 1px solid var(--border-soft);
}

.stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
}

.stat-value {
    font-size: 28px;
    font-weight: 500;
    color: var(--accent-color);
}

.stat-label {
    font-size: 14px;
    color: var(--text-secondary);
}

.membership-section {
    margin-top: 30px;
    padding-top: 20px;
    border-top: 1px solid var(--border-soft);
}

.section-header {
    margin-bottom: 20px;
}

.section-header h3 {
    font-size: 18px;
    font-weight: 500;
    color: var(--text-primary);
    margin: 0;
}

.membership-alert {
    margin-bottom: 16px;
}

/* 暗色主题下的 alert 适配 */
.theme-dark .membership-alert {
    background-color: var(--card-bg) !important;
    border-color: var(--border-soft) !important;
    color: var(--text-primary) !important;
}

.theme-dark .membership-alert .el-alert__title {
    color: var(--text-primary) !important;
}

.theme-dark .membership-alert .el-alert__content {
    color: var(--text-secondary) !important;
}

.theme-dark .membership-alert .el-alert__icon {
    color: var(--accent-color) !important;
}

.membership-info {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.info-row {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
}

.el-tag {
    padding: 8px 16px;
    font-size: 14px;
    border-radius: 6px;
}

/* 新的样式 */
.expire-info {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    background: var(--fill-color-light);
    border-radius: 12px;
    border: 1px solid var(--border-soft);
    flex: 1;
}

.expire-date {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
}

.expire-date .el-icon {
    font-size: 18px;
    color: var(--accent-color);
}

.expire-label {
    font-size: 14px;
    color: var(--text-secondary);
}

.expire-value {
    font-size: 16px;
    font-weight: 500;
    color: var(--text-primary);
}

.surplus-tag {
    flex-shrink: 0;
    font-size: 13px;
    font-weight: 500;
}

.stats-grid {
    display: flex;
    gap: 16px;
    width: 100%;
}

.stat-card {
    flex: 1;
    padding: 20px;
    background: var(--fill-color-light);
    border-radius: 12px;
    border: 1px solid var(--border-soft);
    text-align: center;
    transition: all 0.3s ease;
}

.stat-card:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-light);
}

.stat-card .stat-number {
    font-size: 28px;
    font-weight: 600;
    color: var(--accent-color);
    margin-bottom: 8px;
}

.stat-card .stat-text {
    font-size: 14px;
    color: var(--text-secondary);
}

.savings-highlight {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: linear-gradient(135deg, rgba(67, 160, 71, 0.1) 0%, rgba(67, 160, 71, 0.05) 100%);
    border-radius: 12px;
    border: 1px solid rgba(67, 160, 71, 0.2);
    width: 100%;
}

.savings-highlight .el-icon {
    font-size: 24px;
    color: #43a047;
}

.savings-text {
    font-size: 14px;
    color: var(--text-secondary);
}

.savings-amount {
    font-size: 18px;
    font-weight: 600;
    color: #43a047;
    margin-left: auto;
}

:deep(.el-alert) {
    border-radius: 8px;
    margin: 16px 0;
}

:deep(.el-tag .el-icon) {
    margin-right: 4px;
}

/* 暗色主题适配 */
.theme-dark .study-stats {
    background: var(--fill-color-light) !important;
    border-color: var(--border-soft) !important;
}

.theme-dark .stat-value {
    color: var(--accent-color) !important;
}

.theme-dark .stat-label {
    color: var(--text-secondary) !important;
}

.theme-dark .teacher-badge {
    background: var(--accent-color) !important;
}

.theme-dark .slogan {
    color: var(--text-secondary) !important;
}

.theme-dark .section-header h3 {
    color: var(--text-primary) !important;
}

.theme-dark .membership-section {
    border-top-color: var(--border-soft) !important;
}

.theme-dark .expire-info {
    background: var(--fill-color-light) !important;
    border-color: var(--border-soft) !important;
}

.theme-dark .stat-card {
    background: var(--fill-color-light) !important;
    border-color: var(--border-soft) !important;
}

.theme-dark .savings-highlight {
    background: linear-gradient(135deg, rgba(67, 160, 71, 0.15) 0%, rgba(67, 160, 71, 0.08) 100%) !important;
    border-color: rgba(67, 160, 71, 0.25) !important;
}

.theme-dark .savings-highlight .el-icon {
    color: #66bb6a !important;
}

.logout-section {
    margin-top: 30px;
    padding-top: 20px;
    border-top: 1px solid var(--border-soft);
    display: flex;
    justify-content: center;
}

.logout-btn {
    padding: 12px 32px;
    font-size: 16px;
    border-radius: 8px;
}

.logout-btn .el-icon {
    margin-right: 8px;
}

.logout-btn:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-medium);
}

.logout-btn:active {
    transform: translateY(0);
}

/* 暗色主题退出按钮样式 */
.theme-dark .logout-btn {
    --el-button-bg-color: transparent;
    --el-button-border-color: #f56c6c;
    --el-button-text-color: #f56c6c;
    --el-button-hover-bg-color: rgba(245, 108, 108, 0.1);
    --el-button-hover-border-color: #ff5252;
    --el-button-hover-text-color: #ff5252;
    --el-button-active-bg-color: rgba(245, 108, 108, 0.2);
    --el-button-active-border-color: #ff5252;
    --el-button-active-text-color: #ff5252;
}
</style>