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
    path: '/',
    component: () => import('@/views/layout/Layout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '数据统计', icon: 'DataLine' }
      },
      {
        path: 'appointments',
        name: 'Appointments',
        component: () => import('@/views/appointment/AppointmentList.vue'),
        meta: { title: '预约管理', icon: 'Calendar' }
      },
      {
        path: 'consultations',
        name: 'Consultations',
        component: () => import('@/views/consultation/ConsultationList.vue'),
        meta: { title: '就诊管理', icon: 'Document' }
      },
      {
        path: 'doctors',
        name: 'Doctors',
        component: () => import('@/views/doctor/DoctorList.vue'),
        meta: { title: '医生管理', icon: 'User' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/user/UserList.vue'),
        meta: { title: '用户管理', icon: 'UserFilled' }
      },
      {
        path: 'departments',
        name: 'Departments',
        component: () => import('@/views/department/DepartmentList.vue'),
        meta: { title: '科室管理', icon: 'OfficeBuilding' }
      },
      {
        path: 'hospitals',
        name: 'Hospitals',
        component: () => import('@/views/hospital/HospitalList.vue'),
        meta: { title: '医院管理', icon: 'Building' }
      },
      {
        path: 'medicines',
        name: 'Medicines',
        component: () => import('@/views/medicine/MedicineList.vue'),
        meta: { title: '药品管理', icon: 'MedicineBox' }
      },
      {
        path: 'medical-projects',
        name: 'MedicalProjects',
        component: () => import('@/views/medical-project/MedicalProjectList.vue'),
        meta: { title: '诊疗项目', icon: 'Connection' }
      },
      {
        path: 'schedules',
        name: 'Schedules',
        component: () => import('@/views/schedule/ScheduleList.vue'),
        meta: { title: '排班管理', icon: 'Clock' }
      },
      {
        path: 'carousels',
        name: 'Carousels',
        component: () => import('@/views/carousel/CarouselList.vue'),
        meta: { title: '轮播图管理', icon: 'Picture' }
      },
      {
        path: 'announcements',
        name: 'Announcements',
        component: () => import('@/views/announcement/AnnouncementList.vue'),
        meta: { title: '公告管理', icon: 'Bell' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 医院预约挂号管理后台` : '医院预约挂号管理后台'
  
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
