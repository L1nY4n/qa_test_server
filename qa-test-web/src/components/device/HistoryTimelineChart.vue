<template>
  <v-chart class="timeline-chart" :option="option" autoresize />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  DataZoomComponent,
  ToolboxComponent,
} from 'echarts/components'
import type { EChartsOption } from 'echarts'
import VChart from 'vue-echarts'

use([
  CanvasRenderer,
  LineChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  DataZoomComponent,
  ToolboxComponent,
])

interface TimelineSeries {
  name: string
  data: Array<[string | number, number]>
  color?: string
  yAxisIndex?: number
  unit?: string
  dashed?: boolean
}

interface TimelineYAxis {
  name: string
  min?: number | 'auto'
  max?: number | 'auto'
  position?: 'left' | 'right'
  offset?: number
}

const props = withDefaults(
  defineProps<{
    title?: string
    unit?: string
    series: TimelineSeries[]
    yAxisMin?: number | 'auto'
    yAxes?: TimelineYAxis[]
  }>(),
  {
    title: '',
    unit: '',
    yAxisMin: 'auto',
    yAxes: () => [],
  }
)

const option = computed<EChartsOption>(() => {
  const unit = props.unit ? ` (${props.unit})` : ''
  const yAxes = props.yAxes.length
    ? props.yAxes.map((axis) => ({
        type: 'value' as const,
        name: axis.name || '',
        min: axis.min === 'auto' || axis.min === undefined ? undefined : axis.min,
        max: axis.max === 'auto' || axis.max === undefined ? undefined : axis.max,
        position: axis.position || 'left',
        offset: axis.offset || 0,
        splitLine: {
          show: !axis.position || axis.position === 'left',
          lineStyle: {
            color: '#edf2f8',
          },
        },
        axisLabel: {
          color: '#7b8da0',
        },
      }))
    : [
        {
          type: 'value' as const,
          name: props.unit || '',
          min: props.yAxisMin === 'auto' ? undefined : props.yAxisMin,
          splitLine: {
            lineStyle: {
              color: '#edf2f8',
            },
          },
          axisLabel: {
            color: '#7b8da0',
          },
        },
      ]

  return {
    animation: false,
    grid: {
      left: 52,
      right: props.yAxes.length > 1 ? 90 : 20,
      top: 42,
      bottom: 72,
    },
    legend: {
      top: 8,
      type: 'scroll',
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
      formatter: (raw: any) => {
        const rows = Array.isArray(raw) ? raw : [raw]
        if (!rows.length) {
          return ''
        }
        const ts = rows[0]?.axisValueLabel || rows[0]?.axisValue || ''
        const lines = rows.map((row: any) => {
          const idx = Number(row?.seriesIndex ?? 0)
          const cfg = props.series[idx]
          const suffix = cfg?.unit ? ` ${cfg.unit}` : unit
          const rawValue = row?.data?.[1] ?? row?.value?.[1]
          const value = Number(rawValue)
          if (!Number.isFinite(value)) {
            return `${row.marker || ''}${row.seriesName}: --`
          }
          return `${row.marker || ''}${row.seriesName}: ${value.toFixed(2)}${suffix}`
        })
        return [ts, ...lines].join('<br/>')
      },
    },
    toolbox: {
      right: 8,
      feature: {
        dataZoom: { yAxisIndex: 'none' },
        restore: {},
        saveAsImage: {},
      },
    },
    xAxis: {
      type: 'time',
      axisLabel: {
        color: '#7b8da0',
      },
    },
    yAxis: yAxes,
    dataZoom: [
      {
        type: 'inside',
        throttle: 80,
      },
      {
        type: 'slider',
        height: 26,
        bottom: 18,
      },
    ],
    series: props.series.map((item) => ({
      name: item.name,
      type: 'line',
      data: item.data,
      yAxisIndex: item.yAxisIndex || 0,
      showSymbol: false,
      smooth: false,
      sampling: 'lttb',
      lineStyle: {
        width: 1.6,
        color: item.color,
        type: item.dashed ? 'dashed' : 'solid',
      },
      emphasis: {
        focus: 'series',
      },
      progressive: 800,
      progressiveThreshold: 2000,
    })),
  }
})
</script>

<style scoped>
.timeline-chart {
  width: 100%;
  height: 280px;
}
</style>
