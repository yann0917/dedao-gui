<template>
    <el-pagination background
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
import { ref, reactive, defineEmits, defineProps } from 'vue';

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