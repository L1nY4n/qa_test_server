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
}

export async function health(): Promise<SystemHealth> {
    return await axios.get("/system/health")
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
