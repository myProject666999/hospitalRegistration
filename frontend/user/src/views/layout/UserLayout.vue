<template>
  <el-container>
    <el-header>
      <div class="header-content">
        <div class="logo" @click="router.push('/home')">
          <el-icon size="32"><Hospital /></el-icon>
          <span>医院预约挂号系统</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="nav-menu"
          mode="horizontal"
          router
        >
          <el-menu-item index="/home">首页</el-menu-item>
          <el-menu-item index="/doctors">医生列表</el-menu-item>
          <el-menu-item index="/medicines">药品信息</el-menu-item>
          <el-menu-item index="/medical-projects">诊疗项目</el-menu-item>
          <el-menu-item index="/announcements">公告信息</el-menu-item>
        </el-menu>
        <div class="user-info">
          <template v-if="userStore.isLoggedIn">
            <el-dropdown @command="handleCommand">
              <span class="dropdown-trigger">
                <el-avatar :size="32" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
                <span class="username">{{ userStore.userInfo.real_name || userStore.userInfo.username }}</span>
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                  <el-dropdown-item command="appointments">预约管理</el-dropdown-item>
                  <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" link @click="router.push('/login')">登录</el-button>
            <el-button type="primary" @click="router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
    <el-footer>
      <div class="footer-content">
        <p>© 2024 医院预约挂号系统 - 版权所有</p>
        <p>联系电话：400-123-4567 | 工作时间：周一至周日 8:00-20:00</p>
      </div>
    </el-footer>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Hospital, ArrowDown } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => {
  if (route.path.startsWith('/user')) {
    return '/user'
  }
  return route.path
})

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/user/profile')
      break
    case 'appointments':
      router.push('/user/appointments')
      break
    case 'favorites':
      router.push('/user/favorites')
      break
    case 'logout':
      userStore.logout()
      break
  }
}
</script>

<style scoped>
.el-header {
  background-color: #409eff;
  color: white;
  padding: 0;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: bold;
  cursor: pointer;
}

.nav-menu {
  background-color: transparent;
  border-bottom: none;
}

.nav-menu .el-menu-item {
  color: white;
  border-bottom: 2px solid transparent;
}

.nav-menu .el-menu-item:hover,
.nav-menu .el-menu-item.is-active {
  background-color: rgba(255, 255, 255, 0.1);
  border-bottom-color: white;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.dropdown-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  color: white;
}

.el-footer {
  background-color: #f5f7fa;
  color: #909399;
  padding: 20px 0;
}

.footer-content {
  max-width: 1400px;
  margin: 0 auto;
  text-align: center;
}

.footer-content p {
  margin: 5px 0;
  font-size: 14px;
}

.el-main {
  background-color: #f5f7fa;
  padding: 20px;
  min-height: calc(100vh - 120px - 100px);
}
</style>
