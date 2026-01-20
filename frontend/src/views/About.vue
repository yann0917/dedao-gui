<template>
  <div class="about-page">
    <div class="about-content">
      <h2>Theme Settings</h2>
      <p>Customize your experience</p>
      <div class="color-picker-section">
        <span>System Color: </span>
        <el-color-picker v-model="systemColor" :predefine="predefineColors" show-alpha @change="setThemeColor('')" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import {setThemeColor} from '../utils/utils'

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
<style scoped>
.about-page {
  padding: 40px;
  height: 100%;
  box-sizing: border-box;
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.about-content {
  background: var(--card-bg);
  padding: 40px;
  border-radius: 16px;
  box-shadow: var(--shadow-soft);
  text-align: center;
  min-width: 300px;
}

h2 {
  color: var(--text-primary);
  margin-bottom: 10px;
}

p {
  color: var(--text-secondary);
  margin-bottom: 30px;
}

.color-picker-section {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-regular);
}
</style>