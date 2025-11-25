<template>
  <transition name="breadcrumb-fade">
    <el-breadcrumb v-if="breadcrumbs.length > 0" separator="/" class="breadcrumb-container">
      <transition-group name="breadcrumb-item">
        <el-breadcrumb-item
          v-for="(item, index) in breadcrumbs"
          :key="item.path"
          :to="index < breadcrumbs.length - 1 ? item.path : undefined"
        >
          {{ item.meta.title }}
        </el-breadcrumb-item>
      </transition-group>
    </el-breadcrumb>
  </transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import type { RouteLocationMatched } from 'vue-router'

const route = useRoute()

// 面包屑数据
const breadcrumbs = computed(() => {
  const matched = route.matched.filter(
    (item: RouteLocationMatched) => item.meta && item.meta.title
  )
  return matched
})
</script>

<style scoped>
.breadcrumb-container {
  margin-bottom: 16px;
  padding: 12px 16px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
}

/* 面包屑淡入动画 */
.breadcrumb-fade-enter-active,
.breadcrumb-fade-leave-active {
  transition: opacity 0.3s ease;
}

.breadcrumb-fade-enter-from,
.breadcrumb-fade-leave-to {
  opacity: 0;
}

/* 面包屑项动画 */
.breadcrumb-item-enter-active {
  transition: all 0.3s ease;
}

.breadcrumb-item-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}

:deep(.el-breadcrumb__item) {
  font-size: 14px;
}

:deep(.el-breadcrumb__inner) {
  color: #00000073;
  font-weight: 400;
  transition: color 0.3s ease;
}

:deep(.el-breadcrumb__inner:hover) {
  color: #1890ff;
}

:deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  color: #000000d9;
  font-weight: 500;
}

:deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner:hover) {
  color: #000000d9;
  cursor: default;
}
</style>
