import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/',
    component: () => import('@/views/layout/UserLayout.vue'),
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'doctors',
        name: 'Doctors',
        component: () => import('@/views/Doctors.vue'),
        meta: { title: '医生列表' }
      },
      {
        path: 'doctors/:id',
        name: 'DoctorDetail',
        component: () => import('@/views/DoctorDetail.vue'),
        meta: { title: '医生详情' }
      },
      {
        path: 'medicines',
        name: 'Medicines',
        component: () => import('@/views/Medicines.vue'),
        meta: { title: '药品信息' }
      },
      {
        path: 'medicines/:id',
        name: 'MedicineDetail',
        component: () => import('@/views/MedicineDetail.vue'),
        meta: { title: '药品详情' }
      },
      {
        path: 'medical-projects',
        name: 'MedicalProjects',
        component: () => import('@/views/MedicalProjects.vue'),
        meta: { title: '诊疗项目' }
      },
      {
        path: 'announcements',
        name: 'Announcements',
        component: () => import('@/views/Announcements.vue'),
        meta: { title: '公告信息' }
      },
      {
        path: 'announcements/:id',
        name: 'AnnouncementDetail',
        component: () => import('@/views/AnnouncementDetail.vue'),
        meta: { title: '公告详情' }
      },
      {
        path: 'user/profile',
        name: 'UserProfile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人信息', requiresAuth: true }
      },
      {
        path: 'user/password',
        name: 'UserPassword',
        component: () => import('@/views/user/Password.vue'),
        meta: { title: '修改密码', requiresAuth: true }
      },
      {
        path: 'user/appointments',
        name: 'UserAppointments',
        component: () => import('@/views/user/Appointments.vue'),
        meta: { title: '预约管理', requiresAuth: true }
      },
      {
        path: 'user/consultations',
        name: 'UserConsultations',
        component: () => import('@/views/user/Consultations.vue'),
        meta: { title: '就诊记录', requiresAuth: true }
      },
      {
        path: 'user/favorites',
        name: 'UserFavorites',
        component: () => import('@/views/user/Favorites.vue'),
        meta: { title: '我的收藏', requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 医院预约挂号系统` : '医院预约挂号系统'
  
  const userStore = useUserStore()
  
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  
  if (requiresAuth && !userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    next('/login')
  } else {
    next()
  }
})

export default router
