import request from '@/utils/request'

// Subscription Interface
export interface Subscription {
    id: number
    user_id: number
    token: string
    template_id?: number
    access_count: number
    last_access_at?: string
    created_at: string
    user?: {
        id: number
        username: string
    }
    template?: {
        id: number
        name: string
    }
}

// Create Subscription Request
export interface CreateSubscriptionRequest {
    user_id: number
    template_id?: number
}

// Get Subscription List
export function getSubscriptionList() {
    return request<Subscription[]>({
        url: '/admin/subscriptions',
        method: 'get'
    })
}

// Create Subscription
export function createSubscription(data: CreateSubscriptionRequest) {
    return request<Subscription>({
        url: '/admin/subscriptions',
        method: 'post',
        data
    })
}

// Delete Subscription
export function deleteSubscription(id: number) {
    return request({
        url: `/admin/subscriptions/${id}`,
        method: 'delete'
    })
}

// Regenerate Token
export function regenerateSubscriptionToken(id: number) {
    return request<{ token: string }>({
        url: `/admin/subscriptions/${id}/regenerate`,
        method: 'post'
    })
}
