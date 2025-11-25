<template>
  <div class="sidebar-container">
    <!-- Logo区域 -->
    <div class="logo-container">
      <transition name="logo-fade" mode="out-in">
        <div v-if="!isCollapsed" key="full" class="logo-full">
          <span class="logo-icon">
            <el-icon><Monitor /></el-icon>
          </span>
          <span class="logo-text">Snell Master</span>
        </div>
        <div v-else key="mini" class="logo-mini">
          <span class="logo-icon">
            <el-icon><Monitor /></el-icon>
          </span>
        </div>
      </transition>
    </div>

    <!-- 菜单 -->
    <el-menu
      :default-active="currentRoute"
      :collapse="isCollapsed"
      :unique-opened="true"
      router
      class="sidebar-menu"
    >
      <el-menu-item index="/dashboard">
        <el-icon><Odometer /></el-icon>
        <template #title>仪表盘</template>
      </el-menu-item>

      <el-sub-menu index="users">
        <template #title>
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </template>
        <el-menu-item index="/users">用户列表</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="nodes">
        <template #title>
          <el-icon><Monitor /></el-icon>
          <span>节点管理</span>
        </template>
        <el-menu-item index="/nodes">节点列表</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="instances">
        <template #title>
          <el-icon><Connection /></el-icon>
          <span>实例管理</span>
        </template>
        <el-menu-item index="/instances">实例列表</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="traffic">
        <template #title>
          <el-icon><TrendCharts /></el-icon>
          <span>流量统计</span>
        </template>
        <el-menu-item index="/traffic">流量概览</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="subscription">
        <template #title>
          <el-icon><Link /></el-icon>
          <span>订阅管理</span>
        </template>
        <el-menu-item index="/subscriptions">订阅列表</el-menu-item>
        <el-menu-item index="/templates">模板管理</el-menu-item>
      </el-sub-menu>

      <el-menu-item index="/logs">
        <el-icon><Document /></el-icon>
        <template #title>操作日志</template>
      </el-menu-item>

      <el-menu-item index="/system/config">
        <el-icon><Setting /></el-icon>
        <template #title>系统配置</template>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  Odometer,
  User,
  Monitor,
  Connection,
  TrendCharts,
  Link,
  Document,
  Setting
} from '@element-plus/icons-vue'

interface Props {
  isCollapsed: boolean
}

defineProps<Props>()

const route = useRoute()

// 当前激活的路由
const currentRoute = computed(() => route.path)
</script>

<style scoped>
.sidebar-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  border-right: 1px solid #e4e7ed;
}

.logo-container {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  border-bottom: 1px solid #f0f2f5;
  overflow: hidden;
}

.logo-full {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-mini {
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon {
  font-size: 24px;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
  letter-spacing: 0.5px;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 滚动条样式 */
.sidebar-menu::-webkit-scrollbar {
  width: 6px;
}

.sidebar-menu::-webkit-scrollbar-thumb {
  background-color: #e4e7ed;
  border-radius: 3px;
}

.sidebar-menu::-webkit-scrollbar-thumb:hover {
  background-color: #c0c4cc;
}

/* Logo动画 */
.logo-fade-enter-active,
.logo-fade-leave-active {
  transition: opacity 0.2s ease;
}

.logo-fade-enter-from,
.logo-fade-leave-to {
  opacity: 0;
}

/* 修复 Element Plus Menu 在折叠时的样式 */
:deep(.el-menu--collapse) {
  width: 64px;
}

:deep(.el-sub-menu__title),
:deep(.el-menu-item) {
  height: 50px;
  line-height: 50px;
  margin: 4px 0;
  color: #606266;
}

:deep(.el-sub-menu__title:hover),
:deep(.el-menu-item:hover) {
  background-color: var(--el-color-primary-light-9) !important;
  color: var(--el-color-primary);
}

:deep(.el-menu-item.is-active) {
  background-color: var(--el-color-primary-light-9) !important;
  color: var(--el-color-primary);
  border-right: 3px solid var(--el-color-primary);
}

:deep(.el-menu-item .el-icon),
:deep(.el-sub-menu__title .el-icon) {
  font-size: 18px;
  color: inherit;
}
</style>
