<template>
  <div class="phm-page">
    <div class="toolbar">
      <div class="title-box">
        <h3>预测性维护与健康管理 (PHM)</h3>
        <p>基于在线率、参数变更、采样连续性和重启行为进行健康评分与风险评估。</p>
      </div>
      <a-space wrap>
        <a-select v-model:value="windowHours" style="width: 150px" @change="refresh">
          <a-select-option :value="24">最近 24 小时</a-select-option>
          <a-select-option :value="72">最近 72 小时</a-select-option>
          <a-select-option :value="168">最近 7 天</a-select-option>
          <a-select-option :value="240">最近 10 天</a-select-option>
        </a-select>
        <a-button type="primary" :loading="state.loading" @click="refresh">刷新评估</a-button>
        <a-button
          danger
          :disabled="selectedRowKeys.length === 0"
          :loading="state.batchDeleting"
          @click="confirmBatchDelete"
        >
          删除选中（{{ selectedRowKeys.length }}）
        </a-button>
      </a-space>
    </div>

    <a-row :gutter="[12, 12]" class="stats-row">
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="设备总数" :value="state.overview?.total ?? 0" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="高风险" :value="riskCount('high') + riskCount('critical')" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="中风险" :value="riskCount('medium')" />
        </a-card>
      </a-col>
      <a-col :xl="6" :lg="12" :md="12" :sm="24" :xs="24">
        <a-card size="small">
          <a-statistic title="低风险" :value="riskCount('low')" />
        </a-card>
      </a-col>
    </a-row>

    <a-card size="small" class="table-card">
      <div class="table-toolbar">
        <a-tag color="default">窗口: {{ windowHours }} 小时</a-tag>
        <a-tag color="blue">当前列表: {{ (state.overview?.items || []).length }} 台</a-tag>
      </div>
      <a-table
        :row-selection="rowSelection"
        :columns="columns"
        :data-source="state.overview?.items || []"
        :loading="state.loading || state.actionLoading"
        row-key="deviceSn"
        size="small"
        :pagination="{ pageSize: 12, showSizeChanger: false }"
        :scroll="{ x: 1480 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'deviceName'">
            <div class="cell-main">{{ record.deviceName || '-' }}</div>
            <div class="cell-sub">{{ record.deviceSn }}</div>
          </template>

          <template v-else-if="column.dataIndex === 'healthScore'">
            <a-progress
              :percent="record.healthScore"
              :stroke-color="riskColor(record.riskLevel)"
              size="small"
              :show-info="false"
            />
            <span class="score-text">{{ record.healthScore }}</span>
          </template>

          <template v-else-if="column.dataIndex === 'riskLevel'">
            <a-tag :color="riskColor(record.riskLevel)">{{ riskLabel(record.riskLevel) }}</a-tag>
          </template>

          <template v-else-if="column.dataIndex === 'onlineRatio'">
            {{ (Number(record.onlineRatio || 0) * 100).toFixed(1) }}%
          </template>

          <template v-else-if="column.dataIndex === 'reasons'">
            <span :title="joinReasons(record.reasons)">
              {{ shortReasons(record.reasons) }}
            </span>
          </template>

          <template v-else-if="column.dataIndex === 'recommendation'">
            <span class="recommend" :title="translateRecommendation(record.recommendation)">
              {{ translateRecommendation(record.recommendation) }}
            </span>
          </template>

          <template v-else-if="column.dataIndex === 'lastSampledAt'">
            {{ formatTime(record.lastSampledAt) }}
          </template>

          <template v-else-if="column.dataIndex === 'action'">
            <a-space>
              <a-button type="link" size="small" @click="openDetail(record.deviceSn)">详情</a-button>
              <a-popconfirm
                title="确认删除该设备 PHM 列表项目及关联数据？"
                ok-text="确认"
                cancel-text="取消"
                @confirm="deleteOne(record.deviceSn)"
              >
                <a-button type="link" danger size="small">删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-drawer
      v-model:visible="detailVisible"
      width="960"
      :title="detail?.summary?.deviceName || detail?.summary?.deviceSn || '设备 PHM 详情'"
      :destroyOnClose="false"
    >
      <a-spin :spinning="detailLoading">
        <template v-if="detail">
          <a-row :gutter="[12, 12]">
            <a-col :span="8">
              <a-card size="small">
                <a-statistic title="健康评分" :value="detail.summary.healthScore" />
              </a-card>
            </a-col>
            <a-col :span="8">
              <a-card size="small">
                <a-statistic
                  title="在线率"
                  :value="Number((detail.summary.onlineRatio * 100).toFixed(1))"
                  suffix="%"
                />
              </a-card>
            </a-col>
            <a-col :span="8">
              <a-card size="small">
                <a-statistic title="参数变更次数" :value="detail.summary.paramChangeCount" />
              </a-card>
            </a-col>
          </a-row>

          <a-card size="small" class="section">
            <a-descriptions :column="2" size="small" bordered>
              <a-descriptions-item label="风险等级">
                <a-tag :color="riskColor(detail.summary.riskLevel)">
                  {{ riskLabel(detail.summary.riskLevel) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="数据新鲜度">
                {{ detail.summary.freshnessMinutes }} 分钟
              </a-descriptions-item>
              <a-descriptions-item label="重启次数">
                {{ detail.summary.rebootCount }}
              </a-descriptions-item>
              <a-descriptions-item label="采样中断次数">
                {{ detail.summary.dataGapCount }}
              </a-descriptions-item>
              <a-descriptions-item label="最近采样时间">
                {{ formatTime(detail.summary.lastSampledAt) }}
              </a-descriptions-item>
              <a-descriptions-item label="维护建议">
                {{ translateRecommendation(detail.summary.recommendation) }}
              </a-descriptions-item>
            </a-descriptions>
          </a-card>

          <a-card size="small" title="关键风险原因" class="section">
            <a-space wrap>
              <a-tag
                v-for="reason in detail.summary.reasons || []"
                :key="reason"
                color="orange"
              >
                {{ translateReason(reason) }}
              </a-tag>
              <span v-if="!(detail.summary.reasons || []).length">无明显风险原因</span>
            </a-space>
          </a-card>

          <a-card size="small" title="最近参数变更" class="section">
            <a-table
              :columns="changeColumns"
              :data-source="detail.recentChanges || []"
              size="small"
              row-key="id"
              :pagination="{ pageSize: 6, showSizeChanger: false }"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.dataIndex === 'changedAt'">
                  {{ formatTime(record.changedAt) }}
                </template>
                <template v-else-if="column.dataIndex === 'oldValue'">
                  <span class="mono-cell">{{ record.oldValue }}</span>
                </template>
                <template v-else-if="column.dataIndex === 'newValue'">
                  <span class="mono-cell">{{ record.newValue }}</span>
                </template>
              </template>
            </a-table>
          </a-card>
        </template>
      </a-spin>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import * as API from '@/api'
import { computed, onMounted, reactive, ref } from 'vue'
import { message, Modal } from 'ant-design-vue'
import type { PHMDeviceDetail, PHMOverviewPayload, PHMRiskLevel } from '@/api/phm'

const windowHours = ref(24)
const detailVisible = ref(false)
const detailLoading = ref(false)
const detail = ref<PHMDeviceDetail | null>(null)
const selectedRowKeys = ref<string[]>([])

const state = reactive<{
  loading: boolean
  actionLoading: boolean
  batchDeleting: boolean
  overview: PHMOverviewPayload | null
}>({
  loading: false,
  actionLoading: false,
  batchDeleting: false,
  overview: null,
})

const columns = [
  { title: '设备', dataIndex: 'deviceName', width: 220, fixed: 'left' as const },
  { title: '健康分', dataIndex: 'healthScore', width: 130 },
  { title: '风险级别', dataIndex: 'riskLevel', width: 100 },
  { title: '在线率', dataIndex: 'onlineRatio', width: 90 },
  { title: '重启', dataIndex: 'rebootCount', width: 70 },
  { title: '参数变更', dataIndex: 'paramChangeCount', width: 100 },
  { title: '最近采样', dataIndex: 'lastSampledAt', width: 170 },
  { title: '主要原因', dataIndex: 'reasons', width: 240 },
  { title: '维护建议', dataIndex: 'recommendation', width: 280 },
  { title: '操作', dataIndex: 'action', width: 140, fixed: 'right' as const },
]

const changeColumns = [
  { title: '时间', dataIndex: 'changedAt', width: 170 },
  { title: '参数路径', dataIndex: 'paramPath', width: 320 },
  { title: '旧值', dataIndex: 'oldValue', width: 220 },
  { title: '新值', dataIndex: 'newValue', width: 220 },
]

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys: Array<string | number>) => {
    selectedRowKeys.value = keys.map((item) => String(item))
  },
}))

const reasonDict: Record<string, string> = {
  'no samples in selected window': '所选窗口内无采样数据',
  'low online ratio': '在线率偏低',
  'stale telemetry': '遥测数据不新鲜',
  'detected reboot or power instability': '存在重启或供电不稳定迹象',
  'hardware version changed': '硬件版本发生变化',
  'sampling pipeline gaps': '采样链路出现断点',
  'frequent system parameter changes': '系统参数变更过于频繁',
}

const recommendationDict: Record<string, string> = {
  'Check collector service, device link, and network connectivity.': '请检查采集服务、设备链路与网络连通性。',
  'Schedule immediate maintenance and inspect power, communication, and key modules.': '建议立即安排维护，重点排查供电、通信与关键模块。',
  'Inspect power stability and thermal conditions due to repeated reboot behavior.': '检测到多次重启，请重点检查供电稳定性与热管理。',
  'Freeze configuration and review parameter strategy to reduce frequent write operations.': '建议冻结配置并复核参数策略，减少频繁写入。',
  'Inspect network quality and telemetry pipeline to reduce data gaps and offline duration.': '请排查网络质量与采集链路，减少数据断点与离线时长。',
  'Device health is stable. Keep routine inspections and trend tracking.': '设备健康稳定，建议保持例行巡检与趋势跟踪。',
}

const translateReason = (reason: string): string => {
  const key = String(reason || '').trim()
  return reasonDict[key] || key || '-'
}

const translateRecommendation = (value: string): string => {
  const key = String(value || '').trim()
  return recommendationDict[key] || key || '-'
}

const joinReasons = (reasons: string[]): string => {
  if (!Array.isArray(reasons) || reasons.length === 0) {
    return '-'
  }
  return reasons.map((item) => translateReason(item)).join('、')
}

const shortReasons = (reasons: string[]): string => {
  if (!Array.isArray(reasons) || reasons.length === 0) {
    return '-'
  }
  return reasons
    .slice(0, 2)
    .map((item) => translateReason(item))
    .join('、')
}

const riskLabel = (risk: PHMRiskLevel): string => {
  if (risk === 'low') {
    return '低风险'
  }
  if (risk === 'medium') {
    return '中风险'
  }
  if (risk === 'high') {
    return '高风险'
  }
  return '严重'
}

const riskColor = (risk: PHMRiskLevel): string => {
  if (risk === 'low') {
    return '#52c41a'
  }
  if (risk === 'medium') {
    return '#faad14'
  }
  if (risk === 'high') {
    return '#fa8c16'
  }
  return '#f5222d'
}

const riskCount = (risk: PHMRiskLevel): number => {
  if (!state.overview) {
    return 0
  }
  return Number(state.overview.riskStats?.[risk] || 0)
}

const formatTime = (value: string): string => {
  if (!value) {
    return '-'
  }
  const time = new Date(value).getTime()
  if (!Number.isFinite(time)) {
    return value
  }
  return new Date(time).toLocaleString()
}

const refresh = async () => {
  state.loading = true
  try {
    const overview = await API.phm.overview(windowHours.value, 200)
    state.overview = overview
    const snSet = new Set((overview.items || []).map((item) => String(item.deviceSn)))
    selectedRowKeys.value = selectedRowKeys.value.filter((sn) => snSet.has(sn))
  } catch (error: any) {
    message.error(error?.message || 'PHM 评估加载失败')
  } finally {
    state.loading = false
  }
}

const openDetail = async (sn: string) => {
  detailVisible.value = true
  detailLoading.value = true
  try {
    detail.value = await API.phm.deviceDetail(sn, windowHours.value, 720, 200)
  } catch (error: any) {
    message.error(error?.message || 'PHM 详情加载失败')
  } finally {
    detailLoading.value = false
  }
}

const deleteOne = async (sn: string) => {
  state.actionLoading = true
  try {
    const res = await API.phm.deleteDevice(sn)
    message.success(
      `删除完成：历史 ${res.historyRowsDeleted} 条，变更 ${res.changeRowsDeleted} 条`
    )
    selectedRowKeys.value = selectedRowKeys.value.filter((item) => item !== sn)
    if (detailVisible.value && detail.value?.summary?.deviceSn === sn) {
      detailVisible.value = false
      detail.value = null
    }
    await refresh()
  } catch (error: any) {
    message.error(error?.message || '删除失败')
  } finally {
    state.actionLoading = false
  }
}

const doBatchDelete = async () => {
  const sns = selectedRowKeys.value.slice()
  if (!sns.length) {
    return
  }
  state.batchDeleting = true
  try {
    const res = await API.phm.deleteDevices(sns)
    message.success(
      `批量删除完成：设备 ${res.requested} 台，历史 ${res.historyRowsDeleted} 条，变更 ${res.changeRowsDeleted} 条`
    )
    selectedRowKeys.value = []
    if (
      detailVisible.value &&
      detail.value?.summary?.deviceSn &&
      sns.includes(detail.value.summary.deviceSn)
    ) {
      detailVisible.value = false
      detail.value = null
    }
    await refresh()
  } catch (error: any) {
    message.error(error?.message || '批量删除失败')
  } finally {
    state.batchDeleting = false
  }
}

const confirmBatchDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    return
  }
  Modal.confirm({
    title: '确认删除选中的 PHM 设备数据吗？',
    content: `将删除 ${selectedRowKeys.value.length} 台设备的历史与参数变更记录，该操作不可恢复。`,
    okText: '确认删除',
    okType: 'danger',
    cancelText: '取消',
    onOk: async () => {
      await doBatchDelete()
    },
  })
}

onMounted(() => {
  void refresh()
})
</script>

<style scoped>
.phm-page {
  padding: 12px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.title-box h3 {
  margin: 0;
  color: #29435c;
  font-size: 16px;
}

.title-box p {
  margin: 4px 0 0;
  color: #6d8296;
  font-size: 12px;
}

.stats-row {
  margin-bottom: 2px;
}

.table-card {
  margin-top: 10px;
}

.table-toolbar {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.cell-main {
  color: #2f455a;
  font-weight: 600;
}

.cell-sub {
  color: #8799aa;
  font-size: 12px;
}

.score-text {
  margin-left: 8px;
  font-weight: 600;
}

.recommend {
  color: #2f455a;
}

.section {
  margin-top: 10px;
}

.mono-cell {
  display: inline-block;
  max-width: 300px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-family: Consolas, 'Courier New', monospace;
}

@media (max-width: 900px) {
  .phm-page {
    padding: 8px;
  }
}
</style>
