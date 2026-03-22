import axios from './axios'
import type { DeviceHistoryPoint, DeviceParamChange } from './device'

export type PHMRiskLevel = 'low' | 'medium' | 'high' | 'critical'

export interface PHMSummary {
  deviceSn: string
  deviceName: string
  windowHours: number
  samples: number
  onlineRatio: number
  lastSampledAt: string
  freshnessMinutes: number
  rebootCount: number
  hardwareShift: number
  dataGapCount: number
  paramChangeCount: number
  healthScore: number
  riskLevel: PHMRiskLevel
  reasons: string[]
  recommendation: string
}

export interface PHMOverviewPayload {
  windowHours: number
  total: number
  riskStats: {
    low: number
    medium: number
    high: number
    critical: number
  }
  items: PHMSummary[]
  generatedAt: string
}

export interface PHMDeviceDetail {
  summary: PHMSummary
  historyPoints: DeviceHistoryPoint[]
  recentChanges: DeviceParamChange[]
  windowStart: string
  windowEnd: string
}

export interface PHMDeleteResult {
  requested: number
  cleared: string[]
  missing: string[]
  cacheRemoved: number
  historyRowsDeleted: number
  changeRowsDeleted: number
  deletedAt: string
}

export async function overview(windowHours = 24, limit = 100): Promise<PHMOverviewPayload> {
  return await axios.get('/phm/overview', {
    params: { windowHours, limit },
  })
}

export async function deviceDetail(sn: string, windowHours = 24, historyLimit = 600, changeLimit = 200): Promise<PHMDeviceDetail> {
  return await axios.get(`/phm/device/${encodeURIComponent(sn)}`, {
    params: { windowHours, historyLimit, changeLimit },
  })
}

export async function deleteDevice(sn: string): Promise<PHMDeleteResult> {
  return await axios.delete(`/phm/device/${encodeURIComponent(sn)}`)
}

export async function deleteDevices(deviceSns: string[]): Promise<PHMDeleteResult> {
  return await axios.post('/phm/devices/delete', { deviceSns })
}
