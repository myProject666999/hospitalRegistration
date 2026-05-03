<template>
  <div class="page-container">
    <div class="page-header">
      <h3>就诊管理</h3>
    </div>

    <el-card>
      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_name" label="用户" />
        <el-table-column prop="doctor_name" label="医生" />
        <el-table-column prop="department_name" label="科室" />
        <el-table-column prop="date" label="日期" />
        <el-table-column prop="diagnosis" label="诊断">
          <template #default="scope">
            <span>{{ scope.row.diagnosis || '暂无' }}</span>
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

    <el-dialog v-model="detailVisible" title="就诊详情" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="用户">{{ currentItem.user_name }}</el-descriptions-item>
        <el-descriptions-item label="医生">{{ currentItem.doctor_name }}</el-descriptions-item>
        <el-descriptions-item label="科室">{{ currentItem.department_name }}</el-descriptions-item>
        <el-descriptions-item label="日期">{{ currentItem.date }}</el-descriptions-item>
        <el-descriptions-item label="诊断">{{ currentItem.diagnosis || '暂无' }}</el-descriptions-item>
        <el-descriptions-item label="处方">{{ currentItem.prescription || '暂无' }}</el-descriptions-item>
        <el-descriptions-item label="医嘱">{{ currentItem.advice || '暂无' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentItem.status)">
            {{ getStatusText(currentItem.status) }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const list = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const detailVisible = ref(false)
const currentItem = ref({})

const getStatusType = (status) => {
  const types = ['info', 'warning', 'success', 'danger']
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = ['待就诊', '就诊中', '已完成', '已取消']
  return texts[status] || '未知'
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/admin/consultations',
      method: 'get',
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
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

const handleView = (row) => {
  currentItem.value = row
  detailVisible.value = true
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
