<template>
  <div class="subscription-management">
    <div class="header">
      <h1>订阅管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>新建订阅
      </el-button>
    </div>

    <!-- Subscription Table -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="subscriptionList" :style="{ width: '100%' }">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户" min-width="120">
          <template #default="{ row }">
            {{ row.user?.username || '未知用户' }}
          </template>
        </el-table-column>
        <el-table-column label="订阅令牌" min-width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.token" placement="top">
              <span class="token-mask">{{ row.token.substring(0, 16) }}...</span>
            </el-tooltip>
            <el-button link size="small" @click="copyText(row.token)">
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="订阅链接" min-width="200">
          <template #default="{ row }">
            <el-button link size="small" @click="copySubscriptionUrl(row.token)">
              <el-icon><Link /></el-icon>复制链接
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="access_count" label="访问次数" width="100" align="center" />
        <el-table-column prop="last_access_at" label="最后访问" width="180">
          <template #default="{ row }">
            {{ row.last_access_at ? formatDate(row.last_access_at) : '从未访问' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="warning" @click="handleRegenerate(row)">重新生成</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Create Dialog -->
    <el-dialog
      v-model="dialogVisible"
      title="新建订阅"
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
      title="新的订阅令牌"
      width="500px"
    >
      <div class="token-content">
        <p>请复制并保存新的订阅令牌，它将只显示一次。</p>
        <el-input v-model="newToken" readonly>
          <template #append>
            <el-button @click="copyText(newToken)">
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </template>
        </el-input>
        <div style="margin-top: 16px">
          <p>订阅链接：</p>
          <el-input v-model="subscriptionUrl" readonly>
            <template #append>
              <el-button @click="copyText(subscriptionUrl)">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { Plus, CopyDocument, Link } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'
import { useClipboard } from '@vueuse/core'
import {
  getSubscriptionList,
  createSubscription,
  deleteSubscription,
  regenerateSubscriptionToken
} from '@/api/subscription'
import type { Subscription } from '@/api/subscription'
import { getUserList } from '@/api/user'
import type { UserInfo } from '@/types/api'

// State
const loading = ref(false)
const submitting = ref(false)
const subscriptionList = ref<Subscription[]>([])

// Options State
const userLoading = ref(false)
const userOptions = ref<UserInfo[]>([])

// Dialog State
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({
  user_id: undefined as number | undefined
})

// Token Dialog State
const tokenDialogVisible = ref(false)
const newToken = ref('')
const { copy } = useClipboard()

// Validation Rules
const rules = reactive<FormRules>({
  user_id: [{ required: true, message: '请选择用户', trigger: 'change' }]
})

// Computed
const subscriptionUrl = computed(() => {
  if (!newToken.value) return ''
  // Assuming the subscription endpoint is at /api/subscribe/:token
  return `${window.location.origin}/api/subscribe/${newToken.value}`
})

// Helpers
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const copyText = async (text: string) => {
  await copy(text)
  ElMessage.success('已复制到剪贴板')
}

const copySubscriptionUrl = async (token: string) => {
  const url = `${window.location.origin}/api/subscribe/${token}`
  await copy(url)
  ElMessage.success('订阅链接已复制到剪贴板')
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

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getSubscriptionList()
    subscriptionList.value = res
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// Actions
const handleCreate = () => {
  form.user_id = undefined
  dialogVisible.value = true
  // Load initial users if empty
  if (userOptions.value.length === 0) {
    searchUsers('')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid && form.user_id) {
      submitting.value = true
      try {
        const res = await createSubscription({
          user_id: form.user_id
        })
        newToken.value = res.token
        tokenDialogVisible.value = true
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

const handleDelete = (row: Subscription) => {
  ElMessageBox.confirm(
    `确定要删除该订阅吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteSubscription(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleRegenerate = (row: Subscription) => {
  ElMessageBox.confirm(
    `确定要重新生成订阅令牌吗？旧的订阅链接将失效。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const res = await regenerateSubscriptionToken(row.id)
      newToken.value = res.token
      tokenDialogVisible.value = true
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.subscription-management {
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

.token-mask {
  font-family: monospace;
  margin-right: 8px;
  color: #606266;
}

.token-content {
  text-align: left;
}

.token-content p {
  margin-bottom: 12px;
  color: #606266;
}
</style>
