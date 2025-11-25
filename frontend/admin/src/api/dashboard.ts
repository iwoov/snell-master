import request from '@/utils/request'
// import type { ApiResponse } from '@/types/api'

// Dashboard Statistics Interface
export interface DashboardStats {
    total_users: number
    active_users: number
    total_nodes: number
    online_nodes: number
    total_instances: number
    total_traffic: number
    today_traffic: number
}

// Get Dashboard Statistics
export function getDashboardStats() {
    return request<DashboardStats>({
        url: '/admin/dashboard/stats',
        method: 'get'
    })
}

// Placeholder for future implementation: Get Traffic Trend
export function getTrafficTrend() {
    // return request<TrafficTrend[]>({
    //     url: '/admin/dashboard/traffic-trend',
    //     method: 'get'
    // })
    return Promise.resolve({ data: [] })
}

// Placeholder for future implementation: Get Top Users
export function getTopUsers() {
    // return request<UserTraffic[]>({
    //     url: '/admin/dashboard/top-users',
    //     method: 'get'
    // })
    return Promise.resolve({ data: [] })
}

// Placeholder for future implementation: Get Recent Logs
export function getRecentLogs() {
    // return request<OperationLog[]>({
    //     url: '/admin/dashboard/recent-logs',
    //     method: 'get'
    // })
    return Promise.resolve({ data: [] })
}
