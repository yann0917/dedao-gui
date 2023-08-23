<template>
    <el-menu router :default-active="activeIndex" class="el-menu" mode="horizontal" text-color="#000"
         @select="handleSelect" :collapse-transition="false">
        <menu-item v-for="(menu, key) in allRoutes" :key="key" :menu="menu" :path="menu.path" />
    </el-menu>
    <!-- active-text-color="#ff5900" -->
</template>
  
<script lang="ts" setup>

import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MenuItem from './MenuItem.vue'

const route = useRoute()
const router = useRouter()

const allRoutes = router.options.routes
const activeIndex = computed(() => {
    return route.path
})
const handleSelect = (key: string, keyPath: string[]) => {
    // console.log(key, keyPath)
}
</script>