<template>
  <div class="consultations-container">
    <el-table :data="consultations" stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
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
          <el-button 
            type="success" 
            link 
            v-if="scope.row.status === 2 && !scope.row.is_rated" 
            @click="handleComment(scope.row)"
          >
            评价
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

    <el-dialog v-model="detailVisible" title="就诊详情" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="医生">{{ currentConsultation.doctor_name }}</el-descriptions-item>
        <el-descriptions-item label="科室">{{ currentConsultation.department_name }}</el-descriptions-item>
        <el-descriptions-item label="日期">{{ currentConsultation.date }}</el-descriptions-item>
        <el-descriptions-item label="诊断">
          <span>{{ currentConsultation.diagnosis || '暂无' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="处方">
          <span>{{ currentConsultation.prescription || '暂无' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="医嘱">
          <span>{{ currentConsultation.advice || '暂无' }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <el-dialog v-model="commentVisible" title="评价医生" width="500px">
      <el-form ref="commentFormRef" :model="commentForm" :rules="commentRules" label-width="80px">
        <el-form-item label="评分" prop="rating">
          <el-rate v-model="commentForm.rating" :max="5" show-score text-color="#ff9900" />
        </el-form-item>
        <el-form-item label="评价内容" prop="content">
          <el-input v-model="commentForm.content" type="textarea" placeholder="请输入评价内容" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="commentVisible = false">取消</el-button>
        <el-button type="primary" :loading="commentLoading" @click="submitComment">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getUserConsultations, commentConsultation } from '@/api/user'

const consultations = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const detailVisible = ref(false)
const commentVisible = ref(false)
const commentLoading = ref(false)
const commentFormRef = ref(null)
const currentConsultation = ref({})

const commentForm = reactive({
  rating: 5,
  content: ''
})

const commentRules = {
  rating: [{ required: true, message: '请选择评分', trigger: 'change' }],
  content: [{ required: true, message: '请输入评价内容', trigger: 'blur' }]
}

const getStatusType = (status) => {
  const types = ['info', 'warning', 'success', 'danger']
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = ['待就诊', '就诊中', '已完成', '已取消']
  return texts[status] || '未知'
}

const loadConsultations = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    const res = await getUserConsultations(params)
    if (res.code === 200) {
      consultations.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载就诊列表失败', error)
  }
  loading.value = false
}

const handleSizeChange = () => {
  loadConsultations()
}

const handleCurrentChange = () => {
  loadConsultations()
}

const handleView = (row) => {
  currentConsultation.value = row
  detailVisible.value = true
}

const handleComment = (row) => {
  currentConsultation.value = row
  commentForm.rating = 5
  commentForm.content = ''
  commentVisible.value = true
}

const submitComment = async () => {
  const valid = await commentFormRef.value.validate().catch(() => false)
  if (!valid) return

  commentLoading.value = true
  try {
    const res = await commentConsultation(currentConsultation.value.id, {
      rating: commentForm.rating,
      content: commentForm.content
    })
    if (res.code === 200) {
      ElMessage.success('评价成功')
      commentVisible.value = false
      loadConsultations()
    } else {
      ElMessage.error(res.message || '评价失败')
    }
  } catch (error) {
    ElMessage.error('评价失败')
  }
  commentLoading.value = false
}

onMounted(() => {
  loadConsultations()
})
</script>

<style scoped>
.consultations-container {
  padding: 20px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
