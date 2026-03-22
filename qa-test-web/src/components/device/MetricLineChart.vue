<template>
  <v-chart class="metric-line-chart" :option="option" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, ScatterChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'
import VChart from 'vue-echarts'

use([
  CanvasRenderer,
  LineChart,
  ScatterChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
])

const props = withDefaults(
  defineProps<{
    title?: string
    seriesName?: string
    values: number[]
    unit?: string
    color?: string
    highlightIndex?: number
  }>(),
  {
    title: '',
    seriesName: 'value',
    unit: '',
    color: '#1890ff',
    highlightIndex: -1,
  }
)

const option = computed<EChartsOption>(() => {
  const points = props.values.map((v, i) => [i, Number(v ?? 0)])
  const categories = props.values.map((_, i) => `CH${i}`)

  const highlight =
    props.highlightIndex >= 0 && props.highlightIndex < props.values.length
      ? [[props.highlightIndex, Number(props.values[props.highlightIndex] ?? 0)]]
      : []

  return {
    title: {
      text: props.title,
      left: 'left',
      textStyle: {
        fontSize: 13,
        fontWeight: 600,
      },
    },
    grid: {
      left: 36,
      right: 20,
      top: 34,
      bottom: 28,
    },
    tooltip: {
      trigger: 'axis',
      formatter: (raw: any) => {
        const rows = Array.isArray(raw) ? raw : [raw]
        if (!rows.length) {
          return ''
        }
        const idx = rows[0].data?.[0] ?? 0
        const value = rows[0].data?.[1] ?? 0
        const suffix = props.unit ? ` ${props.unit}` : ''
        return `CH${idx}: ${value}${suffix}`
      },
    },
    xAxis: {
      type: 'category',
      data: categories,
      boundaryGap: false,
      axisLabel: {
        color: '#8c8c8c',
      },
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#8c8c8c',
      },
      splitLine: {
        lineStyle: {
          color: '#f0f0f0',
        },
      },
    },
    series: [
      {
        name: props.seriesName,
        type: 'line',
        data: points,
        smooth: true,
        showSymbol: false,
        lineStyle: {
          width: 2,
          color: props.color,
        },
        areaStyle: {
          color: props.color,
          opacity: 0.12,
        },
      },
      {
        name: 'selected',
        type: 'scatter',
        data: highlight,
        symbolSize: 8,
        itemStyle: {
          color: '#fa541c',
        },
        z: 9,
      },
    ],
  }
})
</script>

<style scoped>
.metric-line-chart {
  height: 220px;
  width: 100%;
}
</style>
