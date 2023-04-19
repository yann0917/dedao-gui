<template>
  <!-- <div class="theme-color-popper"> -->
    <el-drawer v-model="drawer" title="I am the title" :with-header="false">
      <span>With default value</span>
    <el-color-picker v-model="systemColor" :predefine="predefineColors" show-alpha @change="setThemeColor('')" />
  </el-drawer>
  <!-- </div> -->
</template>
  
<script lang="ts" setup>
import { ref } from 'vue';
import {setThemeColor} from '../utils/utils'

const drawer = ref(true)
const systemColor = ref('rgba(19, 206, 102, 0.8)')
const predefineColors = ref([
  '#ff4500',
  '#ff8c00',
  '#ffd700',
  '#90ee90',
  '#00ced1',
  '#1e90ff',
  '#c71585',
  'rgba(255, 69, 0, 0.68)',
  'rgb(255, 120, 0)',
  'hsv(51, 100, 98)',
  'hsva(120, 40, 94, 0.5)',
  'hsl(181, 100%, 37%)',
  'hsla(209, 100%, 56%, 0.73)',
  '#c7158577',
])

</script>
<style  scoped lang="scss">
.theme-color-popper {
  display: block;
  align-items: center;
  width: 230px;
  padding: 20px;
  background-color: greenyellow;
}
</style>