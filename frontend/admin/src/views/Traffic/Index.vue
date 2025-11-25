<template>
  <div class="traffic-stats">
    <div class="header">
      <h1>流量统计</h1>
    </div>

    <!-- Summary Cards -->
    <el-row :gutter="20" class="summary-row">
      <el-col :span="6">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>总流量</span>
            </div>
          </template>
          <div class="card-value">{{ formatBytes(summary.total_traffic) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>今日流量</span>
            </div>
          </template>
          <div class="card-value">{{ formatBytes(summary.today_traffic) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>本月流量</span>
            </div>
          </template>
          <div class="card-value">{{ formatBytes(summary.month_traffic) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>活跃用户</span>
            </div>
          </template>
          <div class="card-value">{{ summary.active_users }} / {{ summary.total_users }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Tabs -->
    <el-tabs v-model="activeTab" class="traffic-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="用户排行" name="user">
        <el-card shadow="never">
          <div class="tab-header">
            <el-radio-group v-model="rankingType" size="small" @change="fetchUserRanking">
              <el-radio-button label="day">今日</el-radio-button>
              <el-radio-button label="week">本周</el-radio-button>
              <el-radio-button label="month">本月</el-radio-button>
            </el-radio-group>
          </div>
          <el-table :data="userRanking" :style="{ width: '100%' }" v-loading="loading">
            <el-table-column type="index" label="排名" width="80" />
            <el-table-column prop="username" label="用户名" min-width="150" />
            <el-table-column label="上传" width="150">
              <template #default="{ row }">
                {{ formatBytes(row.upload) }}
              </template>
            </el-table-column>
            <el-table-column label="下载" width="150">
              <template #default="{ row }">
                {{ formatBytes(row.download) }}
              </template>
            </el-table-column>
            <el-table-column label="总计" width="150">
              <template #default="{ row }">
                {{ formatBytes(row.total) }}
              </template>
            </el-table-column>
            <el-table-column label="占比" min-width="200">
              <template #default="{ row }">
                <el-progress :percentage="calculatePercentage(row.total, maxUserTraffic)" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="节点流量" name="node">
        <el-card shadow="never">
          <div class="tab-header">
            <el-radio-group v-model="rankingType" size="small" @change="fetchNodeRanking">
              <el-radio-button label="day">今日</el-radio-button>
              <el-radio-button label="week">本周</el-radio-button>
              <el-radio-button label="month">本月</el-radio-button>
            </el-radio-group>
          </div>
          <el-table :data="nodeRanking" :style="{ width: '100%' }" v-loading="loading">
            <el-table-column prop="node_name" label="节点名称" min-width="150" />
            <el-table-column label="上传" width="150">
              <template #default="{ row }">
                {{ formatBytes(row.upload) }}
              </template>
            </el-table-column>
            <el-table-column label="下载" width="150">
              <template #default="{ row }">
                {{ formatBytes(row.download) }}
              </template>
            </el-table-column>
            <el-table-column label="总计" width="150">
              <template #default="{ row }">
                {{ formatBytes(row.total) }}
              </template>
            </el-table-column>
            <el-table-column label="占比" min-width="200">
              <template #default="{ row }">
                <el-progress :percentage="calculatePercentage(row.total, maxNodeTraffic)" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="流量趋势" name="trend">
        <el-card shadow="never">
          <div class="chart-placeholder">
            <el-alert
              title="图表组件加载失败"
              type="warning"
              description="由于依赖安装问题，暂时无法显示图表。以下是最近30天的流量数据列表。"
              show-icon
              :closable="false"
              style="margin-bottom: 20px"
            />
          </div>
          <el-table :data="trendList" :style="{ width: '100%' }" v-loading="loading" height="500">
            <el-table-column prop="date" label="日期" width="180" />
            <el-table-column label="上传" min-width="150">
              <template #default="{ row }">
                {{ formatBytes(row.upload) }}
              </template>
            </el-table-column>
            <el-table-column label="下载" min-width="150">
              <template #default="{ row }">
                {{ formatBytes(row.download) }}
              </template>
            </el-table-column>
            <el-table-column label="总计" min-width="150">
              <template #default="{ row }">
                {{ formatBytes(row.total) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  getTrafficSummary,
  getUserTrafficRanking,
  getNodeTrafficRanking,
  getTrafficTrend
} from '@/api/traffic'
import type { TrafficSummary, UserTrafficRank, NodeTrafficRank, TrafficTrend } from '@/api/traffic'

// State
const loading = ref(false)
const activeTab = ref('user')
const rankingType = ref<'day' | 'week' | 'month'>('day')

const summary = ref<TrafficSummary>({
  total_users: 0,
  active_users: 0,
  total_traffic: 0,
  today_traffic: 0,
  month_traffic: 0
})

const userRanking = ref<UserTrafficRank[]>([])
const nodeRanking = ref<NodeTrafficRank[]>([])
const trendData = ref<TrafficTrend>({ dates: [], upload: [], download: [], total: [] })

// Computed
const maxUserTraffic = computed(() => {
  if (userRanking.value.length === 0) return 0
  return Math.max(...userRanking.value.map(item => item.total))
})

const maxNodeTraffic = computed(() => {
  if (nodeRanking.value.length === 0) return 0
  return Math.max(...nodeRanking.value.map(item => item.total))
})

const trendList = computed(() => {
  return trendData.value.dates.map((date, index) => ({
    date,
    upload: trendData.value.upload[index],
    download: trendData.value.download[index],
    total: trendData.value.total[index]
  })).reverse() // Show newest first
})

// Helpers
const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const calculatePercentage = (value: number, max: number) => {
  if (max === 0) return 0
  return Math.round((value / max) * 100)
}

// Actions
const fetchSummary = async () => {
  try {
    const res = await getTrafficSummary()
    summary.value = res
  } catch (error) {
    console.error(error)
  }
}

const fetchUserRanking = async () => {
  loading.value = true
  try {
    const res = await getUserTrafficRanking(rankingType.value)
    userRanking.value = res
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchNodeRanking = async () => {
  loading.value = true
  try {
    const res = await getNodeTrafficRanking(rankingType.value)
    nodeRanking.value = res
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchTrend = async () => {
  loading.value = true
  try {
    const res = await getTrafficTrend()
    trendData.value = res
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleTabChange = (tab: string | number) => {
  if (tab === 'user') {
    fetchUserRanking()
  } else if (tab === 'node') {
    fetchNodeRanking()
  } else if (tab === 'trend') {
    fetchTrend()
  }
}

onMounted(() => {
  fetchSummary()
  fetchUserRanking() // Default tab
})
</script>

<style scoped>
.traffic-stats {
  padding: 24px;
}

.header {
  margin-bottom: 24px;
}

.header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.summary-row {
  margin-bottom: 24px;
}

.card-header {
  font-size: 14px;
  color: #606266;
}

.card-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-top: 8px;
}

.traffic-tabs {
  background: #fff;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.tab-header {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
