<template>
  <qr-login v-if="dialogVisible" :dialog-visible="dialogVisible" @close="closeDialog"></qr-login>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import QrLogin from "../components/QrLogin.vue";

const dialogVisible = ref(false);

onMounted(() => {
  openDialog()
});

const openDialog = () => {
  dialogVisible.value = true;
};

const closeDialog = () => {
  dialogVisible.value = false;
};
</script>
