import { createRouter, createWebHashHistory } from 'vue-router'
import routes, { menus } from './routes'
import * as API from '@/api'
import { clearSession, getCurrentUser, getToken, hasAnyRole, setCurrentUser } from '@/utils/auth'

const history = createWebHashHistory()
const APP_TITLE = '\u6fc0\u5149\u5668\u7ba1\u7406\u7cfb\u7edf'

const router = createRouter({ history, routes })

let profileLoaded = false
let profilePromise: Promise<boolean> | null = null

const fallbackRouteByRole = (role: string | undefined): string => {
  for (const group of menus) {
    const children = group.children || []
    for (const item of children) {
      if (hasAnyRole(item.meta?.roles, role)) {
        return String(item.name || 'overview')
      }
    }
  }
  return 'overview'
}

const ensureProfile = async (): Promise<boolean> => {
  if (profileLoaded) {
    return true
  }
  if (profilePromise) {
    return profilePromise
  }

  profilePromise = API.auth
    .profile()
    .then((res) => {
      setCurrentUser(res.user)
      profileLoaded = true
      return true
    })
    .catch(() => {
      clearSession()
      profileLoaded = false
      return false
    })
    .finally(() => {
      profilePromise = null
    })

  return profilePromise
}

router.beforeEach(async (to) => {
  const isPublic = Boolean(to.meta.public)
  const token = getToken()

  if (!token) {
    profileLoaded = false
    if (!isPublic) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    return true
  }

  if (to.name === 'login') {
    const user = getCurrentUser()
    return { name: fallbackRouteByRole(user?.role) }
  }

  const profileOk = await ensureProfile()
  if (!profileOk) {
    if (!isPublic) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    return true
  }

  const user = getCurrentUser()
  if (!hasAnyRole(to.meta.roles, user?.role)) {
    return { name: fallbackRouteByRole(user?.role) }
  }

  return true
})

router.afterEach((to) => {
  const items = [import.meta.env.VITE_TITLE || APP_TITLE]
  if (typeof to.meta.title === 'string' && to.meta.title.trim()) {
    items.unshift(to.meta.title)
  }
  document.title = items.join(' | ')
})

export default router
