<template>
  <div class="template-management">
    <div class="header">
      <h1>模板管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>新建模板
      </el-button>
    </div>

    <!-- Template Table -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="templateList" :style="{ width: '100%' }">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="模板名称" min-width="180" />
        <el-table-column prop="description" label="描述" min-width="200">
          <template #default="{ row }">
            {{ row.description || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="默认模板" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="success">是</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button 
              v-if="!row.is_default" 
              link 
              type="success" 
              @click="handleSetDefault(row)"
            >
              设为默认
            </el-button>
            <el-button 
              link 
              type="danger" 
              @click="handleDelete(row)"
              :disabled="row.is_default"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新建模板' : '编辑模板'"
      width="800px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="form.name" placeholder="例如：Surge 默认模板" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="form.description" 
            type="textarea" 
            :rows="2"
            placeholder="简要描述此模板的用途"
          />
        </el-form-item>
        <el-form-item label="模板内容" prop="content">
          <div class="editor-container">
            <div class="editor-toolbar">
              <el-text type="info" size="small" v-text="'支持的变量：{{SERVER_LIST}}, {{USER_AGENT}}, {{UPDATE_INTERVAL}}'" />
            </div>
            <el-input
              v-model="form.content"
              type="textarea"
              :rows="15"
              placeholder="输入模板内容..."
              class="template-editor"
            />
          </div>
        </el-form-item>
        <el-form-item label="默认模板" v-if="dialogType === 'create'">
          <el-switch v-model="form.is_default" />
          <el-text type="info" size="small" style="margin-left: 12px">
            设为默认后，新用户将自动使用此模板
          </el-text>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'
import {
  getTemplateList,
  createTemplate,
  updateTemplate,
  deleteTemplate,
  setDefaultTemplate
} from '@/api/template'
import type { Template } from '@/api/template'

// State
const loading = ref(false)
const submitting = ref(false)
const templateList = ref<Template[]>([])

// Dialog State
const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const form = reactive({
  id: 0,
  name: '',
  content: '',
  description: '',
  is_default: false
})

// Validation Rules
const rules = reactive<FormRules>({
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入模板内容', trigger: 'blur' }
  ]
})

// Helpers
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// Data Fetching
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getTemplateList()
    templateList.value = res
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// Actions
const handleCreate = () => {
  dialogType.value = 'create'
  form.id = 0
  form.name = ''
  form.content = ''
  form.description = ''
  form.is_default = false
  dialogVisible.value = true
}

const handleEdit = (row: Template) => {
  dialogType.value = 'edit'
  form.id = row.id
  form.name = row.name
  form.content = row.content
  form.description = row.description || ''
  form.is_default = row.is_default
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (dialogType.value === 'create') {
          await createTemplate({
            name: form.name,
            content: form.content,
            description: form.description,
            is_default: form.is_default
          })
          ElMessage.success('创建成功')
        } else {
          await updateTemplate(form.id, {
            name: form.name,
            content: form.content,
            description: form.description
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

const handleDelete = (row: Template) => {
  if (row.is_default) {
    ElMessage.warning('默认模板不能删除')
    return
  }
  
  ElMessageBox.confirm(
    `确定要删除模板 "${row.name}" 吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteTemplate(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleSetDefault = (row: Template) => {
  ElMessageBox.confirm(
    `确定要将 "${row.name}" 设为默认模板吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).then(async () => {
    try {
      await setDefaultTemplate(row.id)
      ElMessage.success('设置成功')
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
.template-management {
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

.editor-container {
  width: 100%;
}

.editor-toolbar {
  padding: 8px 12px;
  background: #f5f7fa;
  border: 1px solid #dcdfe6;
  border-bottom: none;
  border-radius: 4px 4px 0 0;
}

.template-editor {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

.template-editor :deep(textarea) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  border-radius: 0 0 4px 4px;
}
</style>
