<template>
  <div class="doctors-container">
    <div class="search-bar">
      <el-select v-model="searchForm.department_id" placeholder="选择科室" clearable style="width: 200px">
        <el-option
          v-for="dept in departments"
          :key="dept.id"
          :label="dept.name"
          :value="dept.id"
        />
      </el-select>
      <el-input v-model="searchForm.keyword" placeholder="搜索医生姓名" clearable style="width: 250px; margin-left: 15px" @keyup.enter="handleSearch">
        <template #append>
          <el-button @click="handleSearch"><el-icon><Search /></el-icon></el-button>
        </template>
      </el-input>
    </div>

    <div class="doctor-list">
      <el-row :gutter="20">
        <el-col :span="8" v-for="doctor in doctors" :key="doctor.id">
          <el-card shadow="hover" class="doctor-card" @click="goToDetail(doctor.id)">
            <div class="doctor-info">
              <el-avatar :size="100" class="doctor-avatar">
                <el-icon :size="50"><UserFilled /></el-icon>
              </el-avatar>
              <div class="doctor-details">
                <h4 class="doctor-name">{{ doctor.name }}</h4>
                <p class="doctor-title">{{ doctor.title || '主任医师' }}</p>
                <p class="doctor-department">
                  <el-icon><OfficeBuilding /></el-icon>
                  {{ doctor.department_name || '内科' }}
                </p>
                <p class="doctor-hospital">
                  <el-icon><Location /></el-icon>
                  {{ doctor.hospital_name || '医院' }}
                </p>
                <div class="rating-section">
                  <el-rate v-model="doctor.rating" disabled :max="5" text-color="#ff9900" />
                  <span class="rating-text">{{ doctor.rating }}分</span>
                  <span class="rating-count">{{ doctor.rating_count || 0 }}人评价</span>
                </div>
              </div>
            </div>
            <div class="doctor-skills">
              <span>擅长：</span>
              <span class="skills-text">{{ doctor.skills || '暂无介绍' }}</span>
            </div>
            <div class="card-footer">
              <span class="price">挂号费：¥{{ doctor.price || 50 }}</span>
              <el-button type="primary" size="small" @click.stop="goToDetail(doctor.id)">
                立即预约
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[9, 18, 27]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getDoctors, getDepartments } from '@/api/public'
import { 
  UserFilled, 
  Search, 
  OfficeBuilding, 
  Location 
} from '@element-plus/icons-vue'

const router = useRouter()

const doctors = ref([])
const departments = ref([])
const currentPage = ref(1)
const pageSize = ref(9)
const total = ref(0)
const loading = ref(false)

const searchForm = reactive({
  department_id: null,
  keyword: ''
})

const goToDetail = (id) => {
  router.push(`/doctors/${id}`)
}

const loadDepartments = async () => {
  try {
    const res = await getDepartments()
    if (res.code === 200) {
      departments.value = res.data || []
    }
  } catch (error) {
    console.error('加载科室失败', error)
  }
}

const loadDoctors = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchForm
    }
    const res = await getDoctors(params)
    if (res.code === 200) {
      doctors.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载医生列表失败', error)
    doctors.value = []
  }
  loading.value = false
}

const handleSearch = () => {
  currentPage.value = 1
  loadDoctors()
}

const handleSizeChange = () => {
  loadDoctors()
}

const handleCurrentChange = () => {
  loadDoctors()
}

onMounted(() => {
  loadDepartments()
  loadDoctors()
})
</script>

<style scoped>
.doctors-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.search-bar {
  margin-bottom: 20px;
}

.doctor-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.doctor-card:hover {
  transform: translateY(-5px);
}

.doctor-info {
  display: flex;
  gap: 20px;
  margin-bottom: 15px;
}

.doctor-avatar {
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.doctor-details {
  flex: 1;
}

.doctor-name {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.doctor-title {
  color: #409eff;
  font-size: 14px;
  margin-bottom: 5px;
}

.doctor-department,
.doctor-hospital {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #606266;
  font-size: 13px;
  margin-bottom: 3px;
}

.rating-section {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 8px;
}

.rating-text {
  color: #ff9900;
  font-weight: bold;
}

.rating-count {
  color: #909399;
  font-size: 12px;
}

.doctor-skills {
  border-top: 1px solid #ebeef5;
  padding-top: 15px;
  margin-bottom: 15px;
  font-size: 13px;
  color: #606266;
}

.skills-text {
  color: #333;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  color: #f56c6c;
  font-weight: bold;
  font-size: 16px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
