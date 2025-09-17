<template>
    <div class="user-center">
        <el-card class="user-card" shadow="hover">
            <!-- 用户信息头部 -->
            <div class="user-header">
                <div class="user-info">
                    <el-avatar :size="80" :src="user.avatar" fit="cover" class="user-avatar" />
                    <div class="user-meta">
                        <h2 class="nickname">{{ user.nickname }}</h2>
                        <p class="slogan">{{ ebookUser.slogan }}</p>
                        <div class="vip-tags">
                            <el-tag v-if="odobUser.user?.is_vip" class="vip-tag" type="warning" effect="dark" round>
                                <el-icon><HotWater /></el-icon>
                                <span>听书会员</span>
                            </el-tag>
                            <el-tag v-if="ebookUser.is_vip" class="vip-tag" type="danger" effect="dark" round>
                                <el-icon><HotWater /></el-icon>
                                <span>电子书会员</span>
                            </el-tag>
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
                <el-tag v-if="user.is_teacher" class="teacher-tag" type="info">教师</el-tag>
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
                        <el-tag type="danger" effect="plain">听书VIP·{{ timestampToTime(odobUser.user.expire_time) }}到期</el-tag>
                        <el-tag type="danger" effect="plain">距到期还有{{ odobUser.user.surplus_time }}天</el-tag>
                    </div>
                    <div class="info-row">
                        <el-tag type="success" effect="plain">
                            本周听书{{ odobUser.user.week_count }}本 / 累计听书{{ odobUser.user.total_count }}本
                        </el-tag>
                    </div>
                    <div class="info-row">
                        <el-tag type="warning" effect="plain">
                            累计为你节省了{{ odobUser.user.save_price }}{{ odobUser.user.price_desc }}
                        </el-tag>
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
                        <el-tag type="danger" effect="plain">电子书VIP·{{ timestampToTime(ebookUser.expire_time) }}到期</el-tag>
                        <el-tag type="danger" effect="plain">距到期还有{{ ebookUser.surplus_time }}天</el-tag>
                    </div>
                    <div class="info-row">
                        <el-tag type="success" effect="plain">
                            本月读书{{ ebookUser.month_count }}本 / 累计读书{{ ebookUser.total_count }}本
                        </el-tag>
                    </div>
                    <div class="info-row">
                        <el-tag type="warning" effect="plain">
                            累计为你节省了{{ ebookUser.save_price }}{{ ebookUser.price_desc }}
                        </el-tag>
                    </div>
                </div>
            </div>
        </el-card>
    </div>
</template>

<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

import { UserInfo, EbookUserInfo, OdobUserInfo } from '../../wailsjs/go/backend/App'
import { services } from '../../wailsjs/go/models'
import { userStore } from '../stores/user';
import { timestampToTime } from '../utils/utils'
const store = userStore()
let user = reactive(new services.User)
let ebookUser = reactive(new services.EbookVIPInfo)
let odobUser = reactive(new services.OdobVip)
odobUser.user = new services.OdobUser

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
    min-height: calc(100vh - 60px);
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
    align-items: flex-start;
    gap: 24px;
}

.user-avatar {
    border: 4px solid rgba(255, 107, 0, 0.1);
    transition: transform 0.3s ease;
}

.user-avatar:hover {
    transform: scale(1.05);
}

.user-meta {
    flex: 1;
}

.nickname {
    font-size: 24px;
    font-weight: 500;
    color: var(--text-primary);
    margin: 0 0 8px;
}

.slogan {
    color: var(--text-secondary);
    font-size: 14px;
    margin: 0 0 12px;
}

.vip-tags {
    display: flex;
    gap: 12px;
}

.vip-tag {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 6px 12px;
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

.teacher-tag {
    margin-left: auto;
    padding: 8px 16px;
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

.theme-dark .section-header h3 {
    color: var(--text-primary) !important;
}

.theme-dark .membership-section {
    border-top-color: var(--border-soft) !important;
}
</style>