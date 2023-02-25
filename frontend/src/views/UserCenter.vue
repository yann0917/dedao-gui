<template>
    <el-space wrap>
        <el-card v-for="i in 1" :key="i" class="box-card" width="100%" shadow="hover">
            <template #header>
                <div class="card-header">
                    <el-row :gutter="5" align="middle">
                        <el-col :span="6">
                            <el-avatar :size="72" :src="user.avatar" fit="fill" />
                            <el-button text>{{ user.nickname }}</el-button>
                        </el-col>
                        <el-col :span="18" style="text-align: left; line-height: 10px;">
                            <p>{{ ebookUser.slogan }}</p>
                            <el-tag class="ml-2" type="warning" v-if="odobUser.user?.is_vip" round>
                                <el-icon>
                                    <HotWater />
                                </el-icon>听书会员
                            </el-tag>
                            <el-tag class="ml-2" type="danger" v-if="ebookUser.is_vip" round>
                                <el-icon>
                                    <HotWater />
                                </el-icon>电子书会员
                            </el-tag>
                        </el-col>
                    </el-row>
                </div>
            </template>
            <div class="study">
                <el-tag class="ml-2" type="success">今日学习{{ (user.today_study_time / 60).toFixed(0) }}分钟</el-tag>
                <el-tag class="ml-2" type="warning">连续{{ user.study_serial_days }}天</el-tag>
                <el-tag class="ml-2" type="info" v-if="user.is_teacher">教师</el-tag>
            </div>

            <div v-if="odobUser.user.is_vip" style="text-align: left;">
                <el-divider content-position="left">听书会员</el-divider>
                <el-alert v-if="odobUser.user.err_tips" :title="odobUser.user.err_tips" type="warning"
                    :closable="false" />
                <el-tag class="ml-2" type="danger">听书VIP·{{ timestampToTime(odobUser.user.expire_time) }}到期</el-tag>
                <el-tag class="ml-2" type="danger">距到期还有{{ odobUser.user.surplus_time }}天</el-tag><br />
                <el-tag class="ml-2" type="success">本周听书{{ odobUser.user.week_count }}本 /
                    累计听书{{ odobUser.user.total_count }}本</el-tag><br />
                <el-tag class="ml-2" type="warning">累计为你节省了{{ odobUser.user.save_price }}{{ odobUser.user.price_desc }}
                </el-tag>
            </div>
            <div v-if="ebookUser.is_vip" style="text-align: left;">
                <el-divider content-position="left">电子书会员</el-divider>
                <el-alert v-if="ebookUser.err_tips" :title="ebookUser.err_tips" type="warning" :closable="false" />
                <el-tag class="ml-2" type="danger">电子书VIP·{{ timestampToTime(ebookUser.expire_time) }}到期</el-tag>
                <el-tag class="ml-2" type="danger">距到期还有{{ ebookUser.surplus_time }}天</el-tag><br />
                <el-tag class="ml-2" type="success">本月读书{{ ebookUser.month_count }}本 / 累计读书{{ ebookUser.total_count }}本
                </el-tag><br />
                <el-tag class="ml-2" type="warning">累计为你节省了{{ ebookUser.save_price }}{{ ebookUser.price_desc }}</el-tag>
            </div>

        </el-card>
    </el-space>

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
.el-tag{
    margin-right: 5px;
    margin-bottom: 5px;
}
</style>