import request from '@/utils/request'

// Node Interface
export interface Node {
    id: number
    name: string
    api_token: string
    endpoint: string
    location?: string
    country_code?: string
    status: string // 'online' | 'offline'
    cpu_usage: number
    memory_usage: number
    disk_usage: number
    bandwidth_usage: number
    instance_count?: number
    last_seen_at?: string
    created_at: string
}

// Create Node Request
export interface CreateNodeRequest {
    name: string
    endpoint: string
    location?: string
    country_code?: string
}

// Update Node Request
export interface UpdateNodeRequest {
    name?: string
    endpoint?: string
    location?: string
    country_code?: string
}

// Get Node List
export function getNodeList() {
    return request<Node[]>({
        url: '/admin/nodes',
        method: 'get'
    })
}

// Create Node
export function createNode(data: CreateNodeRequest) {
    return request<Node>({
        url: '/admin/nodes',
        method: 'post',
        data
    })
}

// Get Node Detail
export function getNode(id: number) {
    return request<Node>({
        url: `/admin/nodes/${id}`,
        method: 'get'
    })
}

// Update Node
export function updateNode(id: number, data: UpdateNodeRequest) {
    return request<Node>({
        url: `/admin/nodes/${id}`,
        method: 'put',
        data
    })
}

// Delete Node
export function deleteNode(id: number) {
    return request({
        url: `/admin/nodes/${id}`,
        method: 'delete'
    })
}

// Regenerate Token
export function regenerateToken(id: number) {
    return request<{ token: string }>({
        url: `/admin/nodes/${id}/token`,
        method: 'post'
    })
}

// Download Install Script
export function downloadInstallScript(id: number, nodeName: string) {
    return request({
        url: `/admin/nodes/${id}/install-script`,
        method: 'get',
        responseType: 'blob'
    }).then((blob: Blob) => {
        // Create download link
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `install-agent-${nodeName}.sh`
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
    })
}
