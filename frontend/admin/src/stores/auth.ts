import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginRequest, LoginResponse, AdminInfo } from '@/types/api'
import { login as apiLogin, getAdminInfo as apiGetInfo } from '@/api/admin'

export const useAuthStore = defineStore('auth', () => {
    // State
    const token = ref<string>(localStorage.getItem('token') || '')
    const userInfo = ref<AdminInfo | null>(null)

    // Getters
    const isLoggedIn = computed(() => !!token.value)

    // Actions
    async function login(data: LoginRequest): Promise<LoginResponse> {
        const res = await apiLogin(data)
        token.value = res.token
        if (res.admin) {
            userInfo.value = res.admin
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
            const res = await apiGetInfo()
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
