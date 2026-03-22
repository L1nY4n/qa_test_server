<template>
  <div class="overview-page">
    <div class="toolbar">
      <div>
        <h3>运行总览</h3>
        <p>聚合设备在线状态与服务健康状态</p>
      </div>
      <a-button type="primary" :loading="state.loading" @click="refresh">刷新</a-button>
    </div>

    <a-row :gutter="[12, 12]">
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="设备总数" :value="state.stats?.total ?? 0" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="在线设备" :value="state.stats?.online ?? 0" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="离线设备" :value="state.stats?.offline ?? 0" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="服务运行时长(秒)" :value="state.health?.uptimeSeconds ?? 0" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[12, 12]" class="section-row">
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card title="服务健康状态" size="small">
          <a-descriptions :column="1" size="small" bordered>
            <a-descriptions-item label="数据库">
              <a-tag :color="state.health?.db?.ok ? 'green' : 'red'">
                {{ state.health?.db?.ok ? '正常' : '异常' }}
              </a-tag>
              <span class="inline-note">{{ state.health?.db?.message || '--' }}</span>
            </a-descriptions-item>
            <a-descriptions-item label="HTTP 地址">
              {{ state.health?.server?.httpAddr || '--' }}
            </a-descriptions-item>
            <a-descriptions-item label="TCP 地址">
              {{ state.health?.server?.tcpAddr || '--' }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card title="菜单能力规划" size="small">
          <a-space direction="vertical" :size="6">
            <a-tag color="blue">监控中心：总览 / 设备监控</a-tag>
            <a-tag color="purple">运维分析：PHM / 告警中心 / 趋势分析</a-tag>
            <a-tag color="geekblue">策略与系统：规则策略 / 系统状态</a-tag>
          </a-space>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import * as API from '@/api'
import { onMounted, reactive } from 'vue'
import type { DeviceStats } from '@/api/device'
import type { SystemHealth } from '@/api/system'

const state = reactive<{
  loading: boolean
  stats: DeviceStats | null
  health: SystemHealth | null
}>({
  loading: false,
  stats: null,
  health: null,
})

const refresh = async () => {
  state.loading = true
  try {
    const [stats, health] = await Promise.all([
      API.device.stats(30),
      API.system.health(),
    ])
    state.stats = stats
    state.health = health
  } finally {
    state.loading = false
  }
}

onMounted(() => {
  void refresh()
})
</script>

<style scoped>
.overview-page {
  padding: 12px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 12px;
}

.toolbar h3 {
  margin: 0;
  color: #27445b;
  font-size: 16px;
}

.toolbar p {
  margin: 4px 0 0;
  color: #7e93a8;
  font-size: 12px;
}

.section-row {
  margin-top: 2px;
}

.inline-note {
  margin-left: 8px;
  color: #687f96;
  font-size: 12px;
}

@media (max-width: 900px) {
  .overview-page {
    padding: 8px;
  }
}
</style>
