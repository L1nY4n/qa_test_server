import axios from './axios'
import type { CurrentUser, UserRole } from '@/utils/auth'

export interface UserListQuery {
  keyword?: string
  offset?: number
  limit?: number
}

export interface UserListMeta {
  items: CurrentUser[]
  total: number
  offset: number
  limit: number
}

export interface CreateUserRequest {
  username: string
  displayName?: string
  password: string
  role: UserRole
}

export interface UpdateUserRequest {
  displayName?: string
  role?: UserRole
  enabled?: boolean
}

export async function list(params: UserListQuery = {}): Promise<UserListMeta> {
  return await axios.get('/users/list', { params })
}

export async function create(data: CreateUserRequest): Promise<CurrentUser> {
  return await axios.post('/users', data)
}

export async function update(id: number, data: UpdateUserRequest): Promise<CurrentUser> {
  return await axios.put(`/users/${id}`, data)
}

export async function resetPassword(id: number, password: string): Promise<{ id: number; message: string }> {
  return await axios.put(`/users/${id}/password`, { password })
}

export async function remove(id: number): Promise<{ id: number; message: string }> {
  return await axios.delete(`/users/${id}`)
}
