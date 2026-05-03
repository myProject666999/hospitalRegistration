<template>
  <div class="home-container">
    <el-carousel height="300px" class="carousel">
      <el-carousel-item v-for="item in carousels" :key="item.id">
        <div class="carousel-item" :style="{ background: getCarouselBg(item.id) }">
          <div class="carousel-content">
            <h2>{{ item.title || '医院预约挂号系统' }}</h2>
            <p>{{ item.description || '便捷预约，轻松就诊' }}</p>
          </div>
        </div>
      </el-carousel-item>
    </el-carousel>

    <div class="quick-entry">
      <h3>快速入口</h3>
      <el-row :gutter="20">
        <el-col :span="6" v-for="entry in quickEntries" :key="entry.name">
          <router-link :to="entry.path">
            <el-card class="entry-card" shadow="hover">
              <el-icon :size="40" :color="entry.color"><component :is="entry.icon" /></el-icon>
              <span>{{ entry.name }}</span>
            </el-card>
          </router-link>
        </el-col>
      </el-row>
    </div>

    <div class="section">
      <h3>推荐医生</h3>
      <el-row :gutter="20">
        <el-col :span="8" v-for="doctor in doctors.slice(0, 3)" :key="doctor.id">
          <el-card shadow="hover" class="doctor-card" @click="goToDoctor(doctor.id)">
            <div class="doctor-info">
              <el-avatar :size="80" class="doctor-avatar">
                <el-icon :size="40"><User /></el-icon>
              </el-avatar>
              <div class="doctor-details">
                <h4>{{ doctor.name }}</h4>
                <p class="title">{{ doctor.title || '主任医师' }}</p>
                <p class="department">{{ doctor.department_name || '内科' }}</p>
                <div class="rating">
                  <el-rate v-model="doctor.rating" disabled :max="5" text-color="#ff9900" />
                  <span>{{ doctor.rating }}分</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <div class="section">
      <h3>最新公告</h3>
      <el-card>
        <el-table :data="announcements" style="width: 100%">
          <el-table-column prop="title" label="标题" />
          <el-table-column prop="created_at" label="发布时间" width="180">
            <template #default="scope">
              <span>{{ scope.row.created_at ? formatDate(scope.row.created_at) : '' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="scope">
              <el-button type="primary" link @click="goToAnnouncement(scope.row.id)">查看详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCarousels, getAnnouncements, getDoctors } from '@/api/public'
import { 
  User, 
  Calendar, 
  MedicineBox, 
  Document, 
  Connection, 
  Bell 
} from '@element-plus/icons-vue'

const router = useRouter()

const carousels = ref([])
const announcements = ref([])
const doctors = ref([])

const quickEntries = [
  { name: '医生预约', path: '/doctors', icon: 'User', color: '#409eff' },
  { name: '我的预约', path: '/user/appointments', icon: 'Calendar', color: '#67c23a' },
  { name: '药品信息', path: '/medicines', icon: 'MedicineBox', color: '#e6a23c' },
  { name: '公告信息', path: '/announcements', icon: 'Bell', color: '#f56c6c' }
]

const getCarouselBg = (id) => {
  const bgs = [
    'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  ]
  return bgs[id % bgs.length]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString('zh-CN')
}

const goToDoctor = (id) => {
  router.push(`/doctors/${id}`)
}

const goToAnnouncement = (id) => {
  router.push(`/announcements/${id}`)
}

onMounted(async () => {
  try {
    const [carouselRes, announceRes, doctorRes] = await Promise.all([
      getCarousels(),
      getAnnouncements({ page: 1, page_size: 5 }),
      getDoctors({ page: 1, page_size: 3 })
    ])
    if (carouselRes.code === 200) {
      carousels.value = carouselRes.data || [
        { id: 1, title: '欢迎使用医院预约挂号系统', description: '便捷预约，轻松就诊' },
        { id: 2, title: '专家在线问诊', description: '名医专家为您解答疑惑' },
        { id: 3, title: '便捷挂号服务', description: '足不出户，轻松挂号' }
      ]
    }
    if (announceRes.code === 200) {
      announcements.value = announceRes.data?.list || [
        { id: 1, title: '关于2024年春节门诊安排的通知', created_at: '2024-01-20' },
        { id: 2, title: '医院新增儿科夜间门诊', created_at: '2024-01-15' }
      ]
    }
    if (doctorRes.code === 200) {
      doctors.value = doctorRes.data?.list || [
        { id: 1, name: '张医生', title: '主任医师', department_name: '内科', rating: 4.8 },
        { id: 2, name: '李医生', title: '副主任医师', department_name: '外科', rating: 4.9 },
        { id: 3, name: '王医生', title: '主任医师', department_name: '儿科', rating: 4.7 }
      ]
    }
  } catch (error) {
    carousels.value = [
      { id: 1, title: '欢迎使用医院预约挂号系统', description: '便捷预约，轻松就诊' },
      { id: 2, title: '专家在线问诊', description: '名医专家为您解答疑惑' },
      { id: 3, title: '便捷挂号服务', description: '足不出户，轻松挂号' }
    ]
    doctors.value = [
      { id: 1, name: '张医生', title: '主任医师', department_name: '内科', rating: 4.8 },
      { id: 2, name: '李医生', title: '副主任医师', department_name: '外科', rating: 4.9 },
      { id: 3, name: '王医生', title: '主任医师', department_name: '儿科', rating: 4.7 }
    ]
    announcements.value = [
      { id: 1, title: '关于2024年春节门诊安排的通知', created_at: '2024-01-20' },
      { id: 2, title: '医院新增儿科夜间门诊', created_at: '2024-01-15' }
    ]
  }
})
</script>

<style scoped>
.home-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.carousel {
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 30px;
}

.carousel-item {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.carousel-content {
  text-align: center;
  color: white;
}

.carousel-content h2 {
  font-size: 36px;
  margin-bottom: 15px;
}

.carousel-content p {
  font-size: 18px;
  opacity: 0.9;
}

.quick-entry {
  margin-bottom: 30px;
}

.quick-entry h3,
.section h3 {
  font-size: 20px;
  color: #333;
  margin-bottom: 20px;
  padding-left: 10px;
  border-left: 4px solid #409eff;
}

.entry-card {
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.entry-card:hover {
  transform: translateY(-5px);
}

.entry-card .el-icon {
  margin-bottom: 10px;
}

.entry-card span {
  display: block;
  font-size: 14px;
  color: #333;
}

.section {
  margin-bottom: 30px;
}

.doctor-card {
  cursor: pointer;
}

.doctor-info {
  display: flex;
  gap: 20px;
}

.doctor-avatar {
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.doctor-details h4 {
  font-size: 18px;
  color: #333;
  margin-bottom: 5px;
}

.doctor-details .title {
  color: #409eff;
  font-size: 14px;
  margin-bottom: 5px;
}

.doctor-details .department {
  color: #909399;
  font-size: 14px;
  margin-bottom: 10px;
}

.doctor-details .rating {
  display: flex;
  align-items: center;
  gap: 10px;
}

.doctor-details .rating span {
  color: #ff9900;
  font-weight: bold;
}
</style>
