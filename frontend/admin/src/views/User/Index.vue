<template>
  <div class="user-management">
    <div class="header">
      <h1>用户管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>新建用户
      </el-button>
    </div>

    <!-- Search & Filter -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="关键词">
          <el-input
            v-model="filterForm.keyword"
            placeholder="用户名 / 邮箱"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- User Table -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="userList">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column label="流量使用" min-width="200">
          <template #default="{ row }">
            <div class="traffic-info">
              <el-progress
                :percentage="calculatePercentage(row.traffic_used_total, row.traffic_limit)"
                :status="getTrafficStatus(row.traffic_used_total, row.traffic_limit)"
              />
              <span class="traffic-text">
                {{ formatBytes(row.traffic_used_total) }} / {{ formatBytes(row.traffic_limit) }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="(val: string | number | boolean) => handleStatusChange(row, val as number)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="warning" @click="handleResetTraffic(row)">重置流量</el-button>
            <el-button link type="success" @click="handleAssignNodes(row)">分配节点</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
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

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新建用户' : '编辑用户'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="dialogType === 'edit'" />
        </el-form-item>
        <el-form-item
          label="密码"
          prop="password"
          :rules="dialogType === 'create' ? rules.password : []"
        >
          <el-input
            v-model="form.password"
            type="password"
            show-password
            :placeholder="dialogType === 'edit' ? '留空则不修改' : ''"
          />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="流量限制" prop="traffic_limit">
          <el-input v-model.number="form.traffic_limit" type="number">
            <template #append>GB</template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Node Assignment Dialog -->
    <el-dialog
      v-model="nodeDialogVisible"
      title="分配节点"
      width="600px"
    >
      <div class="node-transfer">
        <el-transfer
          v-model="selectedNodes"
          :data="nodeList"
          :titles="['可选节点', '已选节点']"
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="nodeDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleNodeSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'
import {
  getUserList,
  createUser,
  updateUser,
  deleteUser,
  resetUserTraffic,
  updateUserStatus,
  assignNodes
} from '@/api/user'
import type { UserInfo } from '@/types/api'

// State
const loading = ref(false)
const submitting = ref(false)
const userList = ref<UserInfo[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const filterForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

// Dialog State
const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const form = reactive({
  id: 0,
  username: '',
  password: '',
  email: '',
  traffic_limit: 100 // Default 100GB
})

// Node Dialog State
const nodeDialogVisible = ref(false)
const selectedNodes = ref<number[]>([])
const nodeList = ref<any[]>([]) // TODO: Fetch from Node API
const currentUserId = ref(0)

// Validation Rules
const rules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  traffic_limit: [
    { required: true, message: '请输入流量限制', trigger: 'blur' },
    { type: 'number', min: 1, message: '流量限制必须大于 0', trigger: 'blur' }
  ]
})

// Formatting Helpers
const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const calculatePercentage = (used: number, limit: number) => {
  if (limit === 0) return 0
  return Math.min(Math.round((used / limit) * 100), 100)
}

const getTrafficStatus = (used: number, limit: number) => {
  const percentage = calculatePercentage(used, limit)
  if (percentage >= 90) return 'exception'
  if (percentage >= 80) return 'warning'
  return 'success'
}

// Actions
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getUserList({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: filterForm.keyword,
      status: filterForm.status
    })
    userList.value = res.items
    total.value = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

const resetFilter = () => {
  filterForm.keyword = ''
  filterForm.status = undefined
  handleSearch()
}

const handleCreate = () => {
  dialogType.value = 'create'
  form.username = ''
  form.password = ''
  form.email = ''
  form.traffic_limit = 100
  dialogVisible.value = true
}

const handleEdit = (row: UserInfo) => {
  dialogType.value = 'edit'
  form.id = row.id
  form.username = row.username
  form.password = '' // Don't show password
  form.email = row.email
  form.traffic_limit = Math.round(row.traffic_limit / 1024 / 1024 / 1024) // Convert to GB
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const trafficLimitBytes = form.traffic_limit * 1024 * 1024 * 1024
        if (dialogType.value === 'create') {
          await createUser({
            username: form.username,
            password: form.password,
            email: form.email,
            traffic_limit: trafficLimitBytes
          })
          ElMessage.success('创建成功')
        } else {
          await updateUser(form.id, {
            email: form.email,
            password: form.password || undefined,
            traffic_limit: trafficLimitBytes
          })
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        console.error(error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDelete = (row: UserInfo) => {
  ElMessageBox.confirm(
    `确定要删除用户 "${row.username}" 吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteUser(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleResetTraffic = (row: UserInfo) => {
  ElMessageBox.confirm(
    `确定要重置用户 "${row.username}" 的流量吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).then(async () => {
    try {
      await resetUserTraffic(row.id)
      ElMessage.success('流量已重置')
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleStatusChange = async (row: UserInfo, status: number) => {
  try {
    await updateUserStatus(row.id, status)
    ElMessage.success(status === 1 ? '用户已启用' : '用户已禁用')
  } catch (error) {
    row.status = status === 1 ? 0 : 1 // Revert on error
    console.error(error)
  }
}

const handleAssignNodes = (row: UserInfo) => {
  currentUserId.value = row.id
  // TODO: Load assigned nodes and all nodes
  // For now, just show dialog with empty list
  nodeDialogVisible.value = true
}

const handleNodeSubmit = async () => {
  submitting.value = true
  try {
    await assignNodes(currentUserId.value, selectedNodes.value)
    ElMessage.success('节点分配成功')
    nodeDialogVisible.value = false
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
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
.user-management {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.traffic-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.traffic-text {
  font-size: 12px;
  color: #909399;
}
</style>
