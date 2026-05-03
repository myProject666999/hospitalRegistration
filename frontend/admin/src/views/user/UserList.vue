<template>
  <div class="page-container">
    <div class="page-header">
      <h3>用户管理</h3>
      <div class="filter-bar">
        <el-input v-model="keyword" placeholder="搜索用户名" clearable style="width: 200px" @keyup.enter="loadList">
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" style="margin-left: 10px" @click="loadList">搜索</el-button>
      </div>
    </div>

    <el-card>
      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="gender" label="性别">
          <template #default="scope">
            <span>{{ scope.row.gender === 1 ? '男' : scope.row.gender === 2 ? '女' : '未知' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="scope">
            <span>{{ scope.row.created_at ? formatDate(scope.row.created_at) : '' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="danger" link @click="handleToggleStatus(scope.row)">
              {{ scope.row.status === 1 ? '禁用' : '启用' }}
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
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="编辑用户" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" disabled />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="form.gender">
            <el-radio :value="1">男</el-radio>
            <el-radio :value="2">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
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
import { Search } from '@element-plus/icons-vue'

const list = ref([])
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const dialogVisible = ref(false)
const submitLoading = ref(false)
const formRef = ref(null)
const currentId = ref(null)

const form = reactive({
  username: '',
  name: '',
  phone: '',
  gender: 1,
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleString('zh-CN')
}

const loadList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: keyword.value
    }
    const res = await request({
      url: '/admin/users',
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

const handleEdit = (row) => {
  currentId.value = row.id
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleToggleStatus = async (row) => {
  const action = row.status === 1 ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await request({
      url: `/admin/users/${row.id}/status`,
      method: 'put',
      data: { status: row.status === 1 ? 0 : 1 }
    })
    if (res.code === 200) {
      ElMessage.success(`${action}成功`)
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
      url: `/admin/users/${currentId.value}`,
      method: 'put',
      data: { ...form }
    })
    if (res.code === 200) {
      ElMessage.success('编辑成功')
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

.filter-bar {
  display: flex;
  align-items: center;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
