<template>
  <div class="device-insight">
    <a-row :gutter="[10, 10]" class="top-cards">
      <a-col :xl="4" :lg="8" :md="8" :sm="12" :xs="24">
        <a-card size="small">
          <a-statistic :title="t.laserStatus" :value="laserStatusText" />
        </a-card>
      </a-col>
      <a-col :xl="4" :lg="8" :md="8" :sm="12" :xs="24">
        <a-card size="small">
          <a-statistic :title="t.readyState" :value="readyStatusText" />
        </a-card>
      </a-col>
      <a-col :xl="4" :lg="8" :md="8" :sm="12" :xs="24">
        <a-card size="small">
          <a-statistic :title="t.wavelengthMode" :value="wavelengthText" />
        </a-card>
      </a-col>
      <a-col :xl="4" :lg="8" :md="8" :sm="12" :xs="24">
        <a-card size="small">
          <a-statistic :title="t.frequency" :value="numberOrDash(userPara.Freq)" suffix="Hz" />
        </a-card>
      </a-col>
      <a-col :xl="4" :lg="8" :md="8" :sm="12" :xs="24">
        <a-card size="small">
          <a-statistic :title="t.pulseWidth" :value="numberOrDash(userPara.Puse_width)" suffix="ns" />
        </a-card>
      </a-col>
      <a-col :xl="4" :lg="8" :md="8" :sm="12" :xs="24">
        <a-card size="small">
          <a-statistic :title="t.outputPower" :value="numberOrDash(userPara.Laser_power)" />
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[10, 10]">
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" :title="t.deviceIdentity">
          <a-descriptions size="small" :column="1" bordered>
            <a-descriptions-item :label="t.modelAscii">{{ modelText }}</a-descriptions-item>
            <a-descriptions-item :label="t.snAscii">{{ snText }}</a-descriptions-item>
            <a-descriptions-item :label="t.pnAscii">{{ pnText }}</a-descriptions-item>
            <a-descriptions-item :label="t.laserSerial">{{ numberOrDash(laserInfo.Laser_serial) }}</a-descriptions-item>
            <a-descriptions-item :label="t.powerLevel">{{ numberOrDash(laserInfo.Laser_Power_level) }}</a-descriptions-item>
            <a-descriptions-item :label="t.hardwareVersion">{{ numberOrDash(inputReg?.Bate?.Hardware_bate) }}</a-descriptions-item>
            <a-descriptions-item :label="t.mcuVersion">{{ joinVersion(inputReg?.Bate?.Mcu_app_bate) }}</a-descriptions-item>
            <a-descriptions-item :label="t.fpgaVersion">{{ joinVersion(inputReg?.Bate?.Fpga_bate) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" :title="t.networkParams">
          <a-descriptions size="small" :column="1" bordered>
            <a-descriptions-item :label="t.espStatus">{{ numberOrDash(espMon?.Status) }}</a-descriptions-item>
            <a-descriptions-item :label="t.espEnable">{{ switchText(espModule?.En) }}</a-descriptions-item>
            <a-descriptions-item :label="t.espMode">{{ numberOrDash(espModule?.Mode) }}</a-descriptions-item>
            <a-descriptions-item :label="t.ssidAscii">{{ ssidText }}</a-descriptions-item>
            <a-descriptions-item :label="t.wifiIp">{{ formatIPv4(espMon?.Wifi_ip) }}</a-descriptions-item>
            <a-descriptions-item :label="t.wifiGateway">{{ formatIPv4(espMon?.Wifi_gateway) }}</a-descriptions-item>
            <a-descriptions-item :label="t.wifiMask">{{ formatIPv4(espMon?.Wifi_netmask) }}</a-descriptions-item>
            <a-descriptions-item :label="t.ethIp">{{ formatIPv4(espMon?.Eth_ip) }}</a-descriptions-item>
            <a-descriptions-item :label="t.ethGateway">{{ formatIPv4(espMon?.Eth_gateway) }}</a-descriptions-item>
            <a-descriptions-item :label="t.ethMask">{{ formatIPv4(espMon?.Eth_netmask) }}</a-descriptions-item>
            <a-descriptions-item :label="t.socketIp">{{ formatIPv4(espMon?.Socket_ip) }}</a-descriptions-item>
            <a-descriptions-item :label="t.socketPort">{{ numberOrDash(espMon?.Socket_port) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[10, 10]">
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" :title="t.tempChannels">
          <MetricLineChart
            :title="t.tempDistribution"
            :series-name="t.tempSeries"
            :values="tempValues"
            unit="0.1C"
            color="#fa8c16"
          />
        </a-card>
      </a-col>

      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" :title="t.voltageChannels">
          <MetricLineChart
            :title="t.voltageDistribution"
            :series-name="t.voltageSeries"
            :values="voltageValues"
            unit="raw"
            color="#13c2c2"
          />
        </a-card>
      </a-col>
    </a-row>

    <a-card size="small" :title="t.pumpChannels" class="pump-card">
      <a-table
        size="small"
        :columns="pumpColumns"
        :data-source="pumpRows"
        :pagination="false"
        row-key="index"
        :scroll="{ y: 240 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'switch'">
            <a-tag :color="record.switch === t.switchOn ? 'green' : 'default'">{{ record.switch }}</a-tag>
          </template>
          <template v-else-if="column.dataIndex === 'delta'">
            <a-tag :color="record.delta > 30 ? 'orange' : record.delta < -30 ? 'red' : 'default'">
              {{ record.delta }}
            </a-tag>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-row :gutter="[10, 10]">
      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" :title="t.runtimeCounters">
          <a-descriptions size="small" :column="1" bordered>
            <a-descriptions-item :label="t.onlineBitmap">{{ formatHexWords(inputReg?.Online) }}</a-descriptions-item>
            <a-descriptions-item :label="t.laserSwCountdown">{{ numberOrDash(timeReg?.Laser_sw_countdown) }}</a-descriptions-item>
            <a-descriptions-item :label="t.singlePumpRuntime">{{ formatU16PairDuration(timeReg?.Pump_sig_work_time) }}</a-descriptions-item>
            <a-descriptions-item :label="t.pumpRuntime">{{ formatU16PairDuration(timeReg?.Pump_work_time) }}</a-descriptions-item>
            <a-descriptions-item :label="t.emissionRuntime">{{ formatU16PairDuration(timeReg?.Emission_time) }}</a-descriptions-item>
            <a-descriptions-item :label="t.uptime">{{ formatU16PairDuration(timeReg?.Uptime) }}</a-descriptions-item>
            <a-descriptions-item :label="t.totalUptime">{{ formatU16PairDuration(timeReg?.Total_uptime) }}</a-descriptions-item>
            <a-descriptions-item :label="t.systemTime">{{ formatTimeDesc(timeReg?.Sys_time) }}</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>

      <a-col :xl="12" :lg="24" :md="24" :sm="24" :xs="24">
        <a-card size="small" :title="t.alarmBitmap">
          <div class="alarm-block">
            <div class="alarm-title">{{ t.currentAlarm }}</div>
            <div class="alarm-tags">
              <a-tag v-for="index in nowAlarmBitsPreview" :key="`now-${index}`" color="red">
                {{ t.bitPrefix }} {{ index }}
              </a-tag>
              <a-tag v-if="nowAlarmBitsRemain > 0" color="default">+{{ nowAlarmBitsRemain }}</a-tag>
              <a-tag v-if="nowAlarmBits.length === 0">{{ t.none }}</a-tag>
            </div>
          </div>
          <div class="alarm-block">
            <div class="alarm-title">{{ t.historyAlarm }}</div>
            <div class="alarm-tags">
              <a-tag v-for="index in historyAlarmBitsPreview" :key="`hist-${index}`" color="orange">
                {{ t.bitPrefix }} {{ index }}
              </a-tag>
              <a-tag v-if="historyAlarmBitsRemain > 0" color="default">+{{ historyAlarmBitsRemain }}</a-tag>
              <a-tag v-if="historyAlarmBits.length === 0">{{ t.none }}</a-tag>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Device } from '@/types/api'
import MetricLineChart from './MetricLineChart.vue'

const props = defineProps<{ info: Device }>()

const t = {
  laserStatus: '\u6fc0\u5149\u72b6\u6001',
  readyState: '\u5c31\u7eea\u72b6\u6001',
  wavelengthMode: '\u6ce2\u957f\u6a21\u5f0f',
  frequency: '\u9891\u7387',
  pulseWidth: '\u8109\u5bbd',
  outputPower: '\u8f93\u51fa\u529f\u7387',
  deviceIdentity: '\u8bbe\u5907\u6807\u8bc6',
  modelAscii: '\u578b\u53f7 (ASCII)',
  snAscii: 'SN (ASCII)',
  pnAscii: 'PN (ASCII)',
  laserSerial: '\u5e8f\u5217\u53f7',
  powerLevel: '\u529f\u7387\u6863\u4f4d',
  hardwareVersion: '\u786c\u4ef6\u7248\u672c',
  mcuVersion: 'MCU \u7248\u672c',
  fpgaVersion: 'FPGA \u7248\u672c',
  networkParams: '\u7f51\u7edc\u53c2\u6570',
  espStatus: 'ESP \u72b6\u6001',
  espEnable: 'ESP \u4f7f\u80fd',
  espMode: 'ESP \u6a21\u5f0f',
  ssidAscii: 'SSID (ASCII)',
  wifiIp: 'WIFI IP',
  wifiGateway: 'WIFI \u7f51\u5173',
  wifiMask: 'WIFI \u63a9\u7801',
  ethIp: 'ETH IP',
  ethGateway: 'ETH \u7f51\u5173',
  ethMask: 'ETH \u63a9\u7801',
  socketIp: 'Socket IP',
  socketPort: 'Socket \u7aef\u53e3',
  tempChannels: '\u6e29\u5ea6\u901a\u9053',
  tempDistribution: '\u6e29\u5ea6\u5206\u5e03',
  tempSeries: '\u6e29\u5ea6',
  voltageChannels: '\u7535\u538b\u901a\u9053',
  voltageDistribution: '\u7535\u538b\u5206\u5e03',
  voltageSeries: '\u7535\u538b',
  pumpChannels: '\u6cf5\u901a\u9053',
  runtimeCounters: '\u8fd0\u884c\u8ba1\u6570',
  onlineBitmap: '\u5728\u7ebf\u4f4d\u56fe',
  laserSwCountdown: '\u5f00\u673a\u5012\u8ba1\u65f6',
  singlePumpRuntime: '\u5355\u6cf5\u5de5\u4f5c\u65f6\u957f',
  pumpRuntime: '\u6cf5\u603b\u5de5\u4f5c\u65f6\u957f',
  emissionRuntime: '\u51fa\u5149\u65f6\u957f',
  uptime: '\u8fd0\u884c\u65f6\u957f',
  totalUptime: '\u7d2f\u8ba1\u8fd0\u884c\u65f6\u957f',
  systemTime: '\u7cfb\u7edf\u65f6\u95f4',
  alarmBitmap: '\u544a\u8b66\u4f4d\u56fe',
  currentAlarm: '\u5f53\u524d\u544a\u8b66 (Now)',
  historyAlarm: '\u5386\u53f2\u544a\u8b66 (History)',
  bitPrefix: '\u4f4d',
  none: '\u65e0',
  switchOn: '\u5f00',
  switchOff: '\u5173',
  preheat: '\u9884\u70ed',
  standby: '\u5f85\u673a',
  ready: '\u5c31\u7eea',
  fault: '\u6545\u969c',
  notReady: '\u672a\u5c31\u7eea',
  warming: '\u9884\u70ed\u4e2d',
  readyDone: '\u5df2\u5c31\u7eea',
  ir: '\u7ea2\u5916',
  green: '\u7eff\u5149',
  uv: '\u7d2b\u5916',
  channel: '\u901a\u9053',
  state: '\u72b6\u6001',
  actualCurrent: '\u5b9e\u9645\u7535\u6d41',
  fpgaCurrent: 'FPGA \u7535\u6d41',
  delta: '\u504f\u5dee',
  second: '\u79d2',
  minute: '\u5206',
  hour: '\u65f6',
  day: '\u5929',
}

const inputReg = computed<any>(() => props.info?.Packet?.Femto_input_reg ?? {})
const holdingReg = computed<any>(() => props.info?.Packet?.Femto_holding_reg ?? {})
const laserPara = computed<any>(() => holdingReg.value?.Laser_para ?? {})
const userPara = computed<any>(() => holdingReg.value?.User_para ?? {})
const laserInfo = computed<any>(() => laserPara.value?.Laser_info ?? {})
const espModule = computed<any>(() => laserPara.value?.Esp_module ?? {})
const espMon = computed<any>(() => inputReg.value?.Esp32 ?? {})
const monReg = computed<any>(() => inputReg.value?.Mon ?? {})
const timeReg = computed<any>(() => inputReg.value?.Time ?? {})
const alarmReg = computed<any>(() => inputReg.value?.Alarm ?? {})

const LASER_STATUS_TEXT: Record<number, string> = {
  0: t.preheat,
  1: t.standby,
  2: t.ready,
  3: t.fault,
}

const READY_STATUS_TEXT: Record<number, string> = {
  0: t.notReady,
  1: t.warming,
  2: t.readyDone,
}

const WAVELENGTH_TEXT: Record<number, string> = {
  0: t.ir,
  1: t.green,
  2: t.uv,
}

const ALARM_PREVIEW_LIMIT = 24

const laserStatusText = computed(
  () => LASER_STATUS_TEXT[toSafeNumber(inputReg.value?.Status)] ?? String(toSafeNumber(inputReg.value?.Status))
)
const readyStatusText = computed(
  () => READY_STATUS_TEXT[toSafeNumber(userPara.value?.Laser_ready)] ?? String(toSafeNumber(userPara.value?.Laser_ready))
)
const wavelengthText = computed(
  () => WAVELENGTH_TEXT[toSafeNumber(userPara.value?.Laser_wavelength)] ?? String(toSafeNumber(userPara.value?.Laser_wavelength))
)

const modelText = computed(() => decodeAsciiByteArray(laserInfo.value?.Model))
const snText = computed(() => decodeAsciiByteArray(laserInfo.value?.SN))
const pnText = computed(() => decodeAsciiByteArray(laserInfo.value?.PN))
const ssidText = computed(() => decodeAsciiWordArray(espModule.value?.Ssid))

const tempValues = computed<number[]>(() => numberArray(monReg.value?.Temp))
const voltageValues = computed<number[]>(() => numberArray(monReg.value?.Vol))

const pumpRows = computed(() => {
  const pumps = Array.isArray(monReg.value?.Pump_mon) ? monReg.value.Pump_mon : []
  return pumps.map((item: any, index: number) => {
    const actual = toSafeNumber(item?.Actual_cur)
    const fpga = toSafeNumber(item?.Fpga_cur)
    return {
      index,
      switch: toSafeNumber(item?.Pump_sw) > 0 ? t.switchOn : t.switchOff,
      actual,
      fpga,
      delta: actual - fpga,
    }
  })
})

const pumpColumns = [
  { title: t.channel, dataIndex: 'index', width: 80, customRender: ({ text }: { text: number }) => `CH${text}` },
  { title: t.state, dataIndex: 'switch', width: 90 },
  { title: t.actualCurrent, dataIndex: 'actual', width: 120 },
  { title: t.fpgaCurrent, dataIndex: 'fpga', width: 120 },
  { title: t.delta, dataIndex: 'delta', width: 100 },
]

const nowAlarmBits = computed<number[]>(() => decodeAlarmBits(alarmReg.value?.Now))
const historyAlarmBits = computed<number[]>(() => decodeAlarmBits(alarmReg.value?.History))

const nowAlarmBitsPreview = computed(() => nowAlarmBits.value.slice(0, ALARM_PREVIEW_LIMIT))
const nowAlarmBitsRemain = computed(() => Math.max(0, nowAlarmBits.value.length - ALARM_PREVIEW_LIMIT))
const historyAlarmBitsPreview = computed(() => historyAlarmBits.value.slice(0, ALARM_PREVIEW_LIMIT))
const historyAlarmBitsRemain = computed(() => Math.max(0, historyAlarmBits.value.length - ALARM_PREVIEW_LIMIT))

function toSafeNumber(value: unknown): number {
  const n = Number(value ?? 0)
  return Number.isFinite(n) ? n : 0
}

function numberArray(value: unknown): number[] {
  if (!Array.isArray(value)) {
    return []
  }
  return value.map((item) => toSafeNumber(item))
}

function toPrintableAsciiChar(code: number): string {
  if (code < 32 || code > 126) {
    return ''
  }
  return String.fromCharCode(code)
}

function decodeAsciiByteArray(value: unknown): string {
  if (!Array.isArray(value)) {
    return '--'
  }
  const chars: string[] = []
  for (const item of value) {
    const code = toSafeNumber(item) & 0xff
    if (code === 0) {
      break
    }
    const ch = toPrintableAsciiChar(code)
    if (ch) {
      chars.push(ch)
    }
  }
  return chars.length ? chars.join('') : '--'
}

function decodeAsciiWordArray(value: unknown): string {
  if (!Array.isArray(value)) {
    return '--'
  }
  const chars: string[] = []
  for (const item of value) {
    const word = toSafeNumber(item) & 0xffff
    const low = word & 0xff
    const high = (word >> 8) & 0xff
    if (low === 0) {
      break
    }
    const lowChar = toPrintableAsciiChar(low)
    if (lowChar) {
      chars.push(lowChar)
    }
    if (high === 0) {
      continue
    }
    const highChar = toPrintableAsciiChar(high)
    if (highChar) {
      chars.push(highChar)
    }
  }
  return chars.length ? chars.join('') : '--'
}

function formatIPv4(value: unknown): string {
  if (!Array.isArray(value) || value.length < 4) {
    return '--'
  }
  const octets = value.slice(0, 4).map((item) => {
    const n = toSafeNumber(item)
    return Math.max(0, Math.min(255, n))
  })
  if (octets.every((v) => v === 0)) {
    return '--'
  }
  return octets.join('.')
}

function joinVersion(value: unknown): string {
  if (!Array.isArray(value)) {
    return '--'
  }
  const vals = value.map((item) => toSafeNumber(item))
  return vals.some((v) => v > 0) ? vals.join('.') : '--'
}

function numberOrDash(value: unknown): string | number {
  const n = Number(value)
  if (!Number.isFinite(n)) {
    return '--'
  }
  return n
}

function switchText(value: unknown): string {
  return toSafeNumber(value) > 0 ? t.switchOn : t.switchOff
}

function mergeU16Pair(value: unknown): number {
  if (!Array.isArray(value) || value.length < 2) {
    return 0
  }
  const high = toSafeNumber(value[0]) & 0xffff
  const low = toSafeNumber(value[1]) & 0xffff
  return ((high << 16) >>> 0) + low
}

function formatSeconds(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds <= 0) {
    return `0${t.second}`
  }
  const total = Math.floor(seconds)
  const day = Math.floor(total / 86400)
  const hour = Math.floor((total % 86400) / 3600)
  const minute = Math.floor((total % 3600) / 60)
  const second = total % 60
  if (day > 0) {
    return `${day}${t.day} ${hour}${t.hour} ${minute}${t.minute}`
  }
  if (hour > 0) {
    return `${hour}${t.hour} ${minute}${t.minute} ${second}${t.second}`
  }
  if (minute > 0) {
    return `${minute}${t.minute} ${second}${t.second}`
  }
  return `${second}${t.second}`
}

function formatU16PairDuration(value: unknown): string {
  return formatSeconds(mergeU16Pair(value))
}

function formatTimeDesc(value: any): string {
  if (!value || typeof value !== 'object') {
    return '--'
  }
  const year = 2000 + toSafeNumber(value.Year)
  const mon = toSafeNumber(value.Mon)
  const day = toSafeNumber(value.Day)
  const hour = toSafeNumber(value.Hour)
  const minute = toSafeNumber(value.Minute ?? value.Minutes)
  const second = toSafeNumber(value.Second)

  if (year === 2000 && mon === 0 && day === 0 && hour === 0 && minute === 0 && second === 0) {
    return '--'
  }

  const pad = (n: number) => String(Math.max(0, n)).padStart(2, '0')
  return `${year}-${pad(mon)}-${pad(day)} ${pad(hour)}:${pad(minute)}:${pad(second)}`
}

function decodeAlarmBits(value: unknown): number[] {
  if (!Array.isArray(value)) {
    return []
  }
  const result: number[] = []
  value.forEach((wordRaw, wordIndex) => {
    const word = toSafeNumber(wordRaw) & 0xffff
    for (let bit = 0; bit < 16; bit += 1) {
      if (((word >> bit) & 0x1) === 1) {
        result.push(wordIndex * 16 + bit)
      }
    }
  })
  return result
}

function formatHexWords(value: unknown): string {
  if (!Array.isArray(value)) {
    return '--'
  }
  const words = value.map((item) => (toSafeNumber(item) & 0xffff).toString(16).padStart(4, '0'))
  return words.length ? words.map((w) => `0x${w}`).join(' ') : '--'
}
</script>

<style scoped>
.device-insight {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.top-cards :deep(.ant-card-body) {
  padding: 10px;
}

.pump-card :deep(.ant-card-body) {
  padding-top: 8px;
}

.alarm-block + .alarm-block {
  margin-top: 12px;
}

.alarm-title {
  color: #47586a;
  font-size: 12px;
  margin-bottom: 8px;
}

.alarm-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
</style>
