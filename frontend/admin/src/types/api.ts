// Generic API Response
export interface ApiResponse<T = any> {
    code: number
    message: string
    data: T
}

// Pagination Response
export interface PaginationResponse<T> {
    list: T[]
    total: number
    page: number
    page_size: number
}

// Login Request
export interface LoginRequest {
    username: string
    password: string
}

// Login Response
export interface LoginResponse {
    token: string
    admin?: AdminInfo
    user?: UserInfo
}

// Admin Info
export interface AdminInfo {
    id: number
    username: string
    email: string
    role: number // 1: admin, 2: super_admin
    created_at: string
}

// User Info
export interface UserInfo {
    id: number
    username: string
    email: string
    traffic_limit: number
    traffic_used_today: number
    traffic_used_month: number
    traffic_used_total: number
    status: number
    created_at: string
}
