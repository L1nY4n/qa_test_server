import axios from './axios'
import type { CurrentUser } from '@/utils/auth'

export interface AuthPayload {
  token: string
  expiresAt: string
  user: CurrentUser
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  displayName?: string
  password: string
}

export interface ProfilePayload {
  user: CurrentUser
}

export async function login(data: LoginRequest): Promise<AuthPayload> {
  return await axios.post('/auth/login', data)
}

export async function register(data: RegisterRequest): Promise<AuthPayload> {
  return await axios.post('/auth/register', data)
}

export async function profile(): Promise<ProfilePayload> {
  return await axios.get('/auth/profile')
}
