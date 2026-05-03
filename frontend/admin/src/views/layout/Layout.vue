<template>
  <el-container class="admin-layout">
    <el-aside :width="isCollapse ? '64px' : '220px'" class="admin-aside">
      <div class="logo">
        <el-icon v-if="isCollapse" :size="28"><Hospital /></el-icon>
        <template v-else>
          <el-icon :size="28"><Hospital /></el-icon>
          <span>医院管理后台</span>
        </template>
      </div>
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :collapse-transition="false"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>
          <template #title>数据统计</template>
        </el-menu-item>
        <el-sub-menu index="appointment">
          <template #title>
            <el-icon><Calendar /></el-icon>
            <span>预约管理</span>
          </template>
          <el-menu-item index="/appointments">预约列表</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="consultation">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>就诊管理</span>
          </template>
          <el-menu-item index="/consultations">就诊列表</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="doctor" v-if="isAdmin">
          <template #title>
            <el-icon><User /></el-icon>
            <span>医生管理</span>
          </template>
          <el-menu-item index="/doctors">医生列表</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="user" v-if="isAdmin">
          <template #title>
            <el-icon><UserFilled /></el-icon>
            <span>用户管理</span>
          </template>
          <el-menu-item index="/users">用户列表</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="system" v-if="isAdmin">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/departments">科室管理</el-menu-item>
          <el-menu-item index="/hospitals">医院管理</el-menu-item>
          <el-menu-item index="/medicines">药品管理</el-menu-item>
          <el-menu-item index="/medical-projects">诊疗项目</el-menu-item>
          <el-menu-item index="/schedules">排班管理</el-menu-item>
          <el-menu-item index="/carousels">轮播图管理</el-menu-item>
          <el-menu-item index="/announcements">公告管理</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="admin-header">
        <div class="header-left">
          <el-icon class="collapse-icon" @click="toggleCollapse">
            <component :is="isCollapse ? 'Expand' : 'Fold'" />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentTitle">{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-icon><User /></el-icon>
              <span>{{ userInfo.name || userInfo.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { 
  Hospital, 
  DataLine, 
  Calendar, 
  Document, 
  User, 
  UserFilled, 
  Setting,
  Fold,
  Expand,
  ArrowDown
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const userInfo = computed(() => userStore.userInfo)
const isAdmin = computed(() => userStore.role === 'admin')

const activeMenu = computed(() => {
  return route.path
})

const currentTitle = computed(() => {
  const metaTitle = route.meta?.title
  if (metaTitle && metaTitle !== '数据统计') {
    return metaTitle
  }
  return ''
})

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
  }
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}

.admin-aside {
  background-color: #304156;
  transition: width 0.3s;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  border-bottom: 1px solid #263445;
}

.logo .el-icon {
  color: #409eff;
}

.admin-header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.collapse-icon {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
  transition: color 0.3s;
}

.collapse-icon:hover {
  color: #409eff;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #606266;
}

.user-info:hover {
  color: #409eff;
}

.admin-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}
</style>
