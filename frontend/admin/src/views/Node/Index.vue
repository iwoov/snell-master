<template>
  <div class="node-management">
    <div class="header">
      <h1>节点管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>新建节点
      </el-button>
    </div>

    <!-- Node Table -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="nodeList" :style="{ width: '100%' }">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column label="地址" min-width="180">
          <template #default="{ row }">
            <div>{{ row.endpoint }}</div>
            <div class="sub-text">{{ row.location || '未知位置' }}</div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'">
              {{ row.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="负载 (CPU/Mem)" width="180">
          <template #default="{ row }">
            <div class="usage-info">
              <div class="usage-item">
                <span>CPU</span>
                <el-progress :percentage="Math.round(row.cpu_usage)" :stroke-width="6" />
              </div>
              <div class="usage-item">
                <span>Mem</span>
                <el-progress :percentage="Math.round(row.memory_usage)" :stroke-width="6" status="warning" />
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="instance_count" label="实例数" width="100" align="center" />
        <el-table-column prop="last_seen_at" label="最后在线" width="180">
          <template #default="{ row }">
            {{ row.last_seen_at ? formatDate(row.last_seen_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="success" @click="handleDownloadScript(row)">下载安装脚本</el-button>
            <el-button link type="warning" @click="handleRegenerateToken(row)">重置Token</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新建节点' : '编辑节点'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="节点名称" />
        </el-form-item>
        <el-form-item label="地址" prop="endpoint">
          <el-input v-model="form.endpoint" placeholder="IP 或域名" />
        </el-form-item>
        <el-form-item label="位置" prop="location">
          <el-input v-model="form.location" placeholder="例如：香港" />
        </el-form-item>
        <el-form-item label="国家代码" prop="country_code">
          <el-input v-model="form.country_code" placeholder="例如：HK" maxlength="2" />
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

    <!-- Token Dialog -->
    <el-dialog
      v-model="tokenDialogVisible"
      title="新的 API Token"
      width="400px"
    >
      <div class="token-content">
        <p>请复制并保存新的 API Token，它将只显示一次。</p>
        <el-input v-model="newToken" readonly>
          <template #append>
            <el-button @click="copyToken">
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </template>
        </el-input>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="tokenDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, CopyDocument } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'
import { useClipboard } from '@vueuse/core'
import {
  getNodeList,
  createNode,
  updateNode,
  deleteNode,
  regenerateToken,
  downloadInstallScript
} from '@/api/node'
import type { Node } from '@/api/node'

// State
const loading = ref(false)
const submitting = ref(false)
const nodeList = ref<Node[]>([])

// Dialog State
const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const form = reactive({
  id: 0,
  name: '',
  endpoint: '',
  location: '',
  country_code: ''
})

// Token Dialog State
const tokenDialogVisible = ref(false)
const newToken = ref('')
const { copy } = useClipboard()

// Validation Rules
const rules = reactive<FormRules>({
  name: [
    { required: true, message: '请输入节点名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  endpoint: [
    { required: true, message: '请输入节点地址', trigger: 'blur' }
  ]
})

// Helpers
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const copyToken = async () => {
  await copy(newToken.value)
  ElMessage.success('已复制到剪贴板')
}

// Actions
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getNodeList()
    nodeList.value = res
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  dialogType.value = 'create'
  form.name = ''
  form.endpoint = ''
  form.location = ''
  form.country_code = ''
  dialogVisible.value = true
}

const handleEdit = (row: Node) => {
  dialogType.value = 'edit'
  form.id = row.id
  form.name = row.name
  form.endpoint = row.endpoint
  form.location = row.location || ''
  form.country_code = row.country_code || ''
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (dialogType.value === 'create') {
          await createNode({
            name: form.name,
            endpoint: form.endpoint,
            location: form.location,
            country_code: form.country_code
          })
          ElMessage.success('创建成功')
        } else {
          await updateNode(form.id, {
            name: form.name,
            endpoint: form.endpoint,
            location: form.location,
            country_code: form.country_code
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

const handleDelete = (row: Node) => {
  ElMessageBox.confirm(
    `确定要删除节点 "${row.name}" 吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteNode(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleRegenerateToken = (row: Node) => {
  ElMessageBox.confirm(
    `确定要重置节点 "${row.name}" 的 API Token 吗？这将导致节点暂时离线，直到更新配置。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const res = await regenerateToken(row.id)
      newToken.value = res.token
      tokenDialogVisible.value = true
    } catch (error) {
      console.error(error)
    }
  })
}

const handleDownloadScript = async (row: Node) => {
  try {
    await downloadInstallScript(row.id, row.name)
    ElMessage.success('安装脚本下载成功')
  } catch (error) {
    console.error(error)
    ElMessage.error('下载失败，请重试')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.node-management {
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

.table-card {
  margin-bottom: 24px;
}

.sub-text {
  font-size: 12px;
  color: #909399;
}

.usage-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.usage-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.usage-item span {
  width: 30px;
  color: #909399;
}

.usage-item .el-progress {
  flex: 1;
}

.token-content {
  text-align: center;
}

.token-content p {
  margin-bottom: 16px;
  color: #606266;
}
</style>
