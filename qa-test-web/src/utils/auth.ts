export type UserRole = 'admin' | 'operator' | 'viewer'

export interface CurrentUser {
  id: number
  username: string
  displayName: string
  role: UserRole
  enabled: boolean
  createdAt: string
  updatedAt: string
}

const TOKEN_KEY = 'qa_auth_token'
const USER_KEY = 'qa_auth_user'
const EXPIRES_AT_KEY = 'qa_auth_expires_at'

const rolePriority: Record<UserRole, number> = {
  admin: 3,
  operator: 2,
  viewer: 1,
}

const parseJSON = <T>(raw: string | null): T | null => {
  if (!raw) {
    return null
  }
  try {
    return JSON.parse(raw) as T
  } catch {
    return null
  }
}

export const clearSession = (): void => {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  localStorage.removeItem(EXPIRES_AT_KEY)
}

export const saveSession = (token: string, user: CurrentUser, expiresAt?: string): void => {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_KEY, JSON.stringify(user))
  if (expiresAt) {
    localStorage.setItem(EXPIRES_AT_KEY, expiresAt)
  } else {
    localStorage.removeItem(EXPIRES_AT_KEY)
  }
}

export const getToken = (): string => {
  const token = localStorage.getItem(TOKEN_KEY) || ''
  if (!token) {
    return ''
  }

  const expiresAt = localStorage.getItem(EXPIRES_AT_KEY)
  if (!expiresAt) {
    return token
  }
  const expires = new Date(expiresAt).getTime()
  if (!Number.isFinite(expires) || Date.now() < expires) {
    return token
  }

  clearSession()
  return ''
}

export const getCurrentUser = (): CurrentUser | null => {
  return parseJSON<CurrentUser>(localStorage.getItem(USER_KEY))
}

export const setCurrentUser = (user: CurrentUser): void => {
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export const hasAnyRole = (required: string[] | undefined, role: string | undefined): boolean => {
  if (!required || required.length === 0) {
    return true
  }
  if (!role) {
    return false
  }

  const normalizedRole = role.toLowerCase() as UserRole
  const score = rolePriority[normalizedRole]
  if (!score) {
    return false
  }

  return required.some((item) => {
    const requiredRole = item.toLowerCase() as UserRole
    const requiredScore = rolePriority[requiredRole]
    if (!requiredScore) {
      return false
    }
    return score >= requiredScore
  })
}
