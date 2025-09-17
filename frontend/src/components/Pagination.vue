<template>
    <el-pagination 
        background
        class="custom-pagination"
        v-model:currentPage="pageParams.page"
        v-model:page-size="pageParams.pageSize"
        :page-sizes="pageSizes"
        :total="total"
        :layout="layout"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
    />
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const emits  = defineEmits(['pageChange']);
const pageSizes = ref([10,15,20,30,50]);

const props = defineProps({
        total: {
            type: Number,
            default: 0
        },
        layout:{
            type:String,
            default: 'total, sizes, prev, pager, next, jumper'
        }
    })
const pageParams = reactive({
    page: 1,
    pageSize: 15,
})


//改变每页展示的数据条数 
const handleSizeChange = (val:any)=>{
    pageParams.pageSize = val;
    emits('pageChange', pageParams);
}
//改变当前页的页码
const handleCurrentChange = (val:any)=>{
    pageParams.page = val;
    emits('pageChange', pageParams);
}
</script>

<style scoped>
/* 自定义分页器样式 */
.custom-pagination {
    margin-top: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
}

.custom-pagination :deep(.el-pagination__total) {
    color: var(--text-secondary);
    font-size: 14px;
    margin-right: 16px;
}

.custom-pagination :deep(.el-pagination__sizes) {
    margin-right: 16px;
}

.custom-pagination :deep(.el-pagination__sizes .el-select .el-input__inner) {
    background-color: var(--card-bg);
    border-color: var(--border-soft);
    color: var(--text-primary);
    transition: all 0.3s ease;
}

.custom-pagination :deep(.el-pagination__sizes .el-select .el-input__inner:hover) {
    border-color: var(--accent-color);
}

.custom-pagination :deep(.el-pagination__sizes .el-select .el-input__inner:focus) {
    border-color: var(--accent-color);
    box-shadow: 0 0 0 2px rgba(255, 107, 0, 0.2);
}

.custom-pagination :deep(.el-pager li) {
    background-color: var(--card-bg);
    color: var(--text-primary);
    border: 1px solid var(--border-soft);
    margin: 0 2px;
    border-radius: 4px;
    transition: all 0.3s ease;
    min-width: 32px;
    height: 32px;
    line-height: 30px;
}

.custom-pagination :deep(.el-pager li:hover) {
    color: var(--accent-color);
    background-color: var(--card-hover-bg);
    border-color: var(--accent-color);
    transform: translateY(-1px);
    box-shadow: var(--shadow-soft);
}

.custom-pagination :deep(.el-pager li.is-active) {
    background-color: var(--accent-color);
    color: #fff;
    border-color: var(--accent-color);
    font-weight: 500;
}

.custom-pagination :deep(.el-pager li.is-active:hover) {
    background-color: var(--accent-hover);
    border-color: var(--accent-hover);
    color: #fff;
}

.custom-pagination :deep(.btn-prev),
.custom-pagination :deep(.btn-next) {
    background-color: var(--card-bg);
    color: var(--text-primary);
    border: 1px solid var(--border-soft);
    margin: 0 2px;
    border-radius: 4px;
    transition: all 0.3s ease;
    min-width: 32px;
    height: 32px;
    line-height: 30px;
}

.custom-pagination :deep(.btn-prev:hover),
.custom-pagination :deep(.btn-next:hover) {
    color: var(--accent-color);
    background-color: var(--card-hover-bg);
    border-color: var(--accent-color);
    transform: translateY(-1px);
    box-shadow: var(--shadow-soft);
}

.custom-pagination :deep(.btn-prev:disabled),
.custom-pagination :deep(.btn-next:disabled) {
    background-color: var(--fill-color);
    color: var(--text-tertiary);
    border-color: var(--border-soft);
    cursor: not-allowed;
}

.custom-pagination :deep(.btn-prev:disabled:hover),
.custom-pagination :deep(.btn-next:disabled:hover) {
    background-color: var(--fill-color);
    color: var(--text-tertiary);
    border-color: var(--border-soft);
    transform: none;
    box-shadow: none;
}

.custom-pagination :deep(.el-pagination__jump) {
    margin-left: 16px;
    color: var(--text-primary);
}

.custom-pagination :deep(.el-pagination__jump .el-input__inner) {
    background-color: var(--card-bg);
    border-color: var(--border-soft);
    color: var(--text-primary);
    width: 50px;
    text-align: center;
    transition: all 0.3s ease;
}

.custom-pagination :deep(.el-pagination__jump .el-input__inner:hover) {
    border-color: var(--accent-color);
}

.custom-pagination :deep(.el-pagination__jump .el-input__inner:focus) {
    border-color: var(--accent-color);
    box-shadow: 0 0 0 2px rgba(255, 107, 0, 0.2);
}

/* 暗色模式下的分页器样式 */
.theme-dark .custom-pagination :deep(.el-pagination__total) {
    color: var(--text-secondary) !important;
}

.theme-dark .custom-pagination :deep(.el-pagination__sizes .el-select .el-input__inner) {
    background-color: var(--card-bg) !important;
    border-color: var(--border-soft) !important;
    color: var(--text-primary) !important;
}

.theme-dark .custom-pagination :deep(.el-pagination__sizes .el-select .el-input__inner:hover) {
    border-color: var(--accent-color) !important;
}

.theme-dark .custom-pagination :deep(.el-pagination__sizes .el-select .el-input__inner:focus) {
    border-color: var(--accent-color) !important;
    box-shadow: 0 0 0 2px rgba(255, 107, 0, 0.2) !important;
}

.theme-dark .custom-pagination :deep(.el-pager li) {
    background-color: var(--card-bg) !important;
    color: var(--text-primary) !important;
    border-color: var(--border-soft) !important;
}

.theme-dark .custom-pagination :deep(.el-pager li:hover) {
    color: var(--accent-color) !important;
    background-color: var(--card-hover-bg) !important;
    border-color: var(--accent-color) !important;
}

.theme-dark .custom-pagination :deep(.el-pager li.is-active) {
    background-color: var(--accent-color) !important;
    color: #fff !important;
    border-color: var(--accent-color) !important;
}

.theme-dark .custom-pagination :deep(.el-pager li.is-active:hover) {
    background-color: var(--accent-hover) !important;
    border-color: var(--accent-hover) !important;
    color: #fff !important;
}

.theme-dark .custom-pagination :deep(.btn-prev),
.theme-dark .custom-pagination :deep(.btn-next) {
    background-color: var(--card-bg) !important;
    color: var(--text-primary) !important;
    border-color: var(--border-soft) !important;
}

.theme-dark .custom-pagination :deep(.btn-prev:hover),
.theme-dark .custom-pagination :deep(.btn-next:hover) {
    color: var(--accent-color) !important;
    background-color: var(--card-hover-bg) !important;
    border-color: var(--accent-color) !important;
}

.theme-dark .custom-pagination :deep(.btn-prev:disabled),
.theme-dark .custom-pagination :deep(.btn-next:disabled) {
    background-color: var(--fill-color) !important;
    color: var(--text-tertiary) !important;
    border-color: var(--border-soft) !important;
}

.theme-dark .custom-pagination :deep(.btn-prev:disabled:hover),
.theme-dark .custom-pagination :deep(.btn-next:disabled:hover) {
    background-color: var(--fill-color) !important;
    color: var(--text-tertiary) !important;
    border-color: var(--border-soft) !important;
}

.theme-dark .custom-pagination :deep(.el-pagination__jump) {
    color: var(--text-primary) !important;
}

.theme-dark .custom-pagination :deep(.el-pagination__jump .el-input__inner) {
    background-color: var(--card-bg) !important;
    border-color: var(--border-soft) !important;
    color: var(--text-primary) !important;
}

.theme-dark .custom-pagination :deep(.el-pagination__jump .el-input__inner:hover) {
    border-color: var(--accent-color) !important;
}

.theme-dark .custom-pagination :deep(.el-pagination__jump .el-input__inner:focus) {
    border-color: var(--accent-color) !important;
    box-shadow: 0 0 0 2px rgba(255, 107, 0, 0.2) !important;
}
</style>