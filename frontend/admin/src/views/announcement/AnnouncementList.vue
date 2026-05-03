<template>
  <div class="page-container">
    <div class="page-header">
      <h3>公告管理</h3>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新增公告
      </el-button>
    </div>

    <el-card>
      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="scope">
            <div class="title-cell">
              <el-tag v-if="scope.row.is_top === 1" type="danger" size="small" style="margin-right: 8px">置顶</el-tag>
              <span>{{ scope.row.title }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="scope">
            <el-tag :type="getTypeTag(scope.row.type)">
              {{ getTypeText(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览量" width="100" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'warning'">
              {{ scope.row.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑公告' : '新增公告'" width="800px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入公告标题" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="类型" prop="type">
              <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%">
                <el-option label="公告" :value="1" />
                <el-option label="新闻" :value="2" />
                <el-option label="通知" :value="3" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="已发布" :value="1" />
                <el-option label="草稿" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="是否置顶" prop="is_top">
          <el-switch v-model="form.is_top" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="摘要" prop="summary">
          <el-input v-model="form.summary" type="textarea" placeholder="请输入摘要（可选）" :rows="2" />
        </el-form-item>
        <el-form-item label="封面图片" prop="cover">
          <el-input v-model="form.cover" placeholder="请输入封面图片URL（可选）" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="form.content" type="textarea" placeholder="请输入公告内容" :rows="10" />
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
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const dialogVisible = ref(false)
const submitLoading = ref(false)
const formRef = ref(null)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({
  title: '',
  content: '',
  summary: '',
  cover: '',
  type: 1,
  is_top: 0,
  status: 1,
  sort: 0
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const getTypeText = (type) => {
  const map = { 1: '公告', 2: '新闻', 3: '通知' }
  return map[type] || '公告'
}

const getTypeTag = (type) => {
  const map = { 1: 'primary', 2: 'success', 3: 'warning' }
  return map[type] || 'primary'
}

const formatTime = (time) => {
  if (!time) return '-'
  const d = new Date(time)
  return d.toLocaleString('zh-CN')
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/admin/announcements',
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

const resetForm = () => {
  Object.assign(form, {
    title: '',
    content: '',
    summary: '',
    cover: '',
    type: 1,
    is_top: 0,
    status: 1,
    sort: 0
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
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该公告吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await request({
      url: `/admin/announcements/${row.id}`,
      method: 'delete'
    })
    if (res.code === 200) {
      ElMessage.success('删除成功')
      loadList()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const res = await request({
      url: isEdit.value ? `/admin/announcements/${currentId.value}` : '/admin/announcements',
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

.title-cell {
  display: flex;
  align-items: center;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
