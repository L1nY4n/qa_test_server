import axios from 'axios'
import { base_url, requestTimeout } from '@/config'
import { clearSession, getToken } from '@/utils/auth'

const instance = axios.create({
  baseURL: base_url,
  timeout: requestTimeout,
  headers: {
    'Content-Type': 'application/json',
  },
})

instance.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

instance.interceptors.response.use(
  (response) => {
    const { data, config, status } = response
    const { success, data: payload, error: message } = data
    if (success) {
      return payload
    }
    return Promise.reject(
      new Error(message || ('请求失败: ' + JSON.stringify({ url: config.url, status, payload })))
    )
  },
  (error) => {
    const status = error?.response?.status
    if (status === 401) {
      clearSession()
      if (!window.location.hash.includes('/login')) {
        window.location.hash = '#/login'
      }
    }

    if (error.response && error.response.data) {
      const message = error.response.data?.error || error.response.data?.message
      if (typeof message === 'string' && message.trim()) {
        return Promise.reject(new Error(message))
      }
      return Promise.reject(error)
    }

    let { message } = error
    if (message === 'Network Error') {
      message = '接口连接异常'
    }
    if (message.includes('timeout')) {
      message = '接口请求超时'
    }
    if (message.includes('Request failed with status code')) {
      const code = message.substring(message.length - 3)
      message = '后端接口 ' + code + ' 异常'
    }
    return Promise.reject(new Error(message))
  }
)

export default instance
