<template>
  <div class="page-container">
    <div class="page-header">
      <h3>医生管理</h3>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新增医生
      </el-button>
    </div>

    <el-card>
      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="title" label="职称" />
        <el-table-column prop="department_name" label="科室" />
        <el-table-column prop="hospital_name" label="医院" />
        <el-table-column prop="rating" label="评分">
          <template #default="scope">
            <el-rate v-model="scope.row.rating" disabled :max="5" />
          </template>
        </el-table-column>
        <el-table-column prop="price" label="挂号费">
          <template #default="scope">
            <span class="price">¥{{ scope.row.price || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑医生' : '新增医生'" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="职称" prop="title">
          <el-input v-model="form.title" placeholder="请输入职称" />
        </el-form-item>
        <el-form-item label="科室" prop="department_id">
          <el-select v-model="form.department_id" placeholder="请选择科室" style="width: 100%">
            <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="医院" prop="hospital_id">
          <el-select v-model="form.hospital_id" placeholder="请选择医院" style="width: 100%">
            <el-option v-for="h in hospitals" :key="h.id" :label="h.name" :value="h.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="挂号费" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="擅长" prop="skills">
          <el-input v-model="form.skills" type="textarea" placeholder="请输入擅长领域" :rows="3" />
        </el-form-item>
        <el-form-item label="简介" prop="introduction">
          <el-input v-model="form.introduction" type="textarea" placeholder="请输入简介" :rows="3" />
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
const departments = ref([])
const hospitals = ref([])
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
  username: '',
  password: '',
  name: '',
  title: '',
  department_id: null,
  hospital_id: null,
  price: 0,
  skills: '',
  introduction: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/admin/doctors',
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

const loadDepartments = async () => {
  try {
    const res = await request({
      url: '/departments',
      method: 'get'
    })
    if (res.code === 200) {
      departments.value = res.data || []
    }
  } catch (error) {
    console.error('加载科室失败', error)
  }
}

const loadHospitals = async () => {
  try {
    const res = await request({
      url: '/hospitals',
      method: 'get'
    })
    if (res.code === 200) {
      hospitals.value = res.data || []
    }
  } catch (error) {
    console.error('加载医院失败', error)
  }
}

const resetForm = () => {
  Object.assign(form, {
    username: '',
    password: '',
    name: '',
    title: '',
    department_id: null,
    hospital_id: null,
    price: 0,
    skills: '',
    introduction: ''
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
    await ElMessageBox.confirm('确定要删除该医生吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await request({
      url: `/admin/doctors/${row.id}`,
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
    const data = { ...form }
    if (!isEdit.value) {
      delete data.password
    }
    const res = await request({
      url: isEdit.value ? `/admin/doctors/${currentId.value}` : '/admin/doctors',
      method: isEdit.value ? 'put' : 'post',
      data
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
  loadDepartments()
  loadHospitals()
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

.price {
  color: #f56c6c;
  font-weight: bold;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
