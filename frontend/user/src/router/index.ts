import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const routes: RouteRecordRaw[] = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: {
            requiresAuth: false,
            title: 'Login'
        }
    },
    {
        path: '/',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: {
            requiresAuth: true,
            title: 'Dashboard'
        }
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/404.vue'),
        meta: {
            requiresAuth: false,
            title: 'Page Not Found'
        }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Route Guard
router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore()

    // Set Title
    document.title = to.meta.title
        ? `${to.meta.title} - Snell Master`
        : 'Snell Master'

    // Auth Check
    if (to.meta.requiresAuth) {
        if (authStore.isLoggedIn) {
            next()
        } else {
            ElMessage.warning('Please login first')
            next({
                path: '/login',
                query: { redirect: to.fullPath }
            })
        }
    }
    // Login Page Check
    else if (to.path === '/login') {
        if (authStore.isLoggedIn) {
            next('/')
        } else {
            next()
        }
    }
    // Others
    else {
        next()
    }
})

export default router
