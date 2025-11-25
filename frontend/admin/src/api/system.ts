import request from '@/utils/request'

/**
 * 系统配置类型
 */
export interface SystemConfig {
    id: number
    key: string
    value: string
    description: string
    updated_at: string
}

/**
 * 配置更新请求
 */
export interface UpdateConfigRequest {
    value: string
}

/**
 * 批量更新配置请求
 */
export type BatchUpdateConfigsRequest = Record<string, string>

/**
 * 获取所有系统配置
 */
export function getSystemConfigs() {
    return request<SystemConfig[]>({
        url: '/admin/system-configs',
        method: 'get'
    })
}

/**
 * 更新单个系统配置
 */
export function updateSystemConfig(key: string, value: string) {
    return request<void>({
        url: `/admin/system-configs/${key}`,
        method: 'put',
        data: { value }
    })
}

/**
 * 批量更新系统配置
 */
export function batchUpdateConfigs(configs: BatchUpdateConfigsRequest) {
    return request<void>({
        url: '/admin/system-configs',
        method: 'put',
        data: configs
    })
}
