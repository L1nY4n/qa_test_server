import axios from 'axios'
import {
    base_url,
    requestTimeout
} from '@/config'


const instance = axios.create({
    baseURL: base_url,
    timeout: requestTimeout,
    headers: {
        'Content-Type': "application/json",
    },
})

instance.interceptors.request.use(
    (config) => {
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

instance.interceptors.response.use(
    (response) => {
        const { data, config } = response
        const { success, data: d } = data
        if (success) {
            return d
        } else {
            return Promise.reject(
                '请求异常:' +
                JSON.stringify({ url: config.url, status, d }) || 'Error'
            )
        }
    },
    (error) => {

        const { response, message } = error
        if (error.response && error.response.data) {
            const { success, data } = response

            return Promise.reject(error)
        } else {
            let { message } = error
            if (message === 'Network Error') {
                message = '接口连接异常'
            }
            if (message.includes('timeout')) {
                message = '接口请求超时'
            }
            if (message.includes('Request failed with status code')) {
                const code = message.substr(message.length - 3)
                message = '后端接口' + code + '异常'
            }
            return Promise.reject(error)
        }
    }
)

export default instance