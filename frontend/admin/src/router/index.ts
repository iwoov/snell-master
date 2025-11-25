import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import Layout from '@/components/Layout/Index.vue'

const routes: RouteRecordRaw[] = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: {
            requiresAuth: false,
            title: '登录'
        }
    },
    {
        path: '/',
        component: Layout,
        redirect: '/dashboard',
        meta: {
            requiresAuth: true
        },
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/Dashboard.vue'),
                meta: {
                    requiresAuth: true,
                    title: '仪表盘'
                }
            },
            {
                path: 'users',
                name: 'Users',
                component: () => import('@/views/User/Index.vue'),
                meta: {
                    requiresAuth: true,
                    title: '用户管理'
                }
            },
            {
                path: 'nodes',
                name: 'Nodes',
                component: () => import('@/views/Node/Index.vue'),
                meta: {
                    requiresAuth: true,
                    title: '节点管理'
                }
            },
            {
                path: 'instances',
                name: 'Instances',
                component: () => import('@/views/Instance/Index.vue'),
                meta: {
                    requiresAuth: true,
                    title: '实例管理'
                }
            },
            {
                path: 'traffic',
                name: 'Traffic',
                component: () => import('@/views/Traffic/Index.vue'),
                meta: {
                    requiresAuth: true,
                    title: '流量统计'
                }
            },
            {
                path: 'subscriptions',
                name: 'Subscriptions',
                component: () => import('@/views/Subscription/Index.vue'),
                meta: {
                    requiresAuth: true,
                    title: '订阅管理'
                }
            },
            {
                path: 'templates',
                name: 'Templates',
                component: () => import('@/views/Template/Index.vue'),
                meta: {
                    requiresAuth: true,
                    title: '模板管理'
                }
            },
            {
                path: 'logs',
                name: 'Logs',
                component: () => import('@/views/Log/Index.vue'),
                meta: {
                    requiresAuth: true,
                    title: '操作日志'
                }
            },
            {
                path: 'system/config',
                name: 'SystemConfig',
                component: () => import('@/views/System/Config.vue'),
                meta: {
                    requiresAuth: true,
                    title: '系统配置'
                }
            }
        ]
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/404.vue'),
        meta: {
            requiresAuth: false,
            title: '页面未找到'
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
