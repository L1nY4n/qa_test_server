import { Device } from '@/types/api'
import axios from './axios'
import rawAxios from 'axios'
import { base_url } from '@/config'
import { getToken } from '@/utils/auth'

export interface DeviceListQuery {
    keyword?: string
    group?: string
    onlineOnly?: boolean
    offset?: number
    limit?: number
    activeWithin?: number
    summary?: boolean
}

export interface DeviceStats {
    total: number
    online: number
    offline: number
    activeWithinSeconds: number
    generatedAt: string
}

export interface DeviceCardInfo {
    Sn: string
    Name: string
    Group: string
    Model: string
    PN: string
    Last_rx_time: string
    Hardware_bate: number
    Uptime: number[]
    Uptime_seconds: number
    Pump_count: number
    Temp_count: number
    Laser_status: number
    Laser_ready: number
    Laser_wavelength: number
    Online: boolean
}

export interface DeviceInfoPayload {
    device: Device
    online: boolean
}

export interface DeviceHistoryQuery {
    start?: string | number
    end?: string | number
    offset?: number
    limit?: number
    raw?: boolean
}

export interface DeviceTimelineQuery {
    start?: string | number
    end?: string | number
    offset?: number
    limit?: number
    maxPoints?: number
}

export interface DeviceMetricTimelineQuery {
    start?: string | number
    end?: string | number
    offset?: number
    limit?: number
    maxPoints?: number
    tempIndex?: number
    voltageIndex?: number
    currentIndex?: number
}

export interface DeviceHistoryPoint {
    sampledAt: string
    online: boolean
    hardwareBate: number
    uptimeSeconds: number
    pumpCount: number
    tempCount: number
}

export interface DeviceHistoryRecord extends DeviceHistoryPoint {
    id: number
    deviceSn: string
    deviceName: string
    rawJson: string
    createdAt: string
}

export interface DeviceHistoryPayload<T = DeviceHistoryPoint> {
    items: T[]
    total: number
    offset: number
    limit: number
    raw: boolean
}

export interface DeviceHistoryTimelinePoint {
    sampledAt: string
    online: boolean
    tempAvg: number
    tempMax: number
    voltageAvg: number
    voltageMax: number
    currentAvg: number
    currentMax: number
}

export interface DeviceHistoryTimelinePayload {
    items: DeviceHistoryTimelinePoint[]
    total: number
    offset: number
    limit: number
    maxPoints: number
}

export interface DeviceHistoryMetricPoint {
    sampledAt: string
    online: boolean
    temp: number | null
    voltage: number | null
    current: number | null
}

export interface DeviceHistoryMetricPayload {
    items: DeviceHistoryMetricPoint[]
    total: number
    offset: number
    limit: number
    maxPoints: number
    tempIndex: number
    voltageIndex: number
    currentIndex: number
}

export interface DeviceParamChangeQuery {
    start?: string | number
    end?: string | number
    path?: string
    offset?: number
    limit?: number
}

export interface DeviceParamChange {
    id: number
    deviceSn: string
    paramPath: string
    oldValue: string
    newValue: string
    changedAt: string
    createdAt: string
}

export interface DeviceParamChangePayload {
    items: DeviceParamChange[]
    total: number
    offset: number
    limit: number
    path: string
}

export interface DeviceListMeta<T = Device> {
    items: T[]
    total: number
    offset: number
    limit: number
    group?: string
    onlineOnly: boolean
    activeWithinSeconds: number
}

export async function list(params: DeviceListQuery = {}): Promise<Device[]> {
    return await axios.get("/device/list", { params })
}

export async function listWithMeta<T = Device>(params: DeviceListQuery = {}): Promise<DeviceListMeta<T>> {
    return await axios.get("/device/list", {
        params: {
            ...params,
            meta: true,
        },
    })
}

export async function info(sn: string, activeWithin = 30): Promise<DeviceInfoPayload> {
    return await axios.get(`/device/info/${encodeURIComponent(sn)}`, {
        params: { activeWithin },
    })
}

export async function stats(activeWithin = 30): Promise<DeviceStats> {
    return await axios.get("/device/stats", {
        params: { activeWithin },
    })
}

export async function history<T = DeviceHistoryPoint>(sn: string, params: DeviceHistoryQuery = {}): Promise<DeviceHistoryPayload<T>> {
    return await axios.get(`/device/history/${encodeURIComponent(sn)}`, {
        params,
    })
}

export async function historyTimeline(sn: string, params: DeviceTimelineQuery = {}): Promise<DeviceHistoryTimelinePayload> {
    return await axios.get(`/device/history/${encodeURIComponent(sn)}/timeline`, {
        params,
    })
}

export async function historyMetrics(sn: string, params: DeviceMetricTimelineQuery = {}): Promise<DeviceHistoryMetricPayload> {
    return await axios.get(`/device/history/${encodeURIComponent(sn)}/metrics`, {
        params,
    })
}

export async function changeLogs(sn: string, params: DeviceParamChangeQuery = {}): Promise<DeviceParamChangePayload> {
    return await axios.get(`/device/changes/${encodeURIComponent(sn)}`, {
        params,
    })
}

export async function exportChangeLogs(sn: string, params: DeviceParamChangeQuery = {}): Promise<Blob> {
    const token = getToken()
    const response = await rawAxios.get(`${base_url}/device/changes/${encodeURIComponent(sn)}/export`, {
        params,
        responseType: 'blob',
        headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
    return response.data as Blob
}
