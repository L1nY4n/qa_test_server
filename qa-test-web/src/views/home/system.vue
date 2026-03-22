<template>
  <div class="page">
    <a-card title="系统状态" size="small">
      <div class="toolbar">
        <a-tag :color="health?.db?.ok ? 'green' : 'red'">
          {{ health?.db?.ok ? '数据库正常' : '数据库异常' }}
        </a-tag>
        <a-space>
          <a-button type="primary" size="small" :loading="loading" @click="refresh()">刷新</a-button>
          <a-button
            v-if="isAdmin"
            danger
            size="small"
            :loading="resetLoading"
            @click="confirmCloudReset"
          >
            云系统重置
          </a-button>
        </a-space>
      </div>

      <a-row :gutter="[12, 12]">
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-statistic title="运行时长(秒)" :value="health?.uptimeSeconds ?? 0" />
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-statistic title="HTTP 地址" :value="health?.server?.httpAddr ?? '--'" />
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-statistic title="TCP 地址" :value="health?.server?.tcpAddr ?? '--'" />
        </a-col>
      </a-row>

      <a-divider />
      <a-collapse ghost>
        <a-collapse-panel key="ws" header="WebSocket 状态明细">
          <ParamViewer :node="health?.websocket || {}" />
        </a-collapse-panel>
      </a-collapse>
    </a-card>

    <a-card title="服务器性能负载监测" size="small" class="card">
      <div class="perf-toolbar">
        <a-space wrap>
          <a-tag color="blue">采样周期 {{ metrics?.sampleEverySeconds ?? 0 }} 秒</a-tag>
          <a-tag color="default">数据点 {{ metrics?.points?.length ?? 0 }}</a-tag>
          <a-tag color="geekblue">磁盘路径 {{ metrics?.diskPath || '--' }}</a-tag>
        </a-space>
        <a-button size="small" :loading="loading" @click="refresh()">刷新负载</a-button>
      </div>

      <a-alert class="perf-host" type="info" show-icon :message="hostSummary" />

      <a-row :gutter="[12, 12]">
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic title="CPU 使用率" :value="toNumber(currentMetric?.cpuPercent)" :precision="2" suffix="%" />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic
            title="内存使用率"
            :value="toNumber(currentMetric?.memoryUsedPercent)"
            :precision="2"
            suffix="%"
          />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic
            title="磁盘使用率"
            :value="toNumber(currentMetric?.diskUsedPercent)"
            :precision="2"
            suffix="%"
          />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic title="Goroutines" :value="toNumber(currentMetric?.goroutines)" />
        </a-col>
      </a-row>

      <a-row :gutter="[12, 12]" class="mt-row">
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic
            title="下行流量"
            :value="toKBps(currentMetric?.netRecvBps)"
            :precision="2"
            suffix="KB/s"
          />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic
            title="上行流量"
            :value="toKBps(currentMetric?.netSentBps)"
            :precision="2"
            suffix="KB/s"
          />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic
            title="进程 CPU"
            :value="toNumber(currentMetric?.processCpuPercent)"
            :precision="2"
            suffix="%"
          />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic title="进程 RSS" :value="toMiB(currentMetric?.processRssBytes)" :precision="2" suffix="MiB" />
        </a-col>
      </a-row>

      <div class="perf-load">
        <a-tag color="purple">Load1: {{ formatLoad(currentMetric?.load1) }}</a-tag>
        <a-tag color="purple">Load5: {{ formatLoad(currentMetric?.load5) }}</a-tag>
        <a-tag color="purple">Load15: {{ formatLoad(currentMetric?.load15) }}</a-tag>
      </div>

      <div class="chart-row">
        <HistoryTimelineChart
          title="CPU / 内存 / 磁盘 / 进程CPU 趋势"
          :series="usageSeries"
          :y-axes="usageYAxes"
        />
      </div>
      <div class="chart-row">
        <HistoryTimelineChart
          title="网络吞吐与进程内存趋势"
          :series="networkSeries"
          :y-axes="networkYAxes"
        />
      </div>
    </a-card>

    <a-card title="虚拟设备调试接口" size="small" class="card">
      <div class="debug-toolbar">
        <a-tag :color="virtualStatus?.running ? 'green' : 'default'">
          {{ virtualStatus?.running ? '运行中' : '已停止' }}
        </a-tag>
        <a-space>
          <a-button
            type="primary"
            size="small"
            :loading="debugLoading"
            :disabled="Boolean(virtualStatus?.running)"
            @click="startVirtual"
          >
            启动虚拟设备
          </a-button>
          <a-button
            danger
            size="small"
            :loading="debugLoading"
            :disabled="!virtualStatus?.running"
            @click="stopVirtual"
          >
            停止并清理
          </a-button>
          <a-button size="small" :loading="pulseLoading" @click="pulseVirtual">单次脉冲压测</a-button>
          <a-button type="dashed" size="small" :loading="stressLoading" @click="stressPulseVirtual">
            压力测试(临时虚拟组)
          </a-button>
        </a-space>
      </div>

      <a-row :gutter="[12, 12]" class="config-row">
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-form-item label="设备数量">
            <a-input-number v-model:value="cfg.count" :min="1" :max="5000" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-form-item label="上报间隔(ms)">
            <a-input-number v-model:value="cfg.intervalMs" :min="10" :max="60000" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-form-item label="SN 前缀">
            <a-input v-model:value="cfg.prefix" placeholder="例如 VDEV" />
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="[12, 12]" class="config-row">
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-form-item label="名称前缀">
            <a-input v-model:value="cfg.namePrefix" placeholder="例如 Virtual Device" />
          </a-form-item>
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-form-item label="起始编号">
            <a-input-number v-model:value="cfg.startIndex" :min="1" :max="100000" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
          <a-form-item label="脉冲轮数">
            <a-input-number v-model:value="cfg.pulseRepeat" :min="1" :max="500" style="width: 100%" />
          </a-form-item>
        </a-col>
      </a-row>

      <a-row :gutter="[12, 12]" class="config-row">
        <a-col :xl="24" :lg="24" :md="24" :sm="24" :xs="24">
          <a-form-item label="行为选项">
            <a-space wrap>
              <a-checkbox v-model:checked="cfg.mutateParam">启用参数扰动</a-checkbox>
              <a-checkbox v-model:checked="cfg.wsBroadcast">推送 WebSocket</a-checkbox>
            </a-space>
          </a-form-item>
        </a-col>
      </a-row>

      <a-divider />

      <a-row :gutter="[12, 12]">
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic title="活跃虚拟设备" :value="virtualStatus?.activeDevices ?? 0" />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic title="累计更新条数" :value="virtualStatus?.updatesGenerated ?? 0" />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic title="更新速率(条/s)" :value="Number((virtualStatus?.updatesPerSecond ?? 0).toFixed(2))" />
        </a-col>
        <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
          <a-statistic title="广播丢弃" :value="virtualStatus?.broadcastDropped ?? 0" />
        </a-col>
      </a-row>

      <a-alert
        class="tip"
        type="info"
        show-icon
        :message="`示例设备: ${virtualStatus?.sampleSn || '-'}，最近心跳: ${formatTime(virtualStatus?.lastTickAt)}`"
      />
      <a-alert
        v-if="lastPulse"
        class="tip"
        type="success"
        show-icon
        :message="`单次脉冲已生成 ${lastPulse.generated} 条，耗时 ${lastPulse.elapsedMs} ms，吞吐 ${lastPulse.updatesPerSecond.toFixed(2)} 条/s`"
      />
      <a-alert
        v-if="lastStressPulse"
        class="tip"
        type="warning"
        show-icon
        :message="`压力测试虚拟组 ${lastStressPulse.group || 'virtual-stress'}：生成 ${lastStressPulse.generated} 条，耗时 ${lastStressPulse.elapsedMs} ms，吞吐 ${lastStressPulse.updatesPerSecond.toFixed(2)} 条/s，测试后已自动清理`"
      />
      <a-alert
        v-if="lastReset"
        class="tip"
        type="success"
        show-icon
        :message="`云系统重置时间 ${formatTime(lastReset.resetAt)}：清空设备 ${lastReset.clearedDevices} 台，历史 ${lastReset.historyRowsDeleted} 条，参数变更 ${lastReset.paramRowsDeleted} 条`"
      />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import * as API from '@/api'
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import ParamViewer from '@/components/device/ParamViewer.vue'
import HistoryTimelineChart from '@/components/device/HistoryTimelineChart.vue'
import type {
  SystemHealth,
  SystemMetricSample,
  SystemMetricsPayload,
  SystemResetResult,
} from '@/api/system'
import type { VirtualDeviceStatus, VirtualPulseResult } from '@/api/debug'
import { getCurrentUser } from '@/utils/auth'

const loading = ref(false)
const debugLoading = ref(false)
const pulseLoading = ref(false)
const stressLoading = ref(false)
const resetLoading = ref(false)

const health = ref<SystemHealth | null>(null)
const metrics = ref<SystemMetricsPayload | null>(null)
const virtualStatus = ref<VirtualDeviceStatus | null>(null)
const lastPulse = ref<VirtualPulseResult | null>(null)
const lastStressPulse = ref<VirtualPulseResult | null>(null)
const lastReset = ref<SystemResetResult | null>(null)

const autoRefreshMs = 5000
const metricPointLimit = 360
let autoTimer: ReturnType<typeof setInterval> | null = null

const isAdmin = computed(() => getCurrentUser()?.role === 'admin')
const currentMetric = computed<SystemMetricSample | null>(() => {
  if (metrics.value?.current) {
    return metrics.value.current
  }
  if (health.value?.performance) {
    return health.value.performance
  }
  return null
})

const hostSummary = computed(() => {
  const host = metrics.value?.host
  if (!host) {
    return '主机信息加载中...'
  }
  const bootText = host.bootTime
    ? new Date(host.bootTime * 1000).toLocaleString()
    : '--'
  const platform = `${host.platform || host.os || '--'} ${host.platformVersion || ''}`.trim()
  return `主机 ${host.hostname || '--'} | 平台 ${platform} | 内核 ${host.kernelVersion || '--'} | 架构 ${host.architecture || '--'} | 启动时间 ${bootText}`
})

const usageSeries = computed(() => {
  const points = metrics.value?.points || []
  return [
    {
      name: 'CPU',
      data: points.map((item) => [item.sampledAt, toNumber(item.cpuPercent)] as [string, number]),
      color: '#1677ff',
      unit: '%',
    },
    {
      name: '内存',
      data: points.map((item) => [item.sampledAt, toNumber(item.memoryUsedPercent)] as [string, number]),
      color: '#fa8c16',
      unit: '%',
    },
    {
      name: '磁盘',
      data: points.map((item) => [item.sampledAt, toNumber(item.diskUsedPercent)] as [string, number]),
      color: '#52c41a',
      unit: '%',
    },
    {
      name: '进程CPU',
      data: points.map((item) => [item.sampledAt, toNumber(item.processCpuPercent)] as [string, number]),
      color: '#722ed1',
      unit: '%',
      dashed: true,
    },
  ]
})

const usageYAxes = computed(() => [
  { name: '使用率(%)', position: 'left' as const, min: 0, max: 100 },
])

const networkSeries = computed(() => {
  const points = metrics.value?.points || []
  return [
    {
      name: '下行',
      data: points.map((item) => [item.sampledAt, toKBps(item.netRecvBps)] as [string, number]),
      color: '#13c2c2',
      unit: 'KB/s',
      yAxisIndex: 0,
    },
    {
      name: '上行',
      data: points.map((item) => [item.sampledAt, toKBps(item.netSentBps)] as [string, number]),
      color: '#eb2f96',
      unit: 'KB/s',
      yAxisIndex: 0,
    },
    {
      name: '进程RSS',
      data: points.map((item) => [item.sampledAt, toMiB(item.processRssBytes)] as [string, number]),
      color: '#2f54eb',
      unit: 'MiB',
      yAxisIndex: 1,
    },
  ]
})

const networkYAxes = computed(() => [
  { name: '网络(KB/s)', position: 'left' as const, min: 0 },
  { name: '进程内存(MiB)', position: 'right' as const, min: 0 },
])

const cfg = reactive({
  count: 200,
  intervalMs: 200,
  prefix: 'VDEV',
  namePrefix: 'Virtual Device',
  startIndex: 1,
  pulseRepeat: 20,
  mutateParam: false,
  wsBroadcast: true,
})

const toNumber = (value: unknown): number => {
  const n = Number(value ?? 0)
  if (!Number.isFinite(n)) {
    return 0
  }
  return Number(n.toFixed(2))
}

const toKBps = (bps: unknown): number => {
  return toNumber(Number(bps ?? 0) / 1024)
}

const toMiB = (bytes: unknown): number => {
  return toNumber(Number(bytes ?? 0) / (1024 * 1024))
}

const formatLoad = (value: unknown): string => {
  const n = Number(value)
  if (!Number.isFinite(n)) {
    return '--'
  }
  return n.toFixed(2)
}

const formatTime = (value?: string): string => {
  if (!value) {
    return '-'
  }
  const t = new Date(value).getTime()
  if (!Number.isFinite(t)) {
    return value
  }
  return new Date(t).toLocaleString()
}

const fetchAll = async () => {
  const [healthRes, virtualRes, metricsRes] = await Promise.all([
    API.system.health(),
    API.debug.virtualStatus(),
    API.system.metrics(metricPointLimit),
  ])
  health.value = healthRes
  virtualStatus.value = virtualRes
  metrics.value = metricsRes
}

const refresh = async (silent = false) => {
  if (!silent) {
    loading.value = true
  }
  try {
    await fetchAll()
  } catch (error: any) {
    if (!silent) {
      message.error(error?.message || '系统状态刷新失败')
    }
  } finally {
    if (!silent) {
      loading.value = false
    }
  }
}

const startVirtual = async () => {
  debugLoading.value = true
  try {
    virtualStatus.value = await API.debug.virtualStart({ ...cfg })
    message.success('虚拟设备已启动')
  } catch (error: any) {
    message.error(error?.message || '启动失败')
  } finally {
    debugLoading.value = false
  }
}

const stopVirtual = async () => {
  debugLoading.value = true
  try {
    virtualStatus.value = await API.debug.virtualStop(true)
    message.success('虚拟设备已停止并清理')
  } catch (error: any) {
    message.error(error?.message || '停止失败')
  } finally {
    debugLoading.value = false
  }
}

const pulseVirtual = async () => {
  pulseLoading.value = true
  try {
    lastPulse.value = await API.debug.virtualPulse({ ...cfg })
    virtualStatus.value = await API.debug.virtualStatus()
    message.success('单次脉冲执行完成')
  } catch (error: any) {
    message.error(error?.message || '脉冲执行失败')
  } finally {
    pulseLoading.value = false
  }
}

const stressPulseVirtual = async () => {
  stressLoading.value = true
  try {
    lastStressPulse.value = await API.debug.virtualStressPulse({ ...cfg })
    virtualStatus.value = await API.debug.virtualStatus()
    message.success('压力测试执行完成，虚拟组已自动清理')
  } catch (error: any) {
    message.error(error?.message || '压力测试执行失败')
  } finally {
    stressLoading.value = false
  }
}

const executeCloudReset = async () => {
  resetLoading.value = true
  try {
    const result = await API.system.resetCloudSystem()
    lastReset.value = result
    lastPulse.value = null
    lastStressPulse.value = null
    await refresh()
    message.success(
      `重置完成：清空设备 ${result.clearedDevices} 台，历史 ${result.historyRowsDeleted} 条，变更日志 ${result.paramRowsDeleted} 条`
    )
  } catch (error: any) {
    message.error(error?.message || '云系统重置失败')
  } finally {
    resetLoading.value = false
  }
}

const confirmCloudReset = () => {
  Modal.confirm({
    title: '确认执行云系统重置吗？',
    content: '将清空当前设备 SN 列表、历史数据和参数变更日志，该操作不可恢复。',
    okText: '立即重置',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      await executeCloudReset()
    },
  })
}

onMounted(() => {
  void refresh()
  autoTimer = setInterval(() => {
    void refresh(true)
  }, autoRefreshMs)
})

onBeforeUnmount(() => {
  if (autoTimer) {
    clearInterval(autoTimer)
    autoTimer = null
  }
})
</script>

<style scoped>
.page {
  padding: 12px;
}

.card {
  margin-top: 12px;
}

.toolbar,
.debug-toolbar,
.perf-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  gap: 10px;
  flex-wrap: wrap;
}

.perf-host {
  margin-bottom: 12px;
}

.perf-load {
  margin: 10px 0 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.chart-row {
  margin-top: 10px;
  border-radius: 10px;
  border: 1px solid #e6eef7;
  background: linear-gradient(180deg, #f9fcff, #f4f9ff);
  padding: 8px;
}

.config-row {
  margin-bottom: 0;
}

.tip {
  margin-top: 10px;
}

.mt-row {
  margin-top: 2px;
}

@media (min-width: 1600px) {
  .page {
    padding: 16px;
  }

  .card {
    margin-top: 14px;
  }

  .chart-row {
    margin-top: 12px;
    padding: 10px;
  }
}

@media (max-width: 900px) {
  .page {
    padding: 8px;
  }
}
</style>
