<template>
  <div class="benchmark-page">
    <a-card size="small" class="toolbar-card">
      <div class="toolbar">
        <div>
          <h3>行业对标扩展</h3>
          <p>参考成熟工业物联网平台能力，新增 10 项运维分析功能并对接本系统数据。</p>
        </div>
        <a-space wrap>
          <a-select v-model:value="windowHours" style="width: 150px" @change="refresh">
            <a-select-option :value="24">最近 24 小时</a-select-option>
            <a-select-option :value="72">最近 72 小时</a-select-option>
            <a-select-option :value="168">最近 7 天</a-select-option>
            <a-select-option :value="240">最近 10 天</a-select-option>
          </a-select>
          <a-button type="primary" :loading="state.loading" @click="refresh">刷新分析</a-button>
          <a-button :loading="state.exporting" @click="downloadReport">导出 CSV 报告</a-button>
        </a-space>
      </div>
    </a-card>

    <a-row :gutter="[12, 12]" class="stat-row">
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="纳入设备" :value="insight?.riskMatrix?.length ?? 0" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="维护候选" :value="insight?.candidates?.length ?? 0" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="固件版本数" :value="insight?.firmware?.length ?? 0" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="生成时间" :value="formatTime(insight?.generatedAt || '')" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[12, 12]" class="section-row">
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" title="SLA 在线率 Top">
          <a-table
            :columns="deviceColumns"
            :data-source="insight?.slaTop || []"
            row-key="deviceSn"
            size="small"
            :pagination="{ pageSize: 6, showSizeChanger: false }"
            :loading="state.loading"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'onlineRatio'">
                {{ toPercent(record.onlineRatio) }}
              </template>
              <template v-else-if="column.dataIndex === 'completeness'">
                {{ toPercent(record.completeness) }}
              </template>
              <template v-else-if="column.dataIndex === 'riskLevel'">
                <a-tag :color="riskColor(record.riskLevel)">{{ riskLabel(record.riskLevel) }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" title="数据完整率 Top">
          <a-table
            :columns="deviceColumns"
            :data-source="insight?.completeTop || []"
            row-key="deviceSn"
            size="small"
            :pagination="{ pageSize: 6, showSizeChanger: false }"
            :loading="state.loading"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'onlineRatio'">
                {{ toPercent(record.onlineRatio) }}
              </template>
              <template v-else-if="column.dataIndex === 'completeness'">
                {{ toPercent(record.completeness) }}
              </template>
              <template v-else-if="column.dataIndex === 'riskLevel'">
                <a-tag :color="riskColor(record.riskLevel)">{{ riskLabel(record.riskLevel) }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[12, 12]" class="section-row">
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" title="重启波动 Top">
          <a-table
            :columns="restartColumns"
            :data-source="insight?.restartTop || []"
            row-key="deviceSn"
            size="small"
            :pagination="{ pageSize: 6, showSizeChanger: false }"
            :loading="state.loading"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'riskLevel'">
                <a-tag :color="riskColor(record.riskLevel)">{{ riskLabel(record.riskLevel) }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" title="参数变更热榜 Top">
          <a-table
            :columns="changeColumns"
            :data-source="insight?.changeTop || []"
            row-key="deviceSn"
            size="small"
            :pagination="{ pageSize: 6, showSizeChanger: false }"
            :loading="state.loading"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'riskLevel'">
                <a-tag :color="riskColor(record.riskLevel)">{{ riskLabel(record.riskLevel) }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[12, 12]" class="section-row">
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" title="固件版本分布">
          <v-chart class="chart" :option="firmwareOption" autoresize />
        </a-card>
      </a-col>
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" title="离线时段热力（24h）">
          <v-chart class="chart" :option="offlineOption" autoresize />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[12, 12]" class="section-row">
      <a-col :span="24">
        <a-card size="small" title="风险矩阵（在线率 vs 完整率）">
          <v-chart class="chart risk-chart" :option="riskMatrixOption" autoresize />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[12, 12]" class="section-row">
      <a-col :xl="8" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" title="运维审计汇总（密钥工具）">
          <a-table
            :columns="auditColumns"
            :data-source="insight?.audit || []"
            row-key="operation"
            size="small"
            :pagination="false"
            :loading="state.loading"
          />
        </a-card>
      </a-col>
      <a-col :xl="16" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" title="维护候选清单">
          <a-table
            :columns="candidateColumns"
            :data-source="insight?.candidates || []"
            row-key="deviceSn"
            size="small"
            :pagination="{ pageSize: 8, showSizeChanger: false }"
            :loading="state.loading"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'onlineRatio'">
                {{ toPercent(record.onlineRatio) }}
              </template>
              <template v-else-if="column.dataIndex === 'completeness'">
                {{ toPercent(record.completeness) }}
              </template>
              <template v-else-if="column.dataIndex === 'priority'">
                <a-tag :color="priorityColor(record.priority)">{{ priorityLabel(record.priority) }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import * as API from '@/api'
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, ScatterChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'
import VChart from 'vue-echarts'
import type { BenchmarkInsightPayload } from '@/api/benchmark'

use([
  CanvasRenderer,
  BarChart,
  ScatterChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
])

const windowHours = ref(24)
const insight = ref<BenchmarkInsightPayload | null>(null)

const state = reactive({
  loading: false,
  exporting: false,
})

const deviceColumns = [
  { title: '设备', dataIndex: 'deviceName', width: 170 },
  { title: 'SN', dataIndex: 'deviceSn', width: 150 },
  { title: '在线率', dataIndex: 'onlineRatio', width: 90 },
  { title: '完整率', dataIndex: 'completeness', width: 90 },
  { title: '风险', dataIndex: 'riskLevel', width: 90 },
]

const restartColumns = [
  { title: '设备', dataIndex: 'deviceName', width: 160 },
  { title: '重启次数', dataIndex: 'restartCount', width: 90 },
  { title: '风险分', dataIndex: 'riskScore', width: 90 },
  { title: '风险等级', dataIndex: 'riskLevel', width: 100 },
]

const changeColumns = [
  { title: '设备', dataIndex: 'deviceName', width: 160 },
  { title: '变更次数', dataIndex: 'paramChangeCount', width: 100 },
  { title: '风险分', dataIndex: 'riskScore', width: 90 },
  { title: '风险等级', dataIndex: 'riskLevel', width: 100 },
]

const auditColumns = [
  { title: '操作', dataIndex: 'operation', width: 120 },
  { title: '成功', dataIndex: 'success', width: 100 },
  { title: '失败', dataIndex: 'failed', width: 100 },
]

const candidateColumns = [
  { title: '设备', dataIndex: 'deviceName', width: 160 },
  { title: 'SN', dataIndex: 'deviceSn', width: 150 },
  { title: '在线率', dataIndex: 'onlineRatio', width: 90 },
  { title: '完整率', dataIndex: 'completeness', width: 90 },
  { title: '重启', dataIndex: 'restartCount', width: 70 },
  { title: '变更', dataIndex: 'paramChangeCnt', width: 70 },
  { title: '风险分', dataIndex: 'riskScore', width: 90 },
  { title: '优先级', dataIndex: 'priority', width: 100 },
  { title: '建议', dataIndex: 'reason' },
]

const riskColor = (risk: string): string => {
  if (risk === 'low') {
    return 'green'
  }
  if (risk === 'medium') {
    return 'gold'
  }
  if (risk === 'high') {
    return 'orange'
  }
  return 'red'
}

const riskLabel = (risk: string): string => {
  if (risk === 'low') {
    return '低'
  }
  if (risk === 'medium') {
    return '中'
  }
  if (risk === 'high') {
    return '高'
  }
  return '严重'
}

const priorityColor = (priority: string): string => {
  if (priority === 'high') {
    return 'red'
  }
  if (priority === 'medium') {
    return 'orange'
  }
  return 'blue'
}

const priorityLabel = (priority: string): string => {
  if (priority === 'high') {
    return '高优先'
  }
  if (priority === 'medium') {
    return '中优先'
  }
  return '低优先'
}

const toPercent = (value: number): string => {
  return `${(Number(value || 0) * 100).toFixed(2)}%`
}

const formatTime = (value: string): string => {
  if (!value) {
    return '--'
  }
  const t = new Date(value).getTime()
  if (!Number.isFinite(t)) {
    return value
  }
  return new Date(t).toLocaleString()
}

const firmwareOption = computed<EChartsOption>(() => {
  const items = insight.value?.firmware || []
  return {
    grid: { left: 48, right: 20, top: 24, bottom: 30 },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: items.map((item) => `HW-${item.hardwareBate}`),
      axisLabel: { color: '#6b7280' },
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#6b7280' },
      splitLine: { lineStyle: { color: '#edf2f7' } },
    },
    series: [
      {
        type: 'bar',
        data: items.map((item) => item.count),
        itemStyle: { color: '#3b82f6' },
        barWidth: 26,
      },
    ],
  }
})

const offlineOption = computed<EChartsOption>(() => {
  const items = insight.value?.offlineHeat || []
  return {
    grid: { left: 48, right: 20, top: 24, bottom: 30 },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: items.map((item) => `${String(item.hour).padStart(2, '0')}:00`),
      axisLabel: { color: '#6b7280' },
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#6b7280',
        formatter: '{value}%',
      },
      splitLine: { lineStyle: { color: '#edf2f7' } },
    },
    series: [
      {
        name: '离线率',
        type: 'bar',
        data: items.map((item) => Number((item.offlineRate * 100).toFixed(2))),
        itemStyle: { color: '#f59e0b' },
        barWidth: 18,
      },
    ],
  }
})

const riskMatrixOption = computed<EChartsOption>(() => {
  const points = insight.value?.riskMatrix || []
  return {
    grid: { left: 56, right: 26, top: 20, bottom: 44 },
    tooltip: {
      trigger: 'item',
      formatter: (raw: any) => {
        const value = raw?.value || []
        return [
          `${raw?.data?.name || '-'}`,
          `在线率: ${Number(value[0] || 0).toFixed(2)}%`,
          `完整率: ${Number(value[1] || 0).toFixed(2)}%`,
          `风险分: ${Number(value[2] || 0)}`,
        ].join('<br/>')
      },
    },
    xAxis: {
      name: '在线率(%)',
      nameLocation: 'middle',
      nameGap: 28,
      min: 0,
      max: 100,
      axisLabel: { color: '#6b7280' },
      splitLine: { lineStyle: { color: '#edf2f7' } },
    },
    yAxis: {
      name: '完整率(%)',
      nameLocation: 'middle',
      nameGap: 36,
      min: 0,
      max: 100,
      axisLabel: { color: '#6b7280' },
      splitLine: { lineStyle: { color: '#edf2f7' } },
    },
    series: [
      {
        type: 'scatter',
        symbolSize: (value: number[]) => {
          const risk = Number(value[2] || 0)
          return 10 + Math.min(24, risk * 0.4)
        },
        data: points.map((item) => ({
          name: `${item.deviceName} (${item.deviceSn})`,
          value: [item.onlineRatioPct, item.completenessPct, item.riskScore],
          itemStyle: {
            color:
              item.riskLevel === 'critical'
                ? '#ef4444'
                : item.riskLevel === 'high'
                  ? '#f97316'
                  : item.riskLevel === 'medium'
                    ? '#f59e0b'
                    : '#10b981',
          },
        })),
      },
    ],
  }
})

const refresh = async () => {
  state.loading = true
  try {
    insight.value = await API.benchmark.insights(windowHours.value)
  } catch (error: any) {
    message.error(error?.message || '行业对标分析加载失败')
  } finally {
    state.loading = false
  }
}

const downloadReport = async () => {
  state.exporting = true
  try {
    const blob = await API.benchmark.exportReport(windowHours.value)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    const stamp = new Date().toISOString().replace(/[:.]/g, '-')
    a.href = url
    a.download = `行业对标报告_${stamp}.csv`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    message.success('报告导出成功')
  } catch (error: any) {
    message.error(error?.message || '报告导出失败')
  } finally {
    state.exporting = false
  }
}

onMounted(() => {
  void refresh()
})
</script>

<style scoped>
.benchmark-page {
  padding: 12px;
}

.toolbar-card {
  margin-bottom: 12px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.toolbar h3 {
  margin: 0;
  color: #294861;
  font-size: 16px;
}

.toolbar p {
  margin: 4px 0 0;
  color: #6d8094;
  font-size: 12px;
}

.stat-row,
.section-row {
  margin-top: 0;
}

.chart {
  height: 280px;
  width: 100%;
}

.risk-chart {
  height: 360px;
}

@media (max-width: 900px) {
  .benchmark-page {
    padding: 8px;
  }

  .chart {
    height: 240px;
  }
}
</style>
