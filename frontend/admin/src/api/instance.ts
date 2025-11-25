import request from '@/utils/request'

// Instance Interface
export interface Instance {
    id: number
    user_id: number
    node_id: number
    port: number
    psk: string
    version: number
    obfs: string
    config_path: string
    service_name: string
    status: string  // 'running' | 'stopped'
    created_at: string
    updated_at: string
    user?: {
        id: number
        username: string
    }
    node?: {
        id: number
        name: string
        endpoint: string
    }
}

// Instance Filter
export interface InstanceFilter {
    user_id?: number
    node_id?: number
    status?: number
    page?: number
    page_size?: number
}

// Create Instance Request
export interface CreateInstanceRequest {
    user_id: number
    node_id: number
    port?: number
    psk?: string
}

// Get Instance List
export function getInstanceList(params: InstanceFilter) {
    return request<{
        items: Instance[]
        total: number
        page: number
        page_size: number
    }>({
        url: '/admin/instances',
        method: 'get',
        params
    })
}

// Create Instance
export function createInstance(data: CreateInstanceRequest) {
    return request<Instance>({
        url: '/admin/instances',
        method: 'post',
        data
    })
}

// Get Instance Detail
export function getInstance(id: number) {
    return request<Instance>({
        url: `/admin/instances/${id}`,
        method: 'get'
    })
}

// Delete Instance
export function deleteInstance(id: number) {
    return request({
        url: `/admin/instances/${id}`,
        method: 'delete'
    })
}

// Restart Instance
export function restartInstance(id: number) {
    return request({
        url: `/admin/instances/${id}/restart`,
        method: 'post'
    })
}

// Update Instance Status
export function updateInstanceStatus(id: number, status: string) {
    return request({
        url: `/admin/instances/${id}/status`,
        method: 'put',
        data: { status }
    })
}

// Get Instance Config
export function getInstanceConfig(id: number) {
    return request<{ config: string }>({
        url: `/admin/instances/${id}/config`,
        method: 'get'
    })
}
