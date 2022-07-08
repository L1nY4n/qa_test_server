<template>
  <!-- <ConfigProvider> -->
  <router-view #="{ Component }">
    <component :is="Component" />
  </router-view>
  <!-- </ConfigProvider> -->

</template>

<script setup lang="ts">
import { watchEffect } from 'vue';
import { useRoute } from 'vue-router';
import { ConfigProvider } from 'ant-design-vue';
</script>