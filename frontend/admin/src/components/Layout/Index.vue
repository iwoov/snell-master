<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapsed ? '64px' : '200px'" class="layout-aside">
      <Sidebar :is-collapsed="isCollapsed" />
    </el-aside>
    <el-container>
      <el-header class="layout-header">
        <Header @toggle-sidebar="toggleSidebar" :is-collapsed="isCollapsed" />
      </el-header>
      <el-main class="layout-main">
        <Breadcrumb />
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './Sidebar.vue'
import Header from './Header.vue'
import Breadcrumb from './Breadcrumb.vue'

// 侧边栏折叠状态
const isCollapsed = ref(false)

// 切换侧边栏
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}
</script>

<style scoped>
.layout-container {
  min-height: 100vh;
}

.layout-aside {
  background-color: #ffffff;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.layout-header {
  background: #fff;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  height: 64px;
}

.layout-main {
  background-color: #f0f2f5;
  padding: 24px;
  min-height: calc(100vh - 64px);
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
