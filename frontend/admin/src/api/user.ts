import request from '@/utils/request'
import type { UserInfo } from '@/types/api'

// User List Filter
export interface UserFilter {
    keyword?: string
    status?: number
    page?: number
    page_size?: number
}

// Create User Request
export interface CreateUserRequest {
    username: string
    password: string
    email?: string
    traffic_limit?: number
}

// Update User Request
export interface UpdateUserRequest {
    email?: string
    password?: string
    traffic_limit?: number
    status?: number
}

// Assign Nodes Request
export interface AssignNodesRequest {
    node_ids: number[]
}

// Get User List
export function getUserList(params: UserFilter) {
    return request<{
        items: UserInfo[]
        total: number
        page: number
        page_size: number
    }>({
        url: '/admin/users',
        method: 'get',
        params
    })
}

// Create User
export function createUser(data: CreateUserRequest) {
    return request<UserInfo>({
        url: '/admin/users',
        method: 'post',
        data
    })
}

// Get User Detail
export function getUser(id: number) {
    return request<UserInfo>({
        url: `/admin/users/${id}`,
        method: 'get'
    })
}

// Update User
export function updateUser(id: number, data: UpdateUserRequest) {
    return request<UserInfo>({
        url: `/admin/users/${id}`,
        method: 'put',
        data
    })
}

// Delete User
export function deleteUser(id: number) {
    return request({
        url: `/admin/users/${id}`,
        method: 'delete'
    })
}

// Reset User Traffic
export function resetUserTraffic(id: number) {
    return request({
        url: `/admin/users/${id}/reset-traffic`,
        method: 'post'
    })
}

// Update User Status
export function updateUserStatus(id: number, status: number) {
    return request({
        url: `/admin/users/${id}/status`,
        method: 'post',
        data: { status }
    })
}

// Assign Nodes
export function assignNodes(id: number, node_ids: number[]) {
    return request({
        url: `/admin/users/${id}/nodes`,
        method: 'post',
        data: { node_ids }
    })
}
