import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'
import { ElMessage, ElLoading } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

let loadingInstance: any = null
let loadingCount = 0

const service: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json'
    }
})

// Show Loading
function showLoading() {
    if (loadingCount === 0) {
        loadingInstance = ElLoading.service({
            lock: true,
            text: 'Loading...',
            background: 'rgba(0, 0, 0, 0.7)'
        })
    }
    loadingCount++
}

// Hide Loading
function hideLoading() {
    loadingCount--
    if (loadingCount === 0 && loadingInstance) {
        loadingInstance.close()
        loadingInstance = null
    }
}

// Request Interceptor
service.interceptors.request.use(
    (config) => {
        const authStore = useAuthStore()

        // Add Token
        if (authStore.token) {
            config.headers['Authorization'] = `Bearer ${authStore.token}`
        }

        // Show Loading (configurable)
        // @ts-ignore
        if (config.loading !== false) {
            showLoading()
        }

        return config
    },
    (error) => {
        hideLoading()
        ElMessage.error('Request failed')
        return Promise.reject(error)
    }
)

// Response Interceptor
service.interceptors.response.use(
    (response) => {
        hideLoading()

        // For blob responses (file downloads), return raw response data
        if (response.config.responseType === 'blob') {
            return response.data
        }

        const res = response.data

        // Business Success
        if (res.code === 0) {
            return res.data
        }

        // Business Failure
        ElMessage.error(res.message || 'Operation failed')

        // 401 Unauthorized
        if (res.code === 401) {
            const authStore = useAuthStore()
            authStore.logout()
            router.push('/login')
            ElMessage.warning('Login expired, please login again')
        }

        // 403 Forbidden
        if (res.code === 403) {
            ElMessage.error('Permission denied')
        }

        return Promise.reject(new Error(res.message || 'Operation failed'))
    },
    (error) => {
        hideLoading()

        if (error.response) {
            const status = error.response.status

            switch (status) {
                case 401:
                    ElMessage.error('Unauthorized, please login')
                    break
                case 403:
                    ElMessage.error('Permission denied')
                    break
                case 404:
                    ElMessage.error('Resource not found')
                    break
                case 500:
                    ElMessage.error('Server error')
                    break
                default:
                    ElMessage.error('Network error')
            }
        } else if (error.request) {
            ElMessage.error('Network connection failed')
        } else {
            ElMessage.error('Request configuration error')
        }

        return Promise.reject(error)
    }
)

// Wrapper function
const request = <T = any>(config: AxiosRequestConfig): Promise<T> => {
    return service.request(config)
}

export default request
