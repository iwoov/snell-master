import request from '@/utils/request'

// Traffic Summary
export interface TrafficSummary {
    total_users: number
    active_users: number
    total_traffic: number
    today_traffic: number
    month_traffic: number
}

// User Traffic Rank
export interface UserTrafficRank {
    user_id: number
    username: string
    upload: number
    download: number
    total: number
}

// Node Traffic Rank
export interface NodeTrafficRank {
    node_id: number
    node_name: number
    upload: number
    download: number
    total: number
}

// Traffic Trend
export interface TrafficTrend {
    dates: string[]
    upload: number[]
    download: number[]
    total: number[]
}

// Get Traffic Summary
export function getTrafficSummary() {
    return request<TrafficSummary>({
        url: '/admin/traffic/summary',
        method: 'get'
    })
}

// Get User Traffic Ranking
export function getUserTrafficRanking(type: 'day' | 'week' | 'month' = 'day', limit: number = 10) {
    return request<UserTrafficRank[]>({
        url: '/admin/traffic/ranking/user',
        method: 'get',
        params: { type, limit }
    })
}

// Get Node Traffic Ranking
export function getNodeTrafficRanking(type: 'day' | 'week' | 'month' = 'day') {
    return request<NodeTrafficRank[]>({
        url: '/admin/traffic/ranking/node',
        method: 'get',
        params: { type }
    })
}

// Get Traffic Trend
export function getTrafficTrend(days: number = 30) {
    return request<TrafficTrend>({
        url: '/admin/traffic/trend',
        method: 'get',
        params: { days }
    })
}
