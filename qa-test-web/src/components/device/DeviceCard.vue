<template>
  <a-card class="device-card" :class="online ? 'online-card' : 'offline-card'" hoverable @click="openDetail">
    <div class="card-head">
      <div class="sn">{{ info.Sn || '-' }}</div>
      <div class="status-pill" :class="online ? 'online' : 'offline'">
        <span class="status-lamp" />
        <span class="status-text">{{ online ? t.online : t.offline }}</span>
      </div>
    </div>

    <div class="name-row">
      <div class="name">{{ info.Name || t.unnamedDevice }}</div>
      <a-tag v-if="groupText" color="blue">{{ groupText }}</a-tag>
    </div>

    <div class="meta-row">
      <span class="meta-label">{{ t.modelLabel }}</span>
      <span class="meta-value">{{ info.Model || '--' }}</span>
    </div>

    <div class="status-row">
      <a-tag color="cyan">{{ t.statusLabel }} {{ laserStatusText }}</a-tag>
      <a-tag color="purple">{{ t.readyLabel }} {{ laserReadyText }}</a-tag>
      <a-tag color="geekblue">{{ t.wavelengthLabel }} {{ wavelengthText }}</a-tag>
    </div>

    <div class="metrics">
      <div class="metric">
        <span class="label">{{ t.hardwareVersion }}</span>
        <span class="value">{{ hardwareVersion }}</span>
      </div>
      <div class="metric">
        <span class="label">{{ t.uptime }}</span>
        <span class="value">{{ uptimeText }}</span>
      </div>
      <div class="metric">
        <span class="label">{{ t.pumpChannel }}</span>
        <span class="value">{{ pumpCount }}</span>
      </div>
      <div class="metric">
        <span class="label">{{ t.tempChannel }}</span>
        <span class="value">{{ tempCount }}</span>
      </div>
    </div>

    <div class="footer">{{ t.updatedAt }} {{ lastSeenText }}</div>
  </a-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import moment from 'moment'
import type { DeviceCardInfo } from '@/api/device'

const props = withDefaults(
  defineProps<{
    info: DeviceCardInfo
    activeWithin?: number
  }>(),
  {
    activeWithin: 30,
  }
)

const emit = defineEmits<{
  (e: 'open', sn: string): void
}>()

const t = {
  online: '\u5728\u7ebf',
  offline: '\u79bb\u7ebf',
  unnamedDevice: '\u672a\u547d\u540d\u8bbe\u5907',
  modelLabel: '\u578b\u53f7:',
  statusLabel: '\u72b6\u6001',
  readyLabel: '\u5c31\u7eea',
  wavelengthLabel: '\u6ce2\u957f',
  hardwareVersion: '\u786c\u4ef6\u7248\u672c',
  uptime: '\u8fd0\u884c\u65f6\u957f',
  pumpChannel: '\u6cf5\u901a\u9053',
  tempChannel: '\u6e29\u5ea6\u901a\u9053',
  updatedAt: '\u66f4\u65b0\u4e8e',
  stressGroup: '\u538b\u529b\u6d4b\u8bd5\u7ec4',
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
  second: '\u79d2',
  minute: '\u5206',
  hour: '\u65f6',
  day: '\u5929',
}

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

const hardwareVersion = computed(() => props.info?.Hardware_bate ?? '--')
const pumpCount = computed(() => Number(props.info?.Pump_count ?? 0))
const tempCount = computed(() => Number(props.info?.Temp_count ?? 0))

const laserStatusText = computed(
  () => LASER_STATUS_TEXT[Number(props.info?.Laser_status ?? -1)] ?? String(props.info?.Laser_status ?? '--')
)
const laserReadyText = computed(
  () => READY_STATUS_TEXT[Number(props.info?.Laser_ready ?? -1)] ?? String(props.info?.Laser_ready ?? '--')
)
const wavelengthText = computed(
  () => WAVELENGTH_TEXT[Number(props.info?.Laser_wavelength ?? -1)] ?? String(props.info?.Laser_wavelength ?? '--')
)

const groupText = computed(() => {
  if (!props.info?.Group) {
    return ''
  }
  if (props.info.Group === 'virtual-stress') {
    return t.stressGroup
  }
  return props.info.Group
})

const uptimeText = computed(() => {
  const sec = Number(props.info?.Uptime_seconds ?? 0)
  if (Number.isFinite(sec) && sec > 0) {
    return formatSeconds(sec)
  }
  const up = props.info?.Uptime
  if (!Array.isArray(up) || up.length < 2) {
    return '--'
  }
  const merged = (Number(up[0] ?? 0) << 16) + Number(up[1] ?? 0)
  return formatSeconds(merged)
})

const online = computed(() => {
  if (typeof props.info?.Online === 'boolean') {
    return props.info.Online
  }
  const ts = props.info?.Last_rx_time
  if (!ts) return false
  const tms = new Date(ts).getTime()
  if (!Number.isFinite(tms)) return false
  return Date.now() - tms <= props.activeWithin * 1000
})

const lastSeenText = computed(() => {
  const ts = props.info?.Last_rx_time
  if (!ts) return '--'
  const m = moment(ts)
  if (!m.isValid()) return ts
  return m.fromNow()
})

const formatSeconds = (seconds: number): string => {
  if (!Number.isFinite(seconds) || seconds <= 0) {
    return `0${t.second}`
  }
  const s = Math.floor(seconds)
  const day = Math.floor(s / 86400)
  const hour = Math.floor((s % 86400) / 3600)
  const minute = Math.floor((s % 3600) / 60)
  const sec = s % 60
  if (day > 0) {
    return `${day}${t.day} ${hour}${t.hour} ${minute}${t.minute}`
  }
  if (hour > 0) {
    return `${hour}${t.hour} ${minute}${t.minute} ${sec}${t.second}`
  }
  if (minute > 0) {
    return `${minute}${t.minute} ${sec}${t.second}`
  }
  return `${sec}${t.second}`
}

const openDetail = () => {
  if (!props.info?.Sn) {
    return
  }
  emit('open', props.info.Sn)
}
</script>

<style scoped>
.device-card {
  position: relative;
  overflow: hidden;
  border-radius: 12px;
  border: 1px solid #e1e9f3;
  background: linear-gradient(165deg, #f8fbff, #f3f8ff 55%, #eef6ff);
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.device-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 24px rgba(28, 56, 85, 0.12);
}

.online-card {
  border-color: #84d8b6;
}

.online-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 12px;
  box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.35);
  animation: card-pulse 2.4s ease-out infinite;
  pointer-events: none;
}

.offline-card {
  border-color: #d9e3ef;
  background: linear-gradient(165deg, #f7f9fc, #f3f6fa 55%, #edf2f8);
}

.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 10px;
  border-radius: 999px;
  border: 1px solid #d6e2ef;
  font-size: 12px;
  font-weight: 600;
}

.status-pill.online {
  color: #0e7a48;
  border-color: #93e3bf;
  background: linear-gradient(145deg, #e8fff3, #def9ed);
}

.status-pill.offline {
  color: #687f96;
  border-color: #d8e1eb;
  background: linear-gradient(145deg, #f7f9fc, #f1f4f8);
}

.status-lamp {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  flex: 0 0 9px;
}

.status-pill.online .status-lamp {
  background: #34d399;
  box-shadow: 0 0 0 0 rgba(52, 211, 153, 0.65);
  animation: lamp-breath 1.8s ease-in-out infinite;
}

.status-pill.offline .status-lamp {
  background: #94a3b8;
  box-shadow: 0 0 0 0 rgba(148, 163, 184, 0.2);
}

.status-text {
  line-height: 1;
}

.sn {
  color: #1677ff;
  font-size: 15px;
  font-weight: 700;
  line-height: 1.3;
  word-break: break-all;
}

.name-row {
  margin-top: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.name {
  color: #4d5f73;
  font-size: 13px;
  line-height: 1.3;
  min-height: 18px;
}

.meta-row {
  margin-top: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-label {
  color: #71879d;
  font-size: 12px;
}

.meta-value {
  color: #24384d;
  font-size: 12px;
  font-weight: 600;
}

.status-row {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.metrics {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.metric {
  border-radius: 8px;
  border: 1px solid #edf2f8;
  background: #f8fbff;
  padding: 8px;
}

.label {
  display: block;
  color: #71879d;
  font-size: 12px;
}

.value {
  display: block;
  margin-top: 2px;
  color: #1f3347;
  font-size: 13px;
  font-weight: 600;
}

.footer {
  margin-top: 10px;
  color: #7f93a7;
  font-size: 12px;
}

@keyframes lamp-breath {
  0% {
    transform: scale(0.9);
    box-shadow: 0 0 0 0 rgba(52, 211, 153, 0.55);
    opacity: 0.75;
  }
  50% {
    transform: scale(1.12);
    box-shadow: 0 0 0 8px rgba(52, 211, 153, 0);
    opacity: 1;
  }
  100% {
    transform: scale(0.9);
    box-shadow: 0 0 0 0 rgba(52, 211, 153, 0);
    opacity: 0.75;
  }
}

@keyframes card-pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.3);
    opacity: 0.35;
  }
  70% {
    box-shadow: 0 0 0 9px rgba(34, 197, 94, 0);
    opacity: 0;
  }
  100% {
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0);
    opacity: 0;
  }
}
</style>
