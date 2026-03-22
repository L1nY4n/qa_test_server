import axios from './axios'
import rawAxios from 'axios'
import { base_url } from '@/config'
import { getToken } from '@/utils/auth'

export interface BenchmarkTopDevice {
  deviceSn: string
  deviceName: string
  samples: number
  onlineRatio: number
  completeness: number
  restartCount: number
  paramChangeCount: number
  hardwareBate: number
  riskScore: number
  riskLevel: string
}

export interface BenchmarkFirmwareDistribution {
  hardwareBate: number
  count: number
}

export interface BenchmarkOfflineHour {
  hour: number
  totalCount: number
  offlineCount: number
  offlineRate: number
}

export interface BenchmarkRiskMatrixPoint {
  deviceSn: string
  deviceName: string
  onlineRatioPct: number
  completenessPct: number
  paramChangeCount: number
  riskScore: number
  riskLevel: string
}

export interface BenchmarkAuditSummary {
  operation: string
  success: number
  failed: number
}

export interface BenchmarkMaintenanceCandidate {
  deviceSn: string
  deviceName: string
  onlineRatio: number
  completeness: number
  restartCount: number
  paramChangeCnt: number
  riskScore: number
  priority: 'high' | 'medium' | 'low'
  reason: string
}

export interface BenchmarkFeatureItem {
  code: string
  name: string
  description: string
  status: string
}

export interface BenchmarkInsightPayload {
  windowHours: number
  generatedAt: string
  slaTop: BenchmarkTopDevice[]
  completeTop: BenchmarkTopDevice[]
  restartTop: BenchmarkTopDevice[]
  changeTop: BenchmarkTopDevice[]
  firmware: BenchmarkFirmwareDistribution[]
  offlineHeat: BenchmarkOfflineHour[]
  riskMatrix: BenchmarkRiskMatrixPoint[]
  audit: BenchmarkAuditSummary[]
  candidates: BenchmarkMaintenanceCandidate[]
  features: BenchmarkFeatureItem[]
}

export async function insights(windowHours = 24): Promise<BenchmarkInsightPayload> {
  return await axios.get('/benchmark/insights', {
    params: { windowHours },
  })
}

export async function exportReport(windowHours = 24): Promise<Blob> {
  const token = getToken()
  const response = await rawAxios.get(`${base_url}/benchmark/export`, {
    params: { windowHours },
    responseType: 'blob',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  })
  return response.data as Blob
}
