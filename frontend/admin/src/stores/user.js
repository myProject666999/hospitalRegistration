import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { adminLogin, doctorLogin } from '@/api/auth'
import { ElMessage } from 'element-plus'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('admin_userInfo') || '{}'))
  const role = ref(localStorage.getItem('admin_role') || '')

  const isLoggedIn = computed(() => !!token.value)

  async function handleAdminLogin(username, password) {
    try {
      const res = await adminLogin({ username, password })
      if (res.code === 200) {
        token.value = res.data.token
        userInfo.value = res.data.admin || {}
        role.value = 'admin'
        localStorage.setItem('admin_token', res.data.token)
        localStorage.setItem('admin_userInfo', JSON.stringify(res.data.admin || {}))
        localStorage.setItem('admin_role', 'admin')
        ElMessage.success('登录成功')
        router.push('/dashboard')
      } else {
        ElMessage.error(res.message || '登录失败')
      }
    } catch (error) {
      console.error('登录错误:', error)
      ElMessage.error('登录失败')
    }
  }

  async function handleDoctorLogin(username, password) {
    try {
      const res = await doctorLogin({ username, password })
      if (res.code === 200) {
        token.value = res.data.token
        userInfo.value = res.data.doctor || {}
        role.value = 'doctor'
        localStorage.setItem('admin_token', res.data.token)
        localStorage.setItem('admin_userInfo', JSON.stringify(res.data.doctor || {}))
        localStorage.setItem('admin_role', 'doctor')
        ElMessage.success('登录成功')
        router.push('/dashboard')
      } else {
        ElMessage.error(res.message || '登录失败')
      }
    } catch (error) {
      console.error('登录错误:', error)
      ElMessage.error('登录失败')
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = {}
    role.value = ''
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_userInfo')
    localStorage.removeItem('admin_role')
    ElMessage.success('已退出登录')
    router.push('/login')
  }

  return {
    token,
    userInfo,
    role,
    isLoggedIn,
    handleAdminLogin,
    handleDoctorLogin,
    logout
  }
})
