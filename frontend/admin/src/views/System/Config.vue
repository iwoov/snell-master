<template>
  <div class="system-config-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="title">系统配置</span>
          <span class="subtitle">管理系统级配置，包括 Snell Server 版本、下载源等</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="160px"
        class="config-form"
      >
        <!-- Snell Server 配置 -->
        <div class="config-section">
          <div class="section-title">
            <el-icon><Box /></el-icon>
            <span>Snell Server 配置</span>
          </div>

          <el-form-item label="Snell 版本" prop="snell_version">
            <el-input
              v-model="formData.snell_version"
              placeholder="例如: 5.0.1"
              clearable
            >
              <template #append>
                <el-button @click="openVersionHelp">版本说明</el-button>
              </template>
            </el-input>
            <div class="form-item-tip">
              当前 Snell Server 版本号，Agent 将自动下载对应版本
            </div>
          </el-form-item>

          <el-form-item label="官方下载源">
            <el-input
              v-model="formData.snell_base_url"
              readonly
              disabled
            >
              <template #prepend>
                <el-icon><Link /></el-icon>
              </template>
            </el-input>
            <div class="form-item-tip">
              Snell Server 官方下载源地址（只读）
            </div>
          </el-form-item>

          <el-form-item label="镜像下载源" prop="snell_mirror_url">
            <el-input
              v-model="formData.snell_mirror_url"
              placeholder="可选，用于国内加速（留空则使用官方源）"
              clearable
            >
              <template #prepend>
                <el-icon><Link /></el-icon>
              </template>
            </el-input>
            <div class="form-item-tip">
              国内镜像源地址，配置后将优先使用镜像源下载 Snell Server
            </div>
          </el-form-item>
        </div>

        <!-- Agent 配置 -->
        <div class="config-section">
          <div class="section-title">
            <el-icon><Monitor /></el-icon>
            <span>Agent 配置</span>
          </div>

          <el-form-item label="Agent 版本" prop="agent_version">
            <el-input
              v-model="formData.agent_version"
              placeholder="例如: 1.0.0"
              clearable
            >
              <template #append>
                <el-button @click="openVersionHelp">版本说明</el-button>
              </template>
            </el-input>
            <div class="form-item-tip">
              Agent 版本号，用于生成部署脚本时下载对应版本的 Agent
            </div>
          </el-form-item>

          <el-form-item label="下载地址模板">
            <el-input
              v-model="formData.agent_download_url"
              type="textarea"
              :rows="2"
              placeholder="Agent 下载地址模板"
            />
            <div class="form-item-tip">
              支持变量：<el-tag size="small" type="info">{version}</el-tag> - Agent 版本号，
              <el-tag size="small" type="info">{arch}</el-tag> - 系统架构（amd64/i386/aarch64/armv7l）
            </div>
          </el-form-item>
        </div>

        <!-- Master 配置 -->
        <div class="config-section">
          <div class="section-title">
            <el-icon><Operation /></el-icon>
            <span>Master 配置</span>
          </div>

          <el-form-item label="Master URL" prop="master_url">
            <el-input
              v-model="formData.master_url"
              placeholder="例如: http://example.com:8080"
              clearable
            >
              <template #prepend>
                <el-icon><Link /></el-icon>
              </template>
            </el-input>
            <div class="form-item-tip">
              Master 服务器地址，用于生成部署脚本，Agent 将连接此地址
            </div>
          </el-form-item>
        </div>

        <!-- 端口范围配置 -->
        <div class="config-section">
          <div class="section-title">
            <el-icon><Connection /></el-icon>
            <span>端口范围配置</span>
          </div>

          <el-form-item label="默认端口起始" prop="default_port_start">
            <el-input-number
              v-model.number="portStart"
              :min="1024"
              :max="65535"
              :step="1"
              controls-position="right"
              style="width: 200px"
            />
            <div class="form-item-tip">
              实例端口分配的起始值（建议范围：10000-20000）
            </div>
          </el-form-item>

          <el-form-item label="默认端口结束" prop="default_port_end">
            <el-input-number
              v-model.number="portEnd"
              :min="1024"
              :max="65535"
              :step="1"
              controls-position="right"
              style="width: 200px"
            />
            <div class="form-item-tip">
              实例端口分配的结束值（必须大于起始值）
            </div>
          </el-form-item>
        </div>

        <!-- 操作按钮 -->
        <el-form-item>
          <div class="form-actions">
            <el-button @click="handleReset">
              <el-icon><RefreshLeft /></el-icon>
              重置
            </el-button>
            <el-button
              type="primary"
              :loading="saving"
              @click="handleSave"
            >
              <el-icon><Select /></el-icon>
              保存配置
            </el-button>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  Box,
  Link,
  Monitor,
  Operation,
  Connection,
  RefreshLeft,
  Select
} from '@element-plus/icons-vue'
import { getSystemConfigs, batchUpdateConfigs } from '@/api/system'
import type { SystemConfig } from '@/api/system'

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  snell_version: '',
  snell_base_url: '',
  snell_mirror_url: '',
  agent_version: '',
  agent_download_url: '',
  master_url: '',
  default_port_start: '',
  default_port_end: ''
})

// 端口数字值
const portStart = computed({
  get: () => parseInt(formData.default_port_start) || 10000,
  set: (val) => {
    formData.default_port_start = val.toString()
  }
})

const portEnd = computed({
  get: () => parseInt(formData.default_port_end) || 20000,
  set: (val) => {
    formData.default_port_end = val.toString()
  }
})

// 原始配置（用于重置）
const originalConfigs = ref<Record<string, string>>({})

// 保存状态
const saving = ref(false)

// URL 验证规则
const validateURL = (_rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请输入 URL'))
    return
  }
  try {
    new URL(value)
    callback()
  } catch {
    callback(new Error('请输入有效的 URL 格式'))
  }
}

// 端口范围验证
const validatePortRange = (_rule: any, _value: string, callback: any) => {
  const start = parseInt(formData.default_port_start)
  const end = parseInt(formData.default_port_end)
  
  if (start >= end) {
    callback(new Error('结束端口必须大于起始端口'))
    return
  }
  
  if (end - start < 100) {
    callback(new Error('端口范围至少需要 100 个端口'))
    return
  }
  
  callback()
}

// 表单验证规则
const rules: FormRules = {
  snell_version: [
    { required: true, message: '请输入 Snell 版本号', trigger: 'blur' },
    { pattern: /^\d+\.\d+\.\d+$/, message: '版本号格式：x.y.z', trigger: 'blur' }
  ],
  snell_mirror_url: [
    {
      validator: (rule, value, callback) => {
        if (value && value.trim()) {
          validateURL(rule, value, callback)
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  agent_version: [
    { required: true, message: '请输入 Agent 版本号', trigger: 'blur' },
    { pattern: /^\d+\.\d+\.\d+$/, message: '版本号格式：x.y.z', trigger: 'blur' }
  ],
  master_url: [
    { required: true, message: '请输入 Master URL', trigger: 'blur' },
    { validator: validateURL, trigger: 'blur' }
  ],
  default_port_start: [
    { required: true, message: '请输入起始端口', trigger: 'blur' }
  ],
  default_port_end: [
    { required: true, message: '请输入结束端口', trigger: 'blur' },
    { validator: validatePortRange, trigger: 'blur' }
  ]
}

// 加载配置
const loadConfigs = async () => {
  try {
    const res = await getSystemConfigs()
    const configs = res || []
    
    // 转换为 key-value 对象
    const configMap: Record<string, string> = {}
    configs.forEach((item: SystemConfig) => {
      configMap[item.key] = item.value
      formData[item.key as keyof typeof formData] = item.value
    })
    
    // 保存原始配置
    originalConfigs.value = { ...configMap }
    
    ElMessage.success('配置加载成功')
  } catch (error: any) {
    ElMessage.error(error.message || '加载配置失败')
  }
}

// 保存配置
const handleSave = async () => {
  if (!formRef.value) return
  
  try {
    // 验证表单
    await formRef.value.validate()
    
    // 确认保存
    await ElMessageBox.confirm(
      '确定要保存系统配置吗？配置将立即生效。',
      '确认保存',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    saving.value = true
    
    // 批量更新
    const updates: Record<string, string> = {
      snell_version: formData.snell_version,
      snell_base_url: formData.snell_base_url,
      snell_mirror_url: formData.snell_mirror_url,
      agent_version: formData.agent_version,
      agent_download_url: formData.agent_download_url,
      master_url: formData.master_url,
      default_port_start: formData.default_port_start,
      default_port_end: formData.default_port_end
    }
    
    await batchUpdateConfigs(updates)
    
    ElMessage.success('配置保存成功')
    
    // 重新加载配置
    await loadConfigs()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '保存配置失败')
    }
  } finally {
    saving.value = false
  }
}

// 重置配置
const handleReset = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重置配置吗？将恢复到最后保存的状态。',
      '确认重置',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 恢复原始配置
    Object.keys(originalConfigs.value).forEach((key) => {
      const value = originalConfigs.value[key]
      if (value !== undefined) {
        formData[key as keyof typeof formData] = value
      }
    })
    
    // 清除验证错误
    formRef.value?.clearValidate()
    
    ElMessage.success('已重置到最后保存的状态')
  } catch (error) {
    // 用户取消
  }
}

// 打开版本说明
const openVersionHelp = () => {
  ElMessageBox.alert(
    '版本号格式为 x.y.z，例如：5.0.1\n\n' +
    '• x: 主版本号（重大更新）\n' +
    '• y: 次版本号（功能更新）\n' +
    '• z: 修订号（bug 修复）\n\n' +
    '修改版本号后，新部署的 Agent 将自动下载对应版本的程序。\n' +
    '已部署的 Agent 需要重新部署才会更新。',
    '版本号说明',
    {
      confirmButtonText: '知道了',
      type: 'info'
    }
  )
}

// 组件挂载时加载配置
onMounted(() => {
  loadConfigs()
})
</script>

<style scoped lang="scss">
.system-config-container {
  padding: 20px;
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
  
  .title {
    font-size: 18px;
    font-weight: 600;
    color: var(--el-text-color-primary);
  }
  
  .subtitle {
    font-size: 13px;
    color: var(--el-text-color-secondary);
  }
}

.config-form {
  max-width: 800px;
  margin-top: 20px;
}

.config-section {
  margin-bottom: 40px;
  
  &:last-of-type {
    margin-bottom: 20px;
  }
  
  .section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 20px;
    padding-bottom: 12px;
    border-bottom: 2px solid var(--el-border-color-lighter);
    font-size: 16px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    
    .el-icon {
      font-size: 18px;
      color: var(--el-color-primary);
    }
  }
}

.form-item-tip {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--el-text-color-secondary);
  
  .el-tag {
    margin: 0 2px;
    font-family: 'Courier New', monospace;
  }
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid var(--el-border-color-lighter);
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-input-number) {
  .el-input__inner {
    text-align: left;
  }
}
</style>
