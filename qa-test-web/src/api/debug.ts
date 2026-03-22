import axios from './axios'

export interface VirtualDeviceConfig {
  count?: number
  intervalMs?: number
  prefix?: string
  namePrefix?: string
  startIndex?: number
  group?: string
  mutateParam?: boolean
  wsBroadcast?: boolean
  pulseRepeat?: number
}

export interface VirtualDeviceStatus {
  running: boolean
  count: number
  intervalMs: number
  prefix: string
  namePrefix: string
  startIndex: number
  group: string
  mutateParam: boolean
  wsBroadcast: boolean
  startedAt: string
  lastTickAt: string
  updatesGenerated: number
  updatesPerSecond: number
  broadcastDropped: number
  activeDevices: number
  sampleSn: string
}

export interface VirtualPulseResult {
  generated: number
  elapsed: number
  elapsedMs: number
  elapsedNs?: number
  updatesPerSecond: number
  group?: string
  cleaned?: boolean
}

export async function virtualStatus(): Promise<VirtualDeviceStatus> {
  return await axios.get('/debug/virtual/status')
}

export async function virtualStart(payload: VirtualDeviceConfig): Promise<VirtualDeviceStatus> {
  return await axios.post('/debug/virtual/start', payload)
}

export async function virtualStop(remove = true): Promise<VirtualDeviceStatus> {
  return await axios.post('/debug/virtual/stop', { remove })
}

export async function virtualPulse(payload: VirtualDeviceConfig): Promise<VirtualPulseResult> {
  return await axios.post('/debug/virtual/pulse', payload)
}

export async function virtualStressPulse(payload: VirtualDeviceConfig): Promise<VirtualPulseResult> {
  return await axios.post('/debug/virtual/stress/pulse', payload)
}

export interface TimeKeyDecryptRequest {
  deviceSn: string
  key?: string
  keyHex?: string
}

export interface TimeKeyDecryptResult {
  deviceSn: string
  valid: boolean
  decodedYear: number
  decodedMonth: number
  decodedDay: number
  decodedHour: number
  decodedMinute: number
  decodedSecond: number
  fullYear: number
  decodedAt: string
  decodedAtText: string
  rawHead: string
  rawTail: string
}

export interface TimeKeyGenerateRequest {
  deviceSn: string
  time: string
}

export interface TimeKeyGenerateResult {
  deviceSn: string
  inputTime: string
  inputTimeText: string
  decodedYear: number
  decodedMonth: number
  decodedDay: number
  decodedHour: number
  decodedMinute: number
  decodedSecond: number
  fullYear: number
  key: string
  keyHex: string
  rawHead: string
  rawTail: string
}

export interface DecryptLogItem {
  id: number
  operatorId: number
  operatorName: string
  operatorRole: string
  operation: string
  deviceSn: string
  inputMode: string
  keyPreview: string
  keyHash: string
  success: boolean
  errorMessage: string
  decodedYear: number
  decodedMonth: number
  decodedDay: number
  decodedHour: number
  decodedMinute: number
  decodedSecond: number
  sourceIp: string
  userAgent: string
  createdAt: string
}

export interface DecryptLogPayload {
  items: DecryptLogItem[]
  total: number
  offset: number
  limit: number
  deviceSn: string
  keyword: string
}

export interface DecryptLogQuery {
  deviceSn?: string
  keyword?: string
  operation?: 'decode' | 'generate'
  offset?: number
  limit?: number
}

export async function decryptTimeKey(payload: TimeKeyDecryptRequest): Promise<TimeKeyDecryptResult> {
  return await axios.post('/debug/decrypt/time-key', payload)
}

export async function generateTimeKey(payload: TimeKeyGenerateRequest): Promise<TimeKeyGenerateResult> {
  return await axios.post('/debug/decrypt/time-key/generate', payload)
}

export async function decryptLogs(params: DecryptLogQuery = {}): Promise<DecryptLogPayload> {
  return await axios.get('/debug/decrypt/logs', { params })
}
