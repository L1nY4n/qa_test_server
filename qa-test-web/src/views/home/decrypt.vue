<template>
  <div class="page">
    <a-card title="时间密钥工具（运维及以上）" size="small">
      <a-alert
        type="info"
        show-icon
        message="支持两种模式：1) 输入 SN+时间解算密钥；2) 输入 SN+密钥反解时间。所有操作会写入审计日志。"
      />

      <div v-if="isOpsOrAbove" class="content">
        <div class="mode-panel">
          <div class="panel-title">SN + 时间 解算密钥</div>
          <a-row :gutter="[12, 12]" class="config-row">
            <a-col :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
              <a-form-item label="设备 SN">
                <a-input v-model:value="generateForm.deviceSn" placeholder="请输入设备 SN" />
              </a-form-item>
            </a-col>
            <a-col :xl="10" :lg="10" :md="10" :sm="24" :xs="24">
              <a-form-item label="目标时间">
                <a-date-picker
                  v-model:value="generateForm.timeText"
                  show-time
                  format="YYYY-MM-DD HH:mm:ss"
                  value-format="YYYY-MM-DD HH:mm:ss"
                  placeholder="请选择目标时间"
                  :allow-clear="false"
                  style="width: 100%"
                />
              </a-form-item>
            </a-col>
            <a-col :xl="6" :lg="6" :md="6" :sm="24" :xs="24" class="action-col">
              <a-space class="action-group" wrap>
                <a-button type="primary" :loading="generateLoading" @click="runGenerate">
                  解算密钥
                </a-button>
                <a-button @click="fillNow">填入当前时间</a-button>
              </a-space>
            </a-col>
          </a-row>
        </div>

        <a-alert
          v-if="generateResult"
          class="tip"
          type="success"
          show-icon
          :message="`解算成功：${generateResult.inputTimeText} -> 密钥 ${generateResult.key}`"
        />
        <a-descriptions v-if="generateResult" size="small" bordered :column="1" class="result-block">
          <a-descriptions-item label="密钥（明文 28 位）">{{ generateResult.key }}</a-descriptions-item>
          <a-descriptions-item label="密钥（HEX）">{{ generateResult.keyHex }}</a-descriptions-item>
          <a-descriptions-item label="Raw Head">{{ generateResult.rawHead }}</a-descriptions-item>
          <a-descriptions-item label="Raw Tail">{{ generateResult.rawTail }}</a-descriptions-item>
        </a-descriptions>

        <div class="mode-panel">
          <div class="panel-title">SN + 密钥 反解时间</div>
          <a-row :gutter="[12, 12]" class="config-row">
            <a-col :xl="7" :lg="7" :md="8" :sm="24" :xs="24">
              <a-form-item label="设备 SN">
                <a-input v-model:value="decryptForm.deviceSn" placeholder="请输入设备 SN" />
              </a-form-item>
            </a-col>
            <a-col :xl="5" :lg="5" :md="6" :sm="24" :xs="24">
              <a-form-item label="输入模式">
                <a-select v-model:value="decryptForm.inputMode">
                  <a-select-option value="plain">明文密钥（28字符）</a-select-option>
                  <a-select-option value="hex">十六进制（56字符）</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :xl="8" :lg="8" :md="10" :sm="24" :xs="24">
              <a-form-item :label="decryptForm.inputMode === 'hex' ? '密钥 HEX' : '密钥'">
                <a-input
                  v-model:value="decryptForm.keyValue"
                  :placeholder="decryptForm.inputMode === 'hex' ? '例如 414243...' : '请输入 28 位密钥'"
                />
              </a-form-item>
            </a-col>
            <a-col :xl="4" :lg="4" :md="24" :sm="24" :xs="24" class="action-col">
              <a-space class="action-group" wrap>
                <a-button type="primary" :loading="decryptLoading" @click="runDecrypt">执行反解</a-button>
              </a-space>
            </a-col>
          </a-row>
        </div>

        <div class="toolbar">
          <div class="toolbar-left">
            <a-space wrap>
            <a-select v-model:value="logFilter.operation" style="width: 160px">
              <a-select-option value="all">全部日志</a-select-option>
              <a-select-option value="generate">仅解算密钥</a-select-option>
              <a-select-option value="decode">仅反解时间</a-select-option>
            </a-select>
            <a-button :loading="decryptLogLoading" @click="loadDecryptLogs">刷新日志</a-button>
            </a-space>
          </div>
          <div class="toolbar-right">
            <a-tag color="default">日志总数: {{ decryptLogTotal }}</a-tag>
          </div>
        </div>

        <a-alert
          v-if="decryptResult"
          class="tip"
          type="success"
          show-icon
          :message="`反解成功：到期时间 ${decryptResult.decodedAtText || `${decryptResult.fullYear}-${decryptResult.decodedMonth}-${decryptResult.decodedDay} ${decryptResult.decodedHour}:${decryptResult.decodedMinute}:${decryptResult.decodedSecond}`}`"
        />

        <a-table
          size="small"
          :columns="decryptLogColumns"
          :data-source="decryptLogs"
          :loading="decryptLogLoading"
          :pagination="false"
          row-key="id"
          :scroll="{ x: 1360, y: 420 }"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'createdAt'">
              {{ formatTime(record.createdAt) }}
            </template>
            <template v-else-if="column.dataIndex === 'operation'">
              <a-tag :color="record.operation === 'generate' ? 'purple' : 'blue'">
                {{ record.operation === 'generate' ? '解算密钥' : '反解时间' }}
              </a-tag>
            </template>
            <template v-else-if="column.dataIndex === 'success'">
              <a-tag :color="record.success ? 'green' : 'red'">{{ record.success ? '成功' : '失败' }}</a-tag>
            </template>
            <template v-else-if="column.dataIndex === 'decodedAt'">
              <span v-if="record.success">
                {{ `${2000 + Number(record.decodedYear || 0)}-${String(record.decodedMonth || 0).padStart(2, '0')}-${String(record.decodedDay || 0).padStart(2, '0')} ${String(record.decodedHour || 0).padStart(2, '0')}:${String(record.decodedMinute || 0).padStart(2, '0')}:${String(record.decodedSecond || 0).padStart(2, '0')}` }}
              </span>
              <span v-else>{{ record.errorMessage || '-' }}</span>
            </template>
          </template>
        </a-table>
      </div>

      <a-empty v-else description="当前账号无权限访问该功能" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import * as API from '@/api'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import type {
  DecryptLogItem,
  TimeKeyDecryptResult,
  TimeKeyGenerateResult,
} from '@/api/debug'
import { getCurrentUser } from '@/utils/auth'
import { useRoute } from 'vue-router'

const route = useRoute()
const userRole = ref(getCurrentUser()?.role || '')

const decryptLoading = ref(false)
const generateLoading = ref(false)
const decryptLogLoading = ref(false)
const decryptResult = ref<TimeKeyDecryptResult | null>(null)
const generateResult = ref<TimeKeyGenerateResult | null>(null)
const decryptLogs = ref<DecryptLogItem[]>([])
const decryptLogTotal = ref(0)

const isOpsOrAbove = computed(() => {
  const role = String(userRole.value || '').toLowerCase()
  return role === 'operator' || role === 'admin'
})

const decryptForm = reactive({
  deviceSn: '',
  inputMode: 'plain' as 'plain' | 'hex',
  keyValue: '',
})

const generateForm = reactive({
  deviceSn: '',
  timeText: '',
})

const logFilter = reactive({
  operation: 'all' as 'all' | 'decode' | 'generate',
})

const decryptLogColumns = [
  { title: '时间', dataIndex: 'createdAt', width: 170 },
  { title: '操作类型', dataIndex: 'operation', width: 100 },
  { title: '操作者', dataIndex: 'operatorName', width: 120 },
  { title: '角色', dataIndex: 'operatorRole', width: 90 },
  { title: '设备SN', dataIndex: 'deviceSn', width: 180 },
  { title: '结果', dataIndex: 'success', width: 80 },
  { title: '时间字段', dataIndex: 'decodedAt', width: 320, ellipsis: true },
  { title: '密钥摘要', dataIndex: 'keyPreview', width: 160, ellipsis: true },
  { title: '来源IP', dataIndex: 'sourceIp', width: 120 },
]

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

const nowText = (): string => {
  const d = new Date()
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

const fillNow = () => {
  generateForm.timeText = nowText()
}

const loadDecryptLogs = async () => {
  if (!isOpsOrAbove.value) {
    return
  }

  decryptLogLoading.value = true
  try {
    const res = await API.debug.decryptLogs({
      offset: 0,
      limit: 50,
      deviceSn: (decryptForm.deviceSn || generateForm.deviceSn).trim() || undefined,
      operation: logFilter.operation === 'all' ? undefined : logFilter.operation,
    })
    decryptLogs.value = res.items || []
    decryptLogTotal.value = Number(res.total || 0)
  } catch (error: any) {
    message.error(error?.message || '加载日志失败')
  } finally {
    decryptLogLoading.value = false
  }
}

const runGenerate = async () => {
  const sn = generateForm.deviceSn.trim()
  const timeText = generateForm.timeText.trim()
  if (!sn) {
    message.warning('请输入设备 SN')
    return
  }
  if (!timeText) {
    message.warning('请输入目标时间')
    return
  }

  generateLoading.value = true
  try {
    generateResult.value = await API.debug.generateTimeKey({
      deviceSn: sn,
      time: timeText,
    })
    decryptForm.deviceSn = sn
    decryptForm.inputMode = 'plain'
    decryptForm.keyValue = generateResult.value.key
    message.success('密钥解算成功')
    await loadDecryptLogs()
  } catch (error: any) {
    generateResult.value = null
    message.error(error?.message || '密钥解算失败')
    await loadDecryptLogs()
  } finally {
    generateLoading.value = false
  }
}

const runDecrypt = async () => {
  const sn = decryptForm.deviceSn.trim()
  if (!sn) {
    message.warning('请输入设备 SN')
    return
  }
  const keyValue = decryptForm.keyValue.trim()
  if (!keyValue) {
    message.warning('请输入密钥')
    return
  }

  decryptLoading.value = true
  try {
    decryptResult.value = await API.debug.decryptTimeKey({
      deviceSn: sn,
      key: decryptForm.inputMode === 'plain' ? keyValue : undefined,
      keyHex: decryptForm.inputMode === 'hex' ? keyValue : undefined,
    })
    message.success('反解成功')
    await loadDecryptLogs()
  } catch (error: any) {
    decryptResult.value = null
    message.error(error?.message || '反解失败')
    await loadDecryptLogs()
  } finally {
    decryptLoading.value = false
  }
}

onMounted(() => {
  userRole.value = getCurrentUser()?.role || ''
  fillNow()
  void loadDecryptLogs()
})

watch(
  () => route.fullPath,
  () => {
    userRole.value = getCurrentUser()?.role || ''
  }
)
</script>

<style scoped>
.page {
  padding: 12px;
}

.content {
  margin-top: 12px;
}

.config-row {
  margin-bottom: 0;
}

.mode-panel {
  margin-top: 12px;
  border: 1px solid #e7edf5;
  border-radius: 10px;
  padding: 12px;
  background: #fbfdff;
}

.panel-title {
  color: #2f455a;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 6px;
}

.action-col {
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
}

.action-group {
  width: 100%;
  justify-content: flex-end;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  gap: 10px;
  flex-wrap: wrap;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.tip {
  margin-bottom: 10px;
}

.result-block {
  margin-bottom: 10px;
}

@media (min-width: 1600px) {
  .page {
    padding: 16px;
  }

  .content {
    margin-top: 14px;
  }

  .mode-panel {
    margin-top: 14px;
    border-radius: 12px;
    padding: 14px;
  }

  .panel-title {
    font-size: 15px;
  }

  .toolbar {
    margin-bottom: 14px;
    gap: 12px;
  }
}

@media (min-width: 2200px) {
  .page {
    padding: 20px;
  }

  .mode-panel {
    padding: 16px;
  }
}

@media (max-width: 900px) {
  .page {
    padding: 8px;
  }

  .mode-panel {
    padding: 10px;
  }

  .action-col {
    justify-content: flex-start;
  }

  .action-group {
    justify-content: flex-start;
  }

  .toolbar-right {
    width: 100%;
    justify-content: flex-start;
  }
}

@media (max-width: 680px) {
  .mode-panel {
    margin-top: 10px;
    padding: 8px;
  }

  .toolbar {
    gap: 8px;
  }

  .toolbar-left,
  .toolbar-right {
    width: 100%;
  }

  .toolbar-left :deep(.ant-space) {
    width: 100%;
  }

  .toolbar-left :deep(.ant-space .ant-space-item) {
    width: 100%;
  }

  .toolbar-left :deep(.ant-select),
  .toolbar-left :deep(.ant-btn) {
    width: 100%;
  }

  .action-group :deep(.ant-btn) {
    width: 100%;
  }
}
</style>
