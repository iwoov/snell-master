<template>
  <div class="not-found">
    <div class="content fade-in">
      <h1 class="code">404</h1>
      <h2 class="message">Page Not Found</h2>
      <p class="description">Sorry, the page you are looking for does not exist.</p>
      <el-button type="primary" @click="goHome">Back to Home</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goHome() {
  router.push('/')
}
</script>

<style scoped lang="scss">
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.content {
  text-align: center;
  
  .code {
    font-size: 120px;
    font-weight: 700;
    margin: 0;
    background: linear-gradient(135deg, #667eea, #764ba2);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
  
  .message {
    font-size: 32px;
    margin: 16px 0;
    color: #303133;
  }
  
  .description {
    font-size: 16px;
    color: #909399;
    margin: 16px 0 32px;
  }
}
</style>
