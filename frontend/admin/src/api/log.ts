import request from '@/utils/request'

// Operation Log Interface
export interface OperationLog {
    id: number
    admin_id: number
    admin_name: string
    action: string
    target_type: string
    target_id?: number
    description: string
    ip_address: string
    user_agent?: string
    request_body?: string
    response_body?: string
    created_at: string
}

// Log Filter
export interface LogFilter {
    admin_id?: number
    action?: string
    target_type?: string
    page?: number
    page_size?: number
}

// Get Operation Log List
export function getOperationLogs(params: LogFilter) {
    return request<{
        items: OperationLog[]
        total: number
        page: number
        page_size: number
    }>({
        url: '/admin/logs',
        method: 'get',
        params
    })
}
