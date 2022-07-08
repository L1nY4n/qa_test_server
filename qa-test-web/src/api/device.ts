import { Device } from '@/types/api'
import axios from './axios'

export async function list(): Promise<Device[]> {
    return await axios.get("/device/list")
}
