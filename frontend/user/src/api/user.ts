import request from '@/utils/request'
import type { LoginRequest, LoginResponse, UserInfo } from '@/types/api'

// User Login
export function login(data: LoginRequest) {
    return request<LoginResponse>({
        url: '/auth/user/login',
        method: 'post',
        data
    })
}

// Get User Profile
export function getUserProfile() {
    return request<UserInfo>({
        url: '/user/profile',
        method: 'get'
    })
}

// Update Profile
export function updateProfile(data: Partial<UserInfo>) {
    return request({
        url: '/user/profile',
        method: 'put',
        data
    })
}
