<template>
  <div class="page-container">
    <div class="page-header">
      <h3>预约管理</h3>
      <div class="filter-bar">
        <el-select v-model="status" placeholder="预约状态" clearable style="width: 150px" @change="handleFilter">
          <el-option :value="0" label="待确认" />
          <el-option :value="1" label="已确认" />
          <el-option :value="2" label="已完成" />
          <el-option :value="3" label="已取消" />
        </el-select>
      </div>
    </div>

    <el-card>
      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_name" label="用户" />
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
            <el-button type="primary" link @click="handleView(scope.row)">查看</el-button>
            <el-button type="success" link v-if="scope.row.status === 0" @click="handleConfirm(scope.row)">确认</el-button>
            <el-button type="danger" link v-if="scope.row.status === 0 || scope.row.status === 1" @click="handleCancel(scope.row)">取消</el-button>
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
    </el-card>

    <el-dialog v-model="detailVisible" title="预约详情" width="500px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="预约号">{{ currentItem.id }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentItem.user_name }}</el-descriptions-item>
        <el-descriptions-item label="医生">{{ currentItem.doctor_name }}</el-descriptions-item>
        <el-descriptions-item label="科室">{{ currentItem.department_name }}</el-descriptions-item>
        <el-descriptions-item label="日期">{{ currentItem.date }}</el-descriptions-item>
        <el-descriptions-item label="时间段">
          {{ currentItem.time_slot === 'morning' ? '上午' : currentItem.time_slot === 'afternoon' ? '下午' : '晚间' }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentItem.status)">
            {{ getStatusText(currentItem.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="备注">{{ currentItem.remark || '无' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const list = ref([])
const status = ref(null)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const detailVisible = ref(false)
const currentItem = ref({})

const getStatusType = (status) => {
  const types = ['warning', 'primary', 'success', 'danger']
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = ['待确认', '已确认', '已完成', '已取消']
  return texts[status] || '未知'
}

const loadList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (status.value !== null) {
      params.status = status.value
    }
    const res = await request({
      url: '/admin/appointments',
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

const handleFilter = () => {
  currentPage.value = 1
  loadList()
}

const handleSizeChange = () => {
  loadList()
}

const handleCurrentChange = () => {
  loadList()
}

const handleView = (row) => {
  currentItem.value = row
  detailVisible.value = true
}

const handleConfirm = async (row) => {
  try {
    await ElMessageBox.confirm('确定要确认该预约吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await request({
      url: `/admin/appointments/${row.id}/confirm`,
      method: 'put'
    })
    if (res.code === 200) {
      ElMessage.success('确认成功')
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

const handleCancel = async (row) => {
  try {
    await ElMessageBox.confirm('确定要取消该预约吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await request({
      url: `/admin/appointments/${row.id}/cancel`,
      method: 'put'
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

onMounted(() => {
  loadList()
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

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
