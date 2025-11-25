<template>
  <div class="dashboard">
    <h1>Dashboard</h1>
    <p>Welcome to Snell Master User Portal</p>
    <el-button type="primary" @click="handleLogout">Logout</el-button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.dashboard {
  padding: 20px;
}
</style>
