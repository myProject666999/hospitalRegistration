import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, register, getUserInfo } from '@/api/user'
import { ElMessage } from 'element-plus'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  const isLoggedIn = computed(() => !!token.value)

  async function handleLogin(username, password) {
    try {
      const res = await login({ username, password })
      if (res.code === 200) {
        token.value = res.data.token
        userInfo.value = res.data.user || {}
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('userInfo', JSON.stringify(res.data.user || {}))
        ElMessage.success('登录成功')
        router.push('/home')
      } else {
        ElMessage.error(res.message || '登录失败')
      }
    } catch (error) {
      console.error('登录错误:', error)
      ElMessage.error('登录失败')
    }
  }

  async function handleRegister(data) {
    try {
      const res = await register(data)
      if (res.code === 200) {
        ElMessage.success('注册成功，请登录')
        router.push('/login')
      } else {
        ElMessage.error(res.message)
      }
    } catch (error) {
      ElMessage.error('注册失败')
    }
  }

  async function fetchUserInfo() {
    try {
      const res = await getUserInfo()
      if (res.code === 200) {
        userInfo.value = res.data
        localStorage.setItem('userInfo', JSON.stringify(res.data))
      }
    } catch (error) {
      console.error('获取用户信息失败', error)
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    ElMessage.success('已退出登录')
    router.push('/login')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    handleLogin,
    handleRegister,
    fetchUserInfo,
    logout
  }
})
