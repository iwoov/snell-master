<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1>仪表盘</h1>
      <p class="subtitle">系统运行概况</p>
    </div>

    <!-- Loading Skeleton -->
    <el-skeleton v-if="loading" :rows="3" animated />

    <!-- Error State -->
    <el-alert
      v-else-if="error"
      title="获取数据失败"
      type="error"
      :description="error"
      show-icon
      :closable="false"
    />

    <!-- Dashboard Content -->
    <div v-else class="dashboard-content">
      <!-- Statistics Cards -->
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="6" v-for="card in statCards" :key="card.title">
          <StatCard
            :title="card.title"
            :value="card.value"
            :icon="card.icon"
            :color="card.color"
            :prefix="card.prefix"
          />
        </el-col>
      </el-row>

      <!-- Placeholder for Future Charts -->
      <el-row :gutter="20" class="mt-4">
        <el-col :span="24">
          <el-card shadow="never" class="placeholder-card">
            <div class="placeholder-content">
              <el-icon class="placeholder-icon"><TrendCharts /></el-icon>
              <p>流量趋势图将在后续版本中推出</p>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { User, Monitor, Connection, DataLine, TrendCharts } from '@element-plus/icons-vue'
import StatCard from '@/components/Dashboard/StatCard.vue'
import { getDashboardStats, type DashboardStats } from '@/api/dashboard'

const loading = ref(true)
const error = ref('')
const stats = ref<DashboardStats | null>(null)

// Format bytes to human readable string
const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const statCards = computed(() => {
  if (!stats.value) return []
  
  const cards: { title: string; value: string | number; icon: any; color: string; prefix?: string }[] = [
    {
      title: '用户总数 / 活跃',
      value: `${stats.value.total_users} / ${stats.value.active_users}`,
      icon: User,
      color: '#409EFF'
    },
    {
      title: '节点总数 / 在线',
      value: `${stats.value.total_nodes} / ${stats.value.online_nodes}`,
      icon: Monitor,
      color: '#67C23A'
    },
    {
      title: '实例总数',
      value: stats.value.total_instances,
      icon: Connection,
      color: '#E6A23C'
    },
    {
      title: '今日流量',
      value: formatBytes(stats.value.today_traffic),
      icon: DataLine,
      color: '#F56C6C'
    }
  ]
  return cards
})

const fetchData = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await getDashboardStats()
    stats.value = res // request interceptor returns data directly
  } catch (err: any) {
    error.value = err.message || '获取数据失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard {
  padding: 24px;
}

.dashboard-header {
  margin-bottom: 24px;
}

.dashboard-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.mt-4 {
  margin-top: 24px;
}

.placeholder-card {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-style: dashed;
}

.placeholder-content {
  text-align: center;
  color: #909399;
}

.placeholder-icon {
  font-size: 48px;
  margin-bottom: 16px;
}
</style>
