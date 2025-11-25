import request from '@/utils/request'

// Template Interface
export interface Template {
    id: number
    name: string
    content: string
    description?: string
    is_default: boolean
    created_at: string
    updated_at: string
}

// Create Template Request
export interface CreateTemplateRequest {
    name: string
    content: string
    description?: string
    is_default?: boolean
}

// Update Template Request
export interface UpdateTemplateRequest {
    name?: string
    content?: string
    description?: string
}

// Get Template List
export function getTemplateList() {
    return request<Template[]>({
        url: '/admin/templates',
        method: 'get'
    })
}

// Create Template
export function createTemplate(data: CreateTemplateRequest) {
    return request<Template>({
        url: '/admin/templates',
        method: 'post',
        data
    })
}

// Update Template
export function updateTemplate(id: number, data: UpdateTemplateRequest) {
    return request<Template>({
        url: `/admin/templates/${id}`,
        method: 'put',
        data
    })
}

// Delete Template
export function deleteTemplate(id: number) {
    return request({
        url: `/admin/templates/${id}`,
        method: 'delete'
    })
}

// Set Default Template
export function setDefaultTemplate(id: number) {
    return request({
        url: `/admin/templates/${id}/default`,
        method: 'post'
    })
}
