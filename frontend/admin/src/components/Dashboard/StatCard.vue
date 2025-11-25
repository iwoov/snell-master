<template>
  <el-card class="stat-card" :body-style="{ padding: '20px' }" shadow="hover">
    <div class="stat-content">
      <div class="stat-icon" :style="{ backgroundColor: color + '15', color: color }">
        <component :is="icon" />
      </div>
      <div class="stat-info">
        <div class="stat-title">{{ title }}</div>
        <div class="stat-value">
          <span v-if="prefix" class="stat-prefix">{{ prefix }}</span>
          {{ value }}
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  title: string
  value: string | number
  icon: Component
  color: string
  prefix?: string
}>()
</script>

<style scoped>
.stat-card {
  height: 100%;
  border: none;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
  overflow: hidden;
}

.stat-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  line-height: 1.2;
}

.stat-prefix {
  font-size: 14px;
  margin-right: 4px;
  color: #909399;
}
</style>
