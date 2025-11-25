import request from '@/utils/request'
import type { LoginRequest, LoginResponse, AdminInfo } from '@/types/api'

// Admin Login
export function login(data: LoginRequest) {
    return request<LoginResponse>({
        url: '/auth/admin/login',
        method: 'post',
        data
    })
}

// Get Admin Info
export function getAdminInfo() {
    return request<AdminInfo>({
        url: '/admin/profile',
        method: 'get'
    })
}

// Change Password
export function changePassword(data: {
    old_password: string
    new_password: string
}) {
    return request({
        url: '/admin/password',
        method: 'post',
        data
    })
}
