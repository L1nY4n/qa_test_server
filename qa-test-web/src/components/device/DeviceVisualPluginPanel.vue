<template>
  <div class="panel-grid">
    <v-chart class="viz-card" :option="healthGaugeOption" autoresize />
    <v-chart class="viz-card" :option="pumpTimelineOption" autoresize />
    <v-chart class="viz-card" :option="alarmPieOption" autoresize />
    <v-chart class="viz-card" :option="heatmapOption" autoresize />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { GaugeChart, LineChart, PieChart, HeatmapChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
  VisualMapComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'
import VChart from 'vue-echarts'

use([
  CanvasRenderer,
  GaugeChart,
  LineChart,
  PieChart,
  HeatmapChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
  VisualMapComponent,
])

interface PumpRowLike {
  index: number
  actual: number
}

const props = withDefaults(
  defineProps<{
    healthScore: number
    pumpRows: PumpRowLike[]
    alarmNowCount: number
    alarmHistoryCount: number
    tempValues: number[]
    voltageValues: number[]
  }>(),
  {
    healthScore: 0,
    pumpRows: () => [],
    alarmNowCount: 0,
    alarmHistoryCount: 0,
    tempValues: () => [],
    voltageValues: () => [],
  }
)

const safeHealth = computed(() => {
  const n = Number(props.healthScore ?? 0)
  if (!Number.isFinite(n)) {
    return 0
  }
  return Math.max(0, Math.min(100, Number(n.toFixed(1))))
})

const healthGaugeOption = computed<EChartsOption>(() => {
  const score = safeHealth.value
  const color: [number, string][] = [
    [0.3, '#f5222d'],
    [0.7, '#faad14'],
    [1, '#52c41a'],
  ]
  return {
    title: {
      text: '健康度仪表盘',
      left: 'center',
      top: 6,
      textStyle: { fontSize: 13, fontWeight: 700 },
    },
    series: [
      {
        type: 'gauge',
        center: ['50%', '58%'],
        radius: '88%',
        min: 0,
        max: 100,
        splitNumber: 10,
        axisLine: {
          lineStyle: {
            width: 12,
            color,
          },
        },
        pointer: {
          length: '58%',
          width: 4,
        },
        axisTick: { show: false },
        splitLine: {
          length: 10,
        },
        axisLabel: {
          distance: -22,
          color: '#617a93',
          fontSize: 10,
        },
        progress: { show: true, width: 12 },
        detail: {
          valueAnimation: true,
          formatter: '{value} 分',
          color: '#20364d',
          fontSize: 18,
          offsetCenter: [0, '56%'],
        },
        data: [{ value: score }],
      },
    ],
  }
})

const MAX_TIMELINE_POINTS = 120
const MAX_CHANNELS = 4
const pumpTimeline = ref<Array<{ at: number; values: number[] }>>([])
const lastPumpSnapshotKey = ref('')

watch(
  () => props.pumpRows.map((item) => `${item.index}:${Number(item.actual ?? 0).toFixed(2)}`).join('|'),
  (snapshotKey) => {
    if (!snapshotKey || snapshotKey === lastPumpSnapshotKey.value) {
      return
    }
    lastPumpSnapshotKey.value = snapshotKey

    const rows = (props.pumpRows || []).slice().sort((a, b) => a.index - b.index)
    const values = rows.map((item) => Number(item.actual ?? 0))
    pumpTimeline.value.push({
      at: Date.now(),
      values,
    })
    if (pumpTimeline.value.length > MAX_TIMELINE_POINTS) {
      pumpTimeline.value.splice(0, pumpTimeline.value.length - MAX_TIMELINE_POINTS)
    }
  },
  { immediate: true }
)

const pumpTimelineOption = computed<EChartsOption>(() => {
  const timeline = pumpTimeline.value
  const channelCount = Math.min(
    MAX_CHANNELS,
    timeline.reduce((max, item) => Math.max(max, item.values.length), 0)
  )
  const series = Array.from({ length: channelCount }, (_, index) => ({
    name: `CH${index}`,
    type: 'line' as const,
    smooth: true,
    showSymbol: false,
    lineStyle: { width: 2 },
    data: timeline.map((item) => [item.at, Number(item.values[index] ?? 0)] as [number, number]),
  }))

  return {
    title: {
      text: '泵浦实际电流时间关系',
      left: 10,
      top: 6,
      textStyle: { fontSize: 13, fontWeight: 700 },
    },
    tooltip: { trigger: 'axis' },
    legend: { top: 26, right: 8 },
    grid: { left: 42, right: 24, top: 56, bottom: 24 },
    xAxis: {
      type: 'time',
      axisLabel: { color: '#6c8196' },
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#6c8196' },
      splitLine: { lineStyle: { color: '#ebf1f8' } },
    },
    series,
  }
})

const alarmPieOption = computed<EChartsOption>(() => {
  const now = Math.max(0, Number(props.alarmNowCount || 0))
  const history = Math.max(0, Number(props.alarmHistoryCount || 0))
  const clearWeight = Math.max(1, 120 - now - history)

  return {
    title: {
      text: '告警分布插件',
      left: 'center',
      top: 6,
      textStyle: { fontSize: 13, fontWeight: 700 },
    },
    tooltip: { trigger: 'item' },
    legend: { bottom: 2, left: 'center' },
    series: [
      {
        type: 'pie',
        radius: ['44%', '72%'],
        center: ['50%', '55%'],
        data: [
          { name: '当前告警', value: now, itemStyle: { color: '#f5222d' } },
          { name: '历史告警', value: history, itemStyle: { color: '#faad14' } },
          { name: '稳定区间', value: clearWeight, itemStyle: { color: '#52c41a' } },
        ],
        label: {
          formatter: '{b}: {c}',
          fontSize: 11,
        },
      },
    ],
  }
})

const heatmapOption = computed<EChartsOption>(() => {
  const temp = (props.tempValues || []).map((item) => Number(item ?? 0))
  const voltage = (props.voltageValues || []).map((item) => Number(item ?? 0))
  const channelCount = Math.max(temp.length, voltage.length)
  const channels = Array.from({ length: channelCount }, (_, i) => `CH${i}`)

  const rows = ['温度', '电压']
  const heatData: Array<[number, number, number]> = []

  for (let i = 0; i < channelCount; i += 1) {
    heatData.push([i, 0, Number.isFinite(temp[i]) ? temp[i] : 0])
    heatData.push([i, 1, Number.isFinite(voltage[i]) ? voltage[i] : 0])
  }

  const values = heatData.map((item) => item[2])
  const maxValue = values.length ? Math.max(...values) : 1

  return {
    title: {
      text: '通道热力图插件',
      left: 10,
      top: 6,
      textStyle: { fontSize: 13, fontWeight: 700 },
    },
    tooltip: {
      position: 'top',
      formatter: (raw: any) => {
        const value = raw?.data?.[2] ?? 0
        const x = raw?.data?.[0] ?? 0
        const y = raw?.data?.[1] ?? 0
        return `${rows[y]} CH${x}: ${value}`
      },
    },
    grid: { left: 54, right: 18, top: 44, bottom: 30 },
    xAxis: {
      type: 'category',
      data: channels,
      splitArea: { show: true },
      axisLabel: { color: '#6c8196' },
    },
    yAxis: {
      type: 'category',
      data: rows,
      splitArea: { show: true },
      axisLabel: { color: '#6c8196' },
    },
    visualMap: {
      min: 0,
      max: maxValue <= 0 ? 1 : maxValue,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: 4,
      inRange: {
        color: ['#e6f4ff', '#91d5ff', '#1890ff', '#0050b3'],
      },
      textStyle: { color: '#5f7287' },
    },
    series: [
      {
        type: 'heatmap',
        data: heatData,
        label: {
          show: false,
        },
      },
    ],
  }
})
</script>

<style scoped>
.panel-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.viz-card {
  height: 260px;
  border: 1px solid #e3ecf7;
  border-radius: 10px;
  background: linear-gradient(180deg, #fbfdff, #f4f9ff);
}

@media (min-width: 1600px) {
  .panel-grid {
    gap: 12px;
  }

  .viz-card {
    height: 290px;
  }
}

@media (max-width: 900px) {
  .panel-grid {
    grid-template-columns: 1fr;
  }
}
</style>
