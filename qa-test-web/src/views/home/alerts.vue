<template>
  <div class="page">
    <a-card title="告警中心（规划版）" size="small">
      <a-alert
        type="info"
        show-icon
        message="当前已完成菜单重构，告警模块建议下一步接入离线阈值、恢复检测与通知渠道。"
      />
      <a-row :gutter="[12, 12]" class="row">
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-statistic title="离线设备候选告警" :value="stats?.offline ?? 0" />
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-statistic title="在线设备" :value="stats?.online ?? 0" />
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-statistic title="总设备数" :value="stats?.total ?? 0" />
        </a-col>
      </a-row>
      <a-divider />
      <a-space direction="vertical" :size="6">
        <a-tag color="red">P1：设备离线告警与恢复告警</a-tag>
        <a-tag color="orange">P2：告警确认、静默窗口与抑制策略</a-tag>
        <a-tag color="gold">P3：邮件/企业微信/Webhook 多通道通知</a-tag>
      </a-space>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import * as API from '@/api'
import { onMounted, ref } from 'vue'
import type { DeviceStats } from '@/api/device'

const stats = ref<DeviceStats | null>(null)

onMounted(async () => {
  stats.value = await API.device.stats(30)
})
</script>

<style scoped>
.page {
  padding: 12px;
}

.row {
  margin-top: 12px;
}

@media (max-width: 900px) {
  .page {
    padding: 8px;
  }
}
</style>
