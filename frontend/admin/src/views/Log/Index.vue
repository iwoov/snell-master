<template>
  <div class="operation-log">
    <div class="header">
      <h1>操作日志</h1>
    </div>

    <!-- Filter -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="操作类型">
          <el-select v-model="filterForm.action" placeholder="全部" clearable style="width: 150px">
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="查询" value="read" />
          </el-select>
        </el-form-item>
        <el-form-item label="对象类型">
          <el-select v-model="filterForm.target_type" placeholder="全部" clearable style="width: 150px">
            <el-option label="用户" value="user" />
            <el-option label="节点" value="node" />
            <el-option label="实例" value="instance" />
            <el-option label="订阅" value="subscription" />
            <el-option label="模板" value="template" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Log Table -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="logList" :style="{ width: '100%' }">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="admin_name" label="操作人" width="120" />
        <el-table-column label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)">
              {{ getActionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="对象类型" width="120">
          <template #default="{ row }">
            {{ getTargetTypeLabel(row.target_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="操作描述" min-width="250" />
        <el-table-column prop="ip_address" label="IP地址" width="150" />
        <el-table-column prop="created_at" label="操作时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- Detail Dialog -->
    <el-dialog
      v-model="detailDialogVisible"
      title="日志详情"
      width="700px"
    >
      <el-descriptions :column="1" border v-if="currentLog">
        <el-descriptions-item label="操作人">
          {{ currentLog.admin_name }}
        </el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getActionType(currentLog.action)">
            {{ getActionLabel(currentLog.action) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="对象类型">
          {{ getTargetTypeLabel(currentLog.target_type) }}
        </el-descriptions-item>
        <el-descriptions-item label="对象ID" v-if="currentLog.target_id">
          {{ currentLog.target_id }}
        </el-descriptions-item>
        <el-descriptions-item label="操作描述">
          {{ currentLog.description }}
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">
          {{ currentLog.ip_address }}
        </el-descriptions-item>
        <el-descriptions-item label="User Agent" v-if="currentLog.user_agent">
          {{ currentLog.user_agent }}
        </el-descriptions-item>
        <el-descriptions-item label="操作时间">
          {{ formatDate(currentLog.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="请求数据" v-if="currentLog.request_body">
          <pre class="json-content">{{ formatJSON(currentLog.request_body) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="响应数据" v-if="currentLog.response_body">
          <pre class="json-content">{{ formatJSON(currentLog.response_body) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="detailDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import dayjs from 'dayjs'
import { getOperationLogs } from '@/api/log'
import type { OperationLog } from '@/api/log'

// State
const loading = ref(false)
const logList = ref<OperationLog[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const filterForm = reactive({
  action: undefined as string | undefined,
  target_type: undefined as string | undefined
})

// Detail Dialog State
const detailDialogVisible = ref(false)
const currentLog = ref<OperationLog | null>(null)

// Helpers
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const formatJSON = (jsonStr: string) => {
  try {
    return JSON.stringify(JSON.parse(jsonStr), null, 2)
  } catch {
    return jsonStr
  }
}

const getActionType = (action: string): 'success' | 'warning' | 'danger' | 'info' => {
  const types: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
    create: 'success',
    update: 'warning',
    delete: 'danger',
    read: 'info'
  }
  return types[action] || 'info'
}

const getActionLabel = (action: string) => {
  const labels: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    read: '查询'
  }
  return labels[action] || action
}

const getTargetTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    user: '用户',
    node: '节点',
    instance: '实例',
    subscription: '订阅',
    template: '模板',
    admin: '管理员'
  }
  return labels[type] || type
}

// Data Fetching
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getOperationLogs({
      page: currentPage.value,
      page_size: pageSize.value,
      action: filterForm.action,
      target_type: filterForm.target_type
    })
    logList.value = res.items
    total.value = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// Actions
const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

const resetFilter = () => {
  filterForm.action = undefined
  filterForm.target_type = undefined
  handleSearch()
}

const handleViewDetail = (row: OperationLog) => {
  currentLog.value = row
  detailDialogVisible.value = true
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchData()
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.operation-log {
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

.filter-card {
  margin-bottom: 24px;
}

.table-card {
  margin-bottom: 24px;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

.json-content {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  max-height: 300px;
  overflow: auto;
  margin: 0;
}
</style>
