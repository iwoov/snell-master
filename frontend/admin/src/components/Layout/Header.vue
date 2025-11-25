<template>
  <div class="header-container">
    <!-- 左侧 -->
    <div class="header-left">
      <el-button
        class="collapse-btn"
        :icon="isCollapsed ? Expand : Fold"
        text
        @click="handleToggleSidebar"
      />
    </div>

    <!-- 右侧 -->
    <div class="header-right">
      <el-dropdown @command="handleCommand">
        <div class="user-info">
          <el-avatar :size="32" class="user-avatar">
            {{ userInitial }}
          </el-avatar>
          <span class="username">{{ authStore.userInfo?.username }}</span>
          <el-icon class="arrow-icon"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item disabled>
              <div class="user-menu-info">
                <div class="user-menu-username">{{ authStore.userInfo?.username }}</div>
                <div class="user-menu-role">管理员</div>
              </div>
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              <span>退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Fold, Expand, ArrowDown, SwitchButton } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

interface Props {
  isCollapsed: boolean
}

interface Emits {
  (e: 'toggle-sidebar'): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const router = useRouter()
const authStore = useAuthStore()

// 用户名首字母
const userInitial = computed(() => {
  const username = authStore.userInfo?.username || 'A'
  return username.charAt(0).toUpperCase()
})

// 切换侧边栏
const handleToggleSidebar = () => {
  emit('toggle-sidebar')
}

// 处理下拉菜单命令
const handleCommand = async (command: string) => {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      authStore.logout()
      ElMessage.success('已退出登录')
      router.push('/login')
    } catch {
      // 用户取消
    }
  }
}
</script>

<style scoped>
.header-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  font-size: 20px;
  color: #000000d9;
  transition: color 0.3s ease;
}

.collapse-btn:hover {
  color: #1890ff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 4px;
  transition: background-color 0.3s ease;
}

.user-info:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.user-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-weight: 600;
}

.username {
  font-size: 14px;
  color: #000000d9;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.arrow-icon {
  font-size: 12px;
  color: #00000073;
  transition: transform 0.3s ease;
}

.user-info:hover .arrow-icon {
  transform: translateY(2px);
}

.user-menu-info {
  padding: 4px 0;
}

.user-menu-username {
  font-size: 14px;
  font-weight: 600;
  color: #000000d9;
}

.user-menu-role {
  font-size: 12px;
  color: #00000073;
  margin-top: 4px;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-dropdown-menu__item:not(.is-disabled):hover) {
  background-color: rgba(24, 144, 255, 0.08);
  color: #1890ff;
}
</style>
