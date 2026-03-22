<template>
  <div class="device-detail">
    <div class="detail-header">
      <div class="header-left">
        <a-tag :color="isOnline ? 'green' : 'default'">{{ isOnline ? '在线' : '离线' }}</a-tag>
        <span class="sn">{{ info.Sn || '-' }}</span>
        <a-tag v-if="groupText" color="blue">{{ groupText }}</a-tag>
      </div>
      <div class="header-right">
        <span class="name">{{ info.Name || '未命名设备' }}</span>
        <span class="name" v-if="modelText">型号 {{ modelText }}</span>
        <span class="name" v-if="pnText">PN {{ pnText }}</span>
        <span class="last-seen">{{ lastSeenText }}</span>
      </div>
    </div>

    <a-row :gutter="[10, 10]" class="summary-row">
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="硬件版本" :value="hardwareVersion" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="运行时长" :value="uptimeText" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="泵通道" :value="pumpCount" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="温度通道" :value="tempCount" />
        </a-card>
      </a-col>
    </a-row>

    <a-tabs v-model:activeKey="activeTab" class="detail-tabs" size="small" :destroyInactiveTabPane="true">
      <a-tab-pane key="insight" tab="结构化详情">
        <DeviceInsight :info="info" />
      </a-tab-pane>

      <a-tab-pane key="packet" tab="原始报文">
        <ParamViewer :node="info.Packet" />
      </a-tab-pane>

      <a-tab-pane key="history" tab="历史复盘">
        <div class="toolbar-wrap">
          <a-select v-model:value="historyRangePreset" style="width: 140px" @change="loadHistory">
            <a-select-option value="1h">最近 1 小时</a-select-option>
            <a-select-option value="6h">最近 6 小时</a-select-option>
            <a-select-option value="24h">最近 24 小时</a-select-option>
            <a-select-option value="3d">最近 3 天</a-select-option>
            <a-select-option value="7d">最近 7 天</a-select-option>
            <a-select-option value="10d">最近 10 天</a-select-option>
          </a-select>
          <span class="metric-label">温度CH</span>
          <a-input-number
            v-model:value="metricSelection.tempIndex"
            :min="0"
            :max="255"
            :step="1"
            style="width: 84px"
          />
          <span class="metric-label">电压CH</span>
          <a-input-number
            v-model:value="metricSelection.voltageIndex"
            :min="0"
            :max="255"
            :step="1"
            style="width: 84px"
          />
          <span class="metric-label">泵浦电流CH</span>
          <a-input-number
            v-model:value="metricSelection.currentIndex"
            :min="0"
            :max="255"
            :step="1"
            style="width: 84px"
          />
          <a-button :loading="historyLoading" @click="loadHistory">刷新历史</a-button>
          <a-button :loading="historyExporting" @click="exportTimelineJson">导出趋势 JSON</a-button>
          <a-tag color="blue">采样周期: 1 分钟</a-tag>
          <a-tag color="default">记录条数: {{ historyTotal }}</a-tag>
          <a-tag color="cyan">图表点数: {{ metricTimelineItems.length }}</a-tag>
        </div>

        <a-spin :spinning="historyLoading">
          <div v-if="metricTimelineItems.length" class="timeline-grid">
            <HistoryTimelineChart
              :title="linkedTimelineTitle"
              :series="linkedTimelineSeries"
              :y-axes="linkedTimelineYAxes"
            />
          </div>
          <a-empty v-else :description="'所选通道暂无可视化趋势数据'" />
        </a-spin>

        <a-table
          :columns="metricColumns"
          :data-source="metricTimelineItems"
          :loading="historyLoading"
          size="small"
          row-key="sampledAt"
          :pagination="false"
          :scroll="{ y: 360 }"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'sampledAt'">
              {{ formatTime(record.sampledAt) }}
            </template>
            <template v-else-if="column.dataIndex === 'online'">
              <a-tag :color="record.online ? 'green' : 'default'">{{ record.online ? '在线' : '离线' }}</a-tag>
            </template>
            <template v-else-if="column.dataIndex === 'temp'">
              {{ formatMetricValue(record.temp, '°C') }}
            </template>
            <template v-else-if="column.dataIndex === 'voltage'">
              {{ formatMetricValue(record.voltage, 'raw') }}
            </template>
            <template v-else-if="column.dataIndex === 'current'">
              {{ formatMetricValue(record.current, 'raw') }}
            </template>
          </template>
        </a-table>
      </a-tab-pane>

      <a-tab-pane key="changes" tab="参数变更日志">
        <div class="toolbar-wrap">
          <a-select v-model:value="changeRangePreset" style="width: 140px" @change="loadChangeLogs">
            <a-select-option value="1h">最近 1 小时</a-select-option>
            <a-select-option value="6h">最近 6 小时</a-select-option>
            <a-select-option value="24h">最近 24 小时</a-select-option>
            <a-select-option value="3d">最近 3 天</a-select-option>
            <a-select-option value="7d">最近 7 天</a-select-option>
            <a-select-option value="10d">最近 10 天</a-select-option>
          </a-select>
          <a-input-search
            v-model:value="changePathKeyword"
            class="path-input"
            allow-clear
            placeholder="按参数路径过滤，如 Pump_module"
            @search="loadChangeLogs"
          />
          <a-button :loading="changeLoading" @click="loadChangeLogs">刷新日志</a-button>
          <a-button type="primary" :loading="exporting" @click="exportChangeLogs">导出 CSV</a-button>
          <a-tag color="default">变更条数: {{ changeTotal }}</a-tag>
        </div>

        <a-table
          :columns="changeColumns"
          :data-source="changeItems"
          :loading="changeLoading"
          size="small"
          row-key="id"
          :pagination="false"
          :scroll="{ y: 360 }"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'changedAt'">
              {{ formatTime(record.changedAt) }}
            </template>
            <template v-else-if="column.dataIndex === 'oldValue'">
              <span class="value-text">{{ record.oldValue }}</span>
            </template>
            <template v-else-if="column.dataIndex === 'newValue'">
              <span class="value-text">{{ record.newValue }}</span>
            </template>
          </template>
        </a-table>
      </a-tab-pane>

      <a-tab-pane key="base" tab="基础信息">
        <a-descriptions :column="1" size="small" bordered>
          <a-descriptions-item label="SN">{{ info.Sn || '-' }}</a-descriptions-item>
          <a-descriptions-item label="名称">{{ info.Name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="分组">{{ groupText || '--' }}</a-descriptions-item>
          <a-descriptions-item label="型号">{{ modelText || '--' }}</a-descriptions-item>
          <a-descriptions-item label="PN">{{ pnText || '--' }}</a-descriptions-item>
          <a-descriptions-item label="状态">{{ isOnline ? '在线' : '离线' }}</a-descriptions-item>
          <a-descriptions-item label="最后更新时间">{{ info.Last_rx_time || '-' }}</a-descriptions-item>
          <a-descriptions-item label="硬件版本">{{ hardwareVersion }}</a-descriptions-item>
          <a-descriptions-item label="运行时长">{{ uptimeText }}</a-descriptions-item>
          <a-descriptions-item label="泵通道数">{{ pumpCount }}</a-descriptions-item>
          <a-descriptions-item label="温度通道数">{{ tempCount }}</a-descriptions-item>
        </a-descriptions>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, toRefs, watch } from 'vue'
import moment from 'moment'
import { message } from 'ant-design-vue'
import * as API from '@/api'
import type { Device } from '@/types/api'
import type { DeviceHistoryMetricPoint, DeviceParamChange } from '@/api/device'
import ParamViewer from './ParamViewer.vue'
import DeviceInsight from './DeviceInsight.vue'
import HistoryTimelineChart from './HistoryTimelineChart.vue'

const props = defineProps<{ info: Device }>()
const { info } = toRefs(props)

type RangePreset = '1h' | '6h' | '24h' | '3d' | '7d' | '10d'

const inputReg = computed<any>(() => info.value?.Packet?.Femto_input_reg ?? {})
const holdingReg = computed<any>(() => info.value?.Packet?.Femto_holding_reg ?? {})
const mon = computed<any>(() => inputReg.value?.Mon ?? {})

const activeTab = ref<'insight' | 'packet' | 'history' | 'changes' | 'base'>('insight')

const historyRangePreset = ref<RangePreset>('24h')
const historyLoading = ref(false)
const historyExporting = ref(false)
const historyTotal = ref(0)
const metricTimelineItems = ref<DeviceHistoryMetricPoint[]>([])
const metricSelection = ref({
  tempIndex: 0,
  voltageIndex: 0,
  currentIndex: 0,
})

const changeRangePreset = ref<RangePreset>('24h')
const changePathKeyword = ref('')
const changeLoading = ref(false)
const exporting = ref(false)
const changeItems = ref<DeviceParamChange[]>([])
const changeTotal = ref(0)

const changeColumns = [
  { title: '变更时间', dataIndex: 'changedAt', width: 170 },
  { title: '参数路径', dataIndex: 'paramPath', width: 320 },
  { title: '旧值', dataIndex: 'oldValue', width: 280 },
  { title: '新值', dataIndex: 'newValue', width: 280 },
]

const metricColumns = computed(() => [
  { title: '采样时间', dataIndex: 'sampledAt', width: 170 },
  { title: '在线状态', dataIndex: 'online', width: 90 },
  { title: `温度 CH${metricSelection.value.tempIndex}(°C)`, dataIndex: 'temp', width: 170 },
  { title: `电压 CH${metricSelection.value.voltageIndex}(raw)`, dataIndex: 'voltage', width: 180 },
  { title: `泵浦电流 CH${metricSelection.value.currentIndex}(raw)`, dataIndex: 'current', width: 210 },
])

const linkedTimelineYAxes = computed(() => [
  { name: '温度(°C)', position: 'left' as const, min: 'auto' as const },
  { name: '电压(raw)', position: 'right' as const, min: 0 },
  { name: '电流(raw)', position: 'right' as const, offset: 58, min: 0 },
])

const linkedTimelineTitle = computed(
  () =>
    `温度 CH${metricSelection.value.tempIndex} / 电压 CH${metricSelection.value.voltageIndex} / 泵浦电流 CH${metricSelection.value.currentIndex} 联动时间线`
)

const linkedTimelineSeries = computed(() => {
  const temp = metricTimelineItems.value.map((item) => [item.sampledAt, toMetricValue(item.temp)] as [string, number])
  const voltage = metricTimelineItems.value.map((item) => [item.sampledAt, toMetricValue(item.voltage)] as [string, number])
  const current = metricTimelineItems.value.map((item) => [item.sampledAt, toMetricValue(item.current)] as [string, number])
  return [
    { name: '温度', data: temp, color: '#fa8c16', yAxisIndex: 0, unit: '°C' },
    { name: '电压', data: voltage, color: '#13c2c2', yAxisIndex: 1, unit: 'raw' },
    { name: '电流', data: current, color: '#722ed1', yAxisIndex: 2, unit: 'raw' },
  ]
})

const isOnline = computed<boolean>(() => {
  const ts = info.value?.Last_rx_time
  if (!ts) {
    return false
  }
  const t = new Date(ts).getTime()
  if (!Number.isFinite(t)) {
    return false
  }
  return Date.now() - t <= 30_000
})

const lastSeenText = computed<string>(() => {
  const ts = info.value?.Last_rx_time
  if (!ts) {
    return '--'
  }
  const m = moment(ts)
  if (!m.isValid()) {
    return ts
  }
  return `更新于 ${m.fromNow()}`
})

const hardwareVersion = computed(() => Number(inputReg.value?.Bate?.Hardware_bate ?? 0))

const groupText = computed<string>(() => {
  const group = String((info.value as any)?.Group || '').trim()
  if (!group) {
    return ''
  }
  if (group === 'virtual-stress') {
    return '压力测试组'
  }
  return group
})

const decodeAsciiBytes = (value: unknown): string => {
  if (!Array.isArray(value)) {
    return ''
  }
  const chars: string[] = []
  for (const item of value) {
    const code = Number(item ?? 0) & 0xff
    if (code === 0) {
      break
    }
    if (code >= 32 && code <= 126) {
      chars.push(String.fromCharCode(code))
    }
  }
  return chars.join('')
}

const modelText = computed<string>(() =>
  decodeAsciiBytes(holdingReg.value?.Laser_para?.Laser_info?.Model)
)

const pnText = computed<string>(() =>
  decodeAsciiBytes(holdingReg.value?.Laser_para?.Laser_info?.PN)
)

const pumpCount = computed<number>(() => {
  const pumps = mon.value?.Pump_mon ?? mon.value?.Femto_input_reg_monitor_pump
  return Array.isArray(pumps) ? pumps.length : 0
})

const tempCount = computed<number>(() => {
  const temp = mon.value?.Temp
  return Array.isArray(temp) ? temp.length : 0
})

const uptimeText = computed<string>(() => {
  const uptime = inputReg.value?.Time?.Uptime
  if (!Array.isArray(uptime) || uptime.length < 2) {
    return '--'
  }
  const high = Number(uptime[0] ?? 0) & 0xffff
  const low = Number(uptime[1] ?? 0) & 0xffff
  const merged = ((high << 16) >>> 0) + low
  return formatSeconds(merged)
})

const getRangeByPreset = (preset: RangePreset) => {
  const end = moment()
  const start = end.clone()
  if (preset === '1h') {
    start.subtract(1, 'hour')
  } else if (preset === '6h') {
    start.subtract(6, 'hour')
  } else if (preset === '24h') {
    start.subtract(24, 'hour')
  } else if (preset === '3d') {
    start.subtract(3, 'day')
  } else if (preset === '7d') {
    start.subtract(7, 'day')
  } else {
    start.subtract(10, 'day')
  }
  return { start, end }
}

const getTimelineLimitByPreset = (preset: RangePreset): number => {
  if (preset === '1h') {
    return 180
  }
  if (preset === '6h') {
    return 800
  }
  if (preset === '24h') {
    return 2200
  }
  if (preset === '3d') {
    return 5200
  }
  if (preset === '7d') {
    return 12500
  }
  return 18000
}

const loadHistory = async () => {
  const sn = info.value?.Sn
  if (!sn) {
    historyTotal.value = 0
    metricTimelineItems.value = []
    return
  }

  historyLoading.value = true
  try {
    const range = getRangeByPreset(historyRangePreset.value)
    const timelineLimit = getTimelineLimitByPreset(historyRangePreset.value)
    const metricsRes = await API.device.historyMetrics(sn, {
      start: range.start.toISOString(),
      end: range.end.toISOString(),
      limit: timelineLimit,
      offset: 0,
      maxPoints: 5000,
      tempIndex: metricSelection.value.tempIndex,
      voltageIndex: metricSelection.value.voltageIndex,
      currentIndex: metricSelection.value.currentIndex,
    })

    historyTotal.value = Number(metricsRes.total || 0)
    metricTimelineItems.value = (metricsRes.items || [])
      .slice()
      .sort((a, b) => new Date(a.sampledAt).getTime() - new Date(b.sampledAt).getTime())
  } catch (error: any) {
    historyTotal.value = 0
    metricTimelineItems.value = []
    message.error(error?.message || '加载历史数据失败')
  } finally {
    historyLoading.value = false
  }
}

const loadChangeLogs = async () => {
  const sn = info.value?.Sn
  if (!sn) {
    changeItems.value = []
    changeTotal.value = 0
    return
  }

  changeLoading.value = true
  try {
    const range = getRangeByPreset(changeRangePreset.value)
    const res = await API.device.changeLogs(sn, {
      start: range.start.toISOString(),
      end: range.end.toISOString(),
      path: changePathKeyword.value.trim() || undefined,
      offset: 0,
      limit: 600,
    })
    changeItems.value = res.items || []
    changeTotal.value = Number(res.total || 0)
  } catch (error: any) {
    changeItems.value = []
    changeTotal.value = 0
    message.error(error?.message || '加载参数变更日志失败')
  } finally {
    changeLoading.value = false
  }
}

const exportChangeLogs = async () => {
  const sn = info.value?.Sn
  if (!sn) {
    message.warning('设备 SN 不存在，无法导出')
    return
  }

  exporting.value = true
  try {
    const range = getRangeByPreset(changeRangePreset.value)
    const blob = await API.device.exportChangeLogs(sn, {
      start: range.start.toISOString(),
      end: range.end.toISOString(),
      path: changePathKeyword.value.trim() || undefined,
      limit: 5000,
    })

    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `device_changes_${sn.replace(/\s+/g, '_')}_${moment().format('YYYYMMDD_HHmmss')}.csv`
    document.body.appendChild(a)
    a.click()
    a.remove()
    window.URL.revokeObjectURL(url)
    message.success('变更日志已导出')
  } catch (error: any) {
    message.error(error?.message || '导出失败')
  } finally {
    exporting.value = false
  }
}

const exportTimelineJson = async () => {
  const sn = info.value?.Sn
  if (!sn) {
    message.warning('设备 SN 不存在，无法导出')
    return
  }
  if (!metricTimelineItems.value.length) {
    message.warning('当前没有可导出的趋势数据')
    return
  }

  historyExporting.value = true
  try {
    const payload = {
      deviceSn: sn,
      range: historyRangePreset.value,
      metricSelection: {
        tempIndex: metricSelection.value.tempIndex,
        voltageIndex: metricSelection.value.voltageIndex,
        currentIndex: metricSelection.value.currentIndex,
      },
      exportedAt: moment().toISOString(),
      points: metricTimelineItems.value,
    }
    const blob = new Blob([JSON.stringify(payload, null, 2)], {
      type: 'application/json;charset=utf-8',
    })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `device_timeline_${sn.replace(/\s+/g, '_')}_${moment().format('YYYYMMDD_HHmmss')}.json`
    document.body.appendChild(a)
    a.click()
    a.remove()
    window.URL.revokeObjectURL(url)
    message.success('趋势数据已导出，可直接用于 ECharts 长时绘图')
  } finally {
    historyExporting.value = false
  }
}

const formatTime = (value: string): string => {
  if (!value) {
    return '-'
  }
  const m = moment(value)
  return m.isValid() ? m.format('YYYY-MM-DD HH:mm:ss') : value
}

const toMetricValue = (value: number | null | undefined): number => {
  if (value === null || value === undefined) {
    return Number.NaN
  }
  const n = Number(value)
  if (!Number.isFinite(n)) {
    return Number.NaN
  }
  return Number(n.toFixed(2))
}

const formatMetricValue = (value: number | null | undefined, unit: string): string => {
  const n = toMetricValue(value)
  if (!Number.isFinite(n)) {
    return '--'
  }
  return `${n.toFixed(2)} ${unit}`
}

const formatSeconds = (seconds: number): string => {
  if (!Number.isFinite(seconds) || seconds <= 0) {
    return '0秒'
  }
  const s = Math.floor(seconds)
  const day = Math.floor(s / 86400)
  const hour = Math.floor((s % 86400) / 3600)
  const minute = Math.floor((s % 3600) / 60)
  const sec = s % 60
  if (day > 0) {
    return `${day}天 ${hour}时 ${minute}分`
  }
  if (hour > 0) {
    return `${hour}时 ${minute}分 ${sec}秒`
  }
  if (minute > 0) {
    return `${minute}分 ${sec}秒`
  }
  return `${sec}秒`
}

watch(
  () => info.value?.Sn,
  () => {
    historyTotal.value = 0
    metricTimelineItems.value = []
    changeItems.value = []
    changeTotal.value = 0
    if (activeTab.value === 'history') {
      void loadHistory()
    }
    if (activeTab.value === 'changes') {
      void loadChangeLogs()
    }
  }
)

watch(activeTab, (tab) => {
  if (tab === 'history') {
    void loadHistory()
  }
  if (tab === 'changes') {
    void loadChangeLogs()
  }
})

watch(
  metricSelection,
  () => {
    if (activeTab.value === 'history') {
      void loadHistory()
    }
  },
  { deep: true }
)
</script>

<style scoped>
.device-detail {
  border-radius: 10px;
  border: 1px solid #e7edf4;
  background: #fff;
  padding: 12px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 10px;
  border: 1px solid #edf3fb;
  background: linear-gradient(180deg, #f8fbff, #ffffff);
  border-radius: 8px;
  padding: 8px 10px;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.sn {
  color: #1677ff;
  font-weight: 700;
}

.name {
  color: #4d5f73;
  font-size: 13px;
}

.last-seen {
  color: #7f93a7;
  font-size: 12px;
}

.summary-row {
  margin-bottom: 10px;
}

.detail-tabs :deep(.ant-tabs-content-holder) {
  padding-top: 6px;
}

.toolbar-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 10px;
}

.metric-label {
  color: #607387;
  font-size: 12px;
}

.path-input {
  width: min(320px, 72vw);
}

.timeline-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
  margin-bottom: 12px;
}

.value-text {
  display: block;
  max-width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-family: Consolas, 'Courier New', monospace;
}

@media (max-width: 900px) {
  .device-detail {
    padding: 8px;
  }

  .detail-header {
    padding: 6px 8px;
  }
}
</style>
