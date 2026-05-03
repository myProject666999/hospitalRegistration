<template>
  <div class="appointments-container">
    <div class="filter-bar">
      <el-select v-model="status" placeholder="预约状态" clearable @change="handleFilter">
        <el-option :value="0" label="待确认" />
        <el-option :value="1" label="已确认" />
        <el-option :value="2" label="已完成" />
        <el-option :value="3" label="已取消" />
      </el-select>
    </div>

    <el-table :data="appointments" stripe v-loading="loading">
      <el-table-column prop="id" label="预约号" width="100" />
      <el-table-column prop="doctor_name" label="医生" />
      <el-table-column prop="department_name" label="科室" />
      <el-table-column prop="date" label="日期" />
      <el-table-column prop="time_slot" label="时间段">
        <template #default="scope">
          {{ scope.row.time_slot === 'morning' ? '上午' : scope.row.time_slot === 'afternoon' ? '下午' : '晚间' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-tag :type="getStatusType(scope.row.status)">
            {{ getStatusText(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button type="primary" link v-if="scope.row.status === 0 || scope.row.status === 1" @click="handleView(scope.row)">
            查看
          </el-button>
          <el-button type="danger" link v-if="scope.row.status === 0" @click="handleCancel(scope.row)">
            取消预约
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <el-dialog v-model="detailVisible" title="预约详情" width="500px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="预约号">{{ currentAppointment.id }}</el-descriptions-item>
        <el-descriptions-item label="医生">{{ currentAppointment.doctor_name }}</el-descriptions-item>
        <el-descriptions-item label="科室">{{ currentAppointment.department_name }}</el-descriptions-item>
        <el-descriptions-item label="日期">{{ currentAppointment.date }}</el-descriptions-item>
        <el-descriptions-item label="时间段">
          {{ currentAppointment.time_slot === 'morning' ? '上午' : currentAppointment.time_slot === 'afternoon' ? '下午' : '晚间' }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentAppointment.status)">
            {{ getStatusText(currentAppointment.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="备注">{{ currentAppointment.remark || '无' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserAppointments, cancelAppointment } from '@/api/user'

const appointments = ref([])
const status = ref(null)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const detailVisible = ref(false)
const currentAppointment = ref({})

const getStatusType = (status) => {
  const types = ['warning', 'primary', 'success', 'danger']
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = ['待确认', '已确认', '已完成', '已取消']
  return texts[status] || '未知'
}

const loadAppointments = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (status.value !== null) {
      params.status = status.value
    }
    const res = await getUserAppointments(params)
    if (res.code === 200) {
      appointments.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载预约列表失败', error)
  }
  loading.value = false
}

const handleFilter = () => {
  currentPage.value = 1
  loadAppointments()
}

const handleSizeChange = () => {
  loadAppointments()
}

const handleCurrentChange = () => {
  loadAppointments()
}

const handleView = (row) => {
  currentAppointment.value = row
  detailVisible.value = true
}

const handleCancel = async (row) => {
  try {
    await ElMessageBox.confirm('确定要取消该预约吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await cancelAppointment(row.id)
    if (res.code === 200) {
      ElMessage.success('取消成功')
      loadAppointments()
    } else {
      ElMessage.error(res.message || '取消失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('取消失败')
    }
  }
}

onMounted(() => {
  loadAppointments()
})
</script>

<style scoped>
.appointments-container {
  padding: 20px;
}

.filter-bar {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
