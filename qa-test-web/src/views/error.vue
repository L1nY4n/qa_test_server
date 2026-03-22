<template>
  <div class="error-page">
    <a-result
      status="404"
      title="页面不存在"
      sub-title="请检查菜单路径或返回监控首页。"
    >
      <template #extra>
        <a-button type="primary" @click="goHome">返回总览</a-button>
      </template>
    </a-result>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  void router.push({ name: 'overview' })
}
</script>

<style scoped>
.error-page {
  height: 100%;
  display: grid;
  place-items: center;
}
</style>
