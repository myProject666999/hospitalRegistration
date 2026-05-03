<template>
  <div class="page-container">
    <div class="page-header">
      <h3>排班管理</h3>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新增排班
      </el-button>
    </div>

    <el-card>
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="医生">
          <el-select v-model="searchForm.doctor_id" placeholder="请选择医生" clearable @change="loadList">
            <el-option v-for="doctor in doctors" :key="doctor.id" :label="doctor.real_name || doctor.name" :value="doctor.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="searchForm.date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" @change="loadList" />
        </el-form-item>
      </el-form>

      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="doctor_name" label="医生" width="120">
          <template #default="scope">
            {{ scope.row.doctor?.real_name || scope.row.doctor?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="doctor_department" label="科室" width="100">
          <template #default="scope">
            {{ scope.row.doctor?.department?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="date" label="日期" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.date) }}
          </template>
        </el-table-column>
        <el-table-column prop="period" label="时段" width="80">
          <template #default="scope">
            <el-tag :type="getPeriodType(scope.row.period)">
              {{ getPeriodText(scope.row.period) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="time" label="时间" width="150">
          <template #default="scope">
            {{ scope.row.start_time || '-' }} - {{ scope.row.end_time || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="count" label="号源" width="120">
          <template #default="scope">
            <span class="count-info">
              剩余: <span class="remain">{{ scope.row.remain_count }}</span> / 总计: {{ scope.row.total_count }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="地点" width="120">
          <template #default="scope">
            {{ scope.row.location || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)" :disabled="scope.row.status === -1">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)" :disabled="scope.row.status === -1">取消</el-button>
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
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑排班' : '新增排班'" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="医生" prop="doctor_id">
          <el-select v-model="form.doctor_id" placeholder="请选择医生" style="width: 100%" :disabled="isEdit">
            <el-option v-for="doctor in doctors" :key="doctor.id" :label="doctor.real_name || doctor.name" :value="doctor.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期" prop="date">
          <el-date-picker v-model="form.date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" :disabled="isEdit" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="时段" prop="period">
              <el-select v-model="form.period" placeholder="请选择时段" style="width: 100%" :disabled="isEdit">
                <el-option label="上午" :value="1" />
                <el-option label="下午" :value="2" />
                <el-option label="晚上" :value="3" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="号源数" prop="total_count">
              <el-input-number v-model="form.total_count" :min="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始时间" prop="start_time">
              <el-time-picker v-model="form.start_time" placeholder="选择开始时间" value-format="HH:mm" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束时间" prop="end_time">
              <el-time-picker v-model="form.end_time" placeholder="选择结束时间" value-format="HH:mm" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="出诊地点" prop="location">
          <el-input v-model="form.location" placeholder="请输入出诊地点" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import { Plus } from '@element-plus/icons-vue'

const list = ref([])
const doctors = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const dialogVisible = ref(false)
const submitLoading = ref(false)
const formRef = ref(null)
const isEdit = ref(false)
const currentId = ref(null)

const searchForm = reactive({
  doctor_id: null,
  date: ''
})

const form = reactive({
  doctor_id: null,
  date: '',
  period: 1,
  start_time: '',
  end_time: '',
  total_count: 20,
  location: ''
})

const rules = {
  doctor_id: [{ required: true, message: '请选择医生', trigger: 'change' }],
  date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  period: [{ required: true, message: '请选择时段', trigger: 'change' }],
  total_count: [{ required: true, message: '请输入号源数', trigger: 'blur' }]
}

const formatDate = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleDateString('zh-CN')
}

const getPeriodText = (period) => {
  const map = { 1: '上午', 2: '下午', 3: '晚上' }
  return map[period] || '-'
}

const getPeriodType = (period) => {
  const map = { 1: 'primary', 2: 'success', 3: 'warning' }
  return map[period] || ''
}

const getStatusText = (status) => {
  const map = { 1: '可预约', 0: '已约满', '-1': '已取消' }
  return map[status] || '-'
}

const getStatusType = (status) => {
  const map = { 1: 'success', 0: 'warning', '-1': 'danger' }
  return map[status] || 'info'
}

const loadList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (searchForm.doctor_id) {
      params.doctor_id = searchForm.doctor_id
    }
    if (searchForm.date) {
      params.date = searchForm.date
    }
    const res = await request({
      url: '/admin/schedules',
      method: 'get',
      params
    })
    if (res.code === 200) {
      list.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载列表失败', error)
  }
  loading.value = false
}

const loadDoctors = async () => {
  try {
    const res = await request({
      url: '/doctors',
      method: 'get'
    })
    if (res.code === 200) {
      doctors.value = res.data?.list || res.data || []
    }
  } catch (error) {
    console.error('加载医生列表失败', error)
  }
}

const resetForm = () => {
  Object.assign(form, {
    doctor_id: null,
    date: '',
    period: 1,
    start_time: '',
    end_time: '',
    total_count: 20,
    location: ''
  })
}

const handleAdd = () => {
  isEdit.value = false
  currentId.value = null
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  currentId.value = row.id
  Object.assign(form, {
    doctor_id: row.doctor_id,
    date: formatDate(row.date),
    period: row.period,
    start_time: row.start_time,
    end_time: row.end_time,
    total_count: row.total_count,
    location: row.location
  })
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要取消该排班吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await request({
      url: `/admin/schedules/${row.id}`,
      method: 'delete'
    })
    if (res.code === 200) {
      ElMessage.success('取消成功')
      loadList()
    } else {
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const res = await request({
      url: isEdit.value ? `/admin/schedules/${currentId.value}` : '/admin/schedules',
      method: isEdit.value ? 'put' : 'post',
      data: form
    })
    if (res.code === 200) {
      ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
      dialogVisible.value = false
      loadList()
    } else {
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
  submitLoading.value = false
}

onMounted(() => {
  loadList()
  loadDoctors()
})
</script>

<style scoped>
.page-container {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h3 {
  margin: 0;
  font-size: 18px;
}

.search-form {
  margin-bottom: 20px;
}

.count-info {
  font-size: 13px;
}

.count-info .remain {
  color: #67c23a;
  font-weight: bold;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
