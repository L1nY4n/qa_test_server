<template>
  <div class="wrapper">
    <div class="toolbar">
      <div class="toolbar-left">
        <a-input-search
          v-model:value="keyword"
          class="keyword-input"
          :placeholder="t.searchPlaceholder"
          allow-clear
          @search="applyFilter"
        />
        <a-select v-model:value="groupFilter" class="group-filter" @change="applyFilter">
          <a-select-option value="">{{ t.allGroups }}</a-select-option>
          <a-select-option value="virtual-stress">{{ t.stressGroup }}</a-select-option>
        </a-select>
        <div class="toolbar-item">
          <a-switch v-model:checked="onlineOnly" @change="applyFilter" />
          <span class="toolbar-label">{{ t.onlineOnly }}</span>
        </div>
      </div>

      <div class="toolbar-right">
        <a-select
          v-model:value="pageSize"
          class="page-size"
          @change="onPageSizeChange"
        >
          <a-select-option :value="8">8 / {{ t.pageUnit }}</a-select-option>
          <a-select-option :value="12">12 / {{ t.pageUnit }}</a-select-option>
          <a-select-option :value="16">16 / {{ t.pageUnit }}</a-select-option>
          <a-select-option :value="24">24 / {{ t.pageUnit }}</a-select-option>
        </a-select>
        <a-button type="primary" :loading="state.loading" @click="refresh">{{ t.refresh }}</a-button>
        <a-tag color="blue">{{ t.total }}: {{ state.stats?.total ?? 0 }}</a-tag>
        <a-tag color="green">{{ t.online }}: {{ state.stats?.online ?? 0 }}</a-tag>
        <a-tag color="default">{{ t.offline }}: {{ state.stats?.offline ?? 0 }}</a-tag>
      </div>
    </div>

    <div class="list">
      <a-spin :spinning="state.loading">
        <div v-if="state.list.length" class="card-grid">
          <DeviceCard
            v-for="dev in state.list"
            :key="dev.Sn"
            :info="dev"
            :active-within="activeWithinSeconds"
            @open="openDetail"
          />
        </div>
        <a-empty v-else :description="t.noMatchDevices" />
      </a-spin>
    </div>

    <div class="pager" v-if="state.total > pageSize">
      <a-pagination
        v-model:current="page"
        :total="state.total"
        :page-size="pageSize"
        size="small"
        @change="onPageChange"
      />
    </div>

    <a-modal
      v-model:visible="detailVisible"
      :title="detailTitle"
      :footer="null"
      width="92vw"
      :bodyStyle="{ padding: '10px 12px', maxHeight: '78vh', overflow: 'auto' }"
      :destroyOnClose="true"
      :maskClosable="false"
      centered
    >
      <a-spin :spinning="detailLoading">
        <component :is="DeviceDetailAsync" v-if="selectedDevice" :info="selectedDevice" />
        <a-empty v-else :description="t.loadingDetail" />
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import * as API from '@/api'
import DeviceCard from '@/components/device/DeviceCard.vue'
import CreateWebSocket from '@/utils/ws'
import { computed, defineAsyncComponent, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import type { Device } from '@/types/api'
import type { DeviceCardInfo, DeviceListMeta, DeviceStats } from '@/api/device'

const DeviceDetailAsync = defineAsyncComponent(() => import('@/components/device/Device.vue'))

const t = {
  searchPlaceholder: '\u6309 SN / \u540d\u79f0 / \u578b\u53f7 / PN \u641c\u7d22',
  allGroups: '\u5168\u90e8\u5206\u7ec4',
  stressGroup: '\u538b\u529b\u6d4b\u8bd5\u7ec4',
  onlineOnly: '\u4ec5\u5728\u7ebf',
  pageUnit: '\u9875',
  refresh: '\u5237\u65b0',
  total: '\u603b\u6570',
  online: '\u5728\u7ebf',
  offline: '\u79bb\u7ebf',
  noMatchDevices: '\u6682\u65e0\u5339\u914d\u8bbe\u5907',
  loadingDetail: '\u6b63\u5728\u52a0\u8f7d\u8bbe\u5907\u8be6\u60c5...',
  deviceDetail: '\u8bbe\u5907\u8be6\u60c5',
}

const page = ref(1)
const pageSize = ref(12)
const keyword = ref('')
const groupFilter = ref('')
const onlineOnly = ref(false)
const activeWithinSeconds = ref(30)

const detailVisible = ref(false)
const detailLoading = ref(false)
const selectedSn = ref('')
const selectedDevice = ref<Device | null>(null)

const state = reactive<{
  list: DeviceCardInfo[]
  stats: DeviceStats | null
  total: number
  loading: boolean
}>({
  list: [],
  stats: null,
  total: 0,
  loading: false,
})

const detailTitle = computed(() => {
  if (!selectedSn.value) {
    return t.deviceDetail
  }
  return `${t.deviceDetail} - ${selectedSn.value}`
})

let ws: WebSocket | null = null
let pollTimer: ReturnType<typeof setInterval> | null = null
let delayedRefreshTimer: ReturnType<typeof setTimeout> | null = null

const queryParams = () => ({
  keyword: keyword.value.trim() || undefined,
  group: groupFilter.value || undefined,
  onlineOnly: onlineOnly.value,
  activeWithin: activeWithinSeconds.value,
  offset: (page.value - 1) * pageSize.value,
  limit: pageSize.value,
  summary: true,
})

const getList = async () => {
  const res = await API.device.listWithMeta<DeviceCardInfo>(queryParams()) as DeviceListMeta<DeviceCardInfo>
  state.list = res.items || []
  state.total = res.total || 0
}

const getStats = async () => {
  state.stats = await API.device.stats(activeWithinSeconds.value)
}

const fetchDetail = async (sn: string) => {
  detailLoading.value = true
  try {
    const res = await API.device.info(sn, activeWithinSeconds.value)
    selectedDevice.value = res.device
  } catch (error) {
    selectedDevice.value = null
    throw error
  } finally {
    detailLoading.value = false
  }
}

const refresh = async () => {
  state.loading = true
  try {
    const tasks: Array<Promise<unknown>> = [getList(), getStats()]
    if (detailVisible.value && selectedSn.value) {
      tasks.push(fetchDetail(selectedSn.value))
    }
    await Promise.all(tasks)
  } finally {
    state.loading = false
  }
}

const scheduleRefresh = () => {
  if (delayedRefreshTimer) {
    return
  }
  delayedRefreshTimer = setTimeout(() => {
    delayedRefreshTimer = null
    void refresh()
  }, 1200)
}

const applyFilter = () => {
  page.value = 1
  void refresh()
}

const onPageSizeChange = () => {
  page.value = 1
  void refresh()
}

const onPageChange = () => {
  void refresh()
}

const openDetail = (sn: string) => {
  selectedSn.value = sn
  selectedDevice.value = null
  detailVisible.value = true
  void fetchDetail(sn)
}

const isOnline = (lastRxTime: string): boolean => {
  if (!lastRxTime) {
    return false
  }
  const last = new Date(lastRxTime).getTime()
  if (!Number.isFinite(last)) {
    return false
  }
  return Date.now() - last <= activeWithinSeconds.value * 1000
}

const toSafeNumber = (value: unknown): number => {
  const n = Number(value ?? 0)
  return Number.isFinite(n) ? n : 0
}

const decodeAsciiBytes = (value: unknown): string => {
  if (!Array.isArray(value)) {
    return ''
  }
  const chars: string[] = []
  for (const item of value) {
    const code = toSafeNumber(item) & 0xff
    if (code === 0) {
      break
    }
    if (code >= 32 && code <= 126) {
      chars.push(String.fromCharCode(code))
    }
  }
  return chars.join('')
}

const toCardInfo = (device: Device): DeviceCardInfo => {
  const raw = device as any
  const inputReg: any = device?.Packet?.Femto_input_reg ?? {}
  const holdingReg: any = device?.Packet?.Femto_holding_reg ?? {}
  const mon: any = inputReg?.Mon ?? {}
  const uptime = inputReg?.Time?.Uptime
  const pumpMon = mon?.Pump_mon ?? mon?.Femto_input_reg_monitor_pump
  const temp = mon?.Temp

  const uptimePairFromPacket = Array.isArray(uptime) ? uptime.map((v: unknown) => toSafeNumber(v)).slice(0, 2) : [0, 0]
  const uptimePairFromSummary = Array.isArray(raw?.Uptime) ? raw.Uptime.map((v: unknown) => toSafeNumber(v)).slice(0, 2) : null
  const uptimePair = uptimePairFromSummary && uptimePairFromSummary.length >= 2 ? uptimePairFromSummary : uptimePairFromPacket

  const mergedUptimeSeconds = ((uptimePair[0] ?? 0) << 16) + (uptimePair[1] ?? 0)
  const summaryUptimeSeconds = toSafeNumber(raw?.Uptime_seconds)
  const uptimeSeconds = summaryUptimeSeconds > 0 ? summaryUptimeSeconds : (mergedUptimeSeconds > 0 ? mergedUptimeSeconds : 0)

  const modelFromPacket = decodeAsciiBytes(holdingReg?.Laser_para?.Laser_info?.Model)
  const pnFromPacket = decodeAsciiBytes(holdingReg?.Laser_para?.Laser_info?.PN)

  const model = String(raw?.Model || modelFromPacket || '')
  const pn = String(raw?.PN || pnFromPacket || '')
  const lastRxTime = String(raw?.Last_rx_time || device.Last_rx_time || '')

  const summaryPumpCount = toSafeNumber(raw?.Pump_count)
  const summaryTempCount = toSafeNumber(raw?.Temp_count)
  const summaryHw = toSafeNumber(raw?.Hardware_bate)
  const summaryLaserStatus = toSafeNumber(raw?.Laser_status)
  const summaryLaserReady = toSafeNumber(raw?.Laser_ready)
  const summaryWavelength = toSafeNumber(raw?.Laser_wavelength)

  return {
    Sn: String(device.Sn || raw?.Sn || ''),
    Name: String(device.Name || raw?.Name || ''),
    Group: String(raw?.Group || ''),
    Model: model,
    PN: pn,
    Last_rx_time: lastRxTime,
    Hardware_bate: summaryHw > 0 ? summaryHw : toSafeNumber(inputReg?.Bate?.Hardware_bate),
    Uptime: uptimePair,
    Uptime_seconds: uptimeSeconds,
    Pump_count: summaryPumpCount > 0 ? summaryPumpCount : (Array.isArray(pumpMon) ? pumpMon.length : 0),
    Temp_count: summaryTempCount > 0 ? summaryTempCount : (Array.isArray(temp) ? temp.length : 0),
    Laser_status: summaryLaserStatus > 0 ? summaryLaserStatus : toSafeNumber(inputReg?.Status),
    Laser_ready: summaryLaserReady > 0 ? summaryLaserReady : toSafeNumber(holdingReg?.User_para?.Laser_ready),
    Laser_wavelength: summaryWavelength > 0 ? summaryWavelength : toSafeNumber(holdingReg?.User_para?.Laser_wavelength),
    Online: typeof raw?.Online === 'boolean' ? raw.Online : isOnline(lastRxTime),
  }
}

const updateSummaryFromWs = (updateDevice: Device) => {
  const index = state.list.findIndex((d) => d.Sn === updateDevice.Sn)

  if (index > -1) {
    state.list.splice(index, 1, toCardInfo(updateDevice))
    return
  }

  scheduleRefresh()
}

const parseWsPayload = async (payload: unknown): Promise<Device | null> => {
  try {
    if (typeof payload === 'string') {
      return JSON.parse(payload) as Device
    }
    if (payload instanceof Blob) {
      return JSON.parse(await payload.text()) as Device
    }
    if (payload instanceof ArrayBuffer) {
      const decoded = new TextDecoder('utf-8').decode(payload)
      return JSON.parse(decoded) as Device
    }
    return null
  } catch {
    return null
  }
}

const handleWsMessage = async (payload: unknown) => {
  const updateDevice = await parseWsPayload(payload)
  if (!updateDevice || !updateDevice.Sn) {
    return
  }

  updateSummaryFromWs(updateDevice)

  if (detailVisible.value && selectedSn.value === updateDevice.Sn) {
    selectedDevice.value = updateDevice
  }
}

const webSocketConnect = () => {
  ws = CreateWebSocket('/realtime/register/device_upload', 'test')
  ws.onmessage = (event) => {
    void handleWsMessage(event.data)
  }

  ws.onclose = () => {
    setTimeout(() => {
      webSocketConnect()
    }, 10 * 1000)
  }
}

watch(detailVisible, (open) => {
  if (!open) {
    selectedDevice.value = null
    return
  }
  if (selectedSn.value && !selectedDevice.value) {
    void fetchDetail(selectedSn.value)
  }
})

onMounted(() => {
  webSocketConnect()
  void refresh()
  pollTimer = setInterval(() => {
    void refresh()
  }, 15000)
})

onBeforeUnmount(() => {
  ws?.close()
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
  if (delayedRefreshTimer) {
    clearTimeout(delayedRefreshTimer)
    delayedRefreshTimer = null
  }
})
</script>

<style lang="less">
.wrapper {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;

  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    flex-wrap: wrap;
    padding: 10px 14px;
    background: linear-gradient(145deg, rgba(236, 244, 253, 0.95), rgba(224, 236, 248, 0.9));
    border-bottom: 1px solid #cfddec;
    backdrop-filter: blur(5px);
  }

  .toolbar-left,
  .toolbar-right {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
  }

  .toolbar-item {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .toolbar-label {
    color: #607387;
    font-size: 12px;
  }

  .keyword-input {
    width: clamp(220px, 30vw, 360px);
  }

  .group-filter {
    width: 170px;
  }

  .page-size {
    width: 92px;
  }

  .list {
    flex: 1;
    overflow: auto;
    padding: 14px;
    background: linear-gradient(180deg, rgba(245, 250, 255, 0.68), rgba(236, 244, 253, 0.48));
  }

  .card-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 12px;
  }

  .pager {
    padding: 8px 14px 12px;
    display: flex;
    justify-content: flex-end;
    border-top: 1px solid #dbe6f2;
    background: linear-gradient(180deg, rgba(240, 247, 255, 0.9), rgba(231, 240, 250, 0.85));
  }

  :deep(.ant-tag-blue) {
    border-color: #76a8e0;
    background: #e4f0ff;
    color: #255991;
  }

  :deep(.ant-tag-green) {
    border-color: #78d6a4;
    background: #e5fbe9;
    color: #1e7148;
  }

  @media (max-width: 900px) {
    .toolbar {
      padding: 8px 10px;
      gap: 8px;
    }

    .list {
      padding: 10px;
    }

    .card-grid {
      grid-template-columns: 1fr;
      gap: 10px;
    }

    .pager {
      justify-content: center;
      padding: 8px 10px;
    }
  }
}
</style>
