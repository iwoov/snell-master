<template>
  <div class="instance-management">
    <div class="header">
      <h1>实例管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>新建实例
      </el-button>
    </div>

    <!-- Search & Filter -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="用户">
          <el-select
            v-model="filterForm.user_id"
            placeholder="选择用户"
            clearable
            filterable
            remote
            :remote-method="searchUsers"
            :loading="userLoading"
            style="width: 200px"
          >
            <el-option
              v-for="item in userOptions"
              :key="item.id"
              :label="item.username"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="节点">
          <el-select
            v-model="filterForm.node_id"
            placeholder="选择节点"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="item in nodeOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="运行中" :value="1" />
            <el-option label="已停止" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Instance Table -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="instanceList" :style="{ width: '100%' }">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户" min-width="120">
          <template #default="{ row }">
            {{ row.user?.username || '未知用户' }}
          </template>
        </el-table-column>
        <el-table-column label="节点" min-width="150">
          <template #default="{ row }">
            {{ row.node?.name || '未知节点' }}
          </template>
        </el-table-column>
        <el-table-column prop="port" label="端口" width="100" />
        <el-table-column label="密码" width="180">
          <template #default="{ row }">
            <el-tooltip :content="row.psk" placement="top">
              <span class="password-mask">******</span>
            </el-tooltip>
            <el-button link size="small" @click="copyText(row.psk)">
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '运行中' : '已停止' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleViewConfig(row)">配置</el-button>
            <el-button 
              link 
              :type="row.status === 1 ? 'warning' : 'success'" 
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '停止' : '启动' }}
            </el-button>
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

    <!-- Create Dialog -->
    <el-dialog
      v-model="dialogVisible"
      title="新建实例"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="用户" prop="user_id">
          <el-select
            v-model="form.user_id"
            placeholder="选择用户"
            filterable
            remote
            :remote-method="searchUsers"
            :loading="userLoading"
            style="width: 100%"
          >
            <el-option
              v-for="item in userOptions"
              :key="item.id"
              :label="item.username"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="节点" prop="node_id">
          <el-select
            v-model="form.node_id"
            placeholder="选择节点"
            style="width: 100%"
          >
            <el-option
              v-for="item in nodeOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input v-model.number="form.port" placeholder="留空自动分配">
            <template #append>
              <el-tooltip content="留空则由系统自动分配可用端口">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="密码" prop="psk">
          <el-input v-model="form.psk" placeholder="留空自动生成">
            <template #append>
              <el-button @click="generatePSK">生成</el-button>
            </template>
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

    <!-- Config Dialog -->
    <el-dialog
      v-model="configDialogVisible"
      title="实例配置"
      width="600px"
    >
      <div class="config-content">
        <el-input
          v-model="configContent"
          type="textarea"
          :rows="10"
          readonly
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="copyText(configContent)">复制配置</el-button>
          <el-button type="primary" @click="configDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, CopyDocument, InfoFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'
import { useClipboard } from '@vueuse/core'
import {
  getInstanceList,
  createInstance,
  deleteInstance,
  updateInstanceStatus,
  getInstanceConfig
} from '@/api/instance'
import type { Instance } from '@/api/instance'
import { getUserList } from '@/api/user'
import { getNodeList } from '@/api/node'
import type { UserInfo } from '@/types/api'
import type { Node } from '@/api/node'

// State
const loading = ref(false)
const submitting = ref(false)
const instanceList = ref<Instance[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const filterForm = reactive({
  user_id: undefined as number | undefined,
  node_id: undefined as number | undefined,
  status: undefined as number | undefined
})

// Options State
const userLoading = ref(false)
const userOptions = ref<UserInfo[]>([])
const nodeOptions = ref<Node[]>([])

// Dialog State
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({
  user_id: undefined as number | undefined,
  node_id: undefined as number | undefined,
  port: undefined as number | undefined,
  psk: ''
})

// Config Dialog State
const configDialogVisible = ref(false)
const configContent = ref('')
const { copy } = useClipboard()

// Validation Rules
const rules = reactive<FormRules>({
  user_id: [{ required: true, message: '请选择用户', trigger: 'change' }],
  node_id: [{ required: true, message: '请选择节点', trigger: 'change' }]
})

// Helpers
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const copyText = async (text: string) => {
  await copy(text)
  ElMessage.success('已复制到剪贴板')
}

const generatePSK = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < 16; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.psk = result
}

// Data Fetching
const searchUsers = async (query: string) => {
  userLoading.value = true
  try {
    const res = await getUserList({ keyword: query, page: 1, page_size: 20 })
    userOptions.value = res.items
  } catch (error) {
    console.error(error)
  } finally {
    userLoading.value = false
  }
}

const fetchNodes = async () => {
  try {
    const res = await getNodeList()
    nodeOptions.value = res
  } catch (error) {
    console.error(error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getInstanceList({
      page: currentPage.value,
      page_size: pageSize.value,
      user_id: filterForm.user_id,
      node_id: filterForm.node_id,
      status: filterForm.status
    })
    instanceList.value = res.items
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
  filterForm.user_id = undefined
  filterForm.node_id = undefined
  filterForm.status = undefined
  handleSearch()
}

const handleCreate = () => {
  form.user_id = undefined
  form.node_id = undefined
  form.port = undefined
  form.psk = ''
  dialogVisible.value = true
  // Load initial users if empty
  if (userOptions.value.length === 0) {
    searchUsers('')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid && form.user_id && form.node_id) {
      submitting.value = true
      try {
        await createInstance({
          user_id: form.user_id,
          node_id: form.node_id,
          port: form.port,
          psk: form.psk || undefined
        })
        ElMessage.success('创建成功')
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

const handleDelete = (row: Instance) => {
  ElMessageBox.confirm(
    `确定要删除该实例吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteInstance(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleToggleStatus = (row: Instance) => {
  const action = row.status === 1 ? '停止' : '启动'
  const newStatus = row.status === 1 ? 0 : 1
  
  ElMessageBox.confirm(
    `确定要${action}该实例吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).then(async () => {
    try {
      await updateInstanceStatus(row.id, newStatus)
      ElMessage.success(`实例已${action}`)
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleViewConfig = async (row: Instance) => {
  try {
    const res = await getInstanceConfig(row.id)
    configContent.value = res.config
    configDialogVisible.value = true
  } catch (error) {
    console.error(error)
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
  fetchNodes()
})
</script>

<style scoped>
.instance-management {
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

.password-mask {
  margin-right: 8px;
  color: #909399;
}

.config-content {
  margin-bottom: 16px;
}
</style>
