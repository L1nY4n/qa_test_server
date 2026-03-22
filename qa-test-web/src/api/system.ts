import axios from './axios'

export interface SystemHealth {
    uptimeSeconds: number
    db: {
        ok: boolean
        message: string
    }
    websocket: Record<string, unknown>
    server: {
        httpAddr: string
        tcpAddr: string
    }
    performance?: SystemMetricSample
}

export async function health(): Promise<SystemHealth> {
    return await axios.get("/system/health")
}

export interface SystemHostInfo {
    hostname: string
    os: string
    platform: string
    platformVersion: string
    kernelVersion: string
    architecture: string
    bootTime: number
}

export interface SystemMetricSample {
    sampledAt: string
    cpuPercent: number
    load1?: number | null
    load5?: number | null
    load15?: number | null
    memoryTotalBytes: number
    memoryUsedBytes: number
    memoryUsedPercent: number
    swapTotalBytes: number
    swapUsedBytes: number
    swapUsedPercent: number
    diskTotalBytes: number
    diskUsedBytes: number
    diskUsedPercent: number
    netSentBytes: number
    netRecvBytes: number
    netSentBps: number
    netRecvBps: number
    processCpuPercent: number
    processRssBytes: number
    goroutines: number
}

export interface SystemMetricsPayload {
    sampleEverySeconds: number
    diskPath: string
    host: SystemHostInfo
    current: SystemMetricSample
    points: SystemMetricSample[]
}

export async function metrics(points = 180): Promise<SystemMetricsPayload> {
    return await axios.get("/system/metrics", {
        params: { points },
    })
}

export interface SystemResetResult {
    resetAt: string
    clearedDevices: number
    historyRowsDeleted: number
    paramRowsDeleted: number
    virtualDevicesStopped: boolean
}

export async function resetCloudSystem(): Promise<SystemResetResult> {
    return await axios.post("/system/reset")
}
