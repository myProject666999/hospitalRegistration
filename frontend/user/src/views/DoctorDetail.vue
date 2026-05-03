<template>
  <div class="doctor-detail-container">
    <el-breadcrumb separator="/" class="breadcrumb">
      <el-breadcrumb-item :to="{ path: '/home' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item :to="{ path: '/doctors' }">医生列表</el-breadcrumb-item>
      <el-breadcrumb-item>医生详情</el-breadcrumb-item>
    </el-breadcrumb>

    <el-card v-loading="loading" class="info-card">
      <div class="doctor-header">
        <el-avatar :size="120" class="doctor-avatar">
          <el-icon :size="60"><UserFilled /></el-icon>
        </el-avatar>
        <div class="doctor-info">
          <h2>{{ doctor.name }}</h2>
          <div class="tags">
            <el-tag type="primary">{{ doctor.title || '主任医师' }}</el-tag>
            <el-tag>{{ doctor.department_name || '内科' }}</el-tag>
            <el-tag type="success">{{ doctor.hospital_name || '医院' }}</el-tag>
          </div>
          <div class="rating">
            <el-rate v-model="doctor.rating" disabled :max="5" text-color="#ff9900" />
            <span>{{ doctor.rating }}分 ({{ doctor.rating_count || 0 }}人评价)</span>
          </div>
          <div class="price">
            <span>挂号费：</span>
            <span class="price-value">¥{{ doctor.price || 50 }}</span>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="section-card">
      <template #header>
        <span>擅长</span>
      </template>
      <p>{{ doctor.skills || '暂无介绍' }}</p>
    </el-card>

    <el-card class="section-card">
      <template #header>
        <span>简介</span>
      </template>
      <p>{{ doctor.introduction || '暂无介绍' }}</p>
    </el-card>

    <el-card class="section-card">
      <template #header>
        <span>可预约排班</span>
      </template>
      <el-table :data="schedules" v-if="schedules.length > 0">
        <el-table-column prop="date" label="日期" />
        <el-table-column prop="time_slot" label="时间段">
          <template #default="scope">
            {{ scope.row.time_slot === 'morning' ? '上午' : scope.row.time_slot === 'afternoon' ? '下午' : '晚间' }}
          </template>
        </el-table-column>
        <el-table-column prop="max_count" label="最大号数" />
        <el-table-column prop="available_count" label="剩余号数" />
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button 
              type="primary" 
              size="small" 
              :disabled="scope.row.available_count <= 0"
              @click="handleAppointment(scope.row)"
            >
              预约
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else description="暂无排班信息" />
    </el-card>

    <el-dialog v-model="appointmentDialogVisible" title="确认预约" width="400px">
      <el-form label-width="80px">
        <el-form-item label="医生">
          <el-input :value="doctor.name" disabled />
        </el-form-item>
        <el-form-item label="日期">
          <el-input :value="selectedSchedule?.date" disabled />
        </el-form-item>
        <el-form-item label="时间段">
          <el-input :value="getTimeSlotText(selectedSchedule?.time_slot)" disabled />
        </el-form-item>
        <el-form-item label="挂号费">
          <el-input :value="'¥' + doctor.price" disabled />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="appointmentForm.remark" type="textarea" placeholder="请输入备注信息（选填）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="appointmentDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="appointmentLoading" @click="submitAppointment">确认预约</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getDoctorDetail, getDoctorSchedules } from '@/api/public'
import { createAppointment } from '@/api/user'
import { useUserStore } from '@/stores/user'
import { UserFilled } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const doctorId = route.params.id
const loading = ref(false)
const doctor = ref({})
const schedules = ref([])

const appointmentDialogVisible = ref(false)
const appointmentLoading = ref(false)
const selectedSchedule = ref(null)
const appointmentForm = reactive({
  remark: ''
})

const getTimeSlotText = (slot) => {
  const map = { morning: '上午', afternoon: '下午', evening: '晚间' }
  return map[slot] || slot
}

const loadDoctorDetail = async () => {
  loading.value = true
  try {
    const res = await getDoctorDetail(doctorId)
    if (res.code === 200) {
      doctor.value = res.data || {}
    }
  } catch (error) {
    console.error('加载医生详情失败', error)
  }
  loading.value = false
}

const loadSchedules = async () => {
  try {
    const res = await getDoctorSchedules(doctorId, { available: true })
    if (res.code === 200) {
      schedules.value = res.data || []
    }
  } catch (error) {
    console.error('加载排班失败', error)
  }
}

const handleAppointment = (schedule) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  selectedSchedule.value = schedule
  appointmentDialogVisible.value = true
}

const submitAppointment = async () => {
  appointmentLoading.value = true
  try {
    const data = {
      doctor_id: doctorId,
      schedule_id: selectedSchedule.value.id,
      date: selectedSchedule.value.date,
      time_slot: selectedSchedule.value.time_slot,
      remark: appointmentForm.remark
    }
    const res = await createAppointment(data)
    if (res.code === 200) {
      ElMessage.success('预约成功')
      appointmentDialogVisible.value = false
      loadSchedules()
    } else {
      ElMessage.error(res.message || '预约失败')
    }
  } catch (error) {
    ElMessage.error('预约失败，请重试')
  }
  appointmentLoading.value = false
}

onMounted(() => {
  loadDoctorDetail()
  loadSchedules()
})
</script>

<style scoped>
.doctor-detail-container {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.breadcrumb {
  margin-bottom: 20px;
}

.info-card {
  margin-bottom: 20px;
}

.doctor-header {
  display: flex;
  gap: 30px;
}

.doctor-avatar {
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.doctor-info {
  flex: 1;
}

.doctor-info h2 {
  font-size: 28px;
  color: #333;
  margin-bottom: 15px;
}

.tags {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.rating {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
}

.rating span {
  color: #909399;
}

.price {
  display: flex;
  align-items: center;
}

.price-value {
  color: #f56c6c;
  font-size: 24px;
  font-weight: bold;
}

.section-card {
  margin-bottom: 20px;
}
</style>
