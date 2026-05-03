<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <el-icon :size="32" color="#409eff"><UserFilled /></el-icon>
          <span>管理后台登录</span>
        </div>
      </template>
      <el-tabs v-model="loginType" class="login-tabs">
        <el-tab-pane label="管理员登录" name="admin">
          <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" label-width="0">
            <el-form-item prop="username">
              <el-input v-model="loginForm.username" placeholder="请输入用户名" prefix-icon="User" size="large" />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                size="large"
                show-password
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" size="large" :loading="loading" class="login-btn" @click="handleLogin">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="医生登录" name="doctor">
          <el-form ref="doctorFormRef" :model="doctorForm" :rules="loginRules" label-width="0">
            <el-form-item prop="username">
              <el-input v-model="doctorForm.username" placeholder="请输入用户名" prefix-icon="User" size="large" />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="doctorForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                size="large"
                show-password
                @keyup.enter="handleDoctorLogin"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" size="large" :loading="doctorLoading" class="login-btn" @click="handleDoctorLogin">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      <div class="login-tip">
        <el-text type="info">默认管理员账号：admin / admin123</el-text>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useUserStore } from '@/stores/user'
import { UserFilled, User, Lock } from '@element-plus/icons-vue'

const userStore = useUserStore()

const loginType = ref('admin')
const loginFormRef = ref(null)
const doctorFormRef = ref(null)
const loading = ref(false)
const doctorLoading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const doctorForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  const valid = await loginFormRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  await userStore.handleAdminLogin(loginForm.username, loginForm.password)
  loading.value = false
}

const handleDoctorLogin = async () => {
  const valid = await doctorFormRef.value.validate().catch(() => false)
  if (!valid) return

  doctorLoading.value = true
  await userStore.handleDoctorLogin(doctorForm.username, doctorForm.password)
  doctorLoading.value = false
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
}

.login-card {
  width: 420px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 20px;
  font-weight: bold;
  color: #333;
}

.login-tabs {
  margin-bottom: 20px;
}

.login-btn {
  width: 100%;
}

.login-tip {
  text-align: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}
</style>
