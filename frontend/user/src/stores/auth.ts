import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginRequest, LoginResponse, UserInfo } from '@/types/api'
import { login as apiLogin, getUserProfile as apiGetProfile } from '@/api/user'

export const useAuthStore = defineStore('auth', () => {
    // State
    const token = ref<string>(localStorage.getItem('token') || '')
    const userInfo = ref<UserInfo | null>(null)

    // Getters
    const isLoggedIn = computed(() => !!token.value)

    // Actions
    async function login(data: LoginRequest): Promise<LoginResponse> {
        const res = await apiLogin(data)
        token.value = res.token
        if (res.user) {
            userInfo.value = res.user
        }
        localStorage.setItem('token', res.token)
        return res
    }

    function logout() {
        token.value = ''
        userInfo.value = null
        localStorage.removeItem('token')
    }

    async function getUserInfo() {
        try {
            const res = await apiGetProfile()
            userInfo.value = res
            return res
        } catch (error) {
            logout()
            throw error
        }
    }

    return {
        token,
        userInfo,
        isLoggedIn,
        login,
        logout,
        getUserInfo
    }
})
