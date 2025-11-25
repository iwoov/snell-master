<template>
  <div class="login-container">
    <!-- Background Decoration -->
    <div class="bg-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>
    
    <!-- Login Card -->
    <el-card class="login-card slide-up">
      <!-- Header -->
      <div class="login-header">
        <h1 class="title">Snell Master</h1>
        <p class="subtitle">Admin Console</p>
      </div>
      
      <!-- Form -->
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        @keyup.enter="handleLogin"
        class="login-form"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="Username"
            :prefix-icon="User"
            size="large"
            clearable
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="Password"
            :prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            class="login-btn"
          >
            {{ loading ? 'Logging in...' : 'Login' }}
          </el-button>
        </el-form-item>
      </el-form>
      
      <!-- Footer -->
      <div class="login-footer">
        <p>© 2024 Snell Master. All rights reserved.</p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules: FormRules = {
  username: [
    { required: true, message: 'Please enter username', trigger: 'blur' },
    { min: 3, max: 20, message: 'Length should be 3 to 20', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please enter password', trigger: 'blur' },
    { min: 6, max: 50, message: 'Length should be 6 to 50', trigger: 'blur' }
  ]
}

// Handle Login
async function handleLogin() {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    
    try {
      await authStore.login(loginForm)
      ElMessage.success('Login successful')
      router.push('/')
    } catch (error: any) {
      // Password is cleared on error
      loginForm.password = ''
    } finally {
      loading.value = false
    }
  })
}

// Auto Focus / Redirect
onMounted(() => {
  if (authStore.isLoggedIn) {
    router.push('/')
  }
})
</script>

<style scoped lang="scss">
.login-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

// Background Decoration
.bg-decoration {
  position: absolute;
  inset: 0;
  overflow: hidden;
  
  .circle {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    animation: float 20s infinite ease-in-out;
  }
  
  .circle-1 {
    width: 300px;
    height: 300px;
    top: -150px;
    left: -150px;
    animation-delay: 0s;
  }
  
  .circle-2 {
    width: 200px;
    height: 200px;
    bottom: -100px;
    right: -100px;
    animation-delay: 5s;
  }
  
  .circle-3 {
    width: 150px;
    height: 150px;
    top: 50%;
    right: 10%;
    animation-delay: 10s;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

// Login Card
.login-card {
  position: relative;
  width: 420px;
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
  z-index: 1;
  
  :deep(.el-card__body) {
    padding: 0;
  }
}

// Header
.login-header {
  text-align: center;
  margin-bottom: 32px;
  
  .title {
    font-size: 32px;
    font-weight: 700;
    background: linear-gradient(135deg, #667eea, #764ba2);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    margin: 0 0 8px;
  }
  
  .subtitle {
    color: #909399;
    font-size: 14px;
    margin: 0;
  }
}

// Form
.login-form {
  .el-form-item {
    margin-bottom: 24px;
  }
  
  .login-btn {
    width: 100%;
    height: 48px;
    font-size: 16px;
    font-weight: 500;
    border-radius: 8px;
    background: linear-gradient(135deg, #667eea, #764ba2);
    border: none;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 16px rgba(102, 126, 234, 0.4);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
}

// Footer
.login-footer {
  text-align: center;
  margin-top: 24px;
  
  p {
    color: #909399;
    font-size: 12px;
    margin: 0;
  }
}

// Input Focus
:deep(.el-input__wrapper) {
  transition: all 0.3s ease;
  
  &:hover,
  &.is-focus {
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
  }
}
</style>
